CREATE TABLE IF NOT EXISTS call_detail_records (
    id BIGSERIAL PRIMARY KEY,
    sim_card_id BIGINT REFERENCES sim_cards(id) ON DELETE SET NULL,
    modem_id BIGINT REFERENCES modems(id) ON DELETE SET NULL,
    asterisk_unique_id VARCHAR(128) UNIQUE NOT NULL,
    call_direction VARCHAR(10) NOT NULL CHECK (call_direction IN ('inbound', 'outbound')),
    source_number VARCHAR(50),
    destination_number VARCHAR(50) NOT NULL,
    call_start_time TIMESTAMPTZ NOT NULL,
    call_answer_time TIMESTAMPTZ,
    call_end_time TIMESTAMPTZ NOT NULL,
    duration_seconds INTEGER, -- Calculated: call_end_time - call_answer_time, or 0 if not answered
    billable_duration_seconds INTEGER, -- Could be different from duration_seconds
    disposition VARCHAR(50) NOT NULL, -- e.g., ANSWERED, NO ANSWER, BUSY, FAILED
    hangup_cause VARCHAR(100),
    cost_per_minute DECIMAL(10, 4),
    total_cost DECIMAL(10, 4),
    customer_id BIGINT, -- Will reference customers(id) later when customers table is created
    is_spam BOOLEAN DEFAULT FALSE,
    spam_reason VARCHAR(255),
    recorded_audio_path VARCHAR(255),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_cdr_sim_card_id ON call_detail_records(sim_card_id);
CREATE INDEX IF NOT EXISTS idx_cdr_modem_id ON call_detail_records(modem_id);
CREATE INDEX IF NOT EXISTS idx_cdr_asterisk_unique_id ON call_detail_records(asterisk_unique_id);
CREATE INDEX IF NOT EXISTS idx_cdr_call_start_time ON call_detail_records(call_start_time);
CREATE INDEX IF NOT EXISTS idx_cdr_destination_number ON call_detail_records(destination_number);
CREATE INDEX IF NOT EXISTS idx_cdr_customer_id ON call_detail_records(customer_id);
CREATE INDEX IF NOT EXISTS idx_cdr_is_spam ON call_detail_records(is_spam);

-- Trigger to update 'updated_at' timestamp using the existing function
CREATE TRIGGER set_cdr_updated_at
BEFORE UPDATE ON call_detail_records
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();