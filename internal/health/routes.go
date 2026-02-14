package health

import (
	"github.com/go-chi/chi/v5"
)

// RegisterSelfRoutes mounts health self-service routes (under /members/me/health).
// All routes require authentication (handled by parent router group).
func (h *Handler) RegisterSelfRoutes(r chi.Router) {
	r.Get("/", h.GetMetrics)
	r.Post("/bmi", h.LogBMI)
	r.Post("/measurements", h.LogMeasurements)
	r.Route("/photos", func(r chi.Router) {
		r.Post("/", h.UploadPhoto)
		r.Get("/", h.ListPhotos)
	})
}
