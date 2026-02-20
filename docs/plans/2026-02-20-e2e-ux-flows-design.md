# E2E UX Flow Testing Design

## Overview

Interactive Playwright MCP-driven E2E walkthroughs of the Urja gym management platform across three user roles: gym owner (admin), staff, and member.

## Prerequisites

### 1. Dev OTP Bypass

Add `DEV_OTP_BYPASS=true` env var support. When set, phone numbers starting with `9800` accept OTP `123456` without sending SMS.

**Files changed:**
- `internal/config/config.go` - add `DevOTPBypass bool` field
- `internal/auth/service.go` - skip SMS for test phones in `RequestOTP()`, accept hardcoded OTP in `VerifyOTP()`

### 2. Test Data Seed

SQL seed script creating test accounts and sample data.

| Phone | Role | Purpose |
|-------|------|---------|
| `9800000001` | admin | Gym owner E2E flow |
| `9800000002` | staff | Staff E2E flow |
| `9800000003` | member | Member E2E flow |

Also seeds: 1 org ("Kathmandu Fitness Hub"), packages, attendance records, health metrics, workout data.

### 3. Member-Facing Pages

New pages under `/[lang]/member/` with dedicated member layout and sidebar:

| Page | Route | API Endpoint |
|------|-------|-------------|
| Dashboard | `/member` | `/members/me` + overview stats |
| My Attendance | `/member/attendance` | `/members/me/attendance` + `/members/me/streaks` |
| My Packages | `/member/packages` | `/members/me/packages` |
| My Health | `/member/health` | `/members/me/health` |
| My Workouts | `/member/workouts` | `/members/me/workout-logs` + `/members/me/workout-plan` |

Same dark cinematic design system as existing dashboard pages.

## E2E Walkthrough Flows

### Flow 1: Gym Owner (admin)

1. Login page -> enter `9800000001` -> OTP `123456`
2. Dashboard with stats overview
3. Members management (list, add member)
4. Packages management (list, create package)
5. Attendance view
6. Staff management

### Flow 2: Staff

1. Login as `9800000002`
2. Dashboard view
3. Members list + search
4. Manual check-in
5. Packages view

### Flow 3: Member

1. Login as `9800000003`
2. Member dashboard with profile
3. My Attendance with streaks
4. My Packages
5. My Health metrics
6. My Workouts

Screenshots captured at each step via Playwright MCP.

## Approach

- Interactive walkthrough using Playwright MCP (not persistent test scripts)
- Live browser navigation with screenshots at each step
- All three flows run sequentially in one session
