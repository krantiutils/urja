-- Urja Gym Platform: Initial Schema
-- This migration creates the core tables per the PRD specification.
-- Run with: psql -U urja -d urja -f migrations/001_initial_schema.sql

BEGIN;

-- Enable UUID generation
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Organizations (gyms)
CREATE TABLE IF NOT EXISTS organizations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) UNIQUE NOT NULL,
    address TEXT,
    phone VARCHAR(20),
    logo_url TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Users
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    phone VARCHAR(20) UNIQUE NOT NULL,
    name VARCHAR(255),
    email VARCHAR(255),
    avatar_url TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- M:N user <-> organization membership
CREATE TABLE IF NOT EXISTS organization_members (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    organization_id UUID REFERENCES organizations(id) ON DELETE CASCADE,
    role VARCHAR(50) NOT NULL DEFAULT 'member',
    joined_at TIMESTAMP DEFAULT NOW(),
    status VARCHAR(50) DEFAULT 'active',
    UNIQUE(user_id, organization_id)
);

-- Attendance
CREATE TABLE IF NOT EXISTS attendance (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    organization_id UUID REFERENCES organizations(id) ON DELETE CASCADE,
    method VARCHAR(20) NOT NULL DEFAULT 'manual',
    check_in_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- NFC cards
CREATE TABLE IF NOT EXISTS nfc_cards (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    organization_id UUID REFERENCES organizations(id) ON DELETE CASCADE,
    hex_code VARCHAR(255) NOT NULL,
    status VARCHAR(50) DEFAULT 'active',
    assigned_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(organization_id, hex_code)
);

-- Packages (membership plans offered by gyms)
CREATE TABLE IF NOT EXISTS packages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID REFERENCES organizations(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    duration_days INT NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    description TEXT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Member package subscriptions
CREATE TABLE IF NOT EXISTS member_packages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    package_id UUID REFERENCES packages(id) ON DELETE CASCADE,
    organization_id UUID REFERENCES organizations(id) ON DELETE CASCADE,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    status VARCHAR(50) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT NOW()
);

-- Streaks
CREATE TABLE IF NOT EXISTS streaks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    member_id UUID REFERENCES users(id) ON DELETE CASCADE,
    organization_id UUID REFERENCES organizations(id),
    current_streak INT DEFAULT 0,
    longest_streak INT DEFAULT 0,
    last_check_in DATE,
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Achievements
CREATE TABLE IF NOT EXISTS achievements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    name_ne VARCHAR(255),
    description TEXT,
    description_ne TEXT,
    icon_url TEXT,
    criteria JSONB NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS member_achievements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    member_id UUID REFERENCES users(id) ON DELETE CASCADE,
    achievement_id UUID REFERENCES achievements(id),
    earned_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(member_id, achievement_id)
);

-- Health metrics
CREATE TABLE IF NOT EXISTS health_metrics (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    member_id UUID REFERENCES users(id) ON DELETE CASCADE,
    metric_type VARCHAR(50) NOT NULL,
    value JSONB NOT NULL,
    recorded_at TIMESTAMP DEFAULT NOW()
);

-- Progress photos
CREATE TABLE IF NOT EXISTS progress_photos (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    member_id UUID REFERENCES users(id) ON DELETE CASCADE,
    photo_url TEXT NOT NULL,
    category VARCHAR(50),
    notes TEXT,
    taken_at TIMESTAMP DEFAULT NOW()
);

-- Workout templates
CREATE TABLE IF NOT EXISTS workout_templates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID REFERENCES organizations(id),
    name VARCHAR(255) NOT NULL,
    name_ne VARCHAR(255),
    description TEXT,
    category VARCHAR(100),
    difficulty VARCHAR(50),
    duration_minutes INT,
    exercises JSONB NOT NULL,
    is_preset BOOLEAN DEFAULT false,
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT NOW()
);

-- Training guides
CREATE TABLE IF NOT EXISTS training_guides (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    title_ne VARCHAR(255),
    content TEXT NOT NULL,
    content_ne TEXT,
    category VARCHAR(100),
    cover_image_url TEXT,
    author_id UUID REFERENCES users(id),
    is_published BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT NOW()
);

-- SaaS billing
CREATE TABLE IF NOT EXISTS saas_plans (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    price_monthly DECIMAL(10,2),
    price_yearly DECIMAL(10,2),
    features JSONB DEFAULT '[]',
    max_members INT,
    is_active BOOLEAN DEFAULT true
);

CREATE TABLE IF NOT EXISTS saas_subscriptions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID REFERENCES organizations(id),
    plan_id UUID REFERENCES saas_plans(id),
    status VARCHAR(50) DEFAULT 'active',
    current_period_start TIMESTAMP,
    current_period_end TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Indexes for common queries
CREATE INDEX IF NOT EXISTS idx_attendance_user ON attendance(user_id, check_in_at DESC);
CREATE INDEX IF NOT EXISTS idx_attendance_org ON attendance(organization_id, check_in_at DESC);
CREATE INDEX IF NOT EXISTS idx_org_members_user ON organization_members(user_id);
CREATE INDEX IF NOT EXISTS idx_org_members_org ON organization_members(organization_id);
CREATE INDEX IF NOT EXISTS idx_member_packages_user ON member_packages(user_id);
CREATE INDEX IF NOT EXISTS idx_health_metrics_member ON health_metrics(member_id, recorded_at DESC);
CREATE INDEX IF NOT EXISTS idx_streaks_member ON streaks(member_id, organization_id);
CREATE INDEX IF NOT EXISTS idx_nfc_cards_org ON nfc_cards(organization_id);

COMMIT;
