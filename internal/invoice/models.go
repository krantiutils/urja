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
