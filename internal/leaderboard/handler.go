package leaderboard

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/urja-gym/urja/pkg/middleware"
)

// Handler handles HTTP requests for leaderboard endpoints.
type Handler struct {
	service *Service
	logger  *slog.Logger
}

// NewHandler creates a new leaderboard handler.
func NewHandler(service *Service, logger *slog.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

// GetLeaderboard handles GET /api/v1/orgs/{orgId}/leaderboard
func (h *Handler) GetLeaderboard(w http.ResponseWriter, r *http.Request) {
	orgID, ok := middleware.OrgIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing organization"})
		return
	}

	period := r.URL.Query().Get("period")
	if period == "" {
		period = "weekly"
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	result, err := h.service.GetLeaderboard(r.Context(), orgID, period, limit, offset)
	if err != nil {
		h.logger.Error("failed to get leaderboard",
			"error", err, "org_id", orgID, "period", period)
		// Surface validation errors (invalid period) as 400.
		if period != "weekly" && period != "monthly" && period != "alltime" {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		return
	}

	writeJSON(w, http.StatusOK, result)
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.Error("failed to encode JSON response", "error", err)
	}
}
