ALTER TABLE neg_plays ADD COLUMN complete BOOLEAN DEFAULT FALSE;

CREATE OR REPLACE FUNCTION set_neg_plays_complete()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.home_errors    <> -1 AND
       NEW.away_errors    <> -1 AND
       NEW.home_pen_con   <> -1 AND
       NEW.away_pen_con   <> -1 AND
       NEW.home_ruck_inf  <> -1 AND
       NEW.away_ruck_inf  <> -1 AND
       NEW.home_inside10  <> -1 AND
       NEW.away_inside10  <> -1 AND
       NEW.home_on_report <> -1 AND
       NEW.away_on_report <> -1 THEN
        NEW.complete := TRUE;
    ELSE
        NEW.complete := FALSE;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER neg_plays_complete_trigger
BEFORE INSERT OR UPDATE ON neg_plays
FOR EACH ROW
EXECUTE FUNCTION set_neg_plays_complete();

ALTER TABLE defence ADD COLUMN complete BOOLEAN DEFAULT FALSE;

CREATE OR REPLACE FUNCTION set_defence_complete()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.home_effec_tackle   <> -1 AND
       NEW.away_effec_tackle   <> -1 AND
       NEW.home_tackles_made   <> -1 AND
       NEW.away_tackles_made   <> -1 AND
       NEW.home_missed_tackles <> -1 AND
       NEW.away_missed_tackles <> -1 AND
       NEW.home_intercepts     <> -1 AND
       NEW.away_intercepts     <> -1 AND
       NEW.home_ineffec_tackles <> -1 AND
       NEW.away_ineffec_tackles <> -1 THEN
        NEW.complete := TRUE;
    ELSE
        NEW.complete := FALSE;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER defence_complete_trigger
BEFORE INSERT OR UPDATE ON defence
FOR EACH ROW
EXECUTE FUNCTION set_defence_complete();

ALTER TABLE kicking ADD COLUMN complete BOOLEAN DEFAULT FALSE;

CREATE OR REPLACE FUNCTION set_kicking_complete()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.home_kicks           <> -1 AND
       NEW.away_kicks           <> -1 AND
       NEW.home_kicking_meters  <> -1 AND
       NEW.away_kicking_meters  <> -1 AND
       NEW.home_forced_drop_outs <> -1 AND
       NEW.away_forced_drop_outs <> -1 AND
       NEW.home_kick_defusal    <> -1 AND
       NEW.away_kick_defusal    <> -1 AND
       NEW.home_bombs           <> -1 AND
       NEW.away_bombs           <> -1 AND
       NEW.home_grubbers        <> -1 AND
       NEW.away_grubbers        <> -1 THEN
        NEW.complete := TRUE;
    ELSE
        NEW.complete := FALSE;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER kicking_complete_trigger
BEFORE INSERT OR UPDATE ON kicking
FOR EACH ROW
EXECUTE FUNCTION set_kicking_complete();

ALTER TABLE passing ADD COLUMN complete BOOLEAN DEFAULT FALSE;

CREATE OR REPLACE FUNCTION set_passing_complete()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.home_offloads      <> -1 AND
       NEW.away_offloads      <> -1 AND
       NEW.home_receipts      <> -1 AND
       NEW.away_receipts      <> -1 AND
       NEW.home_total_passes  <> -1 AND
       NEW.away_total_passes  <> -1 AND
       NEW.home_dummy_passes  <> -1 AND
       NEW.away_dummy_passes  <> -1 THEN
        NEW.complete := TRUE;
    ELSE
        NEW.complete := FALSE;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER passing_complete_trigger
BEFORE INSERT OR UPDATE ON passing
FOR EACH ROW
EXECUTE FUNCTION set_passing_complete();

ALTER TABLE attack ADD COLUMN complete BOOLEAN DEFAULT FALSE;

CREATE OR REPLACE FUNCTION set_attack_complete()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.home_runs                <> -1 AND
       NEW.away_runs                <> -1 AND
       NEW.home_run_meters          <> -1 AND
       NEW.away_run_meters          <> -1 AND
       NEW.home_post_contact_meters <> -1 AND
       NEW.away_post_contact_meters <> -1 AND
       NEW.home_line_breaks         <> -1 AND
       NEW.away_line_breaks         <> -1 AND
       NEW.home_tackle_breaks       <> -1 AND
       NEW.away_tackle_breaks       <> -1 AND
       NEW.home_avg_set_distance    <> -1 AND
       NEW.away_avg_set_distance    <> -1 AND
       NEW.home_kick_return_meters  <> -1 AND
       NEW.away_kick_return_meters  <> -1 AND
       NEW.home_avg_play_the_ball_speed <> -1 AND
       NEW.away_avg_play_the_ball_speed <> -1 THEN
        NEW.complete := TRUE;
    ELSE
        NEW.complete := FALSE;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER attack_complete_trigger
BEFORE INSERT OR UPDATE ON attack
FOR EACH ROW
EXECUTE FUNCTION set_attack_complete();

ALTER TABLE pos_and_comp ADD COLUMN complete BOOLEAN DEFAULT FALSE;

CREATE OR REPLACE FUNCTION set_pos_and_comp_complete()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.home_pos_per         <> -1 AND
       NEW.away_pos_per         <> -1 AND
       NEW.home_pos_time        <> '' AND
       NEW.away_pos_time        <> '' AND
       NEW.home_sets            <> -1 AND
       NEW.home_sets_completed  <> -1 AND
       NEW.away_sets            <> -1 AND
       NEW.away_sets_completed  <> -1 THEN
        NEW.complete := TRUE;
    ELSE
        NEW.complete := FALSE;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER pos_and_comp_complete_trigger
BEFORE INSERT OR UPDATE ON pos_and_comp
FOR EACH ROW
EXECUTE FUNCTION set_pos_and_comp_complete();

ALTER TABLE match ADD COLUMN complete BOOLEAN DEFAULT FALSE;
ALTER TABLE match ADD COLUMN children_complete BOOLEAN DEFAULT FALSE;

CREATE OR REPLACE FUNCTION set_match_complete()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.home_score <> -1 AND
       NEW.away_score <> -1 AND
       NEW.location <> '' AND
       NEW.kickoff_time <> '' AND
       NEW.date_played <> '' AND
       NEW.weather <> '' THEN
        NEW.complete := TRUE;
    ELSE
        NEW.complete := FALSE;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER match_complete_trigger
BEFORE INSERT OR UPDATE ON match
FOR EACH ROW
EXECUTE FUNCTION set_match_complete();

CREATE OR REPLACE FUNCTION set_match_children_complete()
RETURNS TRIGGER AS $$
DECLARE
    has_player BOOLEAN;
    has_official BOOLEAN;
    has_play BOOLEAN;
    pos_and_comp_complete BOOLEAN;
    attack_complete BOOLEAN;
    passing_complete BOOLEAN;
    kicking_complete BOOLEAN;
    defence_complete BOOLEAN;
    neg_plays_complete BOOLEAN;
BEGIN
    -- Check for at least one player
    SELECT EXISTS(SELECT 1 FROM player WHERE match_id = NEW.id) INTO has_player;
    -- Check for at least one match_official
    SELECT EXISTS(SELECT 1 FROM match_official WHERE match_id = NEW.id) INTO has_official;
    -- Check for at least one play_by_play
    SELECT EXISTS(SELECT 1 FROM play_by_play WHERE match_id = NEW.id) INTO has_play;
    SELECT complete FROM pos_and_comp WHERE match_id = NEW.id INTO pos_and_comp_complete;
    SELECT complete FROM attack WHERE match_id = NEW.id INTO attack_complete;
    SELECT complete FROM passing WHERE match_id = NEW.id INTO passing_complete;
    SELECT complete FROM kicking WHERE match_id = NEW.id INTO kicking_complete;
    SELECT complete FROM defence WHERE match_id = NEW.id INTO defence_complete;
    SELECT complete FROM neg_plays WHERE match_id = NEW.id INTO neg_plays_complete;

    IF has_player AND has_official AND has_play AND
       COALESCE(pos_and_comp_complete, FALSE) AND
       COALESCE(attack_complete, FALSE) AND
       COALESCE(passing_complete, FALSE) AND
       COALESCE(kicking_complete, FALSE) AND
       COALESCE(defence_complete, FALSE) AND
       COALESCE(neg_plays_complete, FALSE) THEN
        NEW.children_complete := TRUE;
    ELSE
        NEW.children_complete := FALSE;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER match_children_complete_trigger
BEFORE INSERT OR UPDATE ON match
FOR EACH ROW
EXECUTE FUNCTION set_match_children_complete();

CREATE OR REPLACE FUNCTION refresh_match_children_complete(match_row_id INT)
RETURNS VOID AS $$
BEGIN
    UPDATE match SET children_complete = children_complete WHERE id = match_row_id;
END;
$$ LANGUAGE plpgsql;

-- 1. pos_and_comp
CREATE OR REPLACE FUNCTION trg_refresh_match_on_pos_and_comp()
RETURNS TRIGGER AS $$
BEGIN
    PERFORM refresh_match_children_complete(NEW.match_id);
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER pos_and_comp_complete_refresh
AFTER UPDATE OF complete ON pos_and_comp
FOR EACH ROW
EXECUTE FUNCTION trg_refresh_match_on_pos_and_comp();

CREATE TRIGGER pos_and_comp_insert_refresh
AFTER INSERT ON pos_and_comp
FOR EACH ROW
EXECUTE FUNCTION trg_refresh_match_on_pos_and_comp();

CREATE TRIGGER pos_and_comp_delete_refresh
AFTER DELETE ON pos_and_comp
FOR EACH ROW
EXECUTE FUNCTION trg_refresh_match_on_pos_and_comp();

-- 2. attack
CREATE OR REPLACE FUNCTION trg_refresh_match_on_attack()
RETURNS TRIGGER AS $$
BEGIN
    PERFORM refresh_match_children_complete(NEW.match_id);
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER attack_complete_refresh
AFTER UPDATE OF complete ON attack
FOR EACH ROW
EXECUTE FUNCTION trg_refresh_match_on_attack();

CREATE TRIGGER attack_insert_refresh
AFTER INSERT ON attack
FOR EACH ROW
EXECUTE FUNCTION trg_refresh_match_on_attack();

CREATE TRIGGER attack_delete_refresh
AFTER DELETE ON attack
FOR EACH ROW
EXECUTE FUNCTION trg_refresh_match_on_attack();

-- 3. passing
CREATE OR REPLACE FUNCTION trg_refresh_match_on_passing()
RETURNS TRIGGER AS $$
BEGIN
    PERFORM refresh_match_children_complete(NEW.match_id);
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER passing_complete_refresh
AFTER UPDATE OF complete ON passing
FOR EACH ROW
EXECUTE FUNCTION trg_refresh_match_on_passing();

CREATE TRIGGER passing_insert_refresh
AFTER INSERT ON passing
FOR EACH ROW
EXECUTE FUNCTION trg_refresh_match_on_passing();

CREATE TRIGGER passing_delete_refresh
AFTER DELETE ON passing
FOR EACH ROW
EXECUTE FUNCTION trg_refresh_match_on_passing();

-- 4. kicking
CREATE OR REPLACE FUNCTION trg_refresh_match_on_kicking()
RETURNS TRIGGER AS $$
BEGIN
    PERFORM refresh_match_children_complete(NEW.match_id);
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER kicking_complete_refresh
AFTER UPDATE OF complete ON kicking
FOR EACH ROW
EXECUTE FUNCTION trg_refresh_match_on_kicking();

CREATE TRIGGER kicking_insert_refresh
AFTER INSERT ON kicking
FOR EACH ROW
EXECUTE FUNCTION trg_refresh_match_on_kicking();

CREATE TRIGGER kicking_delete_refresh
AFTER DELETE ON kicking
FOR EACH ROW
EXECUTE FUNCTION trg_refresh_match_on_kicking();

-- 5. defence
CREATE OR REPLACE FUNCTION trg_refresh_match_on_defence()
RETURNS TRIGGER AS $$
BEGIN
    PERFORM refresh_match_children_complete(NEW.match_id);
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER defence_complete_refresh
AFTER UPDATE OF complete ON defence
FOR EACH ROW
EXECUTE FUNCTION trg_refresh_match_on_defence();

CREATE TRIGGER defence_insert_refresh
AFTER INSERT ON defence
FOR EACH ROW
EXECUTE FUNCTION trg_refresh_match_on_defence();

CREATE TRIGGER defence_delete_refresh
AFTER DELETE ON defence
FOR EACH ROW
EXECUTE FUNCTION trg_refresh_match_on_defence();

-- 6. neg_plays
CREATE OR REPLACE FUNCTION trg_refresh_match_on_neg_plays()
RETURNS TRIGGER AS $$
BEGIN
    PERFORM refresh_match_children_complete(NEW.match_id);
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER neg_plays_complete_refresh
AFTER UPDATE OF complete ON neg_plays
FOR EACH ROW
EXECUTE FUNCTION trg_refresh_match_on_neg_plays();

CREATE TRIGGER neg_plays_insert_refresh
AFTER INSERT ON neg_plays
FOR EACH ROW
EXECUTE FUNCTION trg_refresh_match_on_neg_plays();

CREATE TRIGGER neg_plays_delete_refresh
AFTER DELETE ON neg_plays
FOR EACH ROW
EXECUTE FUNCTION trg_refresh_match_on_neg_plays();

-- 7. player (has_player)
CREATE OR REPLACE FUNCTION trg_refresh_match_on_player()
RETURNS TRIGGER AS $$
BEGIN
    PERFORM refresh_match_children_complete(NEW.match_id);
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER player_insert_refresh
AFTER INSERT ON player
FOR EACH ROW
EXECUTE FUNCTION trg_refresh_match_on_player();

CREATE TRIGGER player_delete_refresh
AFTER DELETE ON player
FOR EACH ROW
EXECUTE FUNCTION trg_refresh_match_on_player();

-- 8. match_official (has_official)
CREATE OR REPLACE FUNCTION trg_refresh_match_on_official()
RETURNS TRIGGER AS $$
BEGIN
    PERFORM refresh_match_children_complete(NEW.match_id);
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER match_official_insert_refresh
AFTER INSERT ON match_official
FOR EACH ROW
EXECUTE FUNCTION trg_refresh_match_on_official();

CREATE TRIGGER match_official_delete_refresh
AFTER DELETE ON match_official
FOR EACH ROW
EXECUTE FUNCTION trg_refresh_match_on_official();

-- 9. play_by_play (has_play)
CREATE OR REPLACE FUNCTION trg_refresh_match_on_play()
RETURNS TRIGGER AS $$
BEGIN
    PERFORM refresh_match_children_complete(NEW.match_id);
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER play_by_play_insert_refresh
AFTER INSERT ON play_by_play
FOR EACH ROW
EXECUTE FUNCTION trg_refresh_match_on_play();

CREATE TRIGGER play_by_play_delete_refresh
AFTER DELETE ON play_by_play
FOR EACH ROW
EXECUTE FUNCTION trg_refresh_match_on_play();

ALTER TABLE round ADD COLUMN complete BOOLEAN DEFAULT FALSE;
ALTER TABLE round ADD COLUMN children_complete BOOLEAN DEFAULT FALSE;

-- Set round.complete based on required fields
CREATE OR REPLACE FUNCTION set_round_complete()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.round_name <> '' AND
       NEW.round_index <> -1 AND
       NEW.start_day <> '' AND
       NEW.end_day <> '' THEN
        NEW.complete := TRUE;
    ELSE
        NEW.complete := FALSE;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER round_complete_trigger
BEFORE INSERT OR UPDATE ON round
FOR EACH ROW
EXECUTE FUNCTION set_round_complete();

-- Set round.children_complete based on matches
CREATE OR REPLACE FUNCTION set_round_children_complete()
RETURNS TRIGGER AS $$
DECLARE
    has_match BOOLEAN;
    all_matches_complete BOOLEAN;
    all_matches_children_complete BOOLEAN;
BEGIN
    -- At least one match
    SELECT EXISTS(SELECT 1 FROM match WHERE round_id = NEW.id) INTO has_match;
    -- All matches complete
    SELECT bool_and(complete) FROM match WHERE round_id = NEW.id INTO all_matches_complete;
    -- All matches children_complete
    SELECT bool_and(children_complete) FROM match WHERE round_id = NEW.id INTO all_matches_children_complete;

    IF has_match AND COALESCE(all_matches_complete, FALSE) AND COALESCE(all_matches_children_complete, FALSE) THEN
        NEW.children_complete := TRUE;
    ELSE
        NEW.children_complete := FALSE;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER round_children_complete_trigger
BEFORE INSERT OR UPDATE ON round
FOR EACH ROW
EXECUTE FUNCTION set_round_children_complete();

-- Helper to refresh round.children_complete
CREATE OR REPLACE FUNCTION refresh_round_children_complete(round_row_id UUID)
RETURNS VOID AS $$
BEGIN
    UPDATE round SET children_complete = children_complete WHERE id = round_row_id;
END;
$$ LANGUAGE plpgsql;

-- Triggers on match to refresh round.children_complete
CREATE OR REPLACE FUNCTION trg_refresh_round_on_match()
RETURNS TRIGGER AS $$
BEGIN
    PERFORM refresh_round_children_complete(NEW.round_id);
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION trg_refresh_round_on_match_delete()
RETURNS TRIGGER AS $$
BEGIN
    PERFORM refresh_round_children_complete(OLD.round_id);
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER match_complete_update_refresh_round
AFTER UPDATE OF complete, children_complete ON match
FOR EACH ROW
EXECUTE FUNCTION trg_refresh_round_on_match();

CREATE TRIGGER match_insert_refresh_round
AFTER INSERT ON match
FOR EACH ROW
EXECUTE FUNCTION trg_refresh_round_on_match();

CREATE TRIGGER match_delete_refresh_round
AFTER DELETE ON match
FOR EACH ROW
EXECUTE FUNCTION trg_refresh_round_on_match_delete();

ALTER TABLE season ADD COLUMN complete BOOLEAN DEFAULT FALSE;

CREATE OR REPLACE FUNCTION set_season_complete()
RETURNS TRIGGER AS $$
DECLARE
    has_round BOOLEAN;
    all_rounds_complete BOOLEAN;
    all_rounds_children_complete BOOLEAN;
BEGIN
    -- At least one round
    SELECT EXISTS(SELECT 1 FROM round WHERE season_id = NEW.id) INTO has_round;
    -- All rounds complete
    SELECT bool_and(complete) FROM round WHERE season_id = NEW.id INTO all_rounds_complete;
    -- All rounds children_complete
    SELECT bool_and(children_complete) FROM round WHERE season_id = NEW.id INTO all_rounds_children_complete;

    IF has_round AND COALESCE(all_rounds_complete, FALSE) AND COALESCE(all_rounds_children_complete, FALSE) THEN
        NEW.complete := TRUE;
    ELSE
        NEW.complete := FALSE;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER season_complete_trigger
BEFORE INSERT OR UPDATE ON season
FOR EACH ROW
EXECUTE FUNCTION set_season_complete();

CREATE OR REPLACE FUNCTION refresh_season_complete(season_row_id UUID)
RETURNS VOID AS $$
BEGIN
    UPDATE season SET complete = complete WHERE id = season_row_id;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION trg_refresh_season_on_round()
RETURNS TRIGGER AS $$
BEGIN
    PERFORM refresh_season_complete(NEW.season_id);
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION trg_refresh_season_on_round_delete()
RETURNS TRIGGER AS $$
BEGIN
    PERFORM refresh_season_complete(OLD.season_id);
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER round_complete_update_refresh_season
AFTER UPDATE OF complete, children_complete ON round
FOR EACH ROW
EXECUTE FUNCTION trg_refresh_season_on_round();

CREATE TRIGGER round_insert_refresh_season
AFTER INSERT ON round
FOR EACH ROW
EXECUTE FUNCTION trg_refresh_season_on_round();

CREATE TRIGGER round_delete_refresh_season
AFTER DELETE ON round
FOR EACH ROW
EXECUTE FUNCTION trg_refresh_season_on_round_delete();

