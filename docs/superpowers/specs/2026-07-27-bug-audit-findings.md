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

Findings 39–42 below were not in either audit. They were found by running the
stack and using it — three of them only became visible once there was a tenant
site to load.

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

## Found while restoring the test baseline (not in either audit)

**Workout self-routes accept a missing `organization_id` and write orphaned rows.**
`tests/e2e/workout_test.go` has three tests — `TestWorkout_GetMyPlan_MissingOrgID`,
`TestWorkout_CreateLog_MissingOrgID`, `TestWorkout_ListLogs_MissingOrgID` —
asserting a 400 when `organization_id` is absent. The handlers
(`internal/workout/handler.go:200,220,244`) never validate it: `GetMyPlan` reads
the query param into an empty string, and `CreateMyLog` persists a log with
`organization_id: null`.

**Fixed.** `workout.Service.ResolveOrg` now resolves the organization from real
membership: a member of one gym does not have to name it, a supplied
organization is verified against membership, and it is an error only when the
member belongs to several. Migration `000051_backfill_workout_org` repairs the
rows already stored, limited to owners with exactly one active membership —
a member of two gyms could have trained at either, and guessing would file
their work under the wrong gym. The three tests pass, alongside three new ones
covering inference, a foreign org and the ambiguous case.

## Found while building and running the tenant site

| # | Finding | Location | Status |
|---|---------|----------|--------|
| 39 | `GET /api/v1/packages` ignored the org entirely and returned every gym's active packages to any caller — a cross-tenant price-list leak, and the wrong prices on a tenant site's `membership_plans` | `internal/packages/repository.go:61` | Fixed |
| 40 | CORS matched origins against a static allow-list, so every newly created gym's enquiry form was blocked by the browser with nothing in the server logs | `cmd/api/main.go:246` | Fixed — `pkg/middleware/cors.go` |
| 41 | Internal links written by an admin or stored in section content kept their English path on the Nepali site, sending Nepali visitors to English pages | `web/src/components/site/SiteHeader.tsx`, `lib/site/content.ts` | Fixed — `lib/site/links.ts` |
| 42 | The enquiry form posted cross-origin to the API, making a gym's only conversion path depend on `NEXT_PUBLIC_API_URL` being baked in at image-build time; a bad build fails silently | `web/src/components/site/sections/LeadFormRenderer.tsx` | Fixed — same-origin proxy at `app/api/site-leads` |

Findings 40 and 42 were only reachable by loading a real page in a browser and
clicking the button: both `go build` and `npm run build` were clean throughout,
and 42 in particular would have passed every test that did not involve a
browser actually submitting the form.

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

---

## Before deploying

**Check for multi-org users whose roles differ.** The authorization fix
(finding 1) replaced the global JWT `role` claim with the per-org role. A user
who is an admin of gym A and a plain member of gym B previously carried one
global role into both; they now get gym B's real role there. That is the
intended behavior, but it is a visible change for anyone affected, so it is
worth knowing who they are first:

```sql
SELECT u.phone, u.name, count(DISTINCT om.role) AS distinct_roles,
       array_agg(o.slug || '=' || om.role ORDER BY o.slug) AS roles
FROM organization_members om
JOIN users u ON u.id = om.user_id
JOIN organizations o ON o.id = om.organization_id
WHERE om.status = 'active'
GROUP BY u.id, u.phone, u.name
HAVING count(DISTINCT om.role) > 1;
```

Empty result means nobody's access changes.

**Wildcard DNS already resolves.** Verified: an arbitrary label
(`zzz-does-not-exist.nepalgym.xyz`) returns `178.104.21.224`, so the
`*.nepalgym.xyz` A record exists and points at hetzner-1. Nothing to do here.

**The wildcard certificate does not exist yet.** A TLS handshake against
hetzner-1 with SNI `ibckirtipur.nepalgym.xyz` returns `CN = TRAEFIK DEFAULT
CERT`; the apex serves a real certificate covering only `nepalgym.xyz` and
`www.nepalgym.xyz`. Ignoring the certificate, Traefik answers 404 — no router
matches tenant hosts. So a visitor to a gym site today gets a browser security
warning and then nothing.

Both follow from the same cause: the deployed `docker-compose.prod.yml` predates
the tenant-site work. The router rule, the `letsencrypt-dns` resolver and the
wildcard SANs are all in the repo's copy now, and issuance happens on its own
once Traefik sees them — provided the Porkbun credentials are present in
`/home/ubuntu/traefik/.env`, which is where `deploy.md` says they live (they are
*not* in `~/.fmw`, which holds unrelated keys).

**The deploy workflow never shipped the compose file.** It ran `docker compose
pull && up -d` against whatever copy was already on the server, so a routing or
TLS change merged to this repo would silently never take effect. Fixed by adding
an scp step ahead of the deploy.

**`SITE_BASE_DOMAIN`.** The API reads this at runtime to decide which origins are
tenant sites. Unset, tenant subdomains fall back to exact-match CORS and every
gym's enquiry form is blocked by the browser. It defaults correctly in compose.

**`NEXT_PUBLIC_BASE_DOMAIN` is a build arg, not a runtime variable.** Next inlines
`NEXT_PUBLIC_*` at build time, including into the middleware that does the
subdomain routing. It was briefly set under `environment:` in compose, which
looks like it works and changes nothing — it is now a Docker build arg passed by
the workflow. This only matters for a domain other than the default, e.g.
staging.
