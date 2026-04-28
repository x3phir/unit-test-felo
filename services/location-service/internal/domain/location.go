package domain

import "time"

type Coordinate struct {
	Latitude  float64
	Longitude float64
}

type LocationSample struct {
	DriverID   string
	Position   Coordinate
	RecordedAt time.Time
}

type RouteRequest struct {
	Origin      Coordinate
	Destination Coordinate
}

type RouteEstimate struct {
	DistanceMeters  int
	DurationSeconds int
	Polyline        string
	CalculatedAt    time.Time
}
