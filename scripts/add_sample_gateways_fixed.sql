-- Add sample gateways for testing
-- These represent different E173 gateway boxes with their own Asterisk instances

-- Insert sample gateways
INSERT INTO gateways (id, name, description, location, ami_host, ami_port, ami_user, ami_pass, status, last_seen, created_at, updated_at) VALUES
('550e8400-e29b-41d4-a716-446655440001'::uuid, 'box1', 'Primary Gateway - Local', 'Office', 'localhost', '5038', 'admin', '3omartel580', 'active', NOW(), NOW(), NOW()),
('550e8400-e29b-41d4-a716-446655440002'::uuid, 'box2', 'Secondary Gateway - Remote', 'Data Center 1', '10.0.1.100', '5038', 'admin', '3omartel580', 'offline', NOW() - INTERVAL '1 hour', NOW(), NOW()),
('550e8400-e29b-41d4-a716-446655440003'::uuid, 'box3', 'Tertiary Gateway - Backup', 'Data Center 2', '10.0.2.100', '5038', 'admin', '3omartel580', 'active', NOW() - INTERVAL '5 minutes', NOW(), NOW())
ON CONFLICT (name) DO UPDATE SET
    description = EXCLUDED.description,
    status = EXCLUDED.status,
    last_seen = EXCLUDED.last_seen,
    updated_at = NOW();

-- Show the created gateways
SELECT id, name, description, status, ami_host FROM gateways ORDER BY name;