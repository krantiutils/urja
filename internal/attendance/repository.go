package attendance

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Record represents an attendance check-in.
type Record struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	OrgID     string    `json:"org_id"`
	CheckInAt time.Time `json:"check_in_at"`
	Method    string    `json:"method"`
}

// Streak represents a member's attendance streak at an organization.
type Streak struct {
	ID            string    `json:"id"`
	MemberID      string    `json:"member_id"`
	OrgID         string    `json:"org_id"`
	CurrentStreak int       `json:"current_streak"`
	LongestStreak int       `json:"longest_streak"`
	LastCheckIn   *string   `json:"last_check_in"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// Repository handles attendance data persistence.
type Repository struct {
	db *pgxpool.Pool
}

// NewRepository creates a new attendance repository.
func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

// CheckIn records a new attendance entry.
func (r *Repository) CheckIn(ctx context.Context, userID, orgID, method string) (*Record, error) {
	var rec Record
	err := r.db.QueryRow(ctx,
		`INSERT INTO attendance (user_id, organization_id, check_in_method, check_in_at)
		 VALUES ($1, $2, $3, NOW())
		 RETURNING id, user_id, organization_id, check_in_at, check_in_method`,
		userID, orgID, method,
	).Scan(&rec.ID, &rec.UserID, &rec.OrgID, &rec.CheckInAt, &rec.Method)
	if err != nil {
		return nil, fmt.Errorf("recording check-in: %w", err)
	}
	return &rec, nil
}

// ListByUser retrieves attendance records for a user.
func (r *Repository) ListByUser(ctx context.Context, userID string, limit, offset int) ([]Record, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, user_id, organization_id, check_in_at, check_in_method
		 FROM attendance WHERE user_id = $1
		 ORDER BY check_in_at DESC LIMIT $2 OFFSET $3`,
		userID, limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("querying attendance: %w", err)
	}
	defer rows.Close()

	var records []Record
	for rows.Next() {
		var rec Record
		if err := rows.Scan(&rec.ID, &rec.UserID, &rec.OrgID, &rec.CheckInAt, &rec.Method); err != nil {
			return nil, fmt.Errorf("scanning attendance: %w", err)
		}
		records = append(records, rec)
	}
	return records, rows.Err()
}

// ListByOrg retrieves attendance records for an organization.
func (r *Repository) ListByOrg(ctx context.Context, orgID string, limit, offset int) ([]Record, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, user_id, organization_id, check_in_at, check_in_method
		 FROM attendance WHERE organization_id = $1
		 ORDER BY check_in_at DESC LIMIT $2 OFFSET $3`,
		orgID, limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("querying org attendance: %w", err)
	}
	defer rows.Close()

	var records []Record
	for rows.Next() {
		var rec Record
		if err := rows.Scan(&rec.ID, &rec.UserID, &rec.OrgID, &rec.CheckInAt, &rec.Method); err != nil {
			return nil, fmt.Errorf("scanning attendance: %w", err)
		}
		records = append(records, rec)
	}
	return records, rows.Err()
}

// UpsertStreak updates the streak for a member at an org after a check-in.
// dateStr must be "YYYY-MM-DD" in the gym's local timezone.
// Handles: same-day idempotency, consecutive-day increment, gap reset.
func (r *Repository) UpsertStreak(ctx context.Context, userID, orgID, dateStr string) (*Streak, error) {
	var s Streak
	var lastCheckIn *string

	err := r.db.QueryRow(ctx,
		`INSERT INTO streaks (member_id, organization_id, current_streak, longest_streak, last_check_in)
		 VALUES ($1, $2, 1, 1, $3::date)
		 ON CONFLICT (member_id, organization_id)
		 DO UPDATE SET
		     current_streak = CASE
		         WHEN streaks.last_check_in = $3::date THEN streaks.current_streak
		         WHEN streaks.last_check_in = ($3::date - 1) THEN streaks.current_streak + 1
		         ELSE 1
		     END,
		     longest_streak = GREATEST(streaks.longest_streak,
		         CASE
		             WHEN streaks.last_check_in = $3::date THEN streaks.current_streak
		             WHEN streaks.last_check_in = ($3::date - 1) THEN streaks.current_streak + 1
		             ELSE 1
		         END),
		     last_check_in = $3::date
		 RETURNING id, member_id, organization_id, current_streak, longest_streak,
		           last_check_in::text, updated_at`,
		userID, orgID, dateStr,
	).Scan(&s.ID, &s.MemberID, &s.OrgID, &s.CurrentStreak, &s.LongestStreak, &lastCheckIn, &s.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("upserting streak: %w", err)
	}
	s.LastCheckIn = lastCheckIn
	return &s, nil
}

// GetStreaksByUser retrieves all streaks for a user across organizations.
func (r *Repository) GetStreaksByUser(ctx context.Context, userID string) ([]Streak, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, member_id, organization_id, current_streak, longest_streak,
		        last_check_in::text, updated_at
		 FROM streaks WHERE member_id = $1
		 ORDER BY updated_at DESC`,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("querying streaks: %w", err)
	}
	defer rows.Close()

	var streaks []Streak
	for rows.Next() {
		var s Streak
		var lastCheckIn *string
		if err := rows.Scan(&s.ID, &s.MemberID, &s.OrgID, &s.CurrentStreak, &s.LongestStreak, &lastCheckIn, &s.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scanning streak: %w", err)
		}
		s.LastCheckIn = lastCheckIn
		streaks = append(streaks, s)
	}
	return streaks, rows.Err()
}

// LookupNFCCard finds the user_id for an active NFC card at an organization.
func (r *Repository) LookupNFCCard(ctx context.Context, orgID, cardHex string) (string, error) {
	var userID string
	err := r.db.QueryRow(ctx,
		`SELECT user_id FROM nfc_cards
		 WHERE organization_id = $1 AND card_hex = $2 AND is_active = true`,
		orgID, cardHex,
	).Scan(&userID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", fmt.Errorf("NFC card not found or inactive")
		}
		return "", fmt.Errorf("looking up NFC card: %w", err)
	}
	return userID, nil
}

// IsOrgMember checks if a user is an active member of an organization.
func (r *Repository) IsOrgMember(ctx context.Context, userID, orgID string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx,
		`SELECT EXISTS(
			SELECT 1 FROM organization_members
			WHERE user_id = $1 AND organization_id = $2 AND status = 'active'
		)`,
		userID, orgID,
	).Scan(&exists)
	return exists, err
}

// GetMonthlyCalendar returns the distinct days in a month on which a user checked in.
// month must be in "YYYY-MM" format.
func (r *Repository) GetMonthlyCalendar(ctx context.Context, userID, month string) ([]int, error) {
	rows, err := r.db.Query(ctx,
		`SELECT DISTINCT EXTRACT(DAY FROM check_in_at AT TIME ZONE 'Asia/Kathmandu')::int AS day
		 FROM attendance
		 WHERE user_id = $1
		   AND check_in_at >= ($2 || '-01')::date
		   AND check_in_at < (($2 || '-01')::date + INTERVAL '1 month')
		 ORDER BY day`,
		userID, month,
	)
	if err != nil {
		return nil, fmt.Errorf("querying monthly calendar: %w", err)
	}
	defer rows.Close()

	var days []int
	for rows.Next() {
		var day int
		if err := rows.Scan(&day); err != nil {
			return nil, fmt.Errorf("scanning calendar day: %w", err)
		}
		days = append(days, day)
	}
	return days, rows.Err()
}

// GetOrgRole returns the user's role in an organization.
func (r *Repository) GetOrgRole(ctx context.Context, userID, orgID string) (string, error) {
	var role string
	err := r.db.QueryRow(ctx,
		`SELECT role FROM organization_members
		 WHERE user_id = $1 AND organization_id = $2 AND status = 'active'`,
		userID, orgID,
	).Scan(&role)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", fmt.Errorf("user is not a member of this organization")
		}
		return "", fmt.Errorf("getting org role: %w", err)
	}
	return role, nil
}
