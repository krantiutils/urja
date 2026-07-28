package boxing

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ErrNotFound is returned when a profile or bout does not exist within the
// requested organization.
var ErrNotFound = errors.New("not found")

// Repository handles boxing profile persistence. Every query is scoped by
// organization_id: a member's profile is per-gym, not global.
type Repository struct {
	db *pgxpool.Pool
}

// NewRepository creates a new boxing repository.
func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

const profileColumns = `id, user_id, organization_id, COALESCE(stance, ''), COALESCE(weight_class, ''),
	COALESCE(skill_level, ''), sparring_cleared, sparring_cleared_at,
	COALESCE(sparring_cleared_by::text, ''), reach_cm, COALESCE(notes, ''),
	created_at, updated_at`

func scanProfile(row pgx.Row) (*Profile, error) {
	var p Profile
	err := row.Scan(&p.ID, &p.UserID, &p.OrgID, &p.Stance, &p.WeightClass,
		&p.SkillLevel, &p.SparringCleared, &p.SparringClearedAt,
		&p.SparringClearedBy, &p.ReachCm, &p.Notes, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// GetProfile retrieves a member's boxing profile within an organization.
func (r *Repository) GetProfile(ctx context.Context, orgID, userID string) (*Profile, error) {
	p, err := scanProfile(r.db.QueryRow(ctx,
		`SELECT `+profileColumns+` FROM member_boxing_profiles
		 WHERE organization_id = $1 AND user_id = $2`, orgID, userID))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("getting boxing profile: %w", err)
	}
	return p, nil
}

// ProfileUpdate carries the member-editable fields. Sparring clearance is
// deliberately absent — see SetSparringClearance.
type ProfileUpdate struct {
	Stance      string
	WeightClass string
	SkillLevel  string
	ReachCm     *float64
	Notes       string
}

// UpsertProfile creates or updates the member-editable fields of a profile.
// It never touches the sparring clearance columns, so a member saving their
// profile cannot clear themselves even if the handler were bypassed.
func (r *Repository) UpsertProfile(ctx context.Context, orgID, userID string, in ProfileUpdate) (*Profile, error) {
	p, err := scanProfile(r.db.QueryRow(ctx,
		`INSERT INTO member_boxing_profiles
			(user_id, organization_id, stance, weight_class, skill_level, reach_cm, notes)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)
		 ON CONFLICT (user_id, organization_id) DO UPDATE
		 SET stance = EXCLUDED.stance,
		     weight_class = EXCLUDED.weight_class,
		     skill_level = EXCLUDED.skill_level,
		     reach_cm = EXCLUDED.reach_cm,
		     notes = EXCLUDED.notes,
		     updated_at = NOW()
		 RETURNING `+profileColumns,
		userID, orgID, nilIfEmpty(in.Stance), nilIfEmpty(in.WeightClass),
		nilIfEmpty(in.SkillLevel), in.ReachCm, nilIfEmpty(in.Notes)))
	if err != nil {
		return nil, fmt.Errorf("saving boxing profile: %w", err)
	}
	return p, nil
}

// SetSparringClearance grants or revokes sparring clearance. Separate from
// UpsertProfile because it is a staff-only safety decision and is audited.
func (r *Repository) SetSparringClearance(ctx context.Context, orgID, userID, clearedBy string, cleared bool) (*Profile, error) {
	var clearedAt interface{}
	var by interface{}
	if cleared {
		clearedAt = time.Now()
		by = clearedBy
	}

	p, err := scanProfile(r.db.QueryRow(ctx,
		`INSERT INTO member_boxing_profiles
			(user_id, organization_id, sparring_cleared, sparring_cleared_at, sparring_cleared_by)
		 VALUES ($1, $2, $3, $4, $5)
		 ON CONFLICT (user_id, organization_id) DO UPDATE
		 SET sparring_cleared = EXCLUDED.sparring_cleared,
		     sparring_cleared_at = EXCLUDED.sparring_cleared_at,
		     sparring_cleared_by = EXCLUDED.sparring_cleared_by,
		     updated_at = NOW()
		 RETURNING `+profileColumns,
		userID, orgID, cleared, clearedAt, by))
	if err != nil {
		return nil, fmt.Errorf("setting sparring clearance: %w", err)
	}
	return p, nil
}

// --- Bouts ---

const boutColumns = `id, user_id, organization_id, bout_date, COALESCE(opponent, ''),
	COALESCE(event_name, ''), result, COALESCE(method, ''), rounds,
	COALESCE(weight_class, ''), COALESCE(notes, ''), created_at`

func scanBout(row pgx.Row) (*Bout, error) {
	var b Bout
	err := row.Scan(&b.ID, &b.UserID, &b.OrgID, &b.BoutDate, &b.Opponent,
		&b.EventName, &b.Result, &b.Method, &b.Rounds, &b.WeightClass,
		&b.Notes, &b.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &b, nil
}

// ListBouts returns a member's bouts within an organization, newest first.
func (r *Repository) ListBouts(ctx context.Context, orgID, userID string) ([]Bout, error) {
	rows, err := r.db.Query(ctx,
		`SELECT `+boutColumns+` FROM bout_records
		 WHERE organization_id = $1 AND user_id = $2
		 ORDER BY bout_date DESC, created_at DESC`, orgID, userID)
	if err != nil {
		return nil, fmt.Errorf("listing bouts: %w", err)
	}
	defer rows.Close()

	bouts := []Bout{}
	for rows.Next() {
		var b Bout
		if err := rows.Scan(&b.ID, &b.UserID, &b.OrgID, &b.BoutDate, &b.Opponent,
			&b.EventName, &b.Result, &b.Method, &b.Rounds, &b.WeightClass,
			&b.Notes, &b.CreatedAt); err != nil {
			return nil, fmt.Errorf("scanning bout: %w", err)
		}
		bouts = append(bouts, b)
	}
	return bouts, rows.Err()
}

// CreateBout records a bout.
func (r *Repository) CreateBout(ctx context.Context, orgID, userID string, b *Bout) (*Bout, error) {
	created, err := scanBout(r.db.QueryRow(ctx,
		`INSERT INTO bout_records
			(user_id, organization_id, bout_date, opponent, event_name, result, method, rounds, weight_class, notes)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		 RETURNING `+boutColumns,
		userID, orgID, b.BoutDate, nilIfEmpty(b.Opponent), nilIfEmpty(b.EventName),
		b.Result, nilIfEmpty(b.Method), b.Rounds, nilIfEmpty(b.WeightClass), nilIfEmpty(b.Notes)))
	if err != nil {
		return nil, fmt.Errorf("creating bout: %w", err)
	}
	return created, nil
}

// DeleteBout removes a bout belonging to a member within an organization.
func (r *Repository) DeleteBout(ctx context.Context, orgID, userID, boutID string) error {
	tag, err := r.db.Exec(ctx,
		`DELETE FROM bout_records WHERE id = $1 AND organization_id = $2 AND user_id = $3`,
		boutID, orgID, userID)
	if err != nil {
		return fmt.Errorf("deleting bout: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// LatestWeightKg reads a member's most recent recorded weight from the existing
// health metrics, so weight class can be derived rather than re-entered.
func (r *Repository) LatestWeightKg(ctx context.Context, userID string) (float64, error) {
	var kg float64
	err := r.db.QueryRow(ctx,
		`SELECT (value->>'weight_kg')::numeric FROM health_metrics
		 WHERE member_id = $1 AND metric_type = 'weight'
		   AND value ? 'weight_kg'
		 ORDER BY recorded_at DESC LIMIT 1`, userID).Scan(&kg)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, ErrNotFound
	}
	if err != nil {
		return 0, fmt.Errorf("reading latest weight: %w", err)
	}
	return kg, nil
}

// IsOrgMember reports whether a user has an active membership in an org.
// Used to reject staff writes targeting a member of a different gym.
func (r *Repository) IsOrgMember(ctx context.Context, orgID, userID string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM organization_members
		 WHERE organization_id = $1 AND user_id = $2 AND status = 'active')`,
		orgID, userID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("checking org membership: %w", err)
	}
	return exists, nil
}

func nilIfEmpty(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}
