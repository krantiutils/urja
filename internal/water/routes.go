package water

import (
	"github.com/go-chi/chi/v5"
)

// RegisterSelfRoutes mounts member self-service water tracking routes (under /members/me).
func (h *Handler) RegisterSelfRoutes(r chi.Router) {
	r.Post("/water-logs", h.CreateWaterLog)
	r.Get("/water-logs", h.ListWaterLogs)
	r.Delete("/water-logs/{id}", h.DeleteWaterLog)
	r.Get("/water-logs/summary", h.GetDailySummary)
}
