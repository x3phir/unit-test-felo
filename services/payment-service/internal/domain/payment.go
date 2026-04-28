package domain

import "time"

type RideCompletedEvent struct {
	EventID    string
	TripID     string
	CustomerID string
	Amount     int64
	Currency   string
}

type PaymentResult struct {
	EventID   string
	TripID    string
	InvoiceID string
	PaidAt    time.Time
}

type Event struct {
	Name       string
	TripID     string
	OccurredAt time.Time
}
