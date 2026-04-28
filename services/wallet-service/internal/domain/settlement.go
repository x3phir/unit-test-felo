package domain

import "time"

type Settlement struct {
	IdempotencyKey string
	TripID         string
	DriverID       string
	Amount         int64
	Currency       string
}

type SettlementRecord struct {
	IdempotencyKey string
	TripID         string
	DriverID       string
	Amount         int64
	Currency       string
	BalanceAfter   int64
	ProcessedAt    time.Time
}
