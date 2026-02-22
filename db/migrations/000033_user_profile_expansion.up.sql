-- Add user_type, onboarding state, and personal goals to users table.
ALTER TABLE users ADD COLUMN user_type VARCHAR(20) NOT NULL DEFAULT 'gym_member'
    CHECK (user_type IN ('gym_member', 'fitness_tracker', 'calorie_tracker'));
ALTER TABLE users ADD COLUMN onboarding_completed BOOLEAN NOT NULL DEFAULT false;
ALTER TABLE users ADD COLUMN goal_type VARCHAR(30)
    CHECK (goal_type IN ('lose_weight', 'build_muscle', 'stay_fit', 'general_health'));
ALTER TABLE users ADD COLUMN daily_water_goal_ml INT DEFAULT 2500;
ALTER TABLE users ADD COLUMN daily_step_goal INT DEFAULT 8000;
