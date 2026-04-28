package domain

import "time"

type Coordinate struct {
	Latitude  float64
	Longitude float64
}

type MatchRequest struct {
	RideID           string
	Pickup           Coordinate
	InitialRadiusKM  float64
}

type DriverCandidate struct {
	DriverID   string
	DistanceKM float64
}

type Assignment struct {
	RideID          string
	DriverID        string
	SearchRadiusKM  float64
	AssignedAt      time.Time
}

type Event struct {
	Name       string
	RideID     string
	DriverID   string
	OccurredAt time.Time
}
