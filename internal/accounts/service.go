package accounts

import (
	"context"
	"fmt"
	"log/slog"
	"time"
)

// Service handles accounts/transaction business logic.
type Service struct {
	repo   *Repository
	logger *slog.Logger
}

// NewService creates a new accounts service.
func NewService(repo *Repository, logger *slog.Logger) *Service {
	return &Service{repo: repo, logger: logger}
}

var (
	validTransactionTypes = map[string]bool{
		"income":  true,
		"expense": true,
	}

	validPaymentTypes = map[string]bool{
		"cash":          true,
		"bank_transfer": true,
		"khalti":        true,
	}
)

// List retrieves transactions with filters.
func (s *Service) List(ctx context.Context, orgID, from, to, category, txnType, paymentType string, limit, offset int) ([]Transaction, int, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	if err := validateDateRange(from, to); err != nil {
		return nil, 0, err
	}

	if txnType != "" && !validTransactionTypes[txnType] {
		return nil, 0, fmt.Errorf("invalid transaction type: %s", txnType)
	}

	if paymentType != "" && !validPaymentTypes[paymentType] {
		return nil, 0, fmt.Errorf("invalid payment type: %s", paymentType)
	}

	return s.repo.List(ctx, ListFilter{
		OrgID:           orgID,
		From:            from,
		To:              to,
		Category:        category,
		TransactionType: txnType,
		PaymentType:     paymentType,
		Limit:           limit,
		Offset:          offset,
	})
}

// Create adds a new transaction.
func (s *Service) Create(ctx context.Context, orgID, category, description, txnDate, txnType string, amount float64, paymentType, reference, entryBy string) (*Transaction, error) {
	if category == "" {
		return nil, fmt.Errorf("category is required")
	}

	if txnDate == "" {
		return nil, fmt.Errorf("transaction_date is required")
	}

	if _, err := time.Parse("2006-01-02", txnDate); err != nil {
		return nil, fmt.Errorf("invalid transaction_date format (expected YYYY-MM-DD)")
	}

	if !validTransactionTypes[txnType] {
		return nil, fmt.Errorf("invalid transaction_type: %s (must be income or expense)", txnType)
	}

	if amount <= 0 {
		return nil, fmt.Errorf("amount must be positive")
	}

	if !validPaymentTypes[paymentType] {
		return nil, fmt.Errorf("invalid payment_type: %s (must be cash, bank_transfer, or khalti)", paymentType)
	}

	txn, err := s.repo.Create(ctx, orgID, category, description, txnDate, txnType, amount, paymentType, reference, entryBy)
	if err != nil {
		return nil, fmt.Errorf("creating transaction: %w", err)
	}

	s.logger.Info("transaction created",
		"txn_id", txn.ID,
		"org_id", orgID,
		"type", txnType,
		"amount", amount,
		"entry_by", entryBy,
	)

	return txn, nil
}

// Update modifies an existing transaction.
func (s *Service) Update(ctx context.Context, orgID, txnID, category, description, txnDate, txnType string, amount float64, paymentType, reference string) (*Transaction, error) {
	if category == "" {
		return nil, fmt.Errorf("category is required")
	}

	if txnDate == "" {
		return nil, fmt.Errorf("transaction_date is required")
	}

	if _, err := time.Parse("2006-01-02", txnDate); err != nil {
		return nil, fmt.Errorf("invalid transaction_date format (expected YYYY-MM-DD)")
	}

	if !validTransactionTypes[txnType] {
		return nil, fmt.Errorf("invalid transaction_type: %s (must be income or expense)", txnType)
	}

	if amount <= 0 {
		return nil, fmt.Errorf("amount must be positive")
	}

	if !validPaymentTypes[paymentType] {
		return nil, fmt.Errorf("invalid payment_type: %s (must be cash, bank_transfer, or khalti)", paymentType)
	}

	txn, err := s.repo.Update(ctx, orgID, txnID, category, description, txnDate, txnType, amount, paymentType, reference)
	if err != nil {
		return nil, fmt.Errorf("updating transaction: %w", err)
	}

	s.logger.Info("transaction updated",
		"txn_id", txnID,
		"org_id", orgID,
	)

	return txn, nil
}

// Delete removes a transaction.
func (s *Service) Delete(ctx context.Context, orgID, txnID string) error {
	if err := s.repo.Delete(ctx, orgID, txnID); err != nil {
		return fmt.Errorf("deleting transaction: %w", err)
	}

	s.logger.Info("transaction deleted",
		"txn_id", txnID,
		"org_id", orgID,
	)

	return nil
}

// GetSummary retrieves financial summary for an organization.
func (s *Service) GetSummary(ctx context.Context, orgID, from, to string) (*Summary, error) {
	if err := validateDateRange(from, to); err != nil {
		return nil, err
	}

	return s.repo.GetSummary(ctx, orgID, from, to)
}

// Export retrieves all transactions matching filters for CSV export.
func (s *Service) Export(ctx context.Context, orgID, from, to, category, txnType, paymentType string) ([]Transaction, error) {
	if err := validateDateRange(from, to); err != nil {
		return nil, err
	}

	if txnType != "" && !validTransactionTypes[txnType] {
		return nil, fmt.Errorf("invalid transaction type: %s", txnType)
	}

	if paymentType != "" && !validPaymentTypes[paymentType] {
		return nil, fmt.Errorf("invalid payment type: %s", paymentType)
	}

	return s.repo.ListForExport(ctx, ListFilter{
		OrgID:           orgID,
		From:            from,
		To:              to,
		Category:        category,
		TransactionType: txnType,
		PaymentType:     paymentType,
	})
}

func validateDateRange(from, to string) error {
	if from != "" {
		if _, err := time.Parse("2006-01-02", from); err != nil {
			return fmt.Errorf("invalid from date format (expected YYYY-MM-DD)")
		}
	}
	if to != "" {
		if _, err := time.Parse("2006-01-02", to); err != nil {
			return fmt.Errorf("invalid to date format (expected YYYY-MM-DD)")
		}
	}
	return nil
}
