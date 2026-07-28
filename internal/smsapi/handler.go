package smsapi

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/urja-gym/urja/pkg/middleware"
)

// Handler handles HTTP requests for SMS endpoints.
type Handler struct {
	service *Service
	logger  *slog.Logger
}

// NewHandler creates a new SMS API handler.
func NewHandler(service *Service, logger *slog.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

// Rate and amount are deliberately absent: the server prices the purchase.
// Accepting them from the caller is how ten thousand credits could be bought
// for a recorded one rupee.
type buyRequest struct {
	Quantity      int    `json:"quantity"`
	PaymentMethod string `json:"payment_method"`
}

type sendRequest struct {
	Message   string   `json:"message"`
	MemberIDs []string `json:"member_ids"`
}

// GetBalance handles GET /api/v1/orgs/{orgId}/sms/balance
func (h *Handler) GetBalance(w http.ResponseWriter, r *http.Request) {
	orgID, ok := middleware.OrgIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing organization"})
		return
	}

	balance, err := h.service.GetBalance(r.Context(), orgID)
	if err != nil {
		h.logger.Error("failed to get SMS balance", "error", err, "org_id", orgID)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		return
	}

	writeJSON(w, http.StatusOK, balance)
}

// Buy handles POST /api/v1/orgs/{orgId}/sms/buy
func (h *Handler) Buy(w http.ResponseWriter, r *http.Request) {
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

	var req buyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	purchase, err := h.service.BuyCredits(r.Context(), orgID, req.Quantity, req.PaymentMethod, userID)
	if err != nil {
		h.logger.Error("failed to buy SMS credits", "error", err, "org_id", orgID)
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusCreated, purchase)
}

// Send handles POST /api/v1/orgs/{orgId}/sms/send
func (h *Handler) Send(w http.ResponseWriter, r *http.Request) {
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

	var req sendRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	campaign, err := h.service.SendSMS(r.Context(), orgID, req.Message, userID, req.MemberIDs)
	if err != nil {
		h.logger.Error("failed to send SMS", "error", err, "org_id", orgID)
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusCreated, campaign)
}

// History handles GET /api/v1/orgs/{orgId}/sms/history
func (h *Handler) History(w http.ResponseWriter, r *http.Request) {
	orgID, ok := middleware.OrgIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing organization"})
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	purchases, err := h.service.ListPurchases(r.Context(), orgID, limit, offset)
	if err != nil {
		h.logger.Error("failed to list SMS purchases", "error", err, "org_id", orgID)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{"data": purchases})
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.Error("failed to encode JSON response", "error", err)
	}
}
