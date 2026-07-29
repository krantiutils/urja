package invoice

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/urja-gym/urja/pkg/middleware"
)

// Handler serves the invoice HTTP endpoints.
type Handler struct {
	service *Service
	logger  *slog.Logger
}

// NewHandler creates a new invoice handler.
func NewHandler(service *Service, logger *slog.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.Error("failed to encode JSON response", "error", err)
	}
}

// writeErr maps a domain error onto the status and code the UI expects.
// Cross-tenant reads deliberately return 404: a 403 would confirm the row
// exists in another gym.
func writeErr(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrNotFound):
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "invoice not found", "code": "not_found"})
	case errors.Is(err, ErrPANNotConfigured):
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "set your PAN number in settings before issuing a bill",
			"code":  "pan_not_configured"})
	case errors.Is(err, ErrAlreadyCancelled):
		writeJSON(w, http.StatusConflict, map[string]string{"error": err.Error(), "code": "already_cancelled"})
	case errors.Is(err, ErrInvoiceCancelled):
		writeJSON(w, http.StatusConflict, map[string]string{"error": err.Error(), "code": "invoice_cancelled"})
	case errors.Is(err, ErrAlreadyBilled):
		writeJSON(w, http.StatusConflict, map[string]string{"error": err.Error(), "code": "already_billed"})
	case errors.Is(err, ErrCreditTooLarge):
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error(), "code": "credit_exceeds_balance"})
	case errors.Is(err, ErrInvalidParent):
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error(), "code": "invalid_parent"})
	case errors.Is(err, ErrReasonRequired):
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error(), "code": "reason_required"})
	// 404, not 403: a 403 would confirm the referenced row exists in
	// another org, which is exactly the cross-tenant leak this closes.
	case errors.Is(err, ErrTransactionNotInOrg):
		writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error(), "code": "transaction_not_found"})
	case errors.Is(err, ErrMemberPackageNotInOrg):
		writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error(), "code": "member_package_not_found"})
	case errors.Is(err, ErrCustomerNotInOrg):
		writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error(), "code": "customer_not_found"})
	default:
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
}

// Issue handles POST /api/v1/orgs/{orgId}/invoices
func (h *Handler) Issue(w http.ResponseWriter, r *http.Request) {
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

	var in IssueInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	inv, err := h.service.Issue(r.Context(), orgID, userID, in)
	if err != nil {
		h.logger.Error("issuing invoice failed", "error", err, "org_id", orgID)
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, inv)
}

// Get handles GET /api/v1/orgs/{orgId}/invoices/{invoiceId}
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	orgID, ok := middleware.OrgIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing organization"})
		return
	}
	inv, err := h.service.Get(r.Context(), orgID, chi.URLParam(r, "invoiceId"))
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, inv)
}

// List handles GET /api/v1/orgs/{orgId}/invoices
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	orgID, ok := middleware.OrgIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing organization"})
		return
	}
	q := r.URL.Query()
	limit, _ := strconv.Atoi(q.Get("limit"))
	offset, _ := strconv.Atoi(q.Get("offset"))

	list, total, err := h.service.List(r.Context(), ListFilter{
		OrgID:      orgID,
		Status:     q.Get("status"),
		FiscalYear: q.Get("fiscal_year"),
		CustomerID: q.Get("customer_user_id"),
		From:       q.Get("from"),
		To:         q.Get("to"),
		Query:      q.Get("q"),
		Limit:      limit,
		Offset:     offset,
	})
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"data": list, "total": total})
}

// Cancel handles POST /api/v1/orgs/{orgId}/invoices/{invoiceId}/cancel
func (h *Handler) Cancel(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusNotImplemented, map[string]string{"error": "not implemented"})
}

// CreditNote handles POST /api/v1/orgs/{orgId}/invoices/{invoiceId}/credit-note
func (h *Handler) CreditNote(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusNotImplemented, map[string]string{"error": "not implemented"})
}

// Print handles POST /api/v1/orgs/{orgId}/invoices/{invoiceId}/print
func (h *Handler) Print(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusNotImplemented, map[string]string{"error": "not implemented"})
}

// NextNumber handles GET /api/v1/orgs/{orgId}/invoices/next-number
func (h *Handler) NextNumber(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusNotImplemented, map[string]string{"error": "not implemented"})
}
