package member

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Member represents a user in the system.
type Member struct {
	ID        string    `json:"id"`
	Phone     string    `json:"phone"`
	Name      string    `json:"name,omitempty"`
	Email     string    `json:"email,omitempty"`
	AvatarURL string    `json:"avatar_url,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Repository handles member data persistence.
type Repository struct {
	db *pgxpool.Pool
}

// NewRepository creates a new member repository.
func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

// GetByID retrieves a member by their ID.
// SECURITY: Never returns sensitive fields (password_hash, fcm_token, etc.)
func (r *Repository) GetByID(ctx context.Context, id string) (*Member, error) {
	var m Member
	err := r.db.QueryRow(ctx,
		`SELECT id, phone, COALESCE(name, ''), COALESCE(email, ''),
		        COALESCE(avatar_url, ''), created_at, updated_at
		 FROM users WHERE id = $1`,
		id,
	).Scan(&m.ID, &m.Phone, &m.Name, &m.Email, &m.AvatarURL, &m.CreatedAt, &m.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("member not found: %w", err)
	}
	return &m, nil
}

// Update updates a member's profile fields.
func (r *Repository) Update(ctx context.Context, id string, name, email, avatarURL *string) error {
	_, err := r.db.Exec(ctx,
		`UPDATE users SET
			name = COALESCE($2, name),
			email = COALESCE($3, email),
			avatar_url = COALESCE($4, avatar_url),
			updated_at = NOW()
		 WHERE id = $1`,
		id, name, email, avatarURL,
	)
	return err
}

// ListByOrg retrieves members of an organization with pagination.
func (r *Repository) ListByOrg(ctx context.Context, orgID string, limit, offset int) ([]Member, int, error) {
	var total int
	err := r.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM organization_members WHERE organization_id = $1 AND status = 'active'`,
		orgID,
	).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("counting members: %w", err)
	}

	rows, err := r.db.Query(ctx,
		`SELECT u.id, u.phone, COALESCE(u.name, ''), COALESCE(u.email, ''),
		        COALESCE(u.avatar_url, ''), u.created_at, u.updated_at
		 FROM users u
		 JOIN organization_members om ON u.id = om.user_id
		 WHERE om.organization_id = $1 AND om.status = 'active'
		 ORDER BY u.created_at DESC
		 LIMIT $2 OFFSET $3`,
		orgID, limit, offset,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("querying members: %w", err)
	}
	defer rows.Close()

	var members []Member
	for rows.Next() {
		var m Member
		if err := rows.Scan(&m.ID, &m.Phone, &m.Name, &m.Email, &m.AvatarURL, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, 0, fmt.Errorf("scanning member: %w", err)
		}
		members = append(members, m)
	}

	return members, total, rows.Err()
}
