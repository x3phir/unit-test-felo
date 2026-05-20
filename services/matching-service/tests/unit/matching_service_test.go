package unit_test

import (
	"context"
	"testing"

	"github.com/felo/felo-backend/services/matching-service/internal/domain"
	"github.com/felo/felo-backend/services/matching-service/internal/service"
	"go.uber.org/mock/gomock"
)

func TestMatchingService_AssignDriver_SelectsNearestDriverWithGomock(t *testing.T) {
	ctrl := gomock.NewController(t)
	availability := NewMockAvailabilityReader(ctrl)
	assignments := NewMockAssignmentRepository(ctrl)
	publisher := NewMockEventPublisher(ctrl)
	svc := service.NewMatchingService(availability, assignments, publisher)

	availability.EXPECT().FindAvailableDrivers(gomock.Any(), gomock.Any(), 1.0).Return([]domain.DriverCandidate{
		{DriverID: "driver-2", DistanceKM: 2.4},
		{DriverID: "driver-1", DistanceKM: 1.1},
	}, nil)
	assignments.EXPECT().Save(gomock.Any(), gomock.AssignableToTypeOf(domain.Assignment{})).Return(nil)
	publisher.EXPECT().Publish(gomock.Any(), gomock.AssignableToTypeOf(domain.Event{})).Return(nil)

	assignment, err := svc.AssignDriver(context.Background(), domain.MatchRequest{
		RideID: "ride-1",
		Pickup: domain.Coordinate{Latitude: -6.2, Longitude: 106.8},
		InitialRadiusKM: 1,
	})
	if err != nil {
		t.Fatalf("AssignDriver() error = %v", err)
	}
	if assignment.DriverID != "driver-1" {
		t.Fatalf("assignment.DriverID = %s, want driver-1", assignment.DriverID)
	}
}
