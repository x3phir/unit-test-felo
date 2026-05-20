package unit_test

import (
	"context"
	"testing"
	"time"

	"github.com/felo/felo-backend/services/ride-service/internal/domain"
	"github.com/felo/felo-backend/services/ride-service/internal/service"
	"go.uber.org/mock/gomock"
)

func TestTripService_RequestRide_UsesRepositoryAndPublishesEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := NewMockTripRepository(ctrl)
	publisher := NewMockEventPublisher(ctrl)
	now := time.Date(2026, 5, 5, 10, 0, 0, 0, time.UTC)
	svc := service.NewTripService(repo, publisher, fixedClock{now: now}, &idGen{ids: []string{"ride-001"}})

	repo.EXPECT().Save(gomock.Any(), gomock.AssignableToTypeOf(domain.Trip{})).Return(nil)
	publisher.EXPECT().Publish(gomock.Any(), gomock.AssignableToTypeOf(domain.Event{})).DoAndReturn(
		func(_ context.Context, event domain.Event) error {
			if event.Name != "ride.created.v1" {
				t.Fatalf("event.Name = %s, want ride.created.v1", event.Name)
			}
			return nil
		},
	)

	trip, err := svc.RequestRide(context.Background(), service.RequestRideInput{
		CustomerID:   "cust-1",
		Pickup:       domain.Coordinate{Latitude: -6.2, Longitude: 106.8},
		Destination:  domain.Coordinate{Latitude: -6.3, Longitude: 106.9},
		FareEstimate: 25000,
	})
	if err != nil {
		t.Fatalf("RequestRide() error = %v", err)
	}
	if trip.State != domain.StateMatching {
		t.Fatalf("trip.State = %s, want %s", trip.State, domain.StateMatching)
	}
}

func TestTripService_StartRide_ValidTransition_UpdatesStateAndPublishes(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := NewMockTripRepository(ctrl)
	publisher := NewMockEventPublisher(ctrl)
	now := time.Date(2026, 5, 5, 10, 10, 0, 0, time.UTC)
	svc := service.NewTripService(repo, publisher, fixedClock{now: now}, &idGen{ids: []string{"ride-001"}})

	repo.EXPECT().GetByID(gomock.Any(), "ride-001").Return(domain.Trip{ID: "ride-001", State: domain.StateArrived}, nil)
	repo.EXPECT().Save(gomock.Any(), gomock.AssignableToTypeOf(domain.Trip{})).Return(nil)
	publisher.EXPECT().Publish(gomock.Any(), gomock.AssignableToTypeOf(domain.Event{})).Return(nil)

	trip, err := svc.StartRide(context.Background(), "ride-001")
	if err != nil {
		t.Fatalf("StartRide() error = %v", err)
	}
	if trip.State != domain.StateOnRide {
		t.Fatalf("trip.State = %s, want %s", trip.State, domain.StateOnRide)
	}
}
