package domain

import "time"

type TickedStatus string

const (
	TickedStatusBooked    TickedStatus = "booked"
	TickedStatusPurchased TickedStatus = "purchased"
	TickedStatusCancelled TickedStatus = "cancelled"
)

type Ticket struct {
	ID        int64        `json:"id"`
	UserID    int64        `json:"user_id"`
	EventID   int64        `json:"event_id"`
	Status    TickedStatus `json:"ticked_status"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	DeletedAt *time.Time   `json:"deleted_at,omitempty"`
}
