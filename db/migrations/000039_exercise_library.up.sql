CREATE TABLE IF NOT EXISTS exercises (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(200) NOT NULL,
    name_ne VARCHAR(200),
    description TEXT,
    muscle_groups TEXT[] NOT NULL DEFAULT '{}',
    equipment VARCHAR(50) NOT NULL DEFAULT 'none'
        CHECK (equipment IN ('none', 'dumbbells', 'barbell', 'resistance_band', 'pull_up_bar', 'machine', 'kettlebell')),
    exercise_type VARCHAR(30) NOT NULL
        CHECK (exercise_type IN ('cardio', 'strength', 'flexibility', 'hiit')),
    difficulty VARCHAR(20) NOT NULL DEFAULT 'beginner'
        CHECK (difficulty IN ('beginner', 'intermediate', 'advanced')),
    calories_per_minute DECIMAL(5,2) DEFAULT 0,
    instructions TEXT,
    image_url TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_exercises_type ON exercises(exercise_type);
CREATE INDEX IF NOT EXISTS idx_exercises_equipment ON exercises(equipment);
CREATE INDEX IF NOT EXISTS idx_exercises_difficulty ON exercises(difficulty);
