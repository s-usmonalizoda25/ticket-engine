package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/s-usmonalizoda25/ticket-engine/internal/domain"
)

type TicketRepository struct {
	pool *pgxpool.Pool
}

func NewTicketRepository(pool *pgxpool.Pool) *TicketRepository {
	return &TicketRepository{pool: pool}
}

func (r *TicketRepository) Create(ctx context.Context, ticket *domain.Ticket) error {
	const query = `
		INSERT INTO tickets (event_id, user_id, status, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`
	err := r.pool.QueryRow(ctx, query, ticket.EventID, ticket.UserID, ticket.Status).Scan(
		&ticket.ID, &ticket.CreatedAt, &ticket.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create ticket: %w", err)
	}
	return nil
}

func (r *TicketRepository) GetByID(ctx context.Context, id int64) (*domain.Ticket, error) {
	const query = `
		SELECT id, event_id, user_id, status, created_at, updated_at, deleted_at
		FROM tickets
		WHERE id = $1 AND deleted_at IS NULL
	`
	var ticket domain.Ticket

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&ticket.ID,
		&ticket.EventID,
		&ticket.UserID,
		&ticket.Status,
		&ticket.CreatedAt,
		&ticket.UpdatedAt,
		&ticket.DeletedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get ticket by id: %w", err)
	}
	return &ticket, nil
}

func (r *TicketRepository) BuyTicketTx(ctx context.Context, eventID, userID int64) (*domain.Ticket, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer tx.Rollback(ctx)

	var ticketsLeft int
	const selectQuery = `
		SELECT tickets_left
		FROM events
		WHERE id = $1 AND deleted_at IS NULL
		FOR UPDATE
	`

	err = tx.QueryRow(ctx, selectQuery, eventID).Scan(&ticketsLeft)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to lock event row: %w", err)
	}

	if ticketsLeft <= 0 {
		return nil, domain.ErrNoTicketsLeft
	}

	const updateQuery = `
		UPDATE events
		SET tickets_left = tickets_left - 1, updated_at = NOW()
		WHERE id = $1
	`

	_, err = tx.Exec(ctx, updateQuery, eventID)
	if err != nil {
		return nil, fmt.Errorf("failed to decrement tickets_left: %w", err)
	}

	const insertQuery = `
		INSERT INTO tickets (event_id, user_id, status, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`

	ticket := &domain.Ticket{
		EventID: eventID,
		UserID:  userID,
		Status:  domain.TicketStatusPurchased,
	}

	err = tx.QueryRow(ctx, insertQuery, eventID, userID, ticket.Status).Scan(
		&ticket.ID,
		&ticket.CreatedAt,
		&ticket.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to insert ticket in tx: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return ticket, nil
}
