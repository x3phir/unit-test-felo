package ports

import (
	"context"

	"github.com/felo/felo-backend/services/invoice-service/internal/domain"
)

type InvoiceRepository interface {
	Create(ctx context.Context, invoice *domain.Invoice) error
	GetByID(ctx context.Context, id string) (*domain.Invoice, error)
	GetByOrderID(ctx context.Context, orderID string) ([]domain.Invoice, error)
	UpdateStatus(ctx context.Context, id string, status domain.InvoiceStatus) error
}

type NotificationPublisher interface {
	// Mempublikasikan event agar Notification Service mengirimkan nota digital [cite: 148]
	PublishInvoiceNotification(ctx context.Context, invoice *domain.Invoice) error
}