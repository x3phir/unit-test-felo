package ports

import (
	"context"

	"github.com/felo/felo-backend/services/shipment-service/internal/domain"
)

type ShipmentRepository interface {
	Save(ctx context.Context, shipment domain.Shipment) error
	GetByID(ctx context.Context, shipmentID string) (domain.Shipment, error)
}

type EventPublisher interface {
	Publish(ctx context.Context, event domain.Event) error
}
