# E2E UX Flow Testing — Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Build member-facing web pages, add dev OTP bypass, seed test data, then do interactive Playwright MCP walkthrough of all 3 user roles.

**Architecture:** Add `DEV_OTP_BYPASS` to Go auth service for test phones (9800*). Add member API methods to the existing `ApiClient`. Create `/[lang]/member/` Next.js pages with their own layout/sidebar using the same dark cinematic design. Seed test users (admin/staff/member) + org + sample data via SQL. Add role-based redirect from login. Drive all flows via Playwright MCP.

**Tech Stack:** Go (auth service), Next.js 14 + React 18 + TypeScript + Tailwind CSS (member pages), PostgreSQL (seed data), Playwright MCP (E2E)

---

## Task 1: Dev OTP Bypass

**Files:**
- Modify: `internal/config/config.go` (add `DevOTPBypass` field)
- Modify: `internal/auth/service.go` (bypass logic in `RequestOTP` and `VerifyOTP`)

**Step 1: Add config field**

In `internal/config/config.go`, add to `AuthConfig` struct:

```go
DevOTPBypass bool
```

In the `Load()` function, after the bcrypt cost loading (line ~153), add:

```go
cfg.Auth.DevOTPBypass = envOrDefault("DEV_OTP_BYPASS", "false") == "true"
```

**Step 2: Bypass in RequestOTP**

In `internal/auth/service.go`, in `RequestOTP()`, after the phone normalization and validation (line ~56), add before rate limit check:

```go
// Dev bypass: skip SMS for test phones (9800*)
if s.cfg.DevOTPBypass && strings.HasPrefix(phone, "9800") {
	// Store a known hash for "123456"
	hash, err := bcrypt.GenerateFromPassword([]byte("123456"), s.cfg.BcryptCost)
	if err != nil {
		return fmt.Errorf("hashing dev OTP: %w", err)
	}
	if err := s.repo.StoreOTPHash(ctx, phone, string(hash), s.cfg.OTPExpiry); err != nil {
		return fmt.Errorf("storing dev OTP: %w", err)
	}
	s.logger.Info("dev OTP bypass: stored test OTP", "phone", phone[:4]+"******")
	return nil
}
```

Add `"strings"` to imports.

**Step 3: Update .env and restart API**

Add to `.env`:

```
DEV_OTP_BYPASS=true
```

Rebuild and restart the API server:

```bash
DB_PASSWORD=urja_secret DEV_OTP_BYPASS=true go run ./cmd/api &
```

**Step 4: Verify bypass works**

```bash
# Request OTP
curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"phone":"9800000001"}'
# Expected: {"message":"OTP sent successfully"}

# Verify OTP
curl -s -X POST http://localhost:8080/api/v1/auth/verify-otp \
  -H "Content-Type: application/json" \
  -d '{"phone":"9800000001","otp":"123456"}'
# Expected: {"access_token":"...","refresh_token":"...","token_type":"bearer"}
```

**Step 5: Commit**

```bash
git add internal/config/config.go internal/auth/service.go
git commit -m "feat: add dev OTP bypass for E2E testing"
```

---

## Task 2: Seed Test Data

**Files:**
- Create: `db/seed_e2e.sql`

**Step 1: Write seed SQL**

Create `db/seed_e2e.sql` with:

```sql
-- E2E Test Seed Data
-- Run: psql -U urja -d urja -f db/seed_e2e.sql
-- Test phones (9800*) + OTP 123456 when DEV_OTP_BYPASS=true

BEGIN;

-- 1. Create test org
INSERT INTO organizations (id, name, name_ne, slug, description, phone, email, address)
VALUES (
  'a0000000-0000-0000-0000-000000000001',
  'Kathmandu Fitness Hub',
  'काठमाडौं फिटनेस हब',
  'kathmandu-fitness-hub',
  'Premium gym in the heart of Kathmandu',
  '9800000000',
  'info@ktmfitness.com',
  'Thamel, Kathmandu'
) ON CONFLICT (slug) DO NOTHING;

-- 2. Create test users
-- Admin (gym owner)
INSERT INTO users (id, phone, name, name_ne, email, gender)
VALUES (
  'b0000000-0000-0000-0000-000000000001',
  '9800000001',
  'Rajesh Sharma',
  'राजेश शर्मा',
  'rajesh@ktmfitness.com',
  'male'
) ON CONFLICT (phone) DO NOTHING;

-- Staff
INSERT INTO users (id, phone, name, name_ne, email, gender)
VALUES (
  'b0000000-0000-0000-0000-000000000002',
  '9800000002',
  'Sita Tamang',
  'सिता तामाङ',
  'sita@ktmfitness.com',
  'female'
) ON CONFLICT (phone) DO NOTHING;

-- Member
INSERT INTO users (id, phone, name, name_ne, email, gender, date_of_birth)
VALUES (
  'b0000000-0000-0000-0000-000000000003',
  '9800000003',
  'Hari Prasad',
  'हरि प्रसाद',
  'hari@gmail.com',
  'male',
  '1995-06-15'
) ON CONFLICT (phone) DO NOTHING;

-- 3. Create org memberships
INSERT INTO organization_members (user_id, organization_id, role, staff_role)
VALUES
  ('b0000000-0000-0000-0000-000000000001', 'a0000000-0000-0000-0000-000000000001', 'admin', 'owner'),
  ('b0000000-0000-0000-0000-000000000002', 'a0000000-0000-0000-0000-000000000001', 'staff', 'trainer'),
  ('b0000000-0000-0000-0000-000000000003', 'a0000000-0000-0000-0000-000000000001', 'member', NULL)
ON CONFLICT (user_id, organization_id) DO NOTHING;

-- 4. Create packages
INSERT INTO packages (id, organization_id, name, name_ne, description, duration_days, price, features)
VALUES
  ('c0000000-0000-0000-0000-000000000001', 'a0000000-0000-0000-0000-000000000001',
   '1 Month Basic', '१ महिना बेसिक', 'Gym access with basic equipment', 30, 2000.00, '["gym_access","locker"]'),
  ('c0000000-0000-0000-0000-000000000002', 'a0000000-0000-0000-0000-000000000001',
   '3 Month Premium', '३ महिना प्रिमियम', 'Full access with trainer', 90, 5000.00, '["gym_access","trainer","locker","pool"]'),
  ('c0000000-0000-0000-0000-000000000003', 'a0000000-0000-0000-0000-000000000001',
   '6 Month Gold', '६ महिना गोल्ड', 'All-inclusive membership', 180, 9000.00, '["gym_access","trainer","locker","pool","sauna","supplements"]')
ON CONFLICT DO NOTHING;

-- 5. Give member an active package subscription
INSERT INTO member_packages (id, user_id, package_id, organization_id, start_date, end_date, payment_method, amount_paid, status)
VALUES (
  'd0000000-0000-0000-0000-000000000001',
  'b0000000-0000-0000-0000-000000000003',
  'c0000000-0000-0000-0000-000000000002',
  'a0000000-0000-0000-0000-000000000001',
  CURRENT_DATE - INTERVAL '30 days',
  CURRENT_DATE + INTERVAL '60 days',
  'cash',
  5000.00,
  'active'
) ON CONFLICT DO NOTHING;

-- 6. Seed attendance records for the member (past 7 days)
INSERT INTO attendance (user_id, organization_id, check_in_at, check_in_method)
VALUES
  ('b0000000-0000-0000-0000-000000000003', 'a0000000-0000-0000-0000-000000000001', NOW() - INTERVAL '6 days' + TIME '06:30', 'qr'),
  ('b0000000-0000-0000-0000-000000000003', 'a0000000-0000-0000-0000-000000000001', NOW() - INTERVAL '5 days' + TIME '07:15', 'nfc'),
  ('b0000000-0000-0000-0000-000000000003', 'a0000000-0000-0000-0000-000000000001', NOW() - INTERVAL '4 days' + TIME '06:45', 'qr'),
  ('b0000000-0000-0000-0000-000000000003', 'a0000000-0000-0000-0000-000000000001', NOW() - INTERVAL '3 days' + TIME '07:00', 'manual'),
  ('b0000000-0000-0000-0000-000000000003', 'a0000000-0000-0000-0000-000000000001', NOW() - INTERVAL '1 day' + TIME '06:30', 'qr'),
  ('b0000000-0000-0000-0000-000000000003', 'a0000000-0000-0000-0000-000000000001', NOW() + TIME '07:00', 'nfc');

-- 7. Seed health metrics for the member
INSERT INTO health_metrics (id, user_id, metric_type, value, recorded_at)
VALUES
  ('e0000000-0000-0000-0000-000000000001', 'b0000000-0000-0000-0000-000000000003', 'bmi',
   '{"height_cm": 175, "weight_kg": 78, "bmi": 25.5, "category": "overweight"}',
   NOW() - INTERVAL '30 days'),
  ('e0000000-0000-0000-0000-000000000002', 'b0000000-0000-0000-0000-000000000003', 'bmi',
   '{"height_cm": 175, "weight_kg": 75, "bmi": 24.5, "category": "normal"}',
   NOW() - INTERVAL '7 days'),
  ('e0000000-0000-0000-0000-000000000003', 'b0000000-0000-0000-0000-000000000003', 'measurements',
   '{"chest_cm": 95, "waist_cm": 85, "hips_cm": 100, "arms_cm": 32, "thighs_cm": 58, "whr": 0.85}',
   NOW() - INTERVAL '7 days');

-- 8. Seed a few more members for the admin/staff views
INSERT INTO users (id, phone, name, name_ne, email, gender) VALUES
  ('b0000000-0000-0000-0000-000000000004', '9841111111', 'Bikash Gurung', 'बिकास गुरुङ', NULL, 'male'),
  ('b0000000-0000-0000-0000-000000000005', '9841111112', 'Anita Rai', 'अनिता राई', NULL, 'female'),
  ('b0000000-0000-0000-0000-000000000006', '9841111113', 'Deepa Thapa', 'दीपा थापा', NULL, 'female'),
  ('b0000000-0000-0000-0000-000000000007', '9841111114', 'Sujan KC', 'सुजन केसी', NULL, 'male'),
  ('b0000000-0000-0000-0000-000000000008', '9841111115', 'Priya Adhikari', 'प्रिया अधिकारी', NULL, 'female')
ON CONFLICT (phone) DO NOTHING;

INSERT INTO organization_members (user_id, organization_id, role)
VALUES
  ('b0000000-0000-0000-0000-000000000004', 'a0000000-0000-0000-0000-000000000001', 'member'),
  ('b0000000-0000-0000-0000-000000000005', 'a0000000-0000-0000-0000-000000000001', 'member'),
  ('b0000000-0000-0000-0000-000000000006', 'a0000000-0000-0000-0000-000000000001', 'member'),
  ('b0000000-0000-0000-0000-000000000007', 'a0000000-0000-0000-0000-000000000001', 'member'),
  ('b0000000-0000-0000-0000-000000000008', 'a0000000-0000-0000-0000-000000000001', 'member')
ON CONFLICT (user_id, organization_id) DO NOTHING;

-- Add attendance for extra members (today)
INSERT INTO attendance (user_id, organization_id, check_in_at, check_in_method)
VALUES
  ('b0000000-0000-0000-0000-000000000004', 'a0000000-0000-0000-0000-000000000001', NOW() + TIME '06:15', 'qr'),
  ('b0000000-0000-0000-0000-000000000005', 'a0000000-0000-0000-0000-000000000001', NOW() + TIME '06:42', 'nfc'),
  ('b0000000-0000-0000-0000-000000000007', 'a0000000-0000-0000-0000-000000000001', NOW() + TIME '07:01', 'manual');

COMMIT;
```

**Step 2: Run the seed**

```bash
PGPASSWORD=urja_secret psql -h localhost -U urja -d urja -f db/seed_e2e.sql
```

**Step 3: Verify seed data**

```bash
PGPASSWORD=urja_secret psql -h localhost -U urja -d urja -c "SELECT u.name, om.role, om.staff_role FROM users u JOIN organization_members om ON u.id = om.user_id WHERE u.phone LIKE '9800%' ORDER BY u.phone;"
```

Expected:

```
      name       | role  | staff_role
-----------------+-------+------------
 Rajesh Sharma   | admin | owner
 Sita Tamang     | staff | trainer
 Hari Prasad     | member|
```

**Step 4: Commit**

```bash
git add db/seed_e2e.sql
git commit -m "feat: add E2E test seed data"
```

---

## Task 3: Add Member API Methods to ApiClient

**Files:**
- Modify: `web/src/types/index.ts` (add member types)
- Modify: `web/src/lib/api.ts` (add member endpoints)

**Step 1: Add member types to `web/src/types/index.ts`**

Append after the existing `StaffList` interface:

```typescript
// --- Member Profile (from GET /members/me) ---

export interface MemberProfile {
  id: string;
  phone: string;
  name: string;
  name_ne: string | null;
  email: string | null;
  date_of_birth: string | null;
  gender: string | null;
  avatar_url: string | null;
  emergency_contact_name: string | null;
  emergency_contact_phone: string | null;
  privacy_settings: PrivacySettings;
  organizations: OrgMembership[];
  created_at: string;
  updated_at: string;
}

export interface PrivacySettings {
  show_email: boolean;
  show_phone: boolean;
  show_profile: boolean;
  show_attendance: boolean;
  show_on_leaderboard: boolean;
}

export interface OrgMembership {
  org_id: string;
  org_name: string;
  role: string;
  status: string;
  joined_at: string;
}

// --- Member Attendance ---

export interface MemberAttendanceRecord {
  id: string;
  user_id: string;
  org_id: string;
  check_in_at: string;
  method: "qr" | "nfc" | "manual";
}

export interface MemberStreak {
  id: string;
  member_id: string;
  org_id: string;
  current_streak: number;
  longest_streak: number;
  last_check_in: string | null;
  updated_at: string;
}

// --- Member Packages ---

export interface MemberPackage {
  id: string;
  user_id: string;
  package_id: string;
  organization_id: string;
  start_date: string;
  end_date: string;
  payment_method: string | null;
  payment_reference: string | null;
  amount_paid: string;
  status: "pending" | "active" | "expired" | "cancelled";
  created_at: string;
  package_name: string | null;
  org_name: string | null;
}

// --- Health Metrics ---

export interface HealthMetric {
  id: string;
  member_id: string;
  metric_type: string;
  value: Record<string, unknown>;
  recorded_at: string;
}

// --- Workout ---

export interface WorkoutLog {
  id: string;
  user_id: string;
  workout_template_id: string | null;
  organization_id: string;
  exercises: unknown[];
  duration_minutes: number | null;
  notes: string;
  logged_at: string;
}

export interface WorkoutPlan {
  id: string;
  user_id: string;
  organization_id: string;
  workout_template_id: string;
  assigned_by: string | null;
  assigned_at: string;
  template: WorkoutTemplate | null;
}

export interface WorkoutTemplate {
  id: string;
  name: string;
  name_ne: string | null;
  description: string | null;
  category: string | null;
  difficulty: string | null;
  duration_minutes: number | null;
  exercises: unknown[];
}
```

**Step 2: Add member API methods to `web/src/lib/api.ts`**

Add new imports to the import block at top:

```typescript
import type {
  // ... existing imports ...
  MemberProfile,
  MemberAttendanceRecord,
  MemberStreak,
  MemberPackage,
  HealthMetric,
  WorkoutLog,
  WorkoutPlan,
} from "@/types";
```

Add methods to the `ApiClient` class after the staff section:

```typescript
  // --- Member (self) ---

  async getMyProfile(): Promise<MemberProfile> {
    return this.request("/api/v1/members/me");
  }

  async getMyAttendance(
    params: { limit?: number; offset?: number } = {}
  ): Promise<{ data: MemberAttendanceRecord[] }> {
    const q = new URLSearchParams();
    if (params.limit) q.set("limit", String(params.limit));
    if (params.offset) q.set("offset", String(params.offset));
    const qs = q.toString();
    return this.request(`/api/v1/members/me/attendance${qs ? `?${qs}` : ""}`);
  }

  async getMyStreaks(): Promise<{ data: MemberStreak[] }> {
    return this.request("/api/v1/members/me/streaks");
  }

  async getMyPackages(
    params: { limit?: number; offset?: number } = {}
  ): Promise<{ data: MemberPackage[] }> {
    const q = new URLSearchParams();
    if (params.limit) q.set("limit", String(params.limit));
    if (params.offset) q.set("offset", String(params.offset));
    const qs = q.toString();
    return this.request(`/api/v1/members/me/packages${qs ? `?${qs}` : ""}`);
  }

  async getMyHealth(
    params: { type?: string; limit?: number } = {}
  ): Promise<{ data: HealthMetric[] }> {
    const q = new URLSearchParams();
    if (params.type) q.set("type", params.type);
    if (params.limit) q.set("limit", String(params.limit));
    const qs = q.toString();
    return this.request(`/api/v1/members/me/health${qs ? `?${qs}` : ""}`);
  }

  async getMyWorkoutLogs(
    params: { organization_id?: string; limit?: number; offset?: number } = {}
  ): Promise<{ data: WorkoutLog[] }> {
    const q = new URLSearchParams();
    if (params.organization_id) q.set("organization_id", params.organization_id);
    if (params.limit) q.set("limit", String(params.limit));
    if (params.offset) q.set("offset", String(params.offset));
    const qs = q.toString();
    return this.request(`/api/v1/members/me/workout-logs${qs ? `?${qs}` : ""}`);
  }

  async getMyWorkoutPlan(orgId: string): Promise<WorkoutPlan> {
    return this.request(`/api/v1/members/me/workout-plan?organization_id=${orgId}`);
  }
```

**Step 3: Commit**

```bash
git add web/src/types/index.ts web/src/lib/api.ts
git commit -m "feat: add member-facing API client methods and types"
```

---

## Task 4: Add Member i18n Strings

**Files:**
- Modify: `web/src/lib/i18n.ts`

**Step 1: Add member nav and page strings**

In the `en` dictionary, add a new `memberNav` section and `memberPages` section after the `topbar` key:

```typescript
    memberNav: {
      myDashboard: "My Dashboard",
      myAttendance: "My Attendance",
      myPackages: "My Packages",
      myHealth: "My Health",
      myWorkouts: "My Workouts",
      profile: "Profile",
    },
    memberPages: {
      welcomeBack: "Welcome back",
      memberSince: "Member since",
      currentPackage: "Current Package",
      noActivePackage: "No active package",
      daysRemaining: "days remaining",
      expired: "Expired",
      attendance: "Attendance",
      currentStreak: "Current Streak",
      longestStreak: "Longest Streak",
      days: "days",
      recentCheckIns: "Recent Check-ins",
      noCheckIns: "No check-ins yet",
      packages: "My Packages",
      activePackages: "Active",
      expiredPackages: "Past",
      noPackages: "No packages found",
      health: "Health Metrics",
      bmi: "BMI",
      measurements: "Body Measurements",
      noHealthData: "No health data recorded yet",
      workouts: "My Workouts",
      workoutPlan: "Assigned Plan",
      recentWorkouts: "Recent Workouts",
      noWorkouts: "No workouts logged yet",
      noPlan: "No workout plan assigned",
      duration: "Duration",
      minutes: "min",
    },
```

Add equivalent Nepali translations in the `ne` dictionary:

```typescript
    memberNav: {
      myDashboard: "मेरो ड्यासबोर्ड",
      myAttendance: "मेरो उपस्थिति",
      myPackages: "मेरा प्याकेजहरू",
      myHealth: "मेरो स्वास्थ्य",
      myWorkouts: "मेरा व्यायामहरू",
      profile: "प्रोफाइल",
    },
    memberPages: {
      welcomeBack: "स्वागत छ",
      memberSince: "सदस्य मिति",
      currentPackage: "हालको प्याकेज",
      noActivePackage: "कुनै सक्रिय प्याकेज छैन",
      daysRemaining: "दिन बाँकी",
      expired: "म्याद सकिएको",
      attendance: "उपस्थिति",
      currentStreak: "हालको स्ट्रिक",
      longestStreak: "सबैभन्दा लामो स्ट्रिक",
      days: "दिन",
      recentCheckIns: "हालको चेक-इनहरू",
      noCheckIns: "अहिलेसम्म चेक-इन छैन",
      packages: "मेरा प्याकेजहरू",
      activePackages: "सक्रिय",
      expiredPackages: "विगत",
      noPackages: "कुनै प्याकेज भेटिएन",
      health: "स्वास्थ्य मेट्रिक्स",
      bmi: "BMI",
      measurements: "शरीर नाप",
      noHealthData: "अहिलेसम्म कुनै स्वास्थ्य डाटा छैन",
      workouts: "मेरा व्यायामहरू",
      workoutPlan: "तोकिएको योजना",
      recentWorkouts: "हालको व्यायामहरू",
      noWorkouts: "अहिलेसम्म कुनै व्यायाम छैन",
      noPlan: "कुनै व्यायाम योजना तोकिएको छैन",
      duration: "अवधि",
      minutes: "मिनेट",
    },
```

**Step 2: Commit**

```bash
git add web/src/lib/i18n.ts
git commit -m "feat: add member page i18n strings (en + ne)"
```

---

## Task 5: Role-Based Login Redirect

**Files:**
- Modify: `web/src/app/[lang]/login/page.tsx`
- Modify: `web/src/app/[lang]/page.tsx`

**Step 1: Update login page redirect**

In `web/src/app/[lang]/login/page.tsx`, modify the `handleVerifyOtp` function to route based on role from the JWT:

Replace the line `router.replace(\`/${locale}/dashboard\`);` (line ~90) with:

```typescript
        // Route based on role from JWT
        const token = localStorage.getItem("access_token");
        if (token) {
          const payload = JSON.parse(atob(token.split(".")[1]));
          if (payload.role === "member") {
            router.replace(`/${locale}/member`);
          } else {
            router.replace(`/${locale}/dashboard`);
          }
        } else {
          router.replace(`/${locale}/dashboard`);
        }
```

Also update the `useEffect` redirect for already-authenticated users (lines ~32-35). Replace the `router.replace` line with the same role-based logic:

```typescript
  useEffect(() => {
    if (!authLoading && isAuthenticated) {
      const token = localStorage.getItem("access_token");
      if (token) {
        try {
          const payload = JSON.parse(atob(token.split(".")[1]));
          if (payload.role === "member") {
            router.replace(`/${locale}/member`);
            return;
          }
        } catch {}
      }
      router.replace(`/${locale}/dashboard`);
    }
  }, [isAuthenticated, authLoading, router, locale]);
```

**Step 2: Update root page redirect**

In `web/src/app/[lang]/page.tsx`, this currently redirects everyone to dashboard. Keep it as-is since the AuthGuard/login flow will handle routing.

**Step 3: Commit**

```bash
git add web/src/app/[lang]/login/page.tsx
git commit -m "feat: role-based login redirect (members -> /member)"
```

---

## Task 6: Member Layout + Sidebar

**Files:**
- Create: `web/src/app/[lang]/member/layout.tsx`
- Create: `web/src/components/layout/MemberSidebar.tsx`

**Step 1: Create MemberSidebar component**

Create `web/src/components/layout/MemberSidebar.tsx`. This follows the exact same pattern as the existing `Sidebar.tsx` but with member-specific nav items:

```typescript
"use client";

import { usePathname } from "next/navigation";
import Link from "next/link";
import type { Dictionary } from "@/lib/i18n";
import type { Locale } from "@/types";
import {
  Home,
  CalendarCheck,
  Package,
  Heart,
  Dumbbell,
  X,
} from "lucide-react";

interface NavItem {
  label: string;
  href: string;
  icon: React.ReactNode;
}

function buildNavItems(t: Dictionary, locale: Locale): NavItem[] {
  const base = `/${locale}/member`;
  return [
    { label: t.memberNav.myDashboard, href: base, icon: <Home className="w-4 h-4" /> },
    { label: t.memberNav.myAttendance, href: `${base}/attendance`, icon: <CalendarCheck className="w-4 h-4" /> },
    { label: t.memberNav.myPackages, href: `${base}/packages`, icon: <Package className="w-4 h-4" /> },
    { label: t.memberNav.myHealth, href: `${base}/health`, icon: <Heart className="w-4 h-4" /> },
    { label: t.memberNav.myWorkouts, href: `${base}/workouts`, icon: <Dumbbell className="w-4 h-4" /> },
  ];
}

interface MemberSidebarProps {
  t: Dictionary;
  locale: Locale;
  isOpen: boolean;
  onClose: () => void;
}

export function MemberSidebar({ t, locale, isOpen, onClose }: MemberSidebarProps) {
  const pathname = usePathname();
  const items = buildNavItems(t, locale);

  const isActive = (href: string) => {
    if (href === `/${locale}/member`) return pathname === href;
    return pathname.startsWith(href);
  };

  return (
    <>
      {isOpen && (
        <div
          className="fixed inset-0 bg-black/60 backdrop-blur-sm z-40 lg:hidden"
          onClick={onClose}
        />
      )}

      <aside
        className={`fixed top-0 left-0 h-full w-64 bg-bg-elevated border-r border-white/[0.06] z-50 flex flex-col transition-transform duration-300 lg:translate-x-0 ${
          isOpen ? "translate-x-0" : "-translate-x-full"
        }`}
      >
        <div className="flex items-center justify-between p-4 border-b border-white/[0.06]">
          <div className="flex items-center gap-3">
            <div className="w-8 h-8 rounded-xl bg-accent/20 flex items-center justify-center">
              <Dumbbell className="w-4 h-4 text-accent" />
            </div>
            <div>
              <h2 className="text-sm font-semibold text-fg truncate">
                {t.common.appName}
              </h2>
              <p className="text-xs text-fg-muted truncate">Member Portal</p>
            </div>
          </div>
          <button
            onClick={onClose}
            className="lg:hidden p-1 rounded-lg hover:bg-surface transition-colors"
          >
            <X className="w-4 h-4 text-fg-muted" />
          </button>
        </div>

        <nav className="flex-1 overflow-y-auto py-3 px-3 space-y-0.5">
          {items.map((item) => {
            const active = isActive(item.href);
            return (
              <Link
                key={item.href}
                href={item.href}
                onClick={onClose}
                className={`flex items-center gap-3 px-3 py-2 rounded-lg text-sm transition-all duration-200 ${
                  active
                    ? "bg-accent/10 text-accent font-medium"
                    : "text-fg-muted hover:bg-surface hover:text-fg"
                }`}
              >
                <span className={active ? "text-accent" : ""}>
                  {item.icon}
                </span>
                {item.label}
              </Link>
            );
          })}
        </nav>
      </aside>
    </>
  );
}
```

**Step 2: Create member layout**

Create `web/src/app/[lang]/member/layout.tsx`:

```typescript
"use client";

import { useState } from "react";
import { MemberSidebar } from "@/components/layout/MemberSidebar";
import { TopBar } from "@/components/layout/TopBar";
import { AuthGuard } from "@/components/layout/AuthGuard";
import { getDictionary } from "@/lib/i18n";
import type { Locale } from "@/types";

export default function MemberLayout({
  children,
  params,
}: {
  children: React.ReactNode;
  params: { lang: string };
}) {
  const locale = params.lang as Locale;
  const t = getDictionary(locale);
  const [sidebarOpen, setSidebarOpen] = useState(false);

  return (
    <AuthGuard locale={locale}>
      <div className="min-h-screen bg-bg-base">
        <MemberSidebar
          t={t}
          locale={locale}
          isOpen={sidebarOpen}
          onClose={() => setSidebarOpen(false)}
        />
        <div className="lg:ml-64">
          <TopBar
            t={t}
            locale={locale}
            onMenuToggle={() => setSidebarOpen(!sidebarOpen)}
          />
          <main className="p-4 lg:p-6">{children}</main>
        </div>
      </div>
    </AuthGuard>
  );
}
```

**Step 3: Commit**

```bash
git add web/src/components/layout/MemberSidebar.tsx web/src/app/[lang]/member/layout.tsx
git commit -m "feat: add member layout with sidebar"
```

---

## Task 7: Member Dashboard Page

**Files:**
- Create: `web/src/app/[lang]/member/page.tsx`

**Step 1: Create the member dashboard**

Create `web/src/app/[lang]/member/page.tsx`:

```typescript
"use client";

import { useEffect, useState } from "react";
import { getDictionary } from "@/lib/i18n";
import { useAuth } from "@/lib/auth";
import { api } from "@/lib/api";
import type { Locale, MemberProfile, MemberStreak, MemberPackage } from "@/types";
import { User, CalendarCheck, Package, TrendingUp, Loader2 } from "lucide-react";

export default function MemberDashboard({ params }: { params: { lang: string } }) {
  const locale = params.lang as Locale;
  const t = getDictionary(locale);
  const { user } = useAuth();

  const [profile, setProfile] = useState<MemberProfile | null>(null);
  const [streaks, setStreaks] = useState<MemberStreak[]>([]);
  const [packages, setPackages] = useState<MemberPackage[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    async function load() {
      try {
        const [p, s, pkg] = await Promise.all([
          api.getMyProfile(),
          api.getMyStreaks().catch(() => ({ data: [] })),
          api.getMyPackages({ limit: 5 }).catch(() => ({ data: [] })),
        ]);
        setProfile(p);
        setStreaks(s.data ?? []);
        setPackages(pkg.data ?? []);
      } catch (err) {
        console.error("Failed to load member data:", err);
      } finally {
        setLoading(false);
      }
    }
    load();
  }, []);

  if (loading) {
    return (
      <div className="flex items-center justify-center py-20">
        <Loader2 className="w-8 h-8 text-accent animate-spin" />
      </div>
    );
  }

  const activePackage = packages.find((p) => p.status === "active");
  const streak = streaks[0];

  const daysRemaining = activePackage
    ? Math.max(0, Math.ceil((new Date(activePackage.end_date).getTime() - Date.now()) / 86400000))
    : 0;

  return (
    <div className="space-y-6">
      {/* Welcome header */}
      <div>
        <h1 className="text-2xl font-semibold text-fg">
          {t.memberPages.welcomeBack}, {profile?.name ?? user?.phone}
        </h1>
        {profile?.organizations?.[0] && (
          <p className="text-fg-muted text-sm mt-1">
            {profile.organizations[0].org_name} &middot; {t.memberPages.memberSince}{" "}
            {new Date(profile.organizations[0].joined_at).toLocaleDateString()}
          </p>
        )}
      </div>

      {/* Stat cards */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        {/* Current Package */}
        <div className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl p-5">
          <div className="flex items-start justify-between">
            <div className="w-10 h-10 rounded-xl bg-accent/10 flex items-center justify-center">
              <Package className="w-5 h-5 text-accent" />
            </div>
          </div>
          <p className="mt-3 text-2xl font-semibold text-fg">
            {activePackage?.package_name ?? t.memberPages.noActivePackage}
          </p>
          <p className="text-xs text-fg-muted mt-1">
            {activePackage ? `${daysRemaining} ${t.memberPages.daysRemaining}` : ""}
          </p>
        </div>

        {/* Current Streak */}
        <div className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl p-5">
          <div className="w-10 h-10 rounded-xl bg-accent/10 flex items-center justify-center">
            <TrendingUp className="w-5 h-5 text-accent" />
          </div>
          <p className="mt-3 text-2xl font-semibold text-fg">
            {streak?.current_streak ?? 0} {t.memberPages.days}
          </p>
          <p className="text-xs text-fg-muted mt-1">{t.memberPages.currentStreak}</p>
        </div>

        {/* Longest Streak */}
        <div className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl p-5">
          <div className="w-10 h-10 rounded-xl bg-accent/10 flex items-center justify-center">
            <CalendarCheck className="w-5 h-5 text-accent" />
          </div>
          <p className="mt-3 text-2xl font-semibold text-fg">
            {streak?.longest_streak ?? 0} {t.memberPages.days}
          </p>
          <p className="text-xs text-fg-muted mt-1">{t.memberPages.longestStreak}</p>
        </div>

        {/* Profile */}
        <div className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl p-5">
          <div className="w-10 h-10 rounded-xl bg-accent/10 flex items-center justify-center">
            <User className="w-5 h-5 text-accent" />
          </div>
          <p className="mt-3 text-lg font-semibold text-fg truncate">
            {profile?.name}
          </p>
          <p className="text-xs text-fg-muted mt-1">{profile?.phone}</p>
        </div>
      </div>
    </div>
  );
}
```

**Step 2: Commit**

```bash
git add web/src/app/[lang]/member/page.tsx
git commit -m "feat: add member dashboard page"
```

---

## Task 8: Member Attendance Page

**Files:**
- Create: `web/src/app/[lang]/member/attendance/page.tsx`

**Step 1: Create the attendance page**

Create `web/src/app/[lang]/member/attendance/page.tsx`:

```typescript
"use client";

import { useEffect, useState } from "react";
import { getDictionary } from "@/lib/i18n";
import { api } from "@/lib/api";
import type { Locale, MemberAttendanceRecord, MemberStreak } from "@/types";
import { CalendarCheck, TrendingUp, Loader2 } from "lucide-react";

function MethodBadge({ method }: { method: string }) {
  const colors: Record<string, string> = {
    qr: "bg-blue-500/10 text-blue-400 border-blue-500/20",
    nfc: "bg-purple-500/10 text-purple-400 border-purple-500/20",
    manual: "bg-gray-500/10 text-gray-400 border-gray-500/20",
  };
  return (
    <span className={`inline-flex items-center px-2 py-0.5 rounded-full text-[11px] font-mono uppercase border ${colors[method] ?? colors.manual}`}>
      {method}
    </span>
  );
}

export default function MemberAttendancePage({ params }: { params: { lang: string } }) {
  const locale = params.lang as Locale;
  const t = getDictionary(locale);

  const [records, setRecords] = useState<MemberAttendanceRecord[]>([]);
  const [streaks, setStreaks] = useState<MemberStreak[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    async function load() {
      try {
        const [att, str] = await Promise.all([
          api.getMyAttendance({ limit: 50 }),
          api.getMyStreaks().catch(() => ({ data: [] })),
        ]);
        setRecords(att.data ?? []);
        setStreaks(str.data ?? []);
      } catch (err) {
        console.error("Failed to load attendance:", err);
      } finally {
        setLoading(false);
      }
    }
    load();
  }, []);

  if (loading) {
    return (
      <div className="flex items-center justify-center py-20">
        <Loader2 className="w-8 h-8 text-accent animate-spin" />
      </div>
    );
  }

  const streak = streaks[0];

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-semibold text-fg">{t.memberPages.attendance}</h1>

      {/* Streak cards */}
      <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
        <div className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl p-5">
          <div className="flex items-center gap-3">
            <div className="w-10 h-10 rounded-xl bg-accent/10 flex items-center justify-center">
              <TrendingUp className="w-5 h-5 text-accent" />
            </div>
            <div>
              <p className="text-2xl font-semibold text-fg">{streak?.current_streak ?? 0} {t.memberPages.days}</p>
              <p className="text-xs text-fg-muted">{t.memberPages.currentStreak}</p>
            </div>
          </div>
        </div>
        <div className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl p-5">
          <div className="flex items-center gap-3">
            <div className="w-10 h-10 rounded-xl bg-accent/10 flex items-center justify-center">
              <CalendarCheck className="w-5 h-5 text-accent" />
            </div>
            <div>
              <p className="text-2xl font-semibold text-fg">{streak?.longest_streak ?? 0} {t.memberPages.days}</p>
              <p className="text-xs text-fg-muted">{t.memberPages.longestStreak}</p>
            </div>
          </div>
        </div>
      </div>

      {/* Check-in history table */}
      <div className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl overflow-hidden">
        <div className="px-5 py-4 border-b border-white/[0.06]">
          <h2 className="text-sm font-semibold text-fg">{t.memberPages.recentCheckIns}</h2>
        </div>
        {records.length === 0 ? (
          <p className="px-5 py-8 text-center text-fg-muted text-sm">{t.memberPages.noCheckIns}</p>
        ) : (
          <div className="divide-y divide-white/[0.04]">
            {records.map((r) => (
              <div key={r.id} className="px-5 py-3 flex items-center justify-between">
                <div>
                  <p className="text-sm text-fg">
                    {new Date(r.check_in_at).toLocaleDateString(locale === "ne" ? "ne-NP" : "en-US", {
                      weekday: "short", month: "short", day: "numeric",
                    })}
                  </p>
                  <p className="text-xs text-fg-muted">
                    {new Date(r.check_in_at).toLocaleTimeString(locale === "ne" ? "ne-NP" : "en-US", {
                      hour: "2-digit", minute: "2-digit",
                    })}
                  </p>
                </div>
                <MethodBadge method={r.method} />
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
```

**Step 2: Commit**

```bash
git add web/src/app/[lang]/member/attendance/page.tsx
git commit -m "feat: add member attendance page with streaks"
```

---

## Task 9: Member Packages Page

**Files:**
- Create: `web/src/app/[lang]/member/packages/page.tsx`

**Step 1: Create the packages page**

Create `web/src/app/[lang]/member/packages/page.tsx`:

```typescript
"use client";

import { useEffect, useState } from "react";
import { getDictionary } from "@/lib/i18n";
import { api } from "@/lib/api";
import type { Locale, MemberPackage } from "@/types";
import { Package, Loader2 } from "lucide-react";

function StatusBadge({ status }: { status: string }) {
  const colors: Record<string, string> = {
    active: "bg-accent/10 text-accent border-accent/20",
    expired: "bg-red-500/10 text-red-400 border-red-500/20",
    cancelled: "bg-gray-500/10 text-gray-400 border-gray-500/20",
    pending: "bg-amber-500/10 text-amber-400 border-amber-500/20",
  };
  return (
    <span className={`inline-flex items-center px-2 py-0.5 rounded-full text-[11px] font-mono uppercase border ${colors[status] ?? colors.pending}`}>
      {status}
    </span>
  );
}

export default function MemberPackagesPage({ params }: { params: { lang: string } }) {
  const locale = params.lang as Locale;
  const t = getDictionary(locale);

  const [packages, setPackages] = useState<MemberPackage[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    async function load() {
      try {
        const res = await api.getMyPackages({ limit: 50 });
        setPackages(res.data ?? []);
      } catch (err) {
        console.error("Failed to load packages:", err);
      } finally {
        setLoading(false);
      }
    }
    load();
  }, []);

  if (loading) {
    return (
      <div className="flex items-center justify-center py-20">
        <Loader2 className="w-8 h-8 text-accent animate-spin" />
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-semibold text-fg">{t.memberPages.packages}</h1>

      {packages.length === 0 ? (
        <div className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl px-5 py-12 text-center">
          <Package className="w-12 h-12 text-fg-muted mx-auto mb-3" />
          <p className="text-fg-muted">{t.memberPages.noPackages}</p>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          {packages.map((pkg) => {
            const daysLeft = Math.max(0, Math.ceil((new Date(pkg.end_date).getTime() - Date.now()) / 86400000));
            return (
              <div key={pkg.id} className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl p-5">
                <div className="flex items-start justify-between">
                  <div>
                    <h3 className="text-lg font-semibold text-fg">{pkg.package_name ?? "Package"}</h3>
                    <p className="text-xs text-fg-muted mt-1">{pkg.org_name}</p>
                  </div>
                  <StatusBadge status={pkg.status} />
                </div>
                <div className="mt-4 grid grid-cols-2 gap-3 text-sm">
                  <div>
                    <p className="text-fg-muted text-xs">Start</p>
                    <p className="text-fg">{new Date(pkg.start_date).toLocaleDateString()}</p>
                  </div>
                  <div>
                    <p className="text-fg-muted text-xs">End</p>
                    <p className="text-fg">{new Date(pkg.end_date).toLocaleDateString()}</p>
                  </div>
                  <div>
                    <p className="text-fg-muted text-xs">Paid</p>
                    <p className="text-fg">Rs. {Number(pkg.amount_paid).toLocaleString("en-IN")}</p>
                  </div>
                  <div>
                    <p className="text-fg-muted text-xs">Remaining</p>
                    <p className="text-fg">{pkg.status === "active" ? `${daysLeft} ${t.memberPages.daysRemaining}` : t.memberPages.expired}</p>
                  </div>
                </div>
              </div>
            );
          })}
        </div>
      )}
    </div>
  );
}
```

**Step 2: Commit**

```bash
git add web/src/app/[lang]/member/packages/page.tsx
git commit -m "feat: add member packages page"
```

---

## Task 10: Member Health Page

**Files:**
- Create: `web/src/app/[lang]/member/health/page.tsx`

**Step 1: Create the health page**

Create `web/src/app/[lang]/member/health/page.tsx`:

```typescript
"use client";

import { useEffect, useState } from "react";
import { getDictionary } from "@/lib/i18n";
import { api } from "@/lib/api";
import type { Locale, HealthMetric } from "@/types";
import { Heart, Ruler, Loader2 } from "lucide-react";

export default function MemberHealthPage({ params }: { params: { lang: string } }) {
  const locale = params.lang as Locale;
  const t = getDictionary(locale);

  const [metrics, setMetrics] = useState<HealthMetric[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    async function load() {
      try {
        const res = await api.getMyHealth({ limit: 50 });
        setMetrics(res.data ?? []);
      } catch (err) {
        console.error("Failed to load health data:", err);
      } finally {
        setLoading(false);
      }
    }
    load();
  }, []);

  if (loading) {
    return (
      <div className="flex items-center justify-center py-20">
        <Loader2 className="w-8 h-8 text-accent animate-spin" />
      </div>
    );
  }

  const bmiRecords = metrics.filter((m) => m.metric_type === "bmi");
  const measurementRecords = metrics.filter((m) => m.metric_type === "measurements");

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-semibold text-fg">{t.memberPages.health}</h1>

      {metrics.length === 0 ? (
        <div className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl px-5 py-12 text-center">
          <Heart className="w-12 h-12 text-fg-muted mx-auto mb-3" />
          <p className="text-fg-muted">{t.memberPages.noHealthData}</p>
        </div>
      ) : (
        <>
          {/* BMI section */}
          {bmiRecords.length > 0 && (
            <div className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl overflow-hidden">
              <div className="px-5 py-4 border-b border-white/[0.06] flex items-center gap-2">
                <Heart className="w-4 h-4 text-accent" />
                <h2 className="text-sm font-semibold text-fg">{t.memberPages.bmi}</h2>
              </div>
              <div className="divide-y divide-white/[0.04]">
                {bmiRecords.map((m) => {
                  const v = m.value as Record<string, number | string>;
                  return (
                    <div key={m.id} className="px-5 py-3 flex items-center justify-between">
                      <div>
                        <p className="text-sm text-fg">
                          BMI: <span className="font-semibold">{v.bmi}</span>
                          <span className="ml-2 text-xs text-fg-muted capitalize">({String(v.category)})</span>
                        </p>
                        <p className="text-xs text-fg-muted">
                          {v.height_cm}cm / {v.weight_kg}kg
                        </p>
                      </div>
                      <p className="text-xs text-fg-muted">
                        {new Date(m.recorded_at).toLocaleDateString()}
                      </p>
                    </div>
                  );
                })}
              </div>
            </div>
          )}

          {/* Measurements section */}
          {measurementRecords.length > 0 && (
            <div className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl overflow-hidden">
              <div className="px-5 py-4 border-b border-white/[0.06] flex items-center gap-2">
                <Ruler className="w-4 h-4 text-accent" />
                <h2 className="text-sm font-semibold text-fg">{t.memberPages.measurements}</h2>
              </div>
              <div className="divide-y divide-white/[0.04]">
                {measurementRecords.map((m) => {
                  const v = m.value as Record<string, number>;
                  return (
                    <div key={m.id} className="px-5 py-4">
                      <div className="flex items-center justify-between mb-2">
                        <p className="text-xs text-fg-muted">
                          {new Date(m.recorded_at).toLocaleDateString()}
                        </p>
                        {v.whr && (
                          <span className="text-xs text-fg-muted">WHR: {v.whr}</span>
                        )}
                      </div>
                      <div className="grid grid-cols-5 gap-2 text-center">
                        {[
                          ["Chest", v.chest_cm],
                          ["Waist", v.waist_cm],
                          ["Hips", v.hips_cm],
                          ["Arms", v.arms_cm],
                          ["Thighs", v.thighs_cm],
                        ].map(([label, val]) => (
                          <div key={String(label)}>
                            <p className="text-lg font-semibold text-fg">{val}</p>
                            <p className="text-[10px] text-fg-muted uppercase">{String(label)}</p>
                          </div>
                        ))}
                      </div>
                    </div>
                  );
                })}
              </div>
            </div>
          )}
        </>
      )}
    </div>
  );
}
```

**Step 2: Commit**

```bash
git add web/src/app/[lang]/member/health/page.tsx
git commit -m "feat: add member health metrics page"
```

---

## Task 11: Member Workouts Page

**Files:**
- Create: `web/src/app/[lang]/member/workouts/page.tsx`

**Step 1: Create the workouts page**

Create `web/src/app/[lang]/member/workouts/page.tsx`:

```typescript
"use client";

import { useEffect, useState } from "react";
import { getDictionary } from "@/lib/i18n";
import { useAuth } from "@/lib/auth";
import { api } from "@/lib/api";
import type { Locale, WorkoutLog, WorkoutPlan } from "@/types";
import { Dumbbell, Clock, Loader2 } from "lucide-react";

export default function MemberWorkoutsPage({ params }: { params: { lang: string } }) {
  const locale = params.lang as Locale;
  const t = getDictionary(locale);
  const { user } = useAuth();

  const [logs, setLogs] = useState<WorkoutLog[]>([]);
  const [plan, setPlan] = useState<WorkoutPlan | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    async function load() {
      try {
        const orgId = user?.org_id;
        const [logsRes, planRes] = await Promise.all([
          api.getMyWorkoutLogs({ organization_id: orgId, limit: 20 }),
          orgId ? api.getMyWorkoutPlan(orgId).catch(() => null) : Promise.resolve(null),
        ]);
        setLogs(logsRes.data ?? []);
        setPlan(planRes);
      } catch (err) {
        console.error("Failed to load workouts:", err);
      } finally {
        setLoading(false);
      }
    }
    load();
  }, [user]);

  if (loading) {
    return (
      <div className="flex items-center justify-center py-20">
        <Loader2 className="w-8 h-8 text-accent animate-spin" />
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-semibold text-fg">{t.memberPages.workouts}</h1>

      {/* Assigned plan */}
      <div className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl overflow-hidden">
        <div className="px-5 py-4 border-b border-white/[0.06] flex items-center gap-2">
          <Dumbbell className="w-4 h-4 text-accent" />
          <h2 className="text-sm font-semibold text-fg">{t.memberPages.workoutPlan}</h2>
        </div>
        {plan?.template ? (
          <div className="px-5 py-4">
            <h3 className="text-lg font-semibold text-fg">{plan.template.name}</h3>
            {plan.template.description && (
              <p className="text-sm text-fg-muted mt-1">{plan.template.description}</p>
            )}
            <div className="flex items-center gap-4 mt-2 text-xs text-fg-muted">
              {plan.template.category && <span className="capitalize">{plan.template.category}</span>}
              {plan.template.difficulty && <span className="capitalize">{plan.template.difficulty}</span>}
              {plan.template.duration_minutes && (
                <span className="flex items-center gap-1">
                  <Clock className="w-3 h-3" />
                  {plan.template.duration_minutes} {t.memberPages.minutes}
                </span>
              )}
            </div>
          </div>
        ) : (
          <p className="px-5 py-8 text-center text-fg-muted text-sm">{t.memberPages.noPlan}</p>
        )}
      </div>

      {/* Workout logs */}
      <div className="bg-gradient-to-b from-white/[0.08] to-white/[0.02] border border-white/[0.06] rounded-2xl overflow-hidden">
        <div className="px-5 py-4 border-b border-white/[0.06]">
          <h2 className="text-sm font-semibold text-fg">{t.memberPages.recentWorkouts}</h2>
        </div>
        {logs.length === 0 ? (
          <p className="px-5 py-8 text-center text-fg-muted text-sm">{t.memberPages.noWorkouts}</p>
        ) : (
          <div className="divide-y divide-white/[0.04]">
            {logs.map((log) => (
              <div key={log.id} className="px-5 py-3">
                <div className="flex items-center justify-between">
                  <p className="text-sm text-fg">
                    {new Date(log.logged_at).toLocaleDateString(locale === "ne" ? "ne-NP" : "en-US", {
                      weekday: "short", month: "short", day: "numeric",
                    })}
                  </p>
                  {log.duration_minutes && (
                    <span className="flex items-center gap-1 text-xs text-fg-muted">
                      <Clock className="w-3 h-3" />
                      {log.duration_minutes} {t.memberPages.minutes}
                    </span>
                  )}
                </div>
                {log.notes && <p className="text-xs text-fg-muted mt-1">{log.notes}</p>}
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
```

**Step 2: Commit**

```bash
git add web/src/app/[lang]/member/workouts/page.tsx
git commit -m "feat: add member workouts page"
```

---

## Task 12: E2E Playwright MCP Walkthrough

**Prerequisites:** Tasks 1-11 must be complete. API server running with `DEV_OTP_BYPASS=true`. Seed data loaded.

**Step 1: Start the API server**

```bash
# Kill any existing API process
pkill -f "go run ./cmd/api" || true
# Start with dev bypass
DB_PASSWORD=urja_secret DEV_OTP_BYPASS=true go run ./cmd/api &
```

**Step 2: Verify the web frontend is running**

```bash
curl -s http://localhost:3000 | head -5
# Should return HTML
```

**Step 3: Flow 1 — Gym Owner (admin)**

Using Playwright MCP browser tools:
1. Navigate to `http://localhost:3000/en/login`
2. Take screenshot of login page
3. Fill phone input with `9800000001`
4. Click "Send OTP" button
5. Fill OTP input with `123456`
6. Click "Verify" button
7. Wait for redirect to `/en/dashboard`
8. Take screenshot of dashboard
9. Click "Members" in sidebar
10. Take screenshot of members page
11. Click "Add Member" button
12. Take screenshot of add member modal
13. Close modal, click "Packages" in sidebar
14. Take screenshot of packages page
15. Click "Attendance" in sidebar
16. Take screenshot of attendance page
17. Click "Staff" in sidebar
18. Take screenshot of staff page
19. Logout via profile dropdown

**Step 4: Flow 2 — Staff**

1. Navigate to `http://localhost:3000/en/login`
2. Fill phone `9800000002`, submit OTP `123456`
3. Wait for dashboard redirect
4. Take screenshot of staff dashboard
5. Click Members, take screenshot
6. Click Attendance, take screenshot
7. Logout

**Step 5: Flow 3 — Member**

1. Navigate to `http://localhost:3000/en/login`
2. Fill phone `9800000003`, submit OTP `123456`
3. Wait for redirect to `/en/member`
4. Take screenshot of member dashboard
5. Click "My Attendance" in sidebar
6. Take screenshot of attendance page with streaks
7. Click "My Packages" in sidebar
8. Take screenshot of packages page
9. Click "My Health" in sidebar
10. Take screenshot of health metrics
11. Click "My Workouts" in sidebar
12. Take screenshot of workouts page
13. Logout

**Step 6: Commit screenshots**

```bash
git add web/screenshots/
git commit -m "feat: add E2E walkthrough screenshots for all 3 roles"
```
