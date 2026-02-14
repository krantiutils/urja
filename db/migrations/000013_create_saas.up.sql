CREATE TABLE saas_plans (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    name_ne VARCHAR(255),
    description TEXT,
    description_ne TEXT,
    price_monthly DECIMAL(10,2) NOT NULL CHECK (price_monthly >= 0),
    price_yearly DECIMAL(10,2) NOT NULL CHECK (price_yearly >= 0),
    features JSONB NOT NULL DEFAULT '[]',
    max_members INT,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TRIGGER set_saas_plans_updated_at
    BEFORE UPDATE ON saas_plans
    FOR EACH ROW
    EXECUTE FUNCTION trigger_set_updated_at();

CREATE TABLE saas_subscriptions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE RESTRICT,
    plan_id UUID NOT NULL REFERENCES saas_plans(id) ON DELETE RESTRICT,
    status VARCHAR(50) NOT NULL DEFAULT 'active'
        CHECK (status IN ('active', 'cancelled', 'past_due', 'trialing')),
    current_period_start TIMESTAMPTZ NOT NULL,
    current_period_end TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CHECK (current_period_end > current_period_start)
);

CREATE INDEX idx_saas_subscriptions_org_id ON saas_subscriptions(organization_id);
CREATE INDEX idx_saas_subscriptions_active ON saas_subscriptions(status) WHERE status = 'active';

CREATE TRIGGER set_saas_subscriptions_updated_at
    BEFORE UPDATE ON saas_subscriptions
    FOR EACH ROW
    EXECUTE FUNCTION trigger_set_updated_at();
