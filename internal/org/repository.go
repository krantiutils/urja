package org

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Organization represents a gym/fitness center.
type Organization struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	Address   string    `json:"address,omitempty"`
	Phone     string    `json:"phone,omitempty"`
	LogoURL   string    `json:"logo_url,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// Repository handles organization data persistence.
type Repository struct {
	db *pgxpool.Pool
}

// NewRepository creates a new org repository.
func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

// GetByID retrieves an organization by ID.
func (r *Repository) GetByID(ctx context.Context, id string) (*Organization, error) {
	var o Organization
	err := r.db.QueryRow(ctx,
		`SELECT id, name, slug, COALESCE(address, ''), COALESCE(phone, ''),
		        COALESCE(logo_url, ''), created_at
		 FROM organizations WHERE id = $1`,
		id,
	).Scan(&o.ID, &o.Name, &o.Slug, &o.Address, &o.Phone, &o.LogoURL, &o.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("organization not found: %w", err)
	}
	return &o, nil
}

// List retrieves all organizations (public listing).
func (r *Repository) List(ctx context.Context, limit, offset int) ([]Organization, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, name, slug, COALESCE(address, ''), COALESCE(phone, ''),
		        COALESCE(logo_url, ''), created_at
		 FROM organizations ORDER BY name LIMIT $1 OFFSET $2`,
		limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("listing organizations: %w", err)
	}
	defer rows.Close()

	var orgs []Organization
	for rows.Next() {
		var o Organization
		if err := rows.Scan(&o.ID, &o.Name, &o.Slug, &o.Address, &o.Phone, &o.LogoURL, &o.CreatedAt); err != nil {
			return nil, fmt.Errorf("scanning organization: %w", err)
		}
		orgs = append(orgs, o)
	}
	return orgs, rows.Err()
}

// IsOrgMember checks if a user is a member of an organization.
// This implements the middleware.OrgMembershipChecker interface.
func (r *Repository) IsOrgMember(ctx context.Context, userID, orgID string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx,
		`SELECT EXISTS(
			SELECT 1 FROM organization_members
			WHERE user_id = $1 AND organization_id = $2 AND status = 'active'
		)`,
		userID, orgID,
	).Scan(&exists)
	return exists, err
}
