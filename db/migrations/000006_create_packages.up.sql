CREATE TABLE packages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    name_ne VARCHAR(255),
    description TEXT,
    description_ne TEXT,
    duration_days INT NOT NULL CHECK (duration_days > 0),
    price DECIMAL(10,2) NOT NULL CHECK (price >= 0),
    currency VARCHAR(10) NOT NULL DEFAULT 'NPR',
    max_members INT,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_packages_org_id ON packages(organization_id);
CREATE INDEX idx_packages_org_active ON packages(organization_id, is_active) WHERE is_active = true;

CREATE TRIGGER set_packages_updated_at
    BEFORE UPDATE ON packages
    FOR EACH ROW
    EXECUTE FUNCTION trigger_set_updated_at();

CREATE TABLE member_packages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    package_id UUID NOT NULL REFERENCES packages(id) ON DELETE RESTRICT,
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    payment_method VARCHAR(50)
        CHECK (payment_method IN ('khalti', 'cash', 'manual', 'bank_transfer')),
    payment_reference VARCHAR(255),
    amount_paid DECIMAL(10,2) NOT NULL CHECK (amount_paid >= 0),
    status VARCHAR(50) NOT NULL DEFAULT 'active'
        CHECK (status IN ('active', 'expired', 'cancelled', 'pending')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CHECK (end_date > start_date)
);

CREATE INDEX idx_member_packages_user_id ON member_packages(user_id);
CREATE INDEX idx_member_packages_org_id ON member_packages(organization_id);
CREATE INDEX idx_member_packages_active ON member_packages(status) WHERE status = 'active';
CREATE INDEX idx_member_packages_end_date ON member_packages(end_date) WHERE status = 'active';
