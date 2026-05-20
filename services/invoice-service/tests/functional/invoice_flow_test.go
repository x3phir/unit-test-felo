//go:build functional

package functional_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/felo/felo-backend/services/invoice-service/internal/domain"
	"github.com/felo/felo-backend/services/invoice-service/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestInvoiceFunctional_CreateInvoice_PersistsToDatabase(t *testing.T) {
	ctx := context.Background()
	db := openPG(t, getenv("FELO_INVOICE_PG_DSN", "postgres://felo:felo@127.0.0.1:54326/invoice_db?sslmode=disable"))
	t.Cleanup(func() { db.Close() })

	// Setup table since it might be empty
	_, err := db.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS invoices (
			invoice_id VARCHAR(255) PRIMARY KEY,
			subject_ref VARCHAR(255) NOT NULL,
			customer_id VARCHAR(255) NOT NULL,
			amount BIGINT NOT NULL,
			currency VARCHAR(10) NOT NULL,
			status VARCHAR(50) NOT NULL
		)
	`)
	if err != nil {
		t.Fatalf("failed to create table: %v", err)
	}

	subjectRef := "ord-ft-001"
	_, _ = db.Exec(ctx, "delete from invoices where subject_ref=$1", subjectRef)

	repo := &pgInvoiceRepo{db: db}
	publisher := &noopNotificationPublisher{}
	
	// mock clock to have predictable IDs
	now := time.Now().UTC()
	clock := func() time.Time { return now }

	svc := service.NewInvoiceService(repo, publisher, clock)

	inv, err := svc.CreateInvoice(ctx, subjectRef, 50000, "payer-001")
	if err != nil {
		t.Fatalf("CreateInvoice() error = %v", err)
	}

	var status string
	var amount int64
	if err := db.QueryRow(ctx, "select status, amount from invoices where invoice_id=$1", inv.InvoiceID).Scan(&status, &amount); err != nil {
		t.Fatalf("query persisted invoice: %v", err)
	}
	
	if status != string(domain.StatusIssued) {
		t.Fatalf("persisted status = %s, want %s", status, domain.StatusIssued)
	}
	if amount != 50000 {
		t.Fatalf("persisted amount = %d, want 50000", amount)
	}
}

type pgInvoiceRepo struct{ db *pgxpool.Pool }

func (r *pgInvoiceRepo) Create(ctx context.Context, inv *domain.Invoice) error {
	_, err := r.db.Exec(ctx, `insert into invoices (invoice_id, subject_ref, customer_id, amount, currency, status)
values ($1,$2,$3,$4,$5,$6)`,
		inv.InvoiceID, inv.SubjectRef, inv.CustomerID, inv.Amount, inv.Currency, string(inv.Status))
	return err
}

func (r *pgInvoiceRepo) GetByID(ctx context.Context, id string) (*domain.Invoice, error) {
	return nil, nil // Not needed for this specific test
}

func (r *pgInvoiceRepo) GetByOrderID(ctx context.Context, orderID string) ([]domain.Invoice, error) {
	return nil, nil // Not needed for this specific test
}

func (r *pgInvoiceRepo) UpdateStatus(ctx context.Context, id string, status domain.InvoiceStatus) error {
	return nil // Not needed for this specific test
}

func (r *pgInvoiceRepo) UpdatePaymentReference(ctx context.Context, id string, ref string) error {
	return nil // Not needed for this specific test
}

type noopNotificationPublisher struct{}

func (p *noopNotificationPublisher) PublishInvoiceNotification(ctx context.Context, inv *domain.Invoice) error {
	return nil
}

func openPG(t *testing.T, dsn string) *pgxpool.Pool {
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
