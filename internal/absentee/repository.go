package absentee

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Absentee represents a member who hasn't checked in recently.
type Absentee struct {
	UserID     string    `json:"user_id"`
	Name       string    `json:"name"`
	AvatarURL  string    `json:"avatar_url,omitempty"`
	Phone      string    `json:"phone"`
	AbsentDays int       `json:"absent_days"`
	JoinedDate time.Time `json:"joined_date"`
}

// Repository handles absentee data queries.
type Repository struct {
	db *pgxpool.Pool
}

// NewRepository creates a new absentee repository.
func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

// List retrieves members who haven't checked in for at least minDays days.
func (r *Repository) List(ctx context.Context, orgID string, minDays, limit, offset int) ([]Absentee, error) {
	rows, err := r.db.Query(ctx,
		`SELECT u.id, COALESCE(u.name, ''), COALESCE(u.avatar_url, ''), u.phone,
		        EXTRACT(DAY FROM NOW() - COALESCE(
		            (SELECT MAX(a.check_in_at) FROM attendance a
		             WHERE a.user_id = u.id AND a.organization_id = $1),
		            om.created_at
		        ))::INTEGER AS absent_days,
		        om.created_at AS joined_date
		 FROM users u
		 JOIN organization_members om ON om.user_id = u.id
		 WHERE om.organization_id = $1 AND om.status = 'active'
		   AND EXTRACT(DAY FROM NOW() - COALESCE(
		       (SELECT MAX(a.check_in_at) FROM attendance a
		        WHERE a.user_id = u.id AND a.organization_id = $1),
		       om.created_at
		   )) >= $2
		 ORDER BY absent_days DESC
		 LIMIT $3 OFFSET $4`,
		orgID, minDays, limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("listing absentees: %w", err)
	}
	defer rows.Close()

	var absentees []Absentee
	for rows.Next() {
		var a Absentee
		if err := rows.Scan(&a.UserID, &a.Name, &a.AvatarURL, &a.Phone,
			&a.AbsentDays, &a.JoinedDate); err != nil {
			return nil, fmt.Errorf("scanning absentee: %w", err)
		}
		absentees = append(absentees, a)
	}
	return absentees, rows.Err()
}

// GetMemberPhone retrieves a member's phone number within the org.
func (r *Repository) GetMemberPhone(ctx context.Context, orgID, memberID string) (string, error) {
	var phone string
	err := r.db.QueryRow(ctx,
		`SELECT u.phone FROM users u
		 JOIN organization_members om ON om.user_id = u.id
		 WHERE u.id = $1 AND om.organization_id = $2 AND om.status = 'active'`,
		memberID, orgID,
	).Scan(&phone)
	if err != nil {
		return "", fmt.Errorf("member not found or not active: %w", err)
	}
	return phone, nil
}
