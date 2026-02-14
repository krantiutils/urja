package org

import (
	"context"
	"log/slog"
)

// Service handles organization business logic.
type Service struct {
	repo   *Repository
	logger *slog.Logger
}

// NewService creates a new org service.
func NewService(repo *Repository, logger *slog.Logger) *Service {
	return &Service{repo: repo, logger: logger}
}

// Get retrieves an organization by ID.
func (s *Service) Get(ctx context.Context, id string) (*Organization, error) {
	return s.repo.GetByID(ctx, id)
}

// List retrieves organizations with pagination.
func (s *Service) List(ctx context.Context, limit, offset int) ([]Organization, error) {
	if limit <= 0 || limit > 50 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	return s.repo.List(ctx, limit, offset)
}

// CheckOrgMembership delegates to the repository to check membership and role.
// This implements middleware.OrgMembershipChecker.
func (s *Service) CheckOrgMembership(ctx context.Context, userID, orgID string) (bool, string, error) {
	return s.repo.CheckOrgMembership(ctx, userID, orgID)
}
