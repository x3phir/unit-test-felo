//go:build functional

package functional_test

import (
	"context"
	"os"
	"testing"

	"github.com/felo/felo-backend/services/send-order-service/internal/domain"
	"github.com/felo/felo-backend/services/send-order-service/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestSendOrderFunctional_CreateSendOrder_PersistsToDatabase(t *testing.T) {
	ctx := context.Background()
	db := openSendOrderPG(t, getenv("FELO_SENDORDER_PG_DSN", "postgres://felo:felo@127.0.0.1:54335/sendorder_db?sslmode=disable"))
	t.Cleanup(func() { db.Close() })

	initSendOrderTables(t, db)

	orderID := "sendorder-ft-001"
	_, _ = db.Exec(ctx, "delete from send_orders where send_order_id=$1", orderID)

	svc := service.NewSendOrderService(
		&pgSendOrderRepo{db: db},
		&noopSendPricing{fee: 20000},
		&noopSendInvoice{},
		&noopSendPublisher{},
		&fixedSendOrderIDs{ids: []string{orderID}},
	)

	order, err := svc.CreateSendOrder(ctx, service.CreateSendOrderInput{
		SenderID:      "sender-ft-001",
		ReceiverPhone: "081234567890",
		Origin:        "loc-a",
		Destination:   "loc-b",
		PackageDetails: domain.PackageDetails{
			WeightKG:    2.5,
			Dimensions:  "30x20x10",
			Description: "buku",
		},
		PayerType: domain.PayerSender,
	})
	if err != nil {
		t.Fatalf("CreateSendOrder() error = %v", err)
	}

	var status string
	if err := db.QueryRow(ctx, "select status from send_orders where send_order_id=$1", order.ID).Scan(&status); err != nil {
		t.Fatalf("query persisted send order: %v", err)
	}
	if status != "created" {
		t.Fatalf("persisted status = %s, want created", status)
	}
}

type pgSendOrderRepo struct{ db *pgxpool.Pool }

func (r *pgSendOrderRepo) Save(ctx context.Context, order domain.SendOrder) error {
	_, err := r.db.Exec(ctx, `insert into send_orders (send_order_id, sender_id, status, created_at)
values ($1,$2,$3,$4)
on conflict (send_order_id) do update set
sender_id=excluded.sender_id,
status=excluded.status`,
		order.ID, order.SenderID, order.Status, order.CreatedAt)
	return err
}

type noopSendPricing struct{ fee int64 }

func (n noopSendPricing) CalculateShippingFee(_ context.Context, _ domain.PackageDetails, _, _ string) (int64, error) {
	return n.fee, nil
}

type noopSendInvoice struct{}

func (noopSendInvoice) CreateInvoice(_ context.Context, _, _ string, _ domain.PayerType, _ int64) error {
	return nil
}

type noopSendPublisher struct{}

func (noopSendPublisher) Publish(_ context.Context, _ domain.Event) error { return nil }

type fixedSendOrderIDs struct {
	ids []string
	idx int
}

func (g *fixedSendOrderIDs) NewID() string {
	id := g.ids[g.idx]
	g.idx++
	return id
}

func initSendOrderTables(t *testing.T, db *pgxpool.Pool) {
	t.Helper()
	ctx := context.Background()
	_, err := db.Exec(ctx, `create table if not exists send_orders (
		send_order_id text primary key,
		sender_id text not null,
		receiver_ref text not null default '',
		shipment_ref text not null default '',
		status text not null,
		created_at timestamptz not null
	)`)
	if err != nil {
		t.Fatalf("initSendOrderTables: %v", err)
	}
}

func openSendOrderPG(t *testing.T, dsn string) *pgxpool.Pool {
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
