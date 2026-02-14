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

	"github.com/urja-gym/urja/internal/accounts"
	"github.com/urja-gym/urja/internal/attendance"
	"github.com/urja-gym/urja/internal/auth"
	"github.com/urja-gym/urja/internal/billing"
	"github.com/urja-gym/urja/internal/config"
	"github.com/urja-gym/urja/internal/dues"
	"github.com/urja-gym/urja/internal/member"
	"github.com/urja-gym/urja/internal/org"
	"github.com/urja-gym/urja/internal/packages"
	"github.com/urja-gym/urja/internal/staff"
	"github.com/urja-gym/urja/internal/subscription"
	"github.com/urja-gym/urja/pkg/database"
	"github.com/urja-gym/urja/pkg/khalti"
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

	// Khalti client
	khaltiClient := khalti.NewClient(
		cfg.Khalti.SecretKey, cfg.Khalti.BaseURL,
		cfg.Khalti.WebsiteURL, cfg.Khalti.ReturnURL, logger,
	)

	// --- Domain wiring ---

	// Auth
	authRepo := auth.NewRepository(pool, rdb)
	authService := auth.NewService(authRepo, smsClient, cfg.Auth, logger)
	authHandler := auth.NewHandler(authService, logger)

	// Organization
	orgRepo := org.NewRepository(pool)
	orgService := org.NewService(orgRepo, logger)
	orgHandler := org.NewHandler(orgService, logger)

	// Member
	memberRepo := member.NewRepository(pool)
	memberService := member.NewService(memberRepo, logger)
	memberHandler := member.NewHandler(memberService, logger)

	// Attendance
	attendanceRepo := attendance.NewRepository(pool)
	attendanceService := attendance.NewService(attendanceRepo, []byte(cfg.Auth.JWTSecret), logger)
	attendanceHandler := attendance.NewHandler(attendanceService, logger)

	// Packages
	pkgRepo := packages.NewRepository(pool)
	pkgService := packages.NewService(pkgRepo, khaltiClient, logger)
	pkgHandler := packages.NewHandler(pkgService, logger)

	// Dues
	duesRepo := dues.NewRepository(pool)
	duesService := dues.NewService(duesRepo, logger)
	duesHandler := dues.NewHandler(duesService, logger)

	// Accounts
	accountsRepo := accounts.NewRepository(pool)
	accountsService := accounts.NewService(accountsRepo, logger)
	accountsHandler := accounts.NewHandler(accountsService, logger)

	// Staff
	staffRepo := staff.NewRepository(pool)
	staffService := staff.NewService(staffRepo, logger)
	staffHandler := staff.NewHandler(staffService, logger)

	// Subscription (package lifecycle)
	subscriptionRepo := subscription.NewRepository(pool)
	subscriptionService := subscription.NewService(subscriptionRepo, logger)
	subscriptionHandler := subscription.NewHandler(subscriptionService, logger)

	// Billing
	billingRepo := billing.NewRepository(pool)
	billingService := billing.NewService(billingRepo, khaltiClient, logger)
	billingHandler := billing.NewHandler(billingService, logger)

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
		// Public routes (no auth required)
		authHandler.RegisterRoutes(r, authService)
		r.Route("/gyms", func(r chi.Router) {
			orgHandler.RegisterPublicRoutes(r)
		})

		r.Route("/packages", func(r chi.Router) {
			pkgHandler.RegisterPublicRoutes(r)
		})

		r.Route("/billing", func(r chi.Router) {
			billingHandler.RegisterPublicRoutes(r)
		})

		// Authenticated routes
		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(authService))

			// Org management (super_admin create, org admin update)
			r.Post("/orgs", orgHandler.Create)
			r.Put("/orgs/{orgId}", orgHandler.Update)

			// Billing subscribe (authenticated)
			r.Post("/billing/subscribe", billingHandler.Subscribe)

			r.Route("/members/me", func(r chi.Router) {
				memberHandler.RegisterSelfRoutes(r)
				r.Route("/attendance", func(r chi.Router) {
					attendanceHandler.RegisterSelfRoutes(r)
				})
				r.Route("/packages", func(r chi.Router) {
					pkgHandler.RegisterSelfRoutes(r)
				})
			})

			// Organization-scoped routes
			r.Route("/orgs/{orgId}", func(r chi.Router) {
				r.Use(middleware.OrgScope(orgService))

				r.Route("/packages", func(r chi.Router) {
					pkgHandler.RegisterOrgRoutes(r)
					subscriptionHandler.RegisterPackageRoutes(r)
				})

				r.Route("/members", func(r chi.Router) {
					memberHandler.RegisterOrgRoutes(r)

					r.Route("/{memberId}", func(r chi.Router) {
						subscriptionHandler.RegisterMemberRoutes(r)
					})
				})

				r.Route("/attendance", func(r chi.Router) {
					attendanceHandler.RegisterOrgRoutes(r)
				})

				r.Route("/staff", func(r chi.Router) {
					staffHandler.RegisterOrgRoutes(r)
				})

				r.Route("/dues", func(r chi.Router) {
					duesHandler.RegisterRoutes(r)
				})

				r.Route("/accounts", func(r chi.Router) {
					accountsHandler.RegisterRoutes(r)
				})

				r.Route("/workout-templates", func(r chi.Router) {
					// Workout template management (TODO)
				})

				r.Route("/leaderboard", func(r chi.Router) {
					// Leaderboard (TODO)
				})
			})
		})
	})

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
