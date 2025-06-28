-- Create voice recognition logs table
CREATE TABLE IF NOT EXISTS voice_recognition_logs (
    id SERIAL PRIMARY KEY,
    call_id VARCHAR(100) NOT NULL,
    category VARCHAR(50) NOT NULL,
    action VARCHAR(50) NOT NULL,
    confidence DECIMAL(3,2) NOT NULL,
    reason TEXT,
    risk_score DECIMAL(3,2),
    keywords TEXT[],
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Create indexes for voice recognition logs
CREATE INDEX IF NOT EXISTS idx_voice_logs_call_id ON voice_recognition_logs (call_id);
CREATE INDEX IF NOT EXISTS idx_voice_logs_category ON voice_recognition_logs (category);
CREATE INDEX IF NOT EXISTS idx_voice_logs_created_at ON voice_recognition_logs (created_at);

-- Create call reviews table for manual review
CREATE TABLE IF NOT EXISTS call_reviews (
    id SERIAL PRIMARY KEY,
    call_id VARCHAR(100) NOT NULL UNIQUE,
    category VARCHAR(50) NOT NULL,
    confidence DECIMAL(3,2) NOT NULL,
    reason TEXT,
    risk_score DECIMAL(3,2),
    status VARCHAR(20) NOT NULL DEFAULT 'pending', -- pending, reviewed, escalated
    reviewer_id INTEGER,
    review_notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    reviewed_at TIMESTAMP WITH TIME ZONE
);

-- Create indexes for call reviews
CREATE INDEX IF NOT EXISTS idx_call_reviews_status ON call_reviews (status);
CREATE INDEX IF NOT EXISTS idx_call_reviews_call_id ON call_reviews (call_id);

-- Create SIM replacement queue table
CREATE TABLE IF NOT EXISTS sim_replacement_queue (
    id SERIAL PRIMARY KEY,
    sim_id VARCHAR(100) NOT NULL UNIQUE,
    reason TEXT NOT NULL,
    priority VARCHAR(20) NOT NULL DEFAULT 'normal', -- low, normal, high, urgent
    scheduled_at TIMESTAMP WITH TIME ZONE NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending', -- pending, in_progress, completed, cancelled
    assigned_to INTEGER,
    completed_at TIMESTAMP WITH TIME ZONE,
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Create indexes for SIM replacement queue
CREATE INDEX IF NOT EXISTS idx_sim_replacement_status ON sim_replacement_queue (status);
CREATE INDEX IF NOT EXISTS idx_sim_replacement_priority ON sim_replacement_queue (priority);

-- Create AI agent interactions table
CREATE TABLE IF NOT EXISTS ai_agent_interactions (
    id SERIAL PRIMARY KEY,
    call_id VARCHAR(100) NOT NULL,
    agent_id VARCHAR(50) NOT NULL,
    agent_name VARCHAR(100),
    strategy VARCHAR(50),
    start_time TIMESTAMP WITH TIME ZONE NOT NULL,
    end_time TIMESTAMP WITH TIME ZONE,
    duration_seconds INTEGER,
    transcript JSONB,
    outcome VARCHAR(50), -- time_wasted, info_collected, call_ended, error
    metadata JSONB,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Create indexes for AI agent interactions
CREATE INDEX IF NOT EXISTS idx_ai_interactions_call_id ON ai_agent_interactions (call_id);
CREATE INDEX IF NOT EXISTS idx_ai_interactions_agent_id ON ai_agent_interactions (agent_id);
CREATE INDEX IF NOT EXISTS idx_ai_interactions_created_at ON ai_agent_interactions (created_at);

-- Create voice transcripts table
CREATE TABLE IF NOT EXISTS voice_transcripts (
    id SERIAL PRIMARY KEY,
    call_id VARCHAR(100) NOT NULL,
    direction VARCHAR(20) NOT NULL, -- incoming, outgoing
    text TEXT NOT NULL,
    language VARCHAR(10),
    duration_seconds DECIMAL(10,2),
    confidence DECIMAL(3,2),
    segments JSONB,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Create indexes for voice transcripts
CREATE INDEX IF NOT EXISTS idx_transcripts_call_id ON voice_transcripts (call_id);
CREATE INDEX IF NOT EXISTS idx_transcripts_created_at ON voice_transcripts (created_at);

-- Add voice recognition fields to sip_calls table if it doesn't exist
DO $$ 
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns 
                   WHERE table_name='sip_calls' AND column_name='voice_analyzed') THEN
        ALTER TABLE sip_calls ADD COLUMN voice_analyzed BOOLEAN DEFAULT FALSE;
        ALTER TABLE sip_calls ADD COLUMN voice_category VARCHAR(50);
        ALTER TABLE sip_calls ADD COLUMN voice_action VARCHAR(50);
        ALTER TABLE sip_calls ADD COLUMN routed_to_ai BOOLEAN DEFAULT FALSE;
        ALTER TABLE sip_calls ADD COLUMN ai_agent_id VARCHAR(50);
    END IF;
END $$;

-- Add last_issue and flagged_at to sim_cards if they don't exist
DO $$ 
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns 
                   WHERE table_name='sim_cards' AND column_name='last_issue') THEN
        ALTER TABLE sim_cards ADD COLUMN last_issue TEXT;
        ALTER TABLE sim_cards ADD COLUMN flagged_at TIMESTAMP WITH TIME ZONE;
        ALTER TABLE sim_cards ADD COLUMN action_taken VARCHAR(100);
    END IF;
END $$;

-- Comments for documentation
COMMENT ON TABLE voice_recognition_logs IS 'Logs all voice recognition analysis results';
COMMENT ON TABLE call_reviews IS 'Calls flagged for manual review';
COMMENT ON TABLE sim_replacement_queue IS 'Queue for SIM cards that need replacement';
COMMENT ON TABLE ai_agent_interactions IS 'Records AI agent handling of spam calls';
COMMENT ON TABLE voice_transcripts IS 'Stores speech-to-text transcripts';

COMMENT ON COLUMN voice_recognition_logs.category IS 'SPAM_ROBOCALL, SIM_BLOCKED, VOICEMAIL, NORMAL_CALL, etc.';
COMMENT ON COLUMN voice_recognition_logs.action IS 'ROUTE_TO_AI, FLAG_SIM, NORMAL_ROUTING, BLOCK_CALL, etc.';
COMMENT ON COLUMN ai_agent_interactions.strategy IS 'TIME_WASTER, INFO_COLLECTOR, CONFUSER, etc.';