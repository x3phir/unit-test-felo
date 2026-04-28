package domain

import "time"

type Coordinate struct {
	Latitude  float64
	Longitude float64
}

type SavedAddress struct {
	Name       string
	Coordinate Coordinate
}

type UserProfile struct {
	ID             string
	Name           string
	Phone          string
	Email          string
	PhotoURL       string
	Locale         string
	SavedAddresses []SavedAddress
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
