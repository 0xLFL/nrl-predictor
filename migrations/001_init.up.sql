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
    start_day VARCHAR(255) DEFAULT NULL,
    end_day VARCHAR(255) DEFAULT NULL,

    UNIQUE(season_id, round_index)
);

CREATE TABLE match (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    round_id UUID NOT NULL REFERENCES round(id) ON DELETE CASCADE,
    home_team VARCHAR(100) NOT NULL,
    away_team VARCHAR(100) NOT NULL,
    home_score INT NOT NULL DEFAULT -1,
    away_score INT NOT NULL DEFAULT -1,
    location VARCHAR(255),
    kickoff_time TIME,
    date_played DATE,
    weather VARCHAR(255),

    UNIQUE(round_id, home_team, away_team)
);

CREATE TABLE player (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    match_id UUID NOT NULL REFERENCES match(id),
    name_first VARCHAR(100) NOT NULL,
    name_last VARCHAR(100) NOT NULL,
    position VARCHAR(50),
    number INT,

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
    home_pos_per INT,
    away_pos_per INT,
    home_pos_time VARCHAR(20),
    away_pos_time VARCHAR(20),
    home_sets INT,
    home_sets_completed INT,
    away_sets INT,
    away_sets_completed INT
);

CREATE TABLE attack (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    match_id UUID NOT NULL REFERENCES match(id) ON DELETE CASCADE,
    home_runs INT,
    away_runs INT,
    home_run_meters INT,
    away_run_meters INT,
    home_post_contact_meters INT,
    away_post_contact_meters INT,
    home_line_breaks INT,
    away_line_breaks INT,
    home_tackle_breaks INT,
    away_tackle_breaks INT,
    home_avg_set_distance FLOAT,
    away_avg_set_distance FLOAT,
    home_kick_return_meters INT,
    away_kick_return_meters INT,
    home_avg_play_the_ball_speed FLOAT,
    away_avg_play_the_ball_speed FLOAT
);

CREATE TABLE passing (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    match_id UUID NOT NULL REFERENCES match(id) ON DELETE CASCADE,
    home_offloads INT,
    away_offloads INT,
    home_receipts INT,
    away_receipts INT,
    home_total_passes INT,
    away_total_passes INT,
    home_dummy_passes INT,
    away_dummy_passes INT
);

CREATE TABLE kicking (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    match_id UUID NOT NULL REFERENCES match(id) ON DELETE CASCADE,
    home_kicks INT,
    away_kicks INT,
    home_kicking_meters INT,
    away_kicking_meters INT,
    home_forced_drop_outs INT,
    away_forced_drop_outs INT,
    home_kick_defusal INT,
    away_kick_defusal INT,
    home_bombs INT,
    away_bombs INT,
    home_grubbers INT,
    away_grubbers INT
);

CREATE TABLE defence (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    match_id UUID NOT NULL REFERENCES match(id) ON DELETE CASCADE,
    home_effec_tackle FLOAT,
    away_effec_tackle FLOAT,
    home_tackles_made INT,
    away_tackles_made INT,
    home_missed_tackles INT,
    away_missed_tackles INT,
    home_intercepts INT,
    away_intercepts INT,
    home_ineffec_tackles INT,
    away_ineffec_tackles INT
);

CREATE TABLE neg_plays (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    match_id UUID NOT NULL REFERENCES match(id) ON DELETE CASCADE,
    home_errors INT,
    away_errors INT,
    home_pen_con INT,
    away_pen_con INT,
    home_ruck_inf INT,
    away_ruck_inf INT,
    home_inside10 INT,
    away_inside10 INT,
    home_on_report INT,
    away_on_report INT
);
