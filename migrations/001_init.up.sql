CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE competition (
    id INT PRIMARY KEY,
    "name" VARCHAR(255) NOT NULL
);

CREATE TABLE season (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    competition_id INT NOT NULL REFERENCES competition(id) ON DELETE CASCADE,
    "year" VARCHAR(10) NOT NULL,

    UNIQUE(competition_id, "year")
);

CREATE TABLE round (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    season_id UUID NOT NULL REFERENCES season(id) ON DELETE CASCADE,
    round_name VARCHAR(100) NOT NULL,
    round_index INT NOT NULL,
    start_day VARCHAR(255) NOT NULL DEFAULT '',
    end_day VARCHAR(255) NOT NULL DEFAULT '',

    UNIQUE(season_id, round_index)
);

CREATE TABLE match (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    round_id UUID NOT NULL REFERENCES round(id) ON DELETE CASCADE,
    home_team VARCHAR(100) NOT NULL,
    away_team VARCHAR(100) NOT NULL,
    home_score INT NOT NULL DEFAULT -1,
    away_score INT NOT NULL DEFAULT -1,
    location VARCHAR(255) NOT NULL DEFAULT '',
    kickoff_time VARCHAR(20) NOT NULL DEFAULT '',
    date_played VARCHAR(20) NOT NULL DEFAULT '',
    weather VARCHAR(255) Not NULL DEFAULT '',

    UNIQUE(round_id, home_team, away_team)
);

CREATE TABLE player (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    match_id UUID NOT NULL REFERENCES match(id),
    name_first VARCHAR(100) NOT NULL,
    name_last VARCHAR(100) NOT NULL,
    position VARCHAR(50) NOT NULL,
    number INT NOT NULL,

    UNIQUE (match_id, name_first, name_last)
);


CREATE TABLE match_player (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    match_id UUID NOT NULL REFERENCES match(id) ON DELETE CASCADE,
    player_id UUID NOT NULL REFERENCES player(id) ON DELETE CASCADE,
    team VARCHAR(10) NOT NULL CHECK (team IN ('home', 'away'))
);

CREATE TABLE match_official (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    match_id UUID NOT NULL REFERENCES match(id) ON DELETE CASCADE,
    name_first VARCHAR(100) NOT NULL,
    name_last VARCHAR(100) NOT NULL,
    role VARCHAR(100)
);

CREATE TABLE play_by_play (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    match_id UUID NOT NULL REFERENCES match(id) ON DELETE CASCADE,
    play_index INT NOT NULL,
    time VARCHAR(20) NOT NULL,
    play TEXT NOT NULL,
    team VARCHAR(10),
    notes TEXT,

    UNIQUE (match_id, play_index)
);

CREATE TABLE pos_and_comp (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    match_id UUID NOT NULL REFERENCES match(id) ON DELETE CASCADE,
    home_pos_per INT DEFAULT -1,
    away_pos_per INT DEFAULT -1,
    home_pos_time VARCHAR(20) DEFAULT '',
    away_pos_time VARCHAR(20) DEFAULT '',
    home_sets INT DEFAULT -1,
    home_sets_completed INT DEFAULT -1,
    away_sets INT DEFAULT -1,
    away_sets_completed INT DEFAULT -1,
    UNIQUE(match_id)
);

CREATE TABLE attack (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    match_id UUID NOT NULL REFERENCES match(id) ON DELETE CASCADE,
    home_runs INT DEFAULT -1,
    away_runs INT DEFAULT -1,
    home_run_meters INT DEFAULT -1,
    away_run_meters INT DEFAULT -1,
    home_post_contact_meters INT DEFAULT -1,
    away_post_contact_meters INT DEFAULT -1,
    home_line_breaks INT DEFAULT -1,
    away_line_breaks INT DEFAULT -1,
    home_tackle_breaks INT DEFAULT -1,
    away_tackle_breaks INT DEFAULT -1,
    home_avg_set_distance FLOAT DEFAULT -1,
    away_avg_set_distance FLOAT DEFAULT -1,
    home_kick_return_meters INT DEFAULT -1,
    away_kick_return_meters INT DEFAULT -1,
    home_avg_play_the_ball_speed FLOAT DEFAULT -1,
    away_avg_play_the_ball_speed FLOAT DEFAULT -1,
    UNIQUE(match_id)
);

CREATE TABLE passing (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    match_id UUID NOT NULL REFERENCES match(id) ON DELETE CASCADE,
    home_offloads INT DEFAULT -1,
    away_offloads INT DEFAULT -1,
    home_receipts INT DEFAULT -1,
    away_receipts INT DEFAULT -1,
    home_total_passes INT DEFAULT -1,
    away_total_passes INT DEFAULT -1,
    home_dummy_passes INT DEFAULT -1,
    away_dummy_passes INT DEFAULT -1,
    UNIQUE(match_id)
);

CREATE TABLE kicking (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    match_id UUID NOT NULL REFERENCES match(id) ON DELETE CASCADE,
    home_kicks INT DEFAULT -1,
    away_kicks INT DEFAULT -1,
    home_kicking_meters INT DEFAULT -1,
    away_kicking_meters INT DEFAULT -1,
    home_forced_drop_outs INT DEFAULT -1,
    away_forced_drop_outs INT DEFAULT -1,
    home_kick_defusal INT DEFAULT -1,
    away_kick_defusal INT DEFAULT -1,
    home_bombs INT DEFAULT -1,
    away_bombs INT DEFAULT -1,
    home_grubbers INT DEFAULT -1,
    away_grubbers INT DEFAULT -1,
    UNIQUE(match_id)
);

CREATE TABLE defence (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    match_id UUID NOT NULL REFERENCES match(id) ON DELETE CASCADE,
    home_effec_tackle FLOAT DEFAULT -1,
    away_effec_tackle FLOAT DEFAULT -1,
    home_tackles_made INT DEFAULT -1,
    away_tackles_made INT DEFAULT -1,
    home_missed_tackles INT DEFAULT -1,
    away_missed_tackles INT DEFAULT -1,
    home_intercepts INT DEFAULT -1,
    away_intercepts INT DEFAULT -1,
    home_ineffec_tackles INT DEFAULT -1,
    away_ineffec_tackles INT DEFAULT -1,
    UNIQUE(match_id)
);

CREATE TABLE neg_plays (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    match_id UUID NOT NULL REFERENCES match(id) ON DELETE CASCADE,
    home_errors INT DEFAULT -1,
    away_errors INT DEFAULT -1,
    home_pen_con INT DEFAULT -1,
    away_pen_con INT DEFAULT -1,
    home_ruck_inf INT DEFAULT -1,
    away_ruck_inf INT DEFAULT -1,
    home_inside10 INT DEFAULT -1,
    away_inside10 INT DEFAULT -1,
    home_on_report INT DEFAULT -1,
    away_on_report INT DEFAULT -1,
    UNIQUE(match_id)
);

