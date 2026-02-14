CREATE TABLE health_metrics (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    member_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    metric_type VARCHAR(50) NOT NULL
        CHECK (metric_type IN ('bmi', 'whr', 'weight', 'body_fat', 'measurements')),
    value JSONB NOT NULL,
    recorded_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_health_metrics_member_id ON health_metrics(member_id);
CREATE INDEX idx_health_metrics_member_type ON health_metrics(member_id, metric_type, recorded_at DESC);

CREATE TABLE progress_photos (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    member_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    photo_url TEXT NOT NULL,
    category VARCHAR(50)
        CHECK (category IN ('front', 'side', 'back')),
    notes TEXT,
    taken_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_progress_photos_member_id ON progress_photos(member_id);
CREATE INDEX idx_progress_photos_member_date ON progress_photos(member_id, taken_at DESC);
