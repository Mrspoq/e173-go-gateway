-- Database Updates for Multi-Gateway Support
-- Execute these after waking up to extend the schema

-- Add SIP call tracking
CREATE TABLE IF NOT EXISTS sip_calls (
    id SERIAL PRIMARY KEY,
    call_id VARCHAR(255) UNIQUE NOT NULL,
    caller_number VARCHAR(50) NOT NULL,
    destination_number VARCHAR(50) NOT NULL,
    gateway_id INTEGER REFERENCES gateways(id),
    filter_result JSONB,
    routed_to_ai BOOLEAN DEFAULT FALSE,
    ai_session_id VARCHAR(255),
    billing_seconds INTEGER DEFAULT 0,
    spam_score DECIMAL(3,2),
    operator_detected VARCHAR(50),
    sticky_routing BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW(),
    ended_at TIMESTAMP,
    
    INDEX idx_caller_number (caller_number),
    INDEX idx_destination_number (destination_number),
    INDEX idx_call_time (created_at),
    INDEX idx_gateway_id (gateway_id)
);

-- Extend gateways table for SIP support
ALTER TABLE gateways ADD COLUMN IF NOT EXISTS sip_endpoint VARCHAR(255);
ALTER TABLE gateways ADD COLUMN IF NOT EXISTS sip_port INTEGER DEFAULT 5060;
ALTER TABLE gateways ADD COLUMN IF NOT EXISTS health_status VARCHAR(50) DEFAULT 'healthy';
ALTER TABLE gateways ADD COLUMN IF NOT EXISTS health_last_check TIMESTAMP DEFAULT NOW();
ALTER TABLE gateways ADD COLUMN IF NOT EXISTS concurrent_calls INTEGER DEFAULT 0;
ALTER TABLE gateways ADD COLUMN IF NOT EXISTS max_concurrent_calls INTEGER DEFAULT 50;
ALTER TABLE gateways ADD COLUMN IF NOT EXISTS region VARCHAR(100);
ALTER TABLE gateways ADD COLUMN IF NOT EXISTS operator_preferences JSONB;

-- Call history analysis table
CREATE TABLE IF NOT EXISTS call_patterns (
    id SERIAL PRIMARY KEY,
    phone_number VARCHAR(50) NOT NULL,
    total_calls INTEGER DEFAULT 0,
    answered_calls INTEGER DEFAULT 0,
    short_calls INTEGER DEFAULT 0, -- calls < 10 seconds
    avg_call_duration DECIMAL(8,2) DEFAULT 0,
    spam_score DECIMAL(3,2) DEFAULT 0,
    last_call_time TIMESTAMP,
    pattern_updated TIMESTAMP DEFAULT NOW(),
    
    UNIQUE(phone_number),
    INDEX idx_phone_number (phone_number),
    INDEX idx_spam_score (spam_score),
    INDEX idx_last_call (last_call_time)
);

-- Blacklist enhancements
ALTER TABLE blacklist ADD COLUMN IF NOT EXISTS reason_code VARCHAR(50);
ALTER TABLE blacklist ADD COLUMN IF NOT EXISTS auto_added BOOLEAN DEFAULT FALSE;
ALTER TABLE blacklist ADD COLUMN IF NOT EXISTS expires_at TIMESTAMP;
ALTER TABLE blacklist ADD COLUMN IF NOT EXISTS confidence_score DECIMAL(3,2);

-- Operator routing rules
CREATE TABLE IF NOT EXISTS operator_routing_rules (
    id SERIAL PRIMARY KEY,
    operator_name VARCHAR(100) NOT NULL,
    prefix_pattern VARCHAR(20) NOT NULL,
    preferred_gateway_id INTEGER REFERENCES gateways(id),
    backup_gateway_ids INTEGER[], -- Array of gateway IDs
    cost_per_minute DECIMAL(8,4),
    quality_score DECIMAL(3,2), -- 0.0 to 1.0
    active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW(),
    
    INDEX idx_operator_name (operator_name),
    INDEX idx_prefix_pattern (prefix_pattern)
);

-- AI voice agents table
CREATE TABLE IF NOT EXISTS ai_voice_agents (
    id SERIAL PRIMARY KEY,
    agent_name VARCHAR(100) NOT NULL,
    endpoint_url VARCHAR(255) NOT NULL,
    api_key VARCHAR(255),
    model_type VARCHAR(50), -- 'elevenlabs', 'openai', 'custom'
    voice_id VARCHAR(100),
    language VARCHAR(10) DEFAULT 'en',
    concurrent_limit INTEGER DEFAULT 10,
    current_sessions INTEGER DEFAULT 0,
    cost_per_minute DECIMAL(8,4),
    quality_rating DECIMAL(3,2),
    active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW()
);

-- AI session tracking
CREATE TABLE IF NOT EXISTS ai_voice_sessions (
    id SERIAL PRIMARY KEY,
    session_id VARCHAR(255) UNIQUE NOT NULL,
    agent_id INTEGER REFERENCES ai_voice_agents(id),
    sip_call_id INTEGER REFERENCES sip_calls(id),
    caller_number VARCHAR(50),
    session_duration INTEGER DEFAULT 0, -- seconds
    tokens_used INTEGER DEFAULT 0,
    revenue_generated DECIMAL(10,4) DEFAULT 0,
    conversation_log JSONB,
    started_at TIMESTAMP DEFAULT NOW(),
    ended_at TIMESTAMP,
    
    INDEX idx_session_id (session_id),
    INDEX idx_agent_id (agent_id),
    INDEX idx_started_at (started_at)
);

-- WhatsApp validation cache
CREATE TABLE IF NOT EXISTS whatsapp_validation_cache (
    phone_number VARCHAR(50) PRIMARY KEY,
    has_whatsapp BOOLEAN NOT NULL,
    confidence_score DECIMAL(3,2),
    checked_at TIMESTAMP DEFAULT NOW(),
    expires_at TIMESTAMP DEFAULT (NOW() + INTERVAL '7 days'),
    
    INDEX idx_checked_at (checked_at),
    INDEX idx_expires_at (expires_at)
);

-- Gateway health monitoring
CREATE TABLE IF NOT EXISTS gateway_health_logs (
    id SERIAL PRIMARY KEY,
    gateway_id INTEGER REFERENCES gateways(id),
    status VARCHAR(50) NOT NULL, -- 'healthy', 'degraded', 'offline'
    response_time_ms INTEGER,
    concurrent_calls INTEGER,
    cpu_usage DECIMAL(5,2),
    memory_usage DECIMAL(5,2),
    error_rate DECIMAL(5,2),
    checked_at TIMESTAMP DEFAULT NOW(),
    
    INDEX idx_gateway_id (gateway_id),
    INDEX idx_checked_at (checked_at),
    INDEX idx_status (status)
);

-- Revenue tracking
CREATE TABLE IF NOT EXISTS revenue_tracking (
    id SERIAL PRIMARY KEY,
    call_id VARCHAR(255),
    revenue_source VARCHAR(50), -- 'legitimate_call', 'spam_monetization', 'ai_interaction'
    amount DECIMAL(10,4) NOT NULL,
    currency VARCHAR(3) DEFAULT 'USD',
    customer_id INTEGER REFERENCES customers(id),
    gateway_id INTEGER REFERENCES gateways(id),
    ai_agent_id INTEGER REFERENCES ai_voice_agents(id),
    billing_date DATE DEFAULT CURRENT_DATE,
    created_at TIMESTAMP DEFAULT NOW(),
    
    INDEX idx_billing_date (billing_date),
    INDEX idx_revenue_source (revenue_source),
    INDEX idx_customer_id (customer_id)
);

-- Insert sample data for testing
INSERT INTO operator_routing_rules (operator_name, prefix_pattern, preferred_gateway_id, cost_per_minute, quality_score) VALUES
('MTN Nigeria', '+234803%', 1, 0.045, 0.95),
('MTN Nigeria', '+234806%', 1, 0.045, 0.95),
('Airtel Nigeria', '+234802%', 2, 0.042, 0.92),
('Glo Nigeria', '+234805%', 3, 0.040, 0.88),
('9mobile Nigeria', '+234809%', 4, 0.048, 0.85);

INSERT INTO ai_voice_agents (agent_name, endpoint_url, model_type, voice_id, cost_per_minute, quality_rating) VALUES
('SpamBot Interceptor 1', 'http://ai-service:8080/voice/session1', 'elevenlabs', 'voice_001', 0.02, 0.90),
('SpamBot Interceptor 2', 'http://ai-service:8080/voice/session2', 'elevenlabs', 'voice_002', 0.02, 0.88),
('Customer Service AI', 'http://ai-service:8080/voice/customer', 'openai', 'voice_customer', 0.05, 0.95);

-- Create indexes for performance
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_sip_calls_performance ON sip_calls (created_at, gateway_id, routed_to_ai);
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_call_patterns_analysis ON call_patterns (spam_score DESC, total_calls DESC);
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_revenue_daily ON revenue_tracking (billing_date, revenue_source);

-- Create views for reporting
CREATE OR REPLACE VIEW daily_revenue_summary AS
SELECT 
    billing_date,
    revenue_source,
    COUNT(*) as transaction_count,
    SUM(amount) as total_revenue,
    AVG(amount) as avg_revenue_per_call
FROM revenue_tracking 
GROUP BY billing_date, revenue_source
ORDER BY billing_date DESC;

CREATE OR REPLACE VIEW gateway_performance_summary AS
SELECT 
    g.id,
    g.name,
    g.health_status,
    COUNT(sc.id) as total_calls_today,
    AVG(CASE WHEN sc.routed_to_ai THEN 0 ELSE 1 END) as legitimate_call_ratio,
    SUM(sc.billing_seconds) as total_minutes_today
FROM gateways g
LEFT JOIN sip_calls sc ON g.id = sc.gateway_id 
    AND sc.created_at >= CURRENT_DATE
GROUP BY g.id, g.name, g.health_status;

-- Comments for future development
-- TODO: Add partitioning for sip_calls table by date
-- TODO: Add triggers for automatic spam score calculation
-- TODO: Add stored procedures for gateway failover logic
-- TODO: Add materialized views for real-time dashboard stats

COMMIT;
