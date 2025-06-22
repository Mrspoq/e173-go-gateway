-- Rate plans for per-second billing
CREATE TABLE IF NOT EXISTS rate_plans (
    id BIGSERIAL PRIMARY KEY,
    plan_name VARCHAR(100) NOT NULL,
    plan_code VARCHAR(20) UNIQUE NOT NULL,
    description TEXT,
    currency VARCHAR(3) DEFAULT 'USD',
    rate_per_minute DECIMAL(10, 6) NOT NULL,
    rate_per_second DECIMAL(10, 6) NOT NULL,
    minimum_billing_seconds INTEGER DEFAULT 1,
    connection_fee DECIMAL(10, 4) DEFAULT 0.0,
    daily_cap DECIMAL(10, 4),
    monthly_cap DECIMAL(10, 4),
    effective_from TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    effective_until TIMESTAMPTZ,
    is_active BOOLEAN DEFAULT TRUE,
    created_by BIGINT REFERENCES users(id),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Customer rate plan assignments
CREATE TABLE IF NOT EXISTS customer_rate_plans (
    id BIGSERIAL PRIMARY KEY,
    customer_id BIGINT NOT NULL REFERENCES customers(id) ON DELETE CASCADE,
    rate_plan_id BIGINT NOT NULL REFERENCES rate_plans(id) ON DELETE CASCADE,
    effective_from TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    effective_until TIMESTAMPTZ,
    is_active BOOLEAN DEFAULT TRUE,
    created_by BIGINT REFERENCES users(id),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    
    UNIQUE(customer_id, rate_plan_id, effective_from)
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_rate_plans_plan_code ON rate_plans(plan_code);
CREATE INDEX IF NOT EXISTS idx_rate_plans_is_active ON rate_plans(is_active);
CREATE INDEX IF NOT EXISTS idx_rate_plans_effective_from ON rate_plans(effective_from);
CREATE INDEX IF NOT EXISTS idx_customer_rate_plans_customer_id ON customer_rate_plans(customer_id);
CREATE INDEX IF NOT EXISTS idx_customer_rate_plans_rate_plan_id ON customer_rate_plans(rate_plan_id);
CREATE INDEX IF NOT EXISTS idx_customer_rate_plans_is_active ON customer_rate_plans(is_active);

-- Triggers
CREATE TRIGGER set_rate_plans_updated_at
BEFORE UPDATE ON rate_plans
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();

-- Default rate plan
INSERT INTO rate_plans (plan_name, plan_code, description, rate_per_minute, rate_per_second) 
VALUES ('Standard Plan', 'STD001', 'Default standard rate plan', 0.05, 0.000833)
ON CONFLICT (plan_code) DO NOTHING;
