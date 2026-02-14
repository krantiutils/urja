DROP INDEX IF EXISTS idx_org_members_staff_role;
ALTER TABLE organization_members DROP COLUMN IF EXISTS staff_role;
