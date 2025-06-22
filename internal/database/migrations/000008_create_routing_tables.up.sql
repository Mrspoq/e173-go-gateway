-- Call routing and filtering tables
CREATE TABLE IF NOT EXISTS routing_rules (
    id BIGSERIAL PRIMARY KEY,
    rule_name VARCHAR(100) NOT NULL,
    rule_order INTEGER NOT NULL DEFAULT 1000, -- Lower numbers = higher priority
    prefix_pattern VARCHAR(50) NOT NULL, -- e.g., "+1", "001", "44"
    destination_pattern VARCHAR(50), -- Optional destination constraint
    caller_id_pattern VARCHAR(50), -- Optional caller ID constraint
    route_to_modem_id BIGINT REFERENCES modems(id) ON DELETE SET NULL,
    route_to_pool VARCHAR(50), -- Alternative to specific modem (e.g., "pool_usa", "pool_uk")
    max_channels INTEGER DEFAULT 1, -- Max simultaneous calls for this route
    time_restrictions JSONB, -- JSON: {"weekdays": [1,2,3,4,5], "hours": {"start": "09:00", "end": "17:00"}}
    customer_restrictions BIGINT[] DEFAULT '{}', -- Array of customer IDs allowed to use this route
    cost_markup_percent DECIMAL(5, 2) DEFAULT 0.0, -- Additional markup on base rate
    is_active BOOLEAN DEFAULT TRUE,
    notes TEXT,
    created_by BIGINT REFERENCES users(id),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Call blacklist for spam prevention
CREATE TABLE IF NOT EXISTS blacklist (
    id BIGSERIAL PRIMARY KEY,
    number_pattern VARCHAR(50) NOT NULL, -- e.g., "+1234567890", "+1234*", "*spam*"
    blacklist_type VARCHAR(20) DEFAULT 'number' CHECK (blacklist_type IN ('number', 'pattern', 'prefix')),
    reason VARCHAR(255),
    auto_added BOOLEAN DEFAULT FALSE, -- True if added by spam detection system
    detection_method VARCHAR(50), -- e.g., "short_call", "high_frequency", "manual"
    block_inbound BOOLEAN DEFAULT TRUE,
    block_outbound BOOLEAN DEFAULT FALSE,
    temporary_until TIMESTAMPTZ, -- NULL for permanent blocks
    violation_count INTEGER DEFAULT 1,
    last_violation_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    created_by BIGINT REFERENCES users(id),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- SIM card pool assignments for load balancing
CREATE TABLE IF NOT EXISTS sim_pools (
    id BIGSERIAL PRIMARY KEY,
    pool_name VARCHAR(50) UNIQUE NOT NULL,
    description TEXT,
    load_balance_method VARCHAR(20) DEFAULT 'round_robin' CHECK (load_balance_method IN ('round_robin', 'least_used', 'random', 'failover')),
    max_channels_per_sim INTEGER DEFAULT 1,
    is_active BOOLEAN DEFAULT TRUE,
    created_by BIGINT REFERENCES users(id),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- SIM cards to pool assignments
CREATE TABLE IF NOT EXISTS sim_pool_assignments (
    id BIGSERIAL PRIMARY KEY,
    sim_pool_id BIGINT NOT NULL REFERENCES sim_pools(id) ON DELETE CASCADE,
    sim_card_id BIGINT NOT NULL REFERENCES sim_cards(id) ON DELETE CASCADE,
    priority INTEGER DEFAULT 100, -- Lower = higher priority in pool
    is_active BOOLEAN DEFAULT TRUE,
    assigned_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    assigned_by BIGINT REFERENCES users(id),
    
    UNIQUE(sim_pool_id, sim_card_id)
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_routing_rules_rule_order ON routing_rules(rule_order);
CREATE INDEX IF NOT EXISTS idx_routing_rules_prefix_pattern ON routing_rules(prefix_pattern);
CREATE INDEX IF NOT EXISTS idx_routing_rules_is_active ON routing_rules(is_active);
CREATE INDEX IF NOT EXISTS idx_blacklist_number_pattern ON blacklist(number_pattern);
CREATE INDEX IF NOT EXISTS idx_blacklist_blacklist_type ON blacklist(blacklist_type);
CREATE INDEX IF NOT EXISTS idx_blacklist_auto_added ON blacklist(auto_added);
CREATE INDEX IF NOT EXISTS idx_sim_pools_pool_name ON sim_pools(pool_name);
CREATE INDEX IF NOT EXISTS idx_sim_pool_assignments_sim_pool_id ON sim_pool_assignments(sim_pool_id);
CREATE INDEX IF NOT EXISTS idx_sim_pool_assignments_sim_card_id ON sim_pool_assignments(sim_card_id);

-- Triggers
CREATE TRIGGER set_routing_rules_updated_at
BEFORE UPDATE ON routing_rules
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();

CREATE TRIGGER set_blacklist_updated_at
BEFORE UPDATE ON blacklist
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();

CREATE TRIGGER set_sim_pools_updated_at
BEFORE UPDATE ON sim_pools
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();
