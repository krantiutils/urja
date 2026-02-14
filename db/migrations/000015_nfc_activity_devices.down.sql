DROP TABLE IF EXISTS nfc_devices;
DROP TABLE IF EXISTS activity_logs;

DROP INDEX IF EXISTS idx_nfc_cards_org_card_number;
ALTER TABLE nfc_cards DROP COLUMN IF EXISTS updated_at;
ALTER TABLE nfc_cards DROP COLUMN IF EXISTS card_number;
ALTER TABLE nfc_cards ALTER COLUMN user_id SET NOT NULL;
