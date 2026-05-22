//go:build functional

package functional_test

import (
	"context"
	"testing"
	"time"

	"github.com/felo/felo-backend/services/payment-service/internal/domain"
	"github.com/felo/felo-backend/services/payment-service/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestPaymentFunctional_HandleRideCompleted_PersistsPaymentToDatabase(t *testing.T) {
	ctx := context.Background()
	db := openPaymentPG(t, getenv("FELO_PAYMENT_PG_DSN", "postgres://felo:felo@127.0.0.1:54338/payment_db?sslmode=disable"))
	t.Cleanup(func() { db.Close() })

	initPaymentTables(t, db)

	eventID := "evt-ft-001"
	_, _ = db.Exec(ctx, "delete from payments where event_id=$1", eventID)

	svc := service.NewPaymentService(
		noopWalletClient{},
		noopInvoiceClient{},
		&pgPaymentStore{db: db},
		noopPaymentPublisher{},
	)

	result, err := svc.HandleRideCompleted(ctx, domain.RideCompletedEvent{
		EventID:    eventID,
		TripID:     "ride-ft-001",
		CustomerID: "cust-active-001",
		Amount:     30000,
		Currency:   "IDR",
	})
	if err != nil {
		t.Fatalf("HandleRideCompleted() error = %v", err)
	}
	if result.EventID != eventID {
		t.Fatalf("result.EventID = %s, want %s", result.EventID, eventID)
	}

	var status string
	if err := db.QueryRow(ctx, "select status from payments where event_id=$1", eventID).Scan(&status); err != nil {
		t.Fatalf("query persisted payment: %v", err)
	}
	if status != "completed" {
		t.Fatalf("persisted status = %s, want completed", status)
	}
}

func TestPaymentFunctional_HandleRideCompleted_IdempotentDuplicate(t *testing.T) {
	ctx := context.Background()
	db := openPaymentPG(t, getenv("FELO_PAYMENT_PG_DSN", "postgres://felo:felo@127.0.0.1:54338/payment_db?sslmode=disable"))
	t.Cleanup(func() { db.Close() })

	initPaymentTables(t, db)

	eventID := "evt-ft-002"
	_, _ = db.Exec(ctx, "delete from payments where event_id=$1", eventID)
	_, _ = db.Exec(ctx, `insert into payments (event_id, ride_id, customer_id, amount, currency, status, created_at)
		values ($1,'ride-ft-002','cust-active-002',40000,'IDR','completed',$2)`,
		eventID, time.Now().UTC())

	svc := service.NewPaymentService(
		noopWalletClient{},
		noopInvoiceClient{},
		&pgPaymentStore{db: db},
		noopPaymentPublisher{},
	)

	result, err := svc.HandleRideCompleted(ctx, domain.RideCompletedEvent{
		EventID:    eventID,
		TripID:     "ride-ft-002",
		CustomerID: "cust-active-002",
		Amount:     40000,
		Currency:   "IDR",
	})
	if err != nil {
		t.Fatalf("HandleRideCompleted() error = %v", err)
	}
	if result.EventID != eventID {
		t.Fatalf("result.EventID = %s, want %s", result.EventID, eventID)
	}
}

type pgPaymentStore struct{ db *pgxpool.Pool }

func (s *pgPaymentStore) Get(ctx context.Context, eventID string) (domain.PaymentResult, bool, error) {
	var result domain.PaymentResult
	err := s.db.QueryRow(ctx,
		"select event_id, ride_id, created_at from payments where event_id=$1", eventID).
		Scan(&result.EventID, &result.TripID, &result.PaidAt)
	if err != nil {
		return domain.PaymentResult{}, false, nil
	}
	return result, true, nil
}

func (s *pgPaymentStore) Save(ctx context.Context, result domain.PaymentResult) error {
	_, err := s.db.Exec(ctx,
		`insert into payments (event_id, ride_id, customer_id, amount, currency, status, created_at)
values ($1,$2,$3,$4,$5,$6,$7)`,
		result.EventID, result.TripID, "cust-active-001", 30000, "IDR", "completed", result.PaidAt)
	return err
}

func initPaymentTables(t *testing.T, db *pgxpool.Pool) {
	t.Helper()
	ctx := context.Background()
	_, err := db.Exec(ctx, `create table if not exists payments (
		event_id text primary key,
		ride_id text not null,
		customer_id text not null,
		amount bigint not null,
		currency text not null,
		status text not null,
		created_at timestamptz not null
	)`)
	if err != nil {
		t.Fatalf("initPaymentTables: %v", err)
	}
}
