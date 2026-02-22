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

	r.Route("/nfc-devices", func(r chi.Router) {
		r.Get("/", h.ListDevices)
		r.Post("/", h.RegisterDevice)
	})
}

// RegisterDeviceRoutes mounts device-level routes (no JWT auth, uses X-Device-Key).
// These are mounted outside the authenticated route group.
func (h *Handler) RegisterDeviceRoutes(r chi.Router) {
	r.Post("/check-in", h.DeviceCheckIn)
	r.Post("/heartbeat", h.DeviceHeartbeat)
}
