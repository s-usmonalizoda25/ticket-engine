package usecase

import (
	"context"
	"fmt"

	"github.com/s-usmonalizoda25/ticket-engine/internal/domain"
)

type TicketUseCase struct {
	ticketRepo domain.TicketRepository
	userRepo   domain.UserRepository
	eventRepo  domain.EventRepository
}

func NewTicketUseCase(
	tr domain.TicketRepository,
	ur domain.UserRepository,
	er domain.EventRepository,
) *TicketUseCase {
	return &TicketUseCase{
		ticketRepo: tr,
		userRepo:   ur,
		eventRepo:  er,
	}
}

func (uc *TicketUseCase) BuyTicket(ctx context.Context, eventID, userID int64) (*domain.Ticket, error) {
	if eventID <= 0 || userID <= 0 {
		return nil, fmt.Errorf("invalid event_id or user_id")
	}

	_, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user check failed: %w", err)
	}

	_, err = uc.eventRepo.GetByID(ctx, eventID)
	if err != nil {
		return nil, fmt.Errorf("event check failed: %w", err)
	}

	ticket, err := uc.ticketRepo.BuyTicketTx(ctx, eventID, userID)
	if err != nil {
		return nil, fmt.Errorf("purchase failed: %w", err)
	}

	return ticket, nil
}
