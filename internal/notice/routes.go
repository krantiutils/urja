package notice

import (
	"github.com/go-chi/chi/v5"
	"github.com/urja-gym/urja/pkg/middleware"
)

// RegisterOrgRoutes mounts org-scoped notice routes (under /orgs/{orgId}/notices).
// All members can list/get notices. Only staff/admin can create, update, delete.
func (h *Handler) RegisterOrgRoutes(r chi.Router) {
	r.Get("/", h.List)
	r.Get("/{id}", h.Get)

	// Staff/admin-only operations
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequireOrgRole("staff", "admin"))
		r.Post("/", h.Create)
		r.Put("/{id}", h.Update)
		r.Delete("/{id}", h.Delete)
	})
}
