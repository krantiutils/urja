package staff

import (
	"github.com/go-chi/chi/v5"
	"github.com/urja-gym/urja/pkg/middleware"
)

// RegisterOrgRoutes mounts org-scoped staff management routes (under /orgs/{orgId}/staff).
// All routes require staff or admin role.
func (h *Handler) RegisterOrgRoutes(r chi.Router) {
	r.Use(middleware.RequireOrgRole("staff", "admin"))

	r.Get("/", h.List)
	r.Post("/", h.Create)

	// Attendance sub-routes (must be before /{id} to avoid chi treating "attendance" as an id)
	r.Route("/attendance", func(r chi.Router) {
		r.Get("/", h.ListAttendance)
		r.Post("/", h.RecordAttendance)
	})

	r.Route("/{id}", func(r chi.Router) {
		r.Get("/", h.Get)
		r.Put("/", h.Update)
		r.Delete("/", h.Delete)
	})
}
