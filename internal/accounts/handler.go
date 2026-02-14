package accounts

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/urja-gym/urja/pkg/middleware"
)

// Handler handles HTTP requests for accounts/transaction endpoints.
type Handler struct {
	service *Service
	logger  *slog.Logger
}

// NewHandler creates a new accounts handler.
func NewHandler(service *Service, logger *slog.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

// List handles GET /api/v1/orgs/{orgId}/accounts
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	orgID, ok := middleware.OrgIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing organization"})
		return
	}

	q := r.URL.Query()
	from := q.Get("from")
	to := q.Get("to")
	category := q.Get("category")
	txnType := q.Get("type")
	paymentType := q.Get("payment_type")
	limit, _ := strconv.Atoi(q.Get("limit"))
	offset, _ := strconv.Atoi(q.Get("offset"))

	txns, total, err := h.service.List(r.Context(), orgID, from, to, category, txnType, paymentType, limit, offset)
	if err != nil {
		h.logger.Error("failed to list transactions", "error", err, "org_id", orgID)
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"data":  txns,
		"total": total,
	})
}

type createTransactionRequest struct {
	Category        string  `json:"category"`
	Description     string  `json:"description,omitempty"`
	TransactionDate string  `json:"transaction_date"`
	TransactionType string  `json:"transaction_type"`
	Amount          float64 `json:"amount"`
	PaymentType     string  `json:"payment_type"`
	Reference       string  `json:"reference,omitempty"`
}

// Create handles POST /api/v1/orgs/{orgId}/accounts
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

	var req createTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	txn, err := h.service.Create(
		r.Context(), orgID,
		req.Category, req.Description, req.TransactionDate,
		req.TransactionType, req.Amount, req.PaymentType,
		req.Reference, userID,
	)
	if err != nil {
		h.logger.Error("failed to create transaction", "error", err, "org_id", orgID)
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusCreated, txn)
}

type updateTransactionRequest struct {
	Category        string  `json:"category"`
	Description     string  `json:"description,omitempty"`
	TransactionDate string  `json:"transaction_date"`
	TransactionType string  `json:"transaction_type"`
	Amount          float64 `json:"amount"`
	PaymentType     string  `json:"payment_type"`
	Reference       string  `json:"reference,omitempty"`
}

// Update handles PUT /api/v1/orgs/{orgId}/accounts/{id}
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	orgID, ok := middleware.OrgIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing organization"})
		return
	}

	txnID := chi.URLParam(r, "id")
	if txnID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "transaction ID is required"})
		return
	}

	var req updateTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	txn, err := h.service.Update(
		r.Context(), orgID, txnID,
		req.Category, req.Description, req.TransactionDate,
		req.TransactionType, req.Amount, req.PaymentType,
		req.Reference,
	)
	if err != nil {
		h.logger.Error("failed to update transaction", "error", err, "org_id", orgID, "txn_id", txnID)
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, txn)
}

// Delete handles DELETE /api/v1/orgs/{orgId}/accounts/{id}
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	orgID, ok := middleware.OrgIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing organization"})
		return
	}

	txnID := chi.URLParam(r, "id")
	if txnID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "transaction ID is required"})
		return
	}

	if err := h.service.Delete(r.Context(), orgID, txnID); err != nil {
		h.logger.Error("failed to delete transaction", "error", err, "org_id", orgID, "txn_id", txnID)
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "transaction deleted"})
}

// GetSummary handles GET /api/v1/orgs/{orgId}/accounts/summary
func (h *Handler) GetSummary(w http.ResponseWriter, r *http.Request) {
	orgID, ok := middleware.OrgIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing organization"})
		return
	}

	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")

	summary, err := h.service.GetSummary(r.Context(), orgID, from, to)
	if err != nil {
		h.logger.Error("failed to get summary", "error", err, "org_id", orgID)
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, summary)
}

// Export handles GET /api/v1/orgs/{orgId}/accounts/export
func (h *Handler) Export(w http.ResponseWriter, r *http.Request) {
	orgID, ok := middleware.OrgIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing organization"})
		return
	}

	q := r.URL.Query()
	from := q.Get("from")
	to := q.Get("to")
	category := q.Get("category")
	txnType := q.Get("type")
	paymentType := q.Get("payment_type")

	txns, err := h.service.Export(r.Context(), orgID, from, to, category, txnType, paymentType)
	if err != nil {
		h.logger.Error("failed to export transactions", "error", err, "org_id", orgID)
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment; filename=transactions.csv")
	w.WriteHeader(http.StatusOK)

	// CSV header
	fmt.Fprintf(w, "SN,Category,Description,Date,Type,Amount,Payment Type,Reference,Entry By\n")
	for i, t := range txns {
		fmt.Fprintf(w, "%d,%s,%s,%s,%s,%.2f,%s,%s,%s\n",
			i+1,
			csvEscape(t.Category),
			csvEscape(t.Description),
			t.TransactionDate,
			t.TransactionType,
			t.Amount,
			t.PaymentType,
			csvEscape(t.Reference),
			csvEscape(t.EntryByName),
		)
	}
}

func csvEscape(s string) string {
	if s == "" {
		return ""
	}
	// Wrap in quotes if it contains commas, quotes, or newlines
	for _, c := range s {
		if c == ',' || c == '"' || c == '\n' || c == '\r' {
			return `"` + strings.ReplaceAll(s, `"`, `""`) + `"`
		}
	}
	return s
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.Error("failed to encode JSON response", "error", err)
	}
}
