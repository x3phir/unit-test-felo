package service

import (
	"context"
	"errors"

	"github.com/felo/felo-backend/services/tracking-service/internal/domain"
	"github.com/felo/felo-backend/services/tracking-service/internal/ports"
)

var (
	ErrInvalidInput     = errors.New("invalid tracking input")
	ErrSessionNotFound  = errors.New("tracking session not found")
	ErrSessionNotActive = errors.New("tracking session is not active")
)

type StartTrackingInput struct {
	ShipmentID string
	DriverID   string
}

type RecordLocationInput struct {
	SessionID  string
	Coordinate domain.Coordinate
	Speed      float64
	Heading    float64
}

type TrackingService struct {
	sessionRepo ports.TrackingSessionRepository
	recordRepo  ports.TrackingRecordRepository
	publisher   ports.EventPublisher
	clock       ports.Clock
	ids         ports.IDGenerator
}

func NewTrackingService(
	sessionRepo ports.TrackingSessionRepository,
	recordRepo ports.TrackingRecordRepository,
	publisher ports.EventPublisher,
	clock ports.Clock,
	ids ports.IDGenerator,
) *TrackingService {
	return &TrackingService{
		sessionRepo: sessionRepo,
		recordRepo:  recordRepo,
		publisher:   publisher,
		clock:       clock,
		ids:         ids,
	}
}

func (s *TrackingService) StartTracking(ctx context.Context, input StartTrackingInput) (domain.TrackingSession, error) {
	if input.ShipmentID == "" || input.DriverID == "" {
		return domain.TrackingSession{}, ErrInvalidInput
	}

	now := s.clock.Now()
	session := domain.TrackingSession{
		ID:         s.ids.NewID(),
		ShipmentID: input.ShipmentID,
		DriverID:   input.DriverID,
		Status:     domain.StatusActive,
		StartedAt:  now,
		UpdatedAt:  now,
	}

	if err := s.sessionRepo.Save(ctx, session); err != nil {
		return domain.TrackingSession{}, err
	}
	if err := s.publisher.Publish(ctx, domain.Event{Name: "tracking.started.v1", SessionID: session.ID, OccurredAt: now}); err != nil {
		return domain.TrackingSession{}, err
	}

	return session, nil
}

func (s *TrackingService) StopTracking(ctx context.Context, sessionID string) (domain.TrackingSession, error) {
	session, err := s.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		return domain.TrackingSession{}, ErrSessionNotFound
	}
	if session.Status != domain.StatusActive {
		return domain.TrackingSession{}, ErrSessionNotActive
	}

	now := s.clock.Now()
	session.Status = domain.StatusCompleted
	session.UpdatedAt = now
	session.EndedAt = &now

	if err := s.sessionRepo.Save(ctx, session); err != nil {
		return domain.TrackingSession{}, err
	}
	if err := s.publisher.Publish(ctx, domain.Event{Name: "tracking.completed.v1", SessionID: session.ID, OccurredAt: now}); err != nil {
		return domain.TrackingSession{}, err
	}

	return session, nil
}

func (s *TrackingService) RecordLocation(ctx context.Context, input RecordLocationInput) (domain.TrackingRecord, error) {
	if input.SessionID == "" {
		return domain.TrackingRecord{}, ErrInvalidInput
	}

	session, err := s.sessionRepo.GetByID(ctx, input.SessionID)
	if err != nil {
		return domain.TrackingRecord{}, ErrSessionNotFound
	}
	if session.Status != domain.StatusActive {
		return domain.TrackingRecord{}, ErrSessionNotActive
	}

	now := s.clock.Now()
	record := domain.TrackingRecord{
		ID:         s.ids.NewID(),
		SessionID:  input.SessionID,
		Coordinate: input.Coordinate,
		Speed:      input.Speed,
		Heading:    input.Heading,
		RecordedAt: now,
	}

	if err := s.recordRepo.Save(ctx, record); err != nil {
		return domain.TrackingRecord{}, err
	}

	session.UpdatedAt = now
	if err := s.sessionRepo.Save(ctx, session); err != nil {
		return domain.TrackingRecord{}, err
	}

	return record, nil
}

func (s *TrackingService) GetTrackingHistory(ctx context.Context, sessionID string) ([]domain.TrackingRecord, error) {
	if sessionID == "" {
		return nil, ErrInvalidInput
	}

	_, err := s.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		return nil, ErrSessionNotFound
	}

	return s.recordRepo.ListBySession(ctx, sessionID)
}
