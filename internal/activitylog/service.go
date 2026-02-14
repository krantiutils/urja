package activitylog

import (
	"context"
	"log/slog"
	"time"
)

// Service handles activity log business logic.
type Service struct {
	repo   *Repository
	logger *slog.Logger
}

// NewService creates a new activity log service.
func NewService(repo *Repository, logger *slog.Logger) *Service {
	return &Service{repo: repo, logger: logger}
}

// Log appends a new activity log entry.
func (s *Service) Log(ctx context.Context, orgID string, userID *string, eventType EventType, description string, metadata map[string]interface{}) {
	_, err := s.repo.Create(ctx, orgID, userID, eventType, description, metadata)
	if err != nil {
		s.logger.Error("failed to create activity log",
			"error", err,
			"org_id", orgID,
			"event_type", string(eventType),
		)
	}
}

// DateGroup groups activity log entries by date for timeline display.
type DateGroup struct {
	Date    string  `json:"date"`
	Entries []Entry `json:"entries"`
}

// List retrieves activity logs grouped by date.
func (s *Service) List(ctx context.Context, orgID string, from, to *time.Time, limit, offset int) ([]DateGroup, int, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	entries, total, err := s.repo.List(ctx, orgID, from, to, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	// Group entries by date for timeline UI
	groups := groupByDate(entries)
	return groups, total, nil
}

func groupByDate(entries []Entry) []DateGroup {
	if len(entries) == 0 {
		return []DateGroup{}
	}

	var groups []DateGroup
	var current *DateGroup

	for _, e := range entries {
		dateStr := e.CreatedAt.Format("2006-01-02")
		if current == nil || current.Date != dateStr {
			if current != nil {
				groups = append(groups, *current)
			}
			current = &DateGroup{
				Date:    dateStr,
				Entries: []Entry{e},
			}
		} else {
			current.Entries = append(current.Entries, e)
		}
	}
	if current != nil {
		groups = append(groups, *current)
	}

	return groups
}
