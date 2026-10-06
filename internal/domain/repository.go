package domain

import "context"

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id int64) (*User, error)
}

type EventRepository interface {
	Create(ctx context.Context, event *Event) error
	GetByID(ctx context.Context, id int64) (*Event, error)
}

type TicketRepository interface {
	Create(ctx context.Context, ticket *Ticket) error
	GetByID(ctx context.Context, id int64) (*Ticket, error)
	BuyTicketTx(ctx context.Context, eventID, userID int64) (*Ticket, error)
}
