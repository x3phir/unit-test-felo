package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/felo/felo-backend/internal/demo/app"
	"github.com/felo/felo-backend/internal/demo/config"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg := config.Load()
	application, err := app.New(ctx, cfg)
	if err != nil {
		log.Fatalf("create app: %v", err)
	}

	log.Printf("FELO demo runtime started")
	log.Printf("ride-service: %s", cfg.RideAddr)
	log.Printf("matching-service: %s", cfg.MatchingAddr)
	log.Printf("wallet-service: %s", cfg.WalletAddr)
	log.Printf("payment-service: %s", cfg.PaymentAddr)
	log.Printf("location-service: %s", cfg.LocationAddr)

	if err := application.Run(ctx); err != nil {
		log.Fatalf("run app: %v", err)
	}
}
