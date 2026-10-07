package main

import (
	"context"
	"time"

	"github.com/s-usmonalizoda25/ticket-engine/config"
	"github.com/s-usmonalizoda25/ticket-engine/internal/repository/postgres"
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

	_ = pool
}
