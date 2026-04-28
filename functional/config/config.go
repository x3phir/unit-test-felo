//go:build functional

package config

import (
	"errors"
	"os"
	"time"
)

var ErrFunctionalDisabled = errors.New("functional tests disabled")

type Config struct {
	Suite            string
	RideGRPCAddr     string
	MatchingGRPCAddr string
	WalletGRPCAddr   string
	PaymentGRPCAddr  string
	LocationGRPCAddr string
	AuthJWT          string
	RabbitURL        string
	EventTimeout     time.Duration
	PollInterval     time.Duration
}

func Load() (Config, error) {
	if os.Getenv("FELO_FUNCTIONAL_ENABLED") != "1" {
		return Config{}, ErrFunctionalDisabled
	}

	suite := os.Getenv("FELO_TEST_SUITE")
	if suite == "" {
		suite = "critical-flow"
	}

	return Config{
		Suite:            suite,
		RideGRPCAddr:     getenv("FELO_RIDE_GRPC_ADDR", "127.0.0.1:50051"),
		MatchingGRPCAddr: getenv("FELO_MATCHING_GRPC_ADDR", "127.0.0.1:50052"),
		WalletGRPCAddr:   getenv("FELO_WALLET_GRPC_ADDR", "127.0.0.1:50053"),
		PaymentGRPCAddr:  getenv("FELO_PAYMENT_GRPC_ADDR", "127.0.0.1:50054"),
		LocationGRPCAddr: getenv("FELO_LOCATION_GRPC_ADDR", "127.0.0.1:50055"),
		AuthJWT:          getenv("FELO_AUTH_JWT", "demo-functional-token"),
		RabbitURL:        getenv("FELO_RABBIT_URL", "amqp://felo:felo@127.0.0.1:5672/"),
		EventTimeout:     20 * time.Second,
		PollInterval:     200 * time.Millisecond,
	}, nil
}

func getenv(key string, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
