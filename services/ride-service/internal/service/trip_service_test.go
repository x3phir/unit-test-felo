package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/felo/felo-backend/services/ride-service/internal/domain"
	"github.com/felo/felo-backend/services/ride-service/internal/service"
)

func TestTripService_RequestRide_PublishesRideCreatedEvent(t *testing.T) {
	now := time.Date(2026, 4, 28, 9, 0, 0, 0, time.UTC)
	repo := &rideRepoFake{trips: map[string]domain.Trip{}}
	publisher := &ridePublisherFake{}
	svc := service.NewTripService(repo, publisher, fixedClock{now: now}, sequenceIDs("trip-001"))

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
	if len(publisher.events) != 1 || publisher.events[0].Name != "ride.created.v1" {
		t.Fatalf("publisher.events = %#v, want ride.created.v1", publisher.events)
	}
}

func TestTripService_RequestRide_InvalidInput_ReturnsError(t *testing.T) {
	svc := service.NewTripService(&rideRepoFake{trips: map[string]domain.Trip{}}, &ridePublisherFake{}, fixedClock{now: time.Now()}, sequenceIDs("trip-001"))

	_, err := svc.RequestRide(context.Background(), service.RequestRideInput{})
	if !errors.Is(err, service.ErrInvalidInput) {
		t.Fatalf("RequestRide() error = %v, want ErrInvalidInput", err)
	}
}

func TestTripService_StartRide_InvalidState_ReturnsError(t *testing.T) {
	repo := &rideRepoFake{trips: map[string]domain.Trip{
		"trip-001": {ID: "trip-001", State: domain.StateMatching},
	}}
	svc := service.NewTripService(repo, &ridePublisherFake{}, fixedClock{now: time.Now()}, sequenceIDs("trip-001"))

	_, err := svc.StartRide(context.Background(), "trip-001")
	if !errors.Is(err, service.ErrInvalidTransition) {
		t.Fatalf("StartRide() error = %v, want ErrInvalidTransition", err)
	}
}

func TestTripService_StartRide_ArrivedTrip_PublishesRideStartedEvent(t *testing.T) {
	now := time.Date(2026, 4, 28, 9, 30, 0, 0, time.UTC)
	repo := &rideRepoFake{trips: map[string]domain.Trip{
		"trip-001": {ID: "trip-001", State: domain.StateArrived},
	}}
	publisher := &ridePublisherFake{}
	svc := service.NewTripService(repo, publisher, fixedClock{now: now}, sequenceIDs("trip-001"))

	trip, err := svc.StartRide(context.Background(), "trip-001")
	if err != nil {
		t.Fatalf("StartRide() error = %v", err)
	}

	if trip.State != domain.StateOnRide {
		t.Fatalf("trip.State = %s, want %s", trip.State, domain.StateOnRide)
	}
	if len(publisher.events) != 1 || publisher.events[0].Name != "ride.started.v1" {
		t.Fatalf("publisher.events = %#v, want ride.started.v1", publisher.events)
	}
}

func TestTripService_CompleteRide_PublishesRideCompletedEvent(t *testing.T) {
	now := time.Date(2026, 4, 28, 10, 0, 0, 0, time.UTC)
	repo := &rideRepoFake{trips: map[string]domain.Trip{
		"trip-001": {ID: "trip-001", State: domain.StateOnRide},
	}}
	publisher := &ridePublisherFake{}
	svc := service.NewTripService(repo, publisher, fixedClock{now: now}, sequenceIDs("trip-001"))

	trip, err := svc.CompleteRide(context.Background(), "trip-001")
	if err != nil {
		t.Fatalf("CompleteRide() error = %v", err)
	}

	if trip.State != domain.StateCompleted {
		t.Fatalf("trip.State = %s, want %s", trip.State, domain.StateCompleted)
	}
	if len(publisher.events) != 1 || publisher.events[0].Name != "ride.completed.v1" {
		t.Fatalf("publisher.events = %#v, want ride.completed.v1", publisher.events)
	}
}

func TestTripService_CompleteRide_InvalidState_ReturnsError(t *testing.T) {
	repo := &rideRepoFake{trips: map[string]domain.Trip{
		"trip-001": {ID: "trip-001", State: domain.StateArrived},
	}}
	svc := service.NewTripService(repo, &ridePublisherFake{}, fixedClock{now: time.Now()}, sequenceIDs("trip-001"))

	_, err := svc.CompleteRide(context.Background(), "trip-001")
	if !errors.Is(err, service.ErrInvalidTransition) {
		t.Fatalf("CompleteRide() error = %v, want ErrInvalidTransition", err)
	}
}

func TestTripService_GenerateNowQR_SetsTenMinuteExpiry(t *testing.T) {
	now := time.Date(2026, 4, 28, 11, 0, 0, 0, time.UTC)
	repo := &rideRepoFake{trips: map[string]domain.Trip{}}
	svc := service.NewTripService(repo, &ridePublisherFake{}, fixedClock{now: now}, sequenceIDs("trip-001", "qr-001"))

	trip, err := svc.GenerateNowQR(context.Background(), service.GenerateNowQRInput{
		CustomerID:   "cust-1",
		Destination:  domain.Coordinate{Latitude: -6.3, Longitude: 106.9},
		FareEstimate: 18000,
	})
	if err != nil {
		t.Fatalf("GenerateNowQR() error = %v", err)
	}

	if trip.State != domain.StateQRGenerated {
		t.Fatalf("trip.State = %s, want %s", trip.State, domain.StateQRGenerated)
	}
	if !trip.QRExpiresAt.Equal(now.Add(10 * time.Minute)) {
		t.Fatalf("trip.QRExpiresAt = %v, want %v", trip.QRExpiresAt, now.Add(10*time.Minute))
	}
}

func TestTripService_GenerateNowQR_InvalidInput_ReturnsError(t *testing.T) {
	svc := service.NewTripService(&rideRepoFake{trips: map[string]domain.Trip{}}, &ridePublisherFake{}, fixedClock{now: time.Now()}, sequenceIDs("trip-001", "qr-001"))

	_, err := svc.GenerateNowQR(context.Background(), service.GenerateNowQRInput{})
	if !errors.Is(err, service.ErrInvalidInput) {
		t.Fatalf("GenerateNowQR() error = %v, want ErrInvalidInput", err)
	}
}

type rideRepoFake struct {
	trips map[string]domain.Trip
}

func (r *rideRepoFake) Save(_ context.Context, trip domain.Trip) error {
	r.trips[trip.ID] = trip
	return nil
}

func (r *rideRepoFake) GetByID(_ context.Context, tripID string) (domain.Trip, error) {
	trip, ok := r.trips[tripID]
	if !ok {
		return domain.Trip{}, errors.New("trip not found")
	}
	return trip, nil
}

type ridePublisherFake struct {
	events []domain.Event
}

func (p *ridePublisherFake) Publish(_ context.Context, event domain.Event) error {
	p.events = append(p.events, event)
	return nil
}

type fixedClock struct {
	now time.Time
}

func (c fixedClock) Now() time.Time {
	return c.now
}

type sequenceIDGenerator struct {
	values []string
	index  int
}

func sequenceIDs(values ...string) *sequenceIDGenerator {
	return &sequenceIDGenerator{values: values}
}

func (g *sequenceIDGenerator) NewID() string {
	if g.index >= len(g.values) {
		return "generated-id"
	}
	value := g.values[g.index]
	g.index++
	return value
}
