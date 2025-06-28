-- Migration: Performance Indexes for High-Volume Operations
-- This migration adds indexes to optimize frequently queried columns

-- Call Detail Records (CDR) indexes
CREATE INDEX IF NOT EXISTS idx_cdr_start_time ON call_detail_records(start_time);
CREATE INDEX IF NOT EXISTS idx_cdr_sim_id ON call_detail_records(sim_id);
CREATE INDEX IF NOT EXISTS idx_cdr_disposition ON call_detail_records(disposition);
CREATE INDEX IF NOT EXISTS idx_cdr_voice_category ON call_detail_records(voice_category);
CREATE INDEX IF NOT EXISTS idx_cdr_destination ON call_detail_records(destination_number);
CREATE INDEX IF NOT EXISTS idx_cdr_composite ON call_detail_records(start_time, disposition, sim_id);

-- SIP Calls indexes for spam detection queries
CREATE INDEX IF NOT EXISTS idx_sip_calls_created_at ON sip_calls(created_at);
CREATE INDEX IF NOT EXISTS idx_sip_calls_voice_category ON sip_calls(voice_category);
CREATE INDEX IF NOT EXISTS idx_sip_calls_routed_to_ai ON sip_calls(routed_to_ai);
CREATE INDEX IF NOT EXISTS idx_sip_calls_caller_id ON sip_calls(caller_id_num);
CREATE INDEX IF NOT EXISTS idx_sip_calls_spam_composite ON sip_calls(voice_category, created_at) WHERE voice_category = 'SPAM_ROBOCALL';

-- SIM Cards indexes for status monitoring
CREATE INDEX IF NOT EXISTS idx_sim_cards_status ON sim_cards(status);
CREATE INDEX IF NOT EXISTS idx_sim_cards_operator ON sim_cards(operator_name);
CREATE INDEX IF NOT EXISTS idx_sim_cards_credit ON sim_cards(current_credit);
CREATE INDEX IF NOT EXISTS idx_sim_cards_active_composite ON sim_cards(status, current_credit) WHERE status = 'active';

-- WhatsApp Validation Cache indexes
CREATE INDEX IF NOT EXISTS idx_whatsapp_cache_phone ON whatsapp_validation_cache(phone_number);
CREATE INDEX IF NOT EXISTS idx_whatsapp_cache_checked ON whatsapp_validation_cache(checked_at);
CREATE INDEX IF NOT EXISTS idx_whatsapp_cache_expired ON whatsapp_validation_cache(checked_at) WHERE checked_at < NOW() - INTERVAL '24 hours';

-- Voice Recognition Results indexes
CREATE INDEX IF NOT EXISTS idx_voice_results_call_id ON voice_recognition_results(call_id);
CREATE INDEX IF NOT EXISTS idx_voice_results_created ON voice_recognition_results(created_at);
CREATE INDEX IF NOT EXISTS idx_voice_results_category ON voice_recognition_results(classification_category);
CREATE INDEX IF NOT EXISTS idx_voice_results_sim_status ON voice_recognition_results(sim_status_detected);

-- AI Agent Interactions indexes
CREATE INDEX IF NOT EXISTS idx_ai_interactions_start ON ai_agent_interactions(start_time);
CREATE INDEX IF NOT EXISTS idx_ai_interactions_agent ON ai_agent_interactions(agent_id);
CREATE INDEX IF NOT EXISTS idx_ai_interactions_call ON ai_agent_interactions(call_id);

-- Blacklist indexes for fast lookup
CREATE INDEX IF NOT EXISTS idx_blacklist_number ON blacklist(phone_number);
CREATE INDEX IF NOT EXISTS idx_blacklist_active ON blacklist(is_active) WHERE is_active = true;

-- Customer indexes
CREATE INDEX IF NOT EXISTS idx_customers_active ON customers(is_active);
CREATE INDEX IF NOT EXISTS idx_customers_balance ON customers(balance);

-- SIM Replacement Queue indexes
CREATE INDEX IF NOT EXISTS idx_sim_replacement_status ON sim_replacement_queue(status);
CREATE INDEX IF NOT EXISTS idx_sim_replacement_priority ON sim_replacement_queue(priority, requested_at);

-- Partial indexes for common queries
CREATE INDEX IF NOT EXISTS idx_recent_calls ON call_detail_records(start_time) WHERE start_time > NOW() - INTERVAL '7 days';
CREATE INDEX IF NOT EXISTS idx_low_credit_sims ON sim_cards(current_credit, id) WHERE current_credit < 10 AND status = 'active';
CREATE INDEX IF NOT EXISTS idx_pending_replacements ON sim_replacement_queue(priority, requested_at) WHERE status = 'pending';

-- Function-based index for hour extraction (for hourly analytics)
CREATE INDEX IF NOT EXISTS idx_cdr_hour ON call_detail_records(EXTRACT(HOUR FROM start_time), start_time);

-- Analyze tables to update statistics
ANALYZE call_detail_records;
ANALYZE sip_calls;
ANALYZE sim_cards;
ANALYZE whatsapp_validation_cache;
ANALYZE voice_recognition_results;
ANALYZE ai_agent_interactions;