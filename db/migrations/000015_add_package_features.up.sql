ALTER TABLE packages ADD COLUMN features JSONB NOT NULL DEFAULT '[]'::jsonb;
