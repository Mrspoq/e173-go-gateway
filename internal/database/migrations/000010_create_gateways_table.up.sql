-- Create gateways table for managing remote E173 gateway instances
CREATE TABLE IF NOT EXISTS gateways (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    location VARCHAR(255),
    ami_host VARCHAR(255) NOT NULL,
    ami_port VARCHAR(10) NOT NULL DEFAULT '5038',
    ami_user VARCHAR(255) NOT NULL,
    ami_pass VARCHAR(255) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'offline',
    enabled BOOLEAN NOT NULL DEFAULT true,
    last_seen TIMESTAMP WITH TIME ZONE,
    last_error TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create index for faster lookups
CREATE INDEX idx_gateways_status ON gateways(status);
CREATE INDEX idx_gateways_enabled ON gateways(enabled);
CREATE INDEX idx_gateways_last_seen ON gateways(last_seen);

-- Add unique constraint on name
ALTER TABLE gateways ADD CONSTRAINT unique_gateway_name UNIQUE (name);

-- Add trigger to update updated_at timestamp
CREATE TRIGGER update_gateways_updated_at BEFORE UPDATE ON gateways
    FOR EACH ROW EXECUTE FUNCTION trigger_set_timestamp();
