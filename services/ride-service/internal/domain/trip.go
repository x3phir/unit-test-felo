package domain

import "time"

type TripState string

const (
	StatePricingInquiry TripState = "pricing_inquiry"
	StateMatching       TripState = "matching"
	StateEnRoute        TripState = "en_route"
	StateArrived        TripState = "arrived"
	StateOnRide         TripState = "on_ride"
	StateCompleted      TripState = "completed"
	StateCancelled      TripState = "cancelled"
	StateQRGenerated    TripState = "qr_generated"
)

type Coordinate struct {
	Latitude  float64
	Longitude float64
}

type Trip struct {
	ID           string
	CustomerID   string
	DriverID     string
	Pickup       Coordinate
	Destination  Coordinate
	FareEstimate int64
	State        TripState
	QRCode       string
	QRExpiresAt  time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Event struct {
	Name       string
	TripID     string
	OccurredAt time.Time
}
