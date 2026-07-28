BEGIN;

DROP INDEX IF EXISTS idx_bout_records_org;
DROP INDEX IF EXISTS idx_bout_records_user;
DROP INDEX IF EXISTS idx_boxing_profiles_org;

DROP TABLE IF EXISTS bout_records;
DROP TABLE IF EXISTS member_boxing_profiles;

COMMIT;
