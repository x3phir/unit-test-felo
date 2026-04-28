package ports

import (
	"context"

	"github.com/felo/felo-backend/services/send-order-service/internal/domain"
)

type SendOrderRepository interface {
	Save(ctx context.Context, order domain.SendOrder) error
}

type PricingClient interface {
	CalculateShippingFee(ctx context.Context, pkg domain.PackageDetails, origin string, destination string) (int64, error)
}

type InvoiceClient interface {
	CreateInvoice(ctx context.Context, orderID string, payerID string, payerType domain.PayerType, amount int64) error
}

type EventPublisher interface {
	Publish(ctx context.Context, event domain.Event) error
}

type IDGenerator interface {
	NewID() string
}
