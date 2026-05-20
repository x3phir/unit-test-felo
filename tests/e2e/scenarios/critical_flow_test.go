//go:build e2e

package scenarios

import (
	"context"
	"testing"
	"time"

	"github.com/felo/felo-backend/tests/e2e/harness"
)

func TestCriticalFlow_RegularRideToSettlement(t *testing.T) {
	cfg, sut := loadE2EContext(t)
	requireSuite(t, "critical-flow", "full-regression")
	ctx, cancel := context.WithTimeout(context.Background(), cfg.EventTimeout)
	defer cancel()

	ride, err := sut.Ride.RequestRide(ctx, harness.RequestRideRequest{
		CustomerID: "cust-active-001",
		Pickup: harness.Coordinate{Latitude: -6.200, Longitude: 106.816},
		Destination: harness.Coordinate{Latitude: -6.214, Longitude: 106.845},
		FareEstimate: 25000,
	})
	if err != nil {
		t.Fatalf("sut.Ride.RequestRide() error = %v", err)
	}

	if err := harness.Eventually(ctx, cfg.PollInterval, func() error {
		driverID, err := sut.Matching.GetAssignment(ctx, ride.ID)
		if err != nil {
			return err
		}
		return harness.ExpectEqual(driverID, "driver-active-001")
	}); err != nil {
		t.Fatalf("matching assignment check failed: %v", err)
	}

	if err := sut.Location.ReportDriverLocation(ctx, harness.LocationSample{
		DriverID: "driver-active-001",
		Position: harness.Coordinate{Latitude: -6.201, Longitude: 106.818},
		RecordedAt: time.Now(),
	}); err != nil {
		t.Fatalf("sut.Location.ReportDriverLocation() error = %v", err)
	}

	if _, err := sut.Ride.StartRide(ctx, ride.ID); err != nil {
		t.Fatalf("sut.Ride.StartRide() error = %v", err)
	}
	if _, err := sut.Ride.CompleteRide(ctx, ride.ID); err != nil {
		t.Fatalf("sut.Ride.CompleteRide() error = %v", err)
	}
	if err := harness.Eventually(ctx, cfg.PollInterval, func() error {
		status, err := sut.Payment.GetPaymentStatus(ctx, ride.ID)
		if err != nil {
			return err
		}
		return harness.ExpectEqual(status, "completed")
	}); err != nil {
		t.Fatalf("payment status check failed: %v", err)
	}
	if _, err := sut.Events.WaitForEvent(ctx, "ride.completed.v1", ride.ID); err != nil {
		t.Fatalf("ride completed event not observed: %v", err)
	}
	if _, err := sut.Events.WaitForEvent(ctx, "payment.completed.v1", ride.ID); err != nil {
		t.Fatalf("payment completed event not observed: %v", err)
	}
}

func TestCriticalFlow_FeloNowQRRideToSettlement(t *testing.T) {
	cfg, sut := loadE2EContext(t)
	requireSuite(t, "critical-flow", "full-regression")
	ctx, cancel := context.WithTimeout(context.Background(), cfg.EventTimeout)
	defer cancel()

	qr, err := sut.Ride.GenerateNowQR(ctx, harness.GenerateNowQRRequest{
		CustomerID: "cust-active-001",
		Destination: harness.Coordinate{Latitude: -6.214, Longitude: 106.845},
		FareEstimate: 18000,
	})
	if err != nil {
		t.Fatalf("sut.Ride.GenerateNowQR() error = %v", err)
	}
	if _, err := sut.Ride.ScanNowQR(ctx, qr.QRCode, "driver-active-001"); err != nil {
		t.Fatalf("sut.Ride.ScanNowQR() error = %v", err)
	}
	ride, err := sut.Ride.AcceptNowQR(ctx, qr.TripID, "driver-active-001")
	if err != nil {
		t.Fatalf("sut.Ride.AcceptNowQR() error = %v", err)
	}
	if _, err := sut.Ride.CompleteRide(ctx, ride.ID); err != nil {
		t.Fatalf("sut.Ride.CompleteRide() error = %v", err)
	}
	if _, err := sut.Events.WaitForEvent(ctx, "payment.completed.v1", ride.ID); err != nil {
		t.Fatalf("payment completed event not observed: %v", err)
	}
}
