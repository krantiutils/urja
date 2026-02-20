-- E2E Test Seed Data for Urja Gym Management Platform
-- This file is idempotent: all INSERTs use ON CONFLICT DO NOTHING.
-- Test phones (9800*) can log in with dev OTP bypass code "123456".

BEGIN;

-- ============================================================
-- 1. Test Organization: Kathmandu Fitness Hub
-- ============================================================
INSERT INTO organizations (id, name, name_ne, slug, description, phone, email, address)
VALUES (
  'a0000000-0000-0000-0000-000000000001',
  'Kathmandu Fitness Hub',
  'काठमाडौं फिटनेस हब',
  'kathmandu-fitness-hub',
  'Premier fitness center in the heart of Kathmandu',
  '014-123456',
  'info@ktmfitness.com',
  'Thamel, Kathmandu'
) ON CONFLICT DO NOTHING;

-- ============================================================
-- 2. Three Core Test Users
-- ============================================================

-- Admin / Owner: Rajesh Sharma
INSERT INTO users (id, phone, name, name_ne, email, gender, date_of_birth)
VALUES (
  'b0000000-0000-0000-0000-000000000001',
  '9800000001',
  'Rajesh Sharma',
  'राजेश शर्मा',
  'rajesh@ktmfitness.com',
  'male',
  '1985-03-15'
) ON CONFLICT DO NOTHING;

-- Staff / Trainer: Sita Tamang
INSERT INTO users (id, phone, name, name_ne, email, gender, date_of_birth)
VALUES (
  'b0000000-0000-0000-0000-000000000002',
  '9800000002',
  'Sita Tamang',
  'सिता तामाङ',
  'sita@ktmfitness.com',
  'female',
  '1992-07-22'
) ON CONFLICT DO NOTHING;

-- Member: Hari Prasad
INSERT INTO users (id, phone, name, name_ne, email, gender, date_of_birth)
VALUES (
  'b0000000-0000-0000-0000-000000000003',
  '9800000003',
  'Hari Prasad',
  'हरि प्रसाद',
  'hari@example.com',
  'male',
  '1998-11-05'
) ON CONFLICT DO NOTHING;

-- ============================================================
-- 3. Organization Memberships (roles + staff_roles)
-- ============================================================

-- Rajesh = admin + owner
INSERT INTO organization_members (user_id, organization_id, role, staff_role, status)
VALUES (
  'b0000000-0000-0000-0000-000000000001',
  'a0000000-0000-0000-0000-000000000001',
  'admin',
  'owner',
  'active'
) ON CONFLICT DO NOTHING;

-- Sita = staff + trainer
INSERT INTO organization_members (user_id, organization_id, role, staff_role, status)
VALUES (
  'b0000000-0000-0000-0000-000000000002',
  'a0000000-0000-0000-0000-000000000001',
  'staff',
  'trainer',
  'active'
) ON CONFLICT DO NOTHING;

-- Hari = member (no staff_role)
INSERT INTO organization_members (user_id, organization_id, role, staff_role, status)
VALUES (
  'b0000000-0000-0000-0000-000000000003',
  'a0000000-0000-0000-0000-000000000001',
  'member',
  NULL,
  'active'
) ON CONFLICT DO NOTHING;

-- ============================================================
-- 4. Three Packages
-- ============================================================

-- 1 Month Basic
INSERT INTO packages (id, organization_id, name, name_ne, description, duration_days, price, features)
VALUES (
  'c0000000-0000-0000-0000-000000000001',
  'a0000000-0000-0000-0000-000000000001',
  '1 Month Basic',
  '१ महिना बेसिक',
  'Basic gym access for 1 month',
  30,
  2000.00,
  '["Gym access", "Locker"]'::jsonb
) ON CONFLICT DO NOTHING;

-- 3 Month Premium
INSERT INTO packages (id, organization_id, name, name_ne, description, duration_days, price, features)
VALUES (
  'c0000000-0000-0000-0000-000000000002',
  'a0000000-0000-0000-0000-000000000001',
  '3 Month Premium',
  '३ महिना प्रिमियम',
  'Premium gym access with trainer support for 3 months',
  90,
  5000.00,
  '["Gym access", "Locker", "Trainer support", "Diet plan"]'::jsonb
) ON CONFLICT DO NOTHING;

-- 6 Month Gold
INSERT INTO packages (id, organization_id, name, name_ne, description, duration_days, price, features)
VALUES (
  'c0000000-0000-0000-0000-000000000003',
  'a0000000-0000-0000-0000-000000000001',
  '6 Month Gold',
  '६ महिना गोल्ड',
  'Full access with personal training for 6 months',
  180,
  9000.00,
  '["Gym access", "Locker", "Personal trainer", "Diet plan", "Cardio zone", "Supplements discount"]'::jsonb
) ON CONFLICT DO NOTHING;

-- ============================================================
-- 5. Active Member Package for Hari (3 Month Premium, started 30 days ago)
-- ============================================================
INSERT INTO member_packages (id, user_id, package_id, organization_id, start_date, end_date, payment_method, amount_paid, status)
VALUES (
  'd0000000-0000-0000-0000-000000000001',
  'b0000000-0000-0000-0000-000000000003',
  'c0000000-0000-0000-0000-000000000002',
  'a0000000-0000-0000-0000-000000000001',
  CURRENT_DATE - INTERVAL '30 days',
  CURRENT_DATE - INTERVAL '30 days' + INTERVAL '90 days',
  'cash',
  5000.00,
  'active'
) ON CONFLICT DO NOTHING;

-- ============================================================
-- 6. Attendance Records for Hari (past 6 days + today = 7 records)
-- ============================================================
INSERT INTO attendance (id, user_id, organization_id, check_in_at, check_in_method)
VALUES
  ('e0000000-0000-0000-0000-000000000001', 'b0000000-0000-0000-0000-000000000003', 'a0000000-0000-0000-0000-000000000001',
    (CURRENT_DATE - INTERVAL '6 days') + TIME '07:15:00', 'qr'),
  ('e0000000-0000-0000-0000-000000000002', 'b0000000-0000-0000-0000-000000000003', 'a0000000-0000-0000-0000-000000000001',
    (CURRENT_DATE - INTERVAL '5 days') + TIME '06:45:00', 'nfc'),
  ('e0000000-0000-0000-0000-000000000003', 'b0000000-0000-0000-0000-000000000003', 'a0000000-0000-0000-0000-000000000001',
    (CURRENT_DATE - INTERVAL '4 days') + TIME '07:30:00', 'qr'),
  ('e0000000-0000-0000-0000-000000000004', 'b0000000-0000-0000-0000-000000000003', 'a0000000-0000-0000-0000-000000000001',
    (CURRENT_DATE - INTERVAL '3 days') + TIME '08:00:00', 'manual'),
  ('e0000000-0000-0000-0000-000000000005', 'b0000000-0000-0000-0000-000000000003', 'a0000000-0000-0000-0000-000000000001',
    (CURRENT_DATE - INTERVAL '2 days') + TIME '07:00:00', 'nfc'),
  ('e0000000-0000-0000-0000-000000000006', 'b0000000-0000-0000-0000-000000000003', 'a0000000-0000-0000-0000-000000000001',
    (CURRENT_DATE - INTERVAL '1 day')  + TIME '06:50:00', 'qr'),
  ('e0000000-0000-0000-0000-000000000007', 'b0000000-0000-0000-0000-000000000003', 'a0000000-0000-0000-0000-000000000001',
    CURRENT_DATE + TIME '07:10:00', 'nfc')
ON CONFLICT DO NOTHING;

-- ============================================================
-- 7. Health Metrics for Hari (2 BMI + 1 measurements)
-- ============================================================
INSERT INTO health_metrics (id, member_id, metric_type, value, recorded_at)
VALUES
  -- BMI record 1 (30 days ago)
  ('f0000000-0000-0000-0000-000000000001', 'b0000000-0000-0000-0000-000000000003', 'bmi',
    '{"bmi": 26.5, "weight": 78, "height": 171.5}'::jsonb,
    NOW() - INTERVAL '30 days'),
  -- BMI record 2 (recent)
  ('f0000000-0000-0000-0000-000000000002', 'b0000000-0000-0000-0000-000000000003', 'bmi',
    '{"bmi": 25.1, "weight": 74, "height": 171.5}'::jsonb,
    NOW() - INTERVAL '2 days'),
  -- Measurements record
  ('f0000000-0000-0000-0000-000000000003', 'b0000000-0000-0000-0000-000000000003', 'measurements',
    '{"chest": 96, "waist": 82, "hips": 94, "biceps": 33, "thighs": 55}'::jsonb,
    NOW() - INTERVAL '5 days')
ON CONFLICT DO NOTHING;

-- ============================================================
-- 8. Five Extra Members (for admin/staff list views + attendance today)
-- ============================================================

-- Extra member 1: Anita Gurung
INSERT INTO users (id, phone, name, name_ne, gender, date_of_birth)
VALUES ('b0000000-0000-0000-0000-000000000010', '9800000010', 'Anita Gurung', 'अनिता गुरुङ', 'female', '1995-01-20')
ON CONFLICT DO NOTHING;
INSERT INTO organization_members (user_id, organization_id, role, status)
VALUES ('b0000000-0000-0000-0000-000000000010', 'a0000000-0000-0000-0000-000000000001', 'member', 'active')
ON CONFLICT DO NOTHING;

-- Extra member 2: Bikash Thapa
INSERT INTO users (id, phone, name, name_ne, gender, date_of_birth)
VALUES ('b0000000-0000-0000-0000-000000000011', '9800000011', 'Bikash Thapa', 'बिकास थापा', 'male', '1993-06-10')
ON CONFLICT DO NOTHING;
INSERT INTO organization_members (user_id, organization_id, role, status)
VALUES ('b0000000-0000-0000-0000-000000000011', 'a0000000-0000-0000-0000-000000000001', 'member', 'active')
ON CONFLICT DO NOTHING;

-- Extra member 3: Chandra Rai
INSERT INTO users (id, phone, name, name_ne, gender, date_of_birth)
VALUES ('b0000000-0000-0000-0000-000000000012', '9800000012', 'Chandra Rai', 'चन्द्र राई', 'male', '1990-09-08')
ON CONFLICT DO NOTHING;
INSERT INTO organization_members (user_id, organization_id, role, status)
VALUES ('b0000000-0000-0000-0000-000000000012', 'a0000000-0000-0000-0000-000000000001', 'member', 'active')
ON CONFLICT DO NOTHING;

-- Extra member 4: Deepa Shrestha
INSERT INTO users (id, phone, name, name_ne, gender, date_of_birth)
VALUES ('b0000000-0000-0000-0000-000000000013', '9800000013', 'Deepa Shrestha', 'दीपा श्रेष्ठ', 'female', '1997-12-25')
ON CONFLICT DO NOTHING;
INSERT INTO organization_members (user_id, organization_id, role, status)
VALUES ('b0000000-0000-0000-0000-000000000013', 'a0000000-0000-0000-0000-000000000001', 'member', 'active')
ON CONFLICT DO NOTHING;

-- Extra member 5: Ekaraj Magar
INSERT INTO users (id, phone, name, name_ne, gender, date_of_birth)
VALUES ('b0000000-0000-0000-0000-000000000014', '9800000014', 'Ekaraj Magar', 'एकराज मगर', 'male', '1991-04-14')
ON CONFLICT DO NOTHING;
INSERT INTO organization_members (user_id, organization_id, role, status)
VALUES ('b0000000-0000-0000-0000-000000000014', 'a0000000-0000-0000-0000-000000000001', 'member', 'active')
ON CONFLICT DO NOTHING;

-- Attendance today for all 5 extra members
INSERT INTO attendance (id, user_id, organization_id, check_in_at, check_in_method)
VALUES
  ('e0000000-0000-0000-0000-000000000010', 'b0000000-0000-0000-0000-000000000010', 'a0000000-0000-0000-0000-000000000001',
    CURRENT_DATE + TIME '06:30:00', 'qr'),
  ('e0000000-0000-0000-0000-000000000011', 'b0000000-0000-0000-0000-000000000011', 'a0000000-0000-0000-0000-000000000001',
    CURRENT_DATE + TIME '07:00:00', 'nfc'),
  ('e0000000-0000-0000-0000-000000000012', 'b0000000-0000-0000-0000-000000000012', 'a0000000-0000-0000-0000-000000000001',
    CURRENT_DATE + TIME '07:45:00', 'manual'),
  ('e0000000-0000-0000-0000-000000000013', 'b0000000-0000-0000-0000-000000000013', 'a0000000-0000-0000-0000-000000000001',
    CURRENT_DATE + TIME '08:15:00', 'qr'),
  ('e0000000-0000-0000-0000-000000000014', 'b0000000-0000-0000-0000-000000000014', 'a0000000-0000-0000-0000-000000000001',
    CURRENT_DATE + TIME '08:30:00', 'nfc')
ON CONFLICT DO NOTHING;

COMMIT;
