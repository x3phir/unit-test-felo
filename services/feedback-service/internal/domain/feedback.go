package domain

import "time"

type Feedback struct {
	ID        string
	TripID    string
	UserID    string
	TargetID  string // DriverID or MerchantID
	Score     int    // 1 to 5
	Comment   string
	CreatedAt time.Time
}

type Event struct {
	Name       string
	ReviewID   string
	TargetID   string
	Score      int
	OccurredAt time.Time
}
