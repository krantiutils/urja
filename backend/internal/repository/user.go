package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/urja-gym/backend/internal/model"
)

var ErrUserNotFound = errors.New("user not found")

// UserRepository handles user persistence in PostgreSQL.
type UserRepository struct {
	pool *pgxpool.Pool
}

// NewUserRepository creates a new UserRepository.
func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

// FindByPhone looks up a user by their phone number.
func (r *UserRepository) FindByPhone(ctx context.Context, phone string) (*model.User, error) {
	var u model.User
	err := r.pool.QueryRow(ctx,
		`SELECT id, phone, name, role, is_active, created_at, updated_at
		 FROM users WHERE phone = $1`, phone,
	).Scan(&u.ID, &u.Phone, &u.Name, &u.Role, &u.IsActive, &u.CreatedAt, &u.UpdatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("querying user by phone: %w", err)
	}
	return &u, nil
}

// FindByID looks up a user by their ID.
func (r *UserRepository) FindByID(ctx context.Context, id string) (*model.User, error) {
	var u model.User
	err := r.pool.QueryRow(ctx,
		`SELECT id, phone, name, role, is_active, created_at, updated_at
		 FROM users WHERE id = $1`, id,
	).Scan(&u.ID, &u.Phone, &u.Name, &u.Role, &u.IsActive, &u.CreatedAt, &u.UpdatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("querying user by id: %w", err)
	}
	return &u, nil
}

// Create inserts a new user and returns it with the generated ID.
func (r *UserRepository) Create(ctx context.Context, phone string) (*model.User, error) {
	var u model.User
	err := r.pool.QueryRow(ctx,
		`INSERT INTO users (phone) VALUES ($1)
		 RETURNING id, phone, name, role, is_active, created_at, updated_at`, phone,
	).Scan(&u.ID, &u.Phone, &u.Name, &u.Role, &u.IsActive, &u.CreatedAt, &u.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("creating user: %w", err)
	}
	return &u, nil
}
