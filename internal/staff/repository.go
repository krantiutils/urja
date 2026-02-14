package staff

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// StaffMember represents a staff member within an organization.
type StaffMember struct {
	ID        string    `json:"id"`         // organization_members.id
	UserID    string    `json:"user_id"`    // users.id
	Name      string    `json:"name"`       // users.name
	Phone     string    `json:"phone"`      // users.phone
	Email     string    `json:"email,omitempty"`
	AvatarURL string    `json:"avatar_url,omitempty"`
	Role      string    `json:"role"`       // organization_members.role (staff/admin)
	StaffRole string    `json:"staff_role"` // owner/manager/trainer/receptionist
	Status    string    `json:"status"`
	JoinedAt  time.Time `json:"joined_at"`
}

// AttendanceRecord represents a staff attendance entry with computed hours.
type AttendanceRecord struct {
	ID         string     `json:"id"`
	StaffID    string     `json:"staff_id"`
	StaffName  string     `json:"staff_name"`
	CheckInAt  time.Time  `json:"check_in_at"`
	CheckOutAt *time.Time `json:"check_out_at,omitempty"`
	Method     string     `json:"method"`
	HoursSpent *float64   `json:"hours_spent,omitempty"`
}

// Repository handles staff data persistence.
type Repository struct {
	db *pgxpool.Pool
}

// NewRepository creates a new staff repository.
func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

const staffSelectColumns = `om.id, u.id, COALESCE(u.name, ''), u.phone, COALESCE(u.email, ''),
	COALESCE(u.avatar_url, ''), om.role, COALESCE(om.staff_role, ''),
	om.status, om.joined_at`

const staffBaseWhere = `om.organization_id = $1 AND om.role IN ('staff', 'admin') AND om.status = 'active'`

// ListByOrg retrieves staff members of an organization with optional search and pagination.
func (r *Repository) ListByOrg(ctx context.Context, orgID, search string, limit, offset int) ([]StaffMember, int, error) {
	var total int

	if search != "" {
		searchPattern := "%" + search + "%"
		err := r.db.QueryRow(ctx,
			`SELECT COUNT(*)
			 FROM organization_members om
			 JOIN users u ON u.id = om.user_id
			 WHERE `+staffBaseWhere+`
			   AND (u.name ILIKE $2 OR u.phone ILIKE $2)`,
			orgID, searchPattern,
		).Scan(&total)
		if err != nil {
			return nil, 0, fmt.Errorf("counting staff: %w", err)
		}

		rows, err := r.db.Query(ctx,
			`SELECT `+staffSelectColumns+`
			 FROM organization_members om
			 JOIN users u ON u.id = om.user_id
			 WHERE `+staffBaseWhere+`
			   AND (u.name ILIKE $4 OR u.phone ILIKE $4)
			 ORDER BY om.joined_at DESC
			 LIMIT $2 OFFSET $3`,
			orgID, limit, offset, searchPattern,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("querying staff: %w", err)
		}
		return scanStaffRows(rows, total)
	}

	err := r.db.QueryRow(ctx,
		`SELECT COUNT(*)
		 FROM organization_members om
		 JOIN users u ON u.id = om.user_id
		 WHERE `+staffBaseWhere,
		orgID,
	).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("counting staff: %w", err)
	}

	rows, err := r.db.Query(ctx,
		`SELECT `+staffSelectColumns+`
		 FROM organization_members om
		 JOIN users u ON u.id = om.user_id
		 WHERE `+staffBaseWhere+`
		 ORDER BY om.joined_at DESC
		 LIMIT $2 OFFSET $3`,
		orgID, limit, offset,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("querying staff: %w", err)
	}
	return scanStaffRows(rows, total)
}

func scanStaffRows(rows pgx.Rows, total int) ([]StaffMember, int, error) {
	defer rows.Close()

	var members []StaffMember
	for rows.Next() {
		var m StaffMember
		if err := rows.Scan(&m.ID, &m.UserID, &m.Name, &m.Phone, &m.Email,
			&m.AvatarURL, &m.Role, &m.StaffRole, &m.Status, &m.JoinedAt); err != nil {
			return nil, 0, fmt.Errorf("scanning staff member: %w", err)
		}
		members = append(members, m)
	}
	return members, total, rows.Err()
}

// GetByID retrieves a single staff member by their organization_members ID.
func (r *Repository) GetByID(ctx context.Context, orgID, id string) (*StaffMember, error) {
	var m StaffMember
	err := r.db.QueryRow(ctx,
		`SELECT `+staffSelectColumns+`
		 FROM organization_members om
		 JOIN users u ON u.id = om.user_id
		 WHERE om.id = $1 AND om.organization_id = $2
		   AND om.role IN ('staff', 'admin') AND om.status = 'active'`,
		id, orgID,
	).Scan(&m.ID, &m.UserID, &m.Name, &m.Phone, &m.Email,
		&m.AvatarURL, &m.Role, &m.StaffRole, &m.Status, &m.JoinedAt)
	if err != nil {
		return nil, fmt.Errorf("staff member not found: %w", err)
	}
	return &m, nil
}

// Create adds a new staff member to an organization.
// It finds or creates the user by phone, then creates/updates the org membership.
func (r *Repository) Create(ctx context.Context, orgID, phone, name, email, avatarURL, staffRole string) (*StaffMember, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("beginning transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Find or create user by phone
	var userID string
	err = tx.QueryRow(ctx,
		`SELECT id FROM users WHERE phone = $1`, phone,
	).Scan(&userID)
	if err != nil {
		err = tx.QueryRow(ctx,
			`INSERT INTO users (phone, name, email, avatar_url)
			 VALUES ($1, $2, NULLIF($3, ''), NULLIF($4, ''))
			 RETURNING id`,
			phone, name, email, avatarURL,
		).Scan(&userID)
		if err != nil {
			return nil, fmt.Errorf("creating user: %w", err)
		}
	}

	authRole := "staff"
	if staffRole == "owner" {
		authRole = "admin"
	}

	// Upsert organization membership
	var memberID string
	err = tx.QueryRow(ctx,
		`INSERT INTO organization_members (user_id, organization_id, role, staff_role, status)
		 VALUES ($1, $2, $3, $4, 'active')
		 ON CONFLICT (user_id, organization_id) DO UPDATE
		   SET role = EXCLUDED.role, staff_role = EXCLUDED.staff_role, status = 'active'
		 RETURNING id`,
		userID, orgID, authRole, staffRole,
	).Scan(&memberID)
	if err != nil {
		return nil, fmt.Errorf("creating staff membership: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("committing transaction: %w", err)
	}

	return r.GetByID(ctx, orgID, memberID)
}

// Update updates a staff member's profile and role.
func (r *Repository) Update(ctx context.Context, orgID, id string, name, email, avatarURL, staffRole *string) (*StaffMember, error) {
	var userID string
	err := r.db.QueryRow(ctx,
		`SELECT user_id FROM organization_members
		 WHERE id = $1 AND organization_id = $2 AND role IN ('staff', 'admin') AND status = 'active'`,
		id, orgID,
	).Scan(&userID)
	if err != nil {
		return nil, fmt.Errorf("staff member not found: %w", err)
	}

	_, err = r.db.Exec(ctx,
		`UPDATE users SET
			name = COALESCE($2, name),
			email = COALESCE($3, email),
			avatar_url = COALESCE($4, avatar_url),
			updated_at = NOW()
		 WHERE id = $1`,
		userID, name, email, avatarURL,
	)
	if err != nil {
		return nil, fmt.Errorf("updating user: %w", err)
	}

	if staffRole != nil {
		authRole := "staff"
		if *staffRole == "owner" {
			authRole = "admin"
		}
		_, err = r.db.Exec(ctx,
			`UPDATE organization_members SET role = $1, staff_role = $2 WHERE id = $3`,
			authRole, *staffRole, id,
		)
		if err != nil {
			return nil, fmt.Errorf("updating staff role: %w", err)
		}
	}

	return r.GetByID(ctx, orgID, id)
}

// Delete soft-deletes a staff member by setting their status to 'left'.
func (r *Repository) Delete(ctx context.Context, orgID, id string) error {
	result, err := r.db.Exec(ctx,
		`UPDATE organization_members SET status = 'left'
		 WHERE id = $1 AND organization_id = $2 AND role IN ('staff', 'admin') AND status = 'active'`,
		id, orgID,
	)
	if err != nil {
		return fmt.Errorf("deleting staff member: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("staff member not found")
	}
	return nil
}

// ListAttendance retrieves staff attendance records for an organization.
func (r *Repository) ListAttendance(ctx context.Context, orgID string, from, to *time.Time, limit, offset int) ([]AttendanceRecord, int, error) {
	countSQL := `SELECT COUNT(*)
		FROM attendance a
		JOIN organization_members om ON om.user_id = a.user_id AND om.organization_id = a.organization_id
		WHERE a.organization_id = $1
		  AND om.role IN ('staff', 'admin') AND om.status = 'active'`
	listSQL := `SELECT a.id, om.id, COALESCE(u.name, ''),
		a.check_in_at, a.check_out_at, a.check_in_method
		FROM attendance a
		JOIN organization_members om ON om.user_id = a.user_id AND om.organization_id = a.organization_id
		JOIN users u ON u.id = a.user_id
		WHERE a.organization_id = $1
		  AND om.role IN ('staff', 'admin') AND om.status = 'active'`

	args := []interface{}{orgID}
	argIdx := 2

	if from != nil {
		filter := fmt.Sprintf(` AND a.check_in_at >= $%d`, argIdx)
		countSQL += filter
		listSQL += filter
		args = append(args, *from)
		argIdx++
	}
	if to != nil {
		filter := fmt.Sprintf(` AND a.check_in_at <= $%d`, argIdx)
		countSQL += filter
		listSQL += filter
		args = append(args, *to)
		argIdx++
	}

	var total int
	if err := r.db.QueryRow(ctx, countSQL, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("counting staff attendance: %w", err)
	}

	listSQL += fmt.Sprintf(` ORDER BY a.check_in_at DESC LIMIT $%d OFFSET $%d`, argIdx, argIdx+1)
	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, listSQL, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("querying staff attendance: %w", err)
	}
	defer rows.Close()

	var records []AttendanceRecord
	for rows.Next() {
		var rec AttendanceRecord
		if err := rows.Scan(&rec.ID, &rec.StaffID, &rec.StaffName,
			&rec.CheckInAt, &rec.CheckOutAt, &rec.Method); err != nil {
			return nil, 0, fmt.Errorf("scanning staff attendance: %w", err)
		}
		if rec.CheckOutAt != nil {
			hours := rec.CheckOutAt.Sub(rec.CheckInAt).Hours()
			rec.HoursSpent = &hours
		}
		records = append(records, rec)
	}
	return records, total, rows.Err()
}

// RecordCheckIn records a staff check-in.
func (r *Repository) RecordCheckIn(ctx context.Context, orgID, staffMemberID, method string) (*AttendanceRecord, error) {
	var userID, staffName string
	err := r.db.QueryRow(ctx,
		`SELECT om.user_id, COALESCE(u.name, '')
		 FROM organization_members om
		 JOIN users u ON u.id = om.user_id
		 WHERE om.id = $1 AND om.organization_id = $2
		   AND om.role IN ('staff', 'admin') AND om.status = 'active'`,
		staffMemberID, orgID,
	).Scan(&userID, &staffName)
	if err != nil {
		return nil, fmt.Errorf("staff member not found: %w", err)
	}

	var rec AttendanceRecord
	err = r.db.QueryRow(ctx,
		`INSERT INTO attendance (user_id, organization_id, check_in_method, check_in_at)
		 VALUES ($1, $2, $3, NOW())
		 RETURNING id, check_in_at, check_in_method`,
		userID, orgID, method,
	).Scan(&rec.ID, &rec.CheckInAt, &rec.Method)
	if err != nil {
		return nil, fmt.Errorf("recording staff check-in: %w", err)
	}
	rec.StaffID = staffMemberID
	rec.StaffName = staffName
	return &rec, nil
}

// RecordCheckOut records a staff check-out on the most recent open attendance record.
func (r *Repository) RecordCheckOut(ctx context.Context, orgID, staffMemberID string) (*AttendanceRecord, error) {
	var userID, staffName string
	err := r.db.QueryRow(ctx,
		`SELECT om.user_id, COALESCE(u.name, '')
		 FROM organization_members om
		 JOIN users u ON u.id = om.user_id
		 WHERE om.id = $1 AND om.organization_id = $2
		   AND om.role IN ('staff', 'admin') AND om.status = 'active'`,
		staffMemberID, orgID,
	).Scan(&userID, &staffName)
	if err != nil {
		return nil, fmt.Errorf("staff member not found: %w", err)
	}

	var rec AttendanceRecord
	err = r.db.QueryRow(ctx,
		`UPDATE attendance SET check_out_at = NOW()
		 WHERE id = (
			SELECT id FROM attendance
			WHERE user_id = $1 AND organization_id = $2 AND check_out_at IS NULL
			ORDER BY check_in_at DESC LIMIT 1
		 )
		 RETURNING id, check_in_at, check_out_at, check_in_method`,
		userID, orgID,
	).Scan(&rec.ID, &rec.CheckInAt, &rec.CheckOutAt, &rec.Method)
	if err != nil {
		return nil, fmt.Errorf("no open check-in found for staff member: %w", err)
	}
	rec.StaffID = staffMemberID
	rec.StaffName = staffName
	if rec.CheckOutAt != nil {
		hours := rec.CheckOutAt.Sub(rec.CheckInAt).Hours()
		rec.HoursSpent = &hours
	}
	return &rec, nil
}
