package nfc

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Card represents an NFC card registered to an organization.
// SECURITY: card_hex is never exposed in API responses.
type Card struct {
	ID           string     `json:"id"`
	OrgID        string     `json:"org_id"`
	UserID       *string    `json:"user_id,omitempty"`
	CardNumber   int        `json:"card_number"`
	UserName     string     `json:"user_name,omitempty"`
	IsActive     bool       `json:"is_active"`
	AssignedAt   *time.Time `json:"assigned_at,omitempty"`
	DeactivatedAt *time.Time `json:"deactivated_at,omitempty"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// DisplayCardNumber returns the formatted card number (e.g. #CG0001).
func (c Card) DisplayCardNumber() string {
	return fmt.Sprintf("#CG%04d", c.CardNumber)
}

// CardResponse is the API response for a card. Never includes card_hex.
type CardResponse struct {
	ID            string     `json:"id"`
	OrgID         string     `json:"org_id"`
	CardNumber    string     `json:"card_number"`
	AssignedUser  *string    `json:"assigned_member,omitempty"`
	UserName      string     `json:"full_name,omitempty"`
	Status        string     `json:"status"`
	AssignedAt    *time.Time `json:"assigned_at,omitempty"`
	LastUpdated   time.Time  `json:"last_updated"`
}

// ToResponse converts a Card to its API representation.
func (c Card) ToResponse() CardResponse {
	status := "active"
	if !c.IsActive {
		status = "inactive"
	} else if c.UserID == nil {
		status = "unassigned"
	} else {
		status = "assigned"
	}

	return CardResponse{
		ID:           c.ID,
		OrgID:        c.OrgID,
		CardNumber:   c.DisplayCardNumber(),
		AssignedUser: c.UserID,
		UserName:     c.UserName,
		Status:       status,
		AssignedAt:   c.AssignedAt,
		LastUpdated:  c.UpdatedAt,
	}
}

// Device represents an NFC reader device.
type Device struct {
	ID               string     `json:"id"`
	OrgID            string     `json:"org_id"`
	Name             string     `json:"name"`
	DeviceIdentifier string     `json:"device_identifier"`
	Status           string     `json:"status"`
	DoorState        string     `json:"door_state"`
	LastSyncAt       *time.Time `json:"last_sync_at,omitempty"`
	UptimeSeconds    int64      `json:"uptime_seconds"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

// Repository handles NFC card and device data persistence.
type Repository struct {
	db *pgxpool.Pool
}

// NewRepository creates a new NFC repository.
func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

// ListCards retrieves NFC cards for an organization with user names.
func (r *Repository) ListCards(ctx context.Context, orgID string, limit, offset int) ([]Card, int, error) {
	var total int
	err := r.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM nfc_cards WHERE organization_id = $1`,
		orgID,
	).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("counting nfc cards: %w", err)
	}

	rows, err := r.db.Query(ctx,
		`SELECT c.id, c.organization_id, c.user_id, c.card_number,
		        COALESCE(u.name, ''), c.is_active, c.assigned_at,
		        c.deactivated_at, c.updated_at
		 FROM nfc_cards c
		 LEFT JOIN users u ON c.user_id = u.id
		 WHERE c.organization_id = $1
		 ORDER BY c.card_number ASC
		 LIMIT $2 OFFSET $3`,
		orgID, limit, offset,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("listing nfc cards: %w", err)
	}
	defer rows.Close()

	var cards []Card
	for rows.Next() {
		var c Card
		if err := rows.Scan(
			&c.ID, &c.OrgID, &c.UserID, &c.CardNumber,
			&c.UserName, &c.IsActive, &c.AssignedAt,
			&c.DeactivatedAt, &c.UpdatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("scanning nfc card: %w", err)
		}
		cards = append(cards, c)
	}
	return cards, total, rows.Err()
}

// RegisterCard creates a new NFC card for the organization.
// The card_hex is stored but NEVER exposed in API responses.
func (r *Repository) RegisterCard(ctx context.Context, orgID, cardHex string) (*Card, error) {
	var c Card
	err := r.db.QueryRow(ctx,
		`INSERT INTO nfc_cards (organization_id, card_hex)
		 VALUES ($1, $2)
		 RETURNING id, organization_id, user_id, card_number, is_active, assigned_at, deactivated_at, updated_at`,
		orgID, cardHex,
	).Scan(&c.ID, &c.OrgID, &c.UserID, &c.CardNumber, &c.IsActive, &c.AssignedAt, &c.DeactivatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("registering nfc card: %w", err)
	}
	return &c, nil
}

// AssignCard assigns a card to a member.
func (r *Repository) AssignCard(ctx context.Context, cardID, orgID, userID string) (*Card, error) {
	var c Card
	err := r.db.QueryRow(ctx,
		`UPDATE nfc_cards
		 SET user_id = $1, assigned_at = NOW(), updated_at = NOW()
		 WHERE id = $2 AND organization_id = $3 AND is_active = true
		 RETURNING id, organization_id, user_id, card_number, is_active, assigned_at, deactivated_at, updated_at`,
		userID, cardID, orgID,
	).Scan(&c.ID, &c.OrgID, &c.UserID, &c.CardNumber, &c.IsActive, &c.AssignedAt, &c.DeactivatedAt, &c.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("card not found or inactive")
		}
		return nil, fmt.Errorf("assigning nfc card: %w", err)
	}
	return &c, nil
}

// UnassignCard removes user assignment from a card.
func (r *Repository) UnassignCard(ctx context.Context, cardID, orgID string) (*Card, error) {
	var c Card
	err := r.db.QueryRow(ctx,
		`UPDATE nfc_cards
		 SET user_id = NULL, assigned_at = NULL, updated_at = NOW()
		 WHERE id = $1 AND organization_id = $2 AND is_active = true
		 RETURNING id, organization_id, user_id, card_number, is_active, assigned_at, deactivated_at, updated_at`,
		cardID, orgID,
	).Scan(&c.ID, &c.OrgID, &c.UserID, &c.CardNumber, &c.IsActive, &c.AssignedAt, &c.DeactivatedAt, &c.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("card not found or inactive")
		}
		return nil, fmt.Errorf("unassigning nfc card: %w", err)
	}
	return &c, nil
}

// GetCardByID retrieves a single card (for validation).
func (r *Repository) GetCardByID(ctx context.Context, cardID, orgID string) (*Card, error) {
	var c Card
	err := r.db.QueryRow(ctx,
		`SELECT c.id, c.organization_id, c.user_id, c.card_number,
		        COALESCE(u.name, ''), c.is_active, c.assigned_at,
		        c.deactivated_at, c.updated_at
		 FROM nfc_cards c
		 LEFT JOIN users u ON c.user_id = u.id
		 WHERE c.id = $1 AND c.organization_id = $2`,
		cardID, orgID,
	).Scan(&c.ID, &c.OrgID, &c.UserID, &c.CardNumber,
		&c.UserName, &c.IsActive, &c.AssignedAt,
		&c.DeactivatedAt, &c.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("card not found")
		}
		return nil, fmt.Errorf("getting nfc card: %w", err)
	}
	return &c, nil
}

// ListDevices retrieves NFC devices for an organization.
func (r *Repository) ListDevices(ctx context.Context, orgID string) ([]Device, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, organization_id, name, device_identifier,
		        status, door_state, last_sync_at, uptime_seconds,
		        created_at, updated_at
		 FROM nfc_devices
		 WHERE organization_id = $1
		 ORDER BY name ASC`,
		orgID,
	)
	if err != nil {
		return nil, fmt.Errorf("listing nfc devices: %w", err)
	}
	defer rows.Close()

	var devices []Device
	for rows.Next() {
		var d Device
		if err := rows.Scan(
			&d.ID, &d.OrgID, &d.Name, &d.DeviceIdentifier,
			&d.Status, &d.DoorState, &d.LastSyncAt, &d.UptimeSeconds,
			&d.CreatedAt, &d.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scanning nfc device: %w", err)
		}
		devices = append(devices, d)
	}
	return devices, rows.Err()
}
