package unit_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/felo/felo-backend/services/wallet-service/internal/domain"
	"github.com/felo/felo-backend/services/wallet-service/internal/service"
)

func TestSettlementService_ApplyRideSettlement_CreditsDriverAndStoresRecord(t *testing.T) {
	store := &settlementStoreFake{records: map[string]domain.SettlementRecord{}}
	wallets := &walletStoreFake{balances: map[string]int64{"driver-1": 100000}}
	svc := service.NewSettlementService(store, wallets)

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

func TestSettlementService_ApplyRideSettlement_IdempotentKeyReturnsExistingRecord(t *testing.T) {
	now := time.Date(2026, 4, 28, 12, 0, 0, 0, time.UTC)
	store := &settlementStoreFake{records: map[string]domain.SettlementRecord{
		"settlement-1": {
			IdempotencyKey: "settlement-1",
			TripID:         "trip-1",
			DriverID:       "driver-1",
			Amount:         25000,
			Currency:       "IDR",
			BalanceAfter:   125000,
			ProcessedAt:    now,
		},
	}}
	wallets := &walletStoreFake{balances: map[string]int64{"driver-1": 100000}}
	svc := service.NewSettlementService(store, wallets)

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
	if wallets.creditCalls != 0 {
		t.Fatalf("wallets.creditCalls = %d, want 0", wallets.creditCalls)
	}
}

func TestSettlementService_ApplyRideSettlement_InvalidAmountReturnsError(t *testing.T) {
	svc := service.NewSettlementService(&settlementStoreFake{records: map[string]domain.SettlementRecord{}}, &walletStoreFake{balances: map[string]int64{}})

	_, err := svc.ApplyRideSettlement(context.Background(), domain.Settlement{
		IdempotencyKey: "settlement-1",
		TripID:         "trip-1",
		DriverID:       "driver-1",
		Amount:         0,
		Currency:       "IDR",
	})
	if err == nil {
		t.Fatal("ApplyRideSettlement() error = nil, want error")
	}
}

func TestSettlementService_ApplyRideSettlement_WalletFailureReturnsError(t *testing.T) {
	store := &settlementStoreFake{records: map[string]domain.SettlementRecord{}}
	wallets := &walletStoreFake{balances: map[string]int64{}, err: errors.New("wallet unavailable")}
	svc := service.NewSettlementService(store, wallets)

	_, err := svc.ApplyRideSettlement(context.Background(), domain.Settlement{
		IdempotencyKey: "settlement-1",
		TripID:         "trip-1",
		DriverID:       "driver-1",
		Amount:         25000,
		Currency:       "IDR",
	})
	if err == nil {
		t.Fatal("ApplyRideSettlement() error = nil, want error")
	}
	if _, ok := store.records["settlement-1"]; ok {
		t.Fatal("store.records should not contain failed settlement")
	}
}

type settlementStoreFake struct {
	records map[string]domain.SettlementRecord
}

func (f *settlementStoreFake) GetByKey(_ context.Context, key string) (domain.SettlementRecord, bool, error) {
	record, ok := f.records[key]
	return record, ok, nil
}

func (f *settlementStoreFake) Save(_ context.Context, record domain.SettlementRecord) error {
	f.records[record.IdempotencyKey] = record
	return nil
}

type walletStoreFake struct {
	balances    map[string]int64
	creditCalls int
	err         error
}

func (f *walletStoreFake) Credit(_ context.Context, ownerID string, amount int64) (int64, error) {
	f.creditCalls++
	if f.err != nil {
		return 0, f.err
	}
	f.balances[ownerID] += amount
	return f.balances[ownerID], nil
}
