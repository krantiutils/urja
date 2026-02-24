package water

import (
	"context"
	"fmt"
	"log/slog"
	"time"
)

// Service handles water tracking business logic.
type Service struct {
	repo   *Repository
	logger *slog.Logger
}

// NewService creates a new water service.
func NewService(repo *Repository, logger *slog.Logger) *Service {
	return &Service{repo: repo, logger: logger}
}

// CreateWaterLogInput holds input for logging water intake.
type CreateWaterLogInput struct {
	AmountML int    `json:"amount_ml"`
	Date     string `json:"date"`
}

// CreateWaterLog validates input and creates a water log entry.
func (s *Service) CreateWaterLog(ctx context.Context, userID string, input *CreateWaterLogInput) (*WaterLog, error) {
	if input.AmountML <= 0 {
		return nil, fmt.Errorf("amount_ml must be positive")
	}
	if input.Date == "" {
		return nil, fmt.Errorf("date is required (YYYY-MM-DD)")
	}
	if _, err := time.Parse("2006-01-02", input.Date); err != nil {
		return nil, fmt.Errorf("date must be in YYYY-MM-DD format")
	}

	l, err := s.repo.CreateWaterLog(ctx, userID, input.AmountML, input.Date)
	if err != nil {
		return nil, err
	}

	s.logger.Info("water logged", "log_id", l.ID, "user_id", userID, "amount_ml", input.AmountML)
	return l, nil
}

// ListWaterLogs retrieves water logs for a user on a given date.
func (s *Service) ListWaterLogs(ctx context.Context, userID, date string) ([]WaterLog, error) {
	if date == "" {
		return nil, fmt.Errorf("date is required (YYYY-MM-DD)")
	}
	return s.repo.ListWaterLogs(ctx, userID, date)
}

// DeleteWaterLog deletes a water log entry owned by the user.
func (s *Service) DeleteWaterLog(ctx context.Context, userID, logID string) error {
	if err := s.repo.DeleteWaterLog(ctx, userID, logID); err != nil {
		return err
	}
	s.logger.Info("water log deleted", "log_id", logID, "user_id", userID)
	return nil
}

// GetDailySummary returns the daily water summary including total intake and goal.
func (s *Service) GetDailySummary(ctx context.Context, userID, date string) (*WaterDailySummary, error) {
	if date == "" {
		return nil, fmt.Errorf("date is required (YYYY-MM-DD)")
	}

	totalML, err := s.repo.GetDailySummary(ctx, userID, date)
	if err != nil {
		return nil, err
	}

	goalML, err := s.repo.GetUserWaterGoal(ctx, userID)
	if err != nil {
		s.logger.Warn("failed to get user water goal, using default", "error", err, "user_id", userID)
		goalML = 2500
	}

	entries, err := s.repo.ListWaterLogs(ctx, userID, date)
	if err != nil {
		return nil, err
	}

	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, fmt.Errorf("parsing date for summary: %w", err)
	}

	return &WaterDailySummary{
		Date:    parsedDate,
		TotalML: totalML,
		GoalML:  goalML,
		Entries: entries,
	}, nil
}
