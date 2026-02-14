-- Notices: gym announcements to members
CREATE TABLE IF NOT EXISTS notices (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    title           TEXT NOT NULL,
    title_ne        TEXT,
    content         TEXT NOT NULL,
    content_ne      TEXT,
    created_by      UUID NOT NULL REFERENCES users(id),
    is_active       BOOLEAN NOT NULL DEFAULT true,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_notices_org_id ON notices(organization_id);
CREATE INDEX idx_notices_org_active ON notices(organization_id, is_active) WHERE is_active = true;
CREATE INDEX idx_notices_created_at ON notices(created_at DESC);

CREATE TRIGGER set_notices_updated_at
    BEFORE UPDATE ON notices
    FOR EACH ROW
    EXECUTE FUNCTION trigger_set_updated_at();
