package attendance

import (
	"context"
	"fmt"
	"log/slog"
)

// Service handles attendance business logic.
type Service struct {
	repo   *Repository
	logger *slog.Logger
}

// NewService creates a new attendance service.
func NewService(repo *Repository, logger *slog.Logger) *Service {
	return &Service{repo: repo, logger: logger}
}

// CheckIn records a member check-in.
func (s *Service) CheckIn(ctx context.Context, userID, orgID, method string) (*Record, error) {
	validMethods := map[string]bool{"qr": true, "nfc": true, "manual": true}
	if !validMethods[method] {
		return nil, fmt.Errorf("invalid check-in method: %s", method)
	}

	rec, err := s.repo.CheckIn(ctx, userID, orgID, method)
	if err != nil {
		return nil, fmt.Errorf("check-in failed: %w", err)
	}

	s.logger.Info("member checked in",
		"user_id", userID,
		"org_id", orgID,
		"method", method,
	)

	return rec, nil
}

// ListByUser lists a user's attendance history.
func (s *Service) ListByUser(ctx context.Context, userID string, limit, offset int) ([]Record, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	return s.repo.ListByUser(ctx, userID, limit, offset)
}

// ListByOrg lists attendance for an organization.
func (s *Service) ListByOrg(ctx context.Context, orgID string, limit, offset int) ([]Record, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	return s.repo.ListByOrg(ctx, orgID, limit, offset)
}
