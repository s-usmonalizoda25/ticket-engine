package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/s-usmonalizoda25/ticket-engine/internal/domain"
	"github.com/s-usmonalizoda25/ticket-engine/internal/usecase"
	"go.uber.org/zap"
)

type TicketHandler struct {
	ticketUC *usecase.TicketUseCase
	log      *zap.Logger
}

func NewTicketHandler(ticketUC *usecase.TicketUseCase, log *zap.Logger) *TicketHandler {
	return &TicketHandler{
		ticketUC: ticketUC,
		log:      log,
	}
}

func (h *TicketHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/tickets/buy", h.BuyTicket)
}

func (h *TicketHandler) BuyTicket(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req BuyTicketRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid request body"})
		return
	}

	ticket, err := h.ticketUC.BuyTicket(r.Context(), req.EventID, req.UserID)
	if err != nil {
		h.log.Error("Failed to process ticket purchase",
			zap.Int64("event_id", req.EventID),
			zap.Int64("user_id", req.UserID),
			zap.Error(err),
		)

		switch {
		case errors.Is(err, domain.ErrNotFound):
			h.writeJSON(w, http.StatusNotFound, ErrorResponse{Error: "user or event not found"})
		case errors.Is(err, domain.ErrNoTicketsLeft):
			h.writeJSON(w, http.StatusConflict, ErrorResponse{Error: "no tickets left for this event"})
		default:
			h.writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "internal server error"})
		}
		return
	}

	h.log.Info("Ticket purchased via HTTP API",
		zap.Int64("ticket_id", ticket.ID),
		zap.Int64("user_id", ticket.UserID),
	)

	h.writeJSON(w, http.StatusCreated, BuyTicketResponse{
		TicketID: ticket.ID,
		EventID:  ticket.EventID,
		UserID:   ticket.UserID,
		Status:   string(ticket.Status),
		Message:  "Ticket successfully purchased",
	})
}

func (h *TicketHandler) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}
