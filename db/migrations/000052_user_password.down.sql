DROP INDEX IF EXISTS idx_users_password_hash_present;

ALTER TABLE users
  DROP COLUMN IF EXISTS password_hash,
  DROP COLUMN IF EXISTS password_set_at;
