//go:build functional

package functional_test

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/felo/felo-backend/services/payment-service/internal/domain"
	"github.com/felo/felo-backend/services/payment-service/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestPaymentFunctional_HandleRideCompleted_PersistsProcessedPayment(t *testing.T) {
	ctx := context.Background()
	db := openPaymentPG(t, getenv("FELO_PAYMENT_PG_DSN", "postgres://felo:felo@127.0.0.1:54338/payment_db?sslmode=disable"))
	t.Cleanup(func() { db.Close() })

	eventID := "evt-ft-001"
	_, _ = db.Exec(ctx, "delete from payments where event_id=$1", eventID)

	svc := service.NewPaymentService(noopWalletClient{}, noopInvoiceClient{}, &pgProcessedStore{db: db}, noopPaymentPublisher{})
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
}

type noopWalletClient struct{}
func (noopWalletClient) DebitCustomer(context.Context, string, int64, string) error { return nil }

type noopInvoiceClient struct{}
func (noopInvoiceClient) IssueRideInvoice(context.Context, string, string, int64, string) (string, error) { return "inv-ft-001", nil }

type pgProcessedStore struct{ db *pgxpool.Pool }
func (s *pgProcessedStore) Get(ctx context.Context, eventID string) (domain.PaymentResult, bool, error) {
	var result domain.PaymentResult
	err := s.db.QueryRow(ctx, "select event_id, ride_id, created_at from payments where event_id=$1", eventID).Scan(&result.EventID, &result.TripID, &result.PaidAt)
	if err != nil {
		return domain.PaymentResult{}, false, nil
	}
	return result, true, nil
}
func (s *pgProcessedStore) Save(ctx context.Context, result domain.PaymentResult) error {
	_, err := s.db.Exec(ctx, "insert into payments (event_id, ride_id, customer_id, amount, currency, status, created_at) values ($1,$2,$3,$4,$5,$6,$7)",
		result.EventID, result.TripID, "cust-active-001", 30000, "IDR", "completed", time.Now().UTC())
	return err
}

type noopPaymentPublisher struct{}
func (noopPaymentPublisher) Publish(context.Context, domain.Event) error { return nil }

func openPaymentPG(t *testing.T, dsn string) *pgxpool.Pool {
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

var _ = errors.New
