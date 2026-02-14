package notice

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Notice represents a gym announcement.
type Notice struct {
	ID        string    `json:"id"`
	OrgID     string    `json:"organization_id"`
	Title     string    `json:"title"`
	TitleNe   string    `json:"title_ne,omitempty"`
	Content   string    `json:"content"`
	ContentNe string    `json:"content_ne,omitempty"`
	CreatedBy string    `json:"created_by"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Repository handles notice data persistence.
type Repository struct {
	db *pgxpool.Pool
}

// NewRepository creates a new notice repository.
func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

// Create inserts a new notice.
func (r *Repository) Create(ctx context.Context, orgID, title, titleNe, content, contentNe, createdBy string) (*Notice, error) {
	var n Notice
	err := r.db.QueryRow(ctx,
		`INSERT INTO notices (organization_id, title, title_ne, content, content_ne, created_by)
		 VALUES ($1, $2, $3, $4, $5, $6)
		 RETURNING id, organization_id, title, COALESCE(title_ne, ''), content, COALESCE(content_ne, ''),
		           created_by, is_active, created_at, updated_at`,
		orgID, title, nilIfEmpty(titleNe), content, nilIfEmpty(contentNe), createdBy,
	).Scan(&n.ID, &n.OrgID, &n.Title, &n.TitleNe, &n.Content, &n.ContentNe,
		&n.CreatedBy, &n.IsActive, &n.CreatedAt, &n.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("creating notice: %w", err)
	}
	return &n, nil
}

// GetByID retrieves a single notice by ID within an organization.
func (r *Repository) GetByID(ctx context.Context, orgID, noticeID string) (*Notice, error) {
	var n Notice
	err := r.db.QueryRow(ctx,
		`SELECT id, organization_id, title, COALESCE(title_ne, ''), content, COALESCE(content_ne, ''),
		        created_by, is_active, created_at, updated_at
		 FROM notices WHERE id = $1 AND organization_id = $2`,
		noticeID, orgID,
	).Scan(&n.ID, &n.OrgID, &n.Title, &n.TitleNe, &n.Content, &n.ContentNe,
		&n.CreatedBy, &n.IsActive, &n.CreatedAt, &n.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("getting notice: %w", err)
	}
	return &n, nil
}

// List retrieves notices for an organization with optional search.
func (r *Repository) List(ctx context.Context, orgID, search string, limit, offset int) ([]Notice, error) {
	var query string
	var args []interface{}

	if search != "" {
		query = `SELECT id, organization_id, title, COALESCE(title_ne, ''), content, COALESCE(content_ne, ''),
		                created_by, is_active, created_at, updated_at
		         FROM notices WHERE organization_id = $1
		           AND (title ILIKE $2 OR content ILIKE $2)
		         ORDER BY created_at DESC LIMIT $3 OFFSET $4`
		args = []interface{}{orgID, "%" + search + "%", limit, offset}
	} else {
		query = `SELECT id, organization_id, title, COALESCE(title_ne, ''), content, COALESCE(content_ne, ''),
		                created_by, is_active, created_at, updated_at
		         FROM notices WHERE organization_id = $1
		         ORDER BY created_at DESC LIMIT $2 OFFSET $3`
		args = []interface{}{orgID, limit, offset}
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("listing notices: %w", err)
	}
	defer rows.Close()

	var notices []Notice
	for rows.Next() {
		var n Notice
		if err := rows.Scan(&n.ID, &n.OrgID, &n.Title, &n.TitleNe, &n.Content, &n.ContentNe,
			&n.CreatedBy, &n.IsActive, &n.CreatedAt, &n.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scanning notice: %w", err)
		}
		notices = append(notices, n)
	}
	return notices, rows.Err()
}

// Update modifies an existing notice.
func (r *Repository) Update(ctx context.Context, orgID, noticeID, title, titleNe, content, contentNe string, isActive bool) (*Notice, error) {
	var n Notice
	err := r.db.QueryRow(ctx,
		`UPDATE notices
		 SET title = $3, title_ne = $4, content = $5, content_ne = $6, is_active = $7
		 WHERE id = $1 AND organization_id = $2
		 RETURNING id, organization_id, title, COALESCE(title_ne, ''), content, COALESCE(content_ne, ''),
		           created_by, is_active, created_at, updated_at`,
		noticeID, orgID, title, nilIfEmpty(titleNe), content, nilIfEmpty(contentNe), isActive,
	).Scan(&n.ID, &n.OrgID, &n.Title, &n.TitleNe, &n.Content, &n.ContentNe,
		&n.CreatedBy, &n.IsActive, &n.CreatedAt, &n.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("updating notice: %w", err)
	}
	return &n, nil
}

// Delete removes a notice.
func (r *Repository) Delete(ctx context.Context, orgID, noticeID string) error {
	tag, err := r.db.Exec(ctx,
		`DELETE FROM notices WHERE id = $1 AND organization_id = $2`,
		noticeID, orgID,
	)
	if err != nil {
		return fmt.Errorf("deleting notice: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("notice not found")
	}
	return nil
}

// nilIfEmpty returns nil for empty strings so the DB stores NULL.
func nilIfEmpty(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}
