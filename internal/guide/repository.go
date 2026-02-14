package guide

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// TrainingGuide represents a training guide with bilingual content.
type TrainingGuide struct {
	ID            string    `json:"id"`
	Title         string    `json:"title"`
	TitleNe       string    `json:"title_ne,omitempty"`
	Content       string    `json:"content"`
	ContentNe     string    `json:"content_ne,omitempty"`
	Category      string    `json:"category,omitempty"`
	CoverImageURL string    `json:"cover_image_url,omitempty"`
	AuthorID      string    `json:"author_id,omitempty"`
	IsPublished   bool      `json:"is_published"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// Repository handles training guide data persistence.
type Repository struct {
	db *pgxpool.Pool
}

// NewRepository creates a new guide repository.
func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

// Create inserts a new training guide.
func (r *Repository) Create(ctx context.Context, title, titleNe, content, contentNe, category, coverImageURL, authorID string) (*TrainingGuide, error) {
	var g TrainingGuide
	err := r.db.QueryRow(ctx,
		`INSERT INTO training_guides (title, title_ne, content, content_ne, category, cover_image_url, author_id)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)
		 RETURNING id, title, COALESCE(title_ne, ''), content, COALESCE(content_ne, ''),
		           COALESCE(category, ''), COALESCE(cover_image_url, ''), COALESCE(author_id::text, ''),
		           is_published, created_at, updated_at`,
		title, nilIfEmpty(titleNe), content, nilIfEmpty(contentNe),
		nilIfEmpty(category), nilIfEmpty(coverImageURL), nilIfEmpty(authorID),
	).Scan(&g.ID, &g.Title, &g.TitleNe, &g.Content, &g.ContentNe,
		&g.Category, &g.CoverImageURL, &g.AuthorID,
		&g.IsPublished, &g.CreatedAt, &g.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("creating training guide: %w", err)
	}
	return &g, nil
}

// GetByID retrieves a training guide by ID.
func (r *Repository) GetByID(ctx context.Context, guideID string) (*TrainingGuide, error) {
	var g TrainingGuide
	err := r.db.QueryRow(ctx,
		`SELECT id, title, COALESCE(title_ne, ''), content, COALESCE(content_ne, ''),
		        COALESCE(category, ''), COALESCE(cover_image_url, ''), COALESCE(author_id::text, ''),
		        is_published, created_at, updated_at
		 FROM training_guides WHERE id = $1`,
		guideID,
	).Scan(&g.ID, &g.Title, &g.TitleNe, &g.Content, &g.ContentNe,
		&g.Category, &g.CoverImageURL, &g.AuthorID,
		&g.IsPublished, &g.CreatedAt, &g.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("getting training guide: %w", err)
	}
	return &g, nil
}

// GetPublishedByID retrieves a published training guide by ID.
func (r *Repository) GetPublishedByID(ctx context.Context, guideID string) (*TrainingGuide, error) {
	var g TrainingGuide
	err := r.db.QueryRow(ctx,
		`SELECT id, title, COALESCE(title_ne, ''), content, COALESCE(content_ne, ''),
		        COALESCE(category, ''), COALESCE(cover_image_url, ''), COALESCE(author_id::text, ''),
		        is_published, created_at, updated_at
		 FROM training_guides WHERE id = $1 AND is_published = true`,
		guideID,
	).Scan(&g.ID, &g.Title, &g.TitleNe, &g.Content, &g.ContentNe,
		&g.Category, &g.CoverImageURL, &g.AuthorID,
		&g.IsPublished, &g.CreatedAt, &g.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("getting published training guide: %w", err)
	}
	return &g, nil
}

// ListPublished retrieves published training guides with optional category filter, search, and pagination.
func (r *Repository) ListPublished(ctx context.Context, category, search string, limit, offset int) ([]TrainingGuide, int, error) {
	conditions := []string{"is_published = true"}
	args := []interface{}{}
	argIdx := 1

	if category != "" {
		conditions = append(conditions, fmt.Sprintf("category = $%d", argIdx))
		args = append(args, category)
		argIdx++
	}

	if search != "" {
		conditions = append(conditions, fmt.Sprintf("(title ILIKE $%d OR content ILIKE $%d OR COALESCE(title_ne, '') ILIKE $%d)", argIdx, argIdx, argIdx))
		args = append(args, "%"+search+"%")
		argIdx++
	}

	where := "WHERE " + conditions[0]
	for i := 1; i < len(conditions); i++ {
		where += " AND " + conditions[i]
	}

	// Count total
	var total int
	countQuery := "SELECT COUNT(*) FROM training_guides " + where
	if err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("counting training guides: %w", err)
	}

	// Fetch page
	query := fmt.Sprintf(
		`SELECT id, title, COALESCE(title_ne, ''), content, COALESCE(content_ne, ''),
		        COALESCE(category, ''), COALESCE(cover_image_url, ''), COALESCE(author_id::text, ''),
		        is_published, created_at, updated_at
		 FROM training_guides %s
		 ORDER BY created_at DESC LIMIT $%d OFFSET $%d`,
		where, argIdx, argIdx+1)
	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("listing published training guides: %w", err)
	}
	defer rows.Close()

	var guides []TrainingGuide
	for rows.Next() {
		var g TrainingGuide
		if err := rows.Scan(&g.ID, &g.Title, &g.TitleNe, &g.Content, &g.ContentNe,
			&g.Category, &g.CoverImageURL, &g.AuthorID,
			&g.IsPublished, &g.CreatedAt, &g.UpdatedAt); err != nil {
			return nil, 0, fmt.Errorf("scanning training guide: %w", err)
		}
		guides = append(guides, g)
	}
	return guides, total, rows.Err()
}

// ListAll retrieves all training guides (for admin) with optional category filter, search, and pagination.
func (r *Repository) ListAll(ctx context.Context, category, search string, limit, offset int) ([]TrainingGuide, int, error) {
	conditions := []string{}
	args := []interface{}{}
	argIdx := 1

	if category != "" {
		conditions = append(conditions, fmt.Sprintf("category = $%d", argIdx))
		args = append(args, category)
		argIdx++
	}

	if search != "" {
		conditions = append(conditions, fmt.Sprintf("(title ILIKE $%d OR content ILIKE $%d OR COALESCE(title_ne, '') ILIKE $%d)", argIdx, argIdx, argIdx))
		args = append(args, "%"+search+"%")
		argIdx++
	}

	where := ""
	if len(conditions) > 0 {
		where = "WHERE " + conditions[0]
		for i := 1; i < len(conditions); i++ {
			where += " AND " + conditions[i]
		}
	}

	// Count total
	var total int
	countQuery := "SELECT COUNT(*) FROM training_guides " + where
	if err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("counting training guides: %w", err)
	}

	// Fetch page
	query := fmt.Sprintf(
		`SELECT id, title, COALESCE(title_ne, ''), content, COALESCE(content_ne, ''),
		        COALESCE(category, ''), COALESCE(cover_image_url, ''), COALESCE(author_id::text, ''),
		        is_published, created_at, updated_at
		 FROM training_guides %s
		 ORDER BY created_at DESC LIMIT $%d OFFSET $%d`,
		where, argIdx, argIdx+1)
	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("listing training guides: %w", err)
	}
	defer rows.Close()

	var guides []TrainingGuide
	for rows.Next() {
		var g TrainingGuide
		if err := rows.Scan(&g.ID, &g.Title, &g.TitleNe, &g.Content, &g.ContentNe,
			&g.Category, &g.CoverImageURL, &g.AuthorID,
			&g.IsPublished, &g.CreatedAt, &g.UpdatedAt); err != nil {
			return nil, 0, fmt.Errorf("scanning training guide: %w", err)
		}
		guides = append(guides, g)
	}
	return guides, total, rows.Err()
}

// Update modifies an existing training guide.
func (r *Repository) Update(ctx context.Context, guideID, title, titleNe, content, contentNe, category, coverImageURL string) (*TrainingGuide, error) {
	var g TrainingGuide
	err := r.db.QueryRow(ctx,
		`UPDATE training_guides
		 SET title = $2, title_ne = $3, content = $4, content_ne = $5, category = $6, cover_image_url = $7
		 WHERE id = $1
		 RETURNING id, title, COALESCE(title_ne, ''), content, COALESCE(content_ne, ''),
		           COALESCE(category, ''), COALESCE(cover_image_url, ''), COALESCE(author_id::text, ''),
		           is_published, created_at, updated_at`,
		guideID, title, nilIfEmpty(titleNe), content, nilIfEmpty(contentNe),
		nilIfEmpty(category), nilIfEmpty(coverImageURL),
	).Scan(&g.ID, &g.Title, &g.TitleNe, &g.Content, &g.ContentNe,
		&g.Category, &g.CoverImageURL, &g.AuthorID,
		&g.IsPublished, &g.CreatedAt, &g.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("updating training guide: %w", err)
	}
	return &g, nil
}

// SetPublished updates the published status of a training guide.
func (r *Repository) SetPublished(ctx context.Context, guideID string, isPublished bool) (*TrainingGuide, error) {
	var g TrainingGuide
	err := r.db.QueryRow(ctx,
		`UPDATE training_guides
		 SET is_published = $2
		 WHERE id = $1
		 RETURNING id, title, COALESCE(title_ne, ''), content, COALESCE(content_ne, ''),
		           COALESCE(category, ''), COALESCE(cover_image_url, ''), COALESCE(author_id::text, ''),
		           is_published, created_at, updated_at`,
		guideID, isPublished,
	).Scan(&g.ID, &g.Title, &g.TitleNe, &g.Content, &g.ContentNe,
		&g.Category, &g.CoverImageURL, &g.AuthorID,
		&g.IsPublished, &g.CreatedAt, &g.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("updating training guide publish status: %w", err)
	}
	return &g, nil
}

// Delete removes a training guide.
func (r *Repository) Delete(ctx context.Context, guideID string) error {
	tag, err := r.db.Exec(ctx,
		`DELETE FROM training_guides WHERE id = $1`,
		guideID,
	)
	if err != nil {
		return fmt.Errorf("deleting training guide: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("training guide not found")
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
