package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/s-usmonalizoda25/ticket-engine/internal/domain"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	const query = `
        INSERT INTO users (email, name, created_at, updated_at)
        VALUES ($1, $2, NOW(), NOW())
        RETURNING id, created_at, updated_at
    `
	err := r.pool.QueryRow(ctx, query, user.Email, user.Name).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (r *UserRepository) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	const query = `
        SELECT id, email, name, created_at, updated_at, deleted_at
        FROM users
        WHERE id = $1 AND deleted_at IS NULL
    `
	var user domain.User
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}
	return &user, nil
}
