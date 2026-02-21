package member

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/urja-gym/urja/internal/leaderboard"
	"github.com/urja-gym/urja/pkg/middleware"
)

// Handler handles HTTP requests for member endpoints.
type Handler struct {
	service            *Service
	leaderboardService *leaderboard.Service
	logger             *slog.Logger
}

// NewHandler creates a new member handler.
func NewHandler(service *Service, leaderboardService *leaderboard.Service, logger *slog.Logger) *Handler {
	return &Handler{service: service, leaderboardService: leaderboardService, logger: logger}
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
		if errors.Is(err, ErrMemberNotFound) {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "member not found"})
			return
		}
		h.logger.Error("failed to get profile", "error", err, "user_id", userID)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		return
	}

	writeJSON(w, http.StatusOK, member)
}

// UpdateMe handles PUT /api/v1/members/me
func (h *Handler) UpdateMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	var req ProfileUpdate
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	if err := h.service.UpdateProfile(r.Context(), userID, &req); err != nil {
		if errors.Is(err, ErrMemberNotFound) {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "member not found"})
			return
		}
		// Validation errors are returned directly from the service
		if isValidationError(err) {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		h.logger.Error("failed to update profile", "error", err, "user_id", userID)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "profile updated"})
}

// UpdatePrivacy handles PUT /api/v1/members/me/privacy
func (h *Handler) UpdatePrivacy(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	var settings PrivacySettings
	if err := json.NewDecoder(r.Body).Decode(&settings); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	if err := h.service.UpdatePrivacy(r.Context(), userID, settings); err != nil {
		if errors.Is(err, ErrMemberNotFound) {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "member not found"})
			return
		}
		h.logger.Error("failed to update privacy", "error", err, "user_id", userID)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "privacy settings updated"})
}

// GetMyLeaderboard handles GET /api/v1/members/me/leaderboard
func (h *Handler) GetMyLeaderboard(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	period := r.URL.Query().Get("period")
	if period == "" {
		period = "monthly"
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 {
		limit = 3
	}

	orgID, err := h.service.GetPrimaryOrgID(r.Context(), userID)
	if err != nil {
		h.logger.Error("failed to get primary org", "error", err, "user_id", userID)
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "no active organization membership found"})
		return
	}

	result, err := h.leaderboardService.GetLeaderboard(r.Context(), orgID, period, limit, 0)
	if err != nil {
		h.logger.Error("failed to get leaderboard", "error", err, "org_id", orgID)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		return
	}

	writeJSON(w, http.StatusOK, result)
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

// CreateOrgMember handles POST /api/v1/orgs/{orgId}/members
func (h *Handler) CreateOrgMember(w http.ResponseWriter, r *http.Request) {
	orgID, ok := middleware.OrgIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing organization"})
		return
	}

	var input CreateMemberInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	member, err := h.service.CreateOrgMember(r.Context(), orgID, &input)
	if err != nil {
		if errors.Is(err, ErrAlreadyOrgMember) {
			writeJSON(w, http.StatusConflict, map[string]string{"error": "user is already a member of this organization"})
			return
		}
		if isValidationError(err) {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		h.logger.Error("failed to create org member", "error", err, "org_id", orgID)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		return
	}

	writeJSON(w, http.StatusCreated, member)
}

// UpdateOrgMember handles PUT /api/v1/orgs/{orgId}/members/{memberId}
func (h *Handler) UpdateOrgMember(w http.ResponseWriter, r *http.Request) {
	orgID, ok := middleware.OrgIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing organization"})
		return
	}

	memberID := chi.URLParam(r, "memberId")
	if memberID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing member ID"})
		return
	}

	var input UpdateOrgMemberInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	if err := h.service.UpdateOrgMember(r.Context(), orgID, memberID, &input); err != nil {
		if errors.Is(err, ErrMemberNotFound) {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "member not found in this organization"})
			return
		}
		if isValidationError(err) {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		h.logger.Error("failed to update org member", "error", err, "org_id", orgID, "member_id", memberID)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "member updated"})
}

// DeleteOrgMember handles DELETE /api/v1/orgs/{orgId}/members/{memberId}
func (h *Handler) DeleteOrgMember(w http.ResponseWriter, r *http.Request) {
	orgID, ok := middleware.OrgIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing organization"})
		return
	}

	memberID := chi.URLParam(r, "memberId")
	if memberID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing member ID"})
		return
	}

	if err := h.service.DeleteOrgMember(r.Context(), orgID, memberID); err != nil {
		if errors.Is(err, ErrMemberNotFound) {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "member not found in this organization"})
			return
		}
		h.logger.Error("failed to delete org member", "error", err, "org_id", orgID, "member_id", memberID)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "member removed from organization"})
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.Error("failed to encode JSON response", "error", err)
	}
}

// isValidationError returns true for errors that should be returned as 400 Bad Request.
// These are errors originating from input validation in the service layer, not wrapped
// database or system errors.
func isValidationError(err error) bool {
	// Validation errors from the service layer are plain errors (not wrapped with fmt.Errorf %w).
	// If the error has no wrapper (Unwrap returns nil), it's a validation error.
	// Wrapped errors from repo/system always use fmt.Errorf("...: %w", err).
	unwrapper, ok := err.(interface{ Unwrap() error })
	if !ok {
		return true
	}
	return unwrapper.Unwrap() == nil
}
