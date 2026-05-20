package unit_test

import (
	"context"
	"errors"
	"testing"

	"github.com/felo/felo-backend/services/matching-service/internal/domain"
	"github.com/felo/felo-backend/services/matching-service/internal/service"
)

func TestMatchingService_AssignDriver_SelectsNearestDriver(t *testing.T) {
	availability := &availabilityFake{
		driversByCall: [][]domain.DriverCandidate{{
			{DriverID: "driver-2", DistanceKM: 2.2},
			{DriverID: "driver-1", DistanceKM: 1.1},
		}},
	}
	assignments := &assignmentRepoFake{}
	publisher := &matchingPublisherFake{}
	svc := service.NewMatchingService(availability, assignments, publisher)

	assignment, err := svc.AssignDriver(context.Background(), domain.MatchRequest{
		RideID:          "ride-1",
		Pickup:          domain.Coordinate{Latitude: -6.2, Longitude: 106.8},
		InitialRadiusKM: 1,
	})
	if err != nil {
		t.Fatalf("AssignDriver() error = %v", err)
	}

	if assignment.DriverID != "driver-1" {
		t.Fatalf("assignment.DriverID = %s, want driver-1", assignment.DriverID)
	}
	if len(publisher.events) != 1 || publisher.events[0].Name != "driver.matched.v1" {
		t.Fatalf("publisher.events = %#v, want driver.matched.v1", publisher.events)
	}
}

func TestMatchingService_AssignDriver_ExpandsRadiusUntilDriverFound(t *testing.T) {
	availability := &availabilityFake{
		driversByCall: [][]domain.DriverCandidate{
			{},
			{},
			{{DriverID: "driver-7", DistanceKM: 3.4}},
		},
	}
	svc := service.NewMatchingService(availability, &assignmentRepoFake{}, &matchingPublisherFake{})

	assignment, err := svc.AssignDriver(context.Background(), domain.MatchRequest{
		RideID:          "ride-1",
		Pickup:          domain.Coordinate{Latitude: -6.2, Longitude: 106.8},
		InitialRadiusKM: 1,
	})
	if err != nil {
		t.Fatalf("AssignDriver() error = %v", err)
	}

	if assignment.SearchRadiusKM <= 1 {
		t.Fatalf("assignment.SearchRadiusKM = %v, want expanded radius", assignment.SearchRadiusKM)
	}
	if availability.calls != 3 {
		t.Fatalf("availability.calls = %d, want 3", availability.calls)
	}
}

func TestMatchingService_AssignDriver_NoDriversAfterRetries_ReturnsError(t *testing.T) {
	availability := &availabilityFake{
		driversByCall: [][]domain.DriverCandidate{{}, {}, {}, {}},
	}
	svc := service.NewMatchingService(availability, &assignmentRepoFake{}, &matchingPublisherFake{})

	_, err := svc.AssignDriver(context.Background(), domain.MatchRequest{
		RideID:          "ride-1",
		Pickup:          domain.Coordinate{Latitude: -6.2, Longitude: 106.8},
		InitialRadiusKM: 1,
	})
	if !errors.Is(err, service.ErrNoDriversAvailable) {
		t.Fatalf("AssignDriver() error = %v, want ErrNoDriversAvailable", err)
	}
}

func TestMatchingService_AssignDriver_ZeroRadiusUsesDefaultRadius(t *testing.T) {
	availability := &availabilityFake{
		driversByCall: [][]domain.DriverCandidate{{{DriverID: "driver-1", DistanceKM: 0.9}}},
	}
	svc := service.NewMatchingService(availability, &assignmentRepoFake{}, &matchingPublisherFake{})

	assignment, err := svc.AssignDriver(context.Background(), domain.MatchRequest{
		RideID: "ride-1",
		Pickup: domain.Coordinate{Latitude: -6.2, Longitude: 106.8},
	})
	if err != nil {
		t.Fatalf("AssignDriver() error = %v", err)
	}
	if assignment.SearchRadiusKM != 1 {
		t.Fatalf("assignment.SearchRadiusKM = %v, want 1", assignment.SearchRadiusKM)
	}
}

type availabilityFake struct {
	driversByCall [][]domain.DriverCandidate
	calls         int
}

func (f *availabilityFake) FindAvailableDrivers(_ context.Context, _ domain.Coordinate, _ float64) ([]domain.DriverCandidate, error) {
	if f.calls >= len(f.driversByCall) {
		return nil, nil
	}
	drivers := f.driversByCall[f.calls]
	f.calls++
	return drivers, nil
}

type assignmentRepoFake struct {
	assignments []domain.Assignment
}

func (f *assignmentRepoFake) Save(_ context.Context, assignment domain.Assignment) error {
	f.assignments = append(f.assignments, assignment)
	return nil
}

type matchingPublisherFake struct {
	events []domain.Event
}

func (f *matchingPublisherFake) Publish(_ context.Context, event domain.Event) error {
	f.events = append(f.events, event)
	return nil
}
