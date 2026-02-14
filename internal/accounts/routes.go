package accounts

import (
	"github.com/go-chi/chi/v5"
	"github.com/urja-gym/urja/pkg/middleware"
)

// RegisterRoutes mounts accounts routes (under /orgs/{orgId}/accounts).
// All accounts routes require admin or staff role.
func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Use(middleware.RequireRole("admin", "staff"))

	r.Get("/", h.List)
	r.Post("/", h.Create)
	r.Get("/summary", h.GetSummary)
	r.Get("/export", h.Export)
	r.Put("/{id}", h.Update)
	r.Delete("/{id}", h.Delete)
}
