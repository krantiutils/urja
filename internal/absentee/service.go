package absentee

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/urja-gym/urja/pkg/sms"
)

// Service handles absentee business logic.
type Service struct {
	repo      *Repository
	smsClient *sms.Client
	logger    *slog.Logger
}

// NewService creates a new absentee service.
func NewService(repo *Repository, smsClient *sms.Client, logger *slog.Logger) *Service {
	return &Service{repo: repo, smsClient: smsClient, logger: logger}
}

// List retrieves absentee members with pagination.
func (s *Service) List(ctx context.Context, orgID string, minDays, limit, offset int) ([]Absentee, error) {
	if minDays <= 0 {
		minDays = 7
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	return s.repo.List(ctx, orgID, minDays, limit, offset)
}

// Notify sends an SMS notification to an absent member.
func (s *Service) Notify(ctx context.Context, orgID, memberID, orgName string) error {
	phone, err := s.repo.GetMemberPhone(ctx, orgID, memberID)
	if err != nil {
		return err
	}

	normalized := sms.NormalizePhone(phone)
	if !sms.ValidatePhone(normalized) {
		return fmt.Errorf("member has invalid phone number")
	}

	message := fmt.Sprintf("We miss you at %s! Come back and continue your fitness journey. - Urja", orgName)
	if err := s.smsClient.Send(ctx, normalized, message); err != nil {
		return fmt.Errorf("sending notification: %w", err)
	}

	s.logger.Info("absentee notified",
		"member_id", memberID,
		"org_id", orgID,
		"phone", normalized[:4]+"******",
	)
	return nil
}
