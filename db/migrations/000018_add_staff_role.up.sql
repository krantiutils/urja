ALTER TABLE organization_members
    ADD COLUMN staff_role VARCHAR(50)
        CHECK (staff_role IS NULL OR staff_role IN ('owner', 'manager', 'trainer', 'receptionist'));

CREATE INDEX idx_org_members_staff_role ON organization_members(organization_id, staff_role)
    WHERE staff_role IS NOT NULL;
