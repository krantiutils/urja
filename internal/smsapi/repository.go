package smsapi

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Balance represents an org's SMS credit balance.
type Balance struct {
	OrgID          string `json:"organization_id"`
	Balance        int    `json:"balance"`
	TotalPurchased int    `json:"total_purchased"`
	TotalUsed      int    `json:"total_used"`
	TotalCampaigns int    `json:"total_campaigns"`
}

// Purchase represents an SMS credit purchase record.
type Purchase struct {
	ID            string    `json:"id"`
	OrgID         string    `json:"organization_id"`
	Quantity      int       `json:"quantity"`
	Rate          float64   `json:"rate"`
	Amount        float64   `json:"amount"`
	PaymentMethod string    `json:"payment_method"`
	PaymentStatus string    `json:"payment_status"`
	PurchasedBy   string    `json:"purchased_by"`
	CreatedAt     time.Time `json:"created_at"`
}

// Campaign represents an SMS campaign (bulk send).
type Campaign struct {
	ID             string    `json:"id"`
	OrgID          string    `json:"organization_id"`
	Message        string    `json:"message"`
	RecipientCount int       `json:"recipient_count"`
	SentBy         string    `json:"sent_by"`
	CreatedAt      time.Time `json:"created_at"`
}

// Repository handles SMS data persistence.
type Repository struct {
	db *pgxpool.Pool
}

// NewRepository creates a new SMS API repository.
func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

// GetBalance retrieves the SMS credit balance for an org.
// Returns zero-balance if no record exists yet.
func (r *Repository) GetBalance(ctx context.Context, orgID string) (*Balance, error) {
	var b Balance
	b.OrgID = orgID

	err := r.db.QueryRow(ctx,
		`SELECT COALESCE(balance, 0), COALESCE(total_purchased, 0), COALESCE(total_used, 0)
		 FROM sms_credits WHERE organization_id = $1`,
		orgID,
	).Scan(&b.Balance, &b.TotalPurchased, &b.TotalUsed)
	if err != nil {
		// No record yet means zero balance
		b.Balance = 0
		b.TotalPurchased = 0
		b.TotalUsed = 0
	}

	// Count campaigns
	err = r.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM sms_campaigns WHERE organization_id = $1`,
		orgID,
	).Scan(&b.TotalCampaigns)
	if err != nil {
		b.TotalCampaigns = 0
	}

	return &b, nil
}

// EnsureCreditsRow creates the sms_credits row for an org if it doesn't exist.
func (r *Repository) EnsureCreditsRow(ctx context.Context, orgID string) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO sms_credits (organization_id) VALUES ($1) ON CONFLICT (organization_id) DO NOTHING`,
		orgID,
	)
	if err != nil {
		return fmt.Errorf("ensuring credits row: %w", err)
	}
	return nil
}

// AddCredits atomically adds credits and records the purchase.
func (r *Repository) AddCredits(ctx context.Context, orgID string, quantity int, rate, amount float64, paymentMethod, purchasedBy string) (*Purchase, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("starting transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Ensure credits row exists
	_, err = tx.Exec(ctx,
		`INSERT INTO sms_credits (organization_id) VALUES ($1) ON CONFLICT (organization_id) DO NOTHING`,
		orgID,
	)
	if err != nil {
		return nil, fmt.Errorf("ensuring credits row: %w", err)
	}

	// Update balance
	_, err = tx.Exec(ctx,
		`UPDATE sms_credits SET balance = balance + $2, total_purchased = total_purchased + $2
		 WHERE organization_id = $1`,
		orgID, quantity,
	)
	if err != nil {
		return nil, fmt.Errorf("updating credits: %w", err)
	}

	// Record purchase
	var p Purchase
	err = tx.QueryRow(ctx,
		`INSERT INTO sms_purchases (organization_id, quantity, rate, amount, payment_method, payment_status, purchased_by)
		 VALUES ($1, $2, $3, $4, $5, 'completed', $6)
		 RETURNING id, organization_id, quantity, rate, amount, payment_method, payment_status, purchased_by, created_at`,
		orgID, quantity, rate, amount, paymentMethod, purchasedBy,
	).Scan(&p.ID, &p.OrgID, &p.Quantity, &p.Rate, &p.Amount, &p.PaymentMethod, &p.PaymentStatus, &p.PurchasedBy, &p.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("recording purchase: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("committing transaction: %w", err)
	}

	return &p, nil
}

// DeductCredits atomically deducts credits and records a campaign.
// Returns an error if insufficient balance.
func (r *Repository) DeductCredits(ctx context.Context, orgID, message, sentBy string, recipientCount int) (*Campaign, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("starting transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Check and deduct balance
	var newBalance int
	err = tx.QueryRow(ctx,
		`UPDATE sms_credits SET balance = balance - $2, total_used = total_used + $2
		 WHERE organization_id = $1 AND balance >= $2
		 RETURNING balance`,
		orgID, recipientCount,
	).Scan(&newBalance)
	if err != nil {
		return nil, fmt.Errorf("insufficient SMS credits")
	}

	// Record campaign
	var c Campaign
	err = tx.QueryRow(ctx,
		`INSERT INTO sms_campaigns (organization_id, message, recipient_count, sent_by)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, organization_id, message, recipient_count, sent_by, created_at`,
		orgID, message, recipientCount, sentBy,
	).Scan(&c.ID, &c.OrgID, &c.Message, &c.RecipientCount, &c.SentBy, &c.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("recording campaign: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("committing transaction: %w", err)
	}

	return &c, nil
}

// ListPurchases retrieves SMS purchase history for an org.
func (r *Repository) ListPurchases(ctx context.Context, orgID string, limit, offset int) ([]Purchase, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, organization_id, quantity, rate, amount, payment_method, payment_status, purchased_by, created_at
		 FROM sms_purchases WHERE organization_id = $1
		 ORDER BY created_at DESC LIMIT $2 OFFSET $3`,
		orgID, limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("listing purchases: %w", err)
	}
	defer rows.Close()

	var purchases []Purchase
	for rows.Next() {
		var p Purchase
		if err := rows.Scan(&p.ID, &p.OrgID, &p.Quantity, &p.Rate, &p.Amount,
			&p.PaymentMethod, &p.PaymentStatus, &p.PurchasedBy, &p.CreatedAt); err != nil {
			return nil, fmt.Errorf("scanning purchase: %w", err)
		}
		purchases = append(purchases, p)
	}
	return purchases, rows.Err()
}

// GetOrgMemberPhones retrieves phone numbers of active members in the org.
// If memberIDs is non-empty, only retrieves those specific members.
func (r *Repository) GetOrgMemberPhones(ctx context.Context, orgID string, memberIDs []string) ([]string, error) {
	var query string
	var args []interface{}

	if len(memberIDs) > 0 {
		query = `SELECT u.phone FROM users u
			 JOIN organization_members om ON om.user_id = u.id
			 WHERE om.organization_id = $1 AND om.status = 'active' AND u.id = ANY($2)`
		args = []interface{}{orgID, memberIDs}
	} else {
		query = `SELECT u.phone FROM users u
			 JOIN organization_members om ON om.user_id = u.id
			 WHERE om.organization_id = $1 AND om.status = 'active'`
		args = []interface{}{orgID}
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("querying member phones: %w", err)
	}
	defer rows.Close()

	var phones []string
	for rows.Next() {
		var phone string
		if err := rows.Scan(&phone); err != nil {
			return nil, fmt.Errorf("scanning phone: %w", err)
		}
		phones = append(phones, phone)
	}
	return phones, rows.Err()
}
