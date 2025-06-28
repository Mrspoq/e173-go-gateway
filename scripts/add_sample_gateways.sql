-- Add sample gateways for testing
-- These represent different E173 gateway boxes with their own Asterisk instances

-- Insert sample gateways
INSERT INTO gateways (id, name, description, location, ip_address, port, api_key, ami_host, ami_port, ami_username, ami_password, status, last_heartbeat, created_at, updated_at) VALUES
('550e8400-e29b-41d4-a716-446655440001'::uuid, 'box1', 'Primary Gateway - Local', 'Office', '127.0.0.1', 8081, 'box1-api-key-123', 'localhost', 5038, 'admin', '3omartel580', 'active', NOW(), NOW(), NOW()),
('550e8400-e29b-41d4-a716-446655440002'::uuid, 'box2', 'Secondary Gateway - Remote', 'Data Center 1', '10.0.1.100', 8082, 'box2-api-key-456', '10.0.1.100', 5038, 'admin', '3omartel580', 'offline', NOW() - INTERVAL '1 hour', NOW(), NOW()),
('550e8400-e29b-41d4-a716-446655440003'::uuid, 'box3', 'Tertiary Gateway - Backup', 'Data Center 2', '10.0.2.100', 8083, 'box3-api-key-789', '10.0.2.100', 5038, 'admin', '3omartel580', 'active', NOW() - INTERVAL '5 minutes', NOW(), NOW())
ON CONFLICT (id) DO UPDATE SET
    name = EXCLUDED.name,
    description = EXCLUDED.description,
    status = EXCLUDED.status,
    updated_at = NOW();

-- Add sample modems for box1 (local gateway)
INSERT INTO modems (gateway_id, imei, imsi, phone_number, operator, signal_strength, status, last_seen) VALUES
('550e8400-e29b-41d4-a716-446655440001'::uuid, '869385031234561', '234150000000001', '+1234567001', 'Operator1', 4, 'online', NOW()),
('550e8400-e29b-41d4-a716-446655440001'::uuid, '869385031234562', '234150000000002', '+1234567002', 'Operator1', 3, 'online', NOW()),
('550e8400-e29b-41d4-a716-446655440001'::uuid, '869385031234563', '234150000000003', '+1234567003', 'Operator2', 5, 'online', NOW()),
('550e8400-e29b-41d4-a716-446655440001'::uuid, '869385031234564', '234150000000004', '+1234567004', 'Operator2', 2, 'offline', NOW() - INTERVAL '10 minutes'),
('550e8400-e29b-41d4-a716-446655440001'::uuid, '869385031234565', '234150000000005', '+1234567005', 'Operator3', 4, 'online', NOW())
ON CONFLICT (imei) DO UPDATE SET
    status = EXCLUDED.status,
    last_seen = EXCLUDED.last_seen;

-- Add sample modems for box3 (remote gateway)
INSERT INTO modems (gateway_id, imei, imsi, phone_number, operator, signal_strength, status, last_seen) VALUES
('550e8400-e29b-41d4-a716-446655440003'::uuid, '869385031234571', '234150000000011', '+1234567011', 'Operator1', 5, 'online', NOW()),
('550e8400-e29b-41d4-a716-446655440003'::uuid, '869385031234572', '234150000000012', '+1234567012', 'Operator2', 4, 'online', NOW()),
('550e8400-e29b-41d4-a716-446655440003'::uuid, '869385031234573', '234150000000013', '+1234567013', 'Operator3', 3, 'online', NOW())
ON CONFLICT (imei) DO UPDATE SET
    status = EXCLUDED.status,
    last_seen = EXCLUDED.last_seen;