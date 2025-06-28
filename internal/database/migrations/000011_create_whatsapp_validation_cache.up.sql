-- Create WhatsApp validation cache table
CREATE TABLE IF NOT EXISTS whatsapp_validation_cache (
    phone_number VARCHAR(50) PRIMARY KEY,
    has_whatsapp BOOLEAN NOT NULL DEFAULT false,
    profile_name VARCHAR(255),
    is_business BOOLEAN NOT NULL DEFAULT false,
    confidence DECIMAL(3,2) NOT NULL DEFAULT 0.00,
    checked_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    raw_response JSONB,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_whatsapp_checked ON whatsapp_validation_cache (checked_at);
CREATE INDEX IF NOT EXISTS idx_whatsapp_expires ON whatsapp_validation_cache (expires_at);
CREATE INDEX IF NOT EXISTS idx_whatsapp_has_whatsapp ON whatsapp_validation_cache (has_whatsapp);

-- Create a trigger to update the updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_whatsapp_validation_updated_at BEFORE UPDATE
ON whatsapp_validation_cache FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Add comment to table
COMMENT ON TABLE whatsapp_validation_cache IS 'Cache for WhatsApp validation API results to reduce API calls';
COMMENT ON COLUMN whatsapp_validation_cache.phone_number IS 'Phone number in international format';
COMMENT ON COLUMN whatsapp_validation_cache.has_whatsapp IS 'Whether the number has WhatsApp';
COMMENT ON COLUMN whatsapp_validation_cache.profile_name IS 'WhatsApp profile name if available';
COMMENT ON COLUMN whatsapp_validation_cache.is_business IS 'Whether this is a WhatsApp Business account';
COMMENT ON COLUMN whatsapp_validation_cache.confidence IS 'Confidence score of the validation (0.00-1.00)';
COMMENT ON COLUMN whatsapp_validation_cache.checked_at IS 'When the validation was performed';
COMMENT ON COLUMN whatsapp_validation_cache.expires_at IS 'When this cache entry expires';
COMMENT ON COLUMN whatsapp_validation_cache.raw_response IS 'Raw API response in JSON format';