-- Migration: Add gateway_id to modems table for multi-gateway support
-- This allows each modem to be associated with a specific gateway (box1, box2, etc)

-- Add gateway_id column to modems table
ALTER TABLE modems 
ADD COLUMN IF NOT EXISTS gateway_id UUID REFERENCES gateways(id);

-- Add index for faster gateway-based queries
CREATE INDEX IF NOT EXISTS idx_modems_gateway_id ON modems(gateway_id);

-- For existing modems, associate them with the first gateway (if any exists)
-- In production, you'd want to manually assign these
UPDATE modems 
SET gateway_id = (SELECT id FROM gateways ORDER BY created_at LIMIT 1)
WHERE gateway_id IS NULL;