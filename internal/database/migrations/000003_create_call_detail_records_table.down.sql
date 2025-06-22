DROP TRIGGER IF EXISTS set_cdr_updated_at ON call_detail_records;
-- Note: We do not drop trigger_set_timestamp() here as it is shared by other tables.
DROP TABLE IF EXISTS call_detail_records;