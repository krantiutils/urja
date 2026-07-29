package invoice

import (
	"context"
	"fmt"
	"log/slog"
	"math"
	"strings"
	"time"

	"github.com/urja-gym/urja/pkg/moneywords"
	"github.com/urja-gym/urja/pkg/nepalidate"
)

// Service holds invoice business logic.
type Service struct {
	repo   *Repository
	logger *slog.Logger
}

// NewService creates a new invoice service.
func NewService(repo *Repository, logger *slog.Logger) *Service {
	return &Service{repo: repo, logger: logger}
}

func round2(f float64) float64 {
	return math.Round(f*100) / 100
}

var validPaymentMethods = map[string]bool{
	"cash": true, "bank_transfer": true, "khalti": true,
}

// validateIssue checks an issue request and returns the computed money.
func validateIssue(in IssueInput) (subtotal, taxable float64, err error) {
	if strings.TrimSpace(in.CustomerName) == "" {
		return 0, 0, fmt.Errorf("customer name is required")
	}
	if len(in.Items) == 0 {
		return 0, 0, ErrNoItems
	}
	if in.PaymentMethod != "" && !validPaymentMethods[in.PaymentMethod] {
		return 0, 0, fmt.Errorf("invalid payment method: %s", in.PaymentMethod)
	}
	if in.CustomerPAN != "" && !isNineDigits(in.CustomerPAN) {
		return 0, 0, fmt.Errorf("customer PAN must be exactly 9 digits")
	}

	for i, it := range in.Items {
		if strings.TrimSpace(it.Description) == "" {
			return 0, 0, fmt.Errorf("line %d: description is required", i+1)
		}
		if it.Quantity <= 0 {
			return 0, 0, fmt.Errorf("line %d: quantity must be greater than 0", i+1)
		}
		if it.UnitPrice < 0 {
			return 0, 0, fmt.Errorf("line %d: unit price cannot be negative", i+1)
		}
		subtotal += round2(it.Quantity * it.UnitPrice)
	}
	subtotal = round2(subtotal)

	if in.Discount < 0 {
		return 0, 0, fmt.Errorf("discount cannot be negative")
	}
	if in.Discount > subtotal {
		return 0, 0, ErrInvalidDiscount
	}

	return subtotal, round2(subtotal - in.Discount), nil
}

func isNineDigits(s string) bool {
	if len(s) != 9 {
		return false
	}
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

// Issue validates and writes a new bill.
func (s *Service) Issue(ctx context.Context, orgID, issuedBy string, in IssueInput) (*Invoice, error) {
	subtotal, taxable, err := validateIssue(in)
	if err != nil {
		return nil, err
	}

	bs, err := nepalidate.Today()
	if err != nil {
		return nil, fmt.Errorf("determining Nepali date: %w", err)
	}
	loc, _ := time.LoadLocation("Asia/Kathmandu")

	inv, err := s.repo.Issue(ctx, issueParams{
		OrgID:         orgID,
		FiscalYear:    bs.FiscalYear(),
		DocType:       "invoice",
		IssuedDate:    time.Now().In(loc).Format("2006-01-02"),
		IssuedDateBS:  bs.String(),
		Subtotal:      subtotal,
		Discount:      round2(in.Discount),
		TaxableAmount: taxable,
		// PAN-only: no VAT is added, so the total is the taxable amount.
		Total:         taxable,
		AmountInWords: moneywords.Rupees(taxable),
		IssuedBy:      issuedBy,
		In:            in,
	})
	if err != nil {
		return nil, err
	}

	s.logger.Info("invoice issued",
		"invoice_number", inv.InvoiceNumber, "org_id", orgID,
		"total", inv.Total, "issued_by", issuedBy)
	return inv, nil
}

// Cancel withdraws a bill issued in error. It does not touch the ledger — if
// money actually moved, a credit note is the correct instrument.
func (s *Service) Cancel(ctx context.Context, orgID, id, cancelledBy, reason string) (*Invoice, error) {
	reason = strings.TrimSpace(reason)
	if reason == "" {
		return nil, ErrReasonRequired
	}
	inv, err := s.repo.Cancel(ctx, orgID, id, cancelledBy, reason)
	if err != nil {
		return nil, err
	}
	s.logger.Info("invoice cancelled",
		"invoice_number", inv.InvoiceNumber, "org_id", orgID,
		"cancelled_by", cancelledBy, "reason", reason)
	return inv, nil
}

// CreditNote reverses part or all of an issued bill. Unlike cancelling, this
// writes a reversing row to the ledger: the money genuinely goes back.
func (s *Service) CreditNote(ctx context.Context, orgID, parentID, issuedBy string, in CreditInput) (*Invoice, error) {
	if strings.TrimSpace(in.Reason) == "" {
		return nil, ErrReasonRequired
	}

	parent, err := s.repo.Get(ctx, orgID, parentID)
	if err != nil {
		return nil, err
	}

	issue := IssueInput{
		CustomerUserID:  parent.CustomerUserID,
		CustomerName:    parent.CustomerName,
		CustomerPAN:     parent.CustomerPAN,
		CustomerAddress: parent.CustomerAddress,
		CustomerPhone:   parent.CustomerPhone,
		PaymentMethod:   parent.PaymentMethod,
		Items:           in.Items,
	}
	subtotal, taxable, err := validateIssue(issue)
	if err != nil {
		return nil, err
	}

	bs, err := nepalidate.Today()
	if err != nil {
		return nil, fmt.Errorf("determining Nepali date: %w", err)
	}
	loc, _ := time.LoadLocation("Asia/Kathmandu")

	note, err := s.repo.CreditNote(ctx, issueParams{
		OrgID:         orgID,
		FiscalYear:    bs.FiscalYear(),
		DocType:       "credit_note",
		IssuedDate:    time.Now().In(loc).Format("2006-01-02"),
		IssuedDateBS:  bs.String(),
		Subtotal:      subtotal,
		TaxableAmount: taxable,
		Total:         taxable,
		AmountInWords: moneywords.Rupees(taxable),
		IssuedBy:      issuedBy,
		In:            issue,
	}, parentID)
	if err != nil {
		return nil, err
	}

	s.logger.Info("credit note raised",
		"credit_note", note.InvoiceNumber, "against", parent.InvoiceNumber,
		"org_id", orgID, "amount", note.Total, "reason", in.Reason)
	return note, nil
}

// Get reads one invoice, scoped to the org.
func (s *Service) Get(ctx context.Context, orgID, id string) (*Invoice, error) {
	return s.repo.Get(ctx, orgID, id)
}

// List returns invoices for an org.
func (s *Service) List(ctx context.Context, f ListFilter) ([]Invoice, int, error) {
	if f.Limit <= 0 || f.Limit > 100 {
		f.Limit = 20
	}
	if f.Offset < 0 {
		f.Offset = 0
	}
	if f.Status != "" && f.Status != "issued" && f.Status != "cancelled" {
		return nil, 0, fmt.Errorf("invalid status filter: %s", f.Status)
	}
	return s.repo.List(ctx, f)
}
