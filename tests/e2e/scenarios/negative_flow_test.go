//go:build e2e

package scenarios

import (
	"context"
	"testing"

	"github.com/felo/felo-backend/tests/e2e/harness"
)

func TestNegativeFlow_NoDriverAvailable_AfterRetryNoMatch(t *testing.T) {
	cfg, sut := loadE2EContext(t)
	requireSuite(t, "full-regression")
	ctx, cancel := context.WithTimeout(context.Background(), cfg.EventTimeout)
	defer cancel()

	ride, err := sut.Ride.RequestRide(ctx, harness.RequestRideRequest{
		CustomerID: "cust-active-002",
		Pickup: harness.Coordinate{Latitude: -6.300, Longitude: 106.900},
		Destination: harness.Coordinate{Latitude: -6.350, Longitude: 106.950},
		FareEstimate: 25000,
	})
	if err != nil {
		t.Fatalf("sut.Ride.RequestRide() error = %v", err)
	}

	if err := harness.Eventually(ctx, cfg.PollInterval, func() error {
		driverID, err := sut.Matching.GetAssignment(ctx, ride.ID)
		if err == nil && driverID != "" {
			return harness.ExpectEqual(driverID, "")
		}
		return nil
	}); err != nil {
		t.Fatalf("unexpected driver assignment observed: %v", err)
	}
}

func TestNegativeFlow_QRExpired_CannotBeScanned(t *testing.T) {
	_, _ = loadE2EContext(t)
	requireSuite(t, "full-regression")
	t.Skip("enable once ride-service exposes QR expiry test fixture or time control hook")
}

func TestNegativeFlow_PaymentFailure_PublishesFailedEvent(t *testing.T) {
	cfg, sut := loadE2EContext(t)
	requireSuite(t, "full-regression")
	ctx, cancel := context.WithTimeout(context.Background(), cfg.EventTimeout)
	defer cancel()

	ride, err := sut.Ride.RequestRide(ctx, harness.RequestRideRequest{
		CustomerID: "cust-low-balance-001",
		Pickup: harness.Coordinate{Latitude: -6.200, Longitude: 106.816},
		Destination: harness.Coordinate{Latitude: -6.214, Longitude: 106.845},
		FareEstimate: 950000,
	})
	if err != nil {
		t.Fatalf("sut.Ride.RequestRide() error = %v", err)
	}
	if _, err := sut.Ride.StartRide(ctx, ride.ID); err != nil {
		t.Fatalf("sut.Ride.StartRide() error = %v", err)
	}
	if _, err := sut.Ride.CompleteRide(ctx, ride.ID); err != nil {
		t.Fatalf("sut.Ride.CompleteRide() error = %v", err)
	}
	if _, err := sut.Events.WaitForEvent(ctx, "payment.failed.v1", ride.ID); err != nil {
		t.Fatalf("payment failed event not observed: %v", err)
	}
}

func TestNegativeFlow_DuplicateSettlement_DoesNotDoubleCreditDriver(t *testing.T) {
	_, _ = loadE2EContext(t)
	requireSuite(t, "full-regression")
	t.Skip("enable once wallet-service exposes deterministic duplicate event replay hook")
}
