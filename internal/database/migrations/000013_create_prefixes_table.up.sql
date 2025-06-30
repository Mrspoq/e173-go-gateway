-- Create prefixes table for routing
CREATE TABLE IF NOT EXISTS prefixes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    prefix VARCHAR(20) NOT NULL UNIQUE,
    country VARCHAR(100) NOT NULL,
    operator VARCHAR(100),
    gateway_id UUID REFERENCES gateways(id) ON DELETE SET NULL,
    rate_per_minute DECIMAL(10,4) DEFAULT 0.0000,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create index for prefix matching
CREATE INDEX idx_prefixes_prefix ON prefixes(prefix);
CREATE INDEX idx_prefixes_active ON prefixes(is_active);
CREATE INDEX idx_prefixes_gateway ON prefixes(gateway_id);

-- Create prefix routes table for advanced routing
CREATE TABLE IF NOT EXISTS prefix_routes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    prefix_id UUID NOT NULL REFERENCES prefixes(id) ON DELETE CASCADE,
    gateway_id UUID NOT NULL REFERENCES gateways(id) ON DELETE CASCADE,
    priority INT DEFAULT 1,
    weight INT DEFAULT 100,
    max_concurrent INT DEFAULT 0, -- 0 means unlimited
    current_active INT DEFAULT 0,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(prefix_id, gateway_id)
);

-- Create indexes
CREATE INDEX idx_prefix_routes_prefix ON prefix_routes(prefix_id);
CREATE INDEX idx_prefix_routes_gateway ON prefix_routes(gateway_id);
CREATE INDEX idx_prefix_routes_active ON prefix_routes(is_active);

-- Insert some sample prefixes
INSERT INTO prefixes (prefix, country, operator, is_active) VALUES
('1', 'United States', 'Various', true),
('44', 'United Kingdom', 'Various', true),
('33', 'France', 'Various', true),
('49', 'Germany', 'Various', true),
('34', 'Spain', 'Various', true),
('39', 'Italy', 'Various', true),
('86', 'China', 'Various', true),
('91', 'India', 'Various', true),
('234', 'Nigeria', 'Various', true),
('254', 'Kenya', 'Various', true),
('27', 'South Africa', 'Various', true),
('55', 'Brazil', 'Various', true),
('52', 'Mexico', 'Various', true),
('966', 'Saudi Arabia', 'Various', true),
('971', 'UAE', 'Various', true),
('20', 'Egypt', 'Various', true),
('212', 'Morocco', 'Various', true),
('90', 'Turkey', 'Various', true),
('7', 'Russia', 'Various', true),
('380', 'Ukraine', 'Various', true);