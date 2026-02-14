package leaderboard

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// RankedMember represents a single entry in the leaderboard.
type RankedMember struct {
	Rank      int    `json:"rank"`
	MemberID  string `json:"member_id"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url,omitempty"`
	Value     int    `json:"value"`
	Metric    string `json:"metric"`
}

// Repository handles leaderboard data access.
type Repository struct {
	db *pgxpool.Pool
}

// NewRepository creates a new leaderboard repository.
func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

// GetWeeklyRanking ranks members by current_streak descending.
func (r *Repository) GetWeeklyRanking(ctx context.Context, orgID string, limit, offset int) ([]RankedMember, error) {
	rows, err := r.db.Query(ctx, `
		WITH ranked AS (
			SELECT
				u.id AS member_id,
				u.name,
				COALESCE(u.avatar_url, '') AS avatar_url,
				s.current_streak AS value,
				ROW_NUMBER() OVER (ORDER BY s.current_streak DESC, u.name ASC) AS rank
			FROM streaks s
			JOIN users u ON s.member_id = u.id
			JOIN organization_members om
				ON om.user_id = u.id AND om.organization_id = s.organization_id
			WHERE s.organization_id = $1
				AND om.status = 'active'
				AND COALESCE((u.privacy_settings->>'show_on_leaderboard')::boolean, true) = true
				AND s.current_streak > 0
		)
		SELECT member_id, name, avatar_url, value, rank
		FROM ranked
		ORDER BY rank
		LIMIT $2 OFFSET $3`,
		orgID, limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("querying weekly leaderboard: %w", err)
	}
	defer rows.Close()

	var results []RankedMember
	for rows.Next() {
		var m RankedMember
		if err := rows.Scan(&m.MemberID, &m.Name, &m.AvatarURL, &m.Value, &m.Rank); err != nil {
			return nil, fmt.Errorf("scanning weekly leaderboard: %w", err)
		}
		m.Metric = "current_streak"
		results = append(results, m)
	}
	return results, rows.Err()
}

// GetMonthlyRanking ranks members by total check-ins this calendar month (NPT).
func (r *Repository) GetMonthlyRanking(ctx context.Context, orgID string, limit, offset int) ([]RankedMember, error) {
	rows, err := r.db.Query(ctx, `
		WITH ranked AS (
			SELECT
				u.id AS member_id,
				u.name,
				COALESCE(u.avatar_url, '') AS avatar_url,
				COUNT(a.id)::int AS value,
				ROW_NUMBER() OVER (ORDER BY COUNT(a.id) DESC, u.name ASC) AS rank
			FROM attendance a
			JOIN users u ON a.user_id = u.id
			JOIN organization_members om
				ON om.user_id = u.id AND om.organization_id = a.organization_id
			WHERE a.organization_id = $1
				AND a.check_in_at >= date_trunc('month', NOW() AT TIME ZONE 'Asia/Kathmandu')
					AT TIME ZONE 'Asia/Kathmandu'
				AND om.status = 'active'
				AND COALESCE((u.privacy_settings->>'show_on_leaderboard')::boolean, true) = true
			GROUP BY u.id, u.name, u.avatar_url
			HAVING COUNT(a.id) > 0
		)
		SELECT member_id, name, avatar_url, value, rank
		FROM ranked
		ORDER BY rank
		LIMIT $2 OFFSET $3`,
		orgID, limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("querying monthly leaderboard: %w", err)
	}
	defer rows.Close()

	var results []RankedMember
	for rows.Next() {
		var m RankedMember
		if err := rows.Scan(&m.MemberID, &m.Name, &m.AvatarURL, &m.Value, &m.Rank); err != nil {
			return nil, fmt.Errorf("scanning monthly leaderboard: %w", err)
		}
		m.Metric = "check_ins"
		results = append(results, m)
	}
	return results, rows.Err()
}

// GetAllTimeRanking ranks members by longest_streak descending.
func (r *Repository) GetAllTimeRanking(ctx context.Context, orgID string, limit, offset int) ([]RankedMember, error) {
	rows, err := r.db.Query(ctx, `
		WITH ranked AS (
			SELECT
				u.id AS member_id,
				u.name,
				COALESCE(u.avatar_url, '') AS avatar_url,
				s.longest_streak AS value,
				ROW_NUMBER() OVER (ORDER BY s.longest_streak DESC, u.name ASC) AS rank
			FROM streaks s
			JOIN users u ON s.member_id = u.id
			JOIN organization_members om
				ON om.user_id = u.id AND om.organization_id = s.organization_id
			WHERE s.organization_id = $1
				AND om.status = 'active'
				AND COALESCE((u.privacy_settings->>'show_on_leaderboard')::boolean, true) = true
				AND s.longest_streak > 0
		)
		SELECT member_id, name, avatar_url, value, rank
		FROM ranked
		ORDER BY rank
		LIMIT $2 OFFSET $3`,
		orgID, limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("querying all-time leaderboard: %w", err)
	}
	defer rows.Close()

	var results []RankedMember
	for rows.Next() {
		var m RankedMember
		if err := rows.Scan(&m.MemberID, &m.Name, &m.AvatarURL, &m.Value, &m.Rank); err != nil {
			return nil, fmt.Errorf("scanning all-time leaderboard: %w", err)
		}
		m.Metric = "longest_streak"
		results = append(results, m)
	}
	return results, rows.Err()
}
