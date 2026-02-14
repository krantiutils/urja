package guide

import (
	"github.com/go-chi/chi/v5"
	"github.com/urja-gym/urja/pkg/middleware"
)

// RegisterPublicRoutes mounts public training guide routes (under /training-guides).
// No authentication required. Only published guides are visible.
func (h *Handler) RegisterPublicRoutes(r chi.Router) {
	r.Get("/", h.ListPublished)
	r.Get("/{id}", h.GetPublished)
}

// RegisterOrgRoutes mounts org-scoped training guide management routes (under /orgs/{orgId}/training-guides).
// All operations require staff or admin role.
func (h *Handler) RegisterOrgRoutes(r chi.Router) {
	r.Use(middleware.RequireRole("staff", "admin"))

	r.Get("/", h.ListAll)
	r.Get("/{id}", h.GetAdmin)
	r.Post("/", h.Create)
	r.Put("/{id}", h.Update)
	r.Patch("/{id}/publish", h.Publish)
	r.Delete("/{id}", h.Delete)
}
