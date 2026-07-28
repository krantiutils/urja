-- Combat-sports member profile: stance, weight class, sparring clearance and
-- bout history. Extends the generic member/health modules for boxing gyms.

BEGIN;

CREATE TABLE IF NOT EXISTS member_boxing_profiles (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id             UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    organization_id     UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    stance              VARCHAR(20),
    weight_class        VARCHAR(40),
    skill_level         VARCHAR(30),
    -- Sparring clearance is a safety gate: only staff may grant it, and we
    -- record who did and when so it is auditable.
    sparring_cleared    BOOLEAN NOT NULL DEFAULT false,
    sparring_cleared_at TIMESTAMPTZ,
    sparring_cleared_by UUID REFERENCES users(id) ON DELETE SET NULL,
    reach_cm            NUMERIC(5,1),
    notes               TEXT,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (user_id, organization_id)
);

CREATE TABLE IF NOT EXISTS bout_records (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    bout_date       DATE NOT NULL,
    opponent        VARCHAR(200),
    event_name      VARCHAR(200),
    result          VARCHAR(20) NOT NULL,
    method          VARCHAR(30),
    rounds          INT,
    weight_class    VARCHAR(40),
    notes           TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_boxing_profiles_org ON member_boxing_profiles(organization_id);
CREATE INDEX IF NOT EXISTS idx_bout_records_user ON bout_records(user_id, bout_date DESC);
CREATE INDEX IF NOT EXISTS idx_bout_records_org ON bout_records(organization_id, bout_date DESC);

COMMIT;
