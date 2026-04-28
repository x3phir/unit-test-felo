package ports

import (
	"context"

	"github.com/felo/felo-backend/services/payment-service/internal/domain"
)

type WalletClient interface {
	DebitCustomer(ctx context.Context, customerID string, amount int64, idempotencyKey string) error
}

type InvoiceClient interface {
	IssueRideInvoice(ctx context.Context, tripID string, customerID string, amount int64, currency string) (string, error)
}

type ProcessedEventStore interface {
	Get(ctx context.Context, eventID string) (domain.PaymentResult, bool, error)
	Save(ctx context.Context, result domain.PaymentResult) error
}

type EventPublisher interface {
	Publish(ctx context.Context, event domain.Event) error
}
