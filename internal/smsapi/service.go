package smsapi

import (
	"math"
	"os"
	"strconv"
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/urja-gym/urja/pkg/sms"
)

// Service handles SMS business logic.
// Server-owned pricing. A gym chooses how many credits to buy; what they cost
// is not the caller's to state.
const (
	defaultCreditRate     = 1.50 // NPR per SMS
	maxCreditsPerPurchase = 100000
)

type Service struct {
	repo       *Repository
	smsClient  *sms.Client
	logger     *slog.Logger
	creditRate float64
}

// NewService creates a new SMS API service.
func NewService(repo *Repository, smsClient *sms.Client, logger *slog.Logger) *Service {
	rate := defaultCreditRate
	// Overridable for a gym on different terms, but still never by the client.
	if v := os.Getenv("SMS_CREDIT_RATE"); v != "" {
		if parsed, err := strconv.ParseFloat(v, 64); err == nil && parsed > 0 {
			rate = parsed
		}
	}
	return &Service{repo: repo, smsClient: smsClient, logger: logger, creditRate: rate}
}

// GetBalance retrieves SMS balance and stats.
func (s *Service) GetBalance(ctx context.Context, orgID string) (*Balance, error) {
	return s.repo.GetBalance(ctx, orgID)
}

// BuyCredits purchases SMS credits for an org.
// BuyCredits records a purchase of SMS credits.
//
// The price is decided here, not by the caller. This previously took the rate
// and the total from the request body and stored both, so anybody who could
// reach the endpoint could award themselves ten thousand credits for one rupee
// and leave a matching entry in the books. Quantity is the only thing the
// caller gets to choose.
func (s *Service) BuyCredits(ctx context.Context, orgID string, quantity int, paymentMethod, purchasedBy string) (*Purchase, error) {
	if quantity <= 0 {
		return nil, fmt.Errorf("quantity must be positive")
	}
	if quantity > maxCreditsPerPurchase {
		return nil, fmt.Errorf("quantity must be %d or fewer", maxCreditsPerPurchase)
	}

	rate := s.creditRate
	amount := math.Round(float64(quantity)*rate*100) / 100

	validMethods := map[string]bool{"khalti": true, "cash": true, "manual": true, "bank_transfer": true}
	if !validMethods[paymentMethod] {
		return nil, fmt.Errorf("invalid payment method: %s", paymentMethod)
	}

	p, err := s.repo.AddCredits(ctx, orgID, quantity, rate, amount, paymentMethod, purchasedBy)
	if err != nil {
		return nil, err
	}

	s.logger.Info("SMS credits purchased",
		"org_id", orgID,
		"quantity", quantity,
		"purchased_by", purchasedBy,
	)
	return p, nil
}

// CreditRate is the price the gym pays per SMS.
func (s *Service) CreditRate() float64 { return s.creditRate }

// SendSMS sends SMS to specified members or all members in the org.
func (s *Service) SendSMS(ctx context.Context, orgID, message, sentBy string, memberIDs []string) (*Campaign, error) {
	message = strings.TrimSpace(message)
	if message == "" {
		return nil, fmt.Errorf("message is required")
	}

	phones, err := s.repo.GetOrgMemberPhones(ctx, orgID, memberIDs)
	if err != nil {
		return nil, fmt.Errorf("getting member phones: %w", err)
	}
	if len(phones) == 0 {
		return nil, fmt.Errorf("no recipients found")
	}

	// Deduct credits first (atomic with campaign recording)
	campaign, err := s.repo.DeductCredits(ctx, orgID, message, sentBy, len(phones))
	if err != nil {
		return nil, err
	}

	// Send SMS to each recipient (best-effort after credits are deducted)
	var sendErrors int
	for _, phone := range phones {
		normalized := sms.NormalizePhone(phone)
		if err := s.smsClient.Send(ctx, normalized, message); err != nil {
			s.logger.Error("failed to send SMS",
				"phone", normalized[:4]+"******",
				"campaign_id", campaign.ID,
				"error", err,
			)
			sendErrors++
		}
	}

	if sendErrors > 0 {
		s.logger.Warn("SMS campaign completed with errors",
			"campaign_id", campaign.ID,
			"total", len(phones),
			"failed", sendErrors,
		)
	} else {
		s.logger.Info("SMS campaign sent successfully",
			"campaign_id", campaign.ID,
			"recipients", len(phones),
		)
	}

	return campaign, nil
}

// ListPurchases retrieves purchase history with pagination.
func (s *Service) ListPurchases(ctx context.Context, orgID string, limit, offset int) ([]Purchase, error) {
	if limit <= 0 || limit > 50 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	return s.repo.ListPurchases(ctx, orgID, limit, offset)
}
