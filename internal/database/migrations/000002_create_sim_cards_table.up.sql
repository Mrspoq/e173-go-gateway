-- 000002_create_sim_cards_table.up.sql

-- Function to update updated_at column (if not already created by modem migration, idempotent)
CREATE OR REPLACE FUNCTION trigger_set_timestamp()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TABLE IF NOT EXISTS sim_cards (
    id SERIAL PRIMARY KEY,
    modem_id INTEGER REFERENCES modems(id) ON DELETE SET NULL, -- Link to modems, set NULL if modem is deleted
    iccid VARCHAR(25) UNIQUE NOT NULL,
    imsi VARCHAR(20) UNIQUE,
    msisdn VARCHAR(20) UNIQUE,
    operator_name VARCHAR(100),
    network_country_code VARCHAR(10),
    balance DECIMAL(10, 4) DEFAULT 0.00,
    balance_currency VARCHAR(10),
    balance_last_checked_at TIMESTAMPTZ,
    data_allowance_mb INTEGER,
    data_used_mb INTEGER,
    status VARCHAR(50) NOT NULL DEFAULT 'unknown', -- e.g., 'active', 'inactive', 'blocked', 'low_credit', 'needs_recharge', 'in_use', 'available', 'error'
    pin1 VARCHAR(10),
    puk1 VARCHAR(10),
    pin2 VARCHAR(10),
    puk2 VARCHAR(10),
    activation_date DATE,
    expiry_date DATE,
    recharge_history JSONB,
    notes TEXT,
    cell_id VARCHAR(50),
    lac VARCHAR(50), -- Location Area Code
    psc VARCHAR(50), -- Primary Scrambling Code (3G)
    rscp INTEGER,    -- Received Signal Code Power (3G)
    ecio INTEGER,    -- Ec/Io or Ec/No (3G)
    bts_info_history JSONB, -- Historical BTS data
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Trigger to automatically update updated_at timestamp
CREATE TRIGGER set_sim_cards_updated_at
BEFORE UPDATE ON sim_cards
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();

-- Indexes for faster lookups
CREATE INDEX IF NOT EXISTS idx_sim_cards_modem_id ON sim_cards(modem_id);
CREATE INDEX IF NOT EXISTS idx_sim_cards_iccid ON sim_cards(iccid);
CREATE INDEX IF NOT EXISTS idx_sim_cards_msisdn ON sim_cards(msisdn);
CREATE INDEX IF NOT EXISTS idx_sim_cards_status ON sim_cards(status);