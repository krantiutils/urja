package billing

import (
	"context"
	"fmt"
	"log/slog"
	"math"
	"time"

	"github.com/urja-gym/urja/pkg/khalti"
)

// Service handles billing business logic.
type Service struct {
	repo         *Repository
	khaltiClient *khalti.Client
	logger       *slog.Logger
}

// NewService creates a new billing service.
func NewService(repo *Repository, khaltiClient *khalti.Client, logger *slog.Logger) *Service {
	return &Service{
		repo:         repo,
		khaltiClient: khaltiClient,
		logger:       logger,
	}
}

// ListPlans returns all active SaaS plans.
func (s *Service) ListPlans(ctx context.Context) ([]SaaSPlan, error) {
	return s.repo.ListActivePlans(ctx)
}

// SubscribeInput holds input for subscribing an org to a plan.
type SubscribeInput struct {
	OrganizationID string
	PlanID         string
	BillingPeriod  string // "monthly" or "yearly"
	KhaltiPidx     string
}

// Subscribe verifies a Khalti payment and creates a subscription.
func (s *Service) Subscribe(ctx context.Context, userID string, in *SubscribeInput) (*SaaSSubscription, error) {
	// Verify the user is an admin of the organization (or super_admin).
	isAdmin, err := s.repo.IsOrgAdmin(ctx, userID, in.OrganizationID)
	if err != nil {
		return nil, fmt.Errorf("checking admin: %w", err)
	}
	if !isAdmin {
		isSA, saErr := s.repo.IsSuperAdmin(ctx, userID)
		if saErr != nil || !isSA {
			return nil, fmt.Errorf("forbidden: only org admins can subscribe")
		}
	}

	// Validate billing period.
	if in.BillingPeriod != "monthly" && in.BillingPeriod != "yearly" {
		return nil, fmt.Errorf("billing_period must be 'monthly' or 'yearly'")
	}

	// Get the plan.
	plan, err := s.repo.GetPlanByID(ctx, in.PlanID)
	if err != nil {
		return nil, fmt.Errorf("plan not found")
	}
	if !plan.IsActive {
		return nil, fmt.Errorf("plan is not active")
	}

	// Check for existing active subscription.
	existing, _ := s.repo.GetActiveSubscription(ctx, in.OrganizationID)
	if existing != nil {
		return nil, fmt.Errorf("organization already has an active subscription")
	}

	// Determine expected amount (in paisa = NPR * 100).
	var expectedAmount float64
	if in.BillingPeriod == "monthly" {
		expectedAmount = plan.PriceMonthly
	} else {
		expectedAmount = plan.PriceYearly
	}
	expectedPaisa := int(math.Round(expectedAmount * 100))

	// Verify Khalti payment.
	lookup, err := s.khaltiClient.LookupPayment(ctx, in.KhaltiPidx)
	if err != nil {
		return nil, fmt.Errorf("khalti payment verification failed: %w", err)
	}

	if lookup.Status != "Completed" {
		s.logger.Warn("khalti payment not completed",
			"pidx", in.KhaltiPidx,
			"status", lookup.Status,
			"org_id", in.OrganizationID,
		)
		return nil, fmt.Errorf("payment not completed, status: %s", lookup.Status)
	}

	if lookup.TotalAmount != expectedPaisa {
		s.logger.Error("khalti payment amount mismatch",
			"expected_paisa", expectedPaisa,
			"actual_paisa", lookup.TotalAmount,
			"pidx", in.KhaltiPidx,
		)
		return nil, fmt.Errorf("payment amount mismatch: expected %d paisa, got %d paisa", expectedPaisa, lookup.TotalAmount)
	}

	// Calculate subscription period.
	now := time.Now().UTC()
	var periodEnd time.Time
	if in.BillingPeriod == "monthly" {
		periodEnd = now.AddDate(0, 1, 0)
	} else {
		periodEnd = now.AddDate(1, 0, 0)
	}

	sub, err := s.repo.CreateSubscription(ctx, in.OrganizationID, in.PlanID, "active", now, periodEnd)
	if err != nil {
		return nil, fmt.Errorf("creating subscription: %w", err)
	}

	s.logger.Info("subscription created",
		"org_id", in.OrganizationID,
		"plan_id", in.PlanID,
		"period", in.BillingPeriod,
		"khalti_txn", lookup.TransactionID,
	)

	return sub, nil
}
