-- Restore NOT NULL on assigned_at (only safe if all rows have a value).
UPDATE nfc_cards SET assigned_at = NOW() WHERE assigned_at IS NULL;
ALTER TABLE nfc_cards ALTER COLUMN assigned_at SET NOT NULL;
