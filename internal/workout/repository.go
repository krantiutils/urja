package workout

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// WorkoutTemplate represents a workout template.
type WorkoutTemplate struct {
	ID              string          `json:"id"`
	OrgID           *string         `json:"organization_id"`
	Name            string          `json:"name"`
	NameNe          string          `json:"name_ne,omitempty"`
	Description     string          `json:"description,omitempty"`
	DescriptionNe   string          `json:"description_ne,omitempty"`
	Category        string          `json:"category,omitempty"`
	Difficulty      string          `json:"difficulty,omitempty"`
	DurationMinutes *int            `json:"duration_minutes,omitempty"`
	Exercises       json.RawMessage `json:"exercises"`
	IsPreset        bool            `json:"is_preset"`
	CreatedBy       *string         `json:"created_by,omitempty"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
}

// WorkoutLog represents a logged workout session.
type WorkoutLog struct {
	ID                string          `json:"id"`
	UserID            string          `json:"user_id"`
	WorkoutTemplateID *string         `json:"workout_template_id,omitempty"`
	OrgID             string          `json:"organization_id"`
	Exercises         json.RawMessage `json:"exercises"`
	DurationMinutes   *int            `json:"duration_minutes,omitempty"`
	Notes             string          `json:"notes,omitempty"`
	LoggedAt          time.Time       `json:"logged_at"`
}

// MemberWorkoutPlan represents an assigned workout plan for a member.
type MemberWorkoutPlan struct {
	ID                string          `json:"id"`
	UserID            string          `json:"user_id"`
	OrgID             string          `json:"organization_id"`
	WorkoutTemplateID string          `json:"workout_template_id"`
	AssignedBy        *string         `json:"assigned_by,omitempty"`
	AssignedAt        time.Time       `json:"assigned_at"`
	Template          *WorkoutTemplate `json:"template,omitempty"`
}

// Repository handles workout data persistence.
type Repository struct {
	db *pgxpool.Pool
}

// NewRepository creates a new workout repository.
func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

// --- Template CRUD ---

// ListTemplates retrieves templates for an org plus global presets.
func (r *Repository) ListTemplates(ctx context.Context, orgID string, limit, offset int) ([]WorkoutTemplate, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, organization_id, name, COALESCE(name_ne, ''), COALESCE(description, ''),
		        COALESCE(description_ne, ''), COALESCE(category, ''), COALESCE(difficulty, ''),
		        duration_minutes, exercises, is_preset, created_by, created_at, updated_at
		 FROM workout_templates
		 WHERE organization_id = $1 OR is_preset = true
		 ORDER BY is_preset DESC, created_at DESC
		 LIMIT $2 OFFSET $3`,
		orgID, limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("listing templates: %w", err)
	}
	defer rows.Close()

	var templates []WorkoutTemplate
	for rows.Next() {
		var t WorkoutTemplate
		if err := rows.Scan(&t.ID, &t.OrgID, &t.Name, &t.NameNe, &t.Description,
			&t.DescriptionNe, &t.Category, &t.Difficulty,
			&t.DurationMinutes, &t.Exercises, &t.IsPreset, &t.CreatedBy, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scanning template: %w", err)
		}
		templates = append(templates, t)
	}
	return templates, rows.Err()
}

// GetTemplate retrieves a single template by ID (must be in the org or a preset).
func (r *Repository) GetTemplate(ctx context.Context, orgID, templateID string) (*WorkoutTemplate, error) {
	var t WorkoutTemplate
	err := r.db.QueryRow(ctx,
		`SELECT id, organization_id, name, COALESCE(name_ne, ''), COALESCE(description, ''),
		        COALESCE(description_ne, ''), COALESCE(category, ''), COALESCE(difficulty, ''),
		        duration_minutes, exercises, is_preset, created_by, created_at, updated_at
		 FROM workout_templates
		 WHERE id = $1 AND (organization_id = $2 OR is_preset = true)`,
		templateID, orgID,
	).Scan(&t.ID, &t.OrgID, &t.Name, &t.NameNe, &t.Description,
		&t.DescriptionNe, &t.Category, &t.Difficulty,
		&t.DurationMinutes, &t.Exercises, &t.IsPreset, &t.CreatedBy, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("getting template: %w", err)
	}
	return &t, nil
}

// CreateTemplate inserts a new workout template.
func (r *Repository) CreateTemplate(ctx context.Context, orgID, name, nameNe, description, descriptionNe, category, difficulty string, durationMinutes *int, exercises json.RawMessage, createdBy string) (*WorkoutTemplate, error) {
	var t WorkoutTemplate
	err := r.db.QueryRow(ctx,
		`INSERT INTO workout_templates
		    (organization_id, name, name_ne, description, description_ne, category, difficulty, duration_minutes, exercises, created_by)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		 RETURNING id, organization_id, name, COALESCE(name_ne, ''), COALESCE(description, ''),
		           COALESCE(description_ne, ''), COALESCE(category, ''), COALESCE(difficulty, ''),
		           duration_minutes, exercises, is_preset, created_by, created_at, updated_at`,
		orgID, name, nilIfEmpty(nameNe), nilIfEmpty(description), nilIfEmpty(descriptionNe),
		nilIfEmpty(category), nilIfEmpty(difficulty), durationMinutes, exercises, createdBy,
	).Scan(&t.ID, &t.OrgID, &t.Name, &t.NameNe, &t.Description,
		&t.DescriptionNe, &t.Category, &t.Difficulty,
		&t.DurationMinutes, &t.Exercises, &t.IsPreset, &t.CreatedBy, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("creating template: %w", err)
	}
	return &t, nil
}

// UpdateTemplate modifies an existing template.
func (r *Repository) UpdateTemplate(ctx context.Context, orgID, templateID, name, nameNe, description, descriptionNe, category, difficulty string, durationMinutes *int, exercises json.RawMessage) (*WorkoutTemplate, error) {
	var t WorkoutTemplate
	err := r.db.QueryRow(ctx,
		`UPDATE workout_templates
		 SET name = $3, name_ne = $4, description = $5, description_ne = $6,
		     category = $7, difficulty = $8, duration_minutes = $9, exercises = $10
		 WHERE id = $1 AND organization_id = $2 AND is_preset = false
		 RETURNING id, organization_id, name, COALESCE(name_ne, ''), COALESCE(description, ''),
		           COALESCE(description_ne, ''), COALESCE(category, ''), COALESCE(difficulty, ''),
		           duration_minutes, exercises, is_preset, created_by, created_at, updated_at`,
		templateID, orgID, name, nilIfEmpty(nameNe), nilIfEmpty(description), nilIfEmpty(descriptionNe),
		nilIfEmpty(category), nilIfEmpty(difficulty), durationMinutes, exercises,
	).Scan(&t.ID, &t.OrgID, &t.Name, &t.NameNe, &t.Description,
		&t.DescriptionNe, &t.Category, &t.Difficulty,
		&t.DurationMinutes, &t.Exercises, &t.IsPreset, &t.CreatedBy, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("updating template: %w", err)
	}
	return &t, nil
}

// DeleteTemplate removes a template (only org-owned, not presets).
func (r *Repository) DeleteTemplate(ctx context.Context, orgID, templateID string) error {
	tag, err := r.db.Exec(ctx,
		`DELETE FROM workout_templates WHERE id = $1 AND organization_id = $2 AND is_preset = false`,
		templateID, orgID,
	)
	if err != nil {
		return fmt.Errorf("deleting template: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("template not found")
	}
	return nil
}

// --- Plan assignment ---

// AssignPlan assigns a workout template to a member (upserts on conflict).
func (r *Repository) AssignPlan(ctx context.Context, userID, orgID, templateID, assignedBy string) (*MemberWorkoutPlan, error) {
	var p MemberWorkoutPlan
	err := r.db.QueryRow(ctx,
		`INSERT INTO member_workout_plans (user_id, organization_id, workout_template_id, assigned_by)
		 VALUES ($1, $2, $3, $4)
		 ON CONFLICT (user_id, organization_id)
		 DO UPDATE SET workout_template_id = $3, assigned_by = $4, assigned_at = NOW()
		 RETURNING id, user_id, organization_id, workout_template_id, assigned_by, assigned_at`,
		userID, orgID, templateID, assignedBy,
	).Scan(&p.ID, &p.UserID, &p.OrgID, &p.WorkoutTemplateID, &p.AssignedBy, &p.AssignedAt)
	if err != nil {
		return nil, fmt.Errorf("assigning plan: %w", err)
	}
	return &p, nil
}

// GetPlan retrieves the assigned workout plan for a member, joined with the template.
func (r *Repository) GetPlan(ctx context.Context, userID, orgID string) (*MemberWorkoutPlan, error) {
	var p MemberWorkoutPlan
	var t WorkoutTemplate
	err := r.db.QueryRow(ctx,
		`SELECT p.id, p.user_id, p.organization_id, p.workout_template_id, p.assigned_by, p.assigned_at,
		        t.id, t.organization_id, t.name, COALESCE(t.name_ne, ''), COALESCE(t.description, ''),
		        COALESCE(t.description_ne, ''), COALESCE(t.category, ''), COALESCE(t.difficulty, ''),
		        t.duration_minutes, t.exercises, t.is_preset, t.created_by, t.created_at, t.updated_at
		 FROM member_workout_plans p
		 JOIN workout_templates t ON t.id = p.workout_template_id
		 WHERE p.user_id = $1 AND p.organization_id = $2`,
		userID, orgID,
	).Scan(&p.ID, &p.UserID, &p.OrgID, &p.WorkoutTemplateID, &p.AssignedBy, &p.AssignedAt,
		&t.ID, &t.OrgID, &t.Name, &t.NameNe, &t.Description,
		&t.DescriptionNe, &t.Category, &t.Difficulty,
		&t.DurationMinutes, &t.Exercises, &t.IsPreset, &t.CreatedBy, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("getting plan: %w", err)
	}
	p.Template = &t
	return &p, nil
}

// --- Workout logging ---

// CreateLog inserts a new workout log.
func (r *Repository) CreateLog(ctx context.Context, userID, orgID string, templateID *string, exercises json.RawMessage, durationMinutes *int, notes string) (*WorkoutLog, error) {
	var l WorkoutLog
	err := r.db.QueryRow(ctx,
		`INSERT INTO workout_logs (user_id, organization_id, workout_template_id, exercises, duration_minutes, notes)
		 VALUES ($1, $2, $3, $4, $5, $6)
		 RETURNING id, user_id, workout_template_id, organization_id, exercises, duration_minutes,
		           COALESCE(notes, ''), logged_at`,
		userID, orgID, templateID, exercises, durationMinutes, nilIfEmpty(notes),
	).Scan(&l.ID, &l.UserID, &l.WorkoutTemplateID, &l.OrgID, &l.Exercises, &l.DurationMinutes,
		&l.Notes, &l.LoggedAt)
	if err != nil {
		return nil, fmt.Errorf("creating workout log: %w", err)
	}
	return &l, nil
}

// ListLogs retrieves workout logs for a user with optional date filtering.
func (r *Repository) ListLogs(ctx context.Context, userID, orgID string, from, to *time.Time, limit, offset int) ([]WorkoutLog, error) {
	query := `SELECT id, user_id, workout_template_id, organization_id, exercises, duration_minutes,
	                 COALESCE(notes, ''), logged_at
	          FROM workout_logs
	          WHERE user_id = $1 AND organization_id = $2`
	args := []interface{}{userID, orgID}
	argIdx := 3

	if from != nil {
		query += fmt.Sprintf(" AND logged_at >= $%d", argIdx)
		args = append(args, *from)
		argIdx++
	}
	if to != nil {
		query += fmt.Sprintf(" AND logged_at <= $%d", argIdx)
		args = append(args, *to)
		argIdx++
	}

	query += fmt.Sprintf(" ORDER BY logged_at DESC LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("listing workout logs: %w", err)
	}
	defer rows.Close()

	var logs []WorkoutLog
	for rows.Next() {
		var l WorkoutLog
		if err := rows.Scan(&l.ID, &l.UserID, &l.WorkoutTemplateID, &l.OrgID, &l.Exercises,
			&l.DurationMinutes, &l.Notes, &l.LoggedAt); err != nil {
			return nil, fmt.Errorf("scanning workout log: %w", err)
		}
		logs = append(logs, l)
	}
	return logs, rows.Err()
}

func nilIfEmpty(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}
