package qrcode

import (
	"github.com/go-chi/chi/v5"

	"github.com/urja-gym/urja/pkg/middleware"
)

// RegisterOrgRoutes mounts org-scoped QR code routes (under /orgs/{orgId}).
//
// Staff only: this is the check-in code the gym prints and puts on the desk.
// Possession of it is what stands in for being present, so handing it out over
// the API to anybody who happens to be a member undermines the one check that
// scanning it represents. Members check themselves in through /members/me.
//
// Grouped for the same reason as the other handlers on this shared mux.
func (h *Handler) RegisterOrgRoutes(r chi.Router) {
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequireOrgRole("staff", "admin"))

		r.Get("/qr-code", h.GenerateQR)
	})
}
