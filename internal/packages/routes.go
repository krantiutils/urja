package packages

import (
	"github.com/go-chi/chi/v5"
	"github.com/urja-gym/urja/pkg/middleware"
)

// RegisterPublicRoutes mounts public package routes (under /packages).
func (h *Handler) RegisterPublicRoutes(r chi.Router) {
	r.Get("/", h.ListPublic)
}

// RegisterOrgRoutes mounts org-scoped package routes (under /orgs/{orgId}/packages).
// These require admin role.
func (h *Handler) RegisterOrgRoutes(r chi.Router) {
	r.Use(middleware.RequireOrgRole("admin"))
	r.Get("/", h.ListByOrg)
	r.Post("/", h.Create)
	r.Put("/{id}", h.Update)
	r.Delete("/{id}", h.Delete)
}

// RegisterSelfRoutes mounts member self-service package routes (under /members/me/packages).
func (h *Handler) RegisterSelfRoutes(r chi.Router) {
	r.Get("/", h.ListMyPackages)
	r.Post("/{id}/subscribe", h.Subscribe)
	r.Post("/verify-payment", h.VerifyPayment)
}
