-- Drop voice recognition related columns from sip_calls
ALTER TABLE sip_calls DROP COLUMN IF EXISTS voice_analyzed;
ALTER TABLE sip_calls DROP COLUMN IF EXISTS voice_category;
ALTER TABLE sip_calls DROP COLUMN IF EXISTS voice_action;
ALTER TABLE sip_calls DROP COLUMN IF EXISTS routed_to_ai;
ALTER TABLE sip_calls DROP COLUMN IF EXISTS ai_agent_id;

-- Drop SIM card related columns
ALTER TABLE sim_cards DROP COLUMN IF EXISTS last_issue;
ALTER TABLE sim_cards DROP COLUMN IF EXISTS flagged_at;
ALTER TABLE sim_cards DROP COLUMN IF EXISTS action_taken;

-- Drop all voice recognition tables
DROP TABLE IF EXISTS voice_transcripts;
DROP TABLE IF EXISTS ai_agent_interactions;
DROP TABLE IF EXISTS sim_replacement_queue;
DROP TABLE IF EXISTS call_reviews;
DROP TABLE IF EXISTS voice_recognition_logs;