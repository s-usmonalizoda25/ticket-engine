package main

import (
	"context"
	"net/http"
	"time"

	"github.com/s-usmonalizoda25/ticket-engine/config"
	deliveryHTTP "github.com/s-usmonalizoda25/ticket-engine/internal/delivery/http"
	"github.com/s-usmonalizoda25/ticket-engine/internal/repository/postgres"
	"github.com/s-usmonalizoda25/ticket-engine/internal/usecase"
	"github.com/s-usmonalizoda25/ticket-engine/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	log, err := logger.NewLogger()
	if err != nil {
		panic(err)
	}
	defer log.Sync()

	log.Info("Starting Ticket Engine service...")

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config", zap.Error(err))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := postgres.NewPostgresDB(ctx, cfg.DBConnString)
	if err != nil {
		log.Fatal("Failed to connect to PostgreSQL pool", zap.Error(err))
	}
	defer pool.Close()

	log.Info("Successfully connected to PostgreSQL connection pool!")

	userRepo := postgres.NewUserRepository(pool)
	eventRepo := postgres.NewEventRepository(pool)
	ticketRepo := postgres.NewTicketRepository(pool)

	ticketUC := usecase.NewTicketUseCase(ticketRepo, userRepo, eventRepo)

	mux := http.NewServeMux()
	ticketHandler := deliveryHTTP.NewTicketHandler(ticketUC, log)
	ticketHandler.RegisterRoutes(mux)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	log.Info("HTTP server running on port :8080")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("HTTP server failed", zap.Error(err))
	}
}
