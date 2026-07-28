package subscription

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/urja-gym/urja/pkg/middleware"
)

// Handler handles HTTP requests for subscription endpoints.
type Handler struct {
	service *Service
	logger  *slog.Logger
}

// NewHandler creates a new subscription handler.
func NewHandler(service *Service, logger *slog.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

// PackageSummary handles GET /api/v1/orgs/{orgId}/packages/summary
func (h *Handler) PackageSummary(w http.ResponseWriter, r *http.Request) {
	orgID, ok := middleware.OrgIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing organization"})
		return
	}

	summaries, err := h.service.GetPackageSummary(r.Context(), orgID)
	if err != nil {
		h.logger.Error("failed to get package summary", "error", err, "org_id", orgID)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{"data": summaries})
}

// ListExpiring handles GET /api/v1/orgs/{orgId}/packages/expiring
func (h *Handler) ListExpiring(w http.ResponseWriter, r *http.Request) {
	orgID, ok := middleware.OrgIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing organization"})
		return
	}

	days, _ := strconv.Atoi(r.URL.Query().Get("days"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	entries, total, err := h.service.ListExpiring(r.Context(), orgID, days, limit, offset)
	if err != nil {
		h.logger.Error("failed to list expiring packages", "error", err, "org_id", orgID)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		return
	}
	if entries == nil {
		entries = []ExpiringEntry{}
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"data":  entries,
		"total": total,
	})
}

// ListExpired handles GET /api/v1/orgs/{orgId}/packages/expired
func (h *Handler) ListExpired(w http.ResponseWriter, r *http.Request) {
	orgID, ok := middleware.OrgIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing organization"})
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	entries, total, err := h.service.ListExpired(r.Context(), orgID, limit, offset)
	if err != nil {
		h.logger.Error("failed to list expired packages", "error", err, "org_id", orgID)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		return
	}
	if entries == nil {
		entries = []ExpiredEntry{}
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"data":  entries,
		"total": total,
	})
}

// AssignPackage handles POST /api/v1/orgs/{orgId}/members/{memberId}/packages/assign
func (h *Handler) AssignPackage(w http.ResponseWriter, r *http.Request) {
	// Recorded on the ledger entry as who took the money.
	entryByUserID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

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

	var req AssignRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	sub, pmt, err := h.service.AssignPackage(r.Context(), orgID, memberID, entryByUserID, req)
	if err != nil {
		h.handleServiceError(w, err, "assign package", orgID, memberID)
		return
	}

	writeJSON(w, http.StatusCreated, map[string]interface{}{
		"subscription": sub,
		"payment":      pmt,
	})
}

// RenewPackage handles POST /api/v1/orgs/{orgId}/members/{memberId}/packages/renew
func (h *Handler) RenewPackage(w http.ResponseWriter, r *http.Request) {
	// Recorded on the ledger entry as who took the money.
	entryByUserID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

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

	var req RenewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	sub, pmt, err := h.service.RenewPackage(r.Context(), orgID, memberID, entryByUserID, req)
	if err != nil {
		h.handleServiceError(w, err, "renew package", orgID, memberID)
		return
	}

	writeJSON(w, http.StatusCreated, map[string]interface{}{
		"subscription": sub,
		"payment":      pmt,
	})
}

// ExtendPackage handles POST /api/v1/orgs/{orgId}/members/{memberId}/packages/extend
func (h *Handler) ExtendPackage(w http.ResponseWriter, r *http.Request) {
	// Recorded on the ledger entry as who took the money.
	entryByUserID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

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

	var req ExtendRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	sub, pmt, err := h.service.ExtendPackage(r.Context(), orgID, memberID, entryByUserID, req)
	if err != nil {
		h.handleServiceError(w, err, "extend package", orgID, memberID)
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"subscription": sub,
		"payment":      pmt,
	})
}

// ListPayments handles GET /api/v1/orgs/{orgId}/members/{memberId}/payments
func (h *Handler) ListPayments(w http.ResponseWriter, r *http.Request) {
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

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	payments, total, err := h.service.ListPayments(r.Context(), orgID, memberID, limit, offset)
	if err != nil {
		h.logger.Error("failed to list payments", "error", err, "org_id", orgID, "member_id", memberID)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"data":  payments,
		"total": total,
	})
}

// ListSubscriptions handles GET /api/v1/orgs/{orgId}/members/{memberId}/subscriptions
func (h *Handler) ListSubscriptions(w http.ResponseWriter, r *http.Request) {
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

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	subs, total, err := h.service.ListSubscriptions(r.Context(), orgID, memberID, limit, offset)
	if err != nil {
		h.logger.Error("failed to list subscriptions", "error", err, "org_id", orgID, "member_id", memberID)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"data":  subs,
		"total": total,
	})
}

// handleServiceError maps service-layer errors to HTTP responses.
func (h *Handler) handleServiceError(w http.ResponseWriter, err error, action, orgID, memberID string) {
	switch {
	case errors.Is(err, ErrInvalidRequest):
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
	case errors.Is(err, ErrInvalidPayment):
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
	case errors.Is(err, ErrInvalidPackage):
		writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
	case errors.Is(err, ErrInvalidMember):
		writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
	case errors.Is(err, ErrSubscriptionNotFound):
		writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
	default:
		h.logger.Error("failed to "+action, "error", err, "org_id", orgID, "member_id", memberID)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.Error("failed to encode JSON response", "error", err)
	}
}
