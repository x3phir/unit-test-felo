package domain

import "time"

type ShipmentStatus string

const (
	StatusCreated   ShipmentStatus = "created"
	StatusPickedUp  ShipmentStatus = "picked_up"
	StatusInTransit ShipmentStatus = "in_transit"
	StatusDelivered ShipmentStatus = "delivered"
)

type ProofOfDelivery struct {
	PhotoURL  string
	Signature string
}

type Shipment struct {
	ID              string
	SendOrderID     string
	DriverID        string
	Status          ShipmentStatus
	TrackingNumber  string
	ProofOfDelivery ProofOfDelivery
	UpdatedAt       time.Time
}

type Event struct {
	Name       string
	ShipmentID string
	OccurredAt time.Time
}
