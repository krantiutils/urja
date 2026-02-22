package onboarding

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository handles onboarding data persistence.
type Repository struct {
	db *pgxpool.Pool
}

// NewRepository creates a new onboarding repository.
func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

// CompleteOnboarding updates the user's profile with onboarding data.
func (r *Repository) CompleteOnboarding(ctx context.Context, userID, userType, name, goalType string) error {
	tag, err := r.db.Exec(ctx,
		`UPDATE users SET
			user_type = $2,
			name = $3,
			goal_type = $4,
			onboarding_completed = true
		 WHERE id = $1`,
		userID, userType, name, nilIfEmpty(goalType),
	)
	if err != nil {
		return fmt.Errorf("completing onboarding: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}

// UpsertNutritionGoal creates or updates a nutrition goal for a user (org_id can be nil).
func (r *Repository) UpsertNutritionGoal(ctx context.Context, userID string, calorieGoal, proteinG, carbsG, fatG, weightKg, heightCm float64, age int, gender, activityLevel, goalType string) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO nutrition_goals (user_id, organization_id, calorie_goal, protein_goal_g, carbs_goal_g, fat_goal_g, weight_kg, height_cm, age, gender, activity_level, goal_type)
		 VALUES ($1, NULL, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		 ON CONFLICT ON CONSTRAINT nutrition_goals_user_org_unique
		 DO UPDATE SET calorie_goal = $2, protein_goal_g = $3, carbs_goal_g = $4, fat_goal_g = $5,
		               weight_kg = $6, height_cm = $7, age = $8, gender = $9,
		               activity_level = $10, goal_type = $11`,
		userID, calorieGoal, proteinG, carbsG, fatG, weightKg, heightCm, age, gender, activityLevel, goalType,
	)
	if err != nil {
		return fmt.Errorf("upserting nutrition goal: %w", err)
	}
	return nil
}

// JoinGym creates an organization_members record for the user.
func (r *Repository) JoinGym(ctx context.Context, userID, orgID string) error {
	// Check if membership already exists
	var existingStatus string
	err := r.db.QueryRow(ctx,
		`SELECT status FROM organization_members WHERE user_id = $1 AND organization_id = $2`,
		userID, orgID,
	).Scan(&existingStatus)

	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return fmt.Errorf("checking membership: %w", err)
	}

	if err == nil {
		if existingStatus == "active" {
			return fmt.Errorf("already a member of this gym")
		}
		// Re-activate
		_, err = r.db.Exec(ctx,
			`UPDATE organization_members SET status = 'active', joined_at = NOW()
			 WHERE user_id = $1 AND organization_id = $2`,
			userID, orgID,
		)
		if err != nil {
			return fmt.Errorf("reactivating membership: %w", err)
		}
		return nil
	}

	// Create new membership
	_, err = r.db.Exec(ctx,
		`INSERT INTO organization_members (user_id, organization_id, role) VALUES ($1, $2, 'member')`,
		userID, orgID,
	)
	if err != nil {
		return fmt.Errorf("joining gym: %w", err)
	}
	return nil
}

// LeaveGym sets the user's organization_members status to 'left'.
func (r *Repository) LeaveGym(ctx context.Context, userID, orgID string) error {
	tag, err := r.db.Exec(ctx,
		`UPDATE organization_members SET status = 'left'
		 WHERE user_id = $1 AND organization_id = $2 AND status = 'active'`,
		userID, orgID,
	)
	if err != nil {
		return fmt.Errorf("leaving gym: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("not an active member of this gym")
	}
	return nil
}

// GetUserProfile retrieves basic user profile for the onboarding response.
func (r *Repository) GetUserProfile(ctx context.Context, userID string) (*UserProfile, error) {
	var p UserProfile
	err := r.db.QueryRow(ctx,
		`SELECT id, phone, name, COALESCE(user_type, 'gym_member'), onboarding_completed,
		        COALESCE(goal_type, ''), COALESCE(daily_water_goal_ml, 2500), COALESCE(daily_step_goal, 8000)
		 FROM users WHERE id = $1`,
		userID,
	).Scan(&p.ID, &p.Phone, &p.Name, &p.UserType, &p.OnboardingCompleted,
		&p.GoalType, &p.DailyWaterGoalML, &p.DailyStepGoal)
	if err != nil {
		return nil, fmt.Errorf("getting user profile: %w", err)
	}
	return &p, nil
}

func nilIfEmpty(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}
