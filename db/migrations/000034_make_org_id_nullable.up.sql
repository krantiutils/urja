-- Make organization_id nullable across nutrition and workout tables
-- so independent users (fitness_tracker, calorie_tracker) can use features without a gym.

-- food_logs: drop NOT NULL
ALTER TABLE food_logs ALTER COLUMN organization_id DROP NOT NULL;

-- meal_templates: drop NOT NULL
ALTER TABLE meal_templates ALTER COLUMN organization_id DROP NOT NULL;

-- nutrition_goals: drop unique constraint, drop NOT NULL, recreate with COALESCE unique index
ALTER TABLE nutrition_goals DROP CONSTRAINT IF EXISTS nutrition_goals_user_id_organization_id_key;
ALTER TABLE nutrition_goals ALTER COLUMN organization_id DROP NOT NULL;
CREATE UNIQUE INDEX IF NOT EXISTS nutrition_goals_user_org_unique
    ON nutrition_goals (user_id, COALESCE(organization_id, '00000000-0000-0000-0000-000000000000'));

-- workout_logs: drop NOT NULL
ALTER TABLE workout_logs ALTER COLUMN organization_id DROP NOT NULL;

-- member_workout_plans: drop unique constraint, drop NOT NULL, recreate with COALESCE unique index
ALTER TABLE member_workout_plans DROP CONSTRAINT IF EXISTS member_workout_plans_user_id_organization_id_key;
ALTER TABLE member_workout_plans ALTER COLUMN organization_id DROP NOT NULL;
CREATE UNIQUE INDEX IF NOT EXISTS member_workout_plans_user_org_unique
    ON member_workout_plans (user_id, COALESCE(organization_id, '00000000-0000-0000-0000-000000000000'));

-- Update indexes to handle NULL org_id
DROP INDEX IF EXISTS idx_food_logs_user_date;
CREATE INDEX idx_food_logs_user_date ON food_logs(user_id, logged_date, COALESCE(organization_id, '00000000-0000-0000-0000-000000000000'));

DROP INDEX IF EXISTS idx_meal_templates_user_org;
CREATE INDEX idx_meal_templates_user_org ON meal_templates(user_id, COALESCE(organization_id, '00000000-0000-0000-0000-000000000000'));

DROP INDEX IF EXISTS idx_nutrition_goals_user_org;
CREATE INDEX idx_nutrition_goals_user_org ON nutrition_goals(user_id, COALESCE(organization_id, '00000000-0000-0000-0000-000000000000'));
