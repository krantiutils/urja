package member

import (
	"context"
	"fmt"
	"log/slog"
)

// Service handles member business logic.
type Service struct {
	repo   *Repository
	logger *slog.Logger
}

// NewService creates a new member service.
func NewService(repo *Repository, logger *slog.Logger) *Service {
	return &Service{repo: repo, logger: logger}
}

// GetProfile retrieves the authenticated user's profile.
func (s *Service) GetProfile(ctx context.Context, userID string) (*Member, error) {
	m, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("getting profile: %w", err)
	}
	return m, nil
}

// UpdateProfile updates the authenticated user's profile.
func (s *Service) UpdateProfile(ctx context.Context, userID string, name, email, avatarURL *string) error {
	if err := s.repo.Update(ctx, userID, name, email, avatarURL); err != nil {
		return fmt.Errorf("updating profile: %w", err)
	}
	return nil
}

// ListByOrg lists members of an organization with pagination.
func (s *Service) ListByOrg(ctx context.Context, orgID string, limit, offset int) ([]Member, int, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	return s.repo.ListByOrg(ctx, orgID, limit, offset)
}
