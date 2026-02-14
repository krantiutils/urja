package leaderboard

import (
	"context"
	"fmt"
	"log/slog"
)

// LeaderboardResponse is the top-level response for the leaderboard endpoint.
type LeaderboardResponse struct {
	Period   string         `json:"period"`
	Rankings []RankedMember `json:"rankings"`
}

// Service handles leaderboard business logic.
type Service struct {
	repo   *Repository
	logger *slog.Logger
}

// NewService creates a new leaderboard service.
func NewService(repo *Repository, logger *slog.Logger) *Service {
	return &Service{repo: repo, logger: logger}
}

// GetLeaderboard retrieves the leaderboard for an organization by period.
// Valid periods: "weekly", "monthly", "alltime".
func (s *Service) GetLeaderboard(ctx context.Context, orgID, period string, limit, offset int) (*LeaderboardResponse, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	var rankings []RankedMember
	var err error

	switch period {
	case "weekly":
		rankings, err = s.repo.GetWeeklyRanking(ctx, orgID, limit, offset)
	case "monthly":
		rankings, err = s.repo.GetMonthlyRanking(ctx, orgID, limit, offset)
	case "alltime":
		rankings, err = s.repo.GetAllTimeRanking(ctx, orgID, limit, offset)
	default:
		return nil, fmt.Errorf("invalid period %q: must be weekly, monthly, or alltime", period)
	}

	if err != nil {
		return nil, fmt.Errorf("fetching %s leaderboard: %w", period, err)
	}

	if rankings == nil {
		rankings = []RankedMember{}
	}

	return &LeaderboardResponse{
		Period:   period,
		Rankings: rankings,
	}, nil
}
