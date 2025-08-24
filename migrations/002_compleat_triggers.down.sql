DROP TRIGGER IF EXISTS neg_plays_complete_trigger ON neg_plays;
DROP TRIGGER IF EXISTS neg_plays_complete_refresh ON neg_plays;
DROP TRIGGER IF EXISTS neg_plays_insert_refresh ON neg_plays;
DROP TRIGGER IF EXISTS neg_plays_delete_refresh ON neg_plays;
DROP FUNCTION IF EXISTS set_neg_plays_complete();
DROP FUNCTION IF EXISTS trg_refresh_match_on_neg_plays();
ALTER TABLE neg_plays DROP COLUMN IF EXISTS complete;

DROP TRIGGER IF EXISTS defence_complete_trigger ON defence;
DROP TRIGGER IF EXISTS defence_complete_refresh ON defence;
DROP TRIGGER IF EXISTS defence_insert_refresh ON defence;
DROP TRIGGER IF EXISTS defence_delete_refresh ON defence;
DROP FUNCTION IF EXISTS set_defence_complete();
DROP FUNCTION IF EXISTS trg_refresh_match_on_defence();
ALTER TABLE defence DROP COLUMN IF EXISTS complete;

DROP TRIGGER IF EXISTS kicking_complete_trigger ON kicking;
DROP TRIGGER IF EXISTS kicking_complete_refresh ON kicking;
DROP TRIGGER IF EXISTS kicking_insert_refresh ON kicking;
DROP TRIGGER IF EXISTS kicking_delete_refresh ON kicking;
DROP FUNCTION IF EXISTS set_kicking_complete();
DROP FUNCTION IF EXISTS trg_refresh_match_on_kicking();
ALTER TABLE kicking DROP COLUMN IF EXISTS complete;

DROP TRIGGER IF EXISTS passing_complete_trigger ON passing;
DROP TRIGGER IF EXISTS passing_complete_refresh ON passing;
DROP TRIGGER IF EXISTS passing_insert_refresh ON passing;
DROP TRIGGER IF EXISTS passing_delete_refresh ON passing;
DROP FUNCTION IF EXISTS set_passing_complete();
DROP FUNCTION IF EXISTS trg_refresh_match_on_passing();
ALTER TABLE passing DROP COLUMN IF EXISTS complete;

DROP TRIGGER IF EXISTS attack_complete_trigger ON attack;
DROP TRIGGER IF EXISTS attack_complete_refresh ON attack;
DROP TRIGGER IF EXISTS attack_insert_refresh ON attack;
DROP TRIGGER IF EXISTS attack_delete_refresh ON attack;
DROP FUNCTION IF EXISTS set_attack_complete();
DROP FUNCTION IF EXISTS trg_refresh_match_on_attack();
ALTER TABLE attack DROP COLUMN IF EXISTS complete;

DROP TRIGGER IF EXISTS pos_and_comp_complete_trigger ON pos_and_comp;
DROP TRIGGER IF EXISTS pos_and_comp_complete_refresh ON pos_and_comp;
DROP TRIGGER IF EXISTS pos_and_comp_insert_refresh ON pos_and_comp;
DROP TRIGGER IF EXISTS pos_and_comp_delete_refresh ON pos_and_comp;
DROP FUNCTION IF EXISTS set_pos_and_comp_complete();
DROP FUNCTION IF EXISTS trg_refresh_match_on_pos_and_comp();
ALTER TABLE pos_and_comp DROP COLUMN IF EXISTS complete;

DROP TRIGGER IF EXISTS match_complete_trigger ON match;
DROP TRIGGER IF EXISTS match_children_complete_trigger ON match;
DROP TRIGGER IF EXISTS match_complete_update_refresh_round ON match;
DROP TRIGGER IF EXISTS match_insert_refresh_round ON match;
DROP TRIGGER IF EXISTS match_delete_refresh_round ON match;
DROP FUNCTION IF EXISTS set_match_complete();
DROP FUNCTION IF EXISTS set_match_children_complete();
DROP FUNCTION IF EXISTS refresh_match_children_complete(INT);
DROP FUNCTION IF EXISTS trg_refresh_round_on_match();
DROP FUNCTION IF EXISTS trg_refresh_round_on_match_delete();
ALTER TABLE match DROP COLUMN IF EXISTS complete;
ALTER TABLE match DROP COLUMN IF EXISTS children_complete;

DROP TRIGGER IF EXISTS player_insert_refresh ON player;
DROP TRIGGER IF EXISTS player_delete_refresh ON player;
DROP FUNCTION IF EXISTS trg_refresh_match_on_player();

DROP TRIGGER IF EXISTS match_official_insert_refresh ON match_official;
DROP TRIGGER IF EXISTS match_official_delete_refresh ON match_official;
DROP FUNCTION IF EXISTS trg_refresh_match_on_official();

DROP TRIGGER IF EXISTS play_by_play_insert_refresh ON play_by_play;
DROP TRIGGER IF EXISTS play_by_play_delete_refresh ON play_by_play;
DROP FUNCTION IF EXISTS trg_refresh_match_on_play();

DROP TRIGGER IF EXISTS round_complete_trigger ON round;
DROP TRIGGER IF EXISTS round_children_complete_trigger ON round;
DROP TRIGGER IF EXISTS round_complete_update_refresh_season ON round;
DROP TRIGGER IF EXISTS round_insert_refresh_season ON round;
DROP TRIGGER IF EXISTS round_delete_refresh_season ON round;
DROP FUNCTION IF EXISTS set_round_complete();
DROP FUNCTION IF EXISTS set_round_children_complete();
DROP FUNCTION IF EXISTS refresh_round_children_complete(UUID);
DROP FUNCTION IF EXISTS trg_refresh_season_on_round();
DROP FUNCTION IF EXISTS trg_refresh_season_on_round_delete();
ALTER TABLE round DROP COLUMN IF EXISTS complete;
ALTER TABLE round DROP COLUMN IF EXISTS children_complete;

DROP TRIGGER IF EXISTS season_complete_trigger ON season;
DROP FUNCTION IF EXISTS set_season_complete();
DROP FUNCTION IF EXISTS refresh_season_complete(UUID);
ALTER TABLE season DROP COLUMN IF EXISTS complete;

