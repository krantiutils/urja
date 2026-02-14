package feedback

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
)

// Service handles feedback business logic.
type Service struct {
	repo   *Repository
	logger *slog.Logger
}

// NewService creates a new feedback service.
func NewService(repo *Repository, logger *slog.Logger) *Service {
	return &Service{repo: repo, logger: logger}
}

// Create submits a new feedback.
func (s *Service) Create(ctx context.Context, orgID, userID, message string) (*Feedback, error) {
	message = strings.TrimSpace(message)
	if message == "" {
		return nil, fmt.Errorf("message is required")
	}

	f, err := s.repo.Create(ctx, orgID, userID, message)
	if err != nil {
		return nil, err
	}

	s.logger.Info("feedback submitted", "feedback_id", f.ID, "org_id", orgID, "user_id", userID)
	return f, nil
}

// List retrieves feedbacks with pagination.
func (s *Service) List(ctx context.Context, orgID string, limit, offset int) ([]Feedback, error) {
	if limit <= 0 || limit > 50 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	return s.repo.List(ctx, orgID, limit, offset)
}
