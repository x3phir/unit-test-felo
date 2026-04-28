package ports

import (
	"context"

	"github.com/felo/felo-backend/services/merchant-service/internal/domain"
)

type MerchantRepository interface {
	GetByID(ctx context.Context, id string) (*domain.Merchant, error)
	Create(ctx context.Context, merchant *domain.Merchant) error
	UpdateStatus(ctx context.Context, id string, isClosed bool) error
}

type MenuRepository interface {
	GetByMerchantID(ctx context.Context, merchantID string) ([]domain.Menu, error)
	GetByID(ctx context.Context, id string) (*domain.Menu, error)
	GetByIDs(ctx context.Context, ids []string) ([]domain.Menu, error)
	Create(ctx context.Context, menu *domain.Menu) error
	Update(ctx context.Context, menu *domain.Menu) error
	UpdateAvailability(ctx context.Context, id string, isAvailable bool) error
}