CREATE TABLE IF NOT EXISTS workout_programs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(200) NOT NULL,
    name_ne VARCHAR(200),
    description TEXT,
    program_type VARCHAR(30) CHECK (program_type IN ('hiit', 'bodyweight', 'strength', 'flexibility', 'mixed')),
    difficulty VARCHAR(20) NOT NULL DEFAULT 'beginner'
        CHECK (difficulty IN ('beginner', 'intermediate', 'advanced')),
    duration_weeks INT NOT NULL DEFAULT 4,
    equipment VARCHAR(50) NOT NULL DEFAULT 'none',
    goal VARCHAR(30),
    image_url TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS program_days (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    program_id UUID NOT NULL REFERENCES workout_programs(id) ON DELETE CASCADE,
    week_number INT NOT NULL,
    day_number INT NOT NULL,
    day_name VARCHAR(100),
    exercises JSONB NOT NULL DEFAULT '[]',
    rest_day BOOLEAN DEFAULT false,
    UNIQUE(program_id, week_number, day_number)
);

CREATE TABLE IF NOT EXISTS user_program_enrollments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    program_id UUID NOT NULL REFERENCES workout_programs(id) ON DELETE CASCADE,
    started_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    current_week INT NOT NULL DEFAULT 1,
    current_day INT NOT NULL DEFAULT 1,
    status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'paused', 'completed', 'abandoned')),
    UNIQUE(user_id, program_id)
);

CREATE INDEX IF NOT EXISTS idx_program_days_program ON program_days(program_id);
CREATE INDEX IF NOT EXISTS idx_enrollments_user ON user_program_enrollments(user_id);
CREATE INDEX IF NOT EXISTS idx_enrollments_status ON user_program_enrollments(user_id, status);
