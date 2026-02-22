CREATE TABLE IF NOT EXISTS nutrition_streaks (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    current_streak INT NOT NULL DEFAULT 0,
    longest_streak INT NOT NULL DEFAULT 0,
    last_log_date DATE,
    updated_at TIMESTAMPTZ DEFAULT NOW()
);
