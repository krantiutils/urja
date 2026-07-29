package site

import (
	"github.com/go-chi/chi/v5"

	"github.com/urja-gym/urja/pkg/middleware"
)

// RegisterPublicRoutes mounts unauthenticated site routes (under /sites).
// These serve the tenant subdomains, so they must never expose draft content:
// the service reports an unpublished page or a not-live site as not found.
func (h *Handler) RegisterPublicRoutes(r chi.Router) {
	r.Get("/", h.ListLiveSites)
	r.Get("/templates", h.ListTemplates)
	r.Get("/{slug}", h.GetPublicSite)
	r.Get("/{slug}/pages/{pageSlug}", h.GetPublicPage)
	r.Post("/{slug}/leads", h.SubmitLead)
}

// RegisterOrgRoutes mounts site management routes (under /orgs/{orgId}/site).
// Editing a gym's public website is an owner-level action, so admin only —
// staff who can manage members should not be able to republish the website.
func (h *Handler) RegisterOrgRoutes(r chi.Router) {
	r.Use(middleware.RequireOrgRole("admin"))

	r.Get("/settings", h.GetSettings)
	r.Put("/settings", h.UpdateSettings)

	r.Post("/apply-template", h.ApplyTemplate)

	r.Route("/pages", func(r chi.Router) {
		r.Get("/", h.ListPages)
		r.Post("/", h.CreatePage)
		r.Get("/{pageId}", h.GetPage)
		r.Put("/{pageId}", h.UpdatePage)
		r.Delete("/{pageId}", h.DeletePage)
	})

	r.Route("/leads", func(r chi.Router) {
		r.Get("/", h.ListLeads)
		r.Patch("/{leadId}", h.UpdateLead)
	})
}
