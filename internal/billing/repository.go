package billing

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// SaaSPlan represents a SaaS subscription tier.
type SaaSPlan struct {
	ID            string          `json:"id"`
	Name          string          `json:"name"`
	NameNe        string          `json:"name_ne,omitempty"`
	Description   string          `json:"description,omitempty"`
	DescriptionNe string          `json:"description_ne,omitempty"`
	PriceMonthly  float64         `json:"price_monthly"`
	PriceYearly   float64         `json:"price_yearly"`
	Features      json.RawMessage `json:"features"`
	MaxMembers    *int            `json:"max_members"`
	IsActive      bool            `json:"is_active"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
}

// SaaSSubscription represents an organization's subscription.
type SaaSSubscription struct {
	ID                 string    `json:"id"`
	OrganizationID     string    `json:"organization_id"`
	PlanID             string    `json:"plan_id"`
	Status             string    `json:"status"`
	CurrentPeriodStart time.Time `json:"current_period_start"`
	CurrentPeriodEnd   time.Time `json:"current_period_end"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

// Repository handles billing data persistence.
type Repository struct {
	db *pgxpool.Pool
}

// NewRepository creates a new billing repository.
func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

// ListActivePlans returns all active SaaS plans.
func (r *Repository) ListActivePlans(ctx context.Context) ([]SaaSPlan, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, name, COALESCE(name_ne, ''), COALESCE(description, ''), COALESCE(description_ne, ''),
		        price_monthly, price_yearly, features, max_members, is_active, created_at, updated_at
		 FROM saas_plans
		 WHERE is_active = true
		 ORDER BY price_monthly ASC`,
	)
	if err != nil {
		return nil, fmt.Errorf("listing plans: %w", err)
	}
	defer rows.Close()

	var plans []SaaSPlan
	for rows.Next() {
		var p SaaSPlan
		if err := rows.Scan(
			&p.ID, &p.Name, &p.NameNe, &p.Description, &p.DescriptionNe,
			&p.PriceMonthly, &p.PriceYearly, &p.Features, &p.MaxMembers,
			&p.IsActive, &p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scanning plan: %w", err)
		}
		plans = append(plans, p)
	}
	return plans, rows.Err()
}

// GetPlanByID retrieves a SaaS plan by ID.
func (r *Repository) GetPlanByID(ctx context.Context, id string) (*SaaSPlan, error) {
	var p SaaSPlan
	err := r.db.QueryRow(ctx,
		`SELECT id, name, COALESCE(name_ne, ''), COALESCE(description, ''), COALESCE(description_ne, ''),
		        price_monthly, price_yearly, features, max_members, is_active, created_at, updated_at
		 FROM saas_plans WHERE id = $1`,
		id,
	).Scan(
		&p.ID, &p.Name, &p.NameNe, &p.Description, &p.DescriptionNe,
		&p.PriceMonthly, &p.PriceYearly, &p.Features, &p.MaxMembers,
		&p.IsActive, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("plan not found: %w", err)
	}
	return &p, nil
}

// GetActiveSubscription returns the active subscription for an organization, if any.
func (r *Repository) GetActiveSubscription(ctx context.Context, orgID string) (*SaaSSubscription, error) {
	var s SaaSSubscription
	err := r.db.QueryRow(ctx,
		`SELECT id, organization_id, plan_id, status,
		        current_period_start, current_period_end, created_at, updated_at
		 FROM saas_subscriptions
		 WHERE organization_id = $1 AND status = 'active'
		 ORDER BY created_at DESC LIMIT 1`,
		orgID,
	).Scan(
		&s.ID, &s.OrganizationID, &s.PlanID, &s.Status,
		&s.CurrentPeriodStart, &s.CurrentPeriodEnd, &s.CreatedAt, &s.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("no active subscription: %w", err)
	}
	return &s, nil
}

// CreateSubscription inserts a new SaaS subscription.
func (r *Repository) CreateSubscription(ctx context.Context, orgID, planID, status string, periodStart, periodEnd time.Time) (*SaaSSubscription, error) {
	var s SaaSSubscription
	err := r.db.QueryRow(ctx,
		`INSERT INTO saas_subscriptions
			(organization_id, plan_id, status, current_period_start, current_period_end)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id, organization_id, plan_id, status,
		           current_period_start, current_period_end, created_at, updated_at`,
		orgID, planID, status, periodStart, periodEnd,
	).Scan(
		&s.ID, &s.OrganizationID, &s.PlanID, &s.Status,
		&s.CurrentPeriodStart, &s.CurrentPeriodEnd, &s.CreatedAt, &s.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("creating subscription: %w", err)
	}
	return &s, nil
}

// IsOrgAdmin checks if a user is an admin of the given organization.
func (r *Repository) IsOrgAdmin(ctx context.Context, userID, orgID string) (bool, error) {
	var role string
	err := r.db.QueryRow(ctx,
		`SELECT role FROM organization_members
		 WHERE user_id = $1 AND organization_id = $2 AND status = 'active'`,
		userID, orgID,
	).Scan(&role)
	if err != nil {
		return false, nil // not a member
	}
	return role == "admin", nil
}

// IsSuperAdmin checks if a user has the super_admin flag.
func (r *Repository) IsSuperAdmin(ctx context.Context, userID string) (bool, error) {
	var isSA bool
	err := r.db.QueryRow(ctx,
		`SELECT is_super_admin FROM users WHERE id = $1`, userID,
	).Scan(&isSA)
	if err != nil {
		return false, fmt.Errorf("checking super admin: %w", err)
	}
	return isSA, nil
}
