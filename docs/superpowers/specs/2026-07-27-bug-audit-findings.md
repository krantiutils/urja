# Bug Audit Findings — 2026-07-27

Two independent audits of the repo. Ground truth established first:

- `go build ./...` — passes
- `go vet ./...` — fails on 3 test files only (stale signatures, listed below)
- `npx tsc --noEmit` — zero errors
- `npm run build` — passes, all 27 routes × 2 locales

Both dictionaries in `web/src/lib/i18n.ts` were diffed key-by-key (700 leaf keys)
and are fully symmetric — no missing translations. Every `href` and
`router.push` was diffed against the actual route directories — no broken links.

Criticals and highs are fixed in Phase 0 of
[the tenant sites design](2026-07-27-tenant-sites-design.md). Mediums and lows
are recorded here for later.

---

## Critical

| # | Finding | Location |
|---|---------|----------|
| 1 | Global JWT `role` claim used instead of per-org role across 10 route groups | `internal/auth/repository.go:105`, 10 `routes.go` files |
| 2 | `internal/subscription` member routes have no authorization middleware | `internal/subscription/routes.go:8` |
| 3 | Staff can self-promote to org admin | `internal/member/service.go:134` |
| 4 | `training_guides` has no `organization_id` — cross-org tampering and deletion | `db/migrations/000009`, `internal/guide/repository.go` |
| 5 | SMS credit purchase trusts client-supplied price | `internal/smsapi/service.go:30` |
| 6 | Khalti `pidx` never recorded — payment replayable across orgs | `internal/billing/service.go:43` |
| 7 | Accounts / due-payments send `amount` as a JSON string; Go decodes `float64` → every submit 400s | `dashboard/accounts/page.tsx:156,166`, `due-payments/page.tsx:109` |
| 8 | SMS send hardcodes `member_ids: []`; backend treats empty as "all members" | `dashboard/sms/page.tsx:93` |
| 9 | Account deletion calls a nonexistent endpoint and reads the wrong token field | `[lang]/delete/page.tsx:54,69` |
| 10 | "Get Recommendation" persists a plan reassignment before the user confirms | `internal/workout/service.go:270`, `member/workouts/page.tsx:456` |

## High

| # | Finding | Location |
|---|---------|----------|
| 11 | NFC org routes have no role gate | `internal/nfc/routes.go:8` |
| 12 | `absentee.Notify` bypasses SMS credit metering | `internal/absentee/service.go:37` |
| 13 | Attendance monthly calendar: UTC vs Asia/Kathmandu boundary mismatch | `internal/attendance/repository.go:246` |
| 14 | Nutrition dashboard defaults "today" to server UTC | `internal/nutrition/service.go:267` |
| 15 | Member search filters only the current page of 20 | `dashboard/members/page.tsx:115` |
| 16 | Attendance date filter runs over a fixed `limit: 100` fetch | `dashboard/attendance/page.tsx:96` |
| 17 | Health page's mixed-type `limit: 20` hides BMI/measurement records | `member/health/page.tsx:48` |
| 18 | Package delete has no confirmation dialog | `dashboard/packages/page.tsx:158` |
| 19 | Date-only fields parsed as UTC then rendered locally — off-by-one day | `member/packages/page.tsx`, `member/page.tsx`, `member/progress/page.tsx:92` |
| 20 | Transient profile-fetch failure bounces an onboarded user to onboarding | `lib/auth.tsx:58` |
| 21 | Member dashboard computes "today" in UTC; nutrition page uses local | `member/page.tsx:42` |

## Medium — deferred

| # | Finding | Location |
|---|---------|----------|
| 22 | `nfc.AssignCard` doesn't verify the target user belongs to the org | `internal/nfc/service.go:88` |
| 23 | `workout.AssignPlan` doesn't verify the target member belongs to the org | `internal/workout/service.go:151` |
| 24 | Privacy settings save fails completely silently | `dashboard/settings/page.tsx:239` |
| 25 | Unthrottled per-keystroke search causes stale-result races | `dashboard/staff/page.tsx:108`, `due-payments/page.tsx:53` |
| 26 | Org name and personal name fields allow a blank save | `dashboard/settings/page.tsx:342,576` |
| 27 | `orgName` lacks the `profile.organizations[0]` fallback that `orgId` has | `member/page.tsx:465` |
| 28 | Active package and streak chosen without filtering by org | `member/page.tsx:471`, `member/attendance/page.tsx:53` |
| 29 | Feedback submit silently no-ops when `orgId` is empty | `member/page.tsx:562` |
| 30 | Custom food creation allows negative macros | `member/nutrition/page.tsx:1328` |
| 31 | Nutrition goal form's required check uses string truthiness — `"0"` passes | `member/nutrition/page.tsx:1751` |
| 32 | Mutation handlers across nutrition/programs/workouts swallow errors with only `console.error` | multiple |
| 33 | "Enroll Now" silently discards an in-progress enrollment in another program | `member/programs/page.tsx:596` |
| 34 | Onboarding silently defaults gender to "male" for BMR math, disagreeing with the saved profile | `onboarding/page.tsx:215` |
| 35 | `<html lang="en">` hardcoded — `/ne` pages report English to screen readers and search engines | `app/layout.tsx:27` |

## Low — deferred

| # | Finding | Location |
|---|---------|----------|
| 36 | `nutrition.GetFoodItem` has no ownership filter (minor IDOR on custom foods) | `internal/nutrition/repository.go:290` |
| 37 | `activitylog` org route has no role gate (write path currently unwired) | `internal/activitylog/routes.go:8` |
| 38 | `decodeJWT` calls `atob()` on a base64url segment without char translation | `lib/auth.tsx:32` |

## Stale test signatures — breaks `go vet` and `go test`

| File | Problem |
|------|---------|
| `internal/auth/service_test.go:76` | `generateAccessToken` called with 3 args; signature takes 4 (`isSuperAdmin` added) |
| `internal/qrcode/handler_test.go:21` | `NewHandler` called with 2 args; signature takes 3 (`OrgSlugLookup` added) |
| `tests/e2e/setup_test.go:208` | `member.NewHandler` called with 2 args; signature takes 3 (`*leaderboard.Service` added) |

## Dead file

`migrations/001_initial_schema.sql` at the repo root is stale and unused —
`Makefile` and `Dockerfile:30` both point at `db/migrations`, the live set of 44
migrations. The two disagree on column names (`nfc_cards.hex_code` vs live
`card_hex`; `attendance.method` vs live `check_in_method`) and all Go code uses
the live names. Delete it to prevent confusion.
