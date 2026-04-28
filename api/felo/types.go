package felo

type Coordinate struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Ride struct {
	ID     string `json:"id"`
	State  string `json:"state"`
	Driver string `json:"driver,omitempty"`
}

type QRSession struct {
	TripID      string `json:"trip_id"`
	QRCode      string `json:"qr_code"`
	ExpiresAt   string `json:"expires_at"`
	LockedDriver string `json:"locked_driver,omitempty"`
}

type WalletBalance struct {
	OwnerID  string `json:"owner_id"`
	Balance  int64  `json:"balance"`
	Currency string `json:"currency"`
}

type LocationSample struct {
	DriverID   string     `json:"driver_id"`
	Position   Coordinate `json:"position"`
	RecordedAt string     `json:"recorded_at"`
}

type EventEnvelope struct {
	Name    string            `json:"name"`
	Key     string            `json:"key"`
	Payload map[string]string `json:"payload"`
}
