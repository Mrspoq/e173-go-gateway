-- Database Updates for Multi-Gateway Support (UUID Compatible)

BEGIN;

-- Add SIP call tracking (with UUID foreign key)
CREATE TABLE IF NOT EXISTS sip_calls (
    id SERIAL PRIMARY KEY,
    call_id VARCHAR(255) UNIQUE NOT NULL,
    caller_number VARCHAR(50) NOT NULL,
    destination_number VARCHAR(50) NOT NULL,
    gateway_id UUID REFERENCES gateways(id),
    filter_result JSONB,
    routed_to_ai BOOLEAN DEFAULT FALSE,
    ai_session_id VARCHAR(255),
    billing_seconds INTEGER DEFAULT 0,
    spam_score DECIMAL(3,2),
    operator_detected VARCHAR(50),
    sticky_routing BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW(),
    ended_at TIMESTAMP
);

-- Create indexes for sip_calls
CREATE INDEX IF NOT EXISTS idx_sip_calls_caller ON sip_calls (caller_number);
CREATE INDEX IF NOT EXISTS idx_sip_calls_destination ON sip_calls (destination_number);
CREATE INDEX IF NOT EXISTS idx_sip_calls_time ON sip_calls (created_at);
CREATE INDEX IF NOT EXISTS idx_sip_calls_gateway ON sip_calls (gateway_id);

-- Call history analysis table
CREATE TABLE IF NOT EXISTS call_patterns (
    id SERIAL PRIMARY KEY,
    phone_number VARCHAR(50) NOT NULL,
    total_calls INTEGER DEFAULT 0,
    answered_calls INTEGER DEFAULT 0,
    short_calls INTEGER DEFAULT 0,
    avg_call_duration DECIMAL(8,2) DEFAULT 0,
    spam_score DECIMAL(3,2) DEFAULT 0,
    last_call_time TIMESTAMP,
    pattern_updated TIMESTAMP DEFAULT NOW(),
    UNIQUE(phone_number)
);

-- Create indexes for call_patterns
CREATE INDEX IF NOT EXISTS idx_call_patterns_phone ON call_patterns (phone_number);
CREATE INDEX IF NOT EXISTS idx_call_patterns_spam ON call_patterns (spam_score);
CREATE INDEX IF NOT EXISTS idx_call_patterns_lastcall ON call_patterns (last_call_time);

-- Operator routing rules (with UUID foreign key)
CREATE TABLE IF NOT EXISTS operator_routing_rules (
    id SERIAL PRIMARY KEY,
    operator_name VARCHAR(100) NOT NULL,
    prefix_pattern VARCHAR(20) NOT NULL,
    preferred_gateway_id UUID REFERENCES gateways(id),
    backup_gateway_ids UUID[],
    cost_per_minute DECIMAL(8,4),
    quality_score DECIMAL(3,2),
    active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Create indexes for operator_routing_rules
CREATE INDEX IF NOT EXISTS idx_operator_rules_name ON operator_routing_rules (operator_name);
CREATE INDEX IF NOT EXISTS idx_operator_rules_prefix ON operator_routing_rules (prefix_pattern);

-- AI voice agents table
CREATE TABLE IF NOT EXISTS ai_voice_agents (
    id SERIAL PRIMARY KEY,
    agent_name VARCHAR(100) NOT NULL,
    endpoint_url VARCHAR(255) NOT NULL,
    api_key VARCHAR(255),
    model_type VARCHAR(50),
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
    session_duration INTEGER DEFAULT 0,
    tokens_used INTEGER DEFAULT 0,
    revenue_generated DECIMAL(10,4) DEFAULT 0,
    conversation_log JSONB,
    started_at TIMESTAMP DEFAULT NOW(),
    ended_at TIMESTAMP
);

-- Create indexes for ai_voice_sessions
CREATE INDEX IF NOT EXISTS idx_ai_sessions_id ON ai_voice_sessions (session_id);
CREATE INDEX IF NOT EXISTS idx_ai_sessions_agent ON ai_voice_sessions (agent_id);
CREATE INDEX IF NOT EXISTS idx_ai_sessions_start ON ai_voice_sessions (started_at);

-- WhatsApp validation cache
CREATE TABLE IF NOT EXISTS whatsapp_validation_cache (
    phone_number VARCHAR(50) PRIMARY KEY,
    has_whatsapp BOOLEAN NOT NULL,
    confidence_score DECIMAL(3,2),
    checked_at TIMESTAMP DEFAULT NOW(),
    expires_at TIMESTAMP DEFAULT (NOW() + INTERVAL '7 days')
);

-- Create indexes for whatsapp_validation_cache
CREATE INDEX IF NOT EXISTS idx_whatsapp_checked ON whatsapp_validation_cache (checked_at);
CREATE INDEX IF NOT EXISTS idx_whatsapp_expires ON whatsapp_validation_cache (expires_at);

-- Gateway health monitoring (with UUID foreign key)
CREATE TABLE IF NOT EXISTS gateway_health_logs (
    id SERIAL PRIMARY KEY,
    gateway_id UUID REFERENCES gateways(id),
    status VARCHAR(50) NOT NULL,
    response_time_ms INTEGER,
    concurrent_calls INTEGER,
    cpu_usage DECIMAL(5,2),
    memory_usage DECIMAL(5,2),
    error_rate DECIMAL(5,2),
    checked_at TIMESTAMP DEFAULT NOW()
);

-- Create indexes for gateway_health_logs
CREATE INDEX IF NOT EXISTS idx_health_gateway ON gateway_health_logs (gateway_id);
CREATE INDEX IF NOT EXISTS idx_health_checked ON gateway_health_logs (checked_at);
CREATE INDEX IF NOT EXISTS idx_health_status ON gateway_health_logs (status);

-- Revenue tracking (with UUID foreign key)
CREATE TABLE IF NOT EXISTS revenue_tracking (
    id SERIAL PRIMARY KEY,
    call_id VARCHAR(255),
    revenue_source VARCHAR(50),
    amount DECIMAL(10,4) NOT NULL,
    currency VARCHAR(3) DEFAULT 'USD',
    customer_id INTEGER REFERENCES customers(id),
    gateway_id UUID REFERENCES gateways(id),
    ai_agent_id INTEGER REFERENCES ai_voice_agents(id),
    billing_date DATE DEFAULT CURRENT_DATE,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Create indexes for revenue_tracking
CREATE INDEX IF NOT EXISTS idx_revenue_date ON revenue_tracking (billing_date);
CREATE INDEX IF NOT EXISTS idx_revenue_source ON revenue_tracking (revenue_source);
CREATE INDEX IF NOT EXISTS idx_revenue_customer ON revenue_tracking (customer_id);

-- Get the first gateway ID for sample data
DO $$
DECLARE 
    gateway_uuid UUID;
BEGIN
    SELECT id INTO gateway_uuid FROM gateways LIMIT 1;
    
    IF gateway_uuid IS NOT NULL THEN
        -- Insert sample data for testing
        INSERT INTO operator_routing_rules (operator_name, prefix_pattern, preferred_gateway_id, cost_per_minute, quality_score) VALUES
        ('MTN Nigeria', '+234803%', gateway_uuid, 0.045, 0.95),
        ('MTN Nigeria', '+234806%', gateway_uuid, 0.045, 0.95),
        ('Airtel Nigeria', '+234802%', gateway_uuid, 0.042, 0.92),
        ('Glo Nigeria', '+234805%', gateway_uuid, 0.040, 0.88),
        ('9mobile Nigeria', '+234809%', gateway_uuid, 0.048, 0.85)
        ON CONFLICT DO NOTHING;
    END IF;
END $$;

INSERT INTO ai_voice_agents (agent_name, endpoint_url, model_type, voice_id, cost_per_minute, quality_rating) VALUES
('SpamBot Interceptor 1', 'http://ai-service:8080/voice/session1', 'elevenlabs', 'voice_001', 0.02, 0.90),
('SpamBot Interceptor 2', 'http://ai-service:8080/voice/session2', 'elevenlabs', 'voice_002', 0.02, 0.88),
('Customer Service AI', 'http://ai-service:8080/voice/customer', 'openai', 'voice_customer', 0.05, 0.95)
ON CONFLICT DO NOTHING;

COMMIT;

-- Show completion message
SELECT 'Database schema updated successfully for multi-gateway SIP platform!' as message;
