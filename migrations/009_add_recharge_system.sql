-- Create recharge_codes table
CREATE TABLE IF NOT EXISTS recharge_codes (
    id SERIAL PRIMARY KEY,
    sim_card_id INTEGER NOT NULL REFERENCES sim_cards(id) ON DELETE CASCADE,
    code VARCHAR(50) NOT NULL,
    amount DECIMAL(10,2) NOT NULL,
    operator VARCHAR(50) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    used_at TIMESTAMP,
    expiry_date TIMESTAMP,
    response_message TEXT,
    created_by INTEGER NOT NULL REFERENCES users(id),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_recharge_code UNIQUE (code, operator)
);

-- Create recharge_batches table
CREATE TABLE IF NOT EXISTS recharge_batches (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    total_codes INTEGER NOT NULL DEFAULT 0,
    used_codes INTEGER NOT NULL DEFAULT 0,
    total_amount DECIMAL(10,2) NOT NULL DEFAULT 0,
    status VARCHAR(20) NOT NULL DEFAULT 'draft',
    created_by INTEGER NOT NULL REFERENCES users(id),
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create recharge_history table
CREATE TABLE IF NOT EXISTS recharge_history (
    id SERIAL PRIMARY KEY,
    sim_card_id INTEGER NOT NULL REFERENCES sim_cards(id) ON DELETE CASCADE,
    recharge_code_id INTEGER REFERENCES recharge_codes(id),
    batch_id INTEGER REFERENCES recharge_batches(id),
    phone_number VARCHAR(20) NOT NULL,
    amount DECIMAL(10,2) NOT NULL,
    balance_before DECIMAL(10,2),
    balance_after DECIMAL(10,2),
    method VARCHAR(20) NOT NULL DEFAULT 'ussd',
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    error_message TEXT,
    attempts INTEGER NOT NULL DEFAULT 1,
    processed_by INTEGER NOT NULL REFERENCES users(id),
    processed_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for performance
CREATE INDEX idx_recharge_codes_sim_card_id ON recharge_codes(sim_card_id);
CREATE INDEX idx_recharge_codes_status ON recharge_codes(status);
CREATE INDEX idx_recharge_codes_operator ON recharge_codes(operator);
CREATE INDEX idx_recharge_codes_expiry ON recharge_codes(expiry_date);

CREATE INDEX idx_recharge_batches_status ON recharge_batches(status);
CREATE INDEX idx_recharge_batches_created_by ON recharge_batches(created_by);

CREATE INDEX idx_recharge_history_sim_card_id ON recharge_history(sim_card_id);
CREATE INDEX idx_recharge_history_batch_id ON recharge_history(batch_id);
CREATE INDEX idx_recharge_history_status ON recharge_history(status);
CREATE INDEX idx_recharge_history_processed_at ON recharge_history(processed_at);

-- Create trigger to update updated_at
CREATE OR REPLACE FUNCTION update_recharge_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_recharge_codes_updated_at
    BEFORE UPDATE ON recharge_codes
    FOR EACH ROW
    EXECUTE FUNCTION update_recharge_updated_at();

CREATE TRIGGER update_recharge_batches_updated_at
    BEFORE UPDATE ON recharge_batches
    FOR EACH ROW
    EXECUTE FUNCTION update_recharge_updated_at();

-- Add recharge-related columns to sim_cards if not exists
ALTER TABLE sim_cards ADD COLUMN IF NOT EXISTS auto_recharge_enabled BOOLEAN DEFAULT false;
ALTER TABLE sim_cards ADD COLUMN IF NOT EXISTS auto_recharge_threshold DECIMAL(10,2) DEFAULT 5.00;
ALTER TABLE sim_cards ADD COLUMN IF NOT EXISTS auto_recharge_amount DECIMAL(10,2) DEFAULT 10.00;
ALTER TABLE sim_cards ADD COLUMN IF NOT EXISTS last_recharge_at TIMESTAMP;
ALTER TABLE sim_cards ADD COLUMN IF NOT EXISTS total_recharged DECIMAL(10,2) DEFAULT 0;

-- Create view for recharge statistics
CREATE OR REPLACE VIEW recharge_stats AS
SELECT 
    s.id as sim_card_id,
    s.phone_number,
    s.operator,
    COUNT(rh.id) as total_recharges,
    SUM(CASE WHEN rh.status = 'success' THEN 1 ELSE 0 END) as successful_recharges,
    SUM(CASE WHEN rh.status = 'failed' THEN 1 ELSE 0 END) as failed_recharges,
    SUM(CASE WHEN rh.status = 'success' THEN rh.amount ELSE 0 END) as total_recharged_amount,
    MAX(rh.processed_at) as last_recharge_date
FROM sim_cards s
LEFT JOIN recharge_history rh ON s.id = rh.sim_card_id
GROUP BY s.id, s.phone_number, s.operator;

-- Sample data for testing
INSERT INTO recharge_batches (name, description, created_by) 
VALUES ('Test Batch', 'Initial test batch for development', 1);