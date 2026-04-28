package ports

import (
	"context"

	"github.com/felo/felo-backend/services/feedback-service/internal/domain"
)

type FeedbackRepository interface {
	Save(ctx context.Context, feedback domain.Feedback) error
	GetByTripID(ctx context.Context, tripID string) (domain.Feedback, error)
}

type EventPublisher interface {
	Publish(ctx context.Context, event domain.Event) error
}
