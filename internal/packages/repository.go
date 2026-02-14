package packages

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Package represents a gym membership package.
type Package struct {
	ID             string          `json:"id"`
	OrganizationID string          `json:"organization_id"`
	Name           string          `json:"name"`
	NameNe         *string         `json:"name_ne,omitempty"`
	Description    *string         `json:"description,omitempty"`
	DescriptionNe  *string         `json:"description_ne,omitempty"`
	DurationDays   int             `json:"duration_days"`
	Price          string          `json:"price"` // decimal as string to preserve precision
	Currency       string          `json:"currency"`
	MaxMembers     *int            `json:"max_members,omitempty"`
	Features       json.RawMessage `json:"features"`
	IsActive       bool            `json:"is_active"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
}

// MemberPackage represents a member's subscription to a package.
type MemberPackage struct {
	ID               string    `json:"id"`
	UserID           string    `json:"user_id"`
	PackageID        string    `json:"package_id"`
	OrganizationID   string    `json:"organization_id"`
	StartDate        time.Time `json:"start_date"`
	EndDate          time.Time `json:"end_date"`
	PaymentMethod    *string   `json:"payment_method,omitempty"`
	PaymentReference *string   `json:"payment_reference,omitempty"`
	AmountPaid       string    `json:"amount_paid"`
	Status           string    `json:"status"`
	CreatedAt        time.Time `json:"created_at"`

	// Joined fields for list queries.
	PackageName *string `json:"package_name,omitempty"`
	OrgName     *string `json:"org_name,omitempty"`
}

// Repository handles package data persistence.
type Repository struct {
	db *pgxpool.Pool
}

// NewRepository creates a new packages repository.
func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

// ListPublic returns all active packages across all organizations.
func (r *Repository) ListPublic(ctx context.Context, limit, offset int) ([]Package, error) {
	rows, err := r.db.Query(ctx,
		`SELECT p.id, p.organization_id, p.name, p.name_ne, p.description, p.description_ne,
		        p.duration_days, p.price, p.currency, p.max_members, p.features,
		        p.is_active, p.created_at, p.updated_at
		 FROM packages p
		 JOIN organizations o ON o.id = p.organization_id AND o.is_active = true
		 WHERE p.is_active = true
		 ORDER BY p.created_at DESC
		 LIMIT $1 OFFSET $2`,
		limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("querying public packages: %w", err)
	}
	defer rows.Close()
	return scanPackages(rows)
}

// ListByOrg returns all packages for an organization (including inactive, for admins).
func (r *Repository) ListByOrg(ctx context.Context, orgID string, limit, offset int) ([]Package, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, organization_id, name, name_ne, description, description_ne,
		        duration_days, price, currency, max_members, features,
		        is_active, created_at, updated_at
		 FROM packages
		 WHERE organization_id = $1
		 ORDER BY is_active DESC, created_at DESC
		 LIMIT $2 OFFSET $3`,
		orgID, limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("querying org packages: %w", err)
	}
	defer rows.Close()
	return scanPackages(rows)
}

// GetByID returns a single package by ID.
func (r *Repository) GetByID(ctx context.Context, id string) (*Package, error) {
	var p Package
	err := r.db.QueryRow(ctx,
		`SELECT id, organization_id, name, name_ne, description, description_ne,
		        duration_days, price, currency, max_members, features,
		        is_active, created_at, updated_at
		 FROM packages WHERE id = $1`,
		id,
	).Scan(&p.ID, &p.OrganizationID, &p.Name, &p.NameNe, &p.Description, &p.DescriptionNe,
		&p.DurationDays, &p.Price, &p.Currency, &p.MaxMembers, &p.Features,
		&p.IsActive, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("getting package %s: %w", id, err)
	}
	return &p, nil
}

// Create inserts a new package.
func (r *Repository) Create(ctx context.Context, orgID, name string, nameNe *string, description *string, descriptionNe *string, durationDays int, price string, maxMembers *int, features json.RawMessage) (*Package, error) {
	var p Package
	err := r.db.QueryRow(ctx,
		`INSERT INTO packages (organization_id, name, name_ne, description, description_ne,
		                       duration_days, price, max_members, features)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		 RETURNING id, organization_id, name, name_ne, description, description_ne,
		           duration_days, price, currency, max_members, features,
		           is_active, created_at, updated_at`,
		orgID, name, nameNe, description, descriptionNe, durationDays, price, maxMembers, features,
	).Scan(&p.ID, &p.OrganizationID, &p.Name, &p.NameNe, &p.Description, &p.DescriptionNe,
		&p.DurationDays, &p.Price, &p.Currency, &p.MaxMembers, &p.Features,
		&p.IsActive, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("creating package: %w", err)
	}
	return &p, nil
}

// Update modifies an existing package.
func (r *Repository) Update(ctx context.Context, id, orgID string, name *string, nameNe *string, description *string, descriptionNe *string, durationDays *int, price *string, maxMembers *int, features json.RawMessage, isActive *bool) (*Package, error) {
	var p Package
	err := r.db.QueryRow(ctx,
		`UPDATE packages SET
			name           = COALESCE($3, name),
			name_ne        = COALESCE($4, name_ne),
			description    = COALESCE($5, description),
			description_ne = COALESCE($6, description_ne),
			duration_days  = COALESCE($7, duration_days),
			price          = COALESCE($8, price),
			max_members    = COALESCE($9, max_members),
			features       = COALESCE($10, features),
			is_active      = COALESCE($11, is_active)
		 WHERE id = $1 AND organization_id = $2
		 RETURNING id, organization_id, name, name_ne, description, description_ne,
		           duration_days, price, currency, max_members, features,
		           is_active, created_at, updated_at`,
		id, orgID, name, nameNe, description, descriptionNe, durationDays, price, maxMembers, features, isActive,
	).Scan(&p.ID, &p.OrganizationID, &p.Name, &p.NameNe, &p.Description, &p.DescriptionNe,
		&p.DurationDays, &p.Price, &p.Currency, &p.MaxMembers, &p.Features,
		&p.IsActive, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("updating package %s: %w", id, err)
	}
	return &p, nil
}

// SoftDelete deactivates a package (sets is_active = false).
func (r *Repository) SoftDelete(ctx context.Context, id, orgID string) error {
	tag, err := r.db.Exec(ctx,
		`UPDATE packages SET is_active = false WHERE id = $1 AND organization_id = $2`,
		id, orgID,
	)
	if err != nil {
		return fmt.Errorf("soft-deleting package %s: %w", id, err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("package %s not found in organization", id)
	}
	return nil
}

// CreateMemberPackage creates a pending subscription record.
func (r *Repository) CreateMemberPackage(ctx context.Context, userID, packageID, orgID string, startDate, endDate time.Time, paymentMethod string, paymentReference *string, amountPaid string, status string) (*MemberPackage, error) {
	var mp MemberPackage
	err := r.db.QueryRow(ctx,
		`INSERT INTO member_packages (user_id, package_id, organization_id,
		                              start_date, end_date, payment_method,
		                              payment_reference, amount_paid, status)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		 RETURNING id, user_id, package_id, organization_id,
		           start_date, end_date, payment_method, payment_reference,
		           amount_paid, status, created_at`,
		userID, packageID, orgID, startDate, endDate, paymentMethod, paymentReference, amountPaid, status,
	).Scan(&mp.ID, &mp.UserID, &mp.PackageID, &mp.OrganizationID,
		&mp.StartDate, &mp.EndDate, &mp.PaymentMethod, &mp.PaymentReference,
		&mp.AmountPaid, &mp.Status, &mp.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("creating member package: %w", err)
	}
	return &mp, nil
}

// GetMemberPackageByPaymentRef finds a pending member_package by payment_reference (Khalti pidx).
func (r *Repository) GetMemberPackageByPaymentRef(ctx context.Context, paymentRef string) (*MemberPackage, error) {
	var mp MemberPackage
	err := r.db.QueryRow(ctx,
		`SELECT id, user_id, package_id, organization_id,
		        start_date, end_date, payment_method, payment_reference,
		        amount_paid, status, created_at
		 FROM member_packages
		 WHERE payment_reference = $1 AND status = 'pending'`,
		paymentRef,
	).Scan(&mp.ID, &mp.UserID, &mp.PackageID, &mp.OrganizationID,
		&mp.StartDate, &mp.EndDate, &mp.PaymentMethod, &mp.PaymentReference,
		&mp.AmountPaid, &mp.Status, &mp.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("getting member package by payment ref %s: %w", paymentRef, err)
	}
	return &mp, nil
}

// ActivateMemberPackage sets a member_package to active with the transaction reference.
func (r *Repository) ActivateMemberPackage(ctx context.Context, id, transactionID string) error {
	tag, err := r.db.Exec(ctx,
		`UPDATE member_packages
		 SET status = 'active', payment_reference = $2
		 WHERE id = $1 AND status = 'pending'`,
		id, transactionID,
	)
	if err != nil {
		return fmt.Errorf("activating member package %s: %w", id, err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("member package %s not found or not pending", id)
	}
	return nil
}

// CancelMemberPackage sets a pending member_package to cancelled.
func (r *Repository) CancelMemberPackage(ctx context.Context, id string) error {
	tag, err := r.db.Exec(ctx,
		`UPDATE member_packages SET status = 'cancelled' WHERE id = $1 AND status = 'pending'`,
		id,
	)
	if err != nil {
		return fmt.Errorf("cancelling member package %s: %w", id, err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("member package %s not found or not pending", id)
	}
	return nil
}

// ListMemberPackages returns a user's subscriptions with package and org names.
func (r *Repository) ListMemberPackages(ctx context.Context, userID string, limit, offset int) ([]MemberPackage, error) {
	rows, err := r.db.Query(ctx,
		`SELECT mp.id, mp.user_id, mp.package_id, mp.organization_id,
		        mp.start_date, mp.end_date, mp.payment_method, mp.payment_reference,
		        mp.amount_paid, mp.status, mp.created_at,
		        p.name AS package_name, o.name AS org_name
		 FROM member_packages mp
		 JOIN packages p ON p.id = mp.package_id
		 JOIN organizations o ON o.id = mp.organization_id
		 WHERE mp.user_id = $1
		 ORDER BY mp.created_at DESC
		 LIMIT $2 OFFSET $3`,
		userID, limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("querying member packages: %w", err)
	}
	defer rows.Close()

	var results []MemberPackage
	for rows.Next() {
		var mp MemberPackage
		if err := rows.Scan(
			&mp.ID, &mp.UserID, &mp.PackageID, &mp.OrganizationID,
			&mp.StartDate, &mp.EndDate, &mp.PaymentMethod, &mp.PaymentReference,
			&mp.AmountPaid, &mp.Status, &mp.CreatedAt,
			&mp.PackageName, &mp.OrgName,
		); err != nil {
			return nil, fmt.Errorf("scanning member package: %w", err)
		}
		results = append(results, mp)
	}
	return results, rows.Err()
}

// GetMemberPackageByID returns a single member_package by ID.
func (r *Repository) GetMemberPackageByID(ctx context.Context, id string) (*MemberPackage, error) {
	var mp MemberPackage
	err := r.db.QueryRow(ctx,
		`SELECT mp.id, mp.user_id, mp.package_id, mp.organization_id,
		        mp.start_date, mp.end_date, mp.payment_method, mp.payment_reference,
		        mp.amount_paid, mp.status, mp.created_at,
		        p.name AS package_name, o.name AS org_name
		 FROM member_packages mp
		 JOIN packages p ON p.id = mp.package_id
		 JOIN organizations o ON o.id = mp.organization_id
		 WHERE mp.id = $1`,
		id,
	).Scan(&mp.ID, &mp.UserID, &mp.PackageID, &mp.OrganizationID,
		&mp.StartDate, &mp.EndDate, &mp.PaymentMethod, &mp.PaymentReference,
		&mp.AmountPaid, &mp.Status, &mp.CreatedAt,
		&mp.PackageName, &mp.OrgName)
	if err != nil {
		return nil, fmt.Errorf("getting member package %s: %w", id, err)
	}
	return &mp, nil
}

// HasActiveSub checks whether a user already has an active subscription to a given package.
func (r *Repository) HasActiveSub(ctx context.Context, userID, packageID string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx,
		`SELECT EXISTS(
			SELECT 1 FROM member_packages
			WHERE user_id = $1 AND package_id = $2 AND status IN ('active', 'pending')
		)`,
		userID, packageID,
	).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("checking active subscription: %w", err)
	}
	return exists, nil
}

func scanPackages(rows pgx.Rows) ([]Package, error) {
	var pkgs []Package
	for rows.Next() {
		var p Package
		if err := rows.Scan(
			&p.ID, &p.OrganizationID, &p.Name, &p.NameNe, &p.Description, &p.DescriptionNe,
			&p.DurationDays, &p.Price, &p.Currency, &p.MaxMembers, &p.Features,
			&p.IsActive, &p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scanning package: %w", err)
		}
		pkgs = append(pkgs, p)
	}
	return pkgs, rows.Err()
}
