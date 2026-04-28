package ports

import (
	"context"

	"github.com/felo/felo-backend/services/wallet-service/internal/domain"
)

type SettlementStore interface {
	GetByKey(ctx context.Context, key string) (domain.SettlementRecord, bool, error)
	Save(ctx context.Context, record domain.SettlementRecord) error
}

type WalletStore interface {
	Credit(ctx context.Context, ownerID string, amount int64) (int64, error)
}
