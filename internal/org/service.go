package org

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

// ErrNotFound signals that GetAuthenticated has nothing to return for this
// caller — either the org truly doesn't exist, or the caller has no business
// seeing it. Kept distinct from Update's "forbidden" errors on purpose: an
// outsider probing another gym's org ID must get a 404, not a 403 that
// confirms the org exists.
var ErrNotFound = errors.New("organization not found")

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

// GetAuthenticated retrieves the full organization record — including its tax
// identity — for the org's own admin (or a super_admin) viewing settings.
// Gated the same way Update is gated: role must be admin, or the caller must
// be a super_admin. See ErrNotFound for why a failed check surfaces as
// "not found" rather than "forbidden".
//
// The super_admin check runs first and independently of membership — nothing
// seeds an organization_members row for a super_admin, so a plain
// GetOrgRole-then-fallback (the order Update uses) would 404 a real
// super_admin before ever consulting IsSuperAdmin: this route sits outside
// OrgScope specifically so a super_admin can reach it, so that path must work.
func (s *Service) GetAuthenticated(ctx context.Context, userID, orgID string) (*Organization, error) {
	if isSA, err := s.repo.IsSuperAdmin(ctx, userID); err == nil && isSA {
		return s.repo.GetFullByID(ctx, orgID)
	}

	role, err := s.repo.GetOrgRole(ctx, userID, orgID)
	if err != nil || role != "admin" {
		return nil, ErrNotFound
	}

	return s.repo.GetFullByID(ctx, orgID)
}

// Update updates an existing organization. Caller must be an admin of the org.
//
// Unlike GetAuthenticated above, this checks membership before falling back
// to IsSuperAdmin — which means a super_admin with no organization_members
// row for orgID is rejected here too, same bug, just currently unreachable:
// this route sits inside OrgScope (cmd/api/main.go), which already requires
// membership before this method ever runs, so a bare super_admin never
// reaches Update today. If Update is ever moved outside OrgScope the way
// GetAuthenticated was, it needs the same independent-super_admin-check fix.
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

// CheckOrgMembership delegates to the repository to check membership and role.
// This implements middleware.OrgMembershipChecker.
func (s *Service) CheckOrgMembership(ctx context.Context, userID, orgID string) (bool, string, error) {
	return s.repo.CheckOrgMembership(ctx, userID, orgID)
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
