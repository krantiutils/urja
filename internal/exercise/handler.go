package exercise

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/urja-gym/urja/pkg/middleware"
)

// Handler handles HTTP requests for exercise endpoints.
type Handler struct {
	service *Service
	logger  *slog.Logger
}

// NewHandler creates a new exercise handler.
func NewHandler(service *Service, logger *slog.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

// ListExercises handles GET /exercises
func (h *Handler) ListExercises(w http.ResponseWriter, r *http.Request) {
	exerciseType := r.URL.Query().Get("type")
	equipment := r.URL.Query().Get("equipment")
	muscleGroup := r.URL.Query().Get("muscle_group")
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	exercises, err := h.service.ListExercises(r.Context(), exerciseType, equipment, muscleGroup, limit, offset)
	if err != nil {
		h.logger.Error("failed to list exercises", "error", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{"data": exercises})
}

// GetExercise handles GET /exercises/{id}
func (h *Handler) GetExercise(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing exercise ID"})
		return
	}

	exercise, err := h.service.GetExercise(r.Context(), id)
	if err != nil {
		h.logger.Error("failed to get exercise", "error", err, "exercise_id", id)
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "exercise not found"})
		return
	}

	writeJSON(w, http.StatusOK, exercise)
}

// ListPrograms handles GET /programs
func (h *Handler) ListPrograms(w http.ResponseWriter, r *http.Request) {
	programType := r.URL.Query().Get("type")
	difficulty := r.URL.Query().Get("difficulty")
	equipment := r.URL.Query().Get("equipment")
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	programs, err := h.service.ListPrograms(r.Context(), programType, difficulty, equipment, limit, offset)
	if err != nil {
		h.logger.Error("failed to list programs", "error", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{"data": programs})
}

// GetProgram handles GET /programs/{id}
func (h *Handler) GetProgram(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing program ID"})
		return
	}

	program, err := h.service.GetProgram(r.Context(), id)
	if err != nil {
		h.logger.Error("failed to get program", "error", err, "program_id", id)
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "program not found"})
		return
	}

	writeJSON(w, http.StatusOK, program)
}

// EnrollInProgram handles POST /programs/{id}/enroll
func (h *Handler) EnrollInProgram(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	programID := chi.URLParam(r, "id")
	if programID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing program ID"})
		return
	}

	enrollment, err := h.service.EnrollInProgram(r.Context(), userID, programID)
	if err != nil {
		h.logger.Error("failed to enroll in program", "error", err, "user_id", userID, "program_id", programID)
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, enrollment)
}

// GetCurrentProgram handles GET /programs/current
func (h *Handler) GetCurrentProgram(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	enrollment, err := h.service.GetCurrentProgram(r.Context(), userID)
	if err != nil {
		h.logger.Error("failed to get current program", "error", err, "user_id", userID)
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "no active program enrollment"})
		return
	}

	writeJSON(w, http.StatusOK, enrollment)
}

// CompleteDay handles POST /programs/current/complete-day
func (h *Handler) CompleteDay(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	enrollment, err := h.service.CompleteDay(r.Context(), userID)
	if err != nil {
		h.logger.Error("failed to complete day", "error", err, "user_id", userID)
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, enrollment)
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.Error("failed to encode JSON response", "error", err)
	}
}
