package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/s-usmonalizoda25/ticket-engine/internal/domain"
)

type EventRepository struct {
	pool *pgxpool.Pool
}

func NewEventRepository(pool *pgxpool.Pool) *EventRepository {
	return &EventRepository{pool: pool}
}

func (r *EventRepository) Create(ctx context.Context, event *domain.Event) error {
	const query = `
        INSERT INTO events (title, total_tickets, tickets_left, price, created_at, updated_at)
        VALUES ($1, $2, $3, $4, NOW(), NOW())
        RETURNING id, created_at, updated_at
    `
	err := r.pool.QueryRow(
		ctx,
		query,
		event.Title,
		event.TotalTickets,
		event.TicketsLeft,
		event.Price,
	).Scan(&event.ID, &event.CreatedAt, &event.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create event: %w", err)
	}
	return nil
}

func (r *EventRepository) GetByID(ctx context.Context, id int64) (*domain.Event, error) {
	const query = `
        SELECT id, title, total_tickets, tickets_left, price, created_at, updated_at, deleted_at
        FROM events
        WHERE id = $1 AND deleted_at IS NULL
    `
	var event domain.Event
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&event.ID,
		&event.Title,
		&event.TotalTickets,
		&event.TicketsLeft,
		&event.Price,
		&event.CreatedAt,
		&event.UpdatedAt,
		&event.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get event by id: %w", err)
	}
	return &event, nil
}
