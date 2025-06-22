-- Function to update updated_at column
CREATE OR REPLACE FUNCTION trigger_set_timestamp()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Modems table
CREATE TABLE IF NOT EXISTS modems (
    id SERIAL PRIMARY KEY,
    device_path VARCHAR(255) UNIQUE NOT NULL,
    imei VARCHAR(20) UNIQUE,
    imsi VARCHAR(20) UNIQUE,
    model VARCHAR(100),
    manufacturer VARCHAR(100),
    firmware_version VARCHAR(100),
    signal_strength_dbm INTEGER,
    network_operator_name VARCHAR(100),
    network_registration_status VARCHAR(50),
    status VARCHAR(50) NOT NULL DEFAULT 'unknown',
    last_seen_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Index for faster lookups
CREATE INDEX IF NOT EXISTS idx_modems_status ON modems(status);
CREATE INDEX IF NOT EXISTS idx_modems_imei ON modems(imei);
CREATE INDEX IF NOT EXISTS idx_modems_imsi ON modems(imsi);

-- Trigger for updated_at
CREATE TRIGGER set_modems_updated_at
BEFORE UPDATE ON modems
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();
