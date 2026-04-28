package config

import "os"

type Config struct {
	RideAddr         string
	MatchingAddr     string
	WalletAddr       string
	PaymentAddr      string
	LocationAddr     string
	RidePostgresDSN  string
	MatchPostgresDSN string
	WalletPostgresDSN string
	PaymentPostgresDSN string
	LocationPostgresDSN string
	RedisAddr        string
	RabbitURL        string
	AuthToken        string
}

func Load() Config {
	return Config{
		RideAddr:            env("FELO_RIDE_GRPC_ADDR", "127.0.0.1:50051"),
		MatchingAddr:        env("FELO_MATCHING_GRPC_ADDR", "127.0.0.1:50052"),
		WalletAddr:          env("FELO_WALLET_GRPC_ADDR", "127.0.0.1:50053"),
		PaymentAddr:         env("FELO_PAYMENT_GRPC_ADDR", "127.0.0.1:50054"),
		LocationAddr:        env("FELO_LOCATION_GRPC_ADDR", "127.0.0.1:50055"),
		RidePostgresDSN:     env("FELO_RIDE_PG_DSN", "postgres://felo:felo@127.0.0.1:54321/ride_db?sslmode=disable"),
		MatchPostgresDSN:    env("FELO_MATCHING_PG_DSN", "postgres://felo:felo@127.0.0.1:54322/matching_db?sslmode=disable"),
		WalletPostgresDSN:   env("FELO_WALLET_PG_DSN", "postgres://felo:felo@127.0.0.1:54323/wallet_db?sslmode=disable"),
		PaymentPostgresDSN:  env("FELO_PAYMENT_PG_DSN", "postgres://felo:felo@127.0.0.1:54324/payment_db?sslmode=disable"),
		LocationPostgresDSN: env("FELO_LOCATION_PG_DSN", "postgres://felo:felo@127.0.0.1:54325/location_db?sslmode=disable"),
		RedisAddr:           env("FELO_REDIS_ADDR", "127.0.0.1:6379"),
		RabbitURL:           env("FELO_RABBIT_URL", "amqp://felo:felo@127.0.0.1:5672/"),
		AuthToken:           authToken(),
	}
}

func env(key string, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func authToken() string {
	if value := os.Getenv("FELO_AUTH_TOKEN"); value != "" {
		return value
	}
	if value := os.Getenv("FELO_AUTH_JWT"); value != "" {
		return value
	}
	return "demo-functional-token"
}
