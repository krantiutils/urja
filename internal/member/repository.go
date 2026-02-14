package member

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ErrMemberNotFound is returned when a member lookup finds no result.
var ErrMemberNotFound = errors.New("member not found")

// ErrAlreadyOrgMember is returned when trying to add a user who is already an active member.
var ErrAlreadyOrgMember = errors.New("user is already a member of this organization")

// PrivacySettings controls profile visibility preferences.
type PrivacySettings struct {
	ShowEmail         bool `json:"show_email"`
	ShowPhone         bool `json:"show_phone"`
	ShowProfile       bool `json:"show_profile"`
	ShowAttendance    bool `json:"show_attendance"`
	ShowOnLeaderboard bool `json:"show_on_leaderboard"`
}

// OrgMembership represents a user's membership in a single organization.
type OrgMembership struct {
	OrgID    string    `json:"org_id"`
	OrgName  string    `json:"org_name"`
	Role     string    `json:"role"`
	Status   string    `json:"status"`
	JoinedAt time.Time `json:"joined_at"`
}

// Member represents a user profile with full fields.
type Member struct {
	ID                    string          `json:"id"`
	Phone                 string          `json:"phone"`
	Name                  string          `json:"name"`
	NameNe                string          `json:"name_ne,omitempty"`
	Email                 string          `json:"email,omitempty"`
	DateOfBirth           *string         `json:"date_of_birth,omitempty"`
	Gender                string          `json:"gender,omitempty"`
	AvatarURL             string          `json:"avatar_url,omitempty"`
	EmergencyContactName  string          `json:"emergency_contact_name,omitempty"`
	EmergencyContactPhone string          `json:"emergency_contact_phone,omitempty"`
	PrivacySettings       PrivacySettings `json:"privacy_settings"`
	Organizations         []OrgMembership `json:"organizations,omitempty"`
	CreatedAt             time.Time       `json:"created_at"`
	UpdatedAt             time.Time       `json:"updated_at"`
}

// OrgMember represents a member as seen within an organization context.
type OrgMember struct {
	ID        string    `json:"id"`
	Phone     string    `json:"phone"`
	Name      string    `json:"name"`
	NameNe    string    `json:"name_ne,omitempty"`
	Email     string    `json:"email,omitempty"`
	AvatarURL string    `json:"avatar_url,omitempty"`
	Role      string    `json:"role"`
	Status    string    `json:"status"`
	JoinedAt  time.Time `json:"joined_at"`
}

// Repository handles member data persistence.
type Repository struct {
	db *pgxpool.Pool
}

// NewRepository creates a new member repository.
func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

// GetByID retrieves a member's full profile by their user ID.
// SECURITY: Never returns sensitive fields (fcm_token, etc.)
func (r *Repository) GetByID(ctx context.Context, id string) (*Member, error) {
	var m Member
	var privJSON []byte
	var dob *string
	err := r.db.QueryRow(ctx,
		`SELECT id, phone, name, COALESCE(name_ne, ''), COALESCE(email, ''),
		        date_of_birth::text, COALESCE(gender, ''), COALESCE(avatar_url, ''),
		        COALESCE(emergency_contact_name, ''), COALESCE(emergency_contact_phone, ''),
		        privacy_settings, created_at, updated_at
		 FROM users WHERE id = $1 AND is_active = true`,
		id,
	).Scan(&m.ID, &m.Phone, &m.Name, &m.NameNe, &m.Email,
		&dob, &m.Gender, &m.AvatarURL,
		&m.EmergencyContactName, &m.EmergencyContactPhone,
		&privJSON, &m.CreatedAt, &m.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrMemberNotFound
		}
		return nil, fmt.Errorf("querying member: %w", err)
	}

	if dob != nil && *dob != "" {
		m.DateOfBirth = dob
	}

	if err := json.Unmarshal(privJSON, &m.PrivacySettings); err != nil {
		return nil, fmt.Errorf("parsing privacy settings: %w", err)
	}

	return &m, nil
}

// GetOrgMemberships retrieves all organizations a user belongs to.
func (r *Repository) GetOrgMemberships(ctx context.Context, userID string) ([]OrgMembership, error) {
	rows, err := r.db.Query(ctx,
		`SELECT om.organization_id, o.name, om.role, om.status, om.joined_at
		 FROM organization_members om
		 JOIN organizations o ON o.id = om.organization_id
		 WHERE om.user_id = $1 AND om.status = 'active'
		 ORDER BY om.joined_at DESC`,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("querying org memberships: %w", err)
	}
	defer rows.Close()

	var memberships []OrgMembership
	for rows.Next() {
		var om OrgMembership
		if err := rows.Scan(&om.OrgID, &om.OrgName, &om.Role, &om.Status, &om.JoinedAt); err != nil {
			return nil, fmt.Errorf("scanning org membership: %w", err)
		}
		memberships = append(memberships, om)
	}
	return memberships, rows.Err()
}

// UpdateProfile updates a member's profile fields. Only non-nil fields are updated.
func (r *Repository) UpdateProfile(ctx context.Context, id string, upd *ProfileUpdate) error {
	tag, err := r.db.Exec(ctx,
		`UPDATE users SET
			name = COALESCE($2, name),
			name_ne = COALESCE($3, name_ne),
			email = COALESCE($4, email),
			date_of_birth = COALESCE($5, date_of_birth),
			gender = COALESCE($6, gender),
			avatar_url = COALESCE($7, avatar_url),
			emergency_contact_name = COALESCE($8, emergency_contact_name),
			emergency_contact_phone = COALESCE($9, emergency_contact_phone)
		 WHERE id = $1 AND is_active = true`,
		id, upd.Name, upd.NameNe, upd.Email, upd.DateOfBirth,
		upd.Gender, upd.AvatarURL, upd.EmergencyContactName, upd.EmergencyContactPhone,
	)
	if err != nil {
		return fmt.Errorf("updating profile: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrMemberNotFound
	}
	return nil
}

// UpdatePrivacy updates a member's privacy settings.
func (r *Repository) UpdatePrivacy(ctx context.Context, id string, settings PrivacySettings) error {
	privJSON, err := json.Marshal(settings)
	if err != nil {
		return fmt.Errorf("marshaling privacy settings: %w", err)
	}

	tag, err := r.db.Exec(ctx,
		`UPDATE users SET privacy_settings = $2 WHERE id = $1 AND is_active = true`,
		id, privJSON,
	)
	if err != nil {
		return fmt.Errorf("updating privacy settings: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrMemberNotFound
	}
	return nil
}

// ListByOrg retrieves members of an organization with pagination.
// Returns org-scoped member view including role and status.
func (r *Repository) ListByOrg(ctx context.Context, orgID string, limit, offset int) ([]OrgMember, int, error) {
	var total int
	err := r.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM organization_members om
		 JOIN users u ON u.id = om.user_id
		 WHERE om.organization_id = $1 AND om.status = 'active' AND u.is_active = true`,
		orgID,
	).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("counting members: %w", err)
	}

	rows, err := r.db.Query(ctx,
		`SELECT u.id, u.phone, u.name, COALESCE(u.name_ne, ''),
		        COALESCE(u.email, ''), COALESCE(u.avatar_url, ''),
		        om.role, om.status, om.joined_at
		 FROM users u
		 JOIN organization_members om ON u.id = om.user_id
		 WHERE om.organization_id = $1 AND om.status = 'active' AND u.is_active = true
		 ORDER BY om.joined_at DESC
		 LIMIT $2 OFFSET $3`,
		orgID, limit, offset,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("querying members: %w", err)
	}
	defer rows.Close()

	var members []OrgMember
	for rows.Next() {
		var m OrgMember
		if err := rows.Scan(&m.ID, &m.Phone, &m.Name, &m.NameNe,
			&m.Email, &m.AvatarURL,
			&m.Role, &m.Status, &m.JoinedAt); err != nil {
			return nil, 0, fmt.Errorf("scanning member: %w", err)
		}
		members = append(members, m)
	}

	return members, total, rows.Err()
}

// CreateOrgMember adds a user to an organization. If the user doesn't exist (by phone),
// a new user record is created first.
func (r *Repository) CreateOrgMember(ctx context.Context, orgID string, input *CreateMemberInput) (*OrgMember, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("beginning transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Find or create user by phone
	var userID string
	err = tx.QueryRow(ctx,
		`SELECT id FROM users WHERE phone = $1`,
		input.Phone,
	).Scan(&userID)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("looking up user: %w", err)
		}
		// Create new user
		err = tx.QueryRow(ctx,
			`INSERT INTO users (phone, name, name_ne, email)
			 VALUES ($1, $2, $3, $4)
			 RETURNING id`,
			input.Phone, input.Name, nilIfEmpty(input.NameNe), nilIfEmpty(input.Email),
		).Scan(&userID)
		if err != nil {
			return nil, fmt.Errorf("creating user: %w", err)
		}
	}

	// Check for existing active membership
	var existingStatus string
	err = tx.QueryRow(ctx,
		`SELECT status FROM organization_members
		 WHERE user_id = $1 AND organization_id = $2`,
		userID, orgID,
	).Scan(&existingStatus)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("checking existing membership: %w", err)
	}

	role := input.Role
	if role == "" {
		role = "member"
	}

	if err == nil {
		// Membership row exists
		if existingStatus == "active" {
			return nil, ErrAlreadyOrgMember
		}
		// Re-activate a left/suspended membership
		_, err = tx.Exec(ctx,
			`UPDATE organization_members SET status = 'active', role = $3, joined_at = NOW()
			 WHERE user_id = $1 AND organization_id = $2`,
			userID, orgID, role,
		)
		if err != nil {
			return nil, fmt.Errorf("reactivating membership: %w", err)
		}
	} else {
		// Create new membership
		_, err = tx.Exec(ctx,
			`INSERT INTO organization_members (user_id, organization_id, role)
			 VALUES ($1, $2, $3)`,
			userID, orgID, role,
		)
		if err != nil {
			return nil, fmt.Errorf("creating org membership: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("committing transaction: %w", err)
	}

	// Fetch the created/updated org member
	var m OrgMember
	err = r.db.QueryRow(ctx,
		`SELECT u.id, u.phone, u.name, COALESCE(u.name_ne, ''),
		        COALESCE(u.email, ''), COALESCE(u.avatar_url, ''),
		        om.role, om.status, om.joined_at
		 FROM users u
		 JOIN organization_members om ON u.id = om.user_id
		 WHERE u.id = $1 AND om.organization_id = $2`,
		userID, orgID,
	).Scan(&m.ID, &m.Phone, &m.Name, &m.NameNe,
		&m.Email, &m.AvatarURL,
		&m.Role, &m.Status, &m.JoinedAt)
	if err != nil {
		return nil, fmt.Errorf("fetching created member: %w", err)
	}

	return &m, nil
}

// UpdateOrgMember updates a member's role or status within an organization.
func (r *Repository) UpdateOrgMember(ctx context.Context, orgID, userID string, upd *UpdateOrgMemberInput) error {
	tag, err := r.db.Exec(ctx,
		`UPDATE organization_members SET
			role = COALESCE($3, role),
			status = COALESCE($4, status)
		 WHERE user_id = $1 AND organization_id = $2 AND status != 'left'`,
		userID, orgID, upd.Role, upd.Status,
	)
	if err != nil {
		return fmt.Errorf("updating org member: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrMemberNotFound
	}
	return nil
}

// DeleteOrgMember soft-deletes a member from an organization by setting status to 'left'.
func (r *Repository) DeleteOrgMember(ctx context.Context, orgID, userID string) error {
	tag, err := r.db.Exec(ctx,
		`UPDATE organization_members SET status = 'left'
		 WHERE user_id = $1 AND organization_id = $2 AND status = 'active'`,
		userID, orgID,
	)
	if err != nil {
		return fmt.Errorf("deleting org member: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrMemberNotFound
	}
	return nil
}

func nilIfEmpty(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
