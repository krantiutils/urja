package accounts

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Transaction represents a financial ledger entry.
type Transaction struct {
	ID              string    `json:"id"`
	OrgID           string    `json:"organization_id"`
	Category        string    `json:"category"`
	Description     string    `json:"description,omitempty"`
	TransactionDate string    `json:"transaction_date"`
	TransactionType string    `json:"transaction_type"`
	Amount          float64   `json:"amount"`
	PaymentType     string    `json:"payment_type"`
	Reference       string    `json:"reference,omitempty"`
	DueID           *string   `json:"due_id,omitempty"`
	EntryBy         string    `json:"entry_by"`
	EntryByName     string    `json:"entry_by_name,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// Summary holds aggregated financial data.
type Summary struct {
	TotalIncome   float64 `json:"total_income"`
	TotalExpenses float64 `json:"total_expenses"`
	GrossProfit   float64 `json:"gross_profit"`
	ProfitPercent float64 `json:"profit_percent"`
}

// Repository handles transaction data persistence.
type Repository struct {
	db *pgxpool.Pool
}

// NewRepository creates a new accounts repository.
func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

// ListFilter defines query filters for listing transactions.
type ListFilter struct {
	OrgID           string
	From            string
	To              string
	Category        string
	TransactionType string
	PaymentType     string
	Limit           int
	Offset          int
}

// List retrieves transactions for an organization with optional filters.
func (r *Repository) List(ctx context.Context, f ListFilter) ([]Transaction, int, error) {
	var conditions []string
	var args []interface{}
	argIdx := 1

	conditions = append(conditions, fmt.Sprintf("t.organization_id = $%d", argIdx))
	args = append(args, f.OrgID)
	argIdx++

	if f.From != "" {
		conditions = append(conditions, fmt.Sprintf("t.transaction_date >= $%d", argIdx))
		args = append(args, f.From)
		argIdx++
	}

	if f.To != "" {
		conditions = append(conditions, fmt.Sprintf("t.transaction_date <= $%d", argIdx))
		args = append(args, f.To)
		argIdx++
	}

	if f.Category != "" {
		conditions = append(conditions, fmt.Sprintf("t.category = $%d", argIdx))
		args = append(args, f.Category)
		argIdx++
	}

	if f.TransactionType != "" {
		conditions = append(conditions, fmt.Sprintf("t.transaction_type = $%d", argIdx))
		args = append(args, f.TransactionType)
		argIdx++
	}

	if f.PaymentType != "" {
		conditions = append(conditions, fmt.Sprintf("t.payment_type = $%d", argIdx))
		args = append(args, f.PaymentType)
		argIdx++
	}

	where := strings.Join(conditions, " AND ")

	// Count
	var total int
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM transactions t WHERE %s`, where)
	if err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("counting transactions: %w", err)
	}

	// Data
	dataQuery := fmt.Sprintf(
		`SELECT t.id, t.organization_id, t.category,
		        COALESCE(t.description, ''), t.transaction_date,
		        t.transaction_type, t.amount, t.payment_type,
		        COALESCE(t.reference, ''), t.due_id,
		        t.entry_by, COALESCE(u.name, ''),
		        t.created_at, t.updated_at
		 FROM transactions t
		 JOIN users u ON t.entry_by = u.id
		 WHERE %s
		 ORDER BY t.transaction_date DESC, t.created_at DESC
		 LIMIT $%d OFFSET $%d`,
		where, argIdx, argIdx+1,
	)
	args = append(args, f.Limit, f.Offset)

	rows, err := r.db.Query(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("querying transactions: %w", err)
	}
	defer rows.Close()

	var txns []Transaction
	for rows.Next() {
		var t Transaction
		var txnDate time.Time
		if err := rows.Scan(
			&t.ID, &t.OrgID, &t.Category,
			&t.Description, &txnDate,
			&t.TransactionType, &t.Amount, &t.PaymentType,
			&t.Reference, &t.DueID,
			&t.EntryBy, &t.EntryByName,
			&t.CreatedAt, &t.UpdatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("scanning transaction: %w", err)
		}
		t.TransactionDate = txnDate.Format("2006-01-02")
		txns = append(txns, t)
	}

	return txns, total, rows.Err()
}

// Create inserts a new transaction.
func (r *Repository) Create(ctx context.Context, orgID, category, description, txnDate, txnType string, amount float64, paymentType, reference, entryBy string) (*Transaction, error) {
	var t Transaction
	var dateVal time.Time
	err := r.db.QueryRow(ctx,
		`INSERT INTO transactions
		        (organization_id, category, description, transaction_date,
		         transaction_type, amount, payment_type, reference, entry_by)
		 VALUES ($1, $2, NULLIF($3, ''), $4, $5, $6, $7, NULLIF($8, ''), $9)
		 RETURNING id, organization_id, category, COALESCE(description, ''),
		          transaction_date, transaction_type, amount, payment_type,
		          COALESCE(reference, ''), entry_by, created_at, updated_at`,
		orgID, category, description, txnDate, txnType, amount, paymentType, reference, entryBy,
	).Scan(
		&t.ID, &t.OrgID, &t.Category, &t.Description,
		&dateVal, &t.TransactionType, &t.Amount, &t.PaymentType,
		&t.Reference, &t.EntryBy, &t.CreatedAt, &t.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("creating transaction: %w", err)
	}
	t.TransactionDate = dateVal.Format("2006-01-02")
	return &t, nil
}

// Update modifies an existing transaction.
func (r *Repository) Update(ctx context.Context, orgID, txnID, category, description, txnDate, txnType string, amount float64, paymentType, reference string) (*Transaction, error) {
	var t Transaction
	var dateVal time.Time
	err := r.db.QueryRow(ctx,
		`UPDATE transactions SET
		        category = $3, description = NULLIF($4, ''),
		        transaction_date = $5, transaction_type = $6,
		        amount = $7, payment_type = $8, reference = NULLIF($9, '')
		 WHERE id = $1 AND organization_id = $2
		 RETURNING id, organization_id, category, COALESCE(description, ''),
		          transaction_date, transaction_type, amount, payment_type,
		          COALESCE(reference, ''), entry_by, created_at, updated_at`,
		txnID, orgID, category, description, txnDate, txnType, amount, paymentType, reference,
	).Scan(
		&t.ID, &t.OrgID, &t.Category, &t.Description,
		&dateVal, &t.TransactionType, &t.Amount, &t.PaymentType,
		&t.Reference, &t.EntryBy, &t.CreatedAt, &t.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("updating transaction: %w", err)
	}
	t.TransactionDate = dateVal.Format("2006-01-02")
	return &t, nil
}

// Delete removes a transaction.
func (r *Repository) Delete(ctx context.Context, orgID, txnID string) error {
	tag, err := r.db.Exec(ctx,
		`DELETE FROM transactions WHERE id = $1 AND organization_id = $2`,
		txnID, orgID,
	)
	if err != nil {
		return fmt.Errorf("deleting transaction: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("transaction not found")
	}
	return nil
}

// GetSummary computes income, expense, and profit totals for an organization.
func (r *Repository) GetSummary(ctx context.Context, orgID, from, to string) (*Summary, error) {
	var conditions []string
	var args []interface{}
	argIdx := 1

	conditions = append(conditions, fmt.Sprintf("organization_id = $%d", argIdx))
	args = append(args, orgID)
	argIdx++

	if from != "" {
		conditions = append(conditions, fmt.Sprintf("transaction_date >= $%d", argIdx))
		args = append(args, from)
		argIdx++
	}

	if to != "" {
		conditions = append(conditions, fmt.Sprintf("transaction_date <= $%d", argIdx))
		args = append(args, to)
		argIdx++
	}

	where := strings.Join(conditions, " AND ")

	query := fmt.Sprintf(
		`SELECT
		    COALESCE(SUM(CASE WHEN transaction_type = 'income' THEN amount ELSE 0 END), 0) AS total_income,
		    COALESCE(SUM(CASE WHEN transaction_type = 'expense' THEN amount ELSE 0 END), 0) AS total_expenses
		 FROM transactions
		 WHERE %s`,
		where,
	)

	var s Summary
	if err := r.db.QueryRow(ctx, query, args...).Scan(&s.TotalIncome, &s.TotalExpenses); err != nil {
		return nil, fmt.Errorf("computing summary: %w", err)
	}

	s.GrossProfit = s.TotalIncome - s.TotalExpenses
	if s.TotalIncome > 0 {
		s.ProfitPercent = (s.GrossProfit / s.TotalIncome) * 100
	}

	return &s, nil
}

// ListForExport retrieves all matching transactions without pagination (for CSV export).
func (r *Repository) ListForExport(ctx context.Context, f ListFilter) ([]Transaction, error) {
	var conditions []string
	var args []interface{}
	argIdx := 1

	conditions = append(conditions, fmt.Sprintf("t.organization_id = $%d", argIdx))
	args = append(args, f.OrgID)
	argIdx++

	if f.From != "" {
		conditions = append(conditions, fmt.Sprintf("t.transaction_date >= $%d", argIdx))
		args = append(args, f.From)
		argIdx++
	}

	if f.To != "" {
		conditions = append(conditions, fmt.Sprintf("t.transaction_date <= $%d", argIdx))
		args = append(args, f.To)
		argIdx++
	}

	if f.Category != "" {
		conditions = append(conditions, fmt.Sprintf("t.category = $%d", argIdx))
		args = append(args, f.Category)
		argIdx++
	}

	if f.TransactionType != "" {
		conditions = append(conditions, fmt.Sprintf("t.transaction_type = $%d", argIdx))
		args = append(args, f.TransactionType)
		argIdx++
	}

	if f.PaymentType != "" {
		conditions = append(conditions, fmt.Sprintf("t.payment_type = $%d", argIdx))
		args = append(args, f.PaymentType)
		argIdx++
	}

	where := strings.Join(conditions, " AND ")

	query := fmt.Sprintf(
		`SELECT t.id, t.organization_id, t.category,
		        COALESCE(t.description, ''), t.transaction_date,
		        t.transaction_type, t.amount, t.payment_type,
		        COALESCE(t.reference, ''), t.due_id,
		        t.entry_by, COALESCE(u.name, ''),
		        t.created_at, t.updated_at
		 FROM transactions t
		 JOIN users u ON t.entry_by = u.id
		 WHERE %s
		 ORDER BY t.transaction_date DESC, t.created_at DESC`,
		where,
	)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("querying transactions for export: %w", err)
	}
	defer rows.Close()

	var txns []Transaction
	for rows.Next() {
		var t Transaction
		var txnDate time.Time
		if err := rows.Scan(
			&t.ID, &t.OrgID, &t.Category,
			&t.Description, &txnDate,
			&t.TransactionType, &t.Amount, &t.PaymentType,
			&t.Reference, &t.DueID,
			&t.EntryBy, &t.EntryByName,
			&t.CreatedAt, &t.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scanning transaction: %w", err)
		}
		t.TransactionDate = txnDate.Format("2006-01-02")
		txns = append(txns, t)
	}

	return txns, rows.Err()
}
