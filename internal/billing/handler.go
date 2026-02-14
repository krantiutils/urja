package billing

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/urja-gym/urja/pkg/middleware"
)

// Handler handles HTTP requests for billing endpoints.
type Handler struct {
	service *Service
	logger  *slog.Logger
}

// NewHandler creates a new billing handler.
func NewHandler(service *Service, logger *slog.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

// ListPlans handles GET /api/v1/billing/plans
func (h *Handler) ListPlans(w http.ResponseWriter, r *http.Request) {
	plans, err := h.service.ListPlans(r.Context())
	if err != nil {
		h.logger.Error("failed to list plans", "error", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{"data": plans})
}

// subscribeRequest is the JSON body for POST /api/v1/billing/subscribe.
type subscribeRequest struct {
	OrganizationID string `json:"organization_id"`
	PlanID         string `json:"plan_id"`
	BillingPeriod  string `json:"billing_period"` // "monthly" or "yearly"
	KhaltiPidx     string `json:"khalti_pidx"`
}

// Subscribe handles POST /api/v1/billing/subscribe
func (h *Handler) Subscribe(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	var req subscribeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	if req.OrganizationID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "organization_id is required"})
		return
	}
	if req.PlanID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "plan_id is required"})
		return
	}
	if req.BillingPeriod == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "billing_period is required"})
		return
	}
	if req.KhaltiPidx == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "khalti_pidx is required"})
		return
	}

	in := &SubscribeInput{
		OrganizationID: req.OrganizationID,
		PlanID:         req.PlanID,
		BillingPeriod:  req.BillingPeriod,
		KhaltiPidx:     req.KhaltiPidx,
	}

	sub, err := h.service.Subscribe(r.Context(), userID, in)
	if err != nil {
		if strings.Contains(err.Error(), "forbidden") {
			writeJSON(w, http.StatusForbidden, map[string]string{"error": err.Error()})
			return
		}
		if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "not active") {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		if strings.Contains(err.Error(), "already has an active") {
			writeJSON(w, http.StatusConflict, map[string]string{"error": err.Error()})
			return
		}
		if strings.Contains(err.Error(), "payment") {
			writeJSON(w, http.StatusPaymentRequired, map[string]string{"error": err.Error()})
			return
		}
		h.logger.Error("failed to subscribe", "error", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to process subscription"})
		return
	}

	writeJSON(w, http.StatusCreated, sub)
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.Error("failed to encode JSON response", "error", err)
	}
}
