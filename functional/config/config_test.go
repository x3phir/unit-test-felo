//go:build functional

package config

import (
	"os"
	"testing"
)

func TestLoad_DisabledReturnsError(t *testing.T) {
	t.Setenv("FELO_FUNCTIONAL_ENABLED", "")

	_, err := Load()
	if err == nil {
		t.Fatal("Load() error = nil, want error")
	}
}

func TestLoad_DefaultsSuiteWhenNotProvided(t *testing.T) {
	t.Setenv("FELO_FUNCTIONAL_ENABLED", "1")
	t.Setenv("FELO_TEST_SUITE", "")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.Suite != "critical-flow" {
		t.Fatalf("cfg.Suite = %s, want critical-flow", cfg.Suite)
	}
}

func TestLoad_ReadsConfiguredAddresses(t *testing.T) {
	t.Setenv("FELO_FUNCTIONAL_ENABLED", "1")
	t.Setenv("FELO_TEST_SUITE", "smoke")
	t.Setenv("FELO_RIDE_GRPC_ADDR", "127.0.0.1:50051")
	t.Setenv("FELO_MATCHING_GRPC_ADDR", "127.0.0.1:50052")
	t.Setenv("FELO_WALLET_GRPC_ADDR", "127.0.0.1:50053")
	t.Setenv("FELO_PAYMENT_GRPC_ADDR", "127.0.0.1:50054")
	t.Setenv("FELO_LOCATION_GRPC_ADDR", "127.0.0.1:50055")
	t.Setenv("FELO_AUTH_JWT", "jwt")
	t.Setenv("FELO_RABBIT_URL", "amqp://guest:guest@127.0.0.1:5672/")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.RideGRPCAddr != "127.0.0.1:50051" {
		t.Fatalf("cfg.RideGRPCAddr = %s, want 127.0.0.1:50051", cfg.RideGRPCAddr)
	}
	if os.Getenv("FELO_AUTH_JWT") == "" {
		t.Fatal("FELO_AUTH_JWT should be set")
	}
	if cfg.RabbitURL == "" {
		t.Fatal("cfg.RabbitURL should not be empty")
	}
}
