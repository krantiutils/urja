package subscription

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"
)

var (
	ErrInvalidPackage      = errors.New("package not found or inactive")
	ErrInvalidMember       = errors.New("member not found in this organization")
	ErrInvalidPayment      = errors.New("invalid payment details")
	ErrSubscriptionNotFound = errors.New("subscription not found")
	ErrInvalidRequest      = errors.New("invalid request")
)

var validPaymentMethods = map[string]bool{
	"khalti":        true,
	"cash":          true,
	"manual":        true,
	"bank_transfer": true,
}

// Service handles subscription business logic.
type Service struct {
	repo   *Repository
	logger *slog.Logger
}

// NewService creates a new subscription service.
func NewService(repo *Repository, logger *slog.Logger) *Service {
	return &Service{repo: repo, logger: logger}
}

// GetPackageSummary returns all packages for an org with member counts.
func (s *Service) GetPackageSummary(ctx context.Context, orgID string) ([]PackageSummary, error) {
	summaries, err := s.repo.PackageSummaryList(ctx, orgID)
	if err != nil {
		return nil, fmt.Errorf("fetching package summary: %w", err)
	}
	if summaries == nil {
		summaries = []PackageSummary{}
	}
	return summaries, nil
}

// ListExpiring returns member packages expiring within the given number of days.
func (s *Service) ListExpiring(ctx context.Context, orgID string, withinDays, limit, offset int) ([]ExpiringEntry, int, error) {
	if withinDays <= 0 {
		withinDays = 30
	}
	limit, offset = normalizePagination(limit, offset)
	return s.repo.ListExpiring(ctx, orgID, withinDays, limit, offset)
}

// ListExpired returns member packages that have already expired.
func (s *Service) ListExpired(ctx context.Context, orgID string, limit, offset int) ([]ExpiredEntry, int, error) {
	limit, offset = normalizePagination(limit, offset)
	return s.repo.ListExpired(ctx, orgID, limit, offset)
}

// AssignPackage assigns a package to a member with payment recording.
func (s *Service) AssignPackage(ctx context.Context, orgID, memberID string, req AssignRequest) (*MemberSubscription, *Payment, error) {
	if err := s.validateAssignRequest(req); err != nil {
		return nil, nil, err
	}

	// Verify the member belongs to this org
	if err := s.repo.VerifyMemberInOrg(ctx, orgID, memberID); err != nil {
		return nil, nil, ErrInvalidMember
	}

	// Fetch and validate the package
	pkg, err := s.repo.GetPackage(ctx, orgID, req.PackageID)
	if err != nil {
		return nil, nil, ErrInvalidPackage
	}
	if !pkg.IsActive {
		return nil, nil, ErrInvalidPackage
	}
	if !pkg.RegistrationEnabled {
		return nil, nil, fmt.Errorf("%w: registration is disabled for this package", ErrInvalidPackage)
	}

	sub, pmt, err := s.repo.AssignPackage(ctx, orgID, memberID, req, pkg)
	if err != nil {
		return nil, nil, fmt.Errorf("assigning package: %w", err)
	}

	s.logger.Info("package assigned",
		"org_id", orgID,
		"member_id", memberID,
		"package_id", req.PackageID,
		"subscription_id", sub.ID,
	)

	return sub, pmt, nil
}

// RenewPackage renews an existing subscription.
func (s *Service) RenewPackage(ctx context.Context, orgID, memberID string, req RenewRequest) (*MemberSubscription, *Payment, error) {
	if err := s.validateRenewRequest(req); err != nil {
		return nil, nil, err
	}

	// Verify membership
	if err := s.repo.VerifyMemberInOrg(ctx, orgID, memberID); err != nil {
		return nil, nil, ErrInvalidMember
	}

	// Fetch existing subscription
	oldMP, pkg, err := s.repo.GetMemberPackage(ctx, orgID, req.MemberPackageID)
	if err != nil {
		return nil, nil, ErrSubscriptionNotFound
	}

	// Verify this subscription belongs to the specified member
	// (GetMemberPackage only filters by org, need to also check user)
	if err := s.verifySubscriptionOwner(ctx, orgID, memberID, oldMP.ID); err != nil {
		return nil, nil, err
	}

	sub, pmt, err := s.repo.RenewPackage(ctx, orgID, memberID, oldMP, pkg, req)
	if err != nil {
		return nil, nil, fmt.Errorf("renewing package: %w", err)
	}

	s.logger.Info("package renewed",
		"org_id", orgID,
		"member_id", memberID,
		"old_subscription_id", oldMP.ID,
		"new_subscription_id", sub.ID,
	)

	return sub, pmt, nil
}

// ExtendPackage extends an existing active subscription.
func (s *Service) ExtendPackage(ctx context.Context, orgID, memberID string, req ExtendRequest) (*MemberSubscription, *Payment, error) {
	if err := s.validateExtendRequest(req); err != nil {
		return nil, nil, err
	}

	// Verify membership
	if err := s.repo.VerifyMemberInOrg(ctx, orgID, memberID); err != nil {
		return nil, nil, ErrInvalidMember
	}

	// Fetch existing subscription
	oldMP, pkg, err := s.repo.GetMemberPackage(ctx, orgID, req.MemberPackageID)
	if err != nil {
		return nil, nil, ErrSubscriptionNotFound
	}

	if err := s.verifySubscriptionOwner(ctx, orgID, memberID, oldMP.ID); err != nil {
		return nil, nil, err
	}

	// Can only extend active subscriptions that haven't expired
	if oldMP.Status != "active" || oldMP.EndDate.Before(time.Now().Truncate(24*time.Hour)) {
		return nil, nil, fmt.Errorf("%w: subscription is not active or has expired", ErrInvalidRequest)
	}

	sub, pmt, err := s.repo.ExtendPackage(ctx, orgID, memberID, oldMP, pkg, req)
	if err != nil {
		return nil, nil, fmt.Errorf("extending package: %w", err)
	}

	s.logger.Info("package extended",
		"org_id", orgID,
		"member_id", memberID,
		"subscription_id", sub.ID,
		"extra_days", req.ExtraDays,
	)

	return sub, pmt, nil
}

// ListPayments returns payment history for a member.
func (s *Service) ListPayments(ctx context.Context, orgID, memberID string, limit, offset int) ([]Payment, int, error) {
	limit, offset = normalizePagination(limit, offset)
	payments, total, err := s.repo.ListPayments(ctx, orgID, memberID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	if payments == nil {
		payments = []Payment{}
	}
	return payments, total, nil
}

// ListSubscriptions returns subscription history for a member.
func (s *Service) ListSubscriptions(ctx context.Context, orgID, memberID string, limit, offset int) ([]MemberSubscription, int, error) {
	limit, offset = normalizePagination(limit, offset)
	subs, total, err := s.repo.ListSubscriptions(ctx, orgID, memberID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	if subs == nil {
		subs = []MemberSubscription{}
	}
	return subs, total, nil
}

func (s *Service) validateAssignRequest(req AssignRequest) error {
	if req.PackageID == "" {
		return fmt.Errorf("%w: package_id is required", ErrInvalidRequest)
	}
	if req.StartDate == "" {
		return fmt.Errorf("%w: start_date is required", ErrInvalidRequest)
	}
	if _, err := time.Parse("2006-01-02", req.StartDate); err != nil {
		return fmt.Errorf("%w: start_date must be in YYYY-MM-DD format", ErrInvalidRequest)
	}
	if req.PaymentMethod != "" && !validPaymentMethods[req.PaymentMethod] {
		return fmt.Errorf("%w: invalid payment_method", ErrInvalidPayment)
	}
	if req.AmountPaid < 0 {
		return fmt.Errorf("%w: amount_paid cannot be negative", ErrInvalidPayment)
	}
	if req.Discount < 0 {
		return fmt.Errorf("%w: discount cannot be negative", ErrInvalidPayment)
	}
	return nil
}

func (s *Service) validateRenewRequest(req RenewRequest) error {
	if req.MemberPackageID == "" {
		return fmt.Errorf("%w: member_package_id is required", ErrInvalidRequest)
	}
	if req.PaymentMethod != "" && !validPaymentMethods[req.PaymentMethod] {
		return fmt.Errorf("%w: invalid payment_method", ErrInvalidPayment)
	}
	if req.AmountPaid < 0 {
		return fmt.Errorf("%w: amount_paid cannot be negative", ErrInvalidPayment)
	}
	if req.Discount < 0 {
		return fmt.Errorf("%w: discount cannot be negative", ErrInvalidPayment)
	}
	return nil
}

func (s *Service) validateExtendRequest(req ExtendRequest) error {
	if req.MemberPackageID == "" {
		return fmt.Errorf("%w: member_package_id is required", ErrInvalidRequest)
	}
	if req.ExtraDays <= 0 {
		return fmt.Errorf("%w: extra_days must be positive", ErrInvalidRequest)
	}
	if req.PaymentMethod != "" && !validPaymentMethods[req.PaymentMethod] {
		return fmt.Errorf("%w: invalid payment_method", ErrInvalidPayment)
	}
	if req.AmountPaid < 0 {
		return fmt.Errorf("%w: amount_paid cannot be negative", ErrInvalidPayment)
	}
	if req.Discount < 0 {
		return fmt.Errorf("%w: discount cannot be negative", ErrInvalidPayment)
	}
	return nil
}

// verifySubscriptionOwner checks that a subscription belongs to the specified member.
func (s *Service) verifySubscriptionOwner(ctx context.Context, orgID, memberID, memberPackageID string) error {
	ownerID, err := s.repo.GetSubscriptionOwner(ctx, orgID, memberPackageID)
	if err != nil {
		return ErrSubscriptionNotFound
	}
	if ownerID != memberID {
		return fmt.Errorf("%w: subscription does not belong to this member", ErrInvalidRequest)
	}
	return nil
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
