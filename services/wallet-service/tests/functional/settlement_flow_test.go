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
	db := openWalletPG(t, getenv("FELO_WALLET_PG_DSN", "postgres://felo:felo@127.0.0.1:54339/wallet_db?sslmode=disable"))
	t.Cleanup(func() { db.Close() })

	initWalletTables(t, db)

	ref := "settlement-ft-001"
	_, _ = db.Exec(ctx, "delete from wallet_ledger where reference=$1", ref)
	_, _ = db.Exec(ctx, "delete from wallets where owner_id=$1", "driver-ft-001")
	_, _ = db.Exec(ctx, `insert into wallets (owner_id, owner_type, balance, currency, updated_at)
		values ('driver-ft-001','driver',0,'IDR',$1)`, time.Now().UTC())

	svc := service.NewSettlementService(&pgSettlementStore{db: db}, &pgWalletCredit{db: db})

	record, err := svc.ApplyRideSettlement(ctx, domain.Settlement{
		IdempotencyKey: ref,
		TripID:         "trip-ft-001",
		DriverID:       "driver-ft-001",
		Amount:         25000,
		Currency:       "IDR",
	})
	if err != nil {
		t.Fatalf("ApplyRideSettlement() error = %v", err)
	}
	if record.BalanceAfter != 25000 {
		t.Fatalf("record.BalanceAfter = %d, want 25000", record.BalanceAfter)
	}

	var balance int64
	if err := db.QueryRow(ctx, "select balance from wallets where owner_id=$1", "driver-ft-001").Scan(&balance); err != nil {
		t.Fatalf("query wallet balance: %v", err)
	}
	if balance != 25000 {
		t.Fatalf("wallet balance = %d, want 25000", balance)
	}
}

func TestWalletFunctional_ApplyRideSettlement_IdempotentDuplicate(t *testing.T) {
	ctx := context.Background()
	db := openWalletPG(t, getenv("FELO_WALLET_PG_DSN", "postgres://felo:felo@127.0.0.1:54339/wallet_db?sslmode=disable"))
	t.Cleanup(func() { db.Close() })

	initWalletTables(t, db)

	ref := "settlement-ft-002"
	_, _ = db.Exec(ctx, "delete from wallet_ledger where reference=$1", ref)
	_, _ = db.Exec(ctx, "delete from wallets where owner_id=$1", "driver-ft-002")
	_, _ = db.Exec(ctx, `insert into wallets (owner_id, owner_type, balance, currency, updated_at)
		values ('driver-ft-002','driver',50000,'IDR',$1)`, time.Now().UTC())
	_, _ = db.Exec(ctx, `insert into wallet_ledger (reference, owner_id, delta, reason, created_at)
		values ($1,'driver-ft-002',25000,'settlement',$2)`, ref, time.Now().UTC())

	svc := service.NewSettlementService(&pgSettlementStore{db: db}, &pgWalletCredit{db: db})

	record, err := svc.ApplyRideSettlement(ctx, domain.Settlement{
		IdempotencyKey: ref,
		TripID:         "trip-ft-002",
		DriverID:       "driver-ft-002",
		Amount:         25000,
		Currency:       "IDR",
	})
	if err != nil {
		t.Fatalf("ApplyRideSettlement() error = %v", err)
	}
	if record.Amount != 25000 {
		t.Fatalf("record.Amount = %d, want 25000 (from existing)", record.Amount)
	}
}

type pgSettlementStore struct{ db *pgxpool.Pool }

func (s *pgSettlementStore) GetByKey(ctx context.Context, key string) (domain.SettlementRecord, bool, error) {
	var record domain.SettlementRecord
	err := s.db.QueryRow(ctx,
		"select reference, owner_id, delta, created_at from wallet_ledger where reference=$1", key).
		Scan(&record.IdempotencyKey, &record.DriverID, &record.Amount, &record.ProcessedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return domain.SettlementRecord{}, false, nil
		}
		return domain.SettlementRecord{}, false, err
	}
	return record, true, nil
}

func (s *pgSettlementStore) Save(ctx context.Context, record domain.SettlementRecord) error {
	_, err := s.db.Exec(ctx,
		`insert into wallet_ledger (reference, owner_id, delta, reason, created_at)
values ($1,$2,$3,'settlement',$4)`,
		record.IdempotencyKey, record.DriverID, record.Amount, record.ProcessedAt)
	return err
}

type pgWalletCredit struct{ db *pgxpool.Pool }

func (w *pgWalletCredit) Credit(ctx context.Context, ownerID string, amount int64) (int64, error) {
	tx, err := w.db.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)

	var balance int64
	if err := tx.QueryRow(ctx, "select balance from wallets where owner_id=$1 for update", ownerID).Scan(&balance); err != nil {
		return 0, err
	}
	balance += amount
	if _, err := tx.Exec(ctx,
		"update wallets set balance=$2, updated_at=$3 where owner_id=$1",
		ownerID, balance, time.Now().UTC()); err != nil {
		return 0, err
	}
	if err := tx.Commit(ctx); err != nil {
		return 0, err
	}
	return balance, nil
}

func initWalletTables(t *testing.T, db *pgxpool.Pool) {
	t.Helper()
	ctx := context.Background()
	_, err := db.Exec(ctx, `create table if not exists wallets (
		owner_id text primary key,
		owner_type text not null,
		balance bigint not null,
		currency text not null default 'IDR',
		updated_at timestamptz not null
	)`)
	if err != nil {
		t.Fatalf("initWalletTables wallets: %v", err)
	}
	_, err = db.Exec(ctx, `create table if not exists wallet_ledger (
		reference text primary key,
		owner_id text not null,
		delta bigint not null,
		reason text not null,
		created_at timestamptz not null
	)`)
	if err != nil {
		t.Fatalf("initWalletTables wallet_ledger: %v", err)
	}
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
