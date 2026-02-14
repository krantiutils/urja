package dues

import (
	"context"
	"fmt"
	"log/slog"
)

// Service handles dues business logic.
type Service struct {
	repo   *Repository
	logger *slog.Logger
}

// NewService creates a new dues service.
func NewService(repo *Repository, logger *slog.Logger) *Service {
	return &Service{repo: repo, logger: logger}
}

var validPaymentMethods = map[string]bool{
	"cash":          true,
	"bank_transfer": true,
	"khalti":        true,
}

// List retrieves dues for an organization with optional filters.
func (s *Service) List(ctx context.Context, orgID, search, status string, limit, offset int) ([]Due, int, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	if status != "" {
		switch status {
		case "pending", "paid", "overdue", "waived":
			// valid
		default:
			return nil, 0, fmt.Errorf("invalid status filter: %s", status)
		}
	}

	return s.repo.List(ctx, ListFilter{
		OrgID:  orgID,
		Search: search,
		Status: status,
		Limit:  limit,
		Offset: offset,
	})
}

// Pay records a payment for a due.
func (s *Service) Pay(ctx context.Context, orgID, dueID, entryByUserID string, amount float64, method, reference string) (*PaymentRecord, error) {
	if amount <= 0 {
		return nil, fmt.Errorf("amount must be positive")
	}

	if !validPaymentMethods[method] {
		return nil, fmt.Errorf("invalid payment method: %s (must be cash, bank_transfer, or khalti)", method)
	}

	rec, err := s.repo.Pay(ctx, orgID, dueID, "", amount, method, reference, entryByUserID)
	if err != nil {
		return nil, fmt.Errorf("recording payment: %w", err)
	}

	s.logger.Info("due payment recorded",
		"due_id", dueID,
		"org_id", orgID,
		"amount", amount,
		"method", method,
		"entry_by", entryByUserID,
	)

	return rec, nil
}

// BlockAccess deactivates NFC door access for a member.
func (s *Service) BlockAccess(ctx context.Context, orgID, memberID string) (int64, error) {
	count, err := s.repo.BlockAccess(ctx, orgID, memberID)
	if err != nil {
		return 0, fmt.Errorf("blocking access: %w", err)
	}

	s.logger.Info("NFC access blocked",
		"member_id", memberID,
		"org_id", orgID,
		"cards_deactivated", count,
	)

	return count, nil
}
