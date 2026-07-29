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
	PANNumber     string          `json:"pan_number,omitempty"`
	TaxLegalName  string          `json:"tax_legal_name,omitempty"`
	TaxAddress    string          `json:"tax_address,omitempty"`
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
	PANNumber     *string
	TaxLegalName  *string
	TaxAddress    *string
}

// orgPublicColumns backs the anonymous gym directory (GetByID/List, mounted
// by RegisterPublicRoutes with no auth). It must never include pan_number,
// tax_legal_name or tax_address — those are tax identifiers an org admin
// manages via the authenticated PUT /api/v1/orgs/{orgId}, not public data.
const orgPublicColumns = `id, name, COALESCE(name_ne, ''), slug,
	COALESCE(description, ''), COALESCE(description_ne, ''),
	COALESCE(logo_url, ''), COALESCE(address, ''), COALESCE(address_ne, ''),
	COALESCE(phone, ''), COALESCE(email, ''),
	latitude, longitude,
	timezone, settings, is_active, created_at, updated_at`

// orgSelectColumns additionally includes the gym's tax identity. Use only for
// authenticated, admin-facing reads (Create/Update responses) — never wire
// this into a route reachable without an org-admin session, per
// orgPublicColumns above.
const orgSelectColumns = `id, name, COALESCE(name_ne, ''), slug,
	COALESCE(description, ''), COALESCE(description_ne, ''),
	COALESCE(logo_url, ''), COALESCE(address, ''), COALESCE(address_ne, ''),
	COALESCE(phone, ''), COALESCE(email, ''),
	latitude, longitude,
	COALESCE(pan_number, ''), COALESCE(tax_legal_name, ''), COALESCE(tax_address, ''),
	timezone, settings, is_active, created_at, updated_at`

// Repository handles organization data persistence.
type Repository struct {
	db *pgxpool.Pool
}

// NewRepository creates a new org repository.
func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

// scanOrgPublic scans the orgPublicColumns projection: no tax identity
// fields, so PANNumber/TaxLegalName/TaxAddress stay at their zero value ("")
// and are omitted from the JSON response by their omitempty tags.
func scanOrgPublic(scanner interface{ Scan(dest ...any) error }) (*Organization, error) {
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

func scanOrg(scanner interface{ Scan(dest ...any) error }) (*Organization, error) {
	var o Organization
	err := scanner.Scan(
		&o.ID, &o.Name, &o.NameNe, &o.Slug,
		&o.Description, &o.DescriptionNe,
		&o.LogoURL, &o.Address, &o.AddressNe,
		&o.Phone, &o.Email,
		&o.Latitude, &o.Longitude,
		&o.PANNumber, &o.TaxLegalName, &o.TaxAddress,
		&o.Timezone, &o.Settings, &o.IsActive, &o.CreatedAt, &o.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &o, nil
}

// GetByID retrieves an organization by ID for the public gym directory.
// Deliberately uses orgPublicColumns, not orgSelectColumns: this backs an
// unauthenticated route, and the org's tax identity must not leak there.
func (r *Repository) GetByID(ctx context.Context, id string) (*Organization, error) {
	row := r.db.QueryRow(ctx,
		`SELECT `+orgPublicColumns+` FROM organizations WHERE id = $1`, id,
	)
	o, err := scanOrgPublic(row)
	if err != nil {
		return nil, fmt.Errorf("organization not found: %w", err)
	}
	return o, nil
}

// List retrieves active organizations (public listing). Deliberately uses
// orgPublicColumns — see GetByID.
func (r *Repository) List(ctx context.Context, limit, offset int) ([]Organization, error) {
	rows, err := r.db.Query(ctx,
		`SELECT `+orgPublicColumns+` FROM organizations WHERE is_active = true ORDER BY name LIMIT $1 OFFSET $2`,
		limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("listing organizations: %w", err)
	}
	defer rows.Close()

	var orgs []Organization
	for rows.Next() {
		o, err := scanOrgPublic(rows)
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
			settings      = COALESCE($13, settings),
			-- pan_number/tax_legal_name/tax_address: a field absent from the
			-- request (SQL NULL param) must leave the column unchanged, while a
			-- field explicitly sent as "" must clear it to SQL NULL. Plain
			-- COALESCE($n, column) can't distinguish those, since a missing
			-- request field and "please blank this out" are both carried as
			-- NULL at the parameter level. Resolving COALESCE first (nil param
			-- -> fall back to the existing value; non-nil param, including "" ->
			-- use it) and then NULLIF-ing that *result* against '' turns any
			-- surviving '' — which can only be the "" the caller explicitly
			-- sent, never a stored value, because this statement is the only
			-- writer of these three columns and every write already passes
			-- through this same NULLIF — into SQL NULL.
			pan_number     = NULLIF(COALESCE($14, pan_number), ''),
			tax_legal_name = NULLIF(COALESCE($15, tax_legal_name), ''),
			tax_address    = NULLIF(COALESCE($16, tax_address), '')
		 WHERE id = $1
		 RETURNING `+orgSelectColumns,
		id,
		in.Name, in.NameNe,
		in.Description, in.DescriptionNe,
		in.LogoURL, in.Address, in.AddressNe,
		in.Phone, in.Email,
		in.Latitude, in.Longitude,
		in.Settings,
		in.PANNumber, in.TaxLegalName, in.TaxAddress,
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

// GetSlugByID retrieves an organization's slug by its ID.
func (r *Repository) GetSlugByID(ctx context.Context, orgID string) (string, error) {
	var slug string
	err := r.db.QueryRow(ctx,
		`SELECT slug FROM organizations WHERE id = $1`, orgID,
	).Scan(&slug)
	if err != nil {
		return "", fmt.Errorf("getting org slug: %w", err)
	}
	return slug, nil
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
