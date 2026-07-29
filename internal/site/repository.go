package site

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ErrNotFound is returned when a page, site or lead does not exist, or exists
// but is not visible to the caller. Public handlers translate it to a 404 so an
// unpublished page is indistinguishable from a missing one.
var ErrNotFound = errors.New("not found")

// ErrSlugTaken is returned when a page slug collides within an organization.
var ErrSlugTaken = errors.New("a page with that slug already exists")

// Repository handles site data persistence. Every org-scoped query filters on
// organization_id — there is no method that reads or writes a page without one.
type Repository struct {
	db *pgxpool.Pool
}

// NewRepository creates a new site repository.
func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

const pageColumns = `id, organization_id, slug, title, COALESCE(title_ne, ''), sections,
	COALESCE(seo_description, ''), is_published, show_in_nav, sort_order, created_at, updated_at`

func scanPage(row pgx.Row) (*Page, error) {
	var p Page
	var sections []byte
	err := row.Scan(&p.ID, &p.OrgID, &p.Slug, &p.Title, &p.TitleNe, &sections,
		&p.SEODescription, &p.IsPublished, &p.ShowInNav, &p.SortOrder, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(sections, &p.Sections); err != nil {
		return nil, fmt.Errorf("decoding sections for page %s: %w", p.ID, err)
	}
	if p.Sections == nil {
		p.Sections = []Section{}
	}
	return &p, nil
}

// --- Organization lookup ---

// OrgBySlug resolves a subdomain label to an organization ID and display name.
// It is the entry point for every public request.
func (r *Repository) OrgBySlug(ctx context.Context, slug string) (orgID, name, nameNe string, err error) {
	err = r.db.QueryRow(ctx,
		`SELECT id, name, COALESCE(name_ne, '') FROM organizations
		 WHERE slug = $1 AND is_active = true`, slug,
	).Scan(&orgID, &name, &nameNe)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", "", "", ErrNotFound
	}
	if err != nil {
		return "", "", "", fmt.Errorf("resolving org slug: %w", err)
	}
	return orgID, name, nameNe, nil
}

// --- Pages ---

// CreatePage inserts a new page for an organization.
func (r *Repository) CreatePage(ctx context.Context, orgID string, p *Page) (*Page, error) {
	sections, err := json.Marshal(p.Sections)
	if err != nil {
		return nil, fmt.Errorf("encoding sections: %w", err)
	}

	row := r.db.QueryRow(ctx,
		`INSERT INTO site_pages
			(organization_id, slug, title, title_ne, sections, seo_description,
			 is_published, show_in_nav, sort_order)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		 RETURNING `+pageColumns,
		orgID, p.Slug, p.Title, nilIfEmpty(p.TitleNe), sections, nilIfEmpty(p.SEODescription),
		p.IsPublished, p.ShowInNav, p.SortOrder,
	)

	created, err := scanPage(row)
	if err != nil {
		if isUniqueViolation(err) {
			return nil, ErrSlugTaken
		}
		return nil, fmt.Errorf("creating page: %w", err)
	}
	return created, nil
}

// GetPageByID retrieves one page within an organization, published or not.
func (r *Repository) GetPageByID(ctx context.Context, orgID, pageID string) (*Page, error) {
	p, err := scanPage(r.db.QueryRow(ctx,
		`SELECT `+pageColumns+` FROM site_pages WHERE id = $1 AND organization_id = $2`,
		pageID, orgID,
	))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("getting page: %w", err)
	}
	return p, nil
}

// GetPublishedPage retrieves a published page by slug. Used by the public site.
func (r *Repository) GetPublishedPage(ctx context.Context, orgID, slug string) (*Page, error) {
	p, err := scanPage(r.db.QueryRow(ctx,
		`SELECT `+pageColumns+` FROM site_pages
		 WHERE organization_id = $1 AND slug = $2 AND is_published = true`,
		orgID, slug,
	))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("getting published page: %w", err)
	}
	return p, nil
}

// ListPages returns page summaries for an organization. When publishedOnly is
// true only published pages are returned, for public nav.
func (r *Repository) ListPages(ctx context.Context, orgID string, publishedOnly bool) ([]PageSummary, error) {
	query := `SELECT id, slug, title, COALESCE(title_ne, ''), is_published, show_in_nav, sort_order
	          FROM site_pages WHERE organization_id = $1`
	if publishedOnly {
		query += ` AND is_published = true`
	}
	query += ` ORDER BY sort_order, title`

	rows, err := r.db.Query(ctx, query, orgID)
	if err != nil {
		return nil, fmt.Errorf("listing pages: %w", err)
	}
	defer rows.Close()

	pages := []PageSummary{}
	for rows.Next() {
		var s PageSummary
		if err := rows.Scan(&s.ID, &s.Slug, &s.Title, &s.TitleNe,
			&s.IsPublished, &s.ShowInNav, &s.SortOrder); err != nil {
			return nil, fmt.Errorf("scanning page summary: %w", err)
		}
		pages = append(pages, s)
	}
	return pages, rows.Err()
}

// UpdatePage replaces a page's editable fields.
func (r *Repository) UpdatePage(ctx context.Context, orgID, pageID string, p *Page) (*Page, error) {
	sections, err := json.Marshal(p.Sections)
	if err != nil {
		return nil, fmt.Errorf("encoding sections: %w", err)
	}

	row := r.db.QueryRow(ctx,
		`UPDATE site_pages
		 SET slug = $3, title = $4, title_ne = $5, sections = $6, seo_description = $7,
		     is_published = $8, show_in_nav = $9, sort_order = $10, updated_at = NOW()
		 WHERE id = $1 AND organization_id = $2
		 RETURNING `+pageColumns,
		pageID, orgID, p.Slug, p.Title, nilIfEmpty(p.TitleNe), sections,
		nilIfEmpty(p.SEODescription), p.IsPublished, p.ShowInNav, p.SortOrder,
	)

	updated, err := scanPage(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		if isUniqueViolation(err) {
			return nil, ErrSlugTaken
		}
		return nil, fmt.Errorf("updating page: %w", err)
	}
	return updated, nil
}

// DeletePage removes a page from an organization.
func (r *Repository) DeletePage(ctx context.Context, orgID, pageID string) error {
	tag, err := r.db.Exec(ctx,
		`DELETE FROM site_pages WHERE id = $1 AND organization_id = $2`, pageID, orgID)
	if err != nil {
		return fmt.Errorf("deleting page: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// ReplaceAllPages swaps an organization's entire page set in one transaction.
// Used by apply-template, which is destructive by design.
func (r *Repository) ReplaceAllPages(ctx context.Context, orgID string, pages []Page) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("beginning transaction: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	if _, err := tx.Exec(ctx, `DELETE FROM site_pages WHERE organization_id = $1`, orgID); err != nil {
		return fmt.Errorf("clearing pages: %w", err)
	}

	for _, p := range pages {
		sections, err := json.Marshal(p.Sections)
		if err != nil {
			return fmt.Errorf("encoding sections for %s: %w", p.Slug, err)
		}
		_, err = tx.Exec(ctx,
			`INSERT INTO site_pages
				(organization_id, slug, title, title_ne, sections, is_published, show_in_nav, sort_order)
			 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
			orgID, p.Slug, p.Title, nilIfEmpty(p.TitleNe), sections,
			p.IsPublished, p.ShowInNav, p.SortOrder,
		)
		if err != nil {
			return fmt.Errorf("inserting page %s: %w", p.Slug, err)
		}
	}

	return tx.Commit(ctx)
}

// --- Settings ---

// ListLiveSites returns every gym whose site is published, for the apex
// sitemap index. Scoped to live sites of active gyms: a draft site is not
// discoverable by guessing its subdomain, and listing it here would undo that.
func (r *Repository) ListLiveSites(ctx context.Context) ([]LiveSite, error) {
	rows, err := r.db.Query(ctx,
		`SELECT o.slug, ss.updated_at
		   FROM site_settings ss
		   JOIN organizations o ON o.id = ss.organization_id
		  WHERE ss.is_live = true AND o.is_active = true
		  ORDER BY o.slug`)
	if err != nil {
		return nil, fmt.Errorf("listing live sites: %w", err)
	}
	defer rows.Close()

	out := []LiveSite{}
	for rows.Next() {
		var s LiveSite
		if err := rows.Scan(&s.Slug, &s.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scanning live site: %w", err)
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

// GetSettings retrieves a gym's site settings, returning ErrNotFound if the
// gym has never configured a site.
func (r *Repository) GetSettings(ctx context.Context, orgID string) (*Settings, error) {
	var s Settings
	err := r.db.QueryRow(ctx,
		`SELECT ss.organization_id, o.slug, ss.template, ss.theme, ss.nav, ss.socials,
		        ss.footer, ss.is_live, ss.updated_at
		 FROM site_settings ss
		 JOIN organizations o ON o.id = ss.organization_id
		 WHERE ss.organization_id = $1`, orgID,
	).Scan(&s.OrgID, &s.Slug, &s.Template, &s.Theme, &s.Nav, &s.Socials, &s.Footer, &s.IsLive, &s.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("getting site settings: %w", err)
	}
	return &s, nil
}

// UpsertSettings creates or replaces a gym's site settings.
func (r *Repository) UpsertSettings(ctx context.Context, orgID string, s *Settings) (*Settings, error) {
	var out Settings
	err := r.db.QueryRow(ctx,
		`INSERT INTO site_settings (organization_id, template, theme, nav, footer, socials, is_live)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)
		 ON CONFLICT (organization_id) DO UPDATE
		 SET template = EXCLUDED.template, theme = EXCLUDED.theme, nav = EXCLUDED.nav,
		     footer = EXCLUDED.footer, socials = EXCLUDED.socials, is_live = EXCLUDED.is_live,
		     updated_at = NOW()
		 RETURNING organization_id, template, theme, nav, footer, socials, is_live, updated_at`,
		orgID, s.Template, jsonOrEmpty(s.Theme), jsonOrEmpty(s.Nav),
		jsonOrEmpty(s.Footer), jsonOrEmpty(s.Socials), s.IsLive,
	).Scan(&out.OrgID, &out.Template, &out.Theme, &out.Nav, &out.Footer,
		&out.Socials, &out.IsLive, &out.UpdatedAt)
	if err == nil {
		// RETURNING cannot reach the joined organizations row, and a settings
		// response without the slug leaves the builder unable to tell an owner
		// where their site is published.
		err = r.db.QueryRow(ctx, `SELECT slug FROM organizations WHERE id = $1`, orgID).
			Scan(&out.Slug)
	}
	if err != nil {
		return nil, fmt.Errorf("saving site settings: %w", err)
	}
	return &out, nil
}

// --- Leads ---

// CreateLead records a public enquiry.
func (r *Repository) CreateLead(ctx context.Context, orgID string, l *Lead) (*Lead, error) {
	var out Lead
	err := r.db.QueryRow(ctx,
		`INSERT INTO site_leads (organization_id, name, phone, email, message, interest, source_page)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)
		 RETURNING id, organization_id, name, phone, COALESCE(email, ''), COALESCE(message, ''),
		           COALESCE(interest, ''), COALESCE(source_page, ''), status, created_at`,
		orgID, l.Name, l.Phone, nilIfEmpty(l.Email), nilIfEmpty(l.Message),
		nilIfEmpty(l.Interest), nilIfEmpty(l.SourcePage),
	).Scan(&out.ID, &out.OrgID, &out.Name, &out.Phone, &out.Email, &out.Message,
		&out.Interest, &out.SourcePage, &out.Status, &out.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("creating lead: %w", err)
	}
	return &out, nil
}

// ListLeads returns an organization's enquiries, newest first.
func (r *Repository) ListLeads(ctx context.Context, orgID, status string, limit, offset int) ([]Lead, error) {
	query := `SELECT id, organization_id, name, phone, COALESCE(email, ''), COALESCE(message, ''),
	                 COALESCE(interest, ''), COALESCE(source_page, ''), status, created_at
	          FROM site_leads WHERE organization_id = $1`
	args := []interface{}{orgID}

	if status != "" {
		query += ` AND status = $2 ORDER BY created_at DESC LIMIT $3 OFFSET $4`
		args = append(args, status, limit, offset)
	} else {
		query += ` ORDER BY created_at DESC LIMIT $2 OFFSET $3`
		args = append(args, limit, offset)
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("listing leads: %w", err)
	}
	defer rows.Close()

	leads := []Lead{}
	for rows.Next() {
		var l Lead
		if err := rows.Scan(&l.ID, &l.OrgID, &l.Name, &l.Phone, &l.Email, &l.Message,
			&l.Interest, &l.SourcePage, &l.Status, &l.CreatedAt); err != nil {
			return nil, fmt.Errorf("scanning lead: %w", err)
		}
		leads = append(leads, l)
	}
	return leads, rows.Err()
}

// UpdateLeadStatus moves a lead along the pipeline.
func (r *Repository) UpdateLeadStatus(ctx context.Context, orgID, leadID, status string) error {
	tag, err := r.db.Exec(ctx,
		`UPDATE site_leads SET status = $3 WHERE id = $1 AND organization_id = $2`,
		leadID, orgID, status)
	if err != nil {
		return fmt.Errorf("updating lead status: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// CountRecentLeadsByPhone supports abuse throttling on the public form.
func (r *Repository) CountRecentLeadsByPhone(ctx context.Context, orgID, phone string, withinMinutes int) (int, error) {
	var n int
	err := r.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM site_leads
		 WHERE organization_id = $1 AND phone = $2
		   AND created_at > NOW() - make_interval(mins => $3)`,
		orgID, phone, withinMinutes,
	).Scan(&n)
	if err != nil {
		return 0, fmt.Errorf("counting recent leads: %w", err)
	}
	return n, nil
}

// --- helpers ---

func nilIfEmpty(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}

// jsonOrEmpty guards against writing a NULL or invalid value into a NOT NULL
// JSONB column when the caller left a field unset.
func jsonOrEmpty(raw json.RawMessage) []byte {
	if len(raw) == 0 {
		return []byte(`{}`)
	}
	return raw
}

// isUniqueViolation reports whether an error is a Postgres unique constraint
// violation (SQLSTATE 23505).
func isUniqueViolation(err error) bool {
	var pgErr interface{ SQLState() string }
	if errors.As(err, &pgErr) {
		return pgErr.SQLState() == "23505"
	}
	return false
}
