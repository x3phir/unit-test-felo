package unit_test

import (
	"context"
	"testing"

	"github.com/felo/felo-backend/services/wallet-service/internal/domain"
	"github.com/felo/felo-backend/services/wallet-service/internal/service"
	"go.uber.org/mock/gomock"
)

func TestSettlementService_ApplyRideSettlement_CreditsAndStoresWithGomock(t *testing.T) {
	ctrl := gomock.NewController(t)
	store := NewMockSettlementStore(ctrl)
	wallets := NewMockWalletStore(ctrl)
	svc := service.NewSettlementService(store, wallets)

	store.EXPECT().GetByKey(gomock.Any(), "settlement-1").Return(domain.SettlementRecord{}, false, nil)
	wallets.EXPECT().Credit(gomock.Any(), "driver-1", int64(25000)).Return(int64(125000), nil)
	store.EXPECT().Save(gomock.Any(), gomock.AssignableToTypeOf(domain.SettlementRecord{})).Return(nil)

	record, err := svc.ApplyRideSettlement(context.Background(), domain.Settlement{
		IdempotencyKey: "settlement-1",
		TripID:         "trip-1",
		DriverID:       "driver-1",
		Amount:         25000,
		Currency:       "IDR",
	})
	if err != nil {
		t.Fatalf("ApplyRideSettlement() error = %v", err)
	}
	if record.BalanceAfter != 125000 {
		t.Fatalf("record.BalanceAfter = %d, want 125000", record.BalanceAfter)
	}
}
