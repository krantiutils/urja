package database

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// NewRedisClient creates and verifies a Redis connection.
func NewRedisClient(ctx context.Context, redisURL string) (*redis.Client, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("parsing redis URL: %w", err)
	}

	client := redis.NewClient(opts)

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("connecting to redis: %w", err)
	}

	return client, nil
}
