-- Feedbacks: member feedback submitted from mobile app
CREATE TABLE IF NOT EXISTS feedbacks (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    user_id         UUID NOT NULL REFERENCES users(id),
    message         TEXT NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_feedbacks_org_id ON feedbacks(organization_id);
CREATE INDEX idx_feedbacks_created_at ON feedbacks(created_at DESC);
