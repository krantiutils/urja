package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort       int
	Environment      string
	DatabaseURL      string
	RedisURL         string
	JWTSecret        string
	JWTRefreshSecret string
	AakashSMSToken   string
}

func Load() (*Config, error) {
	// Load .env file if it exists; ignore error (env vars may be set externally)
	_ = godotenv.Load()

	cfg := &Config{
		ServerPort:       8080,
		Environment:      "development",
		DatabaseURL:      os.Getenv("DATABASE_URL"),
		RedisURL:         os.Getenv("REDIS_URL"),
		JWTSecret:        os.Getenv("JWT_SECRET"),
		JWTRefreshSecret: os.Getenv("JWT_REFRESH_SECRET"),
		AakashSMSToken:   os.Getenv("AAKASH_SMS_TOKEN"),
	}

	if port := os.Getenv("SERVER_PORT"); port != "" {
		p, err := strconv.Atoi(port)
		if err != nil {
			return nil, fmt.Errorf("invalid SERVER_PORT %q: %w", port, err)
		}
		cfg.ServerPort = p
	}

	if env := os.Getenv("ENVIRONMENT"); env != "" {
		cfg.Environment = env
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) validate() error {
	if c.DatabaseURL == "" {
		return fmt.Errorf("DATABASE_URL is required")
	}
	if c.RedisURL == "" {
		return fmt.Errorf("REDIS_URL is required")
	}
	if c.JWTSecret == "" {
		return fmt.Errorf("JWT_SECRET is required")
	}
	if c.JWTRefreshSecret == "" {
		return fmt.Errorf("JWT_REFRESH_SECRET is required")
	}
	if c.JWTSecret == c.JWTRefreshSecret {
		return fmt.Errorf("JWT_SECRET and JWT_REFRESH_SECRET must be different")
	}
	if len(c.JWTSecret) < 32 {
		return fmt.Errorf("JWT_SECRET must be at least 32 characters")
	}
	if len(c.JWTRefreshSecret) < 32 {
		return fmt.Errorf("JWT_REFRESH_SECRET must be at least 32 characters")
	}
	if c.AakashSMSToken == "" {
		return fmt.Errorf("AAKASH_SMS_TOKEN is required")
	}
	return nil
}

func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}
