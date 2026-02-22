package water

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/urja-gym/urja/pkg/middleware"
)

// Handler handles HTTP requests for water tracking endpoints.
type Handler struct {
	service *Service
	logger  *slog.Logger
}

// NewHandler creates a new water handler.
func NewHandler(service *Service, logger *slog.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

// CreateWaterLog handles POST /api/v1/members/me/water-logs
func (h *Handler) CreateWaterLog(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	var input CreateWaterLogInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	l, err := h.service.CreateWaterLog(r.Context(), userID, &input)
	if err != nil {
		h.logger.Error("failed to create water log", "error", err, "user_id", userID)
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusCreated, l)
}

// ListWaterLogs handles GET /api/v1/members/me/water-logs
func (h *Handler) ListWaterLogs(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	date := r.URL.Query().Get("date")

	logs, err := h.service.ListWaterLogs(r.Context(), userID, date)
	if err != nil {
		h.logger.Error("failed to list water logs", "error", err, "user_id", userID)
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{"data": logs})
}

// DeleteWaterLog handles DELETE /api/v1/members/me/water-logs/{id}
func (h *Handler) DeleteWaterLog(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	logID := chi.URLParam(r, "id")
	if logID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing water log ID"})
		return
	}

	if err := h.service.DeleteWaterLog(r.Context(), userID, logID); err != nil {
		h.logger.Error("failed to delete water log", "error", err, "log_id", logID, "user_id", userID)
		writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "water log deleted"})
}

// GetDailySummary handles GET /api/v1/members/me/water-logs/summary
func (h *Handler) GetDailySummary(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	date := r.URL.Query().Get("date")

	summary, err := h.service.GetDailySummary(r.Context(), userID, date)
	if err != nil {
		h.logger.Error("failed to get water daily summary", "error", err, "user_id", userID)
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, summary)
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.Error("failed to encode JSON response", "error", err)
	}
}
