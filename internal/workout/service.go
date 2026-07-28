package workout

import (
	"context"
	"encoding/json"
	"errors"
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
	Goal            string          `json:"goal"`
	DaysPerWeek     string          `json:"days_per_week"`
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
		strings.TrimSpace(input.Category), input.Difficulty, input.DurationMinutes, exercises, createdBy,
		strings.TrimSpace(input.Goal), strings.TrimSpace(input.DaysPerWeek))
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
	Goal            string          `json:"goal"`
	DaysPerWeek     string          `json:"days_per_week"`
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
		strings.TrimSpace(input.Category), input.Difficulty, input.DurationMinutes, exercises,
		strings.TrimSpace(input.Goal), strings.TrimSpace(input.DaysPerWeek))
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

// --- Browse and self-assign ---

// BrowseTemplates returns templates filtered by goal and difficulty.
// If orgID is empty, only preset templates are returned.
func (s *Service) BrowseTemplates(ctx context.Context, orgID, goal, difficulty string, limit, offset int) ([]WorkoutTemplate, error) {
	if limit <= 0 || limit > 50 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	if orgID == "" {
		return s.repo.ListPresetTemplatesByFilters(ctx, goal, difficulty, limit, offset)
	}
	return s.repo.ListTemplatesByFilters(ctx, orgID, goal, difficulty, limit, offset)
}

// SelfAssignPlan allows a member to assign a workout template to themselves.
func (s *Service) SelfAssignPlan(ctx context.Context, userID, orgID, templateID string) (*MemberWorkoutPlan, error) {
	if templateID == "" {
		return nil, fmt.Errorf("workout_template_id is required")
	}

	// Verify the template exists and is accessible to this org.
	tmpl, err := s.repo.GetTemplate(ctx, orgID, templateID)
	if err != nil {
		return nil, fmt.Errorf("template not found or not accessible")
	}

	p, err := s.repo.SelfAssignPlan(ctx, userID, orgID, templateID, "self")
	if err != nil {
		return nil, err
	}

	p.Template = tmpl
	s.logger.Info("workout plan self-assigned", "user_id", userID, "template_id", templateID, "org_id", orgID)
	return p, nil
}

var validGoals = map[string]bool{
	"lose_weight":  true,
	"build_muscle": true,
	"stay_fit":     true,
}

var validDaysPerWeek = map[string]bool{
	"2-3": true,
	"4-5": true,
	"6-7": true,
}

// RecommendPlan uses a questionnaire to recommend and assign a workout plan.
func (s *Service) RecommendPlan(ctx context.Context, userID, orgID, goal, experience, daysPerWeek string) (*MemberWorkoutPlan, error) {
	if !validGoals[goal] {
		return nil, fmt.Errorf("goal must be lose_weight, build_muscle, or stay_fit")
	}
	if !validDifficulties[experience] {
		return nil, fmt.Errorf("experience must be beginner, intermediate, or advanced")
	}
	if !validDaysPerWeek[daysPerWeek] {
		return nil, fmt.Errorf("days_per_week must be 2-3, 4-5, or 6-7")
	}

	// Fetch all templates for this org (+ presets) with a high limit.
	templates, err := s.repo.ListTemplates(ctx, orgID, 1000, 0)
	if err != nil {
		return nil, fmt.Errorf("fetching templates: %w", err)
	}

	if len(templates) == 0 {
		return nil, fmt.Errorf("no templates available")
	}

	// Score each template and pick the best match.
	bestIdx := 0
	bestScore := -1
	for i, t := range templates {
		score := 0
		if t.Goal == goal {
			score += 3
		}
		if t.Difficulty == experience {
			score += 3
		}
		if t.DaysPerWeek == daysPerWeek {
			score += 2
		}
		if score > bestScore {
			bestScore = score
			bestIdx = i
		}
	}

	chosen := templates[bestIdx]

	// Assign using questionnaire method.
	p, err := s.repo.SelfAssignPlan(ctx, userID, orgID, chosen.ID, "questionnaire")
	if err != nil {
		return nil, err
	}
	p.Template = &chosen

	s.logger.Info("workout plan recommended and assigned",
		"user_id", userID, "template_id", chosen.ID, "org_id", orgID,
		"goal", goal, "experience", experience, "days_per_week", daysPerWeek,
		"score", bestScore)
	return p, nil
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

// ErrOrgRequired is returned when a member's organization cannot be inferred.
var ErrOrgRequired = errors.New("organization_id is required")

// ResolveOrg determines which organization a self-service request belongs to.
//
// A member who trains at one gym should not have to name it, and older clients
// omit the parameter entirely — but the alternative, writing the row with a
// NULL organization, orphans it: it belongs to no gym, so no gym's queries ever
// return it again. So an omitted parameter falls back to the member's sole
// active membership, and is an error only when there is genuinely a choice to
// make.
//
// A supplied organization is always verified against real membership, so naming
// someone else's gym cannot write into it.
func (s *Service) ResolveOrg(ctx context.Context, userID, requested string) (string, error) {
	orgs, err := s.repo.ActiveOrgIDs(ctx, userID)
	if err != nil {
		return "", err
	}

	if requested != "" {
		for _, id := range orgs {
			if id == requested {
				return requested, nil
			}
		}
		return "", fmt.Errorf("%w: not an active member of that organization", ErrOrgRequired)
	}

	switch len(orgs) {
	case 1:
		return orgs[0], nil
	case 0:
		return "", fmt.Errorf("%w: you are not an active member of any gym", ErrOrgRequired)
	default:
		return "", fmt.Errorf("%w: you belong to several gyms, so it must be specified", ErrOrgRequired)
	}
}
