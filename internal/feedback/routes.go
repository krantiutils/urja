package feedback

import (
	"github.com/go-chi/chi/v5"
	"github.com/urja-gym/urja/pkg/middleware"
)

// RegisterOrgRoutes mounts org-scoped feedback routes (under /orgs/{orgId}/feedbacks).
// Any org member can submit feedback. Only staff/admin can list all feedbacks.
func (h *Handler) RegisterOrgRoutes(r chi.Router) {
	r.Post("/", h.Create)

	r.Group(func(r chi.Router) {
		r.Use(middleware.RequireOrgRole("staff", "admin"))
		r.Get("/", h.List)
	})
}
