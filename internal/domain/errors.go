package domain

import "errors"

var (
	ErrNotFound          = errors.New("resource not found")
	ErrNoTicketsLeft     = errors.New("no tickets left for this event")
	ErrTicketAlreadyPaid = errors.New("ticket is already purchased")
)
