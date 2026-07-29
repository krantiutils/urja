package invoice

import (
	"errors"
	"time"
)

var (
	// ErrPANNotConfigured is returned when a gym tries to bill before setting
	// its PAN. Handled as its own code so the UI can link to settings rather
	// than dead-ending on a message.
	ErrPANNotConfigured = errors.New("pan_not_configured")
	ErrNoItems          = errors.New("an invoice needs at least one line item")
	ErrInvalidDiscount  = errors.New("discount cannot exceed the subtotal")
	ErrAlreadyBilled    = errors.New("this payment already has a bill")
	ErrNotFound         = errors.New("invoice not found")
	ErrAlreadyCancelled = errors.New("invoice is already cancelled")
	ErrReasonRequired   = errors.New("a cancellation reason is required")
	ErrInvoiceCancelled = errors.New("cannot credit a cancelled invoice")
	ErrInvalidParent    = errors.New("a credit note cannot reference another credit note")
	ErrCreditTooLarge   = errors.New("credit exceeds the uncredited balance of the invoice")

	// These three guard against a client-supplied foreign key that resolves
	// in the database but belongs to a different org. Without this check, a
	// cross-org transaction_id would not just misattribute a bill — because
	// idx_invoices_one_per_transaction is a global unique index, it would
	// permanently consume another org's one-bill slot for that transaction.
	ErrTransactionNotInOrg   = errors.New("transaction not found in this organization")
	ErrMemberPackageNotInOrg = errors.New("member package not found in this organization")
	ErrCustomerNotInOrg      = errors.New("customer is not an active member of this organization")

	// ErrTransactionAlreadyOwned guards the mirror image of the bug task 8b
	// fixed: an invoice that created its own income row owns that transaction
	// for the rest of that row's life, cancelled or not. Without this check, a
	// cancelled from-scratch bill's reversed transaction_id could be quoted by
	// a brand-new bill — idx_invoices_one_per_transaction only blocks a second
	// *issued* invoice on the same row, so the cancelled one doesn't count —
	// producing a live invoice with owns_transaction = false and no income of
	// its own, while the ledger already nets that money to zero. 400, not 404:
	// the row genuinely belongs to this org, the caller just may not link it.
	ErrTransactionAlreadyOwned = errors.New("this transaction already belongs to another invoice and cannot be linked")
)

// Item is one line on a bill.
type Item struct {
	LineNo        int     `json:"line_no"`
	Description   string  `json:"description"`
	DescriptionNe string  `json:"description_ne,omitempty"`
	Quantity      float64 `json:"quantity"`
	UnitPrice     float64 `json:"unit_price"`
	Amount        float64 `json:"amount"`
}

// Invoice is an issued tax document. Every field a bill prints is snapshotted
// here rather than joined at read time, so a bill never changes retroactively
// when the gym or the customer edits their details.
type Invoice struct {
	ID            string `json:"id"`
	OrgID         string `json:"organization_id"`
	FiscalYear    string `json:"fiscal_year"`
	Sequence      int    `json:"sequence"`
	InvoiceNumber string `json:"invoice_number"`
	DocType       string `json:"doc_type"`
	CreditNoteFor string `json:"credit_note_for,omitempty"`

	SellerName          string `json:"seller_name"`
	SellerPAN           string `json:"seller_pan"`
	SellerAddress       string `json:"seller_address,omitempty"`
	SellerVATRegistered bool   `json:"seller_vat_registered"`

	CustomerUserID  string `json:"customer_user_id,omitempty"`
	CustomerName    string `json:"customer_name"`
	CustomerPAN     string `json:"customer_pan,omitempty"`
	CustomerAddress string `json:"customer_address,omitempty"`
	CustomerPhone   string `json:"customer_phone,omitempty"`

	IssuedDate   string `json:"issued_date"`
	IssuedDateBS string `json:"issued_date_bs"`

	Subtotal      float64 `json:"subtotal"`
	Discount      float64 `json:"discount"`
	TaxableAmount float64 `json:"taxable_amount"`
	VATRate       float64 `json:"vat_rate"`
	VATAmount     float64 `json:"vat_amount"`
	Total         float64 `json:"total"`
	AmountInWords string  `json:"amount_in_words"`

	PaymentMethod string `json:"payment_method,omitempty"`

	Status             string     `json:"status"`
	CancelledAt        *time.Time `json:"cancelled_at,omitempty"`
	CancellationReason string     `json:"cancellation_reason,omitempty"`

	TransactionID   string `json:"transaction_id,omitempty"`
	MemberPackageID string `json:"member_package_id,omitempty"`

	IssuedBy   string    `json:"issued_by"`
	PrintCount int       `json:"print_count"`
	CreatedAt  time.Time `json:"created_at"`

	Items []Item `json:"items,omitempty"`
}

// ItemInput is a requested line on a new bill.
type ItemInput struct {
	Description   string  `json:"description"`
	DescriptionNe string  `json:"description_ne"`
	Quantity      float64 `json:"quantity"`
	UnitPrice     float64 `json:"unit_price"`
}

// IssueInput is a request to issue a bill.
type IssueInput struct {
	CustomerUserID  string      `json:"customer_user_id"`
	CustomerName    string      `json:"customer_name"`
	CustomerPAN     string      `json:"customer_pan"`
	CustomerAddress string      `json:"customer_address"`
	CustomerPhone   string      `json:"customer_phone"`
	PaymentMethod   string      `json:"payment_method"`
	TransactionID   string      `json:"transaction_id"`
	MemberPackageID string      `json:"member_package_id"`
	Discount        float64     `json:"discount"`
	Items           []ItemInput `json:"items"`
}

// CreditInput is a request to credit part or all of an invoice.
type CreditInput struct {
	Reason string      `json:"reason"`
	Items  []ItemInput `json:"items"`
}

// ListFilter narrows a list query.
type ListFilter struct {
	OrgID      string
	Status     string
	FiscalYear string
	CustomerID string
	From       string
	To         string
	Query      string
	Limit      int
	Offset     int
}
