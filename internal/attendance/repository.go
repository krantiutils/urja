package attendance

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Record represents an attendance check-in.
type Record struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	OrgID     string    `json:"org_id"`
	CheckInAt time.Time `json:"check_in_at"`
	Method    string    `json:"method"` // qr, nfc, manual
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
		`INSERT INTO attendance (user_id, organization_id, method, check_in_at)
		 VALUES ($1, $2, $3, NOW())
		 RETURNING id, user_id, organization_id, check_in_at, method`,
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
		`SELECT id, user_id, organization_id, check_in_at, method
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
		`SELECT id, user_id, organization_id, check_in_at, method
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
