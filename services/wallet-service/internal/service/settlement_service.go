package service

import (
	"context"
	"errors"
	"time"

	"github.com/felo/felo-backend/services/wallet-service/internal/domain"
	"github.com/felo/felo-backend/services/wallet-service/internal/ports"
)

var ErrInvalidSettlement = errors.New("invalid settlement")

type SettlementService struct {
	store   ports.SettlementStore
	wallets ports.WalletStore
	now     func() time.Time
}

func NewSettlementService(store ports.SettlementStore, wallets ports.WalletStore) *SettlementService {
	return &SettlementService{
		store:   store,
		wallets: wallets,
		now:     time.Now,
	}
}

func (s *SettlementService) ApplyRideSettlement(ctx context.Context, settlement domain.Settlement) (domain.SettlementRecord, error) {
	if settlement.IdempotencyKey == "" || settlement.DriverID == "" || settlement.Amount <= 0 {
		return domain.SettlementRecord{}, ErrInvalidSettlement
	}

	existing, found, err := s.store.GetByKey(ctx, settlement.IdempotencyKey)
	if err != nil {
		return domain.SettlementRecord{}, err
	}
	if found {
		return existing, nil
	}

	balance, err := s.wallets.Credit(ctx, settlement.DriverID, settlement.Amount)
	if err != nil {
		return domain.SettlementRecord{}, err
	}

	record := domain.SettlementRecord{
		IdempotencyKey: settlement.IdempotencyKey,
		TripID:         settlement.TripID,
		DriverID:       settlement.DriverID,
		Amount:         settlement.Amount,
		Currency:       settlement.Currency,
		BalanceAfter:   balance,
		ProcessedAt:    s.now(),
	}
	if err := s.store.Save(ctx, record); err != nil {
		return domain.SettlementRecord{}, err
	}

	return record, nil
}
