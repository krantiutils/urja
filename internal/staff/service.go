package staff

import (
	"context"
	"fmt"
	"log/slog"
	"regexp"
	"time"
)

var validStaffRoles = map[string]bool{
	"owner":        true,
	"manager":      true,
	"trainer":      true,
	"receptionist": true,
}

var validCheckInMethods = map[string]bool{
	"qr":     true,
	"nfc":    true,
	"manual": true,
}

var phoneRegex = regexp.MustCompile(`^9[78]\d{8}$`)

// Service handles staff business logic.
type Service struct {
	repo   *Repository
	logger *slog.Logger
}

// NewService creates a new staff service.
func NewService(repo *Repository, logger *slog.Logger) *Service {
	return &Service{repo: repo, logger: logger}
}

// ListByOrg lists staff members of an organization with optional search and pagination.
func (s *Service) ListByOrg(ctx context.Context, orgID, search string, limit, offset int) ([]StaffMember, int, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	return s.repo.ListByOrg(ctx, orgID, search, limit, offset)
}

// GetByID retrieves a single staff member by ID.
func (s *Service) GetByID(ctx context.Context, orgID, id string) (*StaffMember, error) {
	if id == "" {
		return nil, fmt.Errorf("staff member ID is required")
	}
	m, err := s.repo.GetByID(ctx, orgID, id)
	if err != nil {
		return nil, fmt.Errorf("getting staff member: %w", err)
	}
	return m, nil
}

// Create adds a new staff member to the organization.
func (s *Service) Create(ctx context.Context, orgID, phone, name, email, avatarURL, staffRole string) (*StaffMember, error) {
	phone = normalizePhone(phone)
	if !phoneRegex.MatchString(phone) {
		return nil, fmt.Errorf("invalid phone number: must be a 10-digit Nepali mobile number")
	}
	if name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if !validStaffRoles[staffRole] {
		return nil, fmt.Errorf("invalid staff role: must be one of owner, manager, trainer, receptionist")
	}

	m, err := s.repo.Create(ctx, orgID, phone, name, email, avatarURL, staffRole)
	if err != nil {
		return nil, fmt.Errorf("creating staff member: %w", err)
	}

	s.logger.Info("staff member created",
		"org_id", orgID,
		"staff_id", m.ID,
		"staff_role", staffRole,
	)
	return m, nil
}

// Update updates a staff member's details.
func (s *Service) Update(ctx context.Context, orgID, id string, name, email, avatarURL, staffRole *string) (*StaffMember, error) {
	if id == "" {
		return nil, fmt.Errorf("staff member ID is required")
	}
	if staffRole != nil && !validStaffRoles[*staffRole] {
		return nil, fmt.Errorf("invalid staff role: must be one of owner, manager, trainer, receptionist")
	}

	m, err := s.repo.Update(ctx, orgID, id, name, email, avatarURL, staffRole)
	if err != nil {
		return nil, fmt.Errorf("updating staff member: %w", err)
	}

	s.logger.Info("staff member updated",
		"org_id", orgID,
		"staff_id", id,
	)
	return m, nil
}

// Delete soft-deletes a staff member.
func (s *Service) Delete(ctx context.Context, orgID, id string) error {
	if id == "" {
		return fmt.Errorf("staff member ID is required")
	}

	if err := s.repo.Delete(ctx, orgID, id); err != nil {
		return fmt.Errorf("deleting staff member: %w", err)
	}

	s.logger.Info("staff member deleted",
		"org_id", orgID,
		"staff_id", id,
	)
	return nil
}

// ListAttendance lists staff attendance records with optional date range filtering.
func (s *Service) ListAttendance(ctx context.Context, orgID string, from, to *time.Time, limit, offset int) ([]AttendanceRecord, int, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	return s.repo.ListAttendance(ctx, orgID, from, to, limit, offset)
}

// RecordAttendance records a staff check-in or check-out.
func (s *Service) RecordAttendance(ctx context.Context, orgID, staffID, action, method string) (*AttendanceRecord, error) {
	if staffID == "" {
		return nil, fmt.Errorf("staff_id is required")
	}

	switch action {
	case "check_in":
		if !validCheckInMethods[method] {
			return nil, fmt.Errorf("invalid check-in method: must be one of qr, nfc, manual")
		}
		rec, err := s.repo.RecordCheckIn(ctx, orgID, staffID, method)
		if err != nil {
			return nil, fmt.Errorf("recording check-in: %w", err)
		}
		s.logger.Info("staff checked in",
			"org_id", orgID,
			"staff_id", staffID,
			"method", method,
		)
		return rec, nil

	case "check_out":
		rec, err := s.repo.RecordCheckOut(ctx, orgID, staffID)
		if err != nil {
			return nil, fmt.Errorf("recording check-out: %w", err)
		}
		s.logger.Info("staff checked out",
			"org_id", orgID,
			"staff_id", staffID,
		)
		return rec, nil

	default:
		return nil, fmt.Errorf("invalid action: must be check_in or check_out")
	}
}

// normalizePhone strips the +977 prefix if present and trims spaces.
func normalizePhone(phone string) string {
	if len(phone) > 3 && phone[:4] == "+977" {
		phone = phone[4:]
	} else if len(phone) > 2 && phone[:3] == "977" {
		phone = phone[3:]
	}
	return phone
}
