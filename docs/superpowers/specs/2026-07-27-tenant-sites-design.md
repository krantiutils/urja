# Tenant Sites, Page Builder & Boxing Profile — Design

**Date:** 2026-07-27
**Status:** Approved for planning

## Goal

Give every gym on Urja a public website at `<slug>.nepalgym.xyz`, composed by the
gym's own admin from a library of sections, seeded from one of five templates.
First tenant: `ibckirtipur.nepalgym.xyz` (a boxing gym).

Alongside it: fix the critical and high-severity bugs two audits found, add a
combat-sports member profile, and drop the SaaS pricing section from the
marketing landing page.

## Scope

Six phases. Each is independently shippable and verifiable.

| Phase | Deliverable |
|-------|-------------|
| 0 | Critical + high bug fixes (backend authz, broken frontend features, timezone) |
| 1 | Subdomain routing, `internal/site` module, schema |
| 2 | Section library + five templates + public renderer |
| 3 | Admin page builder UI |
| 4 | Member boxing profile |
| 5 | Remove pricing from the marketing landing page |

Explicitly out of scope: custom domains (gym-owned apex domains), a media
library / image upload pipeline beyond URL fields, per-page A/B testing,
analytics, and any change to the mobile/Flutter clients.

---

## Phase 0 — Bug fixes

Two Sonnet audits ran against the repo. `go build ./...` passes and the Next.js
build is clean (`tsc --noEmit` reports zero errors, all 27 routes compile), so
everything below is a runtime or authorization defect. We fix criticals and
highs; mediums and lows are recorded in
`docs/superpowers/specs/2026-07-27-bug-audit-findings.md` for later.

### Backend — authorization

The root cause of six findings is one design error. `auth.Repository.GetUserByID`
(`internal/auth/repository.go:105`) derives the JWT's `role` claim with:

```sql
SELECT role FROM organization_members WHERE user_id = $1 LIMIT 1
```

No `ORDER BY`, no `status = 'active'` filter, and — critically — no relationship
to the org whose data is being accessed. `pkg/middleware/rbac.go` exposes both
`RequireRole` (checks that global claim) and `RequireOrgRole` (checks the
per-org role that `OrgScope` resolved for the URL's `{orgId}`). Ten route
groups mounted under `/orgs/{orgId}` use the former.

**B0.1 — Replace `RequireRole` with `RequireOrgRole`** on every route group
mounted under `/orgs/{orgId}`: `staff`, `packages`, `accounts`, `absentee`,
`feedback`, `workout`, `guide`, `smsapi`, `notice`, `dues` routes files.
`internal/member/routes.go` already does this correctly and is the reference.
Retain the global claim in the JWT only for `is_super_admin`.

**B0.2 — Gate `internal/subscription` routes.** `RegisterPackageRoutes` and
`RegisterMemberRoutes` currently apply no role middleware at all, so any active
org member can `POST .../packages/assign` with a client-chosen `amount_paid`.
Add `RequireOrgRole("staff", "admin")` to both.

**B0.3 — Block privilege self-escalation.** `member.UpdateOrgMember` accepts any
`role` value from a caller holding `staff`. Require `admin` to change the `role`
field, and reject requests where `{memberId}` resolves to the caller's own user
ID.

**B0.4 — Scope training guides to an org.** `training_guides` has no
`organization_id` column, yet `guide.RegisterOrgRoutes` mounts under
`/orgs/{orgId}`. Repository `GetByID`, `Update`, `SetPublished` and `Delete` key
on the guide ID alone, so any gym's staff can delete the platform-wide library.
Migration adds a nullable `organization_id` (NULL = global preset, mirroring
`workout_templates.is_preset`); all mutating repository methods gain an org
filter, and NULL-org guides become editable only by super admins.

**B0.5 — Server-owned SMS pricing.** `POST /orgs/{orgId}/sms/buy` accepts
client-supplied `quantity`, `rate` and `amount` and writes
`payment_status='completed'` with no verification. Move the rate to server
config, compute `amount` server-side, and for `payment_method=khalti` verify via
`khaltiClient.Lookup(pidx)` exactly as `billing.Subscribe` already does.

**B0.6 — Prevent Khalti payment replay.** `billing.Subscribe` verifies the
payment correctly but never records the consumed `pidx`, so one real payment can
be replayed across orgs. Add a `khalti_payments` table with a UNIQUE constraint
on `pidx`, insert inside the subscription transaction, and reject on conflict.

**B0.7 — Gate NFC org routes.** `internal/nfc/routes.go` applies no middleware,
letting any member register access-control readers and assign cards. Add
`RequireOrgRole("staff", "admin")`.

**B0.8 — Meter absentee SMS.** `absentee.Notify` calls the paid Aakash gateway
directly, bypassing the credit ledger that `smsapi.SendSMS` correctly uses.
Route it through credit deduction.

### Backend — timezone

Nepal is UTC+05:45 and the Postgres container runs UTC.

**B0.9 — Attendance calendar.** `GetMonthlyCalendar`
(`internal/attendance/repository.go:246`) bounds the query with a UTC-interpreted
`::date` cast but extracts the day with `AT TIME ZONE 'Asia/Kathmandu'`. Check-ins
between 00:00 and 05:45 Nepal time on the 1st land in the wrong month. Apply the
same zone to the range boundary.

**B0.10 — Nutrition dashboard default date.** `GetDailyDashboard` defaults to
`time.Now().Format(...)` — server UTC. Use the `nptLocation` fixed zone that
`internal/attendance/service.go:21` already defines.

### Frontend — broken features

**F0.1 — Money fields sent as strings.** `dashboard/accounts/page.tsx:156,166`
and `dashboard/due-payments/page.tsx:109` pass `amount` straight from a text
input, producing `"amount":"1500.50"`. The Go handlers decode into `float64`
without a `,string` tag, so every Add Transaction, Edit Transaction and Record
Payment request 400s. Send `Number(formAmount)`; validate it parses and is
positive before submitting.

**F0.2 — SMS sends to everyone.** `dashboard/sms/page.tsx:93` hardcodes
`member_ids: []`, and `smsapi.GetOrgMemberPhones` drops its `WHERE` clause when
the list is empty — so every send mass-texts the whole gym and burns credits.
Add a real member picker, block submit until at least one recipient is selected,
and add a confirmation dialog showing the recipient count and credit cost.

**F0.3 — Account deletion is dead.** `[lang]/delete/page.tsx` posts to
`/api/v1/auth/verify` (the route is `/auth/verify-otp`) and destructures `token`
from a response that returns `access_token`. Since this page exists for Google
Play compliance, it must actually work. Fix both, and add an e2e test.

**F0.4 — Recommendation silently overwrites the active plan.**
`workout.RecommendPlan` (`internal/workout/service.go:270`) calls
`SelfAssignPlan` and persists immediately, while the UI shows a preview with a
"Choose Plan" button that makes no API call. Split the backend into a read-only
`RecommendPlan` and an explicit `selfAssignPlan`; wire the confirm button to the
latter.

### Frontend — high

**F0.5 — Member search only filters the current page.**
`dashboard/members/page.tsx` fetches 20 rows and filters client-side, so
searching from page 2 reports "No members found". Pass `search` as a server-side
query param and reset `offset`, matching the pattern `due-payments` already uses.

**F0.6 — Attendance date filter over a capped fetch.** The page fetches the 100
most recent org-wide records and filters client-side by date, so older dates
show nothing. Pass the selected date to the server.

**F0.7 — Health page hides older metric types.** `getMyHealth({limit: 20})`
mixes all metric types before applying the limit, so a member who logs weight
daily loses their one BMI entry. Fetch per type.

**F0.8 — Package delete has no confirmation.** Every sibling page guards delete
with `confirm(...)`; `dashboard/packages/page.tsx` does not. Add it.

**F0.9 — Date-only values parsed as UTC.** `member/packages/page.tsx`,
`member/page.tsx` and `member/progress/page.tsx` call `new Date("YYYY-MM-DD")`,
which parses as UTC midnight and then renders in local time, shifting dates by a
day. Parse from components, as `nutrition/page.tsx` already does.

**F0.10 — Onboarding redirect loop on a transient failure.** `lib/auth.tsx`
swallows `getMyProfile()` errors, leaving `onboarding_completed` undefined; the
login page treats undefined as false and bounces a fully-onboarded user to
onboarding. Model the field as `boolean | null` and redirect only on explicit
`false`.

**F0.11 — Member dashboard "today" in UTC.** `member/page.tsx:42` computes the
date via `toISOString()` at module scope while the nutrition page uses local
components — during the 00:00–05:45 window the dashboard queries yesterday.
Compute it once, in Nepal time, inside the effect.

---

## Phase 1 — Subdomain routing, schema, API

### Routing

`web/src/middleware.ts` gains subdomain resolution ahead of the existing locale
logic:

1. Read `Host`. Strip the port.
2. If the host is `<sub>.nepalgym.xyz` and `<sub>` is not in
   `{www, api, admin}` — or, in development, `<sub>.localhost` — treat `<sub>`
   as an org slug.
3. Parse the optional locale prefix from the path, then rewrite to
   `/{locale}/site/{sub}/{restOfPath}`, defaulting `restOfPath` to `home`.
4. Set `x-site-slug` on the response headers for downstream use.
5. Apex and `www` hosts fall through to today's behavior unchanged.

Requests to `/api/*`, `/_next/*` and any path containing a dot bypass rewriting,
as they do today. A visitor who reaches `<sub>.nepalgym.xyz/dashboard` is
redirected to `https://nepalgym.xyz/{locale}/dashboard` — the app shell stays on
the apex, tenant subdomains serve public marketing pages only.

Unknown slug renders the site-level `not-found`, not a crash.

### Traefik

`docker-compose.prod.yml` gains a wildcard router for the web service, using the
`letsencrypt-dns` resolver that already exists on hetzner-1 (Porkbun DNS-01,
keys already in `/home/ubuntu/traefik/.env`). `*.nepalgym.xyz` already has an A
record pointing at the server, so no DNS change is needed.

```yaml
- traefik.http.routers.urja-sites.rule=HostRegexp(`.+\.nepalgym\.xyz`)
- traefik.http.routers.urja-sites.entrypoints=websecure
- traefik.http.routers.urja-sites.tls.certresolver=letsencrypt-dns
- traefik.http.routers.urja-sites.tls.domains[0].main=nepalgym.xyz
- traefik.http.routers.urja-sites.tls.domains[0].sans=*.nepalgym.xyz
- traefik.http.routers.urja-sites.priority=1
```

The existing API router keeps priority 2 so `/api/` still wins on the apex. A
matching wildcard rule is added to the API router so tenant pages can call the
public site endpoints same-origin.

### Schema — migration `000045_create_site_pages`

```sql
CREATE TABLE site_pages (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    slug            VARCHAR(120) NOT NULL,
    title           VARCHAR(200) NOT NULL,
    title_ne        VARCHAR(200),
    sections        JSONB NOT NULL DEFAULT '[]',
    seo_description TEXT,
    is_published    BOOLEAN NOT NULL DEFAULT false,
    show_in_nav     BOOLEAN NOT NULL DEFAULT true,
    sort_order      INT NOT NULL DEFAULT 0,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (organization_id, slug)
);

CREATE TABLE site_settings (
    organization_id UUID PRIMARY KEY REFERENCES organizations(id) ON DELETE CASCADE,
    template        VARCHAR(50) NOT NULL DEFAULT 'fight_club',
    theme           JSONB NOT NULL DEFAULT '{}',
    nav             JSONB NOT NULL DEFAULT '{}',
    footer          JSONB NOT NULL DEFAULT '{}',
    socials         JSONB NOT NULL DEFAULT '{}',
    is_live         BOOLEAN NOT NULL DEFAULT false,
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE site_leads (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    name            VARCHAR(200) NOT NULL,
    phone           VARCHAR(30) NOT NULL,
    email           VARCHAR(255),
    message         TEXT,
    interest        VARCHAR(100),
    source_page     VARCHAR(120),
    status          VARCHAR(30) NOT NULL DEFAULT 'new',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_site_pages_org ON site_pages(organization_id, sort_order);
CREATE INDEX idx_site_leads_org ON site_leads(organization_id, created_at DESC);
```

`is_live` on `site_settings` is the master switch: until an admin flips it, the
subdomain serves a "coming soon" page rather than a half-built site.

### `internal/site` module

Follows the repository/service/handler/routes layout every other module uses.

Public, unauthenticated (`/api/v1/sites`):

| Method | Path | Purpose |
|---|---|---|
| GET | `/{slug}` | Settings, nav, and the published page index |
| GET | `/{slug}/pages/{pageSlug}` | One published page with its sections |
| POST | `/{slug}/leads` | Submit an enquiry — rate-limited, honeypot field |

Admin, under `/orgs/{orgId}/site` with `RequireOrgRole("admin")`:

| Method | Path | Purpose |
|---|---|---|
| GET/PUT | `/settings` | Template, theme, nav, footer, go-live |
| GET/POST | `/pages` | List (incl. drafts) / create |
| GET/PUT/DELETE | `/pages/{pageId}` | Read / update sections / delete |
| POST | `/apply-template` | Seed pages and theme from a template |
| GET/PATCH | `/leads` | Inbox and status updates |

Section JSON is validated server-side against the known section types on write —
an unknown `type` or a malformed `content` shape is a 400, so the public
renderer never has to defend against garbage. All rich text is stored as
Markdown and sanitized at render, never as raw HTML.

---

## Phase 2 — Section library, templates, public renderer

### Theming

Each template defines a token set written as CSS custom properties on the site
root (`--site-bg`, `--site-surface`, `--site-fg`, `--site-fg-muted`,
`--site-accent`, `--site-accent-fg`, `--site-radius`, `--site-font-display`,
`--site-font-body`). Section components consume them through Tailwind arbitrary
values (`bg-[var(--site-bg)]`). This keeps five visually distinct templates out
of the Tailwind config and lets an admin override individual tokens.

### Section types

Nineteen types. Every section carries `{ id, type, variant, content, style }`
where `style` is `{ background, padding, width, align }`. All text fields are
bilingual (`title` / `titleNe`).

| Type | Variants |
|---|---|
| `hero` | centered, split, fullbleed, minimal |
| `stats_bar` | inline, cards, bordered |
| `class_timetable` | table, day_tabs, cards |
| `coaches` | cards, list, spotlight |
| `programs_grid` | cards, icons, list, numbered |
| `membership_plans` | cards, table, compact |
| `gallery` | grid, masonry, carousel |
| `testimonials` | cards, carousel, quote |
| `faq` | accordion, list, two_column |
| `cta_banner` | solid, gradient, image, split |
| `contact_info` | list, card, two_column |
| `map_embed` | standard, with_info, full_width |
| `lead_form` | inline, card, split |
| `opening_hours` | table, list, compact |
| `logo_strip` | grid, marquee, simple |
| `fight_record` | timeline, table, cards |
| `rich_text` | standard, two_column, centered |
| `media` | image, video |
| `divider` | line, dots, space |

`class_timetable`, `coaches` and `fight_record` hold their content inline in the
section JSON — they are not backed by dedicated entities in this phase, so an
admin types classes and coaches directly into the builder. `membership_plans`
supports two data sources: `auto`, which reads the org's active packages from
the existing packages API, and `manual`. `fight_record` gains an optional `auto`
source in Phase 4, reading a named member's `bout_records`; its `manual` source
remains the default so the section works standalone.

Each renderer lives in `web/src/components/site/sections/<Type>Renderer.tsx`,
dispatched by `SectionRenderer.tsx`. Renderers are server components except
where interaction demands otherwise (`faq` accordion, `gallery` carousel,
`class_timetable` day tabs, `lead_form`).

### Templates

Five, each a theme token set plus a preset array of pages and sections.

1. **Fight Club** — near-black, crimson accent, condensed display type, sharp
   corners. Fullbleed hero, stats, programs, timetable, coaches, masonry
   gallery, testimonials, image CTA.
2. **Iron & Sweat** — industrial: charcoal on off-white, safety-yellow accent,
   monospace display, zero radius, heavy rules. Split hero, logo strip,
   numbered programs, table timetable, coach list, fight record timeline, FAQ,
   split lead form.
3. **Champion** — light editorial, photo-forward, serif display, gold accent,
   generous whitespace. Centered hero, rich text, coach spotlight, grid gallery,
   membership cards, pull-quote testimonials, gradient CTA.
4. **Community** — warm and approachable: orange and teal, rounded corners.
   Split hero, stat cards, icon programs, opening hours, timetable cards,
   coaches, accordion FAQ, card lead form, map.
5. **Minimal Pro** — typographic and restrained, single accent, heavy whitespace.
   Minimal hero, rich text, list programs, table timetable, coach list,
   membership table, list FAQ, two-column contact.

Each seeds four pages: Home, Classes, Coaches, Contact.

### Public site shell

`web/src/app/[lang]/site/[slug]/` holds the layout (nav from `site_settings.nav`
plus pages with `show_in_nav`, footer, locale switcher, a "Member Login" link
back to the apex) and `[[...page]]/page.tsx` for the page itself. Rendered
server-side with `generateMetadata` driving per-page title, description and
Open Graph tags. Unknown or unpublished slugs render `not-found`. Pages are
statically revalidated with a short TTL and revalidated on publish.

---

## Phase 3 — Admin page builder

`/[lang]/dashboard/site` — three routes:

- **Overview** — go-live toggle, template picker with visual previews, the
  public URL, and the page list with publish state and reordering.
- **Builder** (`/pages/{id}`) — three panes: a section list on the left
  (reorder, duplicate, delete, visibility), a live preview in the middle using
  the exact public renderers, and an inspector on the right (variant picker,
  content fields, style tokens). "Add section" opens a categorized panel.
- **Leads** — the enquiry inbox with status transitions.

State lives in a `useSiteBuilder` hook; changes autosave after a 1.5 s debounce
with an explicit saved/saving/error indicator. Reordering uses native HTML5
drag-and-drop plus keyboard-accessible move up/down buttons — no new dependency.

Applying a template to a site that already has pages requires confirmation and
replaces content wholesale.

---

## Phase 4 — Member boxing profile

Migration `000046_boxing_profile`:

```sql
CREATE TABLE member_boxing_profiles (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    stance          VARCHAR(20),        -- orthodox | southpaw | switch
    weight_class    VARCHAR(40),
    skill_level     VARCHAR(30),        -- beginner | intermediate | amateur | pro
    sparring_cleared BOOLEAN NOT NULL DEFAULT false,
    sparring_cleared_at TIMESTAMPTZ,
    sparring_cleared_by UUID REFERENCES users(id),
    reach_cm        NUMERIC(5,1),
    notes           TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (user_id, organization_id)
);

CREATE TABLE bout_records (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    bout_date       DATE NOT NULL,
    opponent        VARCHAR(200),
    event_name      VARCHAR(200),
    result          VARCHAR(20) NOT NULL,   -- win | loss | draw | no_contest
    method          VARCHAR(30),            -- ko | tko | decision | submission | dq
    rounds          INT,
    weight_class    VARCHAR(40),
    notes           TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

Weight class is derived from the member's latest recorded weight in
`health_metrics` using standard amateur boxing divisions, and can be overridden.
Sparring clearance is staff-only — a member cannot clear themselves.

Routes: `GET/PUT /members/me/boxing` and `GET/POST/DELETE /members/me/bouts` for
self-service; `GET/PUT /orgs/{orgId}/members/{memberId}/boxing` and the matching
bout routes for staff, including the clearance toggle. Member UI extends the
existing profile page; staff UI extends the member detail view. The
`fight_record` site section can read a coach's bouts, tying Phase 4 back to
Phase 2.

---

## Phase 5 — Remove pricing

Delete the `pricingPlans` array, the `<section id="pricing">` block, the desktop
nav link, the mobile nav link in `LandingInteractive.tsx`, and the footer link
from `web/src/app/[lang]/page.tsx`. Remove the now-unused `landing.pricing*`
keys from both dictionaries in `lib/i18n.ts` and the `Check` icon import if it
becomes unused. Gym membership packages are untouched.

---

## Testing

- **Go** — table-driven unit tests for slug/subdomain resolution, section JSON
  validation, weight-class derivation, and the lead rate limiter. Authorization
  regression tests asserting that a member of Org B with an admin role in Org A
  receives 403 on Org B's admin routes — one per fixed route group.
- **Playwright** — the repo already has specs in `web/tests/`. Add: tenant site
  renders for each of the five templates, subdomain rewriting including the
  locale prefix, unknown-slug 404, lead submission landing in the dashboard
  inbox, the page builder's add/reorder/publish cycle, and a regression test per
  Phase 0 frontend fix (transaction amount, SMS recipient guard, account
  deletion).
- **Manual** — drive the app with the `run` skill against `ibckirtipur.localhost:3000`.

## Verification

Each phase must show: `go build ./...`, `go vet ./...`, `go test ./...`,
`npx tsc --noEmit`, `npm run build`, and the relevant Playwright specs passing —
with output, before the phase is called done.

## Risks

- **Wildcard cert issuance** is the one step that can't be verified locally.
  The Porkbun DNS-01 resolver and the `*.nepalgym.xyz` A record already exist for
  `*.doctorsewa.org`'s sake, so the path is proven, but first issuance for a new
  wildcard SAN should be watched in `docker logs traefik | grep -i acme`.
- **Phase 0 authz changes can lock people out.** Anyone currently relying on the
  buggy global role to administer a second gym loses access — correctly, but
  visibly. Worth checking `organization_members` for users with roles in more
  than one org before deploying.
- **Section JSON is a schema without migrations.** Adding a required field to a
  section type later means handling old rows. Renderers must tolerate missing
  fields from day one, and a `version` field on each section leaves room for a
  migration helper.
