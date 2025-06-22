-- Payments and billing transactions table
CREATE TABLE IF NOT EXISTS payments (
    id BIGSERIAL PRIMARY KEY,
    customer_id BIGINT NOT NULL REFERENCES customers(id) ON DELETE CASCADE,
    payment_reference VARCHAR(100) UNIQUE NOT NULL,
    payment_type VARCHAR(20) NOT NULL CHECK (payment_type IN ('credit', 'debit', 'refund', 'adjustment', 'auto_recharge')),
    amount DECIMAL(12, 4) NOT NULL,
    currency VARCHAR(3) DEFAULT 'USD',
    description TEXT,
    payment_method VARCHAR(50), -- e.g., 'credit_card', 'bank_transfer', 'crypto', 'manual'
    transaction_id VARCHAR(255), -- External payment processor transaction ID
    gateway_response JSONB, -- Store full gateway response for auditing
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'completed', 'failed', 'refunded', 'cancelled')),
    processed_at TIMESTAMPTZ,
    processed_by BIGINT REFERENCES users(id),
    notes TEXT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_payments_customer_id ON payments(customer_id);
CREATE INDEX IF NOT EXISTS idx_payments_payment_reference ON payments(payment_reference);
CREATE INDEX IF NOT EXISTS idx_payments_payment_type ON payments(payment_type);
CREATE INDEX IF NOT EXISTS idx_payments_status ON payments(status);
CREATE INDEX IF NOT EXISTS idx_payments_processed_at ON payments(processed_at);

-- Trigger to update 'updated_at' timestamp
CREATE TRIGGER set_payments_updated_at
BEFORE UPDATE ON payments
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();
