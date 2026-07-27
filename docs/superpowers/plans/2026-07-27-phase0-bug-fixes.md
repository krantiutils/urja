# Phase 0: Critical & High Bug Fixes — Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Close six critical authorization holes and repair four dead frontend features before any new work lands on top of them.

**Architecture:** The backend fixes are mostly a middleware swap — `RequireRole` reads a globally-derived JWT claim, `RequireOrgRole` reads the per-org role that `OrgScope` already resolved for the URL's `{orgId}`. Four fixes need real logic changes (training-guide org scoping, server-owned SMS pricing, Khalti replay prevention, SMS credit metering). The frontend fixes are independent and can run in parallel with the backend ones.

**Tech Stack:** Go 1.22 + chi v5 + pgx v5 + Redis; Next.js 14 App Router + TypeScript + Tailwind; golang-migrate for schema; Playwright for e2e.

## Global Constraints

- Live migrations live in `db/migrations/` (44 files, `NNNNNN_name.up.sql` / `.down.sql`). The root `migrations/001_initial_schema.sql` is dead — do not edit it.
- Every new migration needs a matching `.down.sql`.
- Timezone for all user-facing dates is `Asia/Kathmandu` (UTC+05:45). `internal/attendance/service.go:21` defines `nptLocation` — reuse it, do not redefine.
- All user-facing strings need both `en` and `ne` entries in `web/src/lib/i18n.ts`. The dictionaries are currently symmetric at 700 leaf keys; keep them that way.
- Go: no new dependencies. Frontend: no new dependencies.
- Verification gate for every task: `go build ./...`, `go vet ./...`, `go test ./...`, and for frontend tasks `npx tsc --noEmit` and `npm run build`.

---

### Task 1: Repair stale test signatures

Three test files call production functions with outdated arity, so `go vet ./...` and `go test ./...` both fail today. Nothing else in this plan can be verified until this is fixed. **This task must run first.**

**Files:**
- Modify: `internal/auth/service_test.go:76`
- Modify: `internal/qrcode/handler_test.go:21`
- Modify: `tests/e2e/setup_test.go:208`

**Interfaces:**
- Consumes: nothing
- Produces: a green `go test ./...` baseline that every later task depends on

- [ ] **Step 1: Confirm the failure**

```bash
cd /home/cdjk/github/urja && go vet ./... 2>&1 | head -20
```

Expected: three errors about "not enough arguments in call to".

- [ ] **Step 2: Read the current production signatures**

```bash
grep -n "func (s \*Service) generateAccessToken" internal/auth/service.go
grep -n "func NewHandler" internal/qrcode/handler.go internal/member/handler.go
```

- [ ] **Step 3: Update each call site to match**

`internal/auth/service_test.go:76` — `generateAccessToken` gained a fourth parameter `isSuperAdmin bool`. Pass `false` for the existing test cases; add one case passing `true` that asserts the claim round-trips.

`internal/qrcode/handler_test.go:21` — `NewHandler` gained an `OrgSlugLookup` third parameter. Pass a stub implementing the interface that returns a fixed slug.

`tests/e2e/setup_test.go:208` — `member.NewHandler` gained a `*leaderboard.Service` second parameter. Construct one from the existing test pool: `leaderboard.NewService(leaderboard.NewRepository(pool), logger)`.

- [ ] **Step 4: Verify**

```bash
go vet ./... && go build ./... && go test ./... 2>&1 | tail -30
```

Expected: vet and build clean. Record which tests pass; any pre-existing failure that needs a live database is acceptable and must be noted, not "fixed" by deletion.

- [ ] **Step 5: Commit**

```bash
git add internal/auth/service_test.go internal/qrcode/handler_test.go tests/e2e/setup_test.go
git commit -m "fix: update test call sites to current production signatures"
```

---

### Task 2: Replace global role checks with org-scoped role checks

The JWT's `role` claim comes from `SELECT role FROM organization_members WHERE user_id = $1 LIMIT 1` (`internal/auth/repository.go:105` and `:121`) — an arbitrary row, unrelated to the org being accessed. Ten route groups mounted under `/orgs/{orgId}` gate on it. `OrgScope` already puts the correct per-org role in the context; `RequireOrgRole` already reads it. This is a swap, plus a regression test proving the hole is closed.

**Files:**
- Modify: `internal/staff/routes.go:11`
- Modify: `internal/packages/routes.go:16`
- Modify: `internal/accounts/routes.go:11`
- Modify: `internal/absentee/routes.go:11`
- Modify: `internal/feedback/routes.go:14`
- Modify: `internal/workout/routes.go:15,26`
- Modify: `internal/guide/routes.go:18`
- Modify: `internal/smsapi/routes.go:11`
- Modify: `internal/notice/routes.go:16`
- Modify: `internal/dues/routes.go:11,18`
- Modify: `internal/auth/repository.go:105,121`
- Test: `tests/e2e/authz_test.go` (create)

**Interfaces:**
- Consumes: `middleware.RequireOrgRole(allowed ...string)` from `pkg/middleware/rbac.go:35` — reads `OrgRoleFromContext`, set by `OrgScope`
- Produces: nothing new; behavior change only

- [ ] **Step 1: Write the failing regression test**

Create `tests/e2e/authz_test.go`. The test builds two orgs and one user who is `admin` in Org A and `member` in Org B, then asserts that every admin route under Org B returns 403.

```go
func TestCrossOrgRoleEscalation(t *testing.T) {
	env := setupTestEnv(t)
	defer env.Teardown()

	orgA := env.CreateOrg(t, "Gym A")
	orgB := env.CreateOrg(t, "Gym B")
	user := env.CreateUser(t, "9800000001")
	env.AddMember(t, orgA, user, "admin")
	env.AddMember(t, orgB, user, "member")
	token := env.TokenFor(t, user)

	// Every one of these is an admin/staff action inside Org B, where the
	// user is only a member. All must be rejected.
	cases := []struct{ method, path string }{
		{"GET", "/api/v1/orgs/" + orgB + "/staff"},
		{"POST", "/api/v1/orgs/" + orgB + "/packages"},
		{"GET", "/api/v1/orgs/" + orgB + "/accounts"},
		{"GET", "/api/v1/orgs/" + orgB + "/absentees"},
		{"GET", "/api/v1/orgs/" + orgB + "/feedbacks"},
		{"POST", "/api/v1/orgs/" + orgB + "/workout-templates"},
		{"GET", "/api/v1/orgs/" + orgB + "/training-guides"},
		{"GET", "/api/v1/orgs/" + orgB + "/sms/balance"},
		{"POST", "/api/v1/orgs/" + orgB + "/notices"},
		{"GET", "/api/v1/orgs/" + orgB + "/dues"},
	}

	for _, tc := range cases {
		t.Run(tc.method+" "+tc.path, func(t *testing.T) {
			resp := env.Do(t, tc.method, tc.path, token, nil)
			if resp.StatusCode != http.StatusForbidden {
				t.Fatalf("expected 403, got %d — cross-org role escalation is open", resp.StatusCode)
			}
		})
	}
}
```

Match the helper names already used in `tests/e2e/setup_test.go`. If a helper does not exist, add it there rather than inventing a parallel harness.

- [ ] **Step 2: Run it to confirm it fails**

```bash
go test ./tests/e2e/ -run TestCrossOrgRoleEscalation -v 2>&1 | tail -30
```

Expected: most subtests fail with 200/201 instead of 403.

- [ ] **Step 3: Swap the middleware in all ten files**

In each listed file, change `middleware.RequireRole(` to `middleware.RequireOrgRole(` — arguments unchanged. Do **not** touch `RequireRole` call sites that are *not* under `/orgs/{orgId}` (verify each against `cmd/api/main.go:319-387`). Do not delete `RequireRole` itself; it stays for super-admin routes.

- [ ] **Step 4: Stop deriving a meaningless global role**

In `internal/auth/repository.go`, the `COALESCE((SELECT om.role ... LIMIT 1), 'member')` subquery in both `GetUserByID` (line 105) and `FindOrCreateUserByPhone` (line 121) produces a role with no org context. Replace both with the constant `'member'` so the JWT claim can never be mistaken for authority:

```go
// GetUserByID retrieves a user's phone and super admin status by their ID.
// The returned role is always "member": per-organization authority is resolved
// per-request by OrgScope, never from the token. See RequireOrgRole.
func (r *Repository) GetUserByID(ctx context.Context, userID string) (phone, role string, isSuperAdmin bool, err error) {
	err = r.db.QueryRow(ctx,
		`SELECT u.phone, u.is_super_admin FROM users u WHERE u.id = $1`,
		userID,
	).Scan(&phone, &isSuperAdmin)
	if err != nil {
		return "", "", false, fmt.Errorf("user not found: %w", err)
	}
	return phone, "member", isSuperAdmin, nil
}
```

Apply the equivalent change to `FindOrCreateUserByPhone`. Keep the signatures — callers in `internal/auth/service.go` are unchanged.

- [ ] **Step 5: Verify the test now passes**

```bash
go test ./tests/e2e/ -run TestCrossOrgRoleEscalation -v 2>&1 | tail -30
go build ./... && go vet ./... && go test ./... 2>&1 | tail -20
```

Expected: all subtests PASS.

- [ ] **Step 6: Check for collateral damage**

```bash
grep -rn "RequireRole(" internal/ cmd/ | grep -v "RequireOrgRole"
```

Every remaining hit must be a route that is genuinely not org-scoped. List them in the commit body with a one-line justification each.

- [ ] **Step 7: Commit**

```bash
git add internal/ tests/
git commit -m "fix: enforce per-org roles on org-scoped routes

The JWT role claim was derived from an arbitrary organization_members row
with no relation to the org being accessed, so an admin of one gym gained
admin powers in every other gym they belonged to."
```

---

### Task 3: Gate subscription routes

`internal/subscription/routes.go:8-23` applies no role middleware at all. Both groups mount inside `/orgs/{orgId}` (`cmd/api/main.go:327,335`), so any active member can assign themselves a paid subscription with a client-chosen `amount_paid`.

**Files:**
- Modify: `internal/subscription/routes.go`
- Test: `tests/e2e/authz_test.go` (extend Task 2's test)

**Interfaces:**
- Consumes: `middleware.RequireOrgRole` (Task 2)
- Produces: nothing new

- [ ] **Step 1: Add the failing cases to the authz test**

```go
{"POST", "/api/v1/orgs/" + orgB + "/members/" + memberID + "/packages/assign"},
{"POST", "/api/v1/orgs/" + orgB + "/members/" + memberID + "/packages/renew"},
{"GET", "/api/v1/orgs/" + orgB + "/members/" + memberID + "/payments"},
```

Add a second test asserting a plain `member` of Org A cannot self-assign inside Org A either — the member role is never sufficient for these routes.

- [ ] **Step 2: Run to confirm failure**

```bash
go test ./tests/e2e/ -run TestCrossOrgRoleEscalation -v 2>&1 | tail -20
```

- [ ] **Step 3: Add the middleware**

```go
func (h *Handler) RegisterPackageRoutes(r chi.Router) {
	r.Use(middleware.RequireOrgRole("staff", "admin"))
	// ... existing routes unchanged
}

func (h *Handler) RegisterMemberRoutes(r chi.Router) {
	r.Use(middleware.RequireOrgRole("staff", "admin"))
	// ... existing routes unchanged
}
```

Add the `middleware` import.

- [ ] **Step 4: Verify and commit**

```bash
go test ./tests/e2e/ -run TestCrossOrgRoleEscalation -v && go build ./...
git add internal/subscription/routes.go tests/e2e/authz_test.go
git commit -m "fix: require staff or admin role on subscription routes"
```

---

### Task 4: Block privilege self-escalation on member update

`member.UpdateOrgMember` (`internal/member/service.go:134`) accepts any `role` value, and the route is gated at `staff`. A staff user can `PUT /orgs/{orgId}/members/{their-own-id}` with `{"role":"admin"}`.

**Files:**
- Modify: `internal/member/service.go:134-143,204-219`
- Modify: `internal/member/handler.go` (thread the caller's ID and org role through)
- Test: `internal/member/service_test.go` (create if absent)

**Interfaces:**
- Consumes: `middleware.UserIDFromContext`, `middleware.OrgRoleFromContext` from `pkg/middleware`
- Produces: `UpdateOrgMember(ctx, callerID, callerRole, orgID, memberID string, in *UpdateOrgMemberInput) (*OrgMember, error)` — note the two new leading parameters

- [ ] **Step 1: Write the failing tests**

```go
func TestUpdateOrgMember_RejectsSelfRoleChange(t *testing.T) {
	// caller is staff, target is caller
	_, err := svc.UpdateOrgMember(ctx, "user-1", "staff", "org-1", "user-1",
		&UpdateOrgMemberInput{Role: ptr("admin")})
	if err == nil || !strings.Contains(err.Error(), "forbidden") {
		t.Fatalf("staff must not change their own role, got err=%v", err)
	}
}

func TestUpdateOrgMember_StaffCannotChangeAnyRole(t *testing.T) {
	_, err := svc.UpdateOrgMember(ctx, "user-1", "staff", "org-1", "user-2",
		&UpdateOrgMemberInput{Role: ptr("admin")})
	if err == nil || !strings.Contains(err.Error(), "forbidden") {
		t.Fatalf("only admins may change roles, got err=%v", err)
	}
}

func TestUpdateOrgMember_StaffCanEditNonRoleFields(t *testing.T) {
	_, err := svc.UpdateOrgMember(ctx, "user-1", "staff", "org-1", "user-2",
		&UpdateOrgMemberInput{Status: ptr("inactive")})
	if err != nil {
		t.Fatalf("staff may still edit non-role fields: %v", err)
	}
}
```

- [ ] **Step 2: Run to confirm failure**

```bash
go test ./internal/member/ -run TestUpdateOrgMember -v
```

Expected: compile error (signature mismatch), which is the first failure to fix.

- [ ] **Step 3: Implement the guard**

At the top of `UpdateOrgMember`, before any repository call:

```go
if in.Role != nil {
	if callerRole != "admin" {
		return nil, fmt.Errorf("forbidden: only admins can change a member's role")
	}
	if callerID == memberID {
		return nil, fmt.Errorf("forbidden: cannot change your own role")
	}
}
```

Update `handler.go` to read the caller ID and org role from context and pass them.

- [ ] **Step 4: Verify and commit**

```bash
go test ./internal/member/ -v && go build ./...
git add internal/member/
git commit -m "fix: only admins may change member roles, never their own"
```

---

### Task 5: Scope training guides to an organization

`training_guides` has no `organization_id` column, yet `guide.RegisterOrgRoutes` mounts under `/orgs/{orgId}`. `GetByID`, `Update`, `SetPublished` and `Delete` key on the guide ID alone, so any gym's staff can edit or delete the platform-wide library.

**Files:**
- Create: `db/migrations/000045_training_guides_org.up.sql`
- Create: `db/migrations/000045_training_guides_org.down.sql`
- Modify: `internal/guide/repository.go:57,214-233,235-253,255-260`
- Modify: `internal/guide/service.go`, `internal/guide/handler.go` (thread `orgID`)
- Test: `tests/e2e/authz_test.go`

**Interfaces:**
- Consumes: `middleware.OrgIDFromContext`
- Produces: repository methods gain a leading `orgID string` parameter: `GetByID(ctx, orgID, id)`, `Update(ctx, orgID, id, in)`, `SetPublished(ctx, orgID, id, published)`, `Delete(ctx, orgID, id)`

- [ ] **Step 1: Write the migration**

`000045_training_guides_org.up.sql`:

```sql
BEGIN;

ALTER TABLE training_guides
    ADD COLUMN IF NOT EXISTS organization_id UUID REFERENCES organizations(id) ON DELETE CASCADE;

COMMENT ON COLUMN training_guides.organization_id IS
    'Owning organization. NULL means a platform-wide preset, editable only by super admins.';

CREATE INDEX IF NOT EXISTS idx_training_guides_org ON training_guides(organization_id);

COMMIT;
```

`000045_training_guides_org.down.sql`:

```sql
BEGIN;
DROP INDEX IF EXISTS idx_training_guides_org;
ALTER TABLE training_guides DROP COLUMN IF EXISTS organization_id;
COMMIT;
```

Existing rows keep `NULL`, which correctly marks them as the platform presets they are today.

- [ ] **Step 2: Write the failing test**

Extend `tests/e2e/authz_test.go`: staff of Org A create a guide, then staff of Org B attempt `PUT`, `PATCH .../publish` and `DELETE` on it. All three must 404 (not 403 — do not leak the guide's existence).

- [ ] **Step 3: Run to confirm failure**

```bash
make migrate-up && go test ./tests/e2e/ -run TestGuideOrgScope -v
```

- [ ] **Step 4: Add the org filter to every mutating query**

`Create` sets `organization_id` from context. `ListAll` filters `WHERE organization_id = $1 OR organization_id IS NULL`. `GetByID`, `Update`, `SetPublished` and `Delete` all add `AND organization_id = $N` so a NULL-org preset can never be mutated through an org route. `ListPublished`/`GetPublished` (the public routes) stay unfiltered — published guides are public by design.

- [ ] **Step 5: Verify and commit**

```bash
go test ./tests/e2e/ -run TestGuideOrgScope -v && go build ./... && go test ./...
git add db/migrations/000045_* internal/guide/ tests/
git commit -m "fix: scope training guide mutations to the owning organization"
```

---

### Task 6: Server-owned SMS pricing

`POST /orgs/{orgId}/sms/buy` (`internal/smsapi/service.go:30`) accepts client-supplied `quantity`, `rate` and `amount`, then writes `payment_status='completed'` and increments the balance — no price table, no payment verification.

**Files:**
- Modify: `internal/config/config.go` (add `SMS.CreditRateNPR`, default `1.00`)
- Modify: `internal/smsapi/service.go:30-57`
- Modify: `internal/smsapi/handler.go:54-81`
- Modify: `.env.example`
- Test: `internal/smsapi/service_test.go` (create)

**Interfaces:**
- Consumes: `khaltiClient.Lookup(pidx string)` — the same call `internal/billing/service.go:43` already uses
- Produces: `BuyCredits(ctx, orgID string, quantity int, paymentMethod, pidx string) (*CreditPurchase, error)` — `rate` and `amount` no longer accepted from the caller

- [ ] **Step 1: Write the failing tests**

```go
func TestBuyCredits_IgnoresClientRate(t *testing.T) {
	// The handler must not read rate/amount from the body at all.
	// Assert the persisted purchase uses the configured rate.
}

func TestBuyCredits_KhaltiAmountMustMatch(t *testing.T) {
	// Lookup returns 500 paisa; quantity 100 at rate 1.00 NPR = 10000 paisa.
	// Expect an error and no credit grant.
}

func TestBuyCredits_RejectsUnverifiedKhalti(t *testing.T) {
	// Lookup returns status "Pending" -> error, no credits.
}
```

- [ ] **Step 2: Run to confirm failure**

```bash
go test ./internal/smsapi/ -v
```

- [ ] **Step 3: Implement**

Compute `amount := float64(quantity) * s.rate` server-side. For `payment_method == "khalti"`, call `Lookup(pidx)`, require `Status == "Completed"` and `Amount == int(amount*100)`, and only then grant credits. For a cash/manual method, require the caller's org role to be `admin` (already enforced by Task 2) and record it as a manual grant. Reject `quantity <= 0` and any quantity above a sane per-request ceiling.

- [ ] **Step 4: Verify and commit**

```bash
go test ./internal/smsapi/ -v && go build ./...
git add internal/smsapi/ internal/config/ .env.example
git commit -m "fix: compute SMS credit price server-side and verify Khalti payment"
```

---

### Task 7: Prevent Khalti payment replay

`billing.Subscribe` (`internal/billing/service.go:43-131`) verifies the payment correctly but never records the consumed `pidx`. `Lookup` keeps returning the same `Completed` payment, so one purchase can be replayed for many orgs.

**Files:**
- Create: `db/migrations/000046_khalti_payments.up.sql` / `.down.sql`
- Modify: `internal/billing/repository.go`, `internal/billing/service.go:43-131`
- Test: `internal/billing/service_test.go`

**Interfaces:**
- Produces: `Repository.ConsumePidx(ctx context.Context, tx pgx.Tx, pidx, orgID string) error` — returns a sentinel `ErrPidxAlreadyUsed` on unique violation

- [ ] **Step 1: Write the migration**

```sql
BEGIN;

CREATE TABLE IF NOT EXISTS khalti_payments (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    pidx            VARCHAR(120) NOT NULL UNIQUE,
    organization_id UUID REFERENCES organizations(id) ON DELETE SET NULL,
    amount_paisa    BIGINT NOT NULL,
    purpose         VARCHAR(40) NOT NULL,
    consumed_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_khalti_payments_org ON khalti_payments(organization_id);

COMMIT;
```

The UNIQUE constraint on `pidx` is the actual defense — the insert must happen in the same transaction as the subscription write.

- [ ] **Step 2: Write the failing test**

```go
func TestSubscribe_RejectsReplayedPidx(t *testing.T) {
	// Same pidx, two orgs. First call succeeds; second must fail with
	// ErrPidxAlreadyUsed and create no subscription.
}
```

- [ ] **Step 3: Run to confirm failure**

- [ ] **Step 4: Implement** — begin a transaction, `ConsumePidx` first (so the unique violation aborts before any grant), then write the subscription, then commit. Apply the same treatment to the SMS purchase path from Task 6 with `purpose='sms_credits'`.

- [ ] **Step 5: Verify and commit**

```bash
go test ./internal/billing/ -v && go build ./...
git add db/migrations/000046_* internal/billing/ internal/smsapi/
git commit -m "fix: record consumed Khalti pidx to prevent cross-org payment replay"
```

---

### Task 8: Gate NFC routes and verify card assignment targets

`internal/nfc/routes.go:8-19` applies no middleware, so any member can register access-control readers and assign cards. Separately, `nfc.AssignCard` (`internal/nfc/service.go:88`) never checks that the target user belongs to the org.

**Files:**
- Modify: `internal/nfc/routes.go`
- Modify: `internal/nfc/service.go:88-100`
- Modify: `internal/nfc/repository.go` (add an org-membership check)
- Test: `tests/e2e/authz_test.go`

**Interfaces:**
- Produces: `Repository.IsOrgMember(ctx context.Context, userID, orgID string) (bool, error)`

- [ ] **Step 1: Write the failing tests** — a plain member of Org B gets 403 on `POST /orgs/{orgB}/nfc-devices` and on card create/assign; and assigning a card to a user outside the org returns an error.

- [ ] **Step 2: Run to confirm failure**

- [ ] **Step 3: Implement** — `r.Use(middleware.RequireOrgRole("staff", "admin"))` in `RegisterOrgRoutes` (leave `RegisterDeviceRoutes` alone: it authenticates via `X-Device-Key`, not JWT). In `AssignCard`, call `IsOrgMember` on the target and return a "user is not a member of this organization" error when false.

- [ ] **Step 4: Verify and commit**

```bash
go test ./tests/e2e/ -run TestNFC -v && go build ./...
git add internal/nfc/ tests/
git commit -m "fix: require staff role for NFC management and verify card target membership"
```

---

### Task 9: Meter absentee SMS through the credit ledger

`absentee.Notify` (`internal/absentee/service.go:37-58`) calls the raw Aakash client directly, bypassing the credit deduction that `smsapi.SendSMS` performs. There is no rate limit or dedup, and it does not verify the member is actually absent.

**Files:**
- Modify: `internal/absentee/service.go:37-58`
- Modify: `cmd/api/main.go:151-154` (inject the SMS credit service)
- Test: `internal/absentee/service_test.go`

**Interfaces:**
- Consumes: a narrow interface rather than the concrete type, so absentee does not import smsapi's whole surface:

```go
// CreditedSender sends an SMS after deducting one credit from the org's balance.
type CreditedSender interface {
	SendMetered(ctx context.Context, orgID, phone, message string) error
}
```

`smsapi.Service` implements it.

- [ ] **Step 1: Write the failing tests** — notifying with a zero credit balance returns an error and sends nothing; a successful notify deducts exactly one credit; notifying the same member twice within the cooldown sends once.

- [ ] **Step 2: Run to confirm failure**

- [ ] **Step 3: Implement** — swap `smsClient.Send` for `creditedSender.SendMetered`, and add a per-member cooldown (24h) checked against the last notification timestamp.

- [ ] **Step 4: Verify and commit**

```bash
go test ./internal/absentee/ -v && go build ./...
git add internal/absentee/ internal/smsapi/ cmd/api/main.go
git commit -m "fix: deduct SMS credits for absentee notifications and add a cooldown"
```

---

### Task 10: Fix Nepal timezone date boundaries

Two queries mix UTC and `Asia/Kathmandu`. Nepal is UTC+05:45 and the Postgres container runs UTC, so the 00:00–05:45 window lands on the wrong day.

**Files:**
- Modify: `internal/attendance/repository.go:246-254`
- Modify: `internal/nutrition/service.go:267-270`
- Test: `internal/attendance/repository_test.go`, `internal/nutrition/service_test.go`

**Interfaces:**
- Consumes: `nptLocation` from `internal/attendance/service.go:21`. If nutrition cannot import attendance without a cycle, move `nptLocation` to `pkg/timeutil` and update both.

- [ ] **Step 1: Write the failing tests**

```go
func TestGetMonthlyCalendar_NepalMidnightBoundary(t *testing.T) {
	// A check-in at 2026-03-01T00:30:00+05:45 (= 2026-02-28T18:45:00Z)
	// must appear as day 1 of March, not in February's calendar.
}

func TestGetDailyDashboard_DefaultsToNepalToday(t *testing.T) {
	// With the clock at 2026-03-01T02:00:00+05:45 (2026-02-28T20:15:00Z),
	// the default date must be "2026-03-01", not "2026-02-28".
}
```

The nutrition test needs an injectable clock — add a `now func() time.Time` field to the service defaulting to `time.Now`, rather than reaching for a global.

- [ ] **Step 2: Run to confirm failure**

- [ ] **Step 3: Implement**

`GetMonthlyCalendar` currently bounds with `check_in_at >= ($2||'-01')::date`. Apply the same zone used for extraction:

```sql
WHERE check_in_at >= (($2 || '-01')::date AT TIME ZONE 'Asia/Kathmandu')
  AND check_in_at <  ((($2 || '-01')::date + INTERVAL '1 month') AT TIME ZONE 'Asia/Kathmandu')
```

`GetDailyDashboard` uses `s.now().In(nptLocation).Format("2006-01-02")`.

- [ ] **Step 4: Verify and commit**

```bash
go test ./internal/attendance/ ./internal/nutrition/ -v && go build ./...
git add internal/attendance/ internal/nutrition/ pkg/
git commit -m "fix: compute attendance and nutrition day boundaries in Nepal time"
```

---

### Task 11: Send money amounts as numbers

`amount` is bound to a text input and passed through unconverted, so the body is `"amount":"1500.50"`. The Go handlers decode into `float64` with no `,string` tag, so **every** Add Transaction, Edit Transaction and Record Payment request returns 400. These features have never worked.

**Files:**
- Modify: `web/src/app/[lang]/dashboard/accounts/page.tsx:156,166`
- Modify: `web/src/app/[lang]/dashboard/due-payments/page.tsx:109`
- Modify: `web/src/lib/api.ts` (change the `amount` field type from `string` to `number`)
- Test: `web/tests/accounts.spec.ts` (create)

**Interfaces:**
- Produces: `api.createTransaction`/`updateTransaction`/`recordPayment` take `amount: number`

- [ ] **Step 1: Write the failing Playwright test**

```ts
test("adding a transaction persists and appears in the list", async ({ page }) => {
  await loginAsAdmin(page);
  await page.goto("/en/dashboard/accounts");
  await page.getByRole("button", { name: "Add Transaction" }).click();
  await page.getByLabel("Amount").fill("1500.50");
  await page.getByLabel("Description").fill("Playwright test entry");
  await page.getByRole("button", { name: "Save" }).click();
  await expect(page.getByText("Playwright test entry")).toBeVisible();
});
```

- [ ] **Step 2: Run to confirm failure**

```bash
cd web && npx playwright test tests/accounts.spec.ts
```

Expected: fail — the row never appears because the API returns 400.

- [ ] **Step 3: Implement** — change the `amount` type in `api.ts` to `number`, and at each of the three call sites send `Number(formAmount)`. Guard first: reject `Number.isNaN` or `<= 0` with an inline validation message using existing i18n keys, adding new `en`/`ne` keys if none fits.

- [ ] **Step 4: Verify and commit**

```bash
npx tsc --noEmit && npm run build && npx playwright test tests/accounts.spec.ts
git add web/src web/tests
git commit -m "fix: send transaction and payment amounts as numbers, not strings"
```

---

### Task 12: Stop SMS from broadcasting to the whole gym

`dashboard/sms/page.tsx:93` hardcodes `member_ids: []`, and `smsapi.GetOrgMemberPhones` drops its `WHERE` clause when the list is empty — so every click texts every active member and burns credits.

**Files:**
- Modify: `web/src/app/[lang]/dashboard/sms/page.tsx`
- Modify: `internal/smsapi/repository.go:208-222` (defensive fix)
- Modify: `web/src/lib/i18n.ts`
- Test: `web/tests/sms.spec.ts` (create)

**Interfaces:**
- Produces: nothing new

- [ ] **Step 1: Write the failing test** — with no recipients selected the Send button is disabled; selecting two members and sending shows a confirmation naming the count; cancelling sends nothing.

- [ ] **Step 2: Run to confirm failure**

- [ ] **Step 3: Implement** — add a searchable member multi-select (reuse the member list API), a "Select all" that sets the IDs *explicitly* rather than relying on the empty-list default, a disabled Send until at least one recipient is chosen, and a confirmation dialog stating recipient count and credit cost. Server-side, make `GetOrgMemberPhones` return an empty slice for an empty ID list instead of every member — an empty list must never mean "everyone".

- [ ] **Step 4: Verify and commit**

```bash
npx tsc --noEmit && npm run build && npx playwright test tests/sms.spec.ts && cd .. && go test ./internal/smsapi/
git add web/src web/tests internal/smsapi
git commit -m "fix: require explicit SMS recipients instead of silently texting everyone"
```

---

### Task 13: Repair account deletion

`[lang]/delete/page.tsx` posts to `/api/v1/auth/verify` (the route is `/auth/verify-otp`, `internal/auth/routes.go:20`) and destructures `token` from a response shaped `{access_token, refresh_token, ...}`. Every attempt 404s, then would send `Bearer undefined`. This page exists for Google Play compliance, so it has to work.

**Files:**
- Modify: `web/src/app/[lang]/delete/page.tsx:54,69`
- Test: `web/tests/delete-account.spec.ts` (create)

**Interfaces:**
- Consumes: `POST /api/v1/auth/verify-otp` → `{access_token, refresh_token, user}` per `web/src/types/index.ts:3-9`

- [ ] **Step 1: Write the failing test** — request OTP, submit the dev-bypass OTP, confirm deletion, assert the success state renders and a subsequent login with the same phone starts a fresh account.

- [ ] **Step 2: Run to confirm failure**

- [ ] **Step 3: Implement** — correct the endpoint to `/api/v1/auth/verify-otp` and destructure `access_token`. Route the call through `web/src/lib/api.ts` rather than a bare `fetch`, so the base URL and error handling match everything else.

- [ ] **Step 4: Verify and commit**

```bash
npx tsc --noEmit && npm run build && npx playwright test tests/delete-account.spec.ts
git add web/src web/tests
git commit -m "fix: account deletion called a nonexistent endpoint and read the wrong token field"
```

---

### Task 14: Make workout recommendation a preview

`workout.RecommendPlan` (`internal/workout/service.go:270-277`) calls `SelfAssignPlan` and persists immediately, while the UI shows a preview with a "Choose Plan" button that makes no API call. Merely opening "Find My Plan" replaces the member's real plan.

**Files:**
- Modify: `internal/workout/service.go:270-277`
- Modify: `web/src/app/[lang]/member/workouts/page.tsx:456-472,474`
- Test: `internal/workout/service_test.go`, `web/tests/workouts.spec.ts`

**Interfaces:**
- Produces: `RecommendPlan(ctx, userID string, in RecommendInput) (*WorkoutTemplate, error)` — returns the suggestion and writes nothing. Assignment stays with the existing `SelfAssignPlan`.

- [ ] **Step 1: Write the failing tests** — calling `RecommendPlan` leaves the existing assignment untouched; `SelfAssignPlan` still replaces it.

- [ ] **Step 2: Run to confirm failure**

- [ ] **Step 3: Implement** — drop the `SelfAssignPlan` call from `RecommendPlan`. In the UI, wire `handleAssignRecommended` to `api.selfAssignPlan(templateId)` and only then update local state and refetch.

- [ ] **Step 4: Verify and commit**

```bash
go test ./internal/workout/ -v && cd web && npx tsc --noEmit && npm run build
git add internal/workout web/src web/tests
git commit -m "fix: workout recommendation previews instead of silently reassigning the plan"
```

---

### Task 15: Server-side member search and attendance date filtering

Two dashboard pages filter a truncated client-side page. Member search over a 20-row page reports "No members found" for anyone on page 2; the attendance page fetches the 100 most recent org-wide records and filters by date client-side, so any older date shows nothing.

**Files:**
- Modify: `web/src/app/[lang]/dashboard/members/page.tsx:115-121,342-347`
- Modify: `web/src/app/[lang]/dashboard/attendance/page.tsx:96-146`
- Modify: `web/src/lib/api.ts` (add `search` and `date` params)
- Modify: `internal/member/repository.go`, `internal/attendance/repository.go` (accept the params if not already supported)
- Test: `web/tests/members.spec.ts`, `web/tests/attendance.spec.ts` (extend the existing specs)

**Interfaces:**
- Produces: `api.listMembers(orgId, {limit, offset, search?})`, `api.listAttendance(orgId, {limit, offset, date?})`

- [ ] **Step 1: Write the failing tests** — seed 45 members, search from page 2 for a name that only exists on page 1, assert it is found; seed attendance older than the 100 most recent, pick that date, assert the rows render.

- [ ] **Step 2: Run to confirm failure**

- [ ] **Step 3: Implement** — check whether the Go repositories already accept `search`/`date` (the dues repository does, and is the reference pattern); add them where missing with proper parameterization — never string-concatenate into SQL. On the client, debounce input by 300ms, reset `offset` to 0 whenever the query changes, and guard against out-of-order responses with a request counter.

- [ ] **Step 4: Verify and commit**

```bash
go test ./internal/member/ ./internal/attendance/ && cd web && npx playwright test tests/members.spec.ts tests/attendance.spec.ts
git add internal/ web/
git commit -m "fix: filter members and attendance server-side instead of over a truncated page"
```

---

### Task 16: Frontend correctness cleanup

Five independent fixes, small enough to share a task but each with its own test.

**Files:**
- Modify: `web/src/app/[lang]/member/health/page.tsx:48`
- Modify: `web/src/app/[lang]/dashboard/packages/page.tsx:158-168`
- Modify: `web/src/app/[lang]/member/packages/page.tsx:71,104,110`
- Modify: `web/src/app/[lang]/member/progress/page.tsx:92,101`
- Modify: `web/src/app/[lang]/member/page.tsx:42,478,487-488`
- Modify: `web/src/lib/auth.tsx:58-76`
- Create: `web/src/lib/date.ts`
- Test: `web/tests/date.spec.ts`

**Interfaces:**
- Produces: `web/src/lib/date.ts` exporting `parseDateOnly(s: string): Date` (parses `YYYY-MM-DD` from components, never through `new Date(string)`) and `nepalToday(): string` (returns today's `YYYY-MM-DD` in `Asia/Kathmandu` via `Intl.DateTimeFormat`)

- [ ] **Step 1: Write the failing unit tests for `date.ts`**

```ts
test("parseDateOnly does not shift across timezones", () => {
  const d = parseDateOnly("2026-03-01");
  expect(d.getFullYear()).toBe(2026);
  expect(d.getMonth()).toBe(2);
  expect(d.getDate()).toBe(1);
});

test("nepalToday returns the Nepal date, not the UTC date", () => {
  // 2026-02-28T20:15:00Z is 2026-03-01T02:00:00+05:45
  expect(nepalToday(new Date("2026-02-28T20:15:00Z"))).toBe("2026-03-01");
});
```

- [ ] **Step 2: Run to confirm failure**

- [ ] **Step 3: Implement each fix**

1. **Health page** — `getMyHealth({limit: 20})` mixes all metric types before applying the limit, so a member who logs weight daily loses their single BMI entry. Fetch per type: separate calls with `type: "weight"`, `"bmi"`, `"body_measurements"`.
2. **Package delete** — add `if (!confirm(t.packages.confirmDelete)) return;`, matching `staff`/`stories`/`feedbacks`. Add the i18n key to both dictionaries.
3. **Date-only parsing** — replace every `new Date(dateString)` on a `YYYY-MM-DD` value with `parseDateOnly`.
4. **Member dashboard "today"** — replace the module-scope `new Date().toISOString().split("T")[0]` with `nepalToday()` computed inside the effect.
5. **Org-filtered picks** — `member/page.tsx:471` picks `packages.find(p => p.status === "active")` and `streaks[0]` across all orgs; filter both by the resolved `orgId`, and give `orgName` the same `profile.organizations[0]` fallback `orgId` already has.

- [ ] **Step 4: Fix the onboarding redirect loop**

`lib/auth.tsx` swallows `getMyProfile()` failures, leaving `onboarding_completed` undefined; `login/page.tsx:37` reads undefined as false and bounces an onboarded user to onboarding. Type the field `boolean | null`, leave it `null` when the profile fetch fails, and redirect only on an explicit `false`.

- [ ] **Step 5: Verify and commit**

```bash
npx tsc --noEmit && npm run build && npx playwright test
git add web/
git commit -m "fix: date parsing, org-scoped member picks, delete confirmation, onboarding redirect"
```

---

### Task 17: Remove the dead root migration file

`migrations/001_initial_schema.sql` is stale and unused — `Makefile` and `Dockerfile:30` both point at `db/migrations`. The two disagree on column names (`nfc_cards.hex_code` vs live `card_hex`, `attendance.method` vs live `check_in_method`), and every Go query uses the live names. Leaving it invites someone to edit the wrong file.

**Files:**
- Delete: `migrations/001_initial_schema.sql`

- [ ] **Step 1: Confirm nothing references it**

```bash
grep -rn "migrations/001_initial_schema\|migrations\b" Makefile Dockerfile docker-compose*.yml scripts/ .github/ | grep -v "db/migrations"
```

Expected: no hits outside `db/migrations`.

- [ ] **Step 2: Delete and verify the stack still builds**

```bash
git rm -r migrations/
docker compose build api 2>&1 | tail -5
```

- [ ] **Step 3: Commit**

```bash
git commit -m "chore: delete the stale root migration; db/migrations is the live set"
```

---

## Self-Review

**Spec coverage:** Every Phase 0 item in the design maps to a task — B0.1→T2, B0.2→T3, B0.3→T4, B0.4→T5, B0.5→T6, B0.6→T7, B0.7→T8, B0.8→T9, B0.9/B0.10→T10, F0.1→T11, F0.2→T12, F0.3→T13, F0.4→T14, F0.5/F0.6→T15, F0.7/F0.8/F0.9/F0.10/F0.11→T16. Task 1 (stale test signatures) and Task 17 (dead migration) come from the audit's findings table and are prerequisites/cleanup the spec references.

**Ordering:** Task 1 must run first — `go test ./...` does not compile until it does. Tasks 2 and 3 are sequential (3 extends 2's test file). Tasks 5–10 are independent of each other. Tasks 11–16 touch only the frontend and can run in parallel with 5–10. Task 17 is independent.

**Type consistency:** `UpdateOrgMember` gains `callerID, callerRole` as its first two parameters (Task 4) and nothing else calls it. Guide repository methods gain a leading `orgID` (Task 5). `BuyCredits` drops `rate`/`amount` (Task 6). `RecommendPlan` keeps its signature and only loses a side effect (Task 14). `parseDateOnly`/`nepalToday` are defined once in Task 16 and used only there.

**Known risk:** Task 2 revokes access from anyone currently relying on the buggy global role to administer a second gym. Before deploying, run on hetzner-1:

```sql
SELECT user_id, COUNT(DISTINCT organization_id) AS orgs,
       array_agg(DISTINCT role) AS roles
FROM organization_members WHERE status = 'active'
GROUP BY user_id HAVING COUNT(DISTINCT organization_id) > 1;
```

Any user whose roles differ across orgs is affected and should be told.
