package exercise

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Exercise represents an exercise in the library.
type Exercise struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	NameNe         string    `json:"name_ne,omitempty"`
	Description    string    `json:"description,omitempty"`
	MuscleGroups   []string  `json:"muscle_groups"`
	Equipment      string    `json:"equipment"`
	ExerciseType   string    `json:"exercise_type"`
	Difficulty     string    `json:"difficulty"`
	CaloriesPerMin float64   `json:"calories_per_minute"`
	Instructions   string    `json:"instructions,omitempty"`
	ImageURL       string    `json:"image_url,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
}

// WorkoutProgram represents a structured workout program.
type WorkoutProgram struct {
	ID            string       `json:"id"`
	Name          string       `json:"name"`
	NameNe        string       `json:"name_ne,omitempty"`
	Description   string       `json:"description,omitempty"`
	ProgramType   string       `json:"program_type"`
	Difficulty    string       `json:"difficulty"`
	DurationWeeks int          `json:"duration_weeks"`
	Equipment     string       `json:"equipment"`
	Goal          string       `json:"goal,omitempty"`
	ImageURL      string       `json:"image_url,omitempty"`
	CreatedAt     time.Time    `json:"created_at"`
	Days          []ProgramDay `json:"days,omitempty"`
}

// ProgramDay represents a single day in a workout program.
type ProgramDay struct {
	ID         string          `json:"id"`
	ProgramID  string          `json:"program_id"`
	WeekNumber int             `json:"week_number"`
	DayNumber  int             `json:"day_number"`
	DayName    string          `json:"day_name,omitempty"`
	Exercises  json.RawMessage `json:"exercises"`
	RestDay    bool            `json:"rest_day"`
}

// UserProgramEnrollment represents a user's enrollment in a workout program.
type UserProgramEnrollment struct {
	ID           string          `json:"id"`
	UserID       string          `json:"user_id"`
	ProgramID    string          `json:"program_id"`
	StartedAt    time.Time       `json:"started_at"`
	CurrentWeek  int             `json:"current_week"`
	CurrentDay   int             `json:"current_day"`
	Status       string          `json:"status"`
	Program      *WorkoutProgram `json:"program,omitempty"`
	TodayWorkout *ProgramDay     `json:"today_workout,omitempty"`
}

// Repository handles exercise data persistence.
type Repository struct {
	db *pgxpool.Pool
}

// NewRepository creates a new exercise repository.
func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

// ListExercises retrieves exercises with optional filters.
func (r *Repository) ListExercises(ctx context.Context, exerciseType, equipment, muscleGroup string, limit, offset int) ([]Exercise, error) {
	query := `SELECT id, name, COALESCE(name_ne, ''), COALESCE(description, ''),
	                 muscle_groups, equipment, exercise_type, difficulty,
	                 COALESCE(calories_per_minute, 0), COALESCE(instructions, ''),
	                 COALESCE(image_url, ''), created_at
	          FROM exercises
	          WHERE 1=1`
	args := []interface{}{}
	argIdx := 1

	if exerciseType != "" {
		query += fmt.Sprintf(" AND exercise_type = $%d", argIdx)
		args = append(args, exerciseType)
		argIdx++
	}
	if equipment != "" {
		query += fmt.Sprintf(" AND equipment = $%d", argIdx)
		args = append(args, equipment)
		argIdx++
	}
	if muscleGroup != "" {
		query += fmt.Sprintf(" AND $%d = ANY(muscle_groups)", argIdx)
		args = append(args, muscleGroup)
		argIdx++
	}

	query += fmt.Sprintf(" ORDER BY name ASC LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("listing exercises: %w", err)
	}
	defer rows.Close()

	var exercises []Exercise
	for rows.Next() {
		var e Exercise
		if err := rows.Scan(&e.ID, &e.Name, &e.NameNe, &e.Description,
			&e.MuscleGroups, &e.Equipment, &e.ExerciseType, &e.Difficulty,
			&e.CaloriesPerMin, &e.Instructions, &e.ImageURL, &e.CreatedAt); err != nil {
			return nil, fmt.Errorf("scanning exercise: %w", err)
		}
		exercises = append(exercises, e)
	}
	return exercises, rows.Err()
}

// GetExercise retrieves a single exercise by ID.
func (r *Repository) GetExercise(ctx context.Context, id string) (*Exercise, error) {
	var e Exercise
	err := r.db.QueryRow(ctx,
		`SELECT id, name, COALESCE(name_ne, ''), COALESCE(description, ''),
		        muscle_groups, equipment, exercise_type, difficulty,
		        COALESCE(calories_per_minute, 0), COALESCE(instructions, ''),
		        COALESCE(image_url, ''), created_at
		 FROM exercises
		 WHERE id = $1`,
		id,
	).Scan(&e.ID, &e.Name, &e.NameNe, &e.Description,
		&e.MuscleGroups, &e.Equipment, &e.ExerciseType, &e.Difficulty,
		&e.CaloriesPerMin, &e.Instructions, &e.ImageURL, &e.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("getting exercise: %w", err)
	}
	return &e, nil
}

// ListPrograms retrieves workout programs with optional filters.
func (r *Repository) ListPrograms(ctx context.Context, programType, difficulty, equipment string, limit, offset int) ([]WorkoutProgram, error) {
	query := `SELECT id, name, COALESCE(name_ne, ''), COALESCE(description, ''),
	                 COALESCE(program_type, ''), difficulty, duration_weeks,
	                 equipment, COALESCE(goal, ''), COALESCE(image_url, ''), created_at
	          FROM workout_programs
	          WHERE 1=1`
	args := []interface{}{}
	argIdx := 1

	if programType != "" {
		query += fmt.Sprintf(" AND program_type = $%d", argIdx)
		args = append(args, programType)
		argIdx++
	}
	if difficulty != "" {
		query += fmt.Sprintf(" AND difficulty = $%d", argIdx)
		args = append(args, difficulty)
		argIdx++
	}
	if equipment != "" {
		query += fmt.Sprintf(" AND equipment = $%d", argIdx)
		args = append(args, equipment)
		argIdx++
	}

	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("listing programs: %w", err)
	}
	defer rows.Close()

	var programs []WorkoutProgram
	for rows.Next() {
		var p WorkoutProgram
		if err := rows.Scan(&p.ID, &p.Name, &p.NameNe, &p.Description,
			&p.ProgramType, &p.Difficulty, &p.DurationWeeks,
			&p.Equipment, &p.Goal, &p.ImageURL, &p.CreatedAt); err != nil {
			return nil, fmt.Errorf("scanning program: %w", err)
		}
		programs = append(programs, p)
	}
	return programs, rows.Err()
}

// GetProgram retrieves a single workout program by ID, including all days ordered by week and day.
func (r *Repository) GetProgram(ctx context.Context, id string) (*WorkoutProgram, error) {
	var p WorkoutProgram
	err := r.db.QueryRow(ctx,
		`SELECT id, name, COALESCE(name_ne, ''), COALESCE(description, ''),
		        COALESCE(program_type, ''), difficulty, duration_weeks,
		        equipment, COALESCE(goal, ''), COALESCE(image_url, ''), created_at
		 FROM workout_programs
		 WHERE id = $1`,
		id,
	).Scan(&p.ID, &p.Name, &p.NameNe, &p.Description,
		&p.ProgramType, &p.Difficulty, &p.DurationWeeks,
		&p.Equipment, &p.Goal, &p.ImageURL, &p.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("getting program: %w", err)
	}

	// Fetch all program days.
	rows, err := r.db.Query(ctx,
		`SELECT id, program_id, week_number, day_number, COALESCE(day_name, ''),
		        exercises, COALESCE(rest_day, false)
		 FROM program_days
		 WHERE program_id = $1
		 ORDER BY week_number ASC, day_number ASC`,
		id,
	)
	if err != nil {
		return nil, fmt.Errorf("listing program days: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var d ProgramDay
		if err := rows.Scan(&d.ID, &d.ProgramID, &d.WeekNumber, &d.DayNumber,
			&d.DayName, &d.Exercises, &d.RestDay); err != nil {
			return nil, fmt.Errorf("scanning program day: %w", err)
		}
		p.Days = append(p.Days, d)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating program days: %w", err)
	}

	return &p, nil
}

// EnrollUser enrolls a user in a workout program, or re-enrolls if already enrolled.
func (r *Repository) EnrollUser(ctx context.Context, userID, programID string) (*UserProgramEnrollment, error) {
	var e UserProgramEnrollment
	err := r.db.QueryRow(ctx,
		`INSERT INTO user_program_enrollments (user_id, program_id)
		 VALUES ($1, $2)
		 ON CONFLICT (user_id, program_id)
		 DO UPDATE SET status = 'active', current_week = 1, current_day = 1, started_at = NOW()
		 RETURNING id, user_id, program_id, started_at, current_week, current_day, status`,
		userID, programID,
	).Scan(&e.ID, &e.UserID, &e.ProgramID, &e.StartedAt, &e.CurrentWeek, &e.CurrentDay, &e.Status)
	if err != nil {
		return nil, fmt.Errorf("enrolling user: %w", err)
	}
	return &e, nil
}

// GetCurrentEnrollment retrieves the active enrollment for a user, joined with the workout program.
func (r *Repository) GetCurrentEnrollment(ctx context.Context, userID string) (*UserProgramEnrollment, error) {
	var e UserProgramEnrollment
	var p WorkoutProgram
	err := r.db.QueryRow(ctx,
		`SELECT e.id, e.user_id, e.program_id, e.started_at, e.current_week, e.current_day, e.status,
		        wp.id, wp.name, COALESCE(wp.name_ne, ''), COALESCE(wp.description, ''),
		        COALESCE(wp.program_type, ''), wp.difficulty, wp.duration_weeks,
		        wp.equipment, COALESCE(wp.goal, ''), COALESCE(wp.image_url, ''), wp.created_at
		 FROM user_program_enrollments e
		 JOIN workout_programs wp ON wp.id = e.program_id
		 WHERE e.user_id = $1 AND e.status = 'active'
		 ORDER BY e.started_at DESC
		 LIMIT 1`,
		userID,
	).Scan(&e.ID, &e.UserID, &e.ProgramID, &e.StartedAt, &e.CurrentWeek, &e.CurrentDay, &e.Status,
		&p.ID, &p.Name, &p.NameNe, &p.Description,
		&p.ProgramType, &p.Difficulty, &p.DurationWeeks,
		&p.Equipment, &p.Goal, &p.ImageURL, &p.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("getting current enrollment: %w", err)
	}
	e.Program = &p
	return &e, nil
}

// CompleteDay advances the user's progress in their active enrollment.
// It increments the current day, rolling over to the next week when necessary,
// and marks the enrollment as completed when the program is finished.
func (r *Repository) CompleteDay(ctx context.Context, userID string) (*UserProgramEnrollment, error) {
	// Get current enrollment with program info to determine bounds.
	enrollment, err := r.GetCurrentEnrollment(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("no active enrollment found: %w", err)
	}

	// Count days per week for this program.
	var maxDayInWeek int
	err = r.db.QueryRow(ctx,
		`SELECT COALESCE(MAX(day_number), 0)
		 FROM program_days
		 WHERE program_id = $1 AND week_number = $2`,
		enrollment.ProgramID, enrollment.CurrentWeek,
	).Scan(&maxDayInWeek)
	if err != nil {
		return nil, fmt.Errorf("getting max day: %w", err)
	}

	nextWeek := enrollment.CurrentWeek
	nextDay := enrollment.CurrentDay + 1
	newStatus := "active"

	if nextDay > maxDayInWeek {
		// Move to the next week.
		nextWeek++
		nextDay = 1

		if nextWeek > enrollment.Program.DurationWeeks {
			// Program completed.
			newStatus = "completed"
			nextWeek = enrollment.Program.DurationWeeks
			nextDay = maxDayInWeek
		}
	}

	var updated UserProgramEnrollment
	err = r.db.QueryRow(ctx,
		`UPDATE user_program_enrollments
		 SET current_week = $1, current_day = $2, status = $3
		 WHERE user_id = $4 AND status = 'active'
		 RETURNING id, user_id, program_id, started_at, current_week, current_day, status`,
		nextWeek, nextDay, newStatus, userID,
	).Scan(&updated.ID, &updated.UserID, &updated.ProgramID, &updated.StartedAt,
		&updated.CurrentWeek, &updated.CurrentDay, &updated.Status)
	if err != nil {
		return nil, fmt.Errorf("completing day: %w", err)
	}
	updated.Program = enrollment.Program
	return &updated, nil
}

// GetProgramDay retrieves a specific day from a workout program.
func (r *Repository) GetProgramDay(ctx context.Context, programID string, week, day int) (*ProgramDay, error) {
	var d ProgramDay
	err := r.db.QueryRow(ctx,
		`SELECT id, program_id, week_number, day_number, COALESCE(day_name, ''),
		        exercises, COALESCE(rest_day, false)
		 FROM program_days
		 WHERE program_id = $1 AND week_number = $2 AND day_number = $3`,
		programID, week, day,
	).Scan(&d.ID, &d.ProgramID, &d.WeekNumber, &d.DayNumber,
		&d.DayName, &d.Exercises, &d.RestDay)
	if err != nil {
		return nil, fmt.Errorf("getting program day: %w", err)
	}
	return &d, nil
}
