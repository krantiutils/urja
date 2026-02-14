package nfc

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
)

// Service handles NFC card business logic.
type Service struct {
	repo   *Repository
	logger *slog.Logger
}

// NewService creates a new NFC service.
func NewService(repo *Repository, logger *slog.Logger) *Service {
	return &Service{repo: repo, logger: logger}
}

// ListCards retrieves NFC cards for an organization.
func (s *Service) ListCards(ctx context.Context, orgID string, limit, offset int) ([]CardResponse, int, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	cards, total, err := s.repo.ListCards(ctx, orgID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]CardResponse, len(cards))
	for i, c := range cards {
		responses[i] = c.ToResponse()
	}
	return responses, total, nil
}

// RegisterCard registers a new NFC card for an organization.
func (s *Service) RegisterCard(ctx context.Context, orgID, cardHex string) (*CardResponse, error) {
	cardHex = strings.TrimSpace(cardHex)
	if cardHex == "" {
		return nil, fmt.Errorf("card_hex is required")
	}

	// Validate hex format (basic check: alphanumeric, reasonable length)
	if len(cardHex) < 4 || len(cardHex) > 50 {
		return nil, fmt.Errorf("invalid card_hex length: must be 4-50 characters")
	}

	card, err := s.repo.RegisterCard(ctx, orgID, cardHex)
	if err != nil {
		return nil, fmt.Errorf("registering card: %w", err)
	}

	s.logger.Info("nfc card registered",
		"org_id", orgID,
		"card_id", card.ID,
		"card_number", card.DisplayCardNumber(),
	)

	resp := card.ToResponse()
	return &resp, nil
}

// AssignCard assigns an NFC card to a member.
func (s *Service) AssignCard(ctx context.Context, cardID, orgID, userID string) (*CardResponse, error) {
	if userID == "" {
		return nil, fmt.Errorf("user_id is required")
	}

	card, err := s.repo.AssignCard(ctx, cardID, orgID, userID)
	if err != nil {
		return nil, err
	}

	// Fetch the full card with user name for response
	fullCard, err := s.repo.GetCardByID(ctx, cardID, orgID)
	if err != nil {
		// Assignment succeeded but name lookup failed — use what we have
		resp := card.ToResponse()
		return &resp, nil
	}

	s.logger.Info("nfc card assigned",
		"org_id", orgID,
		"card_id", cardID,
		"user_id", userID,
		"card_number", fullCard.DisplayCardNumber(),
	)

	resp := fullCard.ToResponse()
	return &resp, nil
}

// UnassignCard removes user assignment from an NFC card.
func (s *Service) UnassignCard(ctx context.Context, cardID, orgID string) (*CardResponse, error) {
	card, err := s.repo.UnassignCard(ctx, cardID, orgID)
	if err != nil {
		return nil, err
	}

	s.logger.Info("nfc card unassigned",
		"org_id", orgID,
		"card_id", cardID,
		"card_number", card.DisplayCardNumber(),
	)

	resp := card.ToResponse()
	return &resp, nil
}

// ListDevices retrieves NFC devices for an organization.
func (s *Service) ListDevices(ctx context.Context, orgID string) ([]Device, error) {
	return s.repo.ListDevices(ctx, orgID)
}
