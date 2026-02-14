package org

import (
	"context"
	"fmt"
	"log/slog"
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"
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

// Create creates a new organization. Caller must be a super_admin.
func (s *Service) Create(ctx context.Context, userID string, in *CreateOrgInput) (*Organization, error) {
	isSA, err := s.repo.IsSuperAdmin(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("checking admin status: %w", err)
	}
	if !isSA {
		return nil, fmt.Errorf("forbidden: only super admins can create organizations")
	}

	if strings.TrimSpace(in.Name) == "" {
		return nil, fmt.Errorf("name is required")
	}

	// Generate slug if not provided.
	if in.Slug == "" {
		in.Slug = slugify(in.Name)
	}

	// Ensure slug is unique by appending a suffix if needed.
	baseSlug := in.Slug
	for suffix := 1; ; suffix++ {
		exists, err := s.repo.SlugExists(ctx, in.Slug)
		if err != nil {
			return nil, fmt.Errorf("checking slug: %w", err)
		}
		if !exists {
			break
		}
		in.Slug = fmt.Sprintf("%s-%d", baseSlug, suffix)
	}

	return s.repo.Create(ctx, in)
}

// Update updates an existing organization. Caller must be an admin of the org.
func (s *Service) Update(ctx context.Context, userID, orgID string, in *UpdateOrgInput) (*Organization, error) {
	role, err := s.repo.GetOrgRole(ctx, userID, orgID)
	if err != nil {
		return nil, fmt.Errorf("forbidden: not a member of this organization")
	}
	if role != "admin" {
		// Also allow super_admin to update any org.
		isSA, saErr := s.repo.IsSuperAdmin(ctx, userID)
		if saErr != nil || !isSA {
			return nil, fmt.Errorf("forbidden: only org admins can update organizations")
		}
	}

	return s.repo.Update(ctx, orgID, in)
}

// IsOrgMember delegates to the repository to check membership.
// This implements middleware.OrgMembershipChecker.
func (s *Service) IsOrgMember(ctx context.Context, userID, orgID string) (bool, error) {
	return s.repo.IsOrgMember(ctx, userID, orgID)
}

var nonAlphanumRe = regexp.MustCompile(`[^a-z0-9-]+`)
var multiDashRe = regexp.MustCompile(`-{2,}`)

// slugify converts a name to a URL-friendly slug.
func slugify(name string) string {
	// Normalize unicode and lowercase.
	s := strings.ToLower(norm.NFKD.String(name))
	// Keep only ASCII alphanumeric and hyphens.
	var b strings.Builder
	for _, r := range s {
		if r <= unicode.MaxASCII && (unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' || r == ' ') {
			b.WriteRune(r)
		}
	}
	s = strings.ReplaceAll(b.String(), " ", "-")
	s = nonAlphanumRe.ReplaceAllString(s, "")
	s = multiDashRe.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	if s == "" {
		s = "org"
	}
	return s
}
