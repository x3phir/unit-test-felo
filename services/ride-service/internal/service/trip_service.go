package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/felo/felo-backend/services/ride-service/internal/domain"
	"github.com/felo/felo-backend/services/ride-service/internal/ports"
)

var (
	ErrInvalidTransition = errors.New("invalid trip state transition")
	ErrInvalidInput      = errors.New("invalid trip input")
)

type RequestRideInput struct {
	CustomerID   string
	Pickup       domain.Coordinate
	Destination  domain.Coordinate
	FareEstimate int64
}

type GenerateNowQRInput struct {
	CustomerID   string
	Destination  domain.Coordinate
	FareEstimate int64
}

type TripService struct {
	repo      ports.TripRepository
	publisher ports.EventPublisher
	clock     ports.Clock
	ids       ports.IDGenerator
	qrTTL     time.Duration
}

func NewTripService(repo ports.TripRepository, publisher ports.EventPublisher, clock ports.Clock, ids ports.IDGenerator) *TripService {
	return &TripService{
		repo:      repo,
		publisher: publisher,
		clock:     clock,
		ids:       ids,
		qrTTL:     10 * time.Minute,
	}
}

func (s *TripService) RequestRide(ctx context.Context, input RequestRideInput) (domain.Trip, error) {
	if input.CustomerID == "" || input.FareEstimate <= 0 {
		return domain.Trip{}, ErrInvalidInput
	}

	now := s.clock.Now()
	trip := domain.Trip{
		ID:           s.ids.NewID(),
		CustomerID:   input.CustomerID,
		Pickup:       input.Pickup,
		Destination:  input.Destination,
		FareEstimate: input.FareEstimate,
		State:        domain.StateMatching,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := s.repo.Save(ctx, trip); err != nil {
		return domain.Trip{}, err
	}
	if err := s.publisher.Publish(ctx, domain.Event{Name: "ride.created.v1", TripID: trip.ID, OccurredAt: now}); err != nil {
		return domain.Trip{}, err
	}

	return trip, nil
}

func (s *TripService) StartRide(ctx context.Context, tripID string) (domain.Trip, error) {
	trip, err := s.repo.GetByID(ctx, tripID)
	if err != nil {
		return domain.Trip{}, err
	}
	if trip.State != domain.StateArrived {
		return domain.Trip{}, fmt.Errorf("%w: want %s got %s", ErrInvalidTransition, domain.StateArrived, trip.State)
	}

	trip.State = domain.StateOnRide
	trip.UpdatedAt = s.clock.Now()
	if err := s.repo.Save(ctx, trip); err != nil {
		return domain.Trip{}, err
	}
	if err := s.publisher.Publish(ctx, domain.Event{Name: "ride.started.v1", TripID: trip.ID, OccurredAt: trip.UpdatedAt}); err != nil {
		return domain.Trip{}, err
	}

	return trip, nil
}

func (s *TripService) CompleteRide(ctx context.Context, tripID string) (domain.Trip, error) {
	trip, err := s.repo.GetByID(ctx, tripID)
	if err != nil {
		return domain.Trip{}, err
	}
	if trip.State != domain.StateOnRide {
		return domain.Trip{}, fmt.Errorf("%w: want %s got %s", ErrInvalidTransition, domain.StateOnRide, trip.State)
	}

	trip.State = domain.StateCompleted
	trip.UpdatedAt = s.clock.Now()
	if err := s.repo.Save(ctx, trip); err != nil {
		return domain.Trip{}, err
	}
	if err := s.publisher.Publish(ctx, domain.Event{Name: "ride.completed.v1", TripID: trip.ID, OccurredAt: trip.UpdatedAt}); err != nil {
		return domain.Trip{}, err
	}

	return trip, nil
}

func (s *TripService) GenerateNowQR(ctx context.Context, input GenerateNowQRInput) (domain.Trip, error) {
	if input.CustomerID == "" || input.FareEstimate <= 0 {
		return domain.Trip{}, ErrInvalidInput
	}

	now := s.clock.Now()
	trip := domain.Trip{
		ID:           s.ids.NewID(),
		CustomerID:   input.CustomerID,
		Destination:  input.Destination,
		FareEstimate: input.FareEstimate,
		State:        domain.StateQRGenerated,
		QRCode:       "qr-" + s.ids.NewID(),
		QRExpiresAt:  now.Add(s.qrTTL),
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := s.repo.Save(ctx, trip); err != nil {
		return domain.Trip{}, err
	}

	return trip, nil
}
