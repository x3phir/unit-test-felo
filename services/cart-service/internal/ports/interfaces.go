package ports

import (
	"context"

	"github.com/felo/felo-backend/services/cart-service/internal/domain"
)

type CartRepository interface {
	Save(ctx context.Context, cart domain.Cart) error
	GetByUserID(ctx context.Context, userID string) (domain.Cart, bool, error)
	Delete(ctx context.Context, userID string) error
}

type MerchantClient interface {
	GetItemPriceAndAvailability(ctx context.Context, merchantID string, menuItemID string) (int64, bool, error)
}
