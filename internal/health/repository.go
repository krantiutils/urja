package health

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrMetricNotFound = errors.New("health metric not found")

// HealthMetric represents a row in the health_metrics table.
type HealthMetric struct {
	ID         string          `json:"id"`
	MemberID   string          `json:"member_id"`
	MetricType string          `json:"metric_type"`
	Value      json.RawMessage `json:"value"`
	RecordedAt time.Time       `json:"recorded_at"`
}

// ProgressPhoto represents a row in the progress_photos table.
type ProgressPhoto struct {
	ID       string    `json:"id"`
	MemberID string    `json:"member_id"`
	PhotoURL string    `json:"photo_url"`
	Category string    `json:"category,omitempty"`
	Notes    string    `json:"notes,omitempty"`
	TakenAt  time.Time `json:"taken_at"`
}

// Repository handles health data persistence.
type Repository struct {
	db *pgxpool.Pool
}

// NewRepository creates a new health repository.
func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

// InsertMetric inserts a health metric record.
func (r *Repository) InsertMetric(ctx context.Context, memberID, metricType string, value json.RawMessage) (*HealthMetric, error) {
	var m HealthMetric
	err := r.db.QueryRow(ctx,
		`INSERT INTO health_metrics (member_id, metric_type, value)
		 VALUES ($1, $2, $3)
		 RETURNING id, member_id, metric_type, value, recorded_at`,
		memberID, metricType, value,
	).Scan(&m.ID, &m.MemberID, &m.MetricType, &m.Value, &m.RecordedAt)
	if err != nil {
		return nil, fmt.Errorf("inserting health metric: %w", err)
	}
	return &m, nil
}

// ListMetrics retrieves health metrics for a member with optional type and date filtering.
func (r *Repository) ListMetrics(ctx context.Context, memberID string, metricType string, from, to *time.Time, limit int) ([]HealthMetric, error) {
	query := `SELECT id, member_id, metric_type, value, recorded_at
		 FROM health_metrics
		 WHERE member_id = $1`
	args := []interface{}{memberID}
	argIdx := 2

	if metricType != "" {
		query += fmt.Sprintf(" AND metric_type = $%d", argIdx)
		args = append(args, metricType)
		argIdx++
	}
	if from != nil {
		query += fmt.Sprintf(" AND recorded_at >= $%d", argIdx)
		args = append(args, *from)
		argIdx++
	}
	if to != nil {
		query += fmt.Sprintf(" AND recorded_at <= $%d", argIdx)
		args = append(args, *to)
		argIdx++
	}

	query += " ORDER BY recorded_at DESC"

	if limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argIdx)
		args = append(args, limit)
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("querying health metrics: %w", err)
	}
	defer rows.Close()

	var metrics []HealthMetric
	for rows.Next() {
		var m HealthMetric
		if err := rows.Scan(&m.ID, &m.MemberID, &m.MetricType, &m.Value, &m.RecordedAt); err != nil {
			return nil, fmt.Errorf("scanning health metric: %w", err)
		}
		metrics = append(metrics, m)
	}
	return metrics, rows.Err()
}

// InsertPhoto inserts a progress photo record.
func (r *Repository) InsertPhoto(ctx context.Context, memberID, photoURL, category, notes string) (*ProgressPhoto, error) {
	var p ProgressPhoto
	err := r.db.QueryRow(ctx,
		`INSERT INTO progress_photos (member_id, photo_url, category, notes)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, member_id, photo_url, COALESCE(category, ''), COALESCE(notes, ''), taken_at`,
		memberID, photoURL, nilIfEmpty(category), nilIfEmpty(notes),
	).Scan(&p.ID, &p.MemberID, &p.PhotoURL, &p.Category, &p.Notes, &p.TakenAt)
	if err != nil {
		return nil, fmt.Errorf("inserting progress photo: %w", err)
	}
	return &p, nil
}

// ListPhotos retrieves progress photos for a member with optional category filter.
func (r *Repository) ListPhotos(ctx context.Context, memberID, category string, limit, offset int) ([]ProgressPhoto, int, error) {
	countQuery := `SELECT COUNT(*) FROM progress_photos WHERE member_id = $1`
	dataQuery := `SELECT id, member_id, photo_url, COALESCE(category, ''), COALESCE(notes, ''), taken_at
		 FROM progress_photos WHERE member_id = $1`
	args := []interface{}{memberID}
	argIdx := 2

	if category != "" {
		filter := fmt.Sprintf(" AND category = $%d", argIdx)
		countQuery += filter
		dataQuery += filter
		args = append(args, category)
		argIdx++
	}

	var total int
	if err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("counting progress photos: %w", err)
	}

	dataQuery += " ORDER BY taken_at DESC"
	dataQuery += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	dataArgs := append(args, limit, offset)

	rows, err := r.db.Query(ctx, dataQuery, dataArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("querying progress photos: %w", err)
	}
	defer rows.Close()

	var photos []ProgressPhoto
	for rows.Next() {
		var p ProgressPhoto
		if err := rows.Scan(&p.ID, &p.MemberID, &p.PhotoURL, &p.Category, &p.Notes, &p.TakenAt); err != nil {
			return nil, 0, fmt.Errorf("scanning progress photo: %w", err)
		}
		photos = append(photos, p)
	}
	return photos, total, rows.Err()
}

// GetMetricByID retrieves a single health metric by ID, scoped to the member.
func (r *Repository) GetMetricByID(ctx context.Context, memberID, metricID string) (*HealthMetric, error) {
	var m HealthMetric
	err := r.db.QueryRow(ctx,
		`SELECT id, member_id, metric_type, value, recorded_at
		 FROM health_metrics
		 WHERE id = $1 AND member_id = $2`,
		metricID, memberID,
	).Scan(&m.ID, &m.MemberID, &m.MetricType, &m.Value, &m.RecordedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrMetricNotFound
		}
		return nil, fmt.Errorf("querying health metric: %w", err)
	}
	return &m, nil
}

func nilIfEmpty(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
