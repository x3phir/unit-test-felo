package domain

import "time"

type TrackingStatus string

const (
	StatusActive    TrackingStatus = "active"
	StatusPaused    TrackingStatus = "paused"
	StatusCompleted TrackingStatus = "completed"
)

type Coordinate struct {
	Latitude  float64
	Longitude float64
}

type TrackingSession struct {
	ID         string
	ShipmentID string
	DriverID   string
	Status     TrackingStatus
	StartedAt  time.Time
	UpdatedAt  time.Time
	EndedAt    *time.Time
}

type TrackingRecord struct {
	ID         string
	SessionID  string
	Coordinate Coordinate
	Speed      float64
	Heading    float64
	RecordedAt time.Time
}

type Event struct {
	Name       string
	SessionID  string
	OccurredAt time.Time
}
