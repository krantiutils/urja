package nfc

import (
	"github.com/go-chi/chi/v5"
)

// RegisterOrgRoutes mounts org-scoped NFC routes (under /orgs/{orgId}).
func (h *Handler) RegisterOrgRoutes(r chi.Router) {
	r.Route("/nfc-cards", func(r chi.Router) {
		r.Get("/", h.ListCards)
		r.Post("/", h.RegisterCard)
		r.Put("/{id}/assign", h.AssignCard)
		r.Put("/{id}/unassign", h.UnassignCard)
	})

	r.Get("/nfc-devices", h.ListDevices)
}
