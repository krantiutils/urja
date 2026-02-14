package workout

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"time"
)

// Service handles workout business logic.
type Service struct {
	repo   *Repository
	logger *slog.Logger
}

// NewService creates a new workout service.
func NewService(repo *Repository, logger *slog.Logger) *Service {
	return &Service{repo: repo, logger: logger}
}

// --- Template operations ---

// ListTemplates returns templates for an org plus global presets.
func (s *Service) ListTemplates(ctx context.Context, orgID string, limit, offset int) ([]WorkoutTemplate, error) {
	if limit <= 0 || limit > 50 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	return s.repo.ListTemplates(ctx, orgID, limit, offset)
}

// GetTemplate returns a single template.
func (s *Service) GetTemplate(ctx context.Context, orgID, templateID string) (*WorkoutTemplate, error) {
	return s.repo.GetTemplate(ctx, orgID, templateID)
}

var validDifficulties = map[string]bool{
	"beginner":     true,
	"intermediate": true,
	"advanced":     true,
}

// CreateTemplateInput holds input for creating a template.
type CreateTemplateInput struct {
	Name            string          `json:"name"`
	NameNe          string          `json:"name_ne"`
	Description     string          `json:"description"`
	DescriptionNe   string          `json:"description_ne"`
	Category        string          `json:"category"`
	Difficulty      string          `json:"difficulty"`
	DurationMinutes *int            `json:"duration_minutes"`
	Exercises       json.RawMessage `json:"exercises"`
}

// CreateTemplate creates a new org-scoped workout template.
func (s *Service) CreateTemplate(ctx context.Context, orgID string, input *CreateTemplateInput, createdBy string) (*WorkoutTemplate, error) {
	name := strings.TrimSpace(input.Name)
	if name == "" {
		return nil, fmt.Errorf("name is required")
	}

	if input.Difficulty != "" && !validDifficulties[input.Difficulty] {
		return nil, fmt.Errorf("difficulty must be beginner, intermediate, or advanced")
	}

	if input.DurationMinutes != nil && *input.DurationMinutes <= 0 {
		return nil, fmt.Errorf("duration_minutes must be positive")
	}

	exercises := input.Exercises
	if len(exercises) == 0 {
		exercises = json.RawMessage(`[]`)
	}

	t, err := s.repo.CreateTemplate(ctx, orgID, name, strings.TrimSpace(input.NameNe),
		strings.TrimSpace(input.Description), strings.TrimSpace(input.DescriptionNe),
		strings.TrimSpace(input.Category), input.Difficulty, input.DurationMinutes, exercises, createdBy)
	if err != nil {
		return nil, err
	}

	s.logger.Info("workout template created", "template_id", t.ID, "org_id", orgID, "created_by", createdBy)
	return t, nil
}

// UpdateTemplateInput holds input for updating a template.
type UpdateTemplateInput struct {
	Name            string          `json:"name"`
	NameNe          string          `json:"name_ne"`
	Description     string          `json:"description"`
	DescriptionNe   string          `json:"description_ne"`
	Category        string          `json:"category"`
	Difficulty      string          `json:"difficulty"`
	DurationMinutes *int            `json:"duration_minutes"`
	Exercises       json.RawMessage `json:"exercises"`
}

// UpdateTemplate modifies an existing org-scoped template.
func (s *Service) UpdateTemplate(ctx context.Context, orgID, templateID string, input *UpdateTemplateInput) (*WorkoutTemplate, error) {
	name := strings.TrimSpace(input.Name)
	if name == "" {
		return nil, fmt.Errorf("name is required")
	}

	if input.Difficulty != "" && !validDifficulties[input.Difficulty] {
		return nil, fmt.Errorf("difficulty must be beginner, intermediate, or advanced")
	}

	if input.DurationMinutes != nil && *input.DurationMinutes <= 0 {
		return nil, fmt.Errorf("duration_minutes must be positive")
	}

	exercises := input.Exercises
	if len(exercises) == 0 {
		exercises = json.RawMessage(`[]`)
	}

	t, err := s.repo.UpdateTemplate(ctx, orgID, templateID, name, strings.TrimSpace(input.NameNe),
		strings.TrimSpace(input.Description), strings.TrimSpace(input.DescriptionNe),
		strings.TrimSpace(input.Category), input.Difficulty, input.DurationMinutes, exercises)
	if err != nil {
		return nil, err
	}

	s.logger.Info("workout template updated", "template_id", templateID, "org_id", orgID)
	return t, nil
}

// DeleteTemplate removes an org-scoped template.
func (s *Service) DeleteTemplate(ctx context.Context, orgID, templateID string) error {
	if err := s.repo.DeleteTemplate(ctx, orgID, templateID); err != nil {
		return err
	}
	s.logger.Info("workout template deleted", "template_id", templateID, "org_id", orgID)
	return nil
}

// --- Plan assignment ---

// AssignPlan assigns a workout template to a member.
func (s *Service) AssignPlan(ctx context.Context, memberID, orgID, templateID, assignedBy string) (*MemberWorkoutPlan, error) {
	if templateID == "" {
		return nil, fmt.Errorf("workout_template_id is required")
	}

	// Verify the template exists and is accessible to this org.
	_, err := s.repo.GetTemplate(ctx, orgID, templateID)
	if err != nil {
		return nil, fmt.Errorf("template not found or not accessible")
	}

	p, err := s.repo.AssignPlan(ctx, memberID, orgID, templateID, assignedBy)
	if err != nil {
		return nil, err
	}

	s.logger.Info("workout plan assigned", "member_id", memberID, "template_id", templateID,
		"org_id", orgID, "assigned_by", assignedBy)
	return p, nil
}

// GetPlan retrieves the assigned workout plan for a member.
func (s *Service) GetPlan(ctx context.Context, userID, orgID string) (*MemberWorkoutPlan, error) {
	return s.repo.GetPlan(ctx, userID, orgID)
}

// --- Workout logging ---

// CreateLogInput holds input for logging a workout.
type CreateLogInput struct {
	WorkoutTemplateID *string         `json:"workout_template_id"`
	OrgID             string          `json:"organization_id"`
	Exercises         json.RawMessage `json:"exercises"`
	DurationMinutes   *int            `json:"duration_minutes"`
	Notes             string          `json:"notes"`
}

// CreateLog logs a workout session.
func (s *Service) CreateLog(ctx context.Context, userID string, input *CreateLogInput) (*WorkoutLog, error) {
	if input.OrgID == "" {
		return nil, fmt.Errorf("organization_id is required")
	}

	exercises := input.Exercises
	if len(exercises) == 0 {
		exercises = json.RawMessage(`[]`)
	}

	if input.DurationMinutes != nil && *input.DurationMinutes <= 0 {
		return nil, fmt.Errorf("duration_minutes must be positive")
	}

	l, err := s.repo.CreateLog(ctx, userID, input.OrgID, input.WorkoutTemplateID, exercises,
		input.DurationMinutes, strings.TrimSpace(input.Notes))
	if err != nil {
		return nil, err
	}

	s.logger.Info("workout logged", "log_id", l.ID, "user_id", userID, "org_id", input.OrgID)
	return l, nil
}

// ListLogs retrieves workout logs for a user with optional date filtering.
func (s *Service) ListLogs(ctx context.Context, userID, orgID string, from, to *time.Time, limit, offset int) ([]WorkoutLog, error) {
	if limit <= 0 || limit > 50 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	return s.repo.ListLogs(ctx, userID, orgID, from, to, limit, offset)
}
