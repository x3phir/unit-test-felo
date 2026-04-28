package ports

import (
	"context"

	"github.com/felo/felo-backend/services/matching-service/internal/domain"
)

type AvailabilityReader interface {
	FindAvailableDrivers(ctx context.Context, pickup domain.Coordinate, radiusKM float64) ([]domain.DriverCandidate, error)
}

type AssignmentRepository interface {
	Save(ctx context.Context, assignment domain.Assignment) error
}

type EventPublisher interface {
	Publish(ctx context.Context, event domain.Event) error
}
