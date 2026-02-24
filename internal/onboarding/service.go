package onboarding

import (
	"context"
	"fmt"
	"log/slog"
	"math"
	"strings"

	"github.com/google/uuid"
)

// UserProfile is the response shape after onboarding.
type UserProfile struct {
	ID                  string  `json:"id"`
	Phone               string  `json:"phone"`
	Name                string  `json:"name"`
	UserType            string  `json:"user_type"`
	OnboardingCompleted bool    `json:"onboarding_completed"`
	GoalType            string  `json:"goal_type,omitempty"`
	DailyWaterGoalML    int     `json:"daily_water_goal_ml"`
	DailyStepGoal       int     `json:"daily_step_goal"`
	CalorieGoal         float64 `json:"calorie_goal,omitempty"`
	ProteinGoalG        float64 `json:"protein_goal_g,omitempty"`
	CarbsGoalG          float64 `json:"carbs_goal_g,omitempty"`
	FatGoalG            float64 `json:"fat_goal_g,omitempty"`
}

// OnboardingInput holds the data submitted during onboarding.
type OnboardingInput struct {
	UserType      string  `json:"user_type"`
	Name          string  `json:"name"`
	GoalType      string  `json:"goal_type"`
	WeightKg      float64 `json:"weight_kg"`
	HeightCm      float64 `json:"height_cm"`
	Age           int     `json:"age"`
	Gender        string  `json:"gender"`
	ActivityLevel string  `json:"activity_level"`
}

// JoinGymInput holds the org_id for joining a gym.
type JoinGymInput struct {
	OrgID string `json:"org_id"`
}

// LeaveGymInput holds the org_id for leaving a gym.
type LeaveGymInput struct {
	OrgID string `json:"org_id"`
}

var validUserTypes = map[string]bool{
	"gym_member":      true,
	"fitness_tracker": true,
	"calorie_tracker": true,
}

var validGoalTypes = map[string]bool{
	"lose_weight":    true,
	"build_muscle":   true,
	"stay_fit":       true,
	"general_health": true,
}

var activityMultipliers = map[string]float64{
	"sedentary":   1.2,
	"light":       1.375,
	"moderate":    1.55,
	"active":      1.725,
	"very_active": 1.9,
}

var goalCalorieAdjustments = map[string]float64{
	"lose_weight":    -500,
	"build_muscle":   300,
	"stay_fit":       0,
	"general_health": 0,
}

// Service handles onboarding business logic.
type Service struct {
	repo   *Repository
	logger *slog.Logger
}

// NewService creates a new onboarding service.
func NewService(repo *Repository, logger *slog.Logger) *Service {
	return &Service{repo: repo, logger: logger}
}

// CompleteOnboarding processes the onboarding form submission.
func (s *Service) CompleteOnboarding(ctx context.Context, userID string, input *OnboardingInput) (*UserProfile, error) {
	if !validUserTypes[input.UserType] {
		return nil, fmt.Errorf("user_type must be gym_member, fitness_tracker, or calorie_tracker")
	}
	name := strings.TrimSpace(input.Name)
	if name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if input.GoalType != "" && !validGoalTypes[input.GoalType] {
		return nil, fmt.Errorf("goal_type must be lose_weight, build_muscle, stay_fit, or general_health")
	}

	// Update user profile
	if err := s.repo.CompleteOnboarding(ctx, userID, input.UserType, name, input.GoalType); err != nil {
		return nil, err
	}

	// Calculate and set nutrition goal if body stats provided
	var calorieGoal, proteinG, carbsG, fatG float64
	if input.WeightKg > 0 && input.HeightCm > 0 && input.Age > 0 && input.Gender != "" {
		multiplier, ok := activityMultipliers[input.ActivityLevel]
		if !ok {
			multiplier = 1.55 // default to moderate
		}

		adjustment := goalCalorieAdjustments[input.GoalType]

		// Mifflin-St Jeor equation
		var offset float64
		if input.Gender == "male" {
			offset = 5
		} else {
			offset = -161
		}
		bmr := 10*input.WeightKg + 6.25*input.HeightCm - 5*float64(input.Age) + offset
		tdee := bmr * multiplier
		calorieGoal = math.Round(tdee + adjustment)

		// Macros: Protein 30%, Carbs 40%, Fat 30%
		proteinG = math.Round(calorieGoal * 0.30 / 4)
		carbsG = math.Round(calorieGoal * 0.40 / 4)
		fatG = math.Round(calorieGoal * 0.30 / 9)

		if err := s.repo.UpsertNutritionGoal(ctx, userID, calorieGoal, proteinG, carbsG, fatG,
			input.WeightKg, input.HeightCm, input.Age, input.Gender, input.ActivityLevel, input.GoalType); err != nil {
			s.logger.Error("failed to set nutrition goal during onboarding", "error", err, "user_id", userID)
			// Non-fatal: continue with onboarding
		}
	}

	profile, err := s.repo.GetUserProfile(ctx, userID)
	if err != nil {
		return nil, err
	}

	profile.CalorieGoal = calorieGoal
	profile.ProteinGoalG = proteinG
	profile.CarbsGoalG = carbsG
	profile.FatGoalG = fatG

	s.logger.Info("onboarding completed", "user_id", userID, "user_type", input.UserType, "goal_type", input.GoalType)
	return profile, nil
}

// JoinGym adds the user to a gym organization.
// Accepts either a UUID org_id or a slug (e.g. "kathmandu-fitness-hub").
func (s *Service) JoinGym(ctx context.Context, userID string, input *JoinGymInput) error {
	if input.OrgID == "" {
		return fmt.Errorf("org_id is required")
	}

	orgID := input.OrgID
	// If input doesn't look like a UUID, treat it as a slug and resolve it
	if _, err := uuid.Parse(orgID); err != nil {
		resolvedID, slugErr := s.repo.GetOrgIDBySlug(ctx, orgID)
		if slugErr != nil {
			return fmt.Errorf("gym not found: %s", orgID)
		}
		orgID = resolvedID
	}

	if err := s.repo.JoinGym(ctx, userID, orgID); err != nil {
		return err
	}
	s.logger.Info("user joined gym", "user_id", userID, "org_id", orgID)
	return nil
}

// LeaveGym removes the user from a gym organization.
func (s *Service) LeaveGym(ctx context.Context, userID string, input *LeaveGymInput) error {
	if input.OrgID == "" {
		return fmt.Errorf("org_id is required")
	}
	if _, err := uuid.Parse(input.OrgID); err != nil {
		return fmt.Errorf("invalid org_id format")
	}
	if err := s.repo.LeaveGym(ctx, userID, input.OrgID); err != nil {
		return err
	}
	s.logger.Info("user left gym", "user_id", userID, "org_id", input.OrgID)
	return nil
}
