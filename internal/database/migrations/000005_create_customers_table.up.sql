-- Customers table for CRM management
CREATE TABLE IF NOT EXISTS customers (
    id BIGSERIAL PRIMARY KEY,
    customer_code VARCHAR(20) UNIQUE NOT NULL, -- e.g., CUST001, CUST002
    company_name VARCHAR(255),
    contact_person VARCHAR(255),
    email VARCHAR(255),
    phone VARCHAR(50),
    address TEXT,
    city VARCHAR(100),
    state VARCHAR(100),
    country VARCHAR(100),
    postal_code VARCHAR(20),
    billing_address TEXT,
    billing_city VARCHAR(100),
    billing_state VARCHAR(100),
    billing_country VARCHAR(100),
    billing_postal_code VARCHAR(20),
    account_status VARCHAR(20) DEFAULT 'active' CHECK (account_status IN ('active', 'suspended', 'terminated', 'pending')),
    credit_limit DECIMAL(12, 4) DEFAULT 0.0,
    current_balance DECIMAL(12, 4) DEFAULT 0.0,
    monthly_limit DECIMAL(12, 4),
    timezone VARCHAR(50) DEFAULT 'UTC',
    preferred_currency VARCHAR(3) DEFAULT 'USD',
    auto_recharge_enabled BOOLEAN DEFAULT FALSE,
    auto_recharge_threshold DECIMAL(12, 4),
    auto_recharge_amount DECIMAL(12, 4),
    notes TEXT,
    created_by BIGINT REFERENCES users(id),
    assigned_to BIGINT REFERENCES users(id), -- Account manager
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_customers_customer_code ON customers(customer_code);
CREATE INDEX IF NOT EXISTS idx_customers_company_name ON customers(company_name);
CREATE INDEX IF NOT EXISTS idx_customers_account_status ON customers(account_status);
CREATE INDEX IF NOT EXISTS idx_customers_assigned_to ON customers(assigned_to);
CREATE INDEX IF NOT EXISTS idx_customers_created_by ON customers(created_by);

-- Trigger to update 'updated_at' timestamp
CREATE TRIGGER set_customers_updated_at
BEFORE UPDATE ON customers
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();
