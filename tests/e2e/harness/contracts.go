//go:build e2e

package harness

import (
	"context"
	"time"
)

type Coordinate struct {
	Latitude  float64
	Longitude float64
}

type Ride struct {
	ID     string
	State  string
	Driver string
}

type QRSession struct {
	TripID     string
	QRCode     string
	ExpiresAt  time.Time
	DriverLock string
}

type WalletBalance struct {
	OwnerID  string
	Balance  int64
	Currency string
}

type LocationSample struct {
	DriverID   string
	Position   Coordinate
	RecordedAt time.Time
}

type Event struct {
	Name    string
	Key     string
	Payload map[string]string
}

type RideClient interface {
	RequestRide(ctx context.Context, req RequestRideRequest) (Ride, error)
	StartRide(ctx context.Context, rideID string) (Ride, error)
	CompleteRide(ctx context.Context, rideID string) (Ride, error)
	GenerateNowQR(ctx context.Context, req GenerateNowQRRequest) (QRSession, error)
	ScanNowQR(ctx context.Context, qrCode string, driverID string) (QRSession, error)
	AcceptNowQR(ctx context.Context, tripID string, driverID string) (Ride, error)
}

type MatchingClient interface {
	GetAssignment(ctx context.Context, rideID string) (string, error)
}

type WalletClient interface {
	GetDriverBalance(ctx context.Context, driverID string) (WalletBalance, error)
	GetCustomerBalance(ctx context.Context, customerID string) (WalletBalance, error)
}

type PaymentClient interface {
	GetPaymentStatus(ctx context.Context, rideID string) (string, error)
}

type LocationClient interface {
	ReportDriverLocation(ctx context.Context, sample LocationSample) error
	GetDriverHistory(ctx context.Context, driverID string, from time.Time, to time.Time) ([]LocationSample, error)
}

type EventObserver interface {
	WaitForEvent(ctx context.Context, name string, key string) (Event, error)
}

type HealthChecker interface {
	Check(ctx context.Context) error
}

type RequestRideRequest struct {
	CustomerID   string
	Pickup       Coordinate
	Destination  Coordinate
	FareEstimate int64
}

type GenerateNowQRRequest struct {
	CustomerID   string
	Destination  Coordinate
	FareEstimate int64
}
