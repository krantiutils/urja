BEGIN;

DROP INDEX IF EXISTS idx_training_guides_org;
ALTER TABLE training_guides DROP COLUMN IF EXISTS organization_id;

COMMIT;
