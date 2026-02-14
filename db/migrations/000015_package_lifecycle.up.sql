-- Package lifecycle: registration toggle + payments table

-- Add registration toggle to packages
ALTER TABLE packages ADD COLUMN registration_enabled BOOLEAN NOT NULL DEFAULT true;

-- Payments table for tracking payment history per member subscription.
-- Supports discount, due amounts, and multiple payment methods.
CREATE TABLE payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    member_package_id UUID NOT NULL REFERENCES member_packages(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    particular VARCHAR(255) NOT NULL,
    total_amount DECIMAL(10,2) NOT NULL CHECK (total_amount >= 0),
    discount DECIMAL(10,2) NOT NULL DEFAULT 0 CHECK (discount >= 0),
    paid_amount DECIMAL(10,2) NOT NULL CHECK (paid_amount >= 0),
    due DECIMAL(10,2) NOT NULL DEFAULT 0 CHECK (due >= 0),
    payment_method VARCHAR(50)
        CHECK (payment_method IN ('khalti', 'cash', 'manual', 'bank_transfer')),
    payment_reference VARCHAR(255),
    payment_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_payments_user_id ON payments(user_id);
CREATE INDEX idx_payments_org_id ON payments(organization_id);
CREATE INDEX idx_payments_member_package_id ON payments(member_package_id);
