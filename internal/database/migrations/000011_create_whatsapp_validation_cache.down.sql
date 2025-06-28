-- Drop WhatsApp validation cache table
DROP TRIGGER IF EXISTS update_whatsapp_validation_updated_at ON whatsapp_validation_cache;
DROP FUNCTION IF EXISTS update_updated_at_column();
DROP TABLE IF EXISTS whatsapp_validation_cache;