package qrcode

import (
	"github.com/go-chi/chi/v5"
)

// RegisterOrgRoutes mounts org-scoped QR code routes (under /orgs/{orgId}).
func (h *Handler) RegisterOrgRoutes(r chi.Router) {
	r.Get("/qr-code", h.GenerateQR)
}
