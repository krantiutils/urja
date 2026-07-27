package site

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
)

// Service handles site business logic.
type Service struct {
	repo   *Repository
	logger *slog.Logger
}

// NewService creates a new site service.
func NewService(repo *Repository, logger *slog.Logger) *Service {
	return &Service{repo: repo, logger: logger}
}

// PublicSite is everything a tenant site needs to render its chrome.
type PublicSite struct {
	Slug      string          `json:"slug"`
	OrgName   string          `json:"org_name"`
	OrgNameNe string          `json:"org_name_ne,omitempty"`
	Template  string          `json:"template"`
	Theme     json.RawMessage `json:"theme"`
	Nav       json.RawMessage `json:"nav"`
	Footer    json.RawMessage `json:"footer"`
	Socials   json.RawMessage `json:"socials"`
	Pages     []PageSummary   `json:"pages"`
}

// GetPublicSite resolves a subdomain slug to a live site. A gym that has not
// flipped is_live is reported as not found, so an unfinished site never leaks.
func (s *Service) GetPublicSite(ctx context.Context, slug string) (*PublicSite, error) {
	slug = NormalizeSlug(slug)

	orgID, name, nameNe, err := s.repo.OrgBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}

	settings, err := s.repo.GetSettings(ctx, orgID)
	if err != nil {
		// A gym with no site row has no site.
		return nil, ErrNotFound
	}
	if !settings.IsLive {
		return nil, ErrNotFound
	}

	pages, err := s.repo.ListPages(ctx, orgID, true)
	if err != nil {
		return nil, err
	}

	return &PublicSite{
		Slug: slug, OrgName: name, OrgNameNe: nameNe,
		Template: settings.Template, Theme: settings.Theme, Nav: settings.Nav,
		Footer: settings.Footer, Socials: settings.Socials, Pages: pages,
	}, nil
}

// GetPublicPage returns one published page of a live site.
func (s *Service) GetPublicPage(ctx context.Context, slug, pageSlug string) (*Page, error) {
	slug = NormalizeSlug(slug)
	pageSlug = NormalizeSlug(pageSlug)
	if pageSlug == "" {
		pageSlug = "home"
	}

	orgID, _, _, err := s.repo.OrgBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}

	settings, err := s.repo.GetSettings(ctx, orgID)
	if err != nil || !settings.IsLive {
		return nil, ErrNotFound
	}

	return s.repo.GetPublishedPage(ctx, orgID, pageSlug)
}

// --- Admin operations ---

// CreatePageInput is the accepted shape for creating a page.
type CreatePageInput struct {
	Slug           string
	Title          string
	TitleNe        string
	Sections       []Section
	SEODescription string
	IsPublished    bool
	ShowInNav      bool
	SortOrder      int
}

// CreatePage validates and creates a page for an organization.
func (s *Service) CreatePage(ctx context.Context, orgID string, in CreatePageInput) (*Page, error) {
	slug := NormalizeSlug(in.Slug)
	if slug == "" {
		slug = NormalizeSlug(strings.ReplaceAll(in.Title, " ", "-"))
	}
	if err := ValidatePageSlug(slug); err != nil {
		return nil, err
	}

	title := strings.TrimSpace(in.Title)
	if title == "" {
		return nil, fmt.Errorf("title is required")
	}

	sections := in.Sections
	if sections == nil {
		sections = []Section{}
	}
	if err := ValidateSections(sections); err != nil {
		return nil, err
	}

	p, err := s.repo.CreatePage(ctx, orgID, &Page{
		Slug: slug, Title: title, TitleNe: strings.TrimSpace(in.TitleNe),
		Sections: sections, SEODescription: strings.TrimSpace(in.SEODescription),
		IsPublished: in.IsPublished, ShowInNav: in.ShowInNav, SortOrder: in.SortOrder,
	})
	if err != nil {
		return nil, err
	}

	s.logger.Info("site page created", "page_id", p.ID, "org_id", orgID, "slug", p.Slug)
	return p, nil
}

// UpdatePage validates and replaces a page's editable fields.
func (s *Service) UpdatePage(ctx context.Context, orgID, pageID string, in CreatePageInput) (*Page, error) {
	slug := NormalizeSlug(in.Slug)
	if err := ValidatePageSlug(slug); err != nil {
		return nil, err
	}

	title := strings.TrimSpace(in.Title)
	if title == "" {
		return nil, fmt.Errorf("title is required")
	}

	sections := in.Sections
	if sections == nil {
		sections = []Section{}
	}
	if err := ValidateSections(sections); err != nil {
		return nil, err
	}

	p, err := s.repo.UpdatePage(ctx, orgID, pageID, &Page{
		Slug: slug, Title: title, TitleNe: strings.TrimSpace(in.TitleNe),
		Sections: sections, SEODescription: strings.TrimSpace(in.SEODescription),
		IsPublished: in.IsPublished, ShowInNav: in.ShowInNav, SortOrder: in.SortOrder,
	})
	if err != nil {
		return nil, err
	}

	s.logger.Info("site page updated", "page_id", pageID, "org_id", orgID)
	return p, nil
}

// GetPage retrieves one page for editing, published or not.
func (s *Service) GetPage(ctx context.Context, orgID, pageID string) (*Page, error) {
	return s.repo.GetPageByID(ctx, orgID, pageID)
}

// ListPages returns every page of an organization, including drafts.
func (s *Service) ListPages(ctx context.Context, orgID string) ([]PageSummary, error) {
	return s.repo.ListPages(ctx, orgID, false)
}

// DeletePage removes a page.
func (s *Service) DeletePage(ctx context.Context, orgID, pageID string) error {
	if err := s.repo.DeletePage(ctx, orgID, pageID); err != nil {
		return err
	}
	s.logger.Info("site page deleted", "page_id", pageID, "org_id", orgID)
	return nil
}

// GetSettings retrieves a gym's site settings, materializing defaults on first
// access so the builder always has something to edit.
func (s *Service) GetSettings(ctx context.Context, orgID string) (*Settings, error) {
	settings, err := s.repo.GetSettings(ctx, orgID)
	if err == ErrNotFound {
		return s.repo.UpsertSettings(ctx, orgID, &Settings{
			Template: DefaultTemplate, IsLive: false,
		})
	}
	return settings, err
}

// UpdateSettingsInput is the accepted shape for a settings update.
type UpdateSettingsInput struct {
	Template *string
	Theme    *json.RawMessage
	Nav      *json.RawMessage
	Footer   *json.RawMessage
	Socials  *json.RawMessage
	IsLive   *bool
}

// UpdateSettings applies a partial settings update.
func (s *Service) UpdateSettings(ctx context.Context, orgID string, in UpdateSettingsInput) (*Settings, error) {
	current, err := s.GetSettings(ctx, orgID)
	if err != nil {
		return nil, err
	}

	if in.Template != nil {
		if _, ok := Templates[*in.Template]; !ok {
			return nil, fmt.Errorf("unknown template %q", *in.Template)
		}
		current.Template = *in.Template
	}
	if in.Theme != nil {
		current.Theme = *in.Theme
	}
	if in.Nav != nil {
		current.Nav = *in.Nav
	}
	if in.Footer != nil {
		current.Footer = *in.Footer
	}
	if in.Socials != nil {
		current.Socials = *in.Socials
	}
	if in.IsLive != nil {
		current.IsLive = *in.IsLive
	}

	updated, err := s.repo.UpsertSettings(ctx, orgID, current)
	if err != nil {
		return nil, err
	}

	s.logger.Info("site settings updated", "org_id", orgID,
		"template", updated.Template, "is_live", updated.IsLive)
	return updated, nil
}

// ApplyTemplate replaces a gym's entire site with a template's preset pages and
// theme. Destructive by design — the caller confirms before reaching here.
func (s *Service) ApplyTemplate(ctx context.Context, orgID, templateID string) error {
	tpl, ok := Templates[templateID]
	if !ok {
		return fmt.Errorf("unknown template %q", templateID)
	}

	// Validate before touching anything: a template with a bad section must not
	// be able to wipe an existing site and then fail.
	pages := tpl.Pages()
	for _, p := range pages {
		if err := ValidateSections(p.Sections); err != nil {
			return fmt.Errorf("template %q page %q is invalid: %w", templateID, p.Slug, err)
		}
	}

	if err := s.repo.ReplaceAllPages(ctx, orgID, pages); err != nil {
		return err
	}

	current, err := s.GetSettings(ctx, orgID)
	if err != nil {
		return err
	}
	current.Template = templateID
	current.Theme = tpl.ThemeJSON()
	if _, err := s.repo.UpsertSettings(ctx, orgID, current); err != nil {
		return err
	}

	s.logger.Info("site template applied", "org_id", orgID, "template", templateID, "pages", len(pages))
	return nil
}

// --- Leads ---

// CreateLeadInput is the accepted shape for a public enquiry.
type CreateLeadInput struct {
	Name       string
	Phone      string
	Email      string
	Message    string
	Interest   string
	SourcePage string
	// Honeypot must be empty. Bots fill every field they find.
	Honeypot string
}

// leadCooldownMinutes bounds how often one phone number may submit to one gym.
const leadCooldownMinutes = 10

// maxLeadsPerPhonePerWindow is the submission cap within the cooldown window.
const maxLeadsPerPhonePerWindow = 3

// SubmitLead validates and records a public enquiry.
func (s *Service) SubmitLead(ctx context.Context, siteSlug string, in CreateLeadInput) (*Lead, error) {
	// A filled honeypot is a bot. Report success so it does not learn otherwise,
	// but persist nothing.
	if strings.TrimSpace(in.Honeypot) != "" {
		s.logger.Info("lead honeypot triggered", "site", siteSlug)
		return &Lead{Name: in.Name, Status: "new"}, nil
	}

	orgID, _, _, err := s.repo.OrgBySlug(ctx, NormalizeSlug(siteSlug))
	if err != nil {
		return nil, err
	}

	name := strings.TrimSpace(in.Name)
	phone := strings.TrimSpace(in.Phone)
	if name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if phone == "" {
		return nil, fmt.Errorf("phone is required")
	}
	if len(name) > 200 || len(phone) > 30 || len(in.Message) > 2000 {
		return nil, fmt.Errorf("submission is too long")
	}

	recent, err := s.repo.CountRecentLeadsByPhone(ctx, orgID, phone, leadCooldownMinutes)
	if err != nil {
		return nil, err
	}
	if recent >= maxLeadsPerPhonePerWindow {
		return nil, fmt.Errorf("too many submissions, please try again later")
	}

	lead, err := s.repo.CreateLead(ctx, orgID, &Lead{
		Name: name, Phone: phone,
		Email:    strings.TrimSpace(in.Email),
		Message:  strings.TrimSpace(in.Message),
		Interest: strings.TrimSpace(in.Interest),
		// The source page is a slug from our own site, not free text.
		SourcePage: NormalizeSlug(in.SourcePage),
	})
	if err != nil {
		return nil, err
	}

	s.logger.Info("site lead received", "lead_id", lead.ID, "org_id", orgID, "interest", lead.Interest)
	return lead, nil
}

// ListLeads returns an organization's enquiries.
func (s *Service) ListLeads(ctx context.Context, orgID, status string, limit, offset int) ([]Lead, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}
	if status != "" {
		if err := ValidateLeadStatus(status); err != nil {
			return nil, err
		}
	}
	return s.repo.ListLeads(ctx, orgID, status, limit, offset)
}

// UpdateLeadStatus moves a lead along the pipeline.
func (s *Service) UpdateLeadStatus(ctx context.Context, orgID, leadID, status string) error {
	if err := ValidateLeadStatus(status); err != nil {
		return err
	}
	return s.repo.UpdateLeadStatus(ctx, orgID, leadID, status)
}
