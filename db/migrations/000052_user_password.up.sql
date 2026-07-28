-- Optional password login alongside phone + OTP.
--
-- OTP is the primary route and stays that way: it needs no memorised secret,
-- which suits members. But staff and gym admins sign in several times a day,
-- often from a desk, and waiting on an SMS each time is friction they feel.
--
-- Nullable on purpose. A NULL password_hash means the account simply has no
-- password and can only sign in by OTP, which is the state every existing
-- account starts in and most will stay in.
ALTER TABLE users
  ADD COLUMN IF NOT EXISTS password_hash TEXT,
  ADD COLUMN IF NOT EXISTS password_set_at TIMESTAMPTZ;

-- Password login looks the account up by phone, the same identifier OTP uses,
-- so there is one identity per person rather than a separate username space.
CREATE INDEX IF NOT EXISTS idx_users_password_hash_present
  ON users (phone) WHERE password_hash IS NOT NULL;
