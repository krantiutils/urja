-- Migration 000015: NFC card improvements, activity logs, and NFC devices
-- Part of ur-924o: NFC Cards + Activity Logs API

-- 1. Fix nfc_cards: allow unassigned cards (user_id nullable)
--    and add a sequential card_number for display (#CG0001 format).
ALTER TABLE nfc_cards ALTER COLUMN user_id DROP NOT NULL;
ALTER TABLE nfc_cards ADD COLUMN card_number SERIAL;
ALTER TABLE nfc_cards ADD COLUMN updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW();

-- Unique card_number per org for display purposes
CREATE UNIQUE INDEX idx_nfc_cards_org_card_number ON nfc_cards(organization_id, card_number);

-- 2. Activity logs: append-only event store for org timeline
CREATE TABLE activity_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    event_type VARCHAR(50) NOT NULL
        CHECK (event_type IN (
            'check_in',
            'subscription_approved',
            'membership_request',
            'membership_approved',
            'package_changed',
            'nfc_card_assigned',
            'nfc_card_unassigned',
            'nfc_card_registered',
            'member_joined',
            'member_suspended'
        )),
    description TEXT NOT NULL,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_activity_logs_org_created ON activity_logs(organization_id, created_at DESC);
CREATE INDEX idx_activity_logs_org_type ON activity_logs(organization_id, event_type);
CREATE INDEX idx_activity_logs_user ON activity_logs(user_id) WHERE user_id IS NOT NULL;

-- 3. NFC devices: physical reader devices attached to a gym
CREATE TABLE nfc_devices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    device_identifier VARCHAR(255) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'offline'
        CHECK (status IN ('online', 'offline')),
    door_state VARCHAR(20) NOT NULL DEFAULT 'closed'
        CHECK (door_state IN ('open', 'closed')),
    last_sync_at TIMESTAMPTZ,
    uptime_seconds BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(organization_id, device_identifier)
);

CREATE INDEX idx_nfc_devices_org ON nfc_devices(organization_id);
