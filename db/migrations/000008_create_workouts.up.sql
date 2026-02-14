CREATE TABLE workout_templates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID REFERENCES organizations(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    name_ne VARCHAR(255),
    description TEXT,
    description_ne TEXT,
    category VARCHAR(100),
    difficulty VARCHAR(50)
        CHECK (difficulty IN ('beginner', 'intermediate', 'advanced')),
    duration_minutes INT CHECK (duration_minutes > 0),
    exercises JSONB NOT NULL DEFAULT '[]',
    is_preset BOOLEAN NOT NULL DEFAULT false,
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_workout_templates_org_id ON workout_templates(organization_id);
CREATE INDEX idx_workout_templates_preset ON workout_templates(is_preset) WHERE is_preset = true;
CREATE INDEX idx_workout_templates_category ON workout_templates(category) WHERE category IS NOT NULL;

CREATE TRIGGER set_workout_templates_updated_at
    BEFORE UPDATE ON workout_templates
    FOR EACH ROW
    EXECUTE FUNCTION trigger_set_updated_at();

CREATE TABLE workout_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    workout_template_id UUID REFERENCES workout_templates(id) ON DELETE SET NULL,
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    exercises JSONB NOT NULL DEFAULT '[]',
    duration_minutes INT,
    notes TEXT,
    logged_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_workout_logs_user_id ON workout_logs(user_id);
CREATE INDEX idx_workout_logs_org_id ON workout_logs(organization_id);
CREATE INDEX idx_workout_logs_user_date ON workout_logs(user_id, logged_at DESC);
