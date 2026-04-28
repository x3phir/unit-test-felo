package ports

import (
	"context"
	"time"

	"github.com/felo/felo-backend/services/ride-service/internal/domain"
)

type TripRepository interface {
	Save(ctx context.Context, trip domain.Trip) error
	GetByID(ctx context.Context, tripID string) (domain.Trip, error)
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
