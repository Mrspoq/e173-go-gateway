-- Migration: Performance Indexes for High-Volume Operations (Fixed)
-- This migration adds indexes to optimize frequently queried columns

-- Call Detail Records (CDR) indexes
CREATE INDEX IF NOT EXISTS idx_cdr_call_start_time ON call_detail_records(call_start_time);
CREATE INDEX IF NOT EXISTS idx_cdr_sim_card_id ON call_detail_records(sim_card_id);
CREATE INDEX IF NOT EXISTS idx_cdr_disposition ON call_detail_records(disposition);
CREATE INDEX IF NOT EXISTS idx_cdr_destination ON call_detail_records(destination_number);
CREATE INDEX IF NOT EXISTS idx_cdr_composite ON call_detail_records(call_start_time, disposition, sim_card_id);
CREATE INDEX IF NOT EXISTS idx_cdr_customer ON call_detail_records(customer_id);

-- SIP Calls indexes for spam detection queries
CREATE INDEX IF NOT EXISTS idx_sip_calls_created_at ON sip_calls(created_at);
CREATE INDEX IF NOT EXISTS idx_sip_calls_routed_to_ai ON sip_calls(routed_to_ai);
CREATE INDEX IF NOT EXISTS idx_sip_calls_spam_score ON sip_calls(spam_score);
CREATE INDEX IF NOT EXISTS idx_sip_calls_ai_composite ON sip_calls(routed_to_ai, created_at) WHERE routed_to_ai = true;

-- SIM Cards indexes for status monitoring
CREATE INDEX IF NOT EXISTS idx_sim_cards_status ON sim_cards(status);
CREATE INDEX IF NOT EXISTS idx_sim_cards_operator ON sim_cards(operator_name);
CREATE INDEX IF NOT EXISTS idx_sim_cards_balance ON sim_cards(balance);
CREATE INDEX IF NOT EXISTS idx_sim_cards_active_composite ON sim_cards(status, balance) WHERE status = 'active';

-- WhatsApp Validation Cache indexes
CREATE INDEX IF NOT EXISTS idx_whatsapp_cache_phone ON whatsapp_validation_cache(phone_number);
CREATE INDEX IF NOT EXISTS idx_whatsapp_cache_checked ON whatsapp_validation_cache(checked_at);

-- Voice Recognition Logs indexes
CREATE INDEX IF NOT EXISTS idx_voice_logs_call_id ON voice_recognition_logs(call_id);
CREATE INDEX IF NOT EXISTS idx_voice_logs_created ON voice_recognition_logs(created_at);
CREATE INDEX IF NOT EXISTS idx_voice_logs_result ON voice_recognition_logs(recognition_result);

-- AI Agent Interactions indexes
CREATE INDEX IF NOT EXISTS idx_ai_interactions_start ON ai_agent_interactions(start_time);
CREATE INDEX IF NOT EXISTS idx_ai_interactions_agent ON ai_agent_interactions(agent_id);
CREATE INDEX IF NOT EXISTS idx_ai_interactions_call ON ai_agent_interactions(call_id);

-- Blacklist indexes for fast lookup
CREATE INDEX IF NOT EXISTS idx_blacklist_number ON blacklist(number);
CREATE INDEX IF NOT EXISTS idx_blacklist_enabled ON blacklist(enabled) WHERE enabled = true;

-- Customer indexes
CREATE INDEX IF NOT EXISTS idx_customers_balance ON customers(balance);
CREATE INDEX IF NOT EXISTS idx_customers_created ON customers(created_at);

-- SIM Replacement Queue indexes
CREATE INDEX IF NOT EXISTS idx_sim_replacement_status ON sim_replacement_queue(status);
CREATE INDEX IF NOT EXISTS idx_sim_replacement_priority ON sim_replacement_queue(priority, created_at);

-- Partial indexes for common queries
CREATE INDEX IF NOT EXISTS idx_recent_calls ON call_detail_records(call_start_time) 
    WHERE call_start_time > (NOW() - INTERVAL '7 days');
CREATE INDEX IF NOT EXISTS idx_low_balance_sims ON sim_cards(balance, id) 
    WHERE balance < 10 AND status = 'active';
CREATE INDEX IF NOT EXISTS idx_pending_replacements ON sim_replacement_queue(priority, created_at) 
    WHERE status = 'pending';

-- Function-based index for hour extraction (for hourly analytics)
CREATE INDEX IF NOT EXISTS idx_cdr_hour ON call_detail_records(date_part('hour', call_start_time), call_start_time);

-- Gateway health monitoring
CREATE INDEX IF NOT EXISTS idx_gateway_health_timestamp ON gateway_health_logs(timestamp);
CREATE INDEX IF NOT EXISTS idx_gateway_health_gateway ON gateway_health_logs(gateway_id);

-- Revenue tracking
CREATE INDEX IF NOT EXISTS idx_revenue_date ON revenue_tracking(date);
CREATE INDEX IF NOT EXISTS idx_revenue_gateway ON revenue_tracking(gateway_id);

-- Analyze tables to update statistics
ANALYZE call_detail_records;
ANALYZE sip_calls;
ANALYZE sim_cards;
ANALYZE whatsapp_validation_cache;
ANALYZE voice_recognition_logs;
ANALYZE ai_agent_interactions;
ANALYZE blacklist;
ANALYZE customers;
ANALYZE sim_replacement_queue;
ANALYZE gateway_health_logs;
ANALYZE revenue_tracking;