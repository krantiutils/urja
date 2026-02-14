package dues

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Due represents an outstanding payment owed by a member.
type Due struct {
	ID               string     `json:"id"`
	OrgID            string     `json:"organization_id"`
	UserID           string     `json:"user_id"`
	MemberName       string     `json:"member_name"`
	MemberPhone      string     `json:"member_phone"`
	MemberAvatar     string     `json:"member_avatar,omitempty"`
	Amount           float64    `json:"amount"`
	DueDate          string     `json:"due_date"`
	Description      string     `json:"description,omitempty"`
	Status           string     `json:"status"`
	PaidAt           *time.Time `json:"paid_at,omitempty"`
	PaidAmount       *float64   `json:"paid_amount,omitempty"`
	PaymentMethod    *string    `json:"payment_method,omitempty"`
	PaymentReference *string    `json:"payment_reference,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
}

// PaymentRecord is the result of paying a due.
type PaymentRecord struct {
	DueID            string    `json:"due_id"`
	PaidAmount       float64   `json:"paid_amount"`
	PaymentMethod    string    `json:"payment_method"`
	PaymentReference string    `json:"payment_reference,omitempty"`
	TransactionID    string    `json:"transaction_id"`
	PaidAt           time.Time `json:"paid_at"`
}

// Repository handles dues data persistence.
type Repository struct {
	db *pgxpool.Pool
}

// NewRepository creates a new dues repository.
func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

// ListFilter defines query filters for listing dues.
type ListFilter struct {
	OrgID  string
	Search string
	Status string
	Limit  int
	Offset int
}

// List retrieves dues for an organization with optional filters.
func (r *Repository) List(ctx context.Context, f ListFilter) ([]Due, int, error) {
	var conditions []string
	var args []interface{}
	argIdx := 1

	conditions = append(conditions, fmt.Sprintf("d.organization_id = $%d", argIdx))
	args = append(args, f.OrgID)
	argIdx++

	if f.Status != "" {
		conditions = append(conditions, fmt.Sprintf("d.status = $%d", argIdx))
		args = append(args, f.Status)
		argIdx++
	}

	if f.Search != "" {
		conditions = append(conditions, fmt.Sprintf(
			"(u.name ILIKE $%d OR u.phone ILIKE $%d)",
			argIdx, argIdx,
		))
		args = append(args, "%"+f.Search+"%")
		argIdx++
	}

	where := strings.Join(conditions, " AND ")

	// Count query
	countQuery := fmt.Sprintf(
		`SELECT COUNT(*) FROM dues d JOIN users u ON d.user_id = u.id WHERE %s`,
		where,
	)
	var total int
	if err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("counting dues: %w", err)
	}

	// Data query
	dataQuery := fmt.Sprintf(
		`SELECT d.id, d.organization_id, d.user_id,
		        COALESCE(u.name, ''), u.phone, COALESCE(u.avatar_url, ''),
		        d.amount, d.due_date, COALESCE(d.description, ''),
		        d.status, d.paid_at, d.paid_amount,
		        d.payment_method, d.payment_reference, d.created_at
		 FROM dues d
		 JOIN users u ON d.user_id = u.id
		 WHERE %s
		 ORDER BY d.due_date ASC, d.created_at DESC
		 LIMIT $%d OFFSET $%d`,
		where, argIdx, argIdx+1,
	)
	args = append(args, f.Limit, f.Offset)

	rows, err := r.db.Query(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("querying dues: %w", err)
	}
	defer rows.Close()

	var dues []Due
	for rows.Next() {
		var d Due
		var dueDate time.Time
		if err := rows.Scan(
			&d.ID, &d.OrgID, &d.UserID,
			&d.MemberName, &d.MemberPhone, &d.MemberAvatar,
			&d.Amount, &dueDate, &d.Description,
			&d.Status, &d.PaidAt, &d.PaidAmount,
			&d.PaymentMethod, &d.PaymentReference, &d.CreatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("scanning due: %w", err)
		}
		d.DueDate = dueDate.Format("2006-01-02")
		dues = append(dues, d)
	}

	return dues, total, rows.Err()
}

// Pay marks a due as paid and creates an income transaction atomically.
func (r *Repository) Pay(ctx context.Context, orgID, dueID, payerUserID string, amount float64, method, reference string, entryByUserID string) (*PaymentRecord, error) {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, fmt.Errorf("beginning transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Verify the due belongs to this org and is unpaid
	var dueUserID string
	var dueAmount float64
	var dueStatus string
	var dueDesc string
	err = tx.QueryRow(ctx,
		`SELECT user_id, amount, status, COALESCE(description, 'Membership due')
		 FROM dues WHERE id = $1 AND organization_id = $2`,
		dueID, orgID,
	).Scan(&dueUserID, &dueAmount, &dueStatus, &dueDesc)
	if err != nil {
		return nil, fmt.Errorf("due not found: %w", err)
	}

	if dueStatus == "paid" {
		return nil, fmt.Errorf("due is already paid")
	}
	if dueStatus == "waived" {
		return nil, fmt.Errorf("due has been waived")
	}

	now := time.Now()

	// Mark due as paid
	_, err = tx.Exec(ctx,
		`UPDATE dues SET status = 'paid', paid_at = $1, paid_amount = $2,
		        payment_method = $3, payment_reference = $4
		 WHERE id = $5`,
		now, amount, method, nilIfEmpty(reference), dueID,
	)
	if err != nil {
		return nil, fmt.Errorf("updating due: %w", err)
	}

	// Create income transaction
	var txnID string
	err = tx.QueryRow(ctx,
		`INSERT INTO transactions
		        (organization_id, category, description, transaction_date,
		         transaction_type, amount, payment_type, reference, due_id, entry_by)
		 VALUES ($1, 'Subscription', $2, $3, 'income', $4, $5, $6, $7, $8)
		 RETURNING id`,
		orgID, dueDesc, now, amount, method, nilIfEmpty(reference), dueID, entryByUserID,
	).Scan(&txnID)
	if err != nil {
		return nil, fmt.Errorf("creating transaction: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("committing payment: %w", err)
	}

	return &PaymentRecord{
		DueID:            dueID,
		PaidAmount:       amount,
		PaymentMethod:    method,
		PaymentReference: reference,
		TransactionID:    txnID,
		PaidAt:           now,
	}, nil
}

// BlockAccess deactivates all active NFC cards for a member in an organization.
func (r *Repository) BlockAccess(ctx context.Context, orgID, memberID string) (int64, error) {
	tag, err := r.db.Exec(ctx,
		`UPDATE nfc_cards SET is_active = false, deactivated_at = NOW()
		 WHERE organization_id = $1 AND user_id = $2 AND is_active = true`,
		orgID, memberID,
	)
	if err != nil {
		return 0, fmt.Errorf("blocking NFC access: %w", err)
	}
	return tag.RowsAffected(), nil
}

// GetByID retrieves a single due by ID within an organization.
func (r *Repository) GetByID(ctx context.Context, orgID, dueID string) (*Due, error) {
	var d Due
	var dueDate time.Time
	err := r.db.QueryRow(ctx,
		`SELECT d.id, d.organization_id, d.user_id,
		        COALESCE(u.name, ''), u.phone, COALESCE(u.avatar_url, ''),
		        d.amount, d.due_date, COALESCE(d.description, ''),
		        d.status, d.paid_at, d.paid_amount,
		        d.payment_method, d.payment_reference, d.created_at
		 FROM dues d
		 JOIN users u ON d.user_id = u.id
		 WHERE d.id = $1 AND d.organization_id = $2`,
		dueID, orgID,
	).Scan(
		&d.ID, &d.OrgID, &d.UserID,
		&d.MemberName, &d.MemberPhone, &d.MemberAvatar,
		&d.Amount, &dueDate, &d.Description,
		&d.Status, &d.PaidAt, &d.PaidAmount,
		&d.PaymentMethod, &d.PaymentReference, &d.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("due not found: %w", err)
	}
	d.DueDate = dueDate.Format("2006-01-02")
	return &d, nil
}

func nilIfEmpty(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
