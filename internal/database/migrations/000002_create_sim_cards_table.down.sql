-- 000002_create_sim_cards_table.down.sql
DROP TRIGGER IF EXISTS set_sim_cards_updated_at ON sim_cards;
DROP TABLE IF EXISTS sim_cards;
-- The trigger_set_timestamp function is likely shared, so we might not want to drop it here
-- unless it's exclusively for this table and created in this migration's .up.sql.
-- If it was created in 000001_create_modems_table.up.sql, it should be dropped in its .down.sql.
-- For safety, let's assume it's potentially shared and avoid dropping it here.
-- If it's defined idempotently (CREATE OR REPLACE) in each .up.sql, then no drop is needed here.