package org

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Organization represents a gym/fitness center.
type Organization struct {
	ID            string          `json:"id"`
	Name          string          `json:"name"`
	NameNe        string          `json:"name_ne,omitempty"`
	Slug          string          `json:"slug"`
	Description   string          `json:"description,omitempty"`
	DescriptionNe string          `json:"description_ne,omitempty"`
	LogoURL       string          `json:"logo_url,omitempty"`
	Address       string          `json:"address,omitempty"`
	AddressNe     string          `json:"address_ne,omitempty"`
	Phone         string          `json:"phone,omitempty"`
	Email         string          `json:"email,omitempty"`
	Latitude      *float64        `json:"latitude,omitempty"`
	Longitude     *float64        `json:"longitude,omitempty"`
	Timezone      string          `json:"timezone"`
	Settings      json.RawMessage `json:"settings"`
	IsActive      bool            `json:"is_active"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
}

// CreateOrgInput holds input data for creating an organization.
type CreateOrgInput struct {
	Name          string
	NameNe        string
	Slug          string
	Description   string
	DescriptionNe string
	LogoURL       string
	Address       string
	AddressNe     string
	Phone         string
	Email         string
	Latitude      *float64
	Longitude     *float64
	Settings      json.RawMessage
}

// UpdateOrgInput holds input data for updating an organization.
// Pointer fields: nil means "don't change", non-nil means "set to this value".
type UpdateOrgInput struct {
	Name          *string
	NameNe        *string
	Description   *string
	DescriptionNe *string
	LogoURL       *string
	Address       *string
	AddressNe     *string
	Phone         *string
	Email         *string
	Latitude      *float64
	Longitude     *float64
	Settings      *json.RawMessage
}

const orgSelectColumns = `id, name, COALESCE(name_ne, ''), slug,
	COALESCE(description, ''), COALESCE(description_ne, ''),
	COALESCE(logo_url, ''), COALESCE(address, ''), COALESCE(address_ne, ''),
	COALESCE(phone, ''), COALESCE(email, ''),
	latitude, longitude,
	timezone, settings, is_active, created_at, updated_at`

// Repository handles organization data persistence.
type Repository struct {
	db *pgxpool.Pool
}

// NewRepository creates a new org repository.
func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func scanOrg(scanner interface{ Scan(dest ...any) error }) (*Organization, error) {
	var o Organization
	err := scanner.Scan(
		&o.ID, &o.Name, &o.NameNe, &o.Slug,
		&o.Description, &o.DescriptionNe,
		&o.LogoURL, &o.Address, &o.AddressNe,
		&o.Phone, &o.Email,
		&o.Latitude, &o.Longitude,
		&o.Timezone, &o.Settings, &o.IsActive, &o.CreatedAt, &o.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &o, nil
}

// GetByID retrieves an organization by ID.
func (r *Repository) GetByID(ctx context.Context, id string) (*Organization, error) {
	row := r.db.QueryRow(ctx,
		`SELECT `+orgSelectColumns+` FROM organizations WHERE id = $1`, id,
	)
	o, err := scanOrg(row)
	if err != nil {
		return nil, fmt.Errorf("organization not found: %w", err)
	}
	return o, nil
}

// List retrieves active organizations (public listing).
func (r *Repository) List(ctx context.Context, limit, offset int) ([]Organization, error) {
	rows, err := r.db.Query(ctx,
		`SELECT `+orgSelectColumns+` FROM organizations WHERE is_active = true ORDER BY name LIMIT $1 OFFSET $2`,
		limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("listing organizations: %w", err)
	}
	defer rows.Close()

	var orgs []Organization
	for rows.Next() {
		o, err := scanOrg(rows)
		if err != nil {
			return nil, fmt.Errorf("scanning organization: %w", err)
		}
		orgs = append(orgs, *o)
	}
	return orgs, rows.Err()
}

// Create inserts a new organization and returns it.
func (r *Repository) Create(ctx context.Context, in *CreateOrgInput) (*Organization, error) {
	row := r.db.QueryRow(ctx,
		`INSERT INTO organizations
			(name, name_ne, slug, description, description_ne,
			 logo_url, address, address_ne, phone, email,
			 latitude, longitude, settings)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)
		 RETURNING `+orgSelectColumns,
		in.Name, nullIfEmpty(in.NameNe), in.Slug,
		nullIfEmpty(in.Description), nullIfEmpty(in.DescriptionNe),
		nullIfEmpty(in.LogoURL), nullIfEmpty(in.Address), nullIfEmpty(in.AddressNe),
		nullIfEmpty(in.Phone), nullIfEmpty(in.Email),
		in.Latitude, in.Longitude,
		defaultJSON(in.Settings),
	)
	o, err := scanOrg(row)
	if err != nil {
		return nil, fmt.Errorf("creating organization: %w", err)
	}
	return o, nil
}

// Update modifies an existing organization using COALESCE for partial updates.
func (r *Repository) Update(ctx context.Context, id string, in *UpdateOrgInput) (*Organization, error) {
	row := r.db.QueryRow(ctx,
		`UPDATE organizations SET
			name          = COALESCE($2, name),
			name_ne       = COALESCE($3, name_ne),
			description   = COALESCE($4, description),
			description_ne = COALESCE($5, description_ne),
			logo_url      = COALESCE($6, logo_url),
			address       = COALESCE($7, address),
			address_ne    = COALESCE($8, address_ne),
			phone         = COALESCE($9, phone),
			email         = COALESCE($10, email),
			latitude      = COALESCE($11, latitude),
			longitude     = COALESCE($12, longitude),
			settings      = COALESCE($13, settings)
		 WHERE id = $1
		 RETURNING `+orgSelectColumns,
		id,
		in.Name, in.NameNe,
		in.Description, in.DescriptionNe,
		in.LogoURL, in.Address, in.AddressNe,
		in.Phone, in.Email,
		in.Latitude, in.Longitude,
		in.Settings,
	)
	o, err := scanOrg(row)
	if err != nil {
		return nil, fmt.Errorf("updating organization: %w", err)
	}
	return o, nil
}

// CheckOrgMembership checks if a user is a member of an organization and returns their role.
// This implements the middleware.OrgMembershipChecker interface.
func (r *Repository) CheckOrgMembership(ctx context.Context, userID, orgID string) (bool, string, error) {
	var role string
	err := r.db.QueryRow(ctx,
		`SELECT role FROM organization_members
		 WHERE user_id = $1 AND organization_id = $2 AND status = 'active'`,
		userID, orgID,
	).Scan(&role)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, "", nil
		}
		return false, "", err
	}
	return true, role, nil
}

// GetOrgRole returns the user's role within an organization.
// Returns empty string if user is not a member.
func (r *Repository) GetOrgRole(ctx context.Context, userID, orgID string) (string, error) {
	var role string
	err := r.db.QueryRow(ctx,
		`SELECT role FROM organization_members
		 WHERE user_id = $1 AND organization_id = $2 AND status = 'active'`,
		userID, orgID,
	).Scan(&role)
	if err != nil {
		return "", fmt.Errorf("querying org role: %w", err)
	}
	return role, nil
}

// SlugExists checks if a slug is already taken.
func (r *Repository) SlugExists(ctx context.Context, slug string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM organizations WHERE slug = $1)`, slug,
	).Scan(&exists)
	return exists, err
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

func nullIfEmpty(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func defaultJSON(raw json.RawMessage) json.RawMessage {
	if len(raw) == 0 {
		return json.RawMessage(`{}`)
	}
	return raw
}
