-- Tenant public websites: one site per organization, served at <slug>.nepalgym.xyz.
-- Page content is an ordered array of section objects validated by internal/site
-- before it is ever written, so the public renderer never sees malformed data.

BEGIN;

CREATE TABLE IF NOT EXISTS site_pages (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    slug            VARCHAR(120) NOT NULL,
    title           VARCHAR(200) NOT NULL,
    title_ne        VARCHAR(200),
    sections        JSONB NOT NULL DEFAULT '[]',
    seo_description TEXT,
    is_published    BOOLEAN NOT NULL DEFAULT false,
    show_in_nav     BOOLEAN NOT NULL DEFAULT true,
    sort_order      INT NOT NULL DEFAULT 0,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (organization_id, slug)
);

CREATE TABLE IF NOT EXISTS site_settings (
    organization_id UUID PRIMARY KEY REFERENCES organizations(id) ON DELETE CASCADE,
    template        VARCHAR(50) NOT NULL DEFAULT 'fight_club',
    theme           JSONB NOT NULL DEFAULT '{}',
    nav             JSONB NOT NULL DEFAULT '{}',
    footer          JSONB NOT NULL DEFAULT '{}',
    socials         JSONB NOT NULL DEFAULT '{}',
    is_live         BOOLEAN NOT NULL DEFAULT false,
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMENT ON COLUMN site_settings.is_live IS
    'Master switch. While false the subdomain serves a coming-soon page instead of a half-built site.';

CREATE TABLE IF NOT EXISTS site_leads (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    name            VARCHAR(200) NOT NULL,
    phone           VARCHAR(30) NOT NULL,
    email           VARCHAR(255),
    message         TEXT,
    interest        VARCHAR(100),
    source_page     VARCHAR(120),
    status          VARCHAR(30) NOT NULL DEFAULT 'new',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_site_pages_org ON site_pages(organization_id, sort_order);
CREATE INDEX IF NOT EXISTS idx_site_leads_org ON site_leads(organization_id, created_at DESC);

COMMIT;
