package guide

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
)

// Valid categories for training guides.
var validCategories = map[string]bool{
	"strength":    true,
	"cardio":      true,
	"flexibility": true,
	"nutrition":   true,
	"recovery":    true,
	"beginner":    true,
}

// Service handles training guide business logic.
type Service struct {
	repo   *Repository
	logger *slog.Logger
}

// NewService creates a new guide service.
func NewService(repo *Repository, logger *slog.Logger) *Service {
	return &Service{repo: repo, logger: logger}
}

// Create creates a new training guide.
func (s *Service) Create(ctx context.Context, orgID, title, titleNe, content, contentNe, category, coverImageURL, authorID string) (*TrainingGuide, error) {
	title = strings.TrimSpace(title)
	content = strings.TrimSpace(content)
	category = strings.TrimSpace(strings.ToLower(category))

	if title == "" {
		return nil, fmt.Errorf("title is required")
	}
	if content == "" {
		return nil, fmt.Errorf("content is required")
	}
	if category != "" && !validCategories[category] {
		return nil, fmt.Errorf("invalid category: %s", category)
	}

	g, err := s.repo.Create(ctx, orgID, title, strings.TrimSpace(titleNe), content, strings.TrimSpace(contentNe),
		category, strings.TrimSpace(coverImageURL), authorID)
	if err != nil {
		return nil, err
	}

	s.logger.Info("training guide created", "guide_id", g.ID, "author_id", authorID)
	return g, nil
}

// GetByID retrieves a training guide by ID (admin view, includes unpublished).
func (s *Service) GetByID(ctx context.Context, orgID, guideID string) (*TrainingGuide, error) {
	return s.repo.GetByID(ctx, orgID, guideID)
}

// GetPublishedByID retrieves a published training guide by ID (public view).
func (s *Service) GetPublishedByID(ctx context.Context, guideID string) (*TrainingGuide, error) {
	return s.repo.GetPublishedByID(ctx, guideID)
}

// ListPublished retrieves published training guides with optional filters.
// ErrOrgRequired is returned when a public listing does not name a gym.
var ErrGymRequired = errors.New("gym is required")

func (s *Service) ListPublished(ctx context.Context, orgSlug, category, search string, limit, offset int) ([]TrainingGuide, int, error) {
	if orgSlug == "" {
		return nil, 0, ErrGymRequired
	}

	if limit <= 0 || limit > 50 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	return s.repo.ListPublished(ctx, orgSlug, strings.TrimSpace(strings.ToLower(category)), strings.TrimSpace(search), limit, offset)
}

// ListAll retrieves all training guides (admin view).
func (s *Service) ListAll(ctx context.Context, orgID, category, search string, limit, offset int) ([]TrainingGuide, int, error) {
	if limit <= 0 || limit > 50 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	return s.repo.ListAll(ctx, orgID, strings.TrimSpace(strings.ToLower(category)), strings.TrimSpace(search), limit, offset)
}

// Update modifies an existing training guide.
func (s *Service) Update(ctx context.Context, orgID, guideID, title, titleNe, content, contentNe, category, coverImageURL string) (*TrainingGuide, error) {
	title = strings.TrimSpace(title)
	content = strings.TrimSpace(content)
	category = strings.TrimSpace(strings.ToLower(category))

	if title == "" {
		return nil, fmt.Errorf("title is required")
	}
	if content == "" {
		return nil, fmt.Errorf("content is required")
	}
	if category != "" && !validCategories[category] {
		return nil, fmt.Errorf("invalid category: %s", category)
	}

	g, err := s.repo.Update(ctx, orgID, guideID, title, strings.TrimSpace(titleNe), content, strings.TrimSpace(contentNe),
		category, strings.TrimSpace(coverImageURL))
	if err != nil {
		return nil, err
	}

	s.logger.Info("training guide updated", "guide_id", guideID)
	return g, nil
}

// SetPublished updates the publish status of a training guide.
func (s *Service) SetPublished(ctx context.Context, orgID, guideID string, isPublished bool) (*TrainingGuide, error) {
	g, err := s.repo.SetPublished(ctx, orgID, guideID, isPublished)
	if err != nil {
		return nil, err
	}

	action := "unpublished"
	if isPublished {
		action = "published"
	}
	s.logger.Info("training guide "+action, "guide_id", guideID)
	return g, nil
}

// Delete removes a training guide.
func (s *Service) Delete(ctx context.Context, orgID, guideID string) error {
	if err := s.repo.Delete(ctx, orgID, guideID); err != nil {
		return err
	}
	s.logger.Info("training guide deleted", "guide_id", guideID)
	return nil
}
