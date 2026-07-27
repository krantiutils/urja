package site

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/urja-gym/urja/pkg/middleware"
)

// maxBodyBytes bounds a request body. Section arrays are the largest thing the
// admin API accepts, and the public lead form should never be large at all.
const (
	maxAdminBodyBytes  = 2 << 20  // 2 MiB
	maxPublicBodyBytes = 16 << 10 // 16 KiB
)

// Handler handles HTTP requests for site endpoints.
type Handler struct {
	service *Service
	logger  *slog.Logger
}

// NewHandler creates a new site handler.
func NewHandler(service *Service, logger *slog.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

// --- Public endpoints ---

// GetPublicSite handles GET /api/v1/sites/{slug}
func (h *Handler) GetPublicSite(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	if slug == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing site slug"})
		return
	}

	site, err := h.service.GetPublicSite(r.Context(), slug)
	if err != nil {
		// A gym that exists but is not live is reported identically to one that
		// does not exist, so an unfinished site cannot be discovered.
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "site not found"})
		return
	}

	writeJSON(w, http.StatusOK, site)
}

// GetPublicPage handles GET /api/v1/sites/{slug}/pages/{pageSlug}
func (h *Handler) GetPublicPage(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	pageSlug := chi.URLParam(r, "pageSlug")

	page, err := h.service.GetPublicPage(r.Context(), slug, pageSlug)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "page not found"})
		return
	}

	writeJSON(w, http.StatusOK, page)
}

type submitLeadRequest struct {
	Name       string `json:"name"`
	Phone      string `json:"phone"`
	Email      string `json:"email"`
	Message    string `json:"message"`
	Interest   string `json:"interest"`
	SourcePage string `json:"source_page"`
	// Website is a honeypot: a real person never sees or fills this field.
	Website string `json:"website"`
}

// SubmitLead handles POST /api/v1/sites/{slug}/leads
func (h *Handler) SubmitLead(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	r.Body = http.MaxBytesReader(w, r.Body, maxPublicBodyBytes)
	var req submitLeadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	lead, err := h.service.SubmitLead(r.Context(), slug, CreateLeadInput{
		Name: req.Name, Phone: req.Phone, Email: req.Email, Message: req.Message,
		Interest: req.Interest, SourcePage: req.SourcePage, Honeypot: req.Website,
	})
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "site not found"})
			return
		}
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	// Never echo the stored lead back to an anonymous caller — it would let
	// anyone confirm what was recorded.
	writeJSON(w, http.StatusCreated, map[string]string{
		"status": "received", "name": lead.Name,
	})
}

// ListTemplates handles GET /api/v1/sites/templates — public so the builder's
// picker can render previews without an extra auth round trip.
func (h *Handler) ListTemplates(w http.ResponseWriter, r *http.Request) {
	out := make([]map[string]interface{}, 0, len(Templates))
	for _, id := range TemplateIDs() {
		tpl := Templates[id]
		out = append(out, map[string]interface{}{
			"id": tpl.ID, "name": tpl.Name, "name_ne": tpl.NameNe, "theme": tpl.Theme,
		})
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"data": out})
}

// --- Admin endpoints ---

// GetSettings handles GET /api/v1/orgs/{orgId}/site/settings
func (h *Handler) GetSettings(w http.ResponseWriter, r *http.Request) {
	orgID, ok := middleware.OrgIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing organization"})
		return
	}

	settings, err := h.service.GetSettings(r.Context(), orgID)
	if err != nil {
		h.logger.Error("failed to get site settings", "error", err, "org_id", orgID)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		return
	}

	writeJSON(w, http.StatusOK, settings)
}

type updateSettingsRequest struct {
	Template *string          `json:"template"`
	Theme    *json.RawMessage `json:"theme"`
	Nav      *json.RawMessage `json:"nav"`
	Footer   *json.RawMessage `json:"footer"`
	Socials  *json.RawMessage `json:"socials"`
	IsLive   *bool            `json:"is_live"`
}

// UpdateSettings handles PUT /api/v1/orgs/{orgId}/site/settings
func (h *Handler) UpdateSettings(w http.ResponseWriter, r *http.Request) {
	orgID, ok := middleware.OrgIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing organization"})
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxAdminBodyBytes)
	var req updateSettingsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	settings, err := h.service.UpdateSettings(r.Context(), orgID, UpdateSettingsInput{
		Template: req.Template, Theme: req.Theme, Nav: req.Nav,
		Footer: req.Footer, Socials: req.Socials, IsLive: req.IsLive,
	})
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, settings)
}

// ListPages handles GET /api/v1/orgs/{orgId}/site/pages
func (h *Handler) ListPages(w http.ResponseWriter, r *http.Request) {
	orgID, ok := middleware.OrgIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing organization"})
		return
	}

	pages, err := h.service.ListPages(r.Context(), orgID)
	if err != nil {
		h.logger.Error("failed to list site pages", "error", err, "org_id", orgID)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{"data": pages})
}

type pageRequest struct {
	Slug           string    `json:"slug"`
	Title          string    `json:"title"`
	TitleNe        string    `json:"title_ne"`
	Sections       []Section `json:"sections"`
	SEODescription string    `json:"seo_description"`
	IsPublished    bool      `json:"is_published"`
	ShowInNav      bool      `json:"show_in_nav"`
	SortOrder      int       `json:"sort_order"`
}

func (p pageRequest) toInput() CreatePageInput {
	return CreatePageInput{
		Slug: p.Slug, Title: p.Title, TitleNe: p.TitleNe, Sections: p.Sections,
		SEODescription: p.SEODescription, IsPublished: p.IsPublished,
		ShowInNav: p.ShowInNav, SortOrder: p.SortOrder,
	}
}

// CreatePage handles POST /api/v1/orgs/{orgId}/site/pages
func (h *Handler) CreatePage(w http.ResponseWriter, r *http.Request) {
	orgID, ok := middleware.OrgIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing organization"})
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxAdminBodyBytes)
	var req pageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	page, err := h.service.CreatePage(r.Context(), orgID, req.toInput())
	if err != nil {
		if errors.Is(err, ErrSlugTaken) {
			writeJSON(w, http.StatusConflict, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusCreated, page)
}

// GetPage handles GET /api/v1/orgs/{orgId}/site/pages/{pageId}
func (h *Handler) GetPage(w http.ResponseWriter, r *http.Request) {
	orgID, ok := middleware.OrgIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing organization"})
		return
	}

	page, err := h.service.GetPage(r.Context(), orgID, chi.URLParam(r, "pageId"))
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "page not found"})
		return
	}

	writeJSON(w, http.StatusOK, page)
}

// UpdatePage handles PUT /api/v1/orgs/{orgId}/site/pages/{pageId}
func (h *Handler) UpdatePage(w http.ResponseWriter, r *http.Request) {
	orgID, ok := middleware.OrgIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing organization"})
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxAdminBodyBytes)
	var req pageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	page, err := h.service.UpdatePage(r.Context(), orgID, chi.URLParam(r, "pageId"), req.toInput())
	if err != nil {
		switch {
		case errors.Is(err, ErrNotFound):
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "page not found"})
		case errors.Is(err, ErrSlugTaken):
			writeJSON(w, http.StatusConflict, map[string]string{"error": err.Error()})
		default:
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
		return
	}

	writeJSON(w, http.StatusOK, page)
}

// DeletePage handles DELETE /api/v1/orgs/{orgId}/site/pages/{pageId}
func (h *Handler) DeletePage(w http.ResponseWriter, r *http.Request) {
	orgID, ok := middleware.OrgIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing organization"})
		return
	}

	if err := h.service.DeletePage(r.Context(), orgID, chi.URLParam(r, "pageId")); err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "page not found"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "page deleted"})
}

type applyTemplateRequest struct {
	Template string `json:"template"`
	// Confirm must be true. Applying a template destroys every existing page,
	// so the client has to say so explicitly rather than by omission.
	Confirm bool `json:"confirm"`
}

// ApplyTemplate handles POST /api/v1/orgs/{orgId}/site/apply-template
func (h *Handler) ApplyTemplate(w http.ResponseWriter, r *http.Request) {
	orgID, ok := middleware.OrgIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing organization"})
		return
	}

	var req applyTemplateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	if !req.Confirm {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "applying a template replaces every page; set confirm to true",
		})
		return
	}

	if err := h.service.ApplyTemplate(r.Context(), orgID, req.Template); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	pages, err := h.service.ListPages(r.Context(), orgID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{"data": pages})
}

// ListLeads handles GET /api/v1/orgs/{orgId}/site/leads
func (h *Handler) ListLeads(w http.ResponseWriter, r *http.Request) {
	orgID, ok := middleware.OrgIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing organization"})
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	leads, err := h.service.ListLeads(r.Context(), orgID, r.URL.Query().Get("status"), limit, offset)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{"data": leads})
}

type updateLeadRequest struct {
	Status string `json:"status"`
}

// UpdateLead handles PATCH /api/v1/orgs/{orgId}/site/leads/{leadId}
func (h *Handler) UpdateLead(w http.ResponseWriter, r *http.Request) {
	orgID, ok := middleware.OrgIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing organization"})
		return
	}

	var req updateLeadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	if err := h.service.UpdateLeadStatus(r.Context(), orgID, chi.URLParam(r, "leadId"), req.Status); err != nil {
		if errors.Is(err, ErrNotFound) {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "lead not found"})
			return
		}
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "lead updated"})
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.Error("failed to encode JSON response", "error", err)
	}
}
