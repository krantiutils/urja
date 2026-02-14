package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"

	"github.com/urja-gym/urja/internal/attendance"
	"github.com/urja-gym/urja/internal/config"
	"github.com/urja-gym/urja/pkg/database"
	"github.com/urja-gym/urja/pkg/middleware"
	uredis "github.com/urja-gym/urja/pkg/redis"
	"github.com/urja-gym/urja/pkg/sms"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	// Load .env in development; ignore error if file doesn't exist.
	if err := godotenv.Load(); err != nil {
		logger.Info("no .env file found, using environment variables")
	}

	cfg, err := config.Load()
	if err != nil {
		logger.Error("failed to load configuration", "error", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Database
	pool, err := database.New(ctx, cfg.Database, logger)
	if err != nil {
		logger.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	// Redis
	rdb, err := uredis.New(ctx, cfg.Redis, logger)
	if err != nil {
		logger.Error("failed to connect to redis", "error", err)
		os.Exit(1)
	}
	defer rdb.Close()

	// SMS client
	smsClient := sms.NewClient(cfg.SMS.AakashToken, cfg.SMS.AakashAPIURL, logger)

	// Attendance module
	attendanceRepo := attendance.NewRepository(pool)
	attendanceService := attendance.NewService(attendanceRepo, []byte(cfg.Auth.JWTSecret), logger)
	attendanceHandler := attendance.NewHandler(attendanceService, logger)

	// Router
	r := chi.NewRouter()

	// Global middleware
	globalLimiter := middleware.NewRateLimiter(100.0/60.0, 20) // ~100 req/min, burst 20
	r.Use(middleware.Recovery(logger))
	r.Use(middleware.Logging(logger))
	r.Use(globalLimiter.Limit())

	// Health check
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	// API v1 routes
	r.Route("/api/v1", func(r chi.Router) {
		// Public routes
		r.Route("/auth", func(r chi.Router) {
			// Auth routes will be registered by the auth domain
		})

		r.Route("/gyms", func(r chi.Router) {
			// Public gym listing routes
		})

		r.Route("/packages", func(r chi.Router) {
			// Public package listing routes
		})

		// Authenticated routes
		r.Group(func(r chi.Router) {
			// Auth middleware will be added here once TokenValidator is implemented
			// r.Use(middleware.Auth(authService))

			r.Route("/members/me", func(r chi.Router) {
				attendanceHandler.RegisterSelfRoutes(r)
			})

			// Organization-scoped routes
			r.Route("/orgs/{orgId}", func(r chi.Router) {
				// Org scope middleware will be added here
				// r.Use(middleware.OrgScope(orgService))

				r.Route("/members", func(r chi.Router) {
					// Staff/admin member management
				})

				r.Route("/attendance", func(r chi.Router) {
					attendanceHandler.RegisterOrgRoutes(r)
				})

				r.Route("/workout-templates", func(r chi.Router) {
					// Workout template management
				})

				r.Route("/leaderboard", func(r chi.Router) {
					// Leaderboard
				})
			})
		})
	})

	// Suppress unused variable warnings — these will be wired in domain packages.
	_ = rdb
	_ = smsClient

	srv := &http.Server{
		Addr:         cfg.Server.Addr(),
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// Graceful shutdown
	errCh := make(chan error, 1)
	go func() {
		logger.Info("server starting", "addr", cfg.Server.Addr())
		errCh <- srv.ListenAndServe()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-quit:
		logger.Info("shutdown signal received", "signal", sig)
	case err := <-errCh:
		if err != nil && err != http.ErrServerClosed {
			logger.Error("server error", "error", err)
		}
	}

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("server shutdown error", "error", err)
		os.Exit(1)
	}

	logger.Info("server stopped gracefully")
}
