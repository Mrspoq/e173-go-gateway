-- Create SIP accounts table
CREATE TABLE IF NOT EXISTS sip_accounts (
    id BIGSERIAL PRIMARY KEY,
    customer_id BIGINT NOT NULL REFERENCES customers(id) ON DELETE CASCADE,
    account_name VARCHAR(255) NOT NULL,
    username VARCHAR(100) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    domain VARCHAR(255) NOT NULL DEFAULT 'sip.e173gateway.com',
    extension VARCHAR(20),
    caller_id VARCHAR(50),
    caller_id_name VARCHAR(100),
    context VARCHAR(50) NOT NULL DEFAULT 'default',
    transport VARCHAR(10) NOT NULL DEFAULT 'UDP',
    nat_support BOOLEAN NOT NULL DEFAULT true,
    direct_media_support BOOLEAN NOT NULL DEFAULT false,
    encryption_enabled BOOLEAN NOT NULL DEFAULT false,
    codecs_allowed VARCHAR(500) NOT NULL DEFAULT 'g711u,g711a,g729,g722',
    max_concurrent_calls INTEGER NOT NULL DEFAULT 2,
    current_active_calls INTEGER NOT NULL DEFAULT 0,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    last_registered_ip VARCHAR(45),
    last_registered_at TIMESTAMP,
    last_call_at TIMESTAMP,
    total_calls BIGINT NOT NULL DEFAULT 0,
    total_minutes BIGINT NOT NULL DEFAULT 0,
    notes TEXT,
    created_by BIGINT REFERENCES users(id),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT check_status CHECK (status IN ('active', 'suspended', 'disabled', 'pending')),
    CONSTRAINT check_transport CHECK (transport IN ('UDP', 'TCP', 'TLS'))
);

-- Create indexes
CREATE INDEX idx_sip_accounts_customer_id ON sip_accounts(customer_id);
CREATE INDEX idx_sip_accounts_username ON sip_accounts(username);
CREATE INDEX idx_sip_accounts_status ON sip_accounts(status);
CREATE INDEX idx_sip_accounts_last_registered_at ON sip_accounts(last_registered_at);

-- Create SIP account permissions table
CREATE TABLE IF NOT EXISTS sip_account_permissions (
    id BIGSERIAL PRIMARY KEY,
    sip_account_id BIGINT NOT NULL REFERENCES sip_accounts(id) ON DELETE CASCADE,
    allow_international BOOLEAN NOT NULL DEFAULT false,
    allow_premium_numbers BOOLEAN NOT NULL DEFAULT false,
    allow_emergency_calls BOOLEAN NOT NULL DEFAULT true,
    allowed_countries TEXT, -- Comma-separated country codes
    blocked_countries TEXT, -- Comma-separated country codes
    allowed_prefixes TEXT,  -- Comma-separated prefixes
    blocked_prefixes TEXT,  -- Comma-separated prefixes
    time_restrictions JSONB, -- Time-based restrictions
    daily_call_limit INTEGER,
    daily_minute_limit INTEGER,
    monthly_call_limit INTEGER,
    monthly_minute_limit INTEGER,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(sip_account_id)
);

-- Create SIP registrations table
CREATE TABLE IF NOT EXISTS sip_registrations (
    id BIGSERIAL PRIMARY KEY,
    sip_account_id BIGINT NOT NULL REFERENCES sip_accounts(id) ON DELETE CASCADE,
    contact_uri VARCHAR(500) NOT NULL,
    source_ip VARCHAR(45) NOT NULL,
    source_port INTEGER NOT NULL,
    user_agent VARCHAR(255),
    expires_seconds INTEGER NOT NULL,
    registered_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expired_at TIMESTAMP NOT NULL,
    unregistered_at TIMESTAMP,
    is_active BOOLEAN NOT NULL DEFAULT true
);

-- Create indexes for registrations
CREATE INDEX idx_sip_registrations_account_id ON sip_registrations(sip_account_id);
CREATE INDEX idx_sip_registrations_is_active ON sip_registrations(is_active);
CREATE INDEX idx_sip_registrations_registered_at ON sip_registrations(registered_at);

-- Create SIP account usage statistics table
CREATE TABLE IF NOT EXISTS sip_account_usage (
    id BIGSERIAL PRIMARY KEY,
    sip_account_id BIGINT NOT NULL REFERENCES sip_accounts(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    total_calls INTEGER NOT NULL DEFAULT 0,
    successful_calls INTEGER NOT NULL DEFAULT 0,
    failed_calls INTEGER NOT NULL DEFAULT 0,
    total_minutes INTEGER NOT NULL DEFAULT 0,
    incoming_calls INTEGER NOT NULL DEFAULT 0,
    outgoing_calls INTEGER NOT NULL DEFAULT 0,
    international_calls INTEGER NOT NULL DEFAULT 0,
    average_call_duration INTEGER NOT NULL DEFAULT 0,
    peak_concurrent_calls INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(sip_account_id, date)
);

-- Create indexes for usage
CREATE INDEX idx_sip_account_usage_account_id ON sip_account_usage(sip_account_id);
CREATE INDEX idx_sip_account_usage_date ON sip_account_usage(date);

-- Create trigger to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_sip_accounts_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_sip_accounts_updated_at BEFORE UPDATE ON sip_accounts
    FOR EACH ROW EXECUTE FUNCTION update_sip_accounts_updated_at();

CREATE TRIGGER update_sip_account_permissions_updated_at BEFORE UPDATE ON sip_account_permissions
    FOR EACH ROW EXECUTE FUNCTION update_sip_accounts_updated_at();

-- Add some sample SIP accounts for testing
INSERT INTO sip_accounts (customer_id, account_name, username, password, extension, caller_id, caller_id_name)
SELECT 
    id,
    'Main SIP Account',
    LOWER(REPLACE(customer_code, '-', '')) || '001',
    MD5(RANDOM()::TEXT),
    '1001',
    COALESCE(phone, '+1234567890'),
    COALESCE(company_name, contact_person, customer_code)
FROM customers
WHERE account_status = 'active'
LIMIT 3;

-- Grant permissions
GRANT SELECT, INSERT, UPDATE, DELETE ON sip_accounts TO gateway_user;
GRANT SELECT, INSERT, UPDATE, DELETE ON sip_account_permissions TO gateway_user;
GRANT SELECT, INSERT, UPDATE, DELETE ON sip_registrations TO gateway_user;
GRANT SELECT, INSERT, UPDATE, DELETE ON sip_account_usage TO gateway_user;
GRANT USAGE ON SEQUENCE sip_accounts_id_seq TO gateway_user;
GRANT USAGE ON SEQUENCE sip_account_permissions_id_seq TO gateway_user;
GRANT USAGE ON SEQUENCE sip_registrations_id_seq TO gateway_user;
GRANT USAGE ON SEQUENCE sip_account_usage_id_seq TO gateway_user;