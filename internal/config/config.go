package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config holds all application configuration loaded from environment variables.
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	Auth     AuthConfig
	SMS      SMSConfig
	S3       S3Config
	Khalti   KhaltiConfig
}

type ServerConfig struct {
	Host            string
	Port            int
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	SSLMode  string
	MaxConns int32
	MinConns int32
}

type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

type AuthConfig struct {
	JWTSecret          string
	AccessTokenExpiry  time.Duration
	RefreshTokenExpiry time.Duration
	OTPExpiry          time.Duration
	OTPMaxAttempts     int
	OTPCooldown        time.Duration
	OTPHourlyLimit     int
	BcryptCost         int
}

type SMSConfig struct {
	AakashToken  string
	AakashAPIURL string
}

type S3Config struct {
	Endpoint  string
	Bucket    string
	AccessKey string
	SecretKey string
	Region    string
	UseSSL    bool
}

type KhaltiConfig struct {
	SecretKey  string
	BaseURL    string
	WebsiteURL string
	ReturnURL  string
}

// Load reads configuration from environment variables.
// It returns an error if any required variable is missing or invalid.
func Load() (*Config, error) {
	cfg := &Config{}

	// Server
	cfg.Server.Host = envOrDefault("SERVER_HOST", "0.0.0.0")
	port, err := envInt("SERVER_PORT", 8080)
	if err != nil {
		return nil, fmt.Errorf("SERVER_PORT: %w", err)
	}
	cfg.Server.Port = port
	cfg.Server.ReadTimeout = envDuration("SERVER_READ_TIMEOUT", 15*time.Second)
	cfg.Server.WriteTimeout = envDuration("SERVER_WRITE_TIMEOUT", 15*time.Second)
	cfg.Server.ShutdownTimeout = envDuration("SERVER_SHUTDOWN_TIMEOUT", 10*time.Second)

	// Database
	cfg.Database.Host = envRequired("DB_HOST")
	dbPort, err := envInt("DB_PORT", 5432)
	if err != nil {
		return nil, fmt.Errorf("DB_PORT: %w", err)
	}
	cfg.Database.Port = dbPort
	cfg.Database.User = envRequired("DB_USER")
	cfg.Database.Password = envRequired("DB_PASSWORD")
	cfg.Database.Name = envRequired("DB_NAME")
	cfg.Database.SSLMode = envOrDefault("DB_SSLMODE", "disable")
	maxConns, err := envInt("DB_MAX_CONNS", 25)
	if err != nil {
		return nil, fmt.Errorf("DB_MAX_CONNS: %w", err)
	}
	cfg.Database.MaxConns = int32(maxConns)
	minConns, err := envInt("DB_MIN_CONNS", 5)
	if err != nil {
		return nil, fmt.Errorf("DB_MIN_CONNS: %w", err)
	}
	cfg.Database.MinConns = int32(minConns)

	// Redis
	cfg.Redis.Host = envOrDefault("REDIS_HOST", "localhost")
	redisPort, err := envInt("REDIS_PORT", 6379)
	if err != nil {
		return nil, fmt.Errorf("REDIS_PORT: %w", err)
	}
	cfg.Redis.Port = redisPort
	cfg.Redis.Password = os.Getenv("REDIS_PASSWORD")
	redisDB, err := envInt("REDIS_DB", 0)
	if err != nil {
		return nil, fmt.Errorf("REDIS_DB: %w", err)
	}
	cfg.Redis.DB = redisDB

	// Auth
	cfg.Auth.JWTSecret = envRequired("JWT_SECRET")
	cfg.Auth.AccessTokenExpiry = envDuration("ACCESS_TOKEN_EXPIRY", 15*time.Minute)
	cfg.Auth.RefreshTokenExpiry = envDuration("REFRESH_TOKEN_EXPIRY", 7*24*time.Hour)
	cfg.Auth.OTPExpiry = envDuration("OTP_EXPIRY", 5*time.Minute)
	otpMaxAttempts, err := envInt("OTP_MAX_ATTEMPTS", 3)
	if err != nil {
		return nil, fmt.Errorf("OTP_MAX_ATTEMPTS: %w", err)
	}
	cfg.Auth.OTPMaxAttempts = otpMaxAttempts
	cfg.Auth.OTPCooldown = envDuration("OTP_COOLDOWN", 60*time.Second)
	otpHourlyLimit, err := envInt("OTP_HOURLY_LIMIT", 5)
	if err != nil {
		return nil, fmt.Errorf("OTP_HOURLY_LIMIT: %w", err)
	}
	cfg.Auth.OTPHourlyLimit = otpHourlyLimit
	bcryptCost, err := envInt("BCRYPT_COST", 12)
	if err != nil {
		return nil, fmt.Errorf("BCRYPT_COST: %w", err)
	}
	cfg.Auth.BcryptCost = bcryptCost

	// SMS
	cfg.SMS.AakashToken = envRequired("AAKASH_SMS_TOKEN")
	cfg.SMS.AakashAPIURL = envOrDefault("AAKASH_SMS_API_URL", "https://sms.aakashsms.com/sms/v3/send")

	// S3 / MinIO
	cfg.S3.Endpoint = envOrDefault("S3_ENDPOINT", "localhost:9000")
	cfg.S3.Bucket = envOrDefault("S3_BUCKET", "urja")
	cfg.S3.AccessKey = os.Getenv("S3_ACCESS_KEY")
	cfg.S3.SecretKey = os.Getenv("S3_SECRET_KEY")
	cfg.S3.Region = envOrDefault("S3_REGION", "us-east-1")
	cfg.S3.UseSSL = envOrDefault("S3_USE_SSL", "false") == "true"

	// Khalti
	cfg.Khalti.SecretKey = os.Getenv("KHALTI_SECRET_KEY")
	cfg.Khalti.BaseURL = envOrDefault("KHALTI_BASE_URL", "https://dev.khalti.com/api/v2")
	cfg.Khalti.WebsiteURL = os.Getenv("KHALTI_WEBSITE_URL")
	cfg.Khalti.ReturnURL = os.Getenv("KHALTI_RETURN_URL")

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("config validation: %w", err)
	}

	return cfg, nil
}

func (c *Config) validate() error {
	if c.Auth.JWTSecret == "" {
		return fmt.Errorf("JWT_SECRET is required")
	}
	if len(c.Auth.JWTSecret) < 32 {
		return fmt.Errorf("JWT_SECRET must be at least 32 characters")
	}
	if c.Database.Host == "" {
		return fmt.Errorf("DB_HOST is required")
	}
	if c.SMS.AakashToken == "" {
		return fmt.Errorf("AAKASH_SMS_TOKEN is required")
	}
	return nil
}

// DSN returns the PostgreSQL connection string.
func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.User, c.Password, c.Host, c.Port, c.Name, c.SSLMode,
	)
}

// Addr returns the Redis address as host:port.
func (c *RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// Addr returns the server listen address as host:port.
func (c *ServerConfig) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

func envRequired(key string) string {
	val := os.Getenv(key)
	if val == "" {
		// We don't panic here; validation catches it.
		// This allows collecting multiple errors if needed in the future.
	}
	return val
}

func envOrDefault(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func envInt(key string, fallback int) (int, error) {
	val := os.Getenv(key)
	if val == "" {
		return fallback, nil
	}
	n, err := strconv.Atoi(val)
	if err != nil {
		return 0, fmt.Errorf("invalid integer value %q for %s", val, key)
	}
	return n, nil
}

func envDuration(key string, fallback time.Duration) time.Duration {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	d, err := time.ParseDuration(val)
	if err != nil {
		return fallback
	}
	return d
}
