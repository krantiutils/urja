-- Allow self-signup users to be created without a name.
-- Name is provided during onboarding, not at phone login.
ALTER TABLE users ALTER COLUMN name SET DEFAULT '';
