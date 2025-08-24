package main

import (
    "context"
    "database/sql"
    _ "github.com/lib/pq" 
    "fmt"
    "os"
    "strings"
    "time"

	"github.com/google/uuid"
)

type DB struct {
    Conn *sql.DB
}

func NewDB() (*DB, error) {
    host := os.Getenv("DB_HOST")
    user := os.Getenv("DB_USER")
    password := os.Getenv("DB_PASSWORD")
    dbname := os.Getenv("DB_NAME")

    dsn := fmt.Sprintf(
        "host=%s user=%s password=%s dbname=%s sslmode=disable",
        host, user, password, dbname,
    )

    conn, err := sql.Open("postgres", dsn)
    if err != nil {
        return nil, err
    }

    // Retry loop in case of "too many clients" or transient errors
    maxRetries := 10
    for i := 0; i < maxRetries; i++ {
        err = conn.Ping()
        if err == nil {
            return &DB{Conn: conn}, nil
        }

        // Detect "too many clients" explicitly (Postgres code: 53300)
        if strings.Contains(err.Error(), "too many clients") {
            wait := time.Duration(i+1) * time.Second // simple backoff
            fmt.Printf("Database has too many clients, retrying in %v...\n", wait)
            time.Sleep(wait)
            maxRetries++
            continue
        }

        // Other errors -> return immediately
        return nil, err
    }

    return nil, fmt.Errorf("could not connect to database after retries: %w", err)
}

func (db *DB) CreateCompIfNotExist(ctx context.Context, name string, id int) error {
    // Upsert with "ON CONFLICT DO NOTHING" if supported (Postgres)
    query := `
        INSERT INTO competition (id, name)
        VALUES ($1, $2)
        ON CONFLICT (id) DO NOTHING
    `
    res, err := db.Conn.ExecContext(ctx, query, id, name)
    if err != nil {
        return fmt.Errorf("insert competition failed: %w", err)
    }

    rowsAffected, err := res.RowsAffected()
    if err != nil {
        return fmt.Errorf("checking rows affected failed: %w", err)
    }

    if rowsAffected == 0 {
        return nil
    }
    return nil
}

func (db *DB) CreateSeasonIfNotExist(ctx context.Context, competitionID int, year string) (uuid.UUID, error) {
    var seasonID uuid.UUID

    query := `
        INSERT INTO season (competition_id, year)
        VALUES ($1, $2)
        ON CONFLICT (competition_id, year) DO NOTHING
        RETURNING id
    `

    err := db.Conn.QueryRowContext(ctx, query, competitionID, year).Scan(&seasonID)
    if err == sql.ErrNoRows {
        // Season already exists, fetch the existing id
        selectQuery := `
            SELECT id FROM season
            WHERE competition_id = $1 AND year = $2
        `
        err = db.Conn.QueryRowContext(ctx, selectQuery, competitionID, year).Scan(&seasonID)
        if err != nil {
            return uuid.Nil, fmt.Errorf("failed to fetch existing season id: %w", err)
        }
        return seasonID, nil
    }
    if err != nil {
        return uuid.Nil, fmt.Errorf("failed to insert season: %w", err)
    }

    return seasonID, nil
}

func (db *DB) CreateRound(ctx context.Context, roundIndex int, roundName string, seasonID uuid.UUID) (uuid.UUID, bool, error) {
    var (
        roundID      uuid.UUID
        datesAreSet  bool
    )

    err := db.Conn.QueryRowContext(ctx, `
        INSERT INTO round (id, season_id, round_index, round_name)
        VALUES ($1, $2, $3, $4)
        ON CONFLICT (season_id, round_index) 
        DO UPDATE SET round_name = EXCLUDED.round_name
        RETURNING id, (start_day IS NOT NULL AND end_day IS NOT NULL) AS dates_set
    `,
        uuid.New(), seasonID, roundIndex, roundName,
    ).Scan(&roundID, &datesAreSet)

    if err != nil {
        return uuid.Nil, false, err
    }
    return roundID, datesAreSet, nil
}

func (db *DB) SetRoundDates(ctx context.Context, roundID uuid.UUID, startDay, endDay string) error {
    _, err := db.Conn.ExecContext(ctx, `
        UPDATE round
        SET start_day = $1, end_day = $2
        WHERE id = $3
    `, startDay, endDay, roundID)

    return err
}

func (db *DB) CreateMatch(ctx context.Context, roundID uuid.UUID, homeTeam, awayTeam string) (uuid.UUID, error) {
    var matchID uuid.UUID

    err := db.Conn.QueryRowContext(ctx, `
        INSERT INTO match (id, round_id, home_team, away_team)
        VALUES ($1, $2, $3, $4)
        ON CONFLICT (round_id, home_team, away_team)
        DO UPDATE SET home_team = EXCLUDED.home_team
        RETURNING id
    `,
        uuid.New(), roundID, homeTeam, awayTeam,
    ).Scan(&matchID)

    if err != nil {
        return uuid.Nil, err
    }

    return matchID, nil
}

func (db *DB) CreatePlay(ctx context.Context, matchID uuid.UUID, playIndex int, timeStr, playText, team, notes string) (uuid.UUID, error) {
    var playID uuid.UUID

    err := db.Conn.QueryRowContext(ctx, `
        INSERT INTO play_by_play (id, match_id, play_index, time, play, team, notes)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        ON CONFLICT (match_id, play_index) 
        DO UPDATE 
        SET time = EXCLUDED.time,
            play = EXCLUDED.play,
            team = EXCLUDED.team,
            notes = EXCLUDED.notes
        RETURNING id
    `,
        uuid.New(), matchID, playIndex, timeStr, playText, team, notes,
    ).Scan(&playID)

    if err != nil {
        return uuid.Nil, err
    }
    return playID, nil
}

func (db *DB) SetTeamLists(ctx context.Context, matchID uuid.UUID, homeTeamList, awayTeamList []*Player) error {
	tx, err := db.Conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

    // helper to link player to match
    insertMatchPlayer := func(playerID, team string) error {
        _, err := tx.ExecContext(ctx,
            `INSERT INTO match_player (match_id, player_id, team)
             VALUES ($1, $2, $3)
             ON CONFLICT DO NOTHING`, // so we don't double insert
            matchID, playerID, team,
        )
        return err
    }

    for _, p := range homeTeamList {
        pid, err := db.InsertPlayer(ctx, matchID, p.nameFirst, p.nameLast, p.position, p.number)
        if err != nil {
            return err
        }
        if err := insertMatchPlayer(pid, "home"); err != nil {
            return err
        }
    }

    for _, p := range awayTeamList {
        pid, err := db.InsertPlayer(ctx, matchID, p.nameFirst, p.nameLast, p.position, p.number)
        if err != nil {
            return err
        }
        if err := insertMatchPlayer(pid, "away"); err != nil {
            return err
        }
    }

    return tx.Commit()
}

func (db *DB) InsertPlayer(ctx context.Context, matchID uuid.UUID, first, last, position string, number int) (string, error) {
    query := `
        INSERT INTO player (match_id, name_first, name_last, position, number)
        VALUES ($1, $2, $3, $4, $5)
        ON CONFLICT (match_id, name_first, name_last) DO NOTHING
        RETURNING id;
    `

    var id string
    err := db.Conn.QueryRowContext(ctx, query, matchID, first, last, position, number).Scan(&id)
    if err == sql.ErrNoRows {
        // Means conflict happened and DO NOTHING triggered
        return "", nil
    }
    if err != nil {
        return "", fmt.Errorf("insert player failed: %w", err)
    }

    return id, nil
}

func (db *DB) setScore(ctx context.Context, matchID uuid.UUID, column string, score int) error {
    query := fmt.Sprintf(`UPDATE match SET %s = $1 WHERE id = $2;`, column)

    res, err := db.Conn.ExecContext(ctx, query, score, matchID)
    if err != nil {
        return fmt.Errorf("failed to update %s: %w", column, err)
    }

    rowsAffected, err := res.RowsAffected()
    if err != nil {
        return fmt.Errorf("failed to check rows affected for %s: %w", column, err)
    }

    if rowsAffected == 0 {
        return fmt.Errorf("no match found with id %s", matchID)
    }

    return nil
}

func (db *DB) SetHomeScore(ctx context.Context, matchID uuid.UUID, score int) error {
    return db.setScore(ctx, matchID, "home_score", score)
}

func (db *DB) SetAwayScore(ctx context.Context, matchID uuid.UUID, score int) error {
    return db.setScore(ctx, matchID, "away_score", score)
}

func (db *DB) SetLocation(ctx context.Context, matchID uuid.UUID, location string) error {
    query := `
        UPDATE match
        SET location = $1
        WHERE id = $2;
    `

    res, err := db.Conn.ExecContext(ctx, query, location, matchID)
    if err != nil {
        return fmt.Errorf("failed to update location: %w", err)
    }

    rowsAffected, err := res.RowsAffected()
    if err != nil {
        return fmt.Errorf("failed to check rows affected: %w", err)
    }

    if rowsAffected == 0 {
        return fmt.Errorf("no match found with id %s", matchID)
    }

    return nil
}

func (db *DB) SetDatePlayed(ctx context.Context, matchID uuid.UUID, dateStr string) error {
    query := `
        UPDATE match
        SET date_played = $1
        WHERE id = $2;
    `

    res, err := db.Conn.ExecContext(ctx, query, dateStr, matchID)
    if err != nil {
        return fmt.Errorf("failed to update date_played: %w", err)
    }

    rowsAffected, err := res.RowsAffected()
    if err != nil {
        return fmt.Errorf("failed to check rows affected: %w", err)
    }

    if rowsAffected == 0 {
        return fmt.Errorf("no match found with id %s", matchID)
    }

    return nil
}

func (db *DB) SetWeather(ctx context.Context, matchID uuid.UUID, weather string) error {
    query := `
        UPDATE match
        SET weather = $1
        WHERE id = $2;
    `

    res, err := db.Conn.ExecContext(ctx, query, weather, matchID)
    if err != nil {
        return fmt.Errorf("failed to update weather: %w", err)
    }

    rowsAffected, err := res.RowsAffected()
    if err != nil {
        return fmt.Errorf("failed to check rows affected: %w", err)
    }

    if rowsAffected == 0 {
        return fmt.Errorf("no match found with id %s", matchID)
    }

    return nil
}

func (db *DB) SetPosAndCompStats(ctx context.Context, matchID uuid.UUID, stats *PosAndComp) error {
    query := `
        INSERT INTO pos_and_comp (
            match_id,
            home_pos_per, away_pos_per,
            home_pos_time, away_pos_time,
            home_sets, home_sets_completed,
            away_sets, away_sets_completed
        )
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
        ON CONFLICT (match_id) DO UPDATE SET
            home_pos_per = EXCLUDED.home_pos_per,
            away_pos_per = EXCLUDED.away_pos_per,
            home_pos_time = EXCLUDED.home_pos_time,
            away_pos_time = EXCLUDED.away_pos_time,
            home_sets = EXCLUDED.home_sets,
            home_sets_completed = EXCLUDED.home_sets_completed,
            away_sets = EXCLUDED.away_sets,
            away_sets_completed = EXCLUDED.away_sets_completed;
    `

    _, err := db.Conn.ExecContext(ctx, query,
        matchID,
        stats.homePosPer, stats.awayPosPer,
        stats.homePosTime, stats.awayPosTime,
        stats.homeSets, stats.homeSetsCompleated,
        stats.awaySets, stats.awaySetsCompleated,
    )
    if err != nil {
        return fmt.Errorf("failed to set PosAndComp stats: %w", err)
    }

    return nil
}

func (db *DB) SetAttackStats(ctx context.Context, matchID uuid.UUID, a *Attack) error {
    // Using INSERT ... ON CONFLICT to upsert
    query := `
        INSERT INTO attack (
            match_id,
            home_runs, away_runs,
            home_run_meters, away_run_meters,
            home_post_contact_meters, away_post_contact_meters,
            home_line_breaks, away_line_breaks,
            home_tackle_breaks, away_tackle_breaks,
            home_avg_set_distance, away_avg_set_distance,
            home_kick_return_meters, away_kick_return_meters
        )
        VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15)
        ON CONFLICT (match_id)
        DO UPDATE SET
            home_runs = EXCLUDED.home_runs,
            away_runs = EXCLUDED.away_runs,
            home_run_meters = EXCLUDED.home_run_meters,
            away_run_meters = EXCLUDED.away_run_meters,
            home_post_contact_meters = EXCLUDED.home_post_contact_meters,
            away_post_contact_meters = EXCLUDED.away_post_contact_meters,
            home_line_breaks = EXCLUDED.home_line_breaks,
            away_line_breaks = EXCLUDED.away_line_breaks,
            home_tackle_breaks = EXCLUDED.home_tackle_breaks,
            away_tackle_breaks = EXCLUDED.away_tackle_breaks,
            home_avg_set_distance = EXCLUDED.home_avg_set_distance,
            away_avg_set_distance = EXCLUDED.away_avg_set_distance,
            home_kick_return_meters = EXCLUDED.home_kick_return_meters,
            away_kick_return_meters = EXCLUDED.away_kick_return_meters
    `

    _, err := db.Conn.ExecContext(ctx, query,
        matchID,
        a.homeRuns, a.awayRuns,
        a.homeRunMeters, a.awayRunMeters,
        a.homePostContactMeters, a.awayPostContactMeters,
        a.homeLineBreaks, a.awayLineBreaks,
        a.homeTackleBreaks, a.awayTackleBreaks,
        a.homeAvgSetDistance, a.awayAvgSetDistance,
        a.homeKickReturnMeters, a.awayKickReturnMeters,
    )

    return err
}

func (db *DB) SetPassingStats(ctx context.Context, matchID uuid.UUID, p *Passing) error {
    query := `
        INSERT INTO passing (
            match_id,
            home_offloads, away_offloads,
            home_receipts, away_receipts,
            home_total_passes, away_total_passes,
            home_dummy_passes, away_dummy_passes
        )
        VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
        ON CONFLICT (match_id)
        DO UPDATE SET
            home_offloads = EXCLUDED.home_offloads,
            away_offloads = EXCLUDED.away_offloads,
            home_receipts = EXCLUDED.home_receipts,
            away_receipts = EXCLUDED.away_receipts,
            home_total_passes = EXCLUDED.home_total_passes,
            away_total_passes = EXCLUDED.away_total_passes,
            home_dummy_passes = EXCLUDED.home_dummy_passes,
            away_dummy_passes = EXCLUDED.away_dummy_passes
    `

    _, err := db.Conn.ExecContext(ctx, query,
        matchID,
        p.homeOffloads, p.awayOffloads,
        p.homeReceipts, p.awayReceipts,
        p.homeTotalPasses, p.awayTotalPasses,
        p.homeDummyPasses, p.awayDummyPasses,
    )

    return err
}

func (db *DB) SetKickingStats(ctx context.Context, matchID uuid.UUID, k *Kicking) error {
    query := `
        INSERT INTO kicking (
            match_id,
            home_kicks, away_kicks,
            home_kicking_meters, away_kicking_meters,
            home_forced_drop_outs, away_forced_drop_outs,
            home_kick_defusal, away_kick_defusal,
            home_bombs, away_bombs,
            home_grubbers, away_grubbers
        )
        VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)
        ON CONFLICT (match_id)
        DO UPDATE SET
            home_kicks = EXCLUDED.home_kicks,
            away_kicks = EXCLUDED.away_kicks,
            home_kicking_meters = EXCLUDED.home_kicking_meters,
            away_kicking_meters = EXCLUDED.away_kicking_meters,
            home_forced_drop_outs = EXCLUDED.home_forced_drop_outs,
            away_forced_drop_outs = EXCLUDED.away_forced_drop_outs,
            home_kick_defusal = EXCLUDED.home_kick_defusal,
            away_kick_defusal = EXCLUDED.away_kick_defusal,
            home_bombs = EXCLUDED.home_bombs,
            away_bombs = EXCLUDED.away_bombs,
            home_grubbers = EXCLUDED.home_grubbers,
            away_grubbers = EXCLUDED.away_grubbers
    `

    _, err := db.Conn.ExecContext(ctx, query,
        matchID,
        k.homeKicks, k.awayKicks,
        k.homeKickingMeters, k.awayKickingMeters,
        k.homeForcedDropOuts, k.awayForcedDropOuts,
        k.homeKickDefusal, k.awayKickDefusal,
        k.homeBombs, k.awayBombs,
        k.homeGrubbers, k.awayGrubbers,
    )

    return err
}

func (db *DB) SetDefenceStats(ctx context.Context, matchID uuid.UUID, d *Defence) error {
    query := `
        INSERT INTO defence (
            match_id,
            home_effec_tackle, away_effec_tackle,
            home_tackles_made, away_tackles_made,
            home_missed_tackles, away_missed_tackles,
            home_intercepts, away_intercepts,
            home_ineffec_tackles, away_ineffec_tackles
        )
        VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
        ON CONFLICT (match_id)
        DO UPDATE SET
            home_effec_tackle = EXCLUDED.home_effec_tackle,
            away_effec_tackle = EXCLUDED.away_effec_tackle,
            home_tackles_made = EXCLUDED.home_tackles_made,
            away_tackles_made = EXCLUDED.away_tackles_made,
            home_missed_tackles = EXCLUDED.home_missed_tackles,
            away_missed_tackles = EXCLUDED.away_missed_tackles,
            home_intercepts = EXCLUDED.home_intercepts,
            away_intercepts = EXCLUDED.away_intercepts,
            home_ineffec_tackles = EXCLUDED.home_ineffec_tackles,
            away_ineffec_tackles = EXCLUDED.away_ineffec_tackles
    `

    _, err := db.Conn.ExecContext(ctx, query,
        matchID,
        d.homeEffecTackle, d.awayEffecTackle,
        d.homeTacklesMade, d.awayTacklesMade,
        d.homeMissedTackles, d.awayMissedTackles,
        d.homeIntercepts, d.awayIntercepts,
        d.homeIneffecTackles, d.awayIneffecTackles,
    )

    return err
}

func (db *DB) SetNegPlayStats(ctx context.Context, matchID uuid.UUID, ng *NegPlays) error {
    query := `
        INSERT INTO neg_plays (
            match_id,
            home_errors, away_errors,
            home_pen_con, away_pen_con,
            home_ruck_inf, away_ruck_inf,
            home_inside10, away_inside10,
            home_on_report, away_on_report
        )
        VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
        ON CONFLICT (match_id)
        DO UPDATE SET
            home_errors = EXCLUDED.home_errors,
            away_errors = EXCLUDED.away_errors,
            home_pen_con = EXCLUDED.home_pen_con,
            away_pen_con = EXCLUDED.away_pen_con,
            home_ruck_inf = EXCLUDED.home_ruck_inf,
            away_ruck_inf = EXCLUDED.away_ruck_inf,
            home_inside10 = EXCLUDED.home_inside10,
            away_inside10 = EXCLUDED.away_inside10,
            home_on_report = EXCLUDED.home_on_report,
            away_on_report = EXCLUDED.away_on_report
    `

    _, err := db.Conn.ExecContext(ctx, query,
        matchID,
        ng.homeErrors, ng.awayErrors,
        ng.homePenCon, ng.awayPenCon,
        ng.homeRuckInf, ng.awayRuckInf,
        ng.homeInside10, ng.awayInside10,
        ng.homeOnReport, ng.awayOnReport,
    )

    return err
}

func (db *DB) GetCompetition(id int) ([]Competition, error) {
	rows, err := db.Conn.Query("SELECT id, name FROM competition WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comps []Competition
	for rows.Next() {
        var c Competition
        if err := rows.Scan(&c.id, &c.name); err != nil {
            return []Competition{}, err
        }
    
        c.seasons, err = db.GetSeasons(c.id)
        if err != nil {
            return []Competition{}, err
        }
        comps = append(comps, c)
    }
	if err := rows.Err(); err != nil {
		return []Competition{}, err
	}

	return comps, nil
}

func (db *DB) GetSeasons(compId int) ([]*Season, error) {
	rows, err := db.Conn.Query("SELECT year, id FROM season WHERE competition_id = $1", compId)
	if err != nil {
		return []*Season{}, err
	}
	defer rows.Close()

	var seasons []*Season
	for rows.Next() {
		var s Season
		if err := rows.Scan(&s.year, &s.id); err != nil {
			return []*Season{}, err
		}
		s.rounds, err = db.GetRounds(s.id)
		seasons = append(seasons, &s)
	}
	if err := rows.Err(); err != nil {
		return []*Season{}, err
	}

	return seasons, nil
}

func (db *DB) GetRounds(seasonId uuid.UUID) ([]*Round, error) {
	rows, err := db.Conn.Query("SELECT id, start_day, end_day, round_name, round_index FROM round WHERE season_id = $1", seasonId)
	if err != nil {
		return []*Round{}, err
	}
	defer rows.Close()

	var rounds []*Round
	for rows.Next() {
		var r Round
		if err := rows.Scan(&r.id, &r.startDay, &r.endDay, &r.roundName, &r.roundIndex); err != nil {
			return []*Round{}, err
		}
		r.matches, err = db.GetMatches(r.id)
		rounds = append(rounds, &r)
	}
	if err := rows.Err(); err != nil {
		return []*Round{}, err
	}

	return rounds, nil
}

func (db *DB) GetMatches(roundId uuid.UUID) ([]*Match, error) {
	rows, err := db.Conn.Query(`
		SELECT
			id,
			home_team,
			away_team,
			home_score,
			away_score,
			location,
			kickoff_time,
			date_played,
			weather
		FROM
			match
		WHERE
			round_id = $1
	`, roundId)

	if err != nil {
		return []*Match{}, err
	}
	defer rows.Close()

	var matches []*Match
	for rows.Next() {
		var m Match
		if err := rows.Scan(
				&m.id,
				&m.homeTeam,
				&m.awayTeam,
				&m.homeScore,
				&m.awayScore,
				&m.location,
				&m.kickoffTime,
				&m.datePlayed,
				&m.weather,
			); err != nil {
			return []*Match{}, err
		}
	
        m.stats = db.GetMatchStats(m.id)
		matches = append(matches, &m)
	}
	if err := rows.Err(); err != nil {
		return []*Match{}, err
	}

	return matches, nil
}

func (db *DB) GetMatchStats(matchId uuid.UUID) (s *MatchStats) {
    s = &MatchStats{}
    s.posAndComp, _ = db.GetPosAndCompStats(matchId)
	s.attack, _ = db.GetAttackStats(matchId)
	s.passing, _ = db.GetPassingStats(matchId)
	s.kicking, _ = db.GetKickingStats(matchId)
	s.defence, _ = db.GetDefenceStats(matchId)
	s.negPlays, _ = db.GetNegPlaysStats(matchId)

	return
}

func (db *DB) GetPosAndCompStats(matchId uuid.UUID) (*PosAndComp, error) {
	stats := &PosAndComp{}

	err := db.Conn.QueryRow(`
		SELECT
            home_pos_per,
			away_pos_per,
			home_pos_time,
			away_pos_time,
			home_sets,
			home_sets_completed,
			away_sets,
			away_sets_completed
		FROM
            pos_and_comp
		WHERE
			match_id = $1
	`, matchId).Scan(
		&stats.homePosPer,
		&stats.awayPosPer,
		&stats.homePosTime,
		&stats.awayPosTime,
		&stats.homeSets,
		&stats.homeSetsCompleated,
		&stats.awaySets,
		&stats.awaySetsCompleated,
	)

	if err != nil {
		return &PosAndComp{}, err
	}

	return stats, nil
}

func (db *DB) GetAttackStats(matchId uuid.UUID) (*Attack, error) {
	stats := &Attack{}

	err := db.Conn.QueryRow(`
		SELECT
            home_runs,
			away_runs,
			home_run_meters,
			away_run_meters,
			home_post_contact_meters,
			away_post_contact_meters,
			home_line_breaks,
			away_line_breaks,
            home_tackle_breaks,
            away_tackle_breaks,
            home_avg_set_distance,
            away_avg_set_distance,
            home_kick_return_meters,
            away_kick_return_meters,
            home_avg_play_the_ball_speed,
            away_avg_play_the_ball_speed
		FROM
			attack
		WHERE
			match_id = $1
	`, matchId).Scan(
		&stats.homeRuns,
		&stats.awayRuns,
        &stats.homeRunMeters,
		&stats.awayRunMeters,
		&stats.homePostContactMeters,
		&stats.awayPostContactMeters,
		&stats.homeLineBreaks,
		&stats.awayLineBreaks,
		&stats.homeTackleBreaks,
		&stats.awayTackleBreaks,
        &stats.homeAvgSetDistance,
		&stats.awayAvgSetDistance,
        &stats.homeKickReturnMeters,
		&stats.awayKickReturnMeters,
        &stats.homeAvgPlayTheBallSpeed,
		&stats.awayAvgPlayTheBallSpeed,
	)

    fmt.Println(stats, err)
	if err != nil {
		return &Attack{}, err
	}

	return stats, nil
}

func (db *DB) GetPassingStats(matchId uuid.UUID) (*Passing, error) {
	stats := &Passing{}

	err := db.Conn.QueryRow(`
		SELECT
            home_offloads,
            away_offloads,
            home_receipts,
            away_receipts,
            home_total_passes,
            away_total_passes,
            home_dummy_passes,
            away_dummy_passes
		FROM
			passing
		WHERE
			match_id = $1
	`, matchId).Scan(
		&stats.homeOffloads,
		&stats.awayOffloads,
        &stats.homeReceipts,
		&stats.awayReceipts,
		&stats.homeTotalPasses,
		&stats.awayTotalPasses,
		&stats.homeDummyPasses,
		&stats.awayDummyPasses,
	)

	if err != nil {
		return &Passing{}, err
	}

	return stats, nil
}

func (db *DB) GetKickingStats(matchId uuid.UUID) (*Kicking, error) {
	stats := &Kicking{}
    
	err := db.Conn.QueryRow(`
		SELECT
            home_kicks,
            away_kicks,
            home_kicking_meters,
            away_kicking_meters,
            home_forced_drop_outs,
            away_forced_drop_outs,
            home_kick_defusal,
            home_kick_defusal,
            home_bombs,
            away_bombs,
            home_grubbers,
            away_grubbers
		FROM
			kicking
		WHERE
			match_id = $1
	`, matchId).Scan(
		&stats.homeKicks,
		&stats.awayKicks,
        &stats.homeKickingMeters,
		&stats.awayKickingMeters,
		&stats.homeForcedDropOuts,
		&stats.awayForcedDropOuts,
		&stats.homeKickDefusal,
		&stats.awayKickDefusal,
        &stats.homeBombs,
		&stats.awayBombs,
		&stats.homeGrubbers,
		&stats.awayGrubbers,
	)

	if err != nil {
		return &Kicking{}, err
	}

	return stats, nil
}

func (db *DB) GetDefenceStats(matchId uuid.UUID) (*Defence, error) {
	stats := &Defence{}

	err := db.Conn.QueryRow(`
		SELECT
            home_effec_tackle,
            away_effec_tackle,
            home_tackles_made,
            away_tackles_made,
            home_missed_tackles,
            away_missed_tackles,
            home_intercepts,
            away_intercepts,
            home_ineffec_tackles,
            away_ineffec_tackles
		FROM
			defence
		WHERE
			match_id = $1
	`, matchId).Scan(
		&stats.homeEffecTackle,
		&stats.awayEffecTackle,
        &stats.homeTacklesMade,
		&stats.awayTacklesMade,
		&stats.homeMissedTackles,
		&stats.awayMissedTackles,
		&stats.homeIntercepts,
		&stats.awayIntercepts,
        &stats.homeIneffecTackles,
		&stats.awayIneffecTackles,
	)

	if err != nil {
		return &Defence{}, err
	}

	return stats, nil
}

func (db *DB) GetNegPlaysStats(matchId uuid.UUID) (*NegPlays, error) {
	stats := &NegPlays{}

	err := db.Conn.QueryRow(`
		SELECT
            home_errors,
            away_errors,
            home_pen_con,
            away_pen_con,
            home_ruck_inf,
            away_ruck_inf,
            home_inside10,
            away_inside10,
            home_on_report,
            away_on_report
		FROM
            neg_plays
		WHERE
			match_id = $1
	`, matchId).Scan(
		&stats.homeErrors,
		&stats.awayErrors,
        &stats.homePenCon,
		&stats.awayPenCon,
		&stats.homeRuckInf,
		&stats.awayRuckInf,
		&stats.homeInside10,
		&stats.awayInside10,
        &stats.homeOnReport,
		&stats.awayOnReport,
	)

	if err != nil {
		return &NegPlays{}, err
	}

	return stats, nil
}
