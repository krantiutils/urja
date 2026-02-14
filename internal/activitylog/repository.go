package activitylog

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// EventType defines valid activity log event types.
type EventType string

const (
	EventCheckIn             EventType = "check_in"
	EventSubscriptionApproved EventType = "subscription_approved"
	EventMembershipRequest   EventType = "membership_request"
	EventMembershipApproved  EventType = "membership_approved"
	EventPackageChanged      EventType = "package_changed"
	EventNFCCardAssigned     EventType = "nfc_card_assigned"
	EventNFCCardUnassigned   EventType = "nfc_card_unassigned"
	EventNFCCardRegistered   EventType = "nfc_card_registered"
	EventMemberJoined        EventType = "member_joined"
	EventMemberSuspended     EventType = "member_suspended"
)

// Entry represents a single activity log entry.
type Entry struct {
	ID          string            `json:"id"`
	OrgID       string            `json:"org_id"`
	UserID      *string           `json:"user_id,omitempty"`
	EventType   string            `json:"event_type"`
	Description string            `json:"description"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt   time.Time         `json:"created_at"`
}

// Repository handles activity log data persistence.
type Repository struct {
	db *pgxpool.Pool
}

// NewRepository creates a new activity log repository.
func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

// Create appends a new activity log entry.
func (r *Repository) Create(ctx context.Context, orgID string, userID *string, eventType EventType, description string, metadata map[string]interface{}) (*Entry, error) {
	var e Entry
	err := r.db.QueryRow(ctx,
		`INSERT INTO activity_logs (organization_id, user_id, event_type, description, metadata)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id, organization_id, user_id, event_type, description, metadata, created_at`,
		orgID, userID, string(eventType), description, metadata,
	).Scan(&e.ID, &e.OrgID, &e.UserID, &e.EventType, &e.Description, &e.Metadata, &e.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("creating activity log: %w", err)
	}
	return &e, nil
}

// List retrieves activity logs for an organization with optional date filtering.
func (r *Repository) List(ctx context.Context, orgID string, from, to *time.Time, limit, offset int) ([]Entry, int, error) {
	// Build count query
	countQuery := `SELECT COUNT(*) FROM activity_logs WHERE organization_id = $1`
	countArgs := []interface{}{orgID}
	argIdx := 2

	if from != nil {
		countQuery += fmt.Sprintf(` AND created_at >= $%d`, argIdx)
		countArgs = append(countArgs, *from)
		argIdx++
	}
	if to != nil {
		countQuery += fmt.Sprintf(` AND created_at <= $%d`, argIdx)
		countArgs = append(countArgs, *to)
		argIdx++
	}

	var total int
	if err := r.db.QueryRow(ctx, countQuery, countArgs...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("counting activity logs: %w", err)
	}

	// Build data query
	dataQuery := `SELECT id, organization_id, user_id, event_type, description, metadata, created_at
		 FROM activity_logs WHERE organization_id = $1`
	dataArgs := []interface{}{orgID}
	argIdx = 2

	if from != nil {
		dataQuery += fmt.Sprintf(` AND created_at >= $%d`, argIdx)
		dataArgs = append(dataArgs, *from)
		argIdx++
	}
	if to != nil {
		dataQuery += fmt.Sprintf(` AND created_at <= $%d`, argIdx)
		dataArgs = append(dataArgs, *to)
		argIdx++
	}

	dataQuery += fmt.Sprintf(` ORDER BY created_at DESC LIMIT $%d OFFSET $%d`, argIdx, argIdx+1)
	dataArgs = append(dataArgs, limit, offset)

	rows, err := r.db.Query(ctx, dataQuery, dataArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("listing activity logs: %w", err)
	}
	defer rows.Close()

	var entries []Entry
	for rows.Next() {
		var e Entry
		if err := rows.Scan(&e.ID, &e.OrgID, &e.UserID, &e.EventType, &e.Description, &e.Metadata, &e.CreatedAt); err != nil {
			return nil, 0, fmt.Errorf("scanning activity log: %w", err)
		}
		entries = append(entries, e)
	}
	return entries, total, rows.Err()
}
