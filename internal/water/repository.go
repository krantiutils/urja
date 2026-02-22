package water

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// WaterLog represents a single water intake entry.
type WaterLog struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	AmountML   int       `json:"amount_ml"`
	LoggedDate string    `json:"logged_date"`
	LoggedAt   time.Time `json:"logged_at"`
}

// WaterDailySummary represents aggregated water intake for a single day.
type WaterDailySummary struct {
	Date    string     `json:"date"`
	TotalML int        `json:"total_ml"`
	GoalML  int        `json:"goal_ml"`
	Entries []WaterLog `json:"entries"`
}

// Repository handles water tracking data persistence.
type Repository struct {
	db *pgxpool.Pool
}

// NewRepository creates a new water repository.
func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

// CreateWaterLog inserts a new water log entry.
func (r *Repository) CreateWaterLog(ctx context.Context, userID string, amountML int, loggedDate string) (*WaterLog, error) {
	var l WaterLog
	err := r.db.QueryRow(ctx,
		`INSERT INTO water_logs (user_id, amount_ml, logged_date)
		 VALUES ($1, $2, $3)
		 RETURNING id, user_id, amount_ml, logged_date, logged_at`,
		userID, amountML, loggedDate,
	).Scan(&l.ID, &l.UserID, &l.AmountML, &l.LoggedDate, &l.LoggedAt)
	if err != nil {
		return nil, fmt.Errorf("creating water log: %w", err)
	}
	return &l, nil
}

// ListWaterLogs retrieves water logs for a user on a given date.
func (r *Repository) ListWaterLogs(ctx context.Context, userID, date string) ([]WaterLog, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, user_id, amount_ml, logged_date, logged_at
		 FROM water_logs
		 WHERE user_id = $1 AND logged_date = $2
		 ORDER BY logged_at ASC`,
		userID, date,
	)
	if err != nil {
		return nil, fmt.Errorf("listing water logs: %w", err)
	}
	defer rows.Close()

	var logs []WaterLog
	for rows.Next() {
		var l WaterLog
		if err := rows.Scan(&l.ID, &l.UserID, &l.AmountML, &l.LoggedDate, &l.LoggedAt); err != nil {
			return nil, fmt.Errorf("scanning water log: %w", err)
		}
		logs = append(logs, l)
	}
	return logs, rows.Err()
}

// DeleteWaterLog deletes a water log entry owned by the user.
func (r *Repository) DeleteWaterLog(ctx context.Context, userID, logID string) error {
	tag, err := r.db.Exec(ctx,
		`DELETE FROM water_logs WHERE id = $1 AND user_id = $2`,
		logID, userID,
	)
	if err != nil {
		return fmt.Errorf("deleting water log: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("water log not found")
	}
	return nil
}

// GetDailySummary returns total_ml for a user on a given date.
func (r *Repository) GetDailySummary(ctx context.Context, userID, date string) (int, error) {
	var totalML int
	err := r.db.QueryRow(ctx,
		`SELECT COALESCE(SUM(amount_ml), 0)
		 FROM water_logs
		 WHERE user_id = $1 AND logged_date = $2`,
		userID, date,
	).Scan(&totalML)
	if err != nil {
		return 0, fmt.Errorf("getting daily water summary: %w", err)
	}
	return totalML, nil
}

// GetUserWaterGoal returns the daily_water_goal_ml for a user from the users table.
func (r *Repository) GetUserWaterGoal(ctx context.Context, userID string) (int, error) {
	var goalML int
	err := r.db.QueryRow(ctx,
		`SELECT COALESCE(daily_water_goal_ml, 2500)
		 FROM users
		 WHERE id = $1`,
		userID,
	).Scan(&goalML)
	if err != nil {
		return 2500, fmt.Errorf("getting user water goal: %w", err)
	}
	return goalML, nil
}
