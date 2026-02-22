package exercise

import (
	"context"
	"log/slog"
)

// Service handles exercise business logic.
type Service struct {
	repo   *Repository
	logger *slog.Logger
}

// NewService creates a new exercise service.
func NewService(repo *Repository, logger *slog.Logger) *Service {
	return &Service{repo: repo, logger: logger}
}

// ListExercises returns exercises with optional filters, validating limit/offset.
func (s *Service) ListExercises(ctx context.Context, exerciseType, equipment, muscleGroup string, limit, offset int) ([]Exercise, error) {
	if limit <= 0 || limit > 50 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	return s.repo.ListExercises(ctx, exerciseType, equipment, muscleGroup, limit, offset)
}

// GetExercise returns a single exercise by ID.
func (s *Service) GetExercise(ctx context.Context, id string) (*Exercise, error) {
	return s.repo.GetExercise(ctx, id)
}

// ListPrograms returns workout programs with optional filters, validating limit/offset.
func (s *Service) ListPrograms(ctx context.Context, programType, difficulty, equipment string, limit, offset int) ([]WorkoutProgram, error) {
	if limit <= 0 || limit > 50 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	return s.repo.ListPrograms(ctx, programType, difficulty, equipment, limit, offset)
}

// GetProgram returns a single workout program by ID, including all days.
func (s *Service) GetProgram(ctx context.Context, id string) (*WorkoutProgram, error) {
	return s.repo.GetProgram(ctx, id)
}

// EnrollInProgram enrolls a user in a workout program.
func (s *Service) EnrollInProgram(ctx context.Context, userID, programID string) (*UserProgramEnrollment, error) {
	// Verify the program exists.
	_, err := s.repo.GetProgram(ctx, programID)
	if err != nil {
		return nil, err
	}

	enrollment, err := s.repo.EnrollUser(ctx, userID, programID)
	if err != nil {
		return nil, err
	}

	s.logger.Info("user enrolled in program", "user_id", userID, "program_id", programID)
	return enrollment, nil
}

// GetCurrentProgram retrieves the user's active enrollment with today's workout.
func (s *Service) GetCurrentProgram(ctx context.Context, userID string) (*UserProgramEnrollment, error) {
	enrollment, err := s.repo.GetCurrentEnrollment(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Fetch today's workout based on current week and day.
	day, err := s.repo.GetProgramDay(ctx, enrollment.ProgramID, enrollment.CurrentWeek, enrollment.CurrentDay)
	if err != nil {
		s.logger.Warn("could not fetch today's workout", "error", err,
			"program_id", enrollment.ProgramID,
			"week", enrollment.CurrentWeek, "day", enrollment.CurrentDay)
	} else {
		enrollment.TodayWorkout = day
	}

	return enrollment, nil
}

// CompleteDay advances the user's program progress.
func (s *Service) CompleteDay(ctx context.Context, userID string) (*UserProgramEnrollment, error) {
	enrollment, err := s.repo.CompleteDay(ctx, userID)
	if err != nil {
		return nil, err
	}

	s.logger.Info("program day completed", "user_id", userID,
		"program_id", enrollment.ProgramID,
		"current_week", enrollment.CurrentWeek,
		"current_day", enrollment.CurrentDay,
		"status", enrollment.Status)
	return enrollment, nil
}
