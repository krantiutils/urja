package boxing

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/urja-gym/urja/pkg/middleware"
)

// Handler handles HTTP requests for boxing profile endpoints.
type Handler struct {
	service *Service
	logger  *slog.Logger
}

// NewHandler creates a new boxing handler.
func NewHandler(service *Service, logger *slog.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

// updateProfileRequest is the member-editable shape.
//
// sparring_cleared is deliberately absent rather than validated-and-rejected:
// a field a member can never set should not exist in the struct they post to,
// so there is no path by which a stray decode could set it.
type updateProfileRequest struct {
	Stance      string   `json:"stance"`
	WeightClass string   `json:"weight_class"`
	SkillLevel  string   `json:"skill_level"`
	ReachCm     *float64 `json:"reach_cm"`
	Notes       string   `json:"notes"`
}

func (r updateProfileRequest) toInput() UpdateProfileInput {
	return UpdateProfileInput{
		Stance: r.Stance, WeightClass: r.WeightClass, SkillLevel: r.SkillLevel,
		ReachCm: r.ReachCm, Notes: r.Notes,
	}
}

type boutRequest struct {
	BoutDate    string `json:"bout_date"`
	Opponent    string `json:"opponent"`
	EventName   string `json:"event_name"`
	Result      string `json:"result"`
	Method      string `json:"method"`
	Rounds      *int   `json:"rounds"`
	WeightClass string `json:"weight_class"`
	Notes       string `json:"notes"`
}

func (r boutRequest) toInput() CreateBoutInput {
	return CreateBoutInput{
		BoutDate: r.BoutDate, Opponent: r.Opponent, EventName: r.EventName,
		Result: r.Result, Method: r.Method, Rounds: r.Rounds,
		WeightClass: r.WeightClass, Notes: r.Notes,
	}
}

// --- Member self-service (under /members/me) ---

// selfContext resolves the caller and the organization they are acting in.
// Self routes are not mounted under /orgs/{orgId}, so the org arrives as a
// query parameter and must be checked against real membership.
func (h *Handler) selfContext(w http.ResponseWriter, r *http.Request) (userID, orgID string, ok bool) {
	userID, found := middleware.UserIDFromContext(r.Context())
	if !found {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return "", "", false
	}

	orgID = r.URL.Query().Get("organization_id")
	if orgID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "organization_id is required"})
		return "", "", false
	}

	// Without this, any authenticated user could read or write a profile in an
	// organization they do not belong to simply by naming it.
	if err := h.service.EnsureOrgMember(r.Context(), orgID, userID); err != nil {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "not a member of this organization"})
		return "", "", false
	}

	return userID, orgID, true
}

// GetMyProfile handles GET /api/v1/members/me/boxing
func (h *Handler) GetMyProfile(w http.ResponseWriter, r *http.Request) {
	userID, orgID, ok := h.selfContext(w, r)
	if !ok {
		return
	}

	view, err := h.service.GetProfile(r.Context(), orgID, userID)
	if err != nil {
		h.logger.Error("failed to get boxing profile", "error", err, "user_id", userID)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		return
	}

	writeJSON(w, http.StatusOK, view)
}

// UpdateMyProfile handles PUT /api/v1/members/me/boxing
func (h *Handler) UpdateMyProfile(w http.ResponseWriter, r *http.Request) {
	userID, orgID, ok := h.selfContext(w, r)
	if !ok {
		return
	}

	var req updateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	view, err := h.service.UpdateProfile(r.Context(), orgID, userID, req.toInput())
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, view)
}

// ListMyBouts handles GET /api/v1/members/me/bouts
func (h *Handler) ListMyBouts(w http.ResponseWriter, r *http.Request) {
	userID, orgID, ok := h.selfContext(w, r)
	if !ok {
		return
	}

	view, err := h.service.GetProfile(r.Context(), orgID, userID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{"data": view.Bouts, "record": view.Record})
}

// CreateMyBout handles POST /api/v1/members/me/bouts
func (h *Handler) CreateMyBout(w http.ResponseWriter, r *http.Request) {
	userID, orgID, ok := h.selfContext(w, r)
	if !ok {
		return
	}

	var req boutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	bout, err := h.service.CreateBout(r.Context(), orgID, userID, req.toInput())
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusCreated, bout)
}

// DeleteMyBout handles DELETE /api/v1/members/me/bouts/{boutId}
func (h *Handler) DeleteMyBout(w http.ResponseWriter, r *http.Request) {
	userID, orgID, ok := h.selfContext(w, r)
	if !ok {
		return
	}

	if err := h.service.DeleteBout(r.Context(), orgID, userID, chi.URLParam(r, "boutId")); err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "bout not found"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "bout deleted"})
}

// --- Staff routes (under /orgs/{orgId}/members/{memberId}) ---

func (h *Handler) staffContext(w http.ResponseWriter, r *http.Request) (orgID, memberID string, ok bool) {
	orgID, found := middleware.OrgIDFromContext(r.Context())
	if !found {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing organization"})
		return "", "", false
	}

	memberID = chi.URLParam(r, "memberId")
	if memberID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing member ID"})
		return "", "", false
	}

	// Staff of gym A must not be able to write a profile for a stranger by
	// passing an arbitrary user id.
	if err := h.service.EnsureOrgMember(r.Context(), orgID, memberID); err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "member not found in this organization"})
		return "", "", false
	}

	return orgID, memberID, true
}

// GetMemberProfile handles GET /orgs/{orgId}/members/{memberId}/boxing
func (h *Handler) GetMemberProfile(w http.ResponseWriter, r *http.Request) {
	orgID, memberID, ok := h.staffContext(w, r)
	if !ok {
		return
	}

	view, err := h.service.GetProfile(r.Context(), orgID, memberID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		return
	}

	writeJSON(w, http.StatusOK, view)
}

// UpdateMemberProfile handles PUT /orgs/{orgId}/members/{memberId}/boxing
func (h *Handler) UpdateMemberProfile(w http.ResponseWriter, r *http.Request) {
	orgID, memberID, ok := h.staffContext(w, r)
	if !ok {
		return
	}

	var req updateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	view, err := h.service.UpdateProfile(r.Context(), orgID, memberID, req.toInput())
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, view)
}

type clearanceRequest struct {
	SparringCleared bool `json:"sparring_cleared"`
}

// SetSparringClearance handles PUT /orgs/{orgId}/members/{memberId}/boxing/sparring-clearance
func (h *Handler) SetSparringClearance(w http.ResponseWriter, r *http.Request) {
	orgID, memberID, ok := h.staffContext(w, r)
	if !ok {
		return
	}

	callerID, found := middleware.UserIDFromContext(r.Context())
	if !found {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	var req clearanceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	view, err := h.service.SetSparringClearance(r.Context(), orgID, callerID, memberID, req.SparringCleared)
	if err != nil {
		if errors.Is(err, ErrForbidden) {
			writeJSON(w, http.StatusForbidden, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, view)
}

// CreateMemberBout handles POST /orgs/{orgId}/members/{memberId}/bouts
func (h *Handler) CreateMemberBout(w http.ResponseWriter, r *http.Request) {
	orgID, memberID, ok := h.staffContext(w, r)
	if !ok {
		return
	}

	var req boutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	bout, err := h.service.CreateBout(r.Context(), orgID, memberID, req.toInput())
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusCreated, bout)
}

// DeleteMemberBout handles DELETE /orgs/{orgId}/members/{memberId}/bouts/{boutId}
func (h *Handler) DeleteMemberBout(w http.ResponseWriter, r *http.Request) {
	orgID, memberID, ok := h.staffContext(w, r)
	if !ok {
		return
	}

	if err := h.service.DeleteBout(r.Context(), orgID, memberID, chi.URLParam(r, "boutId")); err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "bout not found"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "bout deleted"})
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.Error("failed to encode JSON response", "error", err)
	}
}
