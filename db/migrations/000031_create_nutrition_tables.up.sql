-- Create nutrition tables: food_items, food_logs, meal_templates, nutrition_goals.

-- Enable pg_trgm extension for fuzzy food search.
CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- Food items catalog (global items have organization_id=NULL).
CREATE TABLE IF NOT EXISTS food_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID REFERENCES organizations(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    name_ne VARCHAR(255),
    category VARCHAR(100),
    calories_per_100g NUMERIC(8,2),
    protein_per_100g NUMERIC(8,2),
    carbs_per_100g NUMERIC(8,2),
    fat_per_100g NUMERIC(8,2),
    fiber_per_100g NUMERIC(8,2),
    serving_size_g NUMERIC(8,2),
    serving_label VARCHAR(100),
    serving_label_ne VARCHAR(100),
    is_verified BOOLEAN DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_food_items_org_category ON food_items(organization_id, category);
CREATE INDEX IF NOT EXISTS idx_food_items_name_trgm ON food_items USING GIN (name gin_trgm_ops);

DROP TRIGGER IF EXISTS set_food_items_updated_at ON food_items;
CREATE TRIGGER set_food_items_updated_at
    BEFORE UPDATE ON food_items
    FOR EACH ROW
    EXECUTE FUNCTION trigger_set_updated_at();

-- Food logs: individual meal entries.
CREATE TABLE IF NOT EXISTS food_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    food_item_id UUID NOT NULL REFERENCES food_items(id) ON DELETE CASCADE,
    meal_type VARCHAR(20) CHECK (meal_type IN ('breakfast', 'lunch', 'dinner', 'snack')),
    quantity_grams NUMERIC(8,2) NOT NULL,
    calories NUMERIC(8,2),
    protein NUMERIC(8,2),
    carbs NUMERIC(8,2),
    fat NUMERIC(8,2),
    logged_date DATE NOT NULL,
    logged_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    notes TEXT
);

CREATE INDEX IF NOT EXISTS idx_food_logs_user_date ON food_logs(user_id, logged_date, organization_id);

-- Meal templates: saved meal combinations.
CREATE TABLE IF NOT EXISTS meal_templates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    name_ne VARCHAR(255),
    meal_type VARCHAR(20) CHECK (meal_type IN ('breakfast', 'lunch', 'dinner', 'snack')),
    items JSONB NOT NULL DEFAULT '[]',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_meal_templates_user_org ON meal_templates(user_id, organization_id);

-- Nutrition goals: per-user daily targets.
CREATE TABLE IF NOT EXISTS nutrition_goals (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    calorie_goal NUMERIC(8,2),
    protein_goal_g NUMERIC(8,2),
    carbs_goal_g NUMERIC(8,2),
    fat_goal_g NUMERIC(8,2),
    weight_kg NUMERIC(5,2),
    height_cm NUMERIC(5,2),
    age INT,
    gender VARCHAR(10),
    activity_level VARCHAR(20),
    goal_type VARCHAR(20),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (user_id, organization_id)
);

CREATE INDEX IF NOT EXISTS idx_nutrition_goals_user_org ON nutrition_goals(user_id, organization_id);

DROP TRIGGER IF EXISTS set_nutrition_goals_updated_at ON nutrition_goals;
CREATE TRIGGER set_nutrition_goals_updated_at
    BEFORE UPDATE ON nutrition_goals
    FOR EACH ROW
    EXECUTE FUNCTION trigger_set_updated_at();
