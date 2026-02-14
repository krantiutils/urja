-- Dues: tracks outstanding payments owed by members
CREATE TABLE dues (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    amount DECIMAL(10,2) NOT NULL CHECK (amount > 0),
    due_date DATE NOT NULL,
    description TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'pending'
        CHECK (status IN ('pending', 'paid', 'overdue', 'waived')),
    paid_at TIMESTAMPTZ,
    paid_amount DECIMAL(10,2) CHECK (paid_amount IS NULL OR paid_amount >= 0),
    payment_method VARCHAR(50)
        CHECK (payment_method IS NULL OR payment_method IN ('cash', 'bank_transfer', 'khalti')),
    payment_reference VARCHAR(255),
    member_package_id UUID REFERENCES member_packages(id) ON DELETE SET NULL,
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_dues_org_id ON dues(organization_id);
CREATE INDEX idx_dues_user_id ON dues(user_id);
CREATE INDEX idx_dues_org_status ON dues(organization_id, status) WHERE status IN ('pending', 'overdue');
CREATE INDEX idx_dues_org_user ON dues(organization_id, user_id);
CREATE INDEX idx_dues_due_date ON dues(due_date) WHERE status IN ('pending', 'overdue');

CREATE TRIGGER set_dues_updated_at
    BEFORE UPDATE ON dues
    FOR EACH ROW
    EXECUTE FUNCTION trigger_set_updated_at();

-- Transactions: financial ledger for income and expenses
CREATE TABLE transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    category VARCHAR(100) NOT NULL,
    description TEXT,
    transaction_date DATE NOT NULL,
    transaction_type VARCHAR(20) NOT NULL
        CHECK (transaction_type IN ('income', 'expense')),
    amount DECIMAL(10,2) NOT NULL CHECK (amount > 0),
    payment_type VARCHAR(50) NOT NULL
        CHECK (payment_type IN ('cash', 'bank_transfer', 'khalti')),
    reference VARCHAR(255),
    due_id UUID REFERENCES dues(id) ON DELETE SET NULL,
    entry_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_transactions_org_id ON transactions(organization_id);
CREATE INDEX idx_transactions_org_date ON transactions(organization_id, transaction_date DESC);
CREATE INDEX idx_transactions_org_type ON transactions(organization_id, transaction_type);
CREATE INDEX idx_transactions_org_category ON transactions(organization_id, category);
CREATE INDEX idx_transactions_due_id ON transactions(due_id) WHERE due_id IS NOT NULL;

CREATE TRIGGER set_transactions_updated_at
    BEFORE UPDATE ON transactions
    FOR EACH ROW
    EXECUTE FUNCTION trigger_set_updated_at();
