package packages

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"math"
	"strconv"
	"time"

	"github.com/urja-gym/urja/pkg/khalti"
)

// ErrOrgRequired is returned when a public listing is requested without naming
// a gym.
var ErrOrgRequired = errors.New("gym is required")

// Service handles package business logic.
type Service struct {
	repo   *Repository
	khalti *khalti.Client
	logger *slog.Logger
}

// NewService creates a new packages service.
func NewService(repo *Repository, khaltiClient *khalti.Client, logger *slog.Logger) *Service {
	return &Service{repo: repo, khalti: khaltiClient, logger: logger}
}

// ListPublic returns one gym's active packages.
//
// The gym is mandatory. Pricing is public but it is public *per gym*: an
// unscoped listing returned every organization's packages to anyone, which both
// leaked competitors' price lists and put the wrong prices on a tenant site.
func (s *Service) ListPublic(ctx context.Context, orgSlug string, limit, offset int) ([]Package, error) {
	if orgSlug == "" {
		return nil, ErrOrgRequired
	}
	limit, offset = normalizePagination(limit, offset)
	return s.repo.ListPublic(ctx, orgSlug, limit, offset)
}

// ListByOrg returns all packages for an organization (admin view, includes inactive).
func (s *Service) ListByOrg(ctx context.Context, orgID string, limit, offset int) ([]Package, error) {
	limit, offset = normalizePagination(limit, offset)
	return s.repo.ListByOrg(ctx, orgID, limit, offset)
}

// CreatePackageInput holds validated input for creating a package.
type CreatePackageInput struct {
	Name          string          `json:"name"`
	NameNe        *string         `json:"name_ne"`
	Description   *string         `json:"description"`
	DescriptionNe *string         `json:"description_ne"`
	DurationDays  int             `json:"duration_days"`
	Price         string          `json:"price"`
	MaxMembers    *int            `json:"max_members"`
	Features      json.RawMessage `json:"features"`
}

// Create creates a new package for an organization.
func (s *Service) Create(ctx context.Context, orgID string, input CreatePackageInput) (*Package, error) {
	if input.Name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if input.DurationDays <= 0 {
		return nil, fmt.Errorf("duration_days must be positive")
	}
	priceVal, err := strconv.ParseFloat(input.Price, 64)
	if err != nil || priceVal < 0 {
		return nil, fmt.Errorf("price must be a non-negative number")
	}
	if input.Features == nil {
		input.Features = json.RawMessage("[]")
	}

	pkg, err := s.repo.Create(ctx, orgID, input.Name, input.NameNe, input.Description, input.DescriptionNe, input.DurationDays, input.Price, input.MaxMembers, input.Features)
	if err != nil {
		return nil, fmt.Errorf("creating package: %w", err)
	}

	s.logger.Info("package created", "id", pkg.ID, "org_id", orgID, "name", input.Name)
	return pkg, nil
}

// UpdatePackageInput holds input for updating a package.
type UpdatePackageInput struct {
	Name          *string         `json:"name"`
	NameNe        *string         `json:"name_ne"`
	Description   *string         `json:"description"`
	DescriptionNe *string         `json:"description_ne"`
	DurationDays  *int            `json:"duration_days"`
	Price         *string         `json:"price"`
	MaxMembers    *int            `json:"max_members"`
	Features      json.RawMessage `json:"features"`
	IsActive      *bool           `json:"is_active"`
}

// Update modifies an existing package.
func (s *Service) Update(ctx context.Context, id, orgID string, input UpdatePackageInput) (*Package, error) {
	if input.DurationDays != nil && *input.DurationDays <= 0 {
		return nil, fmt.Errorf("duration_days must be positive")
	}
	if input.Price != nil {
		priceVal, err := strconv.ParseFloat(*input.Price, 64)
		if err != nil || priceVal < 0 {
			return nil, fmt.Errorf("price must be a non-negative number")
		}
	}

	pkg, err := s.repo.Update(ctx, id, orgID, input.Name, input.NameNe, input.Description, input.DescriptionNe, input.DurationDays, input.Price, input.MaxMembers, input.Features, input.IsActive)
	if err != nil {
		return nil, fmt.Errorf("updating package: %w", err)
	}

	s.logger.Info("package updated", "id", id, "org_id", orgID)
	return pkg, nil
}

// Delete soft-deletes a package (sets is_active = false).
func (s *Service) Delete(ctx context.Context, id, orgID string) error {
	if err := s.repo.SoftDelete(ctx, id, orgID); err != nil {
		return fmt.Errorf("deleting package: %w", err)
	}
	s.logger.Info("package soft-deleted", "id", id, "org_id", orgID)
	return nil
}

// ListMemberPackages returns a user's subscriptions.
func (s *Service) ListMemberPackages(ctx context.Context, userID string, limit, offset int) ([]MemberPackage, error) {
	limit, offset = normalizePagination(limit, offset)
	return s.repo.ListMemberPackages(ctx, userID, limit, offset)
}

// SubscribeResult is returned when initiating a Khalti subscription.
type SubscribeResult struct {
	MemberPackageID string `json:"member_package_id"`
	PaymentURL      string `json:"payment_url"`
	Pidx            string `json:"pidx"`
}

// Subscribe initiates a Khalti payment for a package subscription.
func (s *Service) Subscribe(ctx context.Context, userID, packageID string) (*SubscribeResult, error) {
	// Fetch the package.
	pkg, err := s.repo.GetByID(ctx, packageID)
	if err != nil {
		return nil, fmt.Errorf("package not found: %w", err)
	}
	if !pkg.IsActive {
		return nil, fmt.Errorf("package is not active")
	}

	// Check for existing active/pending subscription.
	hasSub, err := s.repo.HasActiveSub(ctx, userID, packageID)
	if err != nil {
		return nil, fmt.Errorf("checking existing subscription: %w", err)
	}
	if hasSub {
		return nil, fmt.Errorf("already subscribed to this package")
	}

	// Convert price (decimal string like "5000.00") to paisa (int).
	priceFloat, err := strconv.ParseFloat(pkg.Price, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid package price: %w", err)
	}
	amountPaisa := int64(math.Round(priceFloat * 100))

	if amountPaisa < 1000 {
		return nil, fmt.Errorf("package price too low for Khalti (minimum Rs. 10)")
	}

	// Calculate subscription dates.
	now := time.Now()
	startDate := now
	endDate := now.AddDate(0, 0, pkg.DurationDays)

	// Initiate Khalti payment.
	khaltiResp, err := s.khalti.Initiate(ctx, fmt.Sprintf("pkg-%s-user-%s", packageID[:8], userID[:8]), pkg.Name, amountPaisa, nil)
	if err != nil {
		return nil, fmt.Errorf("initiating payment: %w", err)
	}

	// Create pending member_package with Khalti pidx as payment_reference.
	paymentRef := khaltiResp.Pidx
	mp, err := s.repo.CreateMemberPackage(ctx, userID, packageID, pkg.OrganizationID, startDate, endDate, "khalti", &paymentRef, pkg.Price, "pending")
	if err != nil {
		return nil, fmt.Errorf("creating subscription record: %w", err)
	}

	s.logger.Info("subscription initiated",
		"member_package_id", mp.ID,
		"user_id", userID,
		"package_id", packageID,
		"pidx", khaltiResp.Pidx,
	)

	return &SubscribeResult{
		MemberPackageID: mp.ID,
		PaymentURL:      khaltiResp.PaymentURL,
		Pidx:            khaltiResp.Pidx,
	}, nil
}

// VerifyPaymentResult is returned after successful payment verification.
type VerifyPaymentResult struct {
	MemberPackage *MemberPackage `json:"member_package"`
	TransactionID string         `json:"transaction_id"`
}

// VerifyPayment verifies a Khalti payment and activates the subscription.
func (s *Service) VerifyPayment(ctx context.Context, userID, pidx string) (*VerifyPaymentResult, error) {
	// Find the pending member_package by pidx.
	mp, err := s.repo.GetMemberPackageByPaymentRef(ctx, pidx)
	if err != nil {
		return nil, fmt.Errorf("subscription not found for payment reference: %w", err)
	}

	// Verify the subscription belongs to this user.
	if mp.UserID != userID {
		return nil, fmt.Errorf("payment does not belong to this user")
	}

	// Look up payment status with Khalti.
	lookup, err := s.khalti.Lookup(ctx, pidx)
	if err != nil {
		return nil, fmt.Errorf("khalti verification failed: %w", err)
	}

	if !lookup.IsCompleted() {
		// If user canceled or payment expired, cancel the member_package.
		switch lookup.Status {
		case "User canceled", "Expired", "Refunded":
			if cancelErr := s.repo.CancelMemberPackage(ctx, mp.ID); cancelErr != nil {
				s.logger.Error("failed to cancel member package after failed payment",
					"id", mp.ID, "error", cancelErr)
			}
			return nil, fmt.Errorf("payment %s: %s", lookup.Status, pidx)
		default:
			// Pending, Initiated — payment still in progress.
			return nil, fmt.Errorf("payment not yet completed (status: %s)", lookup.Status)
		}
	}

	// Verify amount matches.
	pkg, err := s.repo.GetByID(ctx, mp.PackageID)
	if err != nil {
		return nil, fmt.Errorf("fetching package for verification: %w", err)
	}
	expectedPaisa := int64(0)
	priceFloat, _ := strconv.ParseFloat(pkg.Price, 64)
	expectedPaisa = int64(math.Round(priceFloat * 100))

	if lookup.TotalAmount != expectedPaisa {
		s.logger.Error("khalti amount mismatch",
			"expected_paisa", expectedPaisa,
			"received_paisa", lookup.TotalAmount,
			"pidx", pidx,
		)
		return nil, fmt.Errorf("payment amount mismatch: expected %d paisa, got %d paisa", expectedPaisa, lookup.TotalAmount)
	}

	// Activate the subscription, storing the Khalti transaction_id.
	txnRef := lookup.TransactionID
	if txnRef == "" {
		txnRef = pidx
	}
	if err := s.repo.ActivateMemberPackage(ctx, mp.ID, txnRef); err != nil {
		return nil, fmt.Errorf("activating subscription: %w", err)
	}

	s.logger.Info("subscription activated",
		"member_package_id", mp.ID,
		"user_id", userID,
		"transaction_id", lookup.TransactionID,
	)

	// Re-fetch to get updated status.
	mpID := mp.ID
	mp, err = s.repo.GetMemberPackageByID(ctx, mpID)
	if err != nil {
		// Non-fatal: the activation succeeded, just can't return updated record.
		s.logger.Error("failed to re-fetch activated member package", "id", mpID, "error", err)
	}

	return &VerifyPaymentResult{
		MemberPackage: mp,
		TransactionID: lookup.TransactionID,
	}, nil
}

func normalizePagination(limit, offset int) (int, int) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	return limit, offset
}
