package feedback

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Feedback represents member feedback.
type Feedback struct {
	ID         string    `json:"id"`
	OrgID      string    `json:"organization_id"`
	UserID     string    `json:"user_id"`
	MemberName string    `json:"member_name,omitempty"`
	AvatarURL  string    `json:"avatar_url,omitempty"`
	Message    string    `json:"message"`
	CreatedAt  time.Time `json:"created_at"`
}

// Repository handles feedback data persistence.
type Repository struct {
	db *pgxpool.Pool
}

// NewRepository creates a new feedback repository.
func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

// Create inserts a new feedback entry.
func (r *Repository) Create(ctx context.Context, orgID, userID, message string) (*Feedback, error) {
	var f Feedback
	err := r.db.QueryRow(ctx,
		`INSERT INTO feedbacks (organization_id, user_id, message)
		 VALUES ($1, $2, $3)
		 RETURNING id, organization_id, user_id, message, created_at`,
		orgID, userID, message,
	).Scan(&f.ID, &f.OrgID, &f.UserID, &f.Message, &f.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("creating feedback: %w", err)
	}
	return &f, nil
}

// List retrieves feedbacks for an organization with member info (name, avatar).
func (r *Repository) List(ctx context.Context, orgID string, limit, offset int) ([]Feedback, error) {
	rows, err := r.db.Query(ctx,
		`SELECT f.id, f.organization_id, f.user_id,
		        COALESCE(u.name, ''), COALESCE(u.avatar_url, ''),
		        f.message, f.created_at
		 FROM feedbacks f
		 JOIN users u ON u.id = f.user_id
		 WHERE f.organization_id = $1
		 ORDER BY f.created_at DESC LIMIT $2 OFFSET $3`,
		orgID, limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("listing feedbacks: %w", err)
	}
	defer rows.Close()

	var feedbacks []Feedback
	for rows.Next() {
		var f Feedback
		if err := rows.Scan(&f.ID, &f.OrgID, &f.UserID, &f.MemberName, &f.AvatarURL,
			&f.Message, &f.CreatedAt); err != nil {
			return nil, fmt.Errorf("scanning feedback: %w", err)
		}
		feedbacks = append(feedbacks, f)
	}
	return feedbacks, rows.Err()
}
