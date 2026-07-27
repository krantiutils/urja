package nfc

import (
	"github.com/go-chi/chi/v5"
	"github.com/urja-gym/urja/pkg/middleware"
)

// RegisterOrgRoutes mounts org-scoped NFC routes (under /orgs/{orgId}).
// All NFC management (card and device registration/assignment) requires staff or admin role.
//
// This router is the shared /orgs/{orgId} mux (see cmd/api/main.go), which already
// has direct routes registered on it (e.g. org update, other RegisterOrgRoutes
// calls) by the time this runs — chi forbids calling Use() on a mux that already
// has routes registered, so the middleware must be scoped via Group() rather than
// applied at the top level here.
func (h *Handler) RegisterOrgRoutes(r chi.Router) {
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequireOrgRole("staff", "admin"))

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
	})
}

// RegisterDeviceRoutes mounts device-level routes (no JWT auth, uses X-Device-Key).
// These are mounted outside the authenticated route group.
func (h *Handler) RegisterDeviceRoutes(r chi.Router) {
	r.Post("/check-in", h.DeviceCheckIn)
	r.Post("/heartbeat", h.DeviceHeartbeat)
}
