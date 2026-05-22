package ports

import (
	"context"
	"time"

	"github.com/felo/felo-backend/services/tracking-service/internal/domain"
)

type TrackingSessionRepository interface {
	Save(ctx context.Context, session domain.TrackingSession) error
	GetByID(ctx context.Context, sessionID string) (domain.TrackingSession, error)
}

type TrackingRecordRepository interface {
	Save(ctx context.Context, record domain.TrackingRecord) error
	ListBySession(ctx context.Context, sessionID string) ([]domain.TrackingRecord, error)
}

type EventPublisher interface {
	Publish(ctx context.Context, event domain.Event) error
}

type Clock interface {
	Now() time.Time
}

type IDGenerator interface {
	NewID() string
}
