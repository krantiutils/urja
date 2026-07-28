package subscription

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository handles subscription data persistence.
type Repository struct {
	db *pgxpool.Pool
}

// NewRepository creates a new subscription repository.
func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

// GetPackage retrieves a single package by ID within an organization.
func (r *Repository) GetPackage(ctx context.Context, orgID, packageID string) (*Package, error) {
	var p Package
	var maxMembers *int
	err := r.db.QueryRow(ctx,
		`SELECT id, organization_id, name, COALESCE(name_ne, ''),
		        COALESCE(description, ''), COALESCE(description_ne, ''),
		        duration_days, price, currency, max_members,
		        is_active, registration_enabled
		 FROM packages
		 WHERE id = $1 AND organization_id = $2`,
		packageID, orgID,
	).Scan(&p.ID, &p.OrganizationID, &p.Name, &p.NameNe,
		&p.Description, &p.DescriptionNe,
		&p.DurationDays, &p.Price, &p.Currency, &maxMembers,
		&p.IsActive, &p.RegistrationEnabled)
	if err != nil {
		return nil, fmt.Errorf("package not found: %w", err)
	}
	p.MaxMembers = maxMembers
	return &p, nil
}

// PackageSummaryList retrieves all packages for an org with active member counts.
func (r *Repository) PackageSummaryList(ctx context.Context, orgID string) ([]PackageSummary, error) {
	rows, err := r.db.Query(ctx,
		`SELECT p.id, p.organization_id, p.name, COALESCE(p.name_ne, ''),
		        COALESCE(p.description, ''), COALESCE(p.description_ne, ''),
		        p.duration_days, p.price, p.currency, p.max_members,
		        p.is_active, p.registration_enabled,
		        COUNT(mp.id) FILTER (WHERE mp.status = 'active' AND mp.end_date >= CURRENT_DATE) AS member_count
		 FROM packages p
		 LEFT JOIN member_packages mp ON p.id = mp.package_id
		 WHERE p.organization_id = $1
		 GROUP BY p.id
		 ORDER BY p.created_at DESC`,
		orgID,
	)
	if err != nil {
		return nil, fmt.Errorf("querying package summary: %w", err)
	}
	defer rows.Close()

	var summaries []PackageSummary
	for rows.Next() {
		var s PackageSummary
		var maxMembers *int
		if err := rows.Scan(
			&s.ID, &s.OrganizationID, &s.Name, &s.NameNe,
			&s.Description, &s.DescriptionNe,
			&s.DurationDays, &s.Price, &s.Currency, &maxMembers,
			&s.IsActive, &s.RegistrationEnabled,
			&s.MemberCount,
		); err != nil {
			return nil, fmt.Errorf("scanning package summary: %w", err)
		}
		s.MaxMembers = maxMembers
		summaries = append(summaries, s)
	}
	return summaries, rows.Err()
}

// ListExpiring retrieves member packages expiring within the given number of days.
func (r *Repository) ListExpiring(ctx context.Context, orgID string, withinDays, limit, offset int) ([]ExpiringEntry, int, error) {
	var total int
	err := r.db.QueryRow(ctx,
		`SELECT COUNT(*)
		 FROM member_packages mp
		 WHERE mp.organization_id = $1
		   AND mp.status = 'active'
		   AND mp.end_date >= CURRENT_DATE
		   AND mp.end_date <= CURRENT_DATE + $2 * INTERVAL '1 day'`,
		orgID, withinDays,
	).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("counting expiring packages: %w", err)
	}

	rows, err := r.db.Query(ctx,
		`SELECT mp.id, u.id, COALESCE(u.name, ''), u.phone, p.name,
		        mp.end_date, (mp.end_date - CURRENT_DATE)::int AS expires_in
		 FROM member_packages mp
		 JOIN users u ON mp.user_id = u.id
		 JOIN packages p ON mp.package_id = p.id
		 WHERE mp.organization_id = $1
		   AND mp.status = 'active'
		   AND mp.end_date >= CURRENT_DATE
		   AND mp.end_date <= CURRENT_DATE + $2 * INTERVAL '1 day'
		 ORDER BY mp.end_date ASC
		 LIMIT $3 OFFSET $4`,
		orgID, withinDays, limit, offset,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("querying expiring packages: %w", err)
	}
	defer rows.Close()

	var entries []ExpiringEntry
	for rows.Next() {
		var e ExpiringEntry
		var endDate time.Time
		if err := rows.Scan(&e.SubID, &e.MemberID, &e.MemberName, &e.MemberPhone,
			&e.PackageName, &endDate, &e.ExpiresIn); err != nil {
			return nil, 0, fmt.Errorf("scanning expiring entry: %w", err)
		}
		e.ExpiryDate = endDate.Format("2006-01-02")
		entries = append(entries, e)
	}
	return entries, total, rows.Err()
}

// ListExpired retrieves member packages that have already expired.
func (r *Repository) ListExpired(ctx context.Context, orgID string, limit, offset int) ([]ExpiredEntry, int, error) {
	var total int
	err := r.db.QueryRow(ctx,
		`SELECT COUNT(*)
		 FROM member_packages mp
		 WHERE mp.organization_id = $1
		   AND (mp.end_date < CURRENT_DATE AND mp.status IN ('active', 'expired'))`,
		orgID,
	).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("counting expired packages: %w", err)
	}

	rows, err := r.db.Query(ctx,
		`SELECT mp.id, u.id, COALESCE(u.name, ''), u.phone, p.name,
		        mp.end_date, (CURRENT_DATE - mp.end_date)::int AS expired_ago
		 FROM member_packages mp
		 JOIN users u ON mp.user_id = u.id
		 JOIN packages p ON mp.package_id = p.id
		 WHERE mp.organization_id = $1
		   AND (mp.end_date < CURRENT_DATE AND mp.status IN ('active', 'expired'))
		 ORDER BY mp.end_date DESC
		 LIMIT $2 OFFSET $3`,
		orgID, limit, offset,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("querying expired packages: %w", err)
	}
	defer rows.Close()

	var entries []ExpiredEntry
	for rows.Next() {
		var e ExpiredEntry
		var endDate time.Time
		if err := rows.Scan(&e.SubID, &e.MemberID, &e.MemberName, &e.MemberPhone,
			&e.PackageName, &endDate, &e.ExpiredAgo); err != nil {
			return nil, 0, fmt.Errorf("scanning expired entry: %w", err)
		}
		e.ExpiryDate = endDate.Format("2006-01-02")
		entries = append(entries, e)
	}
	return entries, total, rows.Err()
}

// memberPackageRow holds the DB result for a member_packages insert/update.
type memberPackageRow struct {
	ID            string
	PackageID     string
	StartDate     time.Time
	EndDate       time.Time
	Status        string
	AmountPaid    float64
	PaymentMethod string
	CreatedAt     time.Time
}

// AssignPackage creates a new member_package and a corresponding payment record.
// Runs within a transaction.
func (r *Repository) AssignPackage(ctx context.Context, orgID, memberID string, req AssignRequest, pkg *Package, entryByUserID string) (*MemberSubscription, *Payment, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid start_date format: %w", err)
	}
	endDate := startDate.AddDate(0, 0, pkg.DurationDays)

	totalAmount := pkg.Price
	due := totalAmount - req.Discount - req.AmountPaid
	if due < 0 {
		due = 0
	}

	// Insert member_package
	var mp memberPackageRow
	err = tx.QueryRow(ctx,
		`INSERT INTO member_packages (user_id, package_id, organization_id, start_date, end_date, payment_method, payment_reference, amount_paid, status)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, 'active')
		 RETURNING id, package_id, start_date, end_date, status, amount_paid, COALESCE(payment_method, ''), created_at`,
		memberID, req.PackageID, orgID, startDate, endDate,
		nilIfEmpty(req.PaymentMethod), nilIfEmpty(req.PaymentReference), req.AmountPaid,
	).Scan(&mp.ID, &mp.PackageID, &mp.StartDate, &mp.EndDate, &mp.Status, &mp.AmountPaid, &mp.PaymentMethod, &mp.CreatedAt)
	if err != nil {
		return nil, nil, fmt.Errorf("inserting member_package: %w", err)
	}

	// Insert payment record
	particular := fmt.Sprintf("Package: %s (%d days)", pkg.Name, pkg.DurationDays)
	var pmt Payment
	err = tx.QueryRow(ctx,
		`INSERT INTO payments (member_package_id, user_id, organization_id, particular, total_amount, discount, paid_amount, due, payment_method, payment_reference)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		 RETURNING id, member_package_id, particular, total_amount, discount, paid_amount, due, COALESCE(payment_method, ''), COALESCE(payment_reference, ''), payment_date`,
		mp.ID, memberID, orgID, particular, totalAmount, req.Discount, req.AmountPaid, due,
		nilIfEmpty(req.PaymentMethod), nilIfEmpty(req.PaymentReference),
	).Scan(&pmt.ID, &pmt.MemberPackageID, &pmt.Particular, &pmt.TotalAmount, &pmt.Discount, &pmt.PaidAmount, &pmt.Due, &pmt.PaymentMethod, &pmt.PaymentReference, &pmt.PaymentDate)
	if err != nil {
		return nil, nil, fmt.Errorf("inserting payment: %w", err)
	}

	if err := recordPackageIncome(ctx, tx, orgID, particular, req.AmountPaid,
		req.PaymentMethod, req.PaymentReference, entryByUserID); err != nil {
		return nil, nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, nil, fmt.Errorf("commit transaction: %w", err)
	}

	daysRemaining := int(time.Until(endDate).Hours() / 24)
	if daysRemaining < 0 {
		daysRemaining = 0
	}

	sub := &MemberSubscription{
		ID:            mp.ID,
		PackageID:     mp.PackageID,
		PackageName:   pkg.Name,
		StartDate:     mp.StartDate.Format("2006-01-02"),
		EndDate:       mp.EndDate.Format("2006-01-02"),
		Status:        mp.Status,
		DaysRemaining: daysRemaining,
		AmountPaid:    mp.AmountPaid,
		PaymentMethod: mp.PaymentMethod,
		CreatedAt:     mp.CreatedAt,
	}

	return sub, &pmt, nil
}

// GetMemberPackage retrieves a single member_package record.
func (r *Repository) GetMemberPackage(ctx context.Context, orgID, memberPackageID string) (*memberPackageRow, *Package, error) {
	var mp memberPackageRow
	var pkg Package
	var maxMembers *int

	err := r.db.QueryRow(ctx,
		`SELECT mp.id, mp.package_id, mp.start_date, mp.end_date, mp.status, mp.amount_paid,
		        COALESCE(mp.payment_method, ''), mp.created_at,
		        p.id, p.organization_id, p.name, COALESCE(p.name_ne, ''),
		        COALESCE(p.description, ''), COALESCE(p.description_ne, ''),
		        p.duration_days, p.price, p.currency, p.max_members,
		        p.is_active, p.registration_enabled
		 FROM member_packages mp
		 JOIN packages p ON mp.package_id = p.id
		 WHERE mp.id = $1 AND mp.organization_id = $2`,
		memberPackageID, orgID,
	).Scan(&mp.ID, &mp.PackageID, &mp.StartDate, &mp.EndDate, &mp.Status, &mp.AmountPaid,
		&mp.PaymentMethod, &mp.CreatedAt,
		&pkg.ID, &pkg.OrganizationID, &pkg.Name, &pkg.NameNe,
		&pkg.Description, &pkg.DescriptionNe,
		&pkg.DurationDays, &pkg.Price, &pkg.Currency, &maxMembers,
		&pkg.IsActive, &pkg.RegistrationEnabled)
	if err != nil {
		return nil, nil, fmt.Errorf("member package not found: %w", err)
	}
	pkg.MaxMembers = maxMembers
	return &mp, &pkg, nil
}

// RenewPackage expires the old subscription and creates a new one with a payment.
func (r *Repository) RenewPackage(ctx context.Context, orgID, memberID string, oldMP *memberPackageRow, pkg *Package, req RenewRequest, entryByUserID string) (*MemberSubscription, *Payment, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Mark old subscription as expired
	_, err = tx.Exec(ctx,
		`UPDATE member_packages SET status = 'expired' WHERE id = $1`,
		oldMP.ID,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("expiring old subscription: %w", err)
	}

	// Create new subscription starting today
	startDate := time.Now().Truncate(24 * time.Hour)
	endDate := startDate.AddDate(0, 0, pkg.DurationDays)
	totalAmount := pkg.Price
	due := totalAmount - req.Discount - req.AmountPaid
	if due < 0 {
		due = 0
	}

	var mp memberPackageRow
	err = tx.QueryRow(ctx,
		`INSERT INTO member_packages (user_id, package_id, organization_id, start_date, end_date, payment_method, payment_reference, amount_paid, status)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, 'active')
		 RETURNING id, package_id, start_date, end_date, status, amount_paid, COALESCE(payment_method, ''), created_at`,
		memberID, pkg.ID, orgID, startDate, endDate,
		nilIfEmpty(req.PaymentMethod), nilIfEmpty(req.PaymentReference), req.AmountPaid,
	).Scan(&mp.ID, &mp.PackageID, &mp.StartDate, &mp.EndDate, &mp.Status, &mp.AmountPaid, &mp.PaymentMethod, &mp.CreatedAt)
	if err != nil {
		return nil, nil, fmt.Errorf("inserting renewed subscription: %w", err)
	}

	// Insert payment
	particular := fmt.Sprintf("Renewal: %s (%d days)", pkg.Name, pkg.DurationDays)
	var pmt Payment
	err = tx.QueryRow(ctx,
		`INSERT INTO payments (member_package_id, user_id, organization_id, particular, total_amount, discount, paid_amount, due, payment_method, payment_reference)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		 RETURNING id, member_package_id, particular, total_amount, discount, paid_amount, due, COALESCE(payment_method, ''), COALESCE(payment_reference, ''), payment_date`,
		mp.ID, memberID, orgID, particular, totalAmount, req.Discount, req.AmountPaid, due,
		nilIfEmpty(req.PaymentMethod), nilIfEmpty(req.PaymentReference),
	).Scan(&pmt.ID, &pmt.MemberPackageID, &pmt.Particular, &pmt.TotalAmount, &pmt.Discount, &pmt.PaidAmount, &pmt.Due, &pmt.PaymentMethod, &pmt.PaymentReference, &pmt.PaymentDate)
	if err != nil {
		return nil, nil, fmt.Errorf("inserting renewal payment: %w", err)
	}

	if err := recordPackageIncome(ctx, tx, orgID, particular, req.AmountPaid,
		req.PaymentMethod, req.PaymentReference, entryByUserID); err != nil {
		return nil, nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, nil, fmt.Errorf("commit transaction: %w", err)
	}

	daysRemaining := int(time.Until(endDate).Hours() / 24)
	if daysRemaining < 0 {
		daysRemaining = 0
	}

	sub := &MemberSubscription{
		ID:            mp.ID,
		PackageID:     mp.PackageID,
		PackageName:   pkg.Name,
		StartDate:     mp.StartDate.Format("2006-01-02"),
		EndDate:       mp.EndDate.Format("2006-01-02"),
		Status:        mp.Status,
		DaysRemaining: daysRemaining,
		AmountPaid:    mp.AmountPaid,
		PaymentMethod: mp.PaymentMethod,
		CreatedAt:     mp.CreatedAt,
	}

	return sub, &pmt, nil
}

// ExtendPackage adds extra days to an existing subscription and records a payment.
func (r *Repository) ExtendPackage(ctx context.Context, orgID, memberID string, oldMP *memberPackageRow, pkg *Package, req ExtendRequest, entryByUserID string) (*MemberSubscription, *Payment, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	newEndDate := oldMP.EndDate.AddDate(0, 0, req.ExtraDays)

	_, err = tx.Exec(ctx,
		`UPDATE member_packages SET end_date = $1 WHERE id = $2`,
		newEndDate, oldMP.ID,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("extending subscription: %w", err)
	}

	totalAmount := req.AmountPaid + req.Discount
	due := totalAmount - req.Discount - req.AmountPaid
	if due < 0 {
		due = 0
	}

	particular := fmt.Sprintf("Extension: %s (+%d days)", pkg.Name, req.ExtraDays)
	var pmt Payment
	err = tx.QueryRow(ctx,
		`INSERT INTO payments (member_package_id, user_id, organization_id, particular, total_amount, discount, paid_amount, due, payment_method, payment_reference)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		 RETURNING id, member_package_id, particular, total_amount, discount, paid_amount, due, COALESCE(payment_method, ''), COALESCE(payment_reference, ''), payment_date`,
		oldMP.ID, memberID, orgID, particular, totalAmount, req.Discount, req.AmountPaid, due,
		nilIfEmpty(req.PaymentMethod), nilIfEmpty(req.PaymentReference),
	).Scan(&pmt.ID, &pmt.MemberPackageID, &pmt.Particular, &pmt.TotalAmount, &pmt.Discount, &pmt.PaidAmount, &pmt.Due, &pmt.PaymentMethod, &pmt.PaymentReference, &pmt.PaymentDate)
	if err != nil {
		return nil, nil, fmt.Errorf("inserting extension payment: %w", err)
	}

	if err := recordPackageIncome(ctx, tx, orgID, particular, req.AmountPaid,
		req.PaymentMethod, req.PaymentReference, entryByUserID); err != nil {
		return nil, nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, nil, fmt.Errorf("commit transaction: %w", err)
	}

	daysRemaining := int(time.Until(newEndDate).Hours() / 24)
	if daysRemaining < 0 {
		daysRemaining = 0
	}

	sub := &MemberSubscription{
		ID:            oldMP.ID,
		PackageID:     oldMP.PackageID,
		PackageName:   pkg.Name,
		StartDate:     oldMP.StartDate.Format("2006-01-02"),
		EndDate:       newEndDate.Format("2006-01-02"),
		Status:        oldMP.Status,
		DaysRemaining: daysRemaining,
		AmountPaid:    oldMP.AmountPaid + req.AmountPaid,
		PaymentMethod: oldMP.PaymentMethod,
		CreatedAt:     oldMP.CreatedAt,
	}

	return sub, &pmt, nil
}

// ListPayments retrieves payment history for a member within an organization.
func (r *Repository) ListPayments(ctx context.Context, orgID, memberID string, limit, offset int) ([]Payment, int, error) {
	var total int
	err := r.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM payments WHERE user_id = $1 AND organization_id = $2`,
		memberID, orgID,
	).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("counting payments: %w", err)
	}

	rows, err := r.db.Query(ctx,
		`SELECT id, member_package_id, particular, total_amount, discount, paid_amount, due,
		        COALESCE(payment_method, ''), COALESCE(payment_reference, ''), payment_date
		 FROM payments
		 WHERE user_id = $1 AND organization_id = $2
		 ORDER BY payment_date DESC
		 LIMIT $3 OFFSET $4`,
		memberID, orgID, limit, offset,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("querying payments: %w", err)
	}
	defer rows.Close()

	var payments []Payment
	for rows.Next() {
		var p Payment
		if err := rows.Scan(&p.ID, &p.MemberPackageID, &p.Particular, &p.TotalAmount,
			&p.Discount, &p.PaidAmount, &p.Due, &p.PaymentMethod, &p.PaymentReference,
			&p.PaymentDate); err != nil {
			return nil, 0, fmt.Errorf("scanning payment: %w", err)
		}
		payments = append(payments, p)
	}
	return payments, total, rows.Err()
}

// ListSubscriptions retrieves subscription history for a member within an organization.
func (r *Repository) ListSubscriptions(ctx context.Context, orgID, memberID string, limit, offset int) ([]MemberSubscription, int, error) {
	var total int
	err := r.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM member_packages WHERE user_id = $1 AND organization_id = $2`,
		memberID, orgID,
	).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("counting subscriptions: %w", err)
	}

	rows, err := r.db.Query(ctx,
		`SELECT mp.id, mp.package_id, p.name, mp.start_date, mp.end_date, mp.status,
		        GREATEST((mp.end_date - CURRENT_DATE)::int, 0) AS days_remaining,
		        mp.amount_paid, COALESCE(mp.payment_method, ''), mp.created_at
		 FROM member_packages mp
		 JOIN packages p ON mp.package_id = p.id
		 WHERE mp.user_id = $1 AND mp.organization_id = $2
		 ORDER BY mp.created_at DESC
		 LIMIT $3 OFFSET $4`,
		memberID, orgID, limit, offset,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("querying subscriptions: %w", err)
	}
	defer rows.Close()

	var subs []MemberSubscription
	for rows.Next() {
		var s MemberSubscription
		var startDate, endDate time.Time
		if err := rows.Scan(&s.ID, &s.PackageID, &s.PackageName, &startDate, &endDate,
			&s.Status, &s.DaysRemaining, &s.AmountPaid, &s.PaymentMethod, &s.CreatedAt); err != nil {
			return nil, 0, fmt.Errorf("scanning subscription: %w", err)
		}
		s.StartDate = startDate.Format("2006-01-02")
		s.EndDate = endDate.Format("2006-01-02")
		subs = append(subs, s)
	}
	return subs, total, rows.Err()
}

// VerifyMemberInOrg checks that the user is a member of the given org.
func (r *Repository) VerifyMemberInOrg(ctx context.Context, orgID, memberID string) error {
	var exists bool
	err := r.db.QueryRow(ctx,
		`SELECT EXISTS(
			SELECT 1 FROM organization_members
			WHERE user_id = $1 AND organization_id = $2 AND status = 'active'
		)`,
		memberID, orgID,
	).Scan(&exists)
	if err != nil {
		return fmt.Errorf("checking membership: %w", err)
	}
	if !exists {
		return pgx.ErrNoRows
	}
	return nil
}

// GetSubscriptionOwner returns the user_id that owns a member_package.
func (r *Repository) GetSubscriptionOwner(ctx context.Context, orgID, memberPackageID string) (string, error) {
	var ownerID string
	err := r.db.QueryRow(ctx,
		`SELECT user_id FROM member_packages WHERE id = $1 AND organization_id = $2`,
		memberPackageID, orgID,
	).Scan(&ownerID)
	if err != nil {
		return "", fmt.Errorf("subscription not found: %w", err)
	}
	return ownerID, nil
}

func nilIfEmpty(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// RaiseDueForBalance records the unpaid remainder of a package as a due.
//
// Assigning a package takes an amount_paid, and gyms routinely take part of it
// up front. Before this the shortfall went nowhere: the subscription recorded
// what was paid, and the amount still owed existed only in somebody's head.
//
// Written here rather than by the dues package to avoid a cycle, and scoped by
// active membership for the same reason dues.Create is.
func (r *Repository) RaiseDueForBalance(ctx context.Context, orgID, memberID string, balance float64, packageName, dueDate string) (string, error) {
	var dueID string
	err := r.db.QueryRow(ctx,
		`INSERT INTO dues (organization_id, user_id, amount, due_date, description, status)
		 SELECT $1, $2, $3, $4::date, $5, 'pending'
		 WHERE EXISTS (
		   SELECT 1 FROM organization_members om
		   WHERE om.organization_id = $1 AND om.user_id = $2 AND om.status = 'active'
		 )
		 RETURNING id`,
		orgID, memberID, balance, dueDate,
		fmt.Sprintf("Balance for %s", packageName),
	).Scan(&dueID)
	if err != nil {
		return "", fmt.Errorf("raising due for package balance: %w", err)
	}
	return dueID, nil
}

func defaultIfEmpty(s, fallback string) string {
	if s == "" {
		return fallback
	}
	return s
}

// recordPackageIncome posts money taken for a package to the ledger.
//
// Package payments were written to `payments` but never to `transactions`,
// which is what the accounts screen sums — so a gym's reported income showed
// dues collections only and missed most of its revenue. This runs inside the
// caller's transaction: a payment recorded against a member but absent from the
// books is exactly the bug being fixed, so the two must not be able to diverge.
func recordPackageIncome(ctx context.Context, tx pgx.Tx, orgID, description string,
	amount float64, method, reference, entryByUserID string) error {
	if amount <= 0 {
		return nil
	}
	_, err := tx.Exec(ctx,
		`INSERT INTO transactions
		        (organization_id, category, description, transaction_date,
		         transaction_type, amount, payment_type, reference, entry_by)
		 VALUES ($1, 'Subscription', $2, CURRENT_DATE, 'income', $3, $4, $5, $6)`,
		orgID, description, amount, defaultIfEmpty(method, "cash"),
		nilIfEmpty(reference), entryByUserID)
	if err != nil {
		return fmt.Errorf("recording income for package payment: %w", err)
	}
	return nil
}
