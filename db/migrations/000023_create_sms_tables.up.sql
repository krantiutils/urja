-- SMS credits: per-org SMS balance tracking
CREATE TABLE IF NOT EXISTS sms_credits (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL UNIQUE REFERENCES organizations(id) ON DELETE CASCADE,
    balance         INTEGER NOT NULL DEFAULT 0 CHECK (balance >= 0),
    total_purchased INTEGER NOT NULL DEFAULT 0,
    total_used      INTEGER NOT NULL DEFAULT 0,
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TRIGGER set_sms_credits_updated_at
    BEFORE UPDATE ON sms_credits
    FOR EACH ROW
    EXECUTE FUNCTION trigger_set_updated_at();

-- SMS purchases: purchase history for SMS credits
CREATE TABLE IF NOT EXISTS sms_purchases (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    quantity        INTEGER NOT NULL CHECK (quantity > 0),
    rate            NUMERIC(10, 2) NOT NULL CHECK (rate > 0),
    amount          NUMERIC(10, 2) NOT NULL CHECK (amount > 0),
    payment_method  TEXT NOT NULL CHECK (payment_method IN ('khalti', 'cash', 'manual', 'bank_transfer')),
    payment_status  TEXT NOT NULL DEFAULT 'pending' CHECK (payment_status IN ('pending', 'completed', 'failed')),
    purchased_by    UUID NOT NULL REFERENCES users(id),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_sms_purchases_org_id ON sms_purchases(organization_id);
CREATE INDEX idx_sms_purchases_created_at ON sms_purchases(created_at DESC);

-- SMS campaigns: track bulk SMS sends
CREATE TABLE IF NOT EXISTS sms_campaigns (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    message         TEXT NOT NULL,
    recipient_count INTEGER NOT NULL CHECK (recipient_count > 0),
    sent_by         UUID NOT NULL REFERENCES users(id),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_sms_campaigns_org_id ON sms_campaigns(organization_id);
CREATE INDEX idx_sms_campaigns_created_at ON sms_campaigns(created_at DESC);
