package member

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/urja-gym/urja/pkg/middleware"
)

// Handler handles HTTP requests for member endpoints.
type Handler struct {
	service *Service
	logger  *slog.Logger
}

// NewHandler creates a new member handler.
func NewHandler(service *Service, logger *slog.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

// GetMe handles GET /api/v1/members/me
func (h *Handler) GetMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	member, err := h.service.GetProfile(r.Context(), userID)
	if err != nil {
		h.logger.Error("failed to get profile", "error", err, "user_id", userID)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		return
	}

	writeJSON(w, http.StatusOK, member)
}

type updateProfileRequest struct {
	Name      *string `json:"name,omitempty"`
	Email     *string `json:"email,omitempty"`
	AvatarURL *string `json:"avatar_url,omitempty"`
}

// UpdateMe handles PUT /api/v1/members/me
func (h *Handler) UpdateMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	var req updateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	if err := h.service.UpdateProfile(r.Context(), userID, req.Name, req.Email, req.AvatarURL); err != nil {
		h.logger.Error("failed to update profile", "error", err, "user_id", userID)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "profile updated"})
}

// ListByOrg handles GET /api/v1/orgs/{orgId}/members
func (h *Handler) ListByOrg(w http.ResponseWriter, r *http.Request) {
	orgID, ok := middleware.OrgIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing organization"})
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	members, total, err := h.service.ListByOrg(r.Context(), orgID, limit, offset)
	if err != nil {
		h.logger.Error("failed to list members", "error", err, "org_id", orgID)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"data":  members,
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
