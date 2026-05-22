package unit_test

import (
	"context"
	"testing"
	"time"

	"github.com/felo/felo-backend/services/tracking-service/internal/domain"
	"github.com/felo/felo-backend/services/tracking-service/internal/service"
	"go.uber.org/mock/gomock"
)

func TestTrackingService_StartTracking_CreatesSessionAndPublishesEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	sessionRepo := NewMockTrackingSessionRepository(ctrl)
	recordRepo := NewMockTrackingRecordRepository(ctrl)
	publisher := NewMockEventPublisher(ctrl)
	now := time.Date(2026, 5, 5, 10, 0, 0, 0, time.UTC)
	svc := service.NewTrackingService(sessionRepo, recordRepo, publisher, fixedClock{now: now}, &idGen{ids: []string{"track-001"}})

	sessionRepo.EXPECT().Save(gomock.Any(), gomock.AssignableToTypeOf(domain.TrackingSession{})).Return(nil)
	publisher.EXPECT().Publish(gomock.Any(), gomock.AssignableToTypeOf(domain.Event{})).DoAndReturn(
		func(_ context.Context, event domain.Event) error {
			if event.Name != "tracking.started.v1" {
				t.Fatalf("event.Name = %s, want tracking.started.v1", event.Name)
			}
			return nil
		},
	)

	session, err := svc.StartTracking(context.Background(), service.StartTrackingInput{
		ShipmentID: "ship-001",
		DriverID:   "driver-001",
	})
	if err != nil {
		t.Fatalf("StartTracking() error = %v", err)
	}
	if session.Status != domain.StatusActive {
		t.Fatalf("session.Status = %s, want %s", session.Status, domain.StatusActive)
	}
	if session.ID != "track-001" {
		t.Fatalf("session.ID = %s, want track-001", session.ID)
	}
}

func TestTrackingService_StartTracking_InvalidInput_ReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	svc := service.NewTrackingService(
		NewMockTrackingSessionRepository(ctrl),
		NewMockTrackingRecordRepository(ctrl),
		NewMockEventPublisher(ctrl),
		fixedClock{},
		&idGen{},
	)

	_, err := svc.StartTracking(context.Background(), service.StartTrackingInput{})
	if err != service.ErrInvalidInput {
		t.Fatalf("StartTracking() error = %v, want ErrInvalidInput", err)
	}
}

func TestTrackingService_StopTracking_CompletesSessionAndPublishesEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	sessionRepo := NewMockTrackingSessionRepository(ctrl)
	recordRepo := NewMockTrackingRecordRepository(ctrl)
	publisher := NewMockEventPublisher(ctrl)
	now := time.Date(2026, 5, 5, 12, 0, 0, 0, time.UTC)
	svc := service.NewTrackingService(sessionRepo, recordRepo, publisher, fixedClock{now: now}, &idGen{})

	sessionRepo.EXPECT().GetByID(gomock.Any(), "track-001").Return(domain.TrackingSession{
		ID: "track-001", Status: domain.StatusActive,
	}, nil)
	sessionRepo.EXPECT().Save(gomock.Any(), gomock.AssignableToTypeOf(domain.TrackingSession{})).Return(nil)
	publisher.EXPECT().Publish(gomock.Any(), gomock.AssignableToTypeOf(domain.Event{})).DoAndReturn(
		func(_ context.Context, event domain.Event) error {
			if event.Name != "tracking.completed.v1" {
				t.Fatalf("event.Name = %s, want tracking.completed.v1", event.Name)
			}
			return nil
		},
	)

	session, err := svc.StopTracking(context.Background(), "track-001")
	if err != nil {
		t.Fatalf("StopTracking() error = %v", err)
	}
	if session.Status != domain.StatusCompleted {
		t.Fatalf("session.Status = %s, want %s", session.Status, domain.StatusCompleted)
	}
	if session.EndedAt == nil || !session.EndedAt.Equal(now) {
		t.Fatalf("session.EndedAt = %v, want %v", session.EndedAt, now)
	}
}

func TestTrackingService_StopTracking_NotFound_ReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	sessionRepo := NewMockTrackingSessionRepository(ctrl)
	svc := service.NewTrackingService(
		sessionRepo,
		NewMockTrackingRecordRepository(ctrl),
		NewMockEventPublisher(ctrl),
		fixedClock{},
		&idGen{},
	)

	sessionRepo.EXPECT().GetByID(gomock.Any(), "unknown").Return(domain.TrackingSession{}, service.ErrSessionNotFound)

	_, err := svc.StopTracking(context.Background(), "unknown")
	if err != service.ErrSessionNotFound {
		t.Fatalf("StopTracking() error = %v, want ErrSessionNotFound", err)
	}
}

func TestTrackingService_RecordLocation_ValidInput_StoresRecord(t *testing.T) {
	ctrl := gomock.NewController(t)
	sessionRepo := NewMockTrackingSessionRepository(ctrl)
	recordRepo := NewMockTrackingRecordRepository(ctrl)
	publisher := NewMockEventPublisher(ctrl)
	now := time.Date(2026, 5, 5, 10, 30, 0, 0, time.UTC)
	svc := service.NewTrackingService(sessionRepo, recordRepo, publisher, fixedClock{now: now}, &idGen{ids: []string{"rec-001"}})

	sessionRepo.EXPECT().GetByID(gomock.Any(), "track-001").Return(domain.TrackingSession{
		ID: "track-001", Status: domain.StatusActive,
	}, nil)
	recordRepo.EXPECT().Save(gomock.Any(), gomock.AssignableToTypeOf(domain.TrackingRecord{})).Return(nil)
	sessionRepo.EXPECT().Save(gomock.Any(), gomock.AssignableToTypeOf(domain.TrackingSession{})).Return(nil)

	record, err := svc.RecordLocation(context.Background(), service.RecordLocationInput{
		SessionID:  "track-001",
		Coordinate: domain.Coordinate{Latitude: -6.2, Longitude: 106.8},
		Speed:      40.5,
		Heading:    180.0,
	})
	if err != nil {
		t.Fatalf("RecordLocation() error = %v", err)
	}
	if record.ID != "rec-001" {
		t.Fatalf("record.ID = %s, want rec-001", record.ID)
	}
}

func TestTrackingService_GetTrackingHistory_ReturnsRecords(t *testing.T) {
	ctrl := gomock.NewController(t)
	sessionRepo := NewMockTrackingSessionRepository(ctrl)
	recordRepo := NewMockTrackingRecordRepository(ctrl)
	publisher := NewMockEventPublisher(ctrl)
	svc := service.NewTrackingService(sessionRepo, recordRepo, publisher, fixedClock{}, &idGen{})

	sessionRepo.EXPECT().GetByID(gomock.Any(), "track-001").Return(domain.TrackingSession{
		ID: "track-001", Status: domain.StatusCompleted,
	}, nil)
	recordRepo.EXPECT().ListBySession(gomock.Any(), "track-001").Return([]domain.TrackingRecord{
		{ID: "rec-001", SessionID: "track-001"},
		{ID: "rec-002", SessionID: "track-001"},
	}, nil)

	records, err := svc.GetTrackingHistory(context.Background(), "track-001")
	if err != nil {
		t.Fatalf("GetTrackingHistory() error = %v", err)
	}
	if len(records) != 2 {
		t.Fatalf("len(records) = %d, want 2", len(records))
	}
}
