package ports

import (
	"context"

	"github.com/felo/felo-backend/services/order-service/internal/domain"
)

type OrderRepository interface {
	Save(ctx context.Context, order domain.FoodOrder) error
	GetByID(ctx context.Context, orderID string) (domain.FoodOrder, error)
}

type LocationClient interface {
	GetDistanceKM(ctx context.Context, origin string, destination string) (float64, error)
}

type AuthClient interface {
	SendOTP(ctx context.Context, userID string) error
	VerifyOTP(ctx context.Context, userID string, otp string) (bool, error)
}

type EventPublisher interface {
	Publish(ctx context.Context, event domain.Event) error
}

type IDGenerator interface {
	NewID() string
}
