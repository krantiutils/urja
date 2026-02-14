package redis

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/redis/go-redis/v9"

	"github.com/urja-gym/urja/internal/config"
)

// New creates a new Redis client from the given config.
// It verifies connectivity with a ping before returning.
func New(ctx context.Context, cfg config.RedisConfig, logger *slog.Logger) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr(),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("pinging redis: %w", err)
	}

	logger.Info("redis connection established",
		"addr", cfg.Addr(),
		"db", cfg.DB,
	)

	return client, nil
}
