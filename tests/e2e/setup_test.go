package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	"github.com/urja-gym/urja/internal/absentee"
	"github.com/urja-gym/urja/internal/accounts"
	"github.com/urja-gym/urja/internal/activitylog"
	"github.com/urja-gym/urja/internal/attendance"
	"github.com/urja-gym/urja/internal/auth"
	"github.com/urja-gym/urja/internal/billing"
	"github.com/urja-gym/urja/internal/config"
	"github.com/urja-gym/urja/internal/dues"
	"github.com/urja-gym/urja/internal/feedback"
	"github.com/urja-gym/urja/internal/guide"
	"github.com/urja-gym/urja/internal/health"
	"github.com/urja-gym/urja/internal/leaderboard"
	"github.com/urja-gym/urja/internal/member"
	"github.com/urja-gym/urja/internal/nfc"
	"github.com/urja-gym/urja/internal/notice"
	"github.com/urja-gym/urja/internal/org"
	"github.com/urja-gym/urja/internal/packages"
	"github.com/urja-gym/urja/internal/qrcode"
	"github.com/urja-gym/urja/internal/smsapi"
	"github.com/urja-gym/urja/internal/staff"
	"github.com/urja-gym/urja/internal/subscription"
	"github.com/urja-gym/urja/internal/workout"
	"github.com/urja-gym/urja/pkg/khalti"
	"github.com/urja-gym/urja/pkg/middleware"
	"github.com/urja-gym/urja/pkg/sms"
)

// --- Globals shared across all E2E tests ---

var (
	testServer  *httptest.Server
	testPool    *pgxpool.Pool
	testRDB     *redis.Client
	testLogger  *slog.Logger
	jwtSecret   = "test-jwt-secret-key-for-e2e-tests-at-least-32-chars"
	authService *auth.Service
)

// testConfig returns the E2E test configuration.
// It reads from environment variables with fallbacks for local dev.
func testConfig() *config.Config {
	dbHost := envOr("E2E_DB_HOST", "localhost")
	dbPort := envOr("E2E_DB_PORT", "5433")
	dbUser := envOr("E2E_DB_USER", "urja")
	dbPassword := envOr("E2E_DB_PASSWORD", "urja_secret")
	dbName := envOr("E2E_DB_NAME", "urja_test")
	redisHost := envOr("E2E_REDIS_HOST", "localhost")
	redisPort := envOr("E2E_REDIS_PORT", "6379")

	port := 0 // let OS pick
	_ = port

	portInt := 5433
	fmt.Sscanf(dbPort, "%d", &portInt)

	redisPortInt := 6379
	fmt.Sscanf(redisPort, "%d", &redisPortInt)

	return &config.Config{
		Server: config.ServerConfig{
			Host:            "127.0.0.1",
			Port:            0,
			ReadTimeout:     15 * time.Second,
			WriteTimeout:    15 * time.Second,
			ShutdownTimeout: 5 * time.Second,
			BaseURL:         "http://localhost:8080",
		},
		Database: config.DatabaseConfig{
			Host:     dbHost,
			Port:     portInt,
			User:     dbUser,
			Password: dbPassword,
			Name:     dbName,
			SSLMode:  "disable",
			MaxConns: 10,
			MinConns: 2,
		},
		Redis: config.RedisConfig{
			Host: redisHost,
			Port: redisPortInt,
		},
		Auth: config.AuthConfig{
			JWTSecret:          jwtSecret,
			AccessTokenExpiry:  15 * time.Minute,
			RefreshTokenExpiry: 7 * 24 * time.Hour,
			OTPExpiry:          5 * time.Minute,
			OTPMaxAttempts:     5,
			OTPCooldown:        1 * time.Second,
			OTPHourlyLimit:     100,
			BcryptCost:         4, // low cost for fast tests
		},
		SMS: config.SMSConfig{
			AakashToken:  "test-token",
			AakashAPIURL: "", // will be set to mock server URL
		},
		Khalti: config.KhaltiConfig{
			SecretKey:  "test-khalti-secret",
			BaseURL:    "", // will be set to mock server URL
			WebsiteURL: "http://localhost",
			ReturnURL:  "http://localhost/return",
		},
	}
}

func envOr(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

// TestMain is the entry point for all E2E tests.
func TestMain(m *testing.M) {
	testLogger = slog.New(slog.NewTextHandler(io.Discard, nil))

	cfg := testConfig()
	ctx := context.Background()

	// --- Mock SMS server ---
	smsMockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "sent"})
	}))
	defer smsMockServer.Close()
	cfg.SMS.AakashAPIURL = smsMockServer.URL

	// --- Mock Khalti server ---
	khaltiMockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":         "Completed",
			"total_amount":   100000,
			"transaction_id": "test-txn-123",
			"pidx":           "test-pidx-123",
		})
	}))
	defer khaltiMockServer.Close()
	cfg.Khalti.BaseURL = khaltiMockServer.URL

	// --- Database ---
	pool, err := pgxpool.New(ctx, cfg.Database.DSN())
	if err != nil {
		fmt.Fprintf(os.Stderr, "FATAL: cannot connect to test database: %v\n", err)
		fmt.Fprintf(os.Stderr, "DSN: %s\n", cfg.Database.DSN())
		fmt.Fprintf(os.Stderr, "Make sure PostgreSQL is running on port %d with database '%s'\n", cfg.Database.Port, cfg.Database.Name)
		os.Exit(1)
	}
	testPool = pool
	defer pool.Close()

	// --- Run migrations ---
	if err := runMigrations(ctx, pool); err != nil {
		fmt.Fprintf(os.Stderr, "FATAL: migration failed: %v\n", err)
		os.Exit(1)
	}

	// --- Redis ---
	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.Redis.Addr(),
	})
	if err := rdb.Ping(ctx).Err(); err != nil {
		fmt.Fprintf(os.Stderr, "FATAL: cannot connect to Redis: %v\n", err)
		os.Exit(1)
	}
	testRDB = rdb
	defer rdb.Close()

	// --- Wire all domains (mirrors cmd/api/main.go) ---
	smsClient := sms.NewClient(cfg.SMS.AakashToken, cfg.SMS.AakashAPIURL, testLogger)
	khaltiClient := khalti.NewClient(cfg.Khalti.SecretKey, cfg.Khalti.BaseURL, cfg.Khalti.WebsiteURL, cfg.Khalti.ReturnURL, testLogger)

	authRepo := auth.NewRepository(pool, rdb)
	authSvc := auth.NewService(authRepo, smsClient, cfg.Auth, testLogger)
	authHandler := auth.NewHandler(authSvc, testLogger)
	authService = authSvc

	orgRepo := org.NewRepository(pool)
	orgSvc := org.NewService(orgRepo, testLogger)
	orgHandler := org.NewHandler(orgSvc, testLogger)

	memberRepo := member.NewRepository(pool)
	memberSvc := member.NewService(memberRepo, testLogger)
	memberHandler := member.NewHandler(memberSvc, testLogger)

	attendanceRepo := attendance.NewRepository(pool)
	attendanceSvc := attendance.NewService(attendanceRepo, []byte(cfg.Auth.JWTSecret), testLogger)
	attendanceHandler := attendance.NewHandler(attendanceSvc, testLogger)

	pkgRepo := packages.NewRepository(pool)
	pkgSvc := packages.NewService(pkgRepo, khaltiClient, testLogger)
	pkgHandler := packages.NewHandler(pkgSvc, testLogger)

	nfcRepo := nfc.NewRepository(pool)
	nfcSvc := nfc.NewService(nfcRepo, testLogger)
	nfcHandler := nfc.NewHandler(nfcSvc, testLogger)

	activityLogRepo := activitylog.NewRepository(pool)
	activityLogSvc := activitylog.NewService(activityLogRepo, testLogger)
	activityLogHandler := activitylog.NewHandler(activityLogSvc, testLogger)

	qrHandler := qrcode.NewHandler(cfg.Server.BaseURL, testLogger)

	noticeRepo := notice.NewRepository(pool)
	noticeSvc := notice.NewService(noticeRepo, testLogger)
	noticeHandler := notice.NewHandler(noticeSvc, testLogger)

	feedbackRepo := feedback.NewRepository(pool)
	feedbackSvc := feedback.NewService(feedbackRepo, testLogger)
	feedbackHandler := feedback.NewHandler(feedbackSvc, testLogger)

	smsapiRepo := smsapi.NewRepository(pool)
	smsapiSvc := smsapi.NewService(smsapiRepo, smsClient, testLogger)
	smsapiHandler := smsapi.NewHandler(smsapiSvc, testLogger)

	absenteeRepo := absentee.NewRepository(pool)
	absenteeSvc := absentee.NewService(absenteeRepo, smsClient, testLogger)
	absenteeHandler := absentee.NewHandler(absenteeSvc, testLogger)

	duesRepo := dues.NewRepository(pool)
	duesSvc := dues.NewService(duesRepo, testLogger)
	duesHandler := dues.NewHandler(duesSvc, testLogger)

	accountsRepo := accounts.NewRepository(pool)
	accountsSvc := accounts.NewService(accountsRepo, testLogger)
	accountsHandler := accounts.NewHandler(accountsSvc, testLogger)

	staffRepo := staff.NewRepository(pool)
	staffSvc := staff.NewService(staffRepo, testLogger)
	staffHandler := staff.NewHandler(staffSvc, testLogger)

	subscriptionRepo := subscription.NewRepository(pool)
	subscriptionSvc := subscription.NewService(subscriptionRepo, testLogger)
	subscriptionHandler := subscription.NewHandler(subscriptionSvc, testLogger)

	healthRepo := health.NewRepository(pool)
	healthSvc := health.NewService(healthRepo, filepath.Join(os.TempDir(), "urja-test-uploads"), testLogger)
	healthHandler := health.NewHandler(healthSvc, testLogger)

	billingRepo := billing.NewRepository(pool)
	billingSvc := billing.NewService(billingRepo, khaltiClient, testLogger)
	billingHandler := billing.NewHandler(billingSvc, testLogger)

	guideRepo := guide.NewRepository(pool)
	guideSvc := guide.NewService(guideRepo, testLogger)
	guideHandler := guide.NewHandler(guideSvc, testLogger)

	leaderboardRepo := leaderboard.NewRepository(pool)
	leaderboardSvc := leaderboard.NewService(leaderboardRepo, testLogger)
	leaderboardHandler := leaderboard.NewHandler(leaderboardSvc, testLogger)

	workoutRepo := workout.NewRepository(pool)
	workoutSvc := workout.NewService(workoutRepo, testLogger)
	workoutHandler := workout.NewHandler(workoutSvc, testLogger)

	// --- Build chi router (mirrors main.go) ---
	r := chi.NewRouter()

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	r.Route("/api/v1", func(r chi.Router) {
		authHandler.RegisterRoutes(r, authSvc)
		r.Route("/gyms", func(r chi.Router) {
			orgHandler.RegisterPublicRoutes(r)
		})
		r.Route("/packages", func(r chi.Router) {
			pkgHandler.RegisterPublicRoutes(r)
		})
		r.Route("/billing", func(r chi.Router) {
			billingHandler.RegisterPublicRoutes(r)
		})

		r.Route("/training-guides", func(r chi.Router) {
			guideHandler.RegisterPublicRoutes(r)
		})

		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(authSvc))

			r.Post("/orgs", orgHandler.Create)
			r.Post("/billing/subscribe", billingHandler.Subscribe)

			r.Route("/members/me", func(r chi.Router) {
				memberHandler.RegisterSelfRoutes(r)
				workoutHandler.RegisterSelfRoutes(r)
				r.Route("/attendance", func(r chi.Router) {
					attendanceHandler.RegisterSelfRoutes(r)
				})
				r.Route("/packages", func(r chi.Router) {
					pkgHandler.RegisterSelfRoutes(r)
				})
				r.Route("/health", func(r chi.Router) {
					healthHandler.RegisterSelfRoutes(r)
				})
			})

			r.Route("/orgs/{orgId}", func(r chi.Router) {
				r.Use(middleware.OrgScope(orgSvc))

				r.Put("/", orgHandler.Update)

				r.Route("/packages", func(r chi.Router) {
					pkgHandler.RegisterOrgRoutes(r)
					subscriptionHandler.RegisterPackageRoutes(r)
				})

				r.Route("/members", func(r chi.Router) {
					memberHandler.RegisterOrgRoutes(r)
					r.Route("/{memberId}", func(r chi.Router) {
						memberHandler.RegisterMemberRoutes(r)
						subscriptionHandler.RegisterMemberRoutes(r)
						workoutHandler.RegisterMemberRoutes(r)
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

				r.Route("/notices", func(r chi.Router) {
					noticeHandler.RegisterOrgRoutes(r)
				})

				r.Route("/feedbacks", func(r chi.Router) {
					feedbackHandler.RegisterOrgRoutes(r)
				})

				r.Route("/sms", func(r chi.Router) {
					smsapiHandler.RegisterOrgRoutes(r)
				})

				r.Route("/absentees", func(r chi.Router) {
					absenteeHandler.RegisterOrgRoutes(r)
				})

				r.Route("/training-guides", func(r chi.Router) {
					guideHandler.RegisterOrgRoutes(r)
				})

				r.Route("/workout-templates", func(r chi.Router) {
					workoutHandler.RegisterOrgRoutes(r)
				})

				nfcHandler.RegisterOrgRoutes(r)
				activityLogHandler.RegisterOrgRoutes(r)
				qrHandler.RegisterOrgRoutes(r)

				r.Route("/leaderboard", func(r chi.Router) {
					leaderboardHandler.RegisterOrgRoutes(r)
				})
			})
		})
	})

	testServer = httptest.NewServer(r)
	defer testServer.Close()

	os.Exit(m.Run())
}

// runMigrations executes all up migration SQL files against the test database.
func runMigrations(ctx context.Context, pool *pgxpool.Pool) error {
	migrationsDir := filepath.Join(projectRoot(), "db", "migrations")
	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("reading migrations dir: %w", err)
	}

	// Collect and sort only .up.sql files
	var upFiles []string
	for _, e := range entries {
		if strings.HasSuffix(e.Name(), ".up.sql") {
			upFiles = append(upFiles, e.Name())
		}
	}
	sort.Strings(upFiles)

	// Drop all tables first for a clean slate
	_, err = pool.Exec(ctx, `
		DO $$ DECLARE
			r RECORD;
		BEGIN
			FOR r IN (SELECT tablename FROM pg_tables WHERE schemaname = 'public') LOOP
				EXECUTE 'DROP TABLE IF EXISTS public.' || quote_ident(r.tablename) || ' CASCADE';
			END LOOP;
		END $$;
	`)
	if err != nil {
		return fmt.Errorf("dropping existing tables: %w", err)
	}

	// Also drop extensions and functions so migrations can recreate them
	_, _ = pool.Exec(ctx, `DROP FUNCTION IF EXISTS trigger_set_updated_at() CASCADE;`)

	for _, f := range upFiles {
		sqlBytes, err := os.ReadFile(filepath.Join(migrationsDir, f))
		if err != nil {
			return fmt.Errorf("reading %s: %w", f, err)
		}
		sql := string(sqlBytes)
		if strings.TrimSpace(sql) == "" {
			continue
		}
		if _, err := pool.Exec(ctx, sql); err != nil {
			return fmt.Errorf("executing %s: %w", f, err)
		}
	}

	return nil
}

// projectRoot returns the absolute path to the project root.
func projectRoot() string {
	// tests/e2e/ -> ../../
	dir, err := os.Getwd()
	if err != nil {
		panic("cannot get working directory: " + err.Error())
	}
	// Walk up until we find go.mod
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			panic("cannot find project root (go.mod not found)")
		}
		dir = parent
	}
}

// ---------- Test helpers ----------

// generateTestToken creates a valid JWT access token for testing.
func generateTestToken(userID, role string) string {
	now := time.Now()
	claims := jwt.MapClaims{
		"sub":   userID,
		"phone": "9800000000",
		"role":  role,
		"iss":   "urja",
		"iat":   jwt.NewNumericDate(now),
		"exp":   jwt.NewNumericDate(now.Add(15 * time.Minute)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		panic("failed to sign test token: " + err.Error())
	}
	return signed
}

// createTestUser inserts a user into the database and returns the user ID.
func createTestUser(t *testing.T, phone, name string) string {
	t.Helper()
	var id string
	err := testPool.QueryRow(context.Background(),
		`INSERT INTO users (phone, name) VALUES ($1, $2)
		 ON CONFLICT (phone) DO UPDATE SET name = EXCLUDED.name
		 RETURNING id`,
		phone, name,
	).Scan(&id)
	if err != nil {
		t.Fatalf("createTestUser(%s): %v", phone, err)
	}
	return id
}

// createTestSuperAdmin inserts a super admin user and returns the user ID.
func createTestSuperAdmin(t *testing.T, phone, name string) string {
	t.Helper()
	var id string
	err := testPool.QueryRow(context.Background(),
		`INSERT INTO users (phone, name, is_super_admin) VALUES ($1, $2, true)
		 ON CONFLICT (phone) DO UPDATE SET name = EXCLUDED.name, is_super_admin = true
		 RETURNING id`,
		phone, name,
	).Scan(&id)
	if err != nil {
		t.Fatalf("createTestSuperAdmin(%s): %v", phone, err)
	}
	return id
}

// createTestOrg inserts an organization and returns the org ID.
// It also makes the given user an admin of the org.
func createTestOrg(t *testing.T, adminUserID, orgName string) string {
	t.Helper()
	var id string
	slug := strings.ToLower(strings.ReplaceAll(orgName, " ", "-"))
	err := testPool.QueryRow(context.Background(),
		`INSERT INTO organizations (name, slug) VALUES ($1, $2) RETURNING id`,
		orgName, slug,
	).Scan(&id)
	if err != nil {
		t.Fatalf("createTestOrg(%s): %v", orgName, err)
	}

	_, err = testPool.Exec(context.Background(),
		`INSERT INTO organization_members (user_id, organization_id, role, status)
		 VALUES ($1, $2, 'admin', 'active')
		 ON CONFLICT (user_id, organization_id) DO UPDATE SET role = 'admin', status = 'active'`,
		adminUserID, id,
	)
	if err != nil {
		t.Fatalf("createTestOrg: adding admin membership: %v", err)
	}

	return id
}

// createTestOrgMember adds a user as a member of an org with the given role.
func createTestOrgMember(t *testing.T, userID, orgID, role string) {
	t.Helper()
	_, err := testPool.Exec(context.Background(),
		`INSERT INTO organization_members (user_id, organization_id, role, status)
		 VALUES ($1, $2, $3, 'active')
		 ON CONFLICT (user_id, organization_id) DO UPDATE SET role = $3, status = 'active'`,
		userID, orgID, role,
	)
	if err != nil {
		t.Fatalf("createTestOrgMember(%s, %s, %s): %v", userID, orgID, role, err)
	}
}

// createTestPackage inserts a package for the org and returns the package ID.
func createTestPackage(t *testing.T, orgID, name string, durationDays int, price float64) string {
	t.Helper()
	var id string
	err := testPool.QueryRow(context.Background(),
		`INSERT INTO packages (organization_id, name, duration_days, price) VALUES ($1, $2, $3, $4) RETURNING id`,
		orgID, name, durationDays, price,
	).Scan(&id)
	if err != nil {
		t.Fatalf("createTestPackage(%s): %v", name, err)
	}
	return id
}

// assignTestPackage assigns a package to a member and returns the member_package ID.
func assignTestPackage(t *testing.T, userID, pkgID, orgID string) string {
	t.Helper()
	var id string
	err := testPool.QueryRow(context.Background(),
		`INSERT INTO member_packages (user_id, package_id, organization_id, start_date, end_date, amount_paid, status, payment_method)
		 VALUES ($1, $2, $3, CURRENT_DATE, CURRENT_DATE + INTERVAL '30 days', 1000, 'active', 'cash')
		 RETURNING id`,
		userID, pkgID, orgID,
	).Scan(&id)
	if err != nil {
		t.Fatalf("assignTestPackage: %v", err)
	}
	return id
}

// createTestDue inserts a due record and returns its ID.
func createTestDue(t *testing.T, orgID, userID string, amount float64) string {
	t.Helper()
	var id string
	err := testPool.QueryRow(context.Background(),
		`INSERT INTO dues (organization_id, user_id, amount, due_date, description, status)
		 VALUES ($1, $2, $3, CURRENT_DATE, 'Test due', 'pending')
		 RETURNING id`,
		orgID, userID, amount,
	).Scan(&id)
	if err != nil {
		t.Fatalf("createTestDue: %v", err)
	}
	return id
}

// createTestNotice inserts a notice and returns its ID.
func createTestNotice(t *testing.T, orgID, userID, title, content string) string {
	t.Helper()
	var id string
	err := testPool.QueryRow(context.Background(),
		`INSERT INTO notices (organization_id, title, content, created_by)
		 VALUES ($1, $2, $3, $4) RETURNING id`,
		orgID, title, content, userID,
	).Scan(&id)
	if err != nil {
		t.Fatalf("createTestNotice: %v", err)
	}
	return id
}

// createTestGuide inserts a training guide and returns its ID.
func createTestGuide(t *testing.T, authorID, title, content, category string, published bool) string {
	t.Helper()
	var id string
	err := testPool.QueryRow(context.Background(),
		`INSERT INTO training_guides (title, content, category, author_id, is_published)
		 VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		title, content, category, authorID, published,
	).Scan(&id)
	if err != nil {
		t.Fatalf("createTestGuide: %v", err)
	}
	return id
}

// createTestNFCCard inserts an NFC card for an org and returns its ID.
func createTestNFCCard(t *testing.T, orgID, cardHex string) string {
	t.Helper()
	var id string
	err := testPool.QueryRow(context.Background(),
		`INSERT INTO nfc_cards (organization_id, card_hex) VALUES ($1, $2) RETURNING id`,
		orgID, cardHex,
	).Scan(&id)
	if err != nil {
		t.Fatalf("createTestNFCCard: %v", err)
	}
	return id
}

// initSMSCredits sets up SMS credits for an org.
func initSMSCredits(t *testing.T, orgID string, balance int) {
	t.Helper()
	_, err := testPool.Exec(context.Background(),
		`INSERT INTO sms_credits (organization_id, balance, total_purchased)
		 VALUES ($1, $2, $2)
		 ON CONFLICT (organization_id) DO UPDATE SET balance = $2, total_purchased = $2`,
		orgID, balance,
	)
	if err != nil {
		t.Fatalf("initSMSCredits: %v", err)
	}
}

// cleanupTables truncates all data tables, preserving schema.
func cleanupTables(t *testing.T) {
	t.Helper()
	tables := []string{
		"progress_photos", "health_metrics",
		"activity_logs", "nfc_devices", "nfc_cards",
		"sms_campaigns", "sms_purchases", "sms_credits",
		"training_guides",
		"feedbacks", "notices",
		"payments", "dues", "transactions",
		"member_achievements", "achievements", "streaks",
		"attendance",
		"member_packages", "packages",
		"organization_members", "organizations",
		"users",
	}
	// Also handle tables that may not exist yet
	for _, table := range tables {
		_, _ = testPool.Exec(context.Background(),
			fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table))
	}
	// Clear Redis
	testRDB.FlushDB(context.Background())
}

// ---------- HTTP helpers ----------

// doRequest sends an HTTP request to the test server and returns the response.
func doRequest(t *testing.T, method, path string, body interface{}, token string) *http.Response {
	t.Helper()

	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("marshaling request body: %v", err)
		}
		bodyReader = bytes.NewReader(data)
	}

	req, err := http.NewRequest(method, testServer.URL+path, bodyReader)
	if err != nil {
		t.Fatalf("creating request: %v", err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("sending request: %v", err)
	}
	return resp
}

// parseJSON decodes the response body into the given value.
func parseJSON(t *testing.T, resp *http.Response, v interface{}) {
	t.Helper()
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("reading response body: %v", err)
	}
	if err := json.Unmarshal(data, v); err != nil {
		t.Fatalf("decoding JSON: %v\nbody: %s", err, string(data))
	}
}

// assertStatus checks that the response has the expected HTTP status code.
func assertStatus(t *testing.T, resp *http.Response, expected int) {
	t.Helper()
	if resp.StatusCode != expected {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		t.Fatalf("expected status %d, got %d\nbody: %s", expected, resp.StatusCode, string(body))
	}
}

// readBody reads and returns the response body as a string, then closes it.
func readBody(t *testing.T, resp *http.Response) string {
	t.Helper()
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("reading response body: %v", err)
	}
	return string(data)
}
