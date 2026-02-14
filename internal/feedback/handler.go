package feedback

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/urja-gym/urja/pkg/middleware"
)

// Handler handles HTTP requests for feedback endpoints.
type Handler struct {
	service *Service
	logger  *slog.Logger
}

// NewHandler creates a new feedback handler.
func NewHandler(service *Service, logger *slog.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

type createFeedbackRequest struct {
	Message string `json:"message"`
}

// Create handles POST /api/v1/orgs/{orgId}/feedbacks
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	orgID, ok := middleware.OrgIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing organization"})
		return
	}

	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	var req createFeedbackRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	f, err := h.service.Create(r.Context(), orgID, userID, req.Message)
	if err != nil {
		h.logger.Error("failed to create feedback", "error", err, "org_id", orgID)
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusCreated, f)
}

// List handles GET /api/v1/orgs/{orgId}/feedbacks
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	orgID, ok := middleware.OrgIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing organization"})
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	feedbacks, err := h.service.List(r.Context(), orgID, limit, offset)
	if err != nil {
		h.logger.Error("failed to list feedbacks", "error", err, "org_id", orgID)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{"data": feedbacks})
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.Error("failed to encode JSON response", "error", err)
	}
}
