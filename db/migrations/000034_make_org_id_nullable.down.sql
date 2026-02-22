-- Reverse: restore NOT NULL constraints and original unique constraints.
-- WARNING: This will fail if any rows have NULL organization_id.

-- Restore original indexes
DROP INDEX IF EXISTS idx_food_logs_user_date;
CREATE INDEX idx_food_logs_user_date ON food_logs(user_id, logged_date, organization_id);

DROP INDEX IF EXISTS idx_meal_templates_user_org;
CREATE INDEX idx_meal_templates_user_org ON meal_templates(user_id, organization_id);

DROP INDEX IF EXISTS idx_nutrition_goals_user_org;
CREATE INDEX idx_nutrition_goals_user_org ON nutrition_goals(user_id, organization_id);

-- member_workout_plans
ALTER TABLE member_workout_plans DROP CONSTRAINT IF EXISTS member_workout_plans_user_org_unique;
DELETE FROM member_workout_plans WHERE organization_id IS NULL;
ALTER TABLE member_workout_plans ALTER COLUMN organization_id SET NOT NULL;
ALTER TABLE member_workout_plans ADD CONSTRAINT member_workout_plans_user_id_organization_id_key
    UNIQUE (user_id, organization_id);

-- nutrition_goals
ALTER TABLE nutrition_goals DROP CONSTRAINT IF EXISTS nutrition_goals_user_org_unique;
DELETE FROM nutrition_goals WHERE organization_id IS NULL;
ALTER TABLE nutrition_goals ALTER COLUMN organization_id SET NOT NULL;
ALTER TABLE nutrition_goals ADD CONSTRAINT nutrition_goals_user_id_organization_id_key
    UNIQUE (user_id, organization_id);

-- workout_logs
DELETE FROM workout_logs WHERE organization_id IS NULL;
ALTER TABLE workout_logs ALTER COLUMN organization_id SET NOT NULL;

-- meal_templates
DELETE FROM meal_templates WHERE organization_id IS NULL;
ALTER TABLE meal_templates ALTER COLUMN organization_id SET NOT NULL;

-- food_logs
DELETE FROM food_logs WHERE organization_id IS NULL;
ALTER TABLE food_logs ALTER COLUMN organization_id SET NOT NULL;
