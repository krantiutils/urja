package activitylog

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/urja-gym/urja/pkg/middleware"
)

// Handler handles HTTP requests for activity log endpoints.
type Handler struct {
	service *Service
	logger  *slog.Logger
}

// NewHandler creates a new activity log handler.
func NewHandler(service *Service, logger *slog.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

// List handles GET /api/v1/orgs/{orgId}/activity-logs
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	orgID, ok := middleware.OrgIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing organization"})
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	var from, to *time.Time
	if fromStr := r.URL.Query().Get("from"); fromStr != "" {
		t, err := time.Parse("2006-01-02", fromStr)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid 'from' date format, expected YYYY-MM-DD"})
			return
		}
		from = &t
	}
	if toStr := r.URL.Query().Get("to"); toStr != "" {
		t, err := time.Parse("2006-01-02", toStr)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid 'to' date format, expected YYYY-MM-DD"})
			return
		}
		// Include the entire "to" day
		endOfDay := t.Add(24*time.Hour - time.Nanosecond)
		to = &endOfDay
	}

	groups, total, err := h.service.List(r.Context(), orgID, from, to, limit, offset)
	if err != nil {
		h.logger.Error("failed to list activity logs", "error", err, "org_id", orgID)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"data":  groups,
		"total": total,
	})
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.Error("failed to encode JSON response", "error", err)
	}
}
