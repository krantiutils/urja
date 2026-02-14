package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// RunMigrations executes SQL migrations against the database.
// For now this runs the embedded migration directly. In production,
// use a proper migration tool like golang-migrate.
func RunMigrations(ctx context.Context, pool *pgxpool.Pool) error {
	migrations := []struct {
		name string
		sql  string
	}{
		{
			name: "001_create_users",
			sql: `
				CREATE TABLE IF NOT EXISTS users (
					id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
					phone VARCHAR(15) NOT NULL UNIQUE,
					name VARCHAR(255) NOT NULL DEFAULT '',
					role VARCHAR(50) NOT NULL DEFAULT 'member',
					is_active BOOLEAN NOT NULL DEFAULT true,
					created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
					updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
				);

				CREATE INDEX IF NOT EXISTS idx_users_phone ON users(phone);
			`,
		},
	}

	for _, m := range migrations {
		if _, err := pool.Exec(ctx, m.sql); err != nil {
			return fmt.Errorf("running migration %s: %w", m.name, err)
		}
	}

	return nil
}
