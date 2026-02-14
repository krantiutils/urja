package subscription

import (
	"time"
)

// Package represents a gym package offering.
type Package struct {
	ID                  string  `json:"id"`
	OrganizationID      string  `json:"organization_id"`
	Name                string  `json:"name"`
	NameNe              string  `json:"name_ne,omitempty"`
	Description         string  `json:"description,omitempty"`
	DescriptionNe       string  `json:"description_ne,omitempty"`
	DurationDays        int     `json:"duration_days"`
	Price               float64 `json:"price"`
	Currency            string  `json:"currency"`
	MaxMembers          *int    `json:"max_members,omitempty"`
	IsActive            bool    `json:"is_active"`
	RegistrationEnabled bool    `json:"registration_enabled"`
}

// PackageSummary extends Package with active member count.
type PackageSummary struct {
	Package
	MemberCount int `json:"member_count"`
}

// ExpiringEntry represents a member whose package is about to expire.
type ExpiringEntry struct {
	MemberID    string `json:"member_id"`
	MemberName  string `json:"member_name"`
	MemberPhone string `json:"member_phone"`
	PackageName string `json:"package_name"`
	ExpiryDate  string `json:"expiry_date"`  // YYYY-MM-DD
	ExpiresIn   int    `json:"expires_in"`    // days until expiry
	SubID       string `json:"subscription_id"`
}

// ExpiredEntry represents a member whose package has already expired.
type ExpiredEntry struct {
	MemberID    string `json:"member_id"`
	MemberName  string `json:"member_name"`
	MemberPhone string `json:"member_phone"`
	PackageName string `json:"package_name"`
	ExpiryDate  string `json:"expiry_date"`   // YYYY-MM-DD
	ExpiredAgo  int    `json:"expired_ago"`    // days since expiry
	SubID       string `json:"subscription_id"`
}

// MemberSubscription represents a member's package subscription record.
type MemberSubscription struct {
	ID            string    `json:"id"`
	PackageID     string    `json:"package_id"`
	PackageName   string    `json:"package_name"`
	StartDate     string    `json:"start_date"`
	EndDate       string    `json:"end_date"`
	Status        string    `json:"status"`
	DaysRemaining int       `json:"days_remaining"`
	AmountPaid    float64   `json:"amount_paid"`
	PaymentMethod string    `json:"payment_method,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
}

// Payment represents a payment history record.
type Payment struct {
	ID               string    `json:"id"`
	MemberPackageID  string    `json:"member_package_id"`
	Particular       string    `json:"particular"`
	TotalAmount      float64   `json:"total_amount"`
	Discount         float64   `json:"discount"`
	PaidAmount       float64   `json:"paid_amount"`
	Due              float64   `json:"due"`
	PaymentMethod    string    `json:"payment_method,omitempty"`
	PaymentReference string    `json:"payment_reference,omitempty"`
	PaymentDate      time.Time `json:"payment_date"`
}

// AssignRequest is the input for assigning a package to a member.
type AssignRequest struct {
	PackageID        string  `json:"package_id"`
	StartDate        string  `json:"start_date"` // YYYY-MM-DD
	PaymentMethod    string  `json:"payment_method"`
	AmountPaid       float64 `json:"amount_paid"`
	Discount         float64 `json:"discount"`
	PaymentReference string  `json:"payment_reference,omitempty"`
}

// RenewRequest is the input for renewing a member's package.
type RenewRequest struct {
	MemberPackageID  string  `json:"member_package_id"`
	PaymentMethod    string  `json:"payment_method"`
	AmountPaid       float64 `json:"amount_paid"`
	Discount         float64 `json:"discount"`
	PaymentReference string  `json:"payment_reference,omitempty"`
}

// ExtendRequest is the input for extending a member's package.
type ExtendRequest struct {
	MemberPackageID  string  `json:"member_package_id"`
	ExtraDays        int     `json:"extra_days"`
	PaymentMethod    string  `json:"payment_method"`
	AmountPaid       float64 `json:"amount_paid"`
	Discount         float64 `json:"discount"`
	PaymentReference string  `json:"payment_reference,omitempty"`
}
