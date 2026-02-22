package exercise

import (
	"github.com/go-chi/chi/v5"
)

// RegisterSelfRoutes mounts member self-service exercise and program routes (under /members/me).
func (h *Handler) RegisterSelfRoutes(r chi.Router) {
	r.Get("/exercises", h.ListExercises)
	r.Get("/exercises/{id}", h.GetExercise)
	r.Get("/programs", h.ListPrograms)
	r.Get("/programs/{id}", h.GetProgram)
	r.Post("/programs/{id}/enroll", h.EnrollInProgram)
	r.Get("/programs/current", h.GetCurrentProgram)
	r.Post("/programs/current/complete-day", h.CompleteDay)
}
