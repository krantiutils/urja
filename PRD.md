# Urja - Gym Management Platform

## Executive Summary

Secure gym management platform (FitFlow clone) with web and mobile apps, built with security-first architecture to avoid the catastrophic vulnerabilities found in the original FitFlow system.

## Tech Stack (Decided)

| Layer | Technology |
|-------|------------|
| **Web Frontend** | Next.js 14 + TypeScript |
| **Mobile** | Flutter |
| **Backend** | Go (performance-first) |
| **Database** | PostgreSQL |
| **Cache** | Redis |
| **File Storage** | S3-compatible (MinIO) |
| **Push Notifications** | Firebase Cloud Messaging |
| **Payments** | Khalti (Nepal, offline-capable) |
| **SMS/OTP** | Aakash SMS |
| **API Gateway** | Kong or Nginx |
| **i18n** | Bilingual — English + Nepali (ne) |

## Key Decisions

| Decision | Answer |
|----------|--------|
| Workout plans | Staff assigns templates + preset plan library |
| Multi-gym | M:N (member CAN join multiple gyms) |
| Staff roles | Flat staff + admin (no sub-roles for now) |
| QR check-in | Member scans gym's QR code |
| Revenue model | SaaS per gym (like FitFlow) |
| FitFlow migration | Deferred (fresh start for now) |
| Offline mode | Required (queue + sync) |
| Language | Bilingual en/ne with [lang] routing |

## Design System

See UI.md for the full Linear/Modern dark design system.

---

## Part 1: Original FitFlow Security Failures (What NOT To Do)

### Critical Vulnerabilities Found

| # | Vulnerability | CVSS | Root Cause |
|---|--------------|------|------------|
| 1 | OTP Code Leak | 9.8 | OTP returned in login API response |
| 2 | Token Without Verification | 10.0 | Bearer token issued before OTP verification |
| 3 | Password Hash Exposure | 7.5 | Bcrypt hashes in `/organizations` endpoint |
| 4 | No Rate Limiting | 5.3 | No throttling on any endpoint |
| 5 | NFC Card Database Exposure | 8.1 | All NFC hex codes accessible to any staff token |
| 6 | FCM Token Exposure | 9.1 | Firebase tokens in member API response |
| 7 | Staff Token Escalation | 9.8 | Staff tokens accessible via public gym phones |
| 8 | Cross-Org Data Access | 8.5 | No organization-level authorization |

### Scale of Exposure
- **120 gyms** in system
- **35,421 members** extracted
- **5,063 NFC cards** exposed
- **39,000+ FCM tokens** leaked
- **120 password hashes** exposed

---

## Part 2: System Architecture

### High-Level Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                         CLIENTS                                  │
├─────────────────┬─────────────────┬─────────────────────────────┤
│   Web App       │   iOS App       │   Android App               │
│   (Next.js)     │   (Flutter)     │   (Flutter)                 │
└────────┬────────┴────────┬────────┴─────────────┬───────────────┘
         │                 │                      │
         └─────────────────┼──────────────────────┘
                           │
                    ┌──────▼──────┐
                    │   API       │
                    │   Gateway   │
                    │   (Kong/    │
                    │   Nginx)    │
                    └──────┬──────┘
                           │
         ┌─────────────────┼─────────────────┐
         │                 │                 │
    ┌────▼────┐      ┌─────▼─────┐     ┌────▼────┐
    │ Auth    │      │ Core API  │     │ Media   │
    │ Service │      │ (Go)      │     │ Service │
    └────┬────┘      └─────┬─────┘     └────┬────┘
         │                 │                 │
    ┌────▼────┐      ┌─────▼─────┐     ┌────▼────┐
    │ Redis   │      │ PostgreSQL│     │ S3/     │
    │ (Tokens │      │ (Main DB) │     │ MinIO   │
    │  Cache) │      │           │     │ (Files) │
    └─────────┘      └───────────┘     └─────────┘
```

---

## Part 3: Security Architecture

### 3.1 Authentication Flow

- Phone-based OTP via Aakash SMS
- OTP: 6-digit, 5min expiry, max 3 attempts, 60s cooldown, 5/hour rate limit
- OTP stored as bcrypt hash in Redis, NEVER returned in API response
- JWT: 15min access token + 7-day refresh token
- Tokens only issued AFTER OTP verification

### 3.2 Authorization (RBAC)

| Role | Scope |
|------|-------|
| member | Own data only |
| staff | Organization-scoped |
| admin | Organization-scoped (full) |
| super_admin | Platform-wide (Urja team only) |

### 3.3 Data Protection

Sensitive fields NEVER exposed: password_hash, otp, fcm_token, nfc_hex_code, refresh_token, device_id

### 3.4 Rate Limiting

- Global: 100 req/min
- Auth login: 5/15min per IP
- OTP request: 5/hour per phone
- OTP verify: 3/5min
- Sensitive ops: 10/min

---

## Part 4: API Design

### Authentication
```
POST /api/v1/auth/login           - Request OTP
POST /api/v1/auth/verify-otp      - Verify OTP, get tokens
POST /api/v1/auth/refresh         - Refresh access token
POST /api/v1/auth/logout          - Invalidate tokens
POST /api/v1/auth/logout-all      - Logout all devices
```

### Member Profile
```
GET  /api/v1/members/me             - Own profile
PUT  /api/v1/members/me             - Update profile
GET  /api/v1/members/me/attendance  - Attendance history
GET  /api/v1/members/me/packages    - Packages
GET  /api/v1/members/me/streaks     - Check-in streaks
GET  /api/v1/members/me/timeline    - Activity timeline
GET  /api/v1/members/me/health      - BMI/WHR history
PUT  /api/v1/members/me/privacy     - Privacy settings
```

### Gamification & Engagement
```
GET  /api/v1/orgs/:orgId/leaderboard       - Leaderboard (weekly/monthly/all-time)
GET  /api/v1/members/me/achievements        - Badges & milestones
GET  /api/v1/members/me/streaks             - Current & best streaks
GET  /api/v1/members/me/timeline            - Activity feed (check-ins, PRs, milestones)
```

### Workouts & Training
```
GET    /api/v1/orgs/:orgId/workout-templates    - Preset workout plan library
POST   /api/v1/orgs/:orgId/workout-templates    - Staff creates template
GET    /api/v1/members/me/workout-plan          - My assigned plan
POST   /api/v1/orgs/:orgId/members/:id/assign-plan - Staff assigns plan to member
GET    /api/v1/members/me/workout-logs          - Workout history
POST   /api/v1/members/me/workout-logs          - Log a workout session
GET    /api/v1/training-guides                  - Training guide library
GET    /api/v1/training-guides/:id              - Training guide detail
```

### Health & Progress
```
GET  /api/v1/members/me/health              - All health metrics
POST /api/v1/members/me/health/bmi          - Log BMI
POST /api/v1/members/me/health/measurements - Log body measurements
POST /api/v1/members/me/health/photos       - Upload progress photo
GET  /api/v1/members/me/health/photos       - Progress photo timeline
```

### SaaS Billing (Gym Owners)
```
GET  /api/v1/billing/plans                  - Available SaaS plans
POST /api/v1/billing/subscribe              - Subscribe gym to plan
GET  /api/v1/billing/invoices               - Invoice history
POST /api/v1/billing/pay                    - Pay via Khalti
```

### Staff/Admin (Organization-Scoped)
```
GET    /api/v1/orgs/:orgId/members     - List members
GET    /api/v1/orgs/:orgId/members/:id - Member details
POST   /api/v1/orgs/:orgId/members     - Create member
PUT    /api/v1/orgs/:orgId/members/:id - Update member
DELETE /api/v1/orgs/:orgId/members/:id - Delete member
GET    /api/v1/orgs/:orgId/attendance  - Attendance records
GET    /api/v1/orgs/:orgId/nfc-cards   - NFC cards (org scoped)
GET    /api/v1/orgs/:orgId/analytics   - Dashboard analytics
```

### Public
```
GET  /api/v1/gyms       - List gyms
GET  /api/v1/gyms/:id   - Gym details
GET  /api/v1/packages   - Available packages
```

---

## Part 5: Database Schema

### Core Tables
- organizations, users, sessions, packages, member_packages
- attendance, nfc_cards, workouts, workout_logs
- fcm_tokens

### Multi-Gym Membership (M:N)
```sql
-- Junction table for M:N user↔org relationship
CREATE TABLE organization_members (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID REFERENCES users(id) ON DELETE CASCADE,
  organization_id UUID REFERENCES organizations(id) ON DELETE CASCADE,
  role VARCHAR(50) NOT NULL DEFAULT 'member',  -- member, staff, admin
  joined_at TIMESTAMP DEFAULT NOW(),
  status VARCHAR(50) DEFAULT 'active',
  UNIQUE(user_id, organization_id)
);
```

### Engagement Tables
```sql
-- Streaks
CREATE TABLE streaks (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  member_id UUID REFERENCES users(id) ON DELETE CASCADE,
  organization_id UUID REFERENCES organizations(id),
  current_streak INT DEFAULT 0,
  longest_streak INT DEFAULT 0,
  last_check_in DATE,
  updated_at TIMESTAMP DEFAULT NOW()
);

-- Achievements
CREATE TABLE achievements (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name VARCHAR(255) NOT NULL,
  name_ne VARCHAR(255),         -- Nepali translation
  description TEXT,
  description_ne TEXT,
  icon_url TEXT,
  criteria JSONB NOT NULL,      -- e.g. {"type": "streak", "value": 30}
  created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE member_achievements (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  member_id UUID REFERENCES users(id) ON DELETE CASCADE,
  achievement_id UUID REFERENCES achievements(id),
  earned_at TIMESTAMP DEFAULT NOW(),
  UNIQUE(member_id, achievement_id)
);

-- Health Metrics
CREATE TABLE health_metrics (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  member_id UUID REFERENCES users(id) ON DELETE CASCADE,
  metric_type VARCHAR(50) NOT NULL,  -- bmi, whr, weight, body_fat, measurements
  value JSONB NOT NULL,
  recorded_at TIMESTAMP DEFAULT NOW()
);

-- Progress Photos
CREATE TABLE progress_photos (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  member_id UUID REFERENCES users(id) ON DELETE CASCADE,
  photo_url TEXT NOT NULL,
  category VARCHAR(50),          -- front, side, back
  notes TEXT,
  taken_at TIMESTAMP DEFAULT NOW()
);

-- Workout Templates (preset + staff-created)
CREATE TABLE workout_templates (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  organization_id UUID REFERENCES organizations(id),  -- NULL = global preset
  name VARCHAR(255) NOT NULL,
  name_ne VARCHAR(255),
  description TEXT,
  category VARCHAR(100),
  difficulty VARCHAR(50),
  duration_minutes INT,
  exercises JSONB NOT NULL,       -- [{name, sets, reps, rest_seconds, notes}]
  is_preset BOOLEAN DEFAULT false,
  created_by UUID REFERENCES users(id),
  created_at TIMESTAMP DEFAULT NOW()
);

-- Training Guides
CREATE TABLE training_guides (
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

-- SaaS Billing
CREATE TABLE saas_plans (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name VARCHAR(255) NOT NULL,
  price_monthly DECIMAL(10,2),
  price_yearly DECIMAL(10,2),
  features JSONB DEFAULT '[]',
  max_members INT,
  is_active BOOLEAN DEFAULT true
);

CREATE TABLE saas_subscriptions (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  organization_id UUID REFERENCES organizations(id),
  plan_id UUID REFERENCES saas_plans(id),
  status VARCHAR(50) DEFAULT 'active',
  current_period_start TIMESTAMP,
  current_period_end TIMESTAMP,
  created_at TIMESTAMP DEFAULT NOW()
);
```

---

## Part 6: Offline Mode

### Requirements
- Check-ins must work offline (queue and sync)
- Attendance data cached locally
- Member card validation cached
- Sync when connectivity restored
- Conflict resolution: server wins for data, merge for attendance

### Implementation
- Service worker for web (PWA)
- Local SQLite for Flutter
- Background sync queue
- Optimistic UI updates

---

## Part 7: Aakash SMS Integration

### Configuration
- API: `https://sms.aakashsms.com/sms/v3/send`
- Auth: Token-based (`AAKASH_SMS_TOKEN` env var)
- Phone validation: Nepal format (10 digits, starts with 9/6/7/8)
- Existing token available from meds project

### OTP Flow
1. Client sends phone number
2. Server generates 6-digit OTP, stores bcrypt hash in Redis
3. Sends SMS via Aakash API
4. Client submits OTP for verification
5. Server verifies against Redis hash, issues JWT pair

---

## Part 8: Khalti Payment Integration

### Features
- Package purchases
- Subscription renewals
- Offline payment recording (cash, manual)
- Payment history and receipts

---

## Part 9: Frontend

### Web (Next.js) — See UI.md for design system
```
/[lang]/                              - Landing page (en/ne)
/[lang]/login                         - Phone + OTP login
/[lang]/register                      - Registration
/[lang]/dashboard                     - Member dashboard + activity timeline
/[lang]/dashboard/attendance          - Attendance history + streaks
/[lang]/dashboard/workouts            - Assigned workout plan + log
/[lang]/dashboard/packages            - Membership
/[lang]/dashboard/leaderboard         - Leaderboard (weekly/monthly/all-time)
/[lang]/dashboard/health              - BMI/WHR/measurements + progress photos
/[lang]/dashboard/achievements        - Badges & milestones
/[lang]/dashboard/profile             - Profile settings
/[lang]/gyms                          - Gym listing
/[lang]/gyms/[slug]                   - Gym details
/[lang]/guides                        - Training guide library
/[lang]/guides/[slug]                 - Training guide detail
/[lang]/admin                         - Admin panel
/[lang]/admin/members                 - Member management
/[lang]/admin/attendance              - Attendance management
/[lang]/admin/nfc                     - NFC card management
/[lang]/admin/packages                - Package management
/[lang]/admin/workouts                - Workout template management
/[lang]/admin/analytics               - Dashboard analytics
/[lang]/admin/billing                 - SaaS subscription & invoices
/[lang]/admin/settings                - Gym settings
```

### Mobile (Flutter)
```
Auth: Login, OTP, Profile Setup (bilingual)
Member: Home, Check-in (QR scan gym QR / NFC), Attendance + Streaks,
        Workouts + Training Guides, Profile, Notifications,
        Leaderboard, Health + Progress Photos, Achievements
Staff: Dashboard, Members, Check-in Manager, NFC Manager,
       Workout Template Assignment
Offline: Local SQLite cache, background sync queue
```

---

## Part 10: Implementation Phases

### Phase 1: Foundation (Week 1-2)
- Go project scaffolding (monorepo)
- PostgreSQL schema + migrations
- Auth service (OTP via Aakash SMS, JWT)
- Basic API structure with middleware

### Phase 2: Core Features (Week 3-4)
- Organization CRUD
- Member management
- Attendance system (manual + NFC + QR)
- Package management

### Phase 3: Web Application (Week 5-6)
- Landing page (Linear design system)
- Auth UI (OTP flow)
- Member dashboard + activity timeline
- Admin panel

### Phase 4: Mobile Application (Week 7-8)
- Flutter project setup
- Auth screens
- Member features + offline mode
- NFC integration

### Phase 5: Engagement & Advanced (Week 9-10)
- Streak system + leaderboard
- Achievement badges
- Health metrics (BMI/WHR/measurements) + progress photos
- Workout template library + staff assignment
- Training guides (bilingual content)
- Push notifications (FCM)
- Khalti payment integration
- SaaS billing for gym owners

### Phase 6: Testing & Launch (Week 11-12)
- Security testing (penetration test)
- Performance testing
- Offline mode testing
- i18n review (Nepali translations)
- Bug fixes
- Launch
