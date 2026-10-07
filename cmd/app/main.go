package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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

	initCtx, initCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer initCancel()

	pool, err := postgres.NewPostgresDB(initCtx, cfg.DBConnString)
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

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Info("HTTP server running on port :8080")
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal("HTTP server failed", zap.Error(err))
		}
	}()

	<-stop
	log.Info("Shutting down HTTP server gracefully...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Error("HTTP server forced to shutdown", zap.Error(err))
	} else {
		log.Info("HTTP server stopped cleanly")
	}

	log.Info("Ticket Engine stopped successfully")
}
