package service

import (
	"context"
	"errors"
	"sort"
	"time"

	"github.com/felo/felo-backend/services/matching-service/internal/domain"
	"github.com/felo/felo-backend/services/matching-service/internal/ports"
)

var ErrNoDriversAvailable = errors.New("no drivers available")

type MatchingService struct {
	availability ports.AvailabilityReader
	assignments  ports.AssignmentRepository
	publisher    ports.EventPublisher
	maxRetries   int
	radiusStepKM float64
	now          func() time.Time
}

func NewMatchingService(availability ports.AvailabilityReader, assignments ports.AssignmentRepository, publisher ports.EventPublisher) *MatchingService {
	return &MatchingService{
		availability: availability,
		assignments:  assignments,
		publisher:    publisher,
		maxRetries:   3,
		radiusStepKM: 1.5,
		now:          time.Now,
	}
}

func (s *MatchingService) AssignDriver(ctx context.Context, request domain.MatchRequest) (domain.Assignment, error) {
	radius := request.InitialRadiusKM
	if radius <= 0 {
		radius = 1
	}

	for attempt := 0; attempt <= s.maxRetries; attempt++ {
		candidates, err := s.availability.FindAvailableDrivers(ctx, request.Pickup, radius)
		if err != nil {
			return domain.Assignment{}, err
		}
		if len(candidates) > 0 {
			sort.Slice(candidates, func(i, j int) bool {
				return candidates[i].DistanceKM < candidates[j].DistanceKM
			})

			assignment := domain.Assignment{
				RideID:         request.RideID,
				DriverID:       candidates[0].DriverID,
				SearchRadiusKM: radius,
				AssignedAt:     s.now(),
			}
			if err := s.assignments.Save(ctx, assignment); err != nil {
				return domain.Assignment{}, err
			}
			if err := s.publisher.Publish(ctx, domain.Event{
				Name:       "driver.matched.v1",
				RideID:     assignment.RideID,
				DriverID:   assignment.DriverID,
				OccurredAt: assignment.AssignedAt,
			}); err != nil {
				return domain.Assignment{}, err
			}
			return assignment, nil
		}
		radius += s.radiusStepKM
	}

	return domain.Assignment{}, ErrNoDriversAvailable
}
