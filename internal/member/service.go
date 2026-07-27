package member

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"regexp"
)

// ProfileUpdate holds optional fields for profile update.
type ProfileUpdate struct {
	Name                  *string `json:"name,omitempty"`
	NameNe                *string `json:"name_ne,omitempty"`
	Email                 *string `json:"email,omitempty"`
	DateOfBirth           *string `json:"date_of_birth,omitempty"`
	Gender                *string `json:"gender,omitempty"`
	AvatarURL             *string `json:"avatar_url,omitempty"`
	EmergencyContactName  *string `json:"emergency_contact_name,omitempty"`
	EmergencyContactPhone *string `json:"emergency_contact_phone,omitempty"`
}

// CreateMemberInput holds data for creating/adding a member to an org.
type CreateMemberInput struct {
	Phone  string `json:"phone"`
	Name   string `json:"name"`
	NameNe string `json:"name_ne,omitempty"`
	Email  string `json:"email,omitempty"`
	Role   string `json:"role,omitempty"`
}

// UpdateOrgMemberInput holds data for updating a member's org role/status.
type UpdateOrgMemberInput struct {
	Role   *string `json:"role,omitempty"`
	Status *string `json:"status,omitempty"`
}

// ErrForbidden indicates the caller does not have permission to perform the requested change.
var ErrForbidden = errors.New("forbidden")

var (
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	phoneRegex = regexp.MustCompile(`^[9876]\d{9}$`)
	validRoles = map[string]struct{}{
		"member": {},
		"staff":  {},
		"admin":  {},
	}
	validStatuses = map[string]struct{}{
		"active":    {},
		"suspended": {},
		"left":      {},
	}
	validGenders = map[string]struct{}{
		"male":           {},
		"female":         {},
		"other":          {},
		"prefer_not_say": {},
	}
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

// GetProfile retrieves the authenticated user's profile including org memberships.
func (s *Service) GetProfile(ctx context.Context, userID string) (*Member, error) {
	m, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("getting profile: %w", err)
	}

	orgs, err := s.repo.GetOrgMemberships(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("getting org memberships: %w", err)
	}
	if orgs == nil {
		orgs = []OrgMembership{}
	}
	m.Organizations = orgs

	return m, nil
}

// UpdateProfile updates the authenticated user's profile.
func (s *Service) UpdateProfile(ctx context.Context, userID string, upd *ProfileUpdate) error {
	if err := validateProfileUpdate(upd); err != nil {
		return err
	}

	if err := s.repo.UpdateProfile(ctx, userID, upd); err != nil {
		return fmt.Errorf("updating profile: %w", err)
	}
	return nil
}

// UpdatePrivacy updates the authenticated user's privacy settings.
func (s *Service) UpdatePrivacy(ctx context.Context, userID string, settings PrivacySettings) error {
	if err := s.repo.UpdatePrivacy(ctx, userID, settings); err != nil {
		return fmt.Errorf("updating privacy: %w", err)
	}
	return nil
}

// ListByOrg lists members of an organization with pagination.
func (s *Service) ListByOrg(ctx context.Context, orgID string, limit, offset int) ([]OrgMember, int, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	return s.repo.ListByOrg(ctx, orgID, limit, offset)
}

// CreateOrgMember adds a member to an organization.
func (s *Service) CreateOrgMember(ctx context.Context, orgID string, input *CreateMemberInput) (*OrgMember, error) {
	if err := validateCreateMemberInput(input); err != nil {
		return nil, err
	}

	m, err := s.repo.CreateOrgMember(ctx, orgID, input)
	if err != nil {
		return nil, fmt.Errorf("creating org member: %w", err)
	}
	return m, nil
}

// UpdateOrgMember updates a member's role or status within an org.
// callerID and callerRole identify the authenticated caller and their role within
// orgID (as resolved by OrgScope/RequireOrgRole for this request) — never trust a
// role embedded in the JWT itself. Changing the role field requires the caller to
// be an admin of this org, and a caller may never change their own role, even if
// they are an admin themselves.
func (s *Service) UpdateOrgMember(ctx context.Context, callerID, callerRole, orgID, userID string, upd *UpdateOrgMemberInput) error {
	if err := validateUpdateOrgMemberInput(upd); err != nil {
		return err
	}

	if upd.Role != nil {
		if callerRole != "admin" {
			return fmt.Errorf("%w: only admins can change a member's role", ErrForbidden)
		}
		if callerID == userID {
			return fmt.Errorf("%w: cannot change your own role", ErrForbidden)
		}
	}

	if err := s.repo.UpdateOrgMember(ctx, orgID, userID, upd); err != nil {
		return fmt.Errorf("updating org member: %w", err)
	}
	return nil
}

// DeleteOrgMember removes a member from an organization (soft delete).
func (s *Service) DeleteOrgMember(ctx context.Context, orgID, userID string) error {
	if err := s.repo.DeleteOrgMember(ctx, orgID, userID); err != nil {
		return fmt.Errorf("deleting org member: %w", err)
	}
	return nil
}

// DeleteAccount deactivates the user's account and anonymizes their personal data.
func (s *Service) DeleteAccount(ctx context.Context, userID string) error {
	if err := s.repo.DeleteAccount(ctx, userID); err != nil {
		return fmt.Errorf("deleting account: %w", err)
	}
	return nil
}

// GetPrimaryOrgID returns the user's primary (earliest active) organization ID.
func (s *Service) GetPrimaryOrgID(ctx context.Context, userID string) (string, error) {
	return s.repo.GetPrimaryOrgID(ctx, userID)
}

func validateProfileUpdate(upd *ProfileUpdate) error {
	if upd.Name != nil && *upd.Name == "" {
		return errors.New("name cannot be empty")
	}
	if upd.Email != nil && *upd.Email != "" && !emailRegex.MatchString(*upd.Email) {
		return errors.New("invalid email format")
	}
	if upd.Gender != nil && *upd.Gender != "" {
		if _, ok := validGenders[*upd.Gender]; !ok {
			return errors.New("invalid gender: must be male, female, other, or prefer_not_say")
		}
	}
	if upd.EmergencyContactPhone != nil && *upd.EmergencyContactPhone != "" {
		if !phoneRegex.MatchString(*upd.EmergencyContactPhone) {
			return errors.New("invalid emergency contact phone number")
		}
	}
	return nil
}

func validateCreateMemberInput(input *CreateMemberInput) error {
	if !phoneRegex.MatchString(input.Phone) {
		return errors.New("invalid phone number: must be 10 digits starting with 9, 8, 7, or 6")
	}
	if input.Name == "" {
		return errors.New("name is required")
	}
	if input.Email != "" && !emailRegex.MatchString(input.Email) {
		return errors.New("invalid email format")
	}
	if input.Role != "" {
		if _, ok := validRoles[input.Role]; !ok {
			return errors.New("invalid role: must be member, staff, or admin")
		}
	}
	return nil
}

func validateUpdateOrgMemberInput(upd *UpdateOrgMemberInput) error {
	if upd.Role == nil && upd.Status == nil {
		return errors.New("at least one of role or status must be provided")
	}
	if upd.Role != nil {
		if _, ok := validRoles[*upd.Role]; !ok {
			return errors.New("invalid role: must be member, staff, or admin")
		}
	}
	if upd.Status != nil {
		if _, ok := validStatuses[*upd.Status]; !ok {
			return errors.New("invalid status: must be active, suspended, or left")
		}
	}
	return nil
}
