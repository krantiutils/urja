CREATE TABLE attendance (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    check_in_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    check_out_at TIMESTAMPTZ,
    check_in_method VARCHAR(20) NOT NULL
        CHECK (check_in_method IN ('nfc', 'qr', 'manual')),
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_attendance_user_id ON attendance(user_id);
CREATE INDEX idx_attendance_org_id ON attendance(organization_id);
CREATE INDEX idx_attendance_org_checkin ON attendance(organization_id, check_in_at DESC);
CREATE INDEX idx_attendance_user_org_date ON attendance(user_id, organization_id, check_in_at DESC);

CREATE TABLE nfc_cards (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    card_hex VARCHAR(50) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    assigned_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deactivated_at TIMESTAMPTZ,
    UNIQUE(organization_id, card_hex)
);

CREATE INDEX idx_nfc_cards_user_id ON nfc_cards(user_id);
CREATE INDEX idx_nfc_cards_active_hex ON nfc_cards(card_hex) WHERE is_active = true;
CREATE INDEX idx_nfc_cards_org_active ON nfc_cards(organization_id) WHERE is_active = true;
