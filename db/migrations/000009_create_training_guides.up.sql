CREATE TABLE training_guides (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    title_ne VARCHAR(255),
    content TEXT NOT NULL,
    content_ne TEXT,
    category VARCHAR(100),
    cover_image_url TEXT,
    author_id UUID REFERENCES users(id) ON DELETE SET NULL,
    is_published BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_training_guides_category ON training_guides(category) WHERE category IS NOT NULL;
CREATE INDEX idx_training_guides_published ON training_guides(is_published) WHERE is_published = true;

CREATE TRIGGER set_training_guides_updated_at
    BEFORE UPDATE ON training_guides
    FOR EACH ROW
    EXECUTE FUNCTION trigger_set_updated_at();
