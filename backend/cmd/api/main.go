package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/urja-gym/backend/internal/config"
	"github.com/urja-gym/backend/internal/database"
	"github.com/urja-gym/backend/internal/handler"
	"github.com/urja-gym/backend/internal/middleware"
	"github.com/urja-gym/backend/internal/repository"
	"github.com/urja-gym/backend/internal/service"
	ujwt "github.com/urja-gym/backend/pkg/jwt"
)

func main() {
	if err := run(); err != nil {
		log.Fatal().Err(err).Msg("application failed")
	}
}

func run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Load config
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	// Configure logging
	if cfg.IsDevelopment() {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	log.Info().Str("environment", cfg.Environment).Msg("starting urja auth service")

	// Connect to Redis
	rdb, err := database.NewRedisClient(ctx, cfg.RedisURL)
	if err != nil {
		return fmt.Errorf("connecting to redis: %w", err)
	}
	defer rdb.Close()
	log.Info().Msg("connected to redis")

	// Connect to PostgreSQL
	pool, err := database.NewPostgresPool(ctx, cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("connecting to postgres: %w", err)
	}
	defer pool.Close()
	log.Info().Msg("connected to postgres")

	// Run migrations
	if err := database.RunMigrations(ctx, pool); err != nil {
		return fmt.Errorf("running migrations: %w", err)
	}
	log.Info().Msg("migrations complete")

	// Initialize dependencies
	userRepo := repository.NewUserRepository(pool)
	jwtIssuer := ujwt.NewIssuer(cfg.JWTSecret, cfg.JWTRefreshSecret)

	var smsSender service.SMSSender
	if cfg.IsDevelopment() {
		smsSender = service.NewDevSMS()
	} else {
		smsSender = service.NewAakashSMS(cfg.AakashSMSToken)
	}

	authService := service.NewAuthService(rdb, userRepo, smsSender, jwtIssuer)
	authHandler := handler.NewAuthHandler(authService)

	// Build router
	r := chi.NewRouter()

	// Global middleware
	r.Use(middleware.Recovery)
	r.Use(middleware.RequestID)
	r.Use(chimw.RealIP)
	r.Use(chimw.Compress(5))

	// Global rate limit: 100 req/min per IP
	r.Use(middleware.RateLimiter(rdb, middleware.RateLimitConfig{
		Prefix: "global",
		Max:    100,
		Window: time.Minute,
	}))

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	// Auth routes
	r.Route("/api/v1/auth", func(r chi.Router) {
		// POST /api/v1/auth/login — rate limited: 5/15min per IP
		r.With(middleware.RateLimiter(rdb, middleware.RateLimitConfig{
			Prefix: "auth_login",
			Max:    5,
			Window: 15 * time.Minute,
		})).Post("/login", authHandler.Login)

		// POST /api/v1/auth/verify-otp — rate limited: 3/5min per IP
		r.With(middleware.RateLimiter(rdb, middleware.RateLimitConfig{
			Prefix: "otp_verify",
			Max:    3,
			Window: 5 * time.Minute,
		})).Post("/verify-otp", authHandler.VerifyOTP)

		// POST /api/v1/auth/refresh
		r.Post("/refresh", authHandler.Refresh)

		// POST /api/v1/auth/logout
		r.Post("/logout", authHandler.Logout)

		// POST /api/v1/auth/logout-all — requires authentication
		r.With(middleware.Auth(jwtIssuer)).Post("/logout-all", authHandler.LogoutAll)
	})

	// Start server
	addr := fmt.Sprintf(":%d", cfg.ServerPort)
	srv := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	errCh := make(chan error, 1)
	go func() {
		log.Info().Str("addr", addr).Msg("HTTP server listening")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-quit:
		log.Info().Str("signal", sig.String()).Msg("shutting down")
	case err := <-errCh:
		return fmt.Errorf("server error: %w", err)
	}

	shutdownCtx, shutdownCancel := context.WithTimeout(ctx, 10*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("server shutdown: %w", err)
	}

	log.Info().Msg("server stopped gracefully")
	return nil
}
