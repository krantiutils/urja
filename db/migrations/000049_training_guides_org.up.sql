-- Scope training guides to an owning organization.
--
-- training_guides was created with no organization_id, yet guide management is
-- mounted under /orgs/{orgId}/training-guides. Every mutating repository method
-- keyed on the guide id alone, so any gym's staff could edit, unpublish or
-- delete the entire platform's guide library.
--
-- Existing rows keep NULL, which correctly marks them as the platform-wide
-- presets they already are. NULL-org guides are readable by everyone but
-- mutable only outside the org routes.

BEGIN;

ALTER TABLE training_guides
    ADD COLUMN IF NOT EXISTS organization_id UUID REFERENCES organizations(id) ON DELETE CASCADE;

COMMENT ON COLUMN training_guides.organization_id IS
    'Owning organization. NULL means a platform-wide preset, not editable through org routes.';

CREATE INDEX IF NOT EXISTS idx_training_guides_org ON training_guides(organization_id);

COMMIT;
