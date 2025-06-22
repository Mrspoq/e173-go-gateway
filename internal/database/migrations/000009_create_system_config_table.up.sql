-- System configuration and settings table
CREATE TABLE IF NOT EXISTS system_config (
    id BIGSERIAL PRIMARY KEY,
    config_key VARCHAR(100) UNIQUE NOT NULL,
    config_value TEXT,
    config_type VARCHAR(20) DEFAULT 'string' CHECK (config_type IN ('string', 'integer', 'float', 'boolean', 'json')),
    description TEXT,
    is_encrypted BOOLEAN DEFAULT FALSE,
    is_system BOOLEAN DEFAULT FALSE, -- True for system-level configs, false for user-configurable
    category VARCHAR(50), -- e.g., 'billing', 'routing', 'security', 'notifications'
    updated_by BIGINT REFERENCES users(id),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- User sessions for authentication tracking
CREATE TABLE IF NOT EXISTS user_sessions (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    session_token VARCHAR(255) UNIQUE NOT NULL,
    ip_address INET,
    user_agent TEXT,
    expires_at TIMESTAMPTZ NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    last_activity_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Audit log for tracking user actions
CREATE TABLE IF NOT EXISTS audit_logs (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(id) ON DELETE SET NULL,
    action VARCHAR(100) NOT NULL, -- e.g., 'login', 'create_customer', 'delete_modem'
    entity_type VARCHAR(50), -- e.g., 'customer', 'modem', 'sim_card'
    entity_id BIGINT, -- ID of the affected entity
    old_values JSONB, -- Previous values before change
    new_values JSONB, -- New values after change
    ip_address INET,
    user_agent TEXT,
    success BOOLEAN DEFAULT TRUE,
    error_message TEXT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Notification templates and alerts
CREATE TABLE IF NOT EXISTS notification_templates (
    id BIGSERIAL PRIMARY KEY,
    template_name VARCHAR(100) UNIQUE NOT NULL,
    template_type VARCHAR(50) NOT NULL, -- e.g., 'email', 'sms', 'webhook', 'internal'
    subject_template TEXT,
    body_template TEXT NOT NULL,
    variables JSONB, -- Available template variables
    is_active BOOLEAN DEFAULT TRUE,
    created_by BIGINT REFERENCES users(id),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- User notifications and alerts
CREATE TABLE IF NOT EXISTS user_notifications (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    notification_type VARCHAR(50) NOT NULL, -- e.g., 'low_balance', 'modem_offline', 'spam_detected'
    title VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    priority VARCHAR(20) DEFAULT 'normal' CHECK (priority IN ('low', 'normal', 'high', 'critical')),
    is_read BOOLEAN DEFAULT FALSE,
    read_at TIMESTAMPTZ,
    action_url VARCHAR(500), -- Optional URL for action button
    metadata JSONB, -- Additional context data
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_system_config_config_key ON system_config(config_key);
CREATE INDEX IF NOT EXISTS idx_system_config_category ON system_config(category);
CREATE INDEX IF NOT EXISTS idx_user_sessions_user_id ON user_sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_user_sessions_session_token ON user_sessions(session_token);
CREATE INDEX IF NOT EXISTS idx_user_sessions_expires_at ON user_sessions(expires_at);
CREATE INDEX IF NOT EXISTS idx_audit_logs_user_id ON audit_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_action ON audit_logs(action);
CREATE INDEX IF NOT EXISTS idx_audit_logs_entity_type ON audit_logs(entity_type);
CREATE INDEX IF NOT EXISTS idx_audit_logs_created_at ON audit_logs(created_at);
CREATE INDEX IF NOT EXISTS idx_notification_templates_template_name ON notification_templates(template_name);
CREATE INDEX IF NOT EXISTS idx_user_notifications_user_id ON user_notifications(user_id);
CREATE INDEX IF NOT EXISTS idx_user_notifications_is_read ON user_notifications(is_read);
CREATE INDEX IF NOT EXISTS idx_user_notifications_created_at ON user_notifications(created_at);

-- Triggers
CREATE TRIGGER set_system_config_updated_at
BEFORE UPDATE ON system_config
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();

CREATE TRIGGER set_notification_templates_updated_at
BEFORE UPDATE ON notification_templates
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();

-- Default system configuration
INSERT INTO system_config (config_key, config_value, config_type, description, category, is_system) VALUES
    ('system_name', 'E173 Gateway Management System', 'string', 'System display name', 'general', false),
    ('default_timezone', 'UTC', 'string', 'Default system timezone', 'general', false),
    ('session_timeout_minutes', '1440', 'integer', 'User session timeout in minutes', 'security', false),
    ('max_failed_login_attempts', '5', 'integer', 'Maximum failed login attempts before account lock', 'security', false),
    ('account_lockout_minutes', '30', 'integer', 'Account lockout duration in minutes', 'security', false),
    ('spam_detection_enabled', 'true', 'boolean', 'Enable automatic spam detection', 'routing', false),
    ('spam_short_call_threshold', '10', 'integer', 'Calls shorter than this (seconds) are flagged as potential spam', 'routing', false),
    ('auto_blacklist_enabled', 'true', 'boolean', 'Automatically blacklist numbers flagged as spam', 'routing', false),
    ('low_balance_threshold', '10.00', 'float', 'Send alerts when customer balance falls below this amount', 'billing', false),
    ('currency_symbol', '$', 'string', 'Default currency symbol', 'billing', false)
ON CONFLICT (config_key) DO NOTHING;
