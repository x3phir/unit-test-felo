//go:build functional

package functional_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/felo/felo-backend/services/wallet-service/internal/domain"
	"github.com/felo/felo-backend/services/wallet-service/internal/service"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestWalletFunctional_ApplyRideSettlement_UpdatesWalletAndLedger(t *testing.T) {
	ctx := context.Background()
	db := openWalletPG(t, getenv("FELO_WALLET_PG_DSN", "postgres://felo:felo@127.0.0.1:54323/wallet_db?sslmode=disable"))
	t.Cleanup(func() { db.Close() })

	ref := "settlement-ft-001"
	_, _ = db.Exec(ctx, "delete from wallet_ledger where reference=$1", ref)
	_, _ = db.Exec(ctx, "delete from wallet_ledger where owner_id=$1", "driver-active-001")
	_, _ = db.Exec(ctx, "update wallets set balance=0, updated_at=now() where owner_id=$1", "driver-active-001")

	svc := service.NewSettlementService(&pgSettlementStore{db: db}, &pgWalletStore{db: db})
	record, err := svc.ApplyRideSettlement(ctx, domain.Settlement{
		IdempotencyKey: ref,
		TripID:         "trip-ft-001",
		DriverID:       "driver-active-001",
		Amount:         25000,
		Currency:       "IDR",
	})
	if err != nil {
		t.Fatalf("ApplyRideSettlement() error = %v", err)
	}
	if record.BalanceAfter != 25000 {
		t.Fatalf("record.BalanceAfter = %d, want 25000", record.BalanceAfter)
	}
}

type pgSettlementStore struct{ db *pgxpool.Pool }
func (s *pgSettlementStore) GetByKey(ctx context.Context, key string) (domain.SettlementRecord, bool, error) {
	var record domain.SettlementRecord
	err := s.db.QueryRow(ctx, `select reference, owner_id, delta, created_at from wallet_ledger where reference=$1`, key).
		Scan(&record.IdempotencyKey, &record.DriverID, &record.Amount, &record.ProcessedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return domain.SettlementRecord{}, false, nil
		}
		return domain.SettlementRecord{}, false, err
	}
	return record, true, nil
}
func (s *pgSettlementStore) Save(_ context.Context, _ domain.SettlementRecord) error { return nil }

type pgWalletStore struct{ db *pgxpool.Pool }
func (s *pgWalletStore) Credit(ctx context.Context, ownerID string, amount int64) (int64, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)
	var balance int64
	if err := tx.QueryRow(ctx, "select balance from wallets where owner_id=$1 for update", ownerID).Scan(&balance); err != nil {
		return 0, err
	}
	balance += amount
	if _, err := tx.Exec(ctx, "update wallets set balance=$2, updated_at=$3 where owner_id=$1", ownerID, balance, time.Now().UTC()); err != nil {
		return 0, err
	}
	if _, err := tx.Exec(ctx, "insert into wallet_ledger (reference, owner_id, delta, reason, created_at) values ($1,$2,$3,$4,$5)", "settlement-"+ownerID, ownerID, amount, "functional_test", time.Now().UTC()); err != nil {
		return 0, err
	}
	if err := tx.Commit(ctx); err != nil {
		return 0, err
	}
	return balance, nil
}

func openWalletPG(t *testing.T, dsn string) *pgxpool.Pool {
	t.Helper()
	db, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		t.Fatalf("pgxpool.New() error = %v", err)
	}
	return db
}

func getenv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
