package http

type BuyTicketRequest struct {
	EventID int64 `json:"event_id"`
	UserID  int64 `json:"user_id"`
}

type BuyTicketResponse struct {
	TicketID int64  `json:"ticket_id"`
	EventID  int64  `json:"event_id"`
	UserID   int64  `json:"user_id"`
	Status   string `json:"status"`
	Message  string `json:"message"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
