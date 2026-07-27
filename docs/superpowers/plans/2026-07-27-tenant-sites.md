# Tenant Subdomain Sites & Page Builder — Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Every gym gets a public website at `<slug>.nepalgym.xyz`, composed by its own admin from a library of 19 section types, seeded from one of five visual templates.

**Architecture:** The org `slug` already unique in `organizations` *is* the subdomain — no new identity concept. Next.js middleware rewrites `<slug>.nepalgym.xyz/<path>` to the internal route `/{locale}/site/{slug}/{path}`. Page content is an ordered array of section objects in a JSONB column, validated server-side on write so the public renderer never defends against malformed data. Each section renders through one component keyed by `type`, with a `variant` selecting layout and CSS custom properties supplying the template's theme.

**Tech Stack:** Go 1.22 + chi v5 + pgx v5; Next.js 14 App Router + TypeScript + Tailwind; Traefik with the existing Porkbun DNS-01 resolver; Playwright.

## Global Constraints

- Migrations live in `db/migrations/`, format `NNNNNN_name.up.sql` / `.down.sql`. **Reserved numbers: 000045 = training_guides org, 000046 = khalti_payments, 000047 = site tables, 000048 = boxing profile.** Do not reuse.
- Every user-facing string is bilingual. Section content carries both `title` and `titleNe`; the two dictionaries in `web/src/lib/i18n.ts` must stay symmetric.
- No new runtime dependencies, frontend or backend. Drag-and-drop uses native HTML5 DnD.
- Tenant sites are **public marketing pages only**. The authenticated app stays on the apex host.
- Rich text is stored as Markdown and sanitized at render. Never store or render raw HTML from user input.
- Every section renderer must tolerate missing or unknown fields without throwing — section JSON has no migration story, so old rows must keep rendering.
- Verification gate: `go build ./...`, `go vet ./...`, `go test ./...`, `npx tsc --noEmit`, `npm run build`, plus the task's Playwright specs.

---

### Task 1: Site schema and the `internal/site` module

**Files:**
- Create: `db/migrations/000047_create_site_tables.up.sql` / `.down.sql`
- Create: `internal/site/{models,repository,service,handler,routes}.go`
- Create: `internal/site/service_test.go`
- Modify: `cmd/api/main.go` (wire the module, register routes)

**Interfaces:**
- Consumes: `middleware.OrgScope`, `middleware.RequireOrgRole` from `pkg/middleware`
- Produces:
  - `site.NewRepository(pool *pgxpool.Pool) *Repository`
  - `site.NewService(repo *Repository, logger *slog.Logger) *Service`
  - `site.NewHandler(svc *Service, logger *slog.Logger) *Handler`
  - `(*Handler).RegisterPublicRoutes(r chi.Router)` mounted at `/api/v1/sites`
  - `(*Handler).RegisterOrgRoutes(r chi.Router)` mounted at `/orgs/{orgId}/site`
  - `Section struct { ID, Type, Variant string; Content json.RawMessage; Style SectionStyle; Version int }`

- [ ] **Step 1: Write the migration**

```sql
BEGIN;

CREATE TABLE IF NOT EXISTS site_pages (
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

CREATE TABLE IF NOT EXISTS site_settings (
    organization_id UUID PRIMARY KEY REFERENCES organizations(id) ON DELETE CASCADE,
    template        VARCHAR(50) NOT NULL DEFAULT 'fight_club',
    theme           JSONB NOT NULL DEFAULT '{}',
    nav             JSONB NOT NULL DEFAULT '{}',
    footer          JSONB NOT NULL DEFAULT '{}',
    socials         JSONB NOT NULL DEFAULT '{}',
    is_live         BOOLEAN NOT NULL DEFAULT false,
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS site_leads (
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

CREATE INDEX IF NOT EXISTS idx_site_pages_org ON site_pages(organization_id, sort_order);
CREATE INDEX IF NOT EXISTS idx_site_leads_org ON site_leads(organization_id, created_at DESC);

COMMIT;
```

The `.down.sql` drops the three tables in reverse dependency order.

- [ ] **Step 2: Write the failing service tests**

Section validation is the security boundary — it gets the most tests:

```go
func TestValidateSections_RejectsUnknownType(t *testing.T) {
	err := ValidateSections([]Section{{Type: "arbitrary_code", Variant: "x"}})
	if err == nil { t.Fatal("unknown section type must be rejected") }
}

func TestValidateSections_RejectsUnknownVariantForType(t *testing.T) {
	err := ValidateSections([]Section{{Type: "hero", Variant: "nonsense"}})
	if err == nil { t.Fatal("variant must be valid for its type") }
}

func TestValidateSections_RejectsScriptInRichText(t *testing.T) {
	// Markdown is stored raw but must never contain an HTML script block.
}

func TestValidateSections_AcceptsEveryDeclaredTypeAndVariant(t *testing.T) {
	// Table-driven over the full SectionSpecs registry — this is the test that
	// catches a renderer and the validator disagreeing about what exists.
}

func TestReservedSlugsRejected(t *testing.T) {
	// "api", "www", "admin", "_next" must never be usable as a page slug.
}
```

- [ ] **Step 3: Run to confirm failure**

```bash
go test ./internal/site/ -v
```

Expected: package does not compile — nothing exists yet.

- [ ] **Step 4: Implement the module**

Follow the exact repository/service/handler/routes layout of `internal/notice` (a small, recent module — read it first). `models.go` holds the `Section` struct and a `SectionSpecs` registry mapping each type to its allowed variants; the registry is the single source of truth shared by validation and, mirrored in TypeScript, by the builder UI.

Public routes (no auth), mounted at `/api/v1/sites`:

| Method | Path | Notes |
|---|---|---|
| GET | `/{slug}` | Settings, nav, published page index. 404 if `is_live` is false. |
| GET | `/{slug}/pages/{pageSlug}` | One published page. 404 if unpublished. |
| POST | `/{slug}/leads` | Rate-limited per IP, honeypot field, max body size. |

Admin routes under `/orgs/{orgId}/site`, all `RequireOrgRole("admin")`:

| Method | Path |
|---|---|
| GET / PUT | `/settings` |
| GET / POST | `/pages` |
| GET / PUT / DELETE | `/pages/{pageId}` |
| POST | `/apply-template` |
| GET | `/leads` |
| PATCH | `/leads/{leadId}` |

Every repository query filters by `organization_id`. The public lookups resolve the org by slug first, then scope everything to that org ID.

- [ ] **Step 5: Verify and commit**

```bash
make migrate-up && go test ./internal/site/ -v && go build ./... && go vet ./...
git add db/migrations/000047_* internal/site/ cmd/api/main.go
git commit -m "feat: add site pages, settings and leads with server-side section validation"
```

---

### Task 2: Subdomain middleware and Traefik routing

**Files:**
- Modify: `web/src/middleware.ts`
- Create: `web/src/lib/subdomain.ts`
- Create: `web/src/lib/subdomain.test.ts`
- Modify: `docker-compose.prod.yml`

**Interfaces:**
- Produces: `extractSubdomain(host: string): string | null` and `resolveSiteRewrite(host: string, pathname: string): { slug: string, locale: string, pagePath: string } | null`

- [ ] **Step 1: Write the failing unit tests**

Pure functions, so test them directly — no browser needed:

```ts
// Production hosts
expect(extractSubdomain("ibckirtipur.nepalgym.xyz")).toBe("ibckirtipur");
expect(extractSubdomain("nepalgym.xyz")).toBeNull();
expect(extractSubdomain("www.nepalgym.xyz")).toBeNull();
expect(extractSubdomain("api.nepalgym.xyz")).toBeNull();
// Port must be stripped
expect(extractSubdomain("ibckirtipur.localhost:3000")).toBe("ibckirtipur");
// Rewrites
expect(resolveSiteRewrite("ibckirtipur.nepalgym.xyz", "/"))
  .toEqual({ slug: "ibckirtipur", locale: "en", pagePath: "home" });
expect(resolveSiteRewrite("ibckirtipur.nepalgym.xyz", "/ne/coaches"))
  .toEqual({ slug: "ibckirtipur", locale: "ne", pagePath: "coaches" });
expect(resolveSiteRewrite("nepalgym.xyz", "/en/dashboard")).toBeNull();
```

- [ ] **Step 2: Run to confirm failure**

- [ ] **Step 3: Implement**

`extractSubdomain` strips the port, matches against the base domain from `process.env.NEXT_PUBLIC_BASE_DOMAIN` (default `nepalgym.xyz`), also accepts `<sub>.localhost` for development, and returns `null` for the reserved labels `www`, `api`, `admin`.

In `middleware.ts`, run subdomain resolution **before** the existing locale redirect. On a match, rewrite to `/{locale}/site/{slug}/{pagePath}` and set an `x-site-slug` response header. Preserve today's early-outs for `/_next`, `/api` and any path containing a dot. A request to a dashboard path on a tenant subdomain redirects to the apex — the app shell lives only on `nepalgym.xyz`.

- [ ] **Step 4: Add the Traefik wildcard router**

In `docker-compose.prod.yml`, add to the `web` service labels. The `letsencrypt-dns` resolver and the `*.nepalgym.xyz` A record already exist on hetzner-1 (see `~/deploy.md`), so no DNS or resolver setup is needed:

```yaml
- traefik.http.routers.urja-sites.rule=HostRegexp(`.+\.nepalgym\.xyz`)
- traefik.http.routers.urja-sites.entrypoints=websecure
- traefik.http.routers.urja-sites.tls.certresolver=letsencrypt-dns
- traefik.http.routers.urja-sites.tls.domains[0].main=nepalgym.xyz
- traefik.http.routers.urja-sites.tls.domains[0].sans=*.nepalgym.xyz
- traefik.http.routers.urja-sites.priority=1
```

Add a matching wildcard rule to the `api` router at priority 3 so tenant pages can call `/api/` same-origin.

- [ ] **Step 5: Verify and commit**

```bash
cd web && npx tsc --noEmit && npm run build
git add web/src docker-compose.prod.yml
git commit -m "feat: route gym subdomains to tenant site pages"
```

---

### Task 3: Section type registry and theme tokens

**Files:**
- Create: `web/src/types/site.ts`
- Create: `web/src/lib/site/section-specs.ts`
- Create: `web/src/lib/site/themes.ts`
- Create: `web/src/lib/site/__tests__/section-specs.test.ts`

**Interfaces:**
- Produces:
  - `SectionType` — union of the 19 type strings
  - `SECTION_SPECS: Record<SectionType, { variants: string[]; defaultContent: unknown; label: string; labelNe: string; category: string }>`
  - `Section { id: string; type: SectionType; variant: string; content: Record<string, unknown>; style: SectionStyle; version: number }`
  - `SectionStyle { background: 'base'|'surface'|'accent'|'none'; padding: 'none'|'sm'|'md'|'lg'; width: 'full'|'contained'|'narrow'; align: 'left'|'center'|'right' }`
  - `THEMES: Record<TemplateId, ThemeTokens>`

- [ ] **Step 1: Write the failing test**

The registry must agree with the Go validator, or the builder will offer sections the API rejects:

```ts
test("every section type declares at least one variant", () => {
  for (const [type, spec] of Object.entries(SECTION_SPECS)) {
    expect(spec.variants.length).toBeGreaterThan(0);
  }
});

test("registry matches the Go SectionSpecs registry", () => {
  // Parse internal/site/models.go and compare the type+variant sets.
  // This is the guard against the two drifting apart.
});

test("every default content object round-trips through JSON unchanged", () => {
  for (const spec of Object.values(SECTION_SPECS)) {
    expect(JSON.parse(JSON.stringify(spec.defaultContent))).toEqual(spec.defaultContent);
  }
});
```

- [ ] **Step 2: Run to confirm failure**

- [ ] **Step 3: Define the 19 types**

| Type | Variants | Category |
|---|---|---|
| `hero` | centered, split, fullbleed, minimal | header |
| `stats_bar` | inline, cards, bordered | social proof |
| `class_timetable` | table, day_tabs, cards | gym |
| `coaches` | cards, list, spotlight | gym |
| `programs_grid` | cards, icons, list, numbered | gym |
| `membership_plans` | cards, table, compact | gym |
| `gallery` | grid, masonry, carousel | media |
| `testimonials` | cards, carousel, quote | social proof |
| `faq` | accordion, list, two_column | content |
| `cta_banner` | solid, gradient, image, split | conversion |
| `contact_info` | list, card, two_column | contact |
| `map_embed` | standard, with_info, full_width | contact |
| `lead_form` | inline, card, split | conversion |
| `opening_hours` | table, list, compact | contact |
| `logo_strip` | grid, marquee, simple | social proof |
| `fight_record` | timeline, table, cards | gym |
| `rich_text` | standard, two_column, centered | content |
| `media` | image, video | media |
| `divider` | line, dots, space | layout |

- [ ] **Step 4: Define the five theme token sets**

Each theme supplies `--site-bg`, `--site-surface`, `--site-fg`, `--site-fg-muted`, `--site-accent`, `--site-accent-fg`, `--site-border`, `--site-radius`, `--site-font-display`, `--site-font-body`. The five: `fight_club` (near-black, crimson, condensed display, sharp), `iron_sweat` (charcoal on off-white, safety yellow, mono display, zero radius), `champion` (light editorial, serif display, gold, generous whitespace), `community` (warm orange and teal, rounded), `minimal_pro` (typographic, restrained, single accent).

These must be visually distinct at a glance, not five recolorings of the same layout. Vary type scale, weight, border treatment and density — not just hue.

- [ ] **Step 5: Verify and commit**

```bash
npx tsc --noEmit && npm test -- section-specs
git add web/src/types/site.ts web/src/lib/site/
git commit -m "feat: define the site section registry and five theme token sets"
```

---

### Task 4: Section renderers

**Files:**
- Create: `web/src/components/site/SectionRenderer.tsx`
- Create: `web/src/components/site/sections/*Renderer.tsx` (19 files)
- Create: `web/src/components/site/primitives/{SectionShell,SiteImage,SiteButton}.tsx`
- Test: `web/tests/site-sections.spec.ts`

**Interfaces:**
- Consumes: `Section`, `SECTION_SPECS` from Task 3
- Produces: `<SectionRenderer section={Section} locale={Locale} org={OrgSummary} />` — dispatches on `section.type`, renders `null` for an unknown type rather than throwing

- [ ] **Step 1: Write the failing test**

One spec that mounts a page containing every section type in every variant and asserts no console errors and no empty render:

```ts
test("every section type and variant renders without error", async ({ page }) => {
  // Fixture page seeded with all 19 types × all variants.
  const errors: string[] = [];
  page.on("pageerror", (e) => errors.push(e.message));
  await page.goto("/en/site/__fixture__/all-sections");
  expect(errors).toEqual([]);
  // Each section root must be present and non-empty.
});

test("a section with missing content fields still renders", async ({ page }) => {
  // Guards the no-migrations-for-JSON risk called out in the design.
});
```

- [ ] **Step 2: Run to confirm failure**

- [ ] **Step 3: Build the shell primitive first**

`SectionShell` owns background, padding, width and alignment from `section.style` so no individual renderer reimplements them. Every renderer wraps its content in it.

- [ ] **Step 4: Implement the 19 renderers**

Server components by default. Only these need `"use client"`: `faq` (accordion), `gallery` (carousel), `class_timetable` (day tabs), `lead_form` (submission), `logo_strip` (marquee). Every text field reads the `Ne` suffixed variant when `locale === "ne"`, falling back to the base field when the translation is empty.

`membership_plans` with `dataSource: "auto"` fetches the org's active packages server-side; `manual` reads from `content.plans`. `fight_record` behaves the same way once Task 9 lands.

- [ ] **Step 5: Verify and commit**

```bash
npx tsc --noEmit && npm run build && npx playwright test tests/site-sections.spec.ts
git add web/src/components/site/
git commit -m "feat: add 19 section renderers with variant support"
```

---

### Task 5: Public site shell and page route

**Files:**
- Create: `web/src/app/[lang]/site/[slug]/layout.tsx`
- Create: `web/src/app/[lang]/site/[slug]/[[...page]]/page.tsx`
- Create: `web/src/app/[lang]/site/[slug]/not-found.tsx`
- Create: `web/src/components/site/{SiteHeader,SiteFooter}.tsx`
- Test: `web/tests/site-public.spec.ts`

**Interfaces:**
- Consumes: `GET /api/v1/sites/{slug}` and `/{slug}/pages/{pageSlug}` from Task 1; `SectionRenderer` from Task 4

- [ ] **Step 1: Write the failing tests** — a published page renders its sections; an unpublished page 404s; a site with `is_live=false` shows the coming-soon page; an unknown slug 404s; the locale switcher preserves the current page; `generateMetadata` emits the page title and description.

- [ ] **Step 2: Run to confirm failure**

- [ ] **Step 3: Implement** — the layout injects the theme tokens as inline CSS custom properties on a wrapper element, renders `SiteHeader` (nav from settings plus `show_in_nav` pages, locale switcher, a Member Login link to the apex) and `SiteFooter`. The page route resolves `[[...page]]` to a page slug defaulting to `home`, fetches, and maps sections through `SectionRenderer`. Use `revalidate` with a short TTL.

- [ ] **Step 4: Verify and commit**

```bash
npx playwright test tests/site-public.spec.ts && npm run build
git add web/src/app/\[lang\]/site web/src/components/site/
git commit -m "feat: render public tenant sites with per-template theming"
```

---

### Task 6: Five templates and the seed

**Files:**
- Create: `web/src/lib/site/templates.ts`
- Create: `internal/site/templates.go`
- Create: `db/migrations/000049_seed_ibckirtipur.up.sql` / `.down.sql`
- Test: `internal/site/templates_test.go`

**Interfaces:**
- Produces: `TEMPLATES: Record<TemplateId, { id, name, nameNe, theme: ThemeTokens, pages: TemplatePage[] }>` where `TemplatePage { slug, title, titleNe, sections: Section[] }`

- [ ] **Step 1: Write the failing test**

```go
func TestTemplates_AllSectionsValidate(t *testing.T) {
	// Every section in every template must pass ValidateSections. This is the
	// test that catches a template referencing a variant that doesn't exist.
	for id, tpl := range Templates {
		for _, page := range tpl.Pages {
			if err := ValidateSections(page.Sections); err != nil {
				t.Errorf("template %s page %s: %v", id, page.Slug, err)
			}
		}
	}
}

func TestTemplates_EachSeedsHomePage(t *testing.T) {
	// Every template must define a page with slug "home".
}
```

- [ ] **Step 2: Run to confirm failure**

- [ ] **Step 3: Author the five templates** — each seeds Home, Classes, Coaches and Contact with real boxing-gym placeholder copy in both languages, not lorem ipsum. Templates live in Go (so `apply-template` works server-side) and are mirrored in TypeScript for the picker's previews; the shared test in Task 3 Step 1 guards the mirror.

- [ ] **Step 4: Seed the first tenant** — the migration creates `site_settings` and the four pages for the org with slug `ibckirtipur` using the `fight_club` template, guarded by `WHERE EXISTS` so it is a no-op if that org is absent. `is_live` starts `false`.

- [ ] **Step 5: Verify and commit**

```bash
go test ./internal/site/ -v && make migrate-up
git add internal/site/ web/src/lib/site/templates.ts db/migrations/000049_*
git commit -m "feat: add five site templates and seed ibckirtipur"
```

---

### Task 7: Admin page builder

**Files:**
- Create: `web/src/app/[lang]/dashboard/site/page.tsx` (overview)
- Create: `web/src/app/[lang]/dashboard/site/pages/[id]/page.tsx` (builder)
- Create: `web/src/app/[lang]/dashboard/site/leads/page.tsx`
- Create: `web/src/components/site-builder/{BuilderCanvas,SectionList,SectionInspector,AddSectionPanel,TemplatePicker}.tsx`
- Create: `web/src/components/site-builder/useSiteBuilder.ts`
- Modify: `web/src/components/layout/Sidebar.tsx` (add the nav entry)
- Modify: `web/src/lib/i18n.ts`
- Test: `web/tests/site-builder.spec.ts`

**Interfaces:**
- Produces: `useSiteBuilder(pageId)` returning `{ sections, addSection, updateSection, removeSection, moveSection, saveState: 'idle'|'saving'|'saved'|'error' }`

- [ ] **Step 1: Write the failing tests** — add a section and see it in the preview; reorder and see the order persist after reload; edit a field and watch autosave reach the saved state; change a variant and see the preview change; publish and confirm the page appears on the public site; apply a template over existing pages only after confirming.

- [ ] **Step 2: Run to confirm failure**

- [ ] **Step 3: Implement** — three panes: section list (reorder, duplicate, delete), live preview using the real public renderers, and an inspector driven by `SECTION_SPECS`. Autosave debounced 1.5s with a visible state indicator. Reordering uses native HTML5 drag-and-drop **plus** keyboard-accessible move up/down buttons — drag-only is not acceptable.

- [ ] **Step 4: Verify and commit**

```bash
npx tsc --noEmit && npm run build && npx playwright test tests/site-builder.spec.ts
git add web/src
git commit -m "feat: add the admin page builder"
```

---

### Task 8: Remove the marketing pricing section

**Files:**
- Modify: `web/src/app/[lang]/page.tsx`
- Modify: `web/src/components/landing/LandingInteractive.tsx`
- Modify: `web/src/lib/i18n.ts`

- [ ] **Step 1: Write the failing test** — the landing page contains no pricing section, no `#pricing` anchor, and no nav or footer link to one, in both locales.

- [ ] **Step 2: Run to confirm failure**

- [ ] **Step 3: Implement** — delete the `pricingPlans` array, the `<section id="pricing">` block, the desktop nav link, the mobile nav link in `LandingInteractive.tsx`, and the footer link. Remove the now-orphaned `landing.pricing*` keys from BOTH dictionaries and drop the `Check` icon import if nothing else uses it. Gym membership packages are untouched — this is only the SaaS pricing block.

- [ ] **Step 4: Verify and commit**

```bash
grep -rn "pricing" web/src/ || echo "clean"
npx tsc --noEmit && npm run build
git add web/src
git commit -m "feat: remove the SaaS pricing section from the marketing landing page"
```

---

### Task 9: Member boxing profile

**Files:**
- Create: `db/migrations/000048_boxing_profile.up.sql` / `.down.sql`
- Create: `internal/boxing/{models,repository,service,handler,routes}.go`
- Create: `internal/boxing/service_test.go`
- Modify: `cmd/api/main.go`
- Create: `web/src/app/[lang]/member/boxing/page.tsx`
- Modify: `web/src/lib/api.ts`, `web/src/lib/i18n.ts`
- Test: `web/tests/boxing.spec.ts`

**Interfaces:**
- Produces: `WeightClassFor(weightKg float64) string` — standard amateur boxing divisions

- [ ] **Step 1: Write the migration** — `member_boxing_profiles` (stance, weight_class, skill_level, sparring_cleared, sparring_cleared_at, sparring_cleared_by, reach_cm, notes, UNIQUE(user_id, organization_id)) and `bout_records` (bout_date, opponent, event_name, result, method, rounds, weight_class, notes), both org-scoped with cascade deletes. Full DDL is in the design doc.

- [ ] **Step 2: Write the failing tests**

```go
func TestWeightClassFor(t *testing.T) {
	cases := []struct{ kg float64; want string }{
		{48, "light_flyweight"}, {52, "flyweight"}, {57, "featherweight"},
		{63, "light_welterweight"}, {69, "welterweight"}, {75, "middleweight"},
		{81, "light_heavyweight"}, {91, "heavyweight"}, {95, "super_heavyweight"},
	}
	// Boundaries matter — test exactly on each division limit.
}

func TestSparringClearance_MemberCannotClearThemselves(t *testing.T) {
	// A member calling the self route must not be able to set sparring_cleared.
}
```

- [ ] **Step 3: Run to confirm failure**

- [ ] **Step 4: Implement** — self routes `GET/PUT /members/me/boxing` and `GET/POST/DELETE /members/me/bouts`; staff routes under `/orgs/{orgId}/members/{memberId}/boxing` gated by `RequireOrgRole("staff","admin")`. The self-update handler must strip `sparring_cleared` from its accepted input entirely rather than validating it — a field the member cannot set should not be in their request struct. Weight class derives from the member's latest `health_metrics` weight, overridable.

- [ ] **Step 5: Verify and commit**

```bash
go test ./internal/boxing/ -v && go build ./... && cd web && npm run build
git add db/migrations/000048_* internal/boxing/ cmd/api/main.go web/src
git commit -m "feat: add member boxing profiles and bout records"
```

---

## Self-Review

**Spec coverage:** Design Phase 1 → Tasks 1–2; Phase 2 → Tasks 3–6; Phase 3 → Task 7; Phase 4 → Task 9; Phase 5 → Task 8. Every design section maps to a task.

**Ordering:** Task 1 (API) and Task 2 (routing) are independent of each other. Task 3 blocks Tasks 4, 6 and 7. Task 4 blocks Task 5. Task 6 needs Tasks 1 and 3. Task 7 needs Tasks 1, 3 and 4. Tasks 8 and 9 are fully independent and can run any time.

**Type consistency:** `Section` is defined once in Task 3 and mirrored in Go in Task 1 — the drift guard is the registry-comparison test in Task 3 Step 1. `SectionStyle` field names (`background`, `padding`, `width`, `align`) are used identically by `SectionShell` (Task 4) and the inspector (Task 7). `TemplateId` values (`fight_club`, `iron_sweat`, `champion`, `community`, `minimal_pro`) are used identically in Tasks 3, 6 and 7, and as the `site_settings.template` default in Task 1.

**Migration numbers:** 000045 and 000046 belong to the Phase 0 plan. This plan uses 000047 (site tables), 000048 (boxing), 000049 (ibckirtipur seed). Verify no collision before applying.

**Open risk:** the Go and TypeScript section registries are duplicated. The comparison test in Task 3 is the only thing keeping them honest; if it proves brittle, generating the TypeScript from the Go source at build time is the fallback.
