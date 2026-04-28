package ports

import (
	"context"

	"github.com/felo/felo-backend/services/notification-service/internal/domain"
)

type NotificationProvider interface {
	SendPush(ctx context.Context, req domain.NotificationRequest) (domain.DeliveryStatus, error)
	SendWhatsApp(ctx context.Context, req domain.NotificationRequest) (domain.DeliveryStatus, error)
	SendSMS(ctx context.Context, req domain.NotificationRequest) (domain.DeliveryStatus, error)
}