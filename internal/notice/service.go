package notice

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
)

// Service handles notice business logic.
type Service struct {
	repo   *Repository
	logger *slog.Logger
}

// NewService creates a new notice service.
func NewService(repo *Repository, logger *slog.Logger) *Service {
	return &Service{repo: repo, logger: logger}
}

// Create creates a new notice.
func (s *Service) Create(ctx context.Context, orgID, title, titleNe, content, contentNe, createdBy string) (*Notice, error) {
	title = strings.TrimSpace(title)
	content = strings.TrimSpace(content)

	if title == "" {
		return nil, fmt.Errorf("title is required")
	}
	if content == "" {
		return nil, fmt.Errorf("content is required")
	}

	n, err := s.repo.Create(ctx, orgID, title, strings.TrimSpace(titleNe), content, strings.TrimSpace(contentNe), createdBy)
	if err != nil {
		return nil, err
	}

	s.logger.Info("notice created", "notice_id", n.ID, "org_id", orgID, "created_by", createdBy)
	return n, nil
}

// GetByID retrieves a notice by ID.
func (s *Service) GetByID(ctx context.Context, orgID, noticeID string) (*Notice, error) {
	return s.repo.GetByID(ctx, orgID, noticeID)
}

// List retrieves notices with optional search and pagination.
func (s *Service) List(ctx context.Context, orgID, search string, limit, offset int) ([]Notice, error) {
	if limit <= 0 || limit > 50 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	return s.repo.List(ctx, orgID, strings.TrimSpace(search), limit, offset)
}

// Update modifies an existing notice.
func (s *Service) Update(ctx context.Context, orgID, noticeID, title, titleNe, content, contentNe string, isActive bool) (*Notice, error) {
	title = strings.TrimSpace(title)
	content = strings.TrimSpace(content)

	if title == "" {
		return nil, fmt.Errorf("title is required")
	}
	if content == "" {
		return nil, fmt.Errorf("content is required")
	}

	n, err := s.repo.Update(ctx, orgID, noticeID, title, strings.TrimSpace(titleNe), content, strings.TrimSpace(contentNe), isActive)
	if err != nil {
		return nil, err
	}

	s.logger.Info("notice updated", "notice_id", noticeID, "org_id", orgID)
	return n, nil
}

// Delete removes a notice.
func (s *Service) Delete(ctx context.Context, orgID, noticeID string) error {
	if err := s.repo.Delete(ctx, orgID, noticeID); err != nil {
		return err
	}
	s.logger.Info("notice deleted", "notice_id", noticeID, "org_id", orgID)
	return nil
}
