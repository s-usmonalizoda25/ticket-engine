package domain

import "time"

type Event struct {
	ID           int64      `json:"id"`
	Title        string     `json:"title"`
	TotalTickets int        `json:"total_tickets"`
	TicketsLeft  int        `json:"tickets_left"`
	Price        float64    `json:"price"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty"`
}