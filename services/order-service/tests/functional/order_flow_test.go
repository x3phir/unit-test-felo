//go:build functional

package functional_test

import (
	"context"
	"os"
	"testing"

	"github.com/felo/felo-backend/services/order-service/internal/domain"
	"github.com/felo/felo-backend/services/order-service/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestOrderFunctional_CreateOrderWallet_PersistsConfirmedOrder(t *testing.T) {
	ctx := context.Background()
	db := openOrderPG(t, getenv("FELO_ORDER_PG_DSN", "postgres://felo:felo@127.0.0.1:54333/order_db?sslmode=disable"))
	t.Cleanup(func() { db.Close() })

	initOrderTables(t, db)

	orderID := "order-ft-001"
	_, _ = db.Exec(ctx, "delete from orders where order_id=$1", orderID)

	svc := service.NewOrderService(
		&pgOrderRepo{db: db},
		&noopOrderLocation{},
		&noopOrderAuth{},
		&noopOrderPublisher{},
		&fixedOrderIDs{ids: []string{orderID}},
	)

	order, err := svc.CreateOrder(ctx, service.CreateOrderInput{
		UserID:        "cust-ft-001",
		MerchantID:    "resto-ft-001",
		UserLocation:  "loc-a",
		RestoLocation: "loc-b",
		PaymentMethod: domain.PayWallet,
		TotalAmount:   50000,
	})
	if err != nil {
		t.Fatalf("CreateOrder() error = %v", err)
	}

	var status string
	if err := db.QueryRow(ctx, "select status from orders where order_id=$1", order.ID).Scan(&status); err != nil {
		t.Fatalf("query persisted order: %v", err)
	}
	if status != string(domain.StateConfirmed) {
		t.Fatalf("persisted status = %s, want %s", status, domain.StateConfirmed)
	}
}

func TestOrderFunctional_CreateOrderCash_PersistsPendingOrder(t *testing.T) {
	ctx := context.Background()
	db := openOrderPG(t, getenv("FELO_ORDER_PG_DSN", "postgres://felo:felo@127.0.0.1:54333/order_db?sslmode=disable"))
	t.Cleanup(func() { db.Close() })

	initOrderTables(t, db)

	orderID := "order-ft-002"
	_, _ = db.Exec(ctx, "delete from orders where order_id=$1", orderID)

	svc := service.NewOrderService(
		&pgOrderRepo{db: db},
		&noopOrderLocation{distance: 0.5},
		&noopOrderAuth{},
		&noopOrderPublisher{},
		&fixedOrderIDs{ids: []string{orderID}},
	)

	order, err := svc.CreateOrder(ctx, service.CreateOrderInput{
		UserID:        "cust-ft-001",
		MerchantID:    "resto-ft-001",
		UserLocation:  "loc-a",
		RestoLocation: "loc-b",
		PaymentMethod: domain.PayCash,
		TotalAmount:   30000,
	})
	if err != nil {
		t.Fatalf("CreateOrder() error = %v", err)
	}

	var status string
	if err := db.QueryRow(ctx, "select status from orders where order_id=$1", order.ID).Scan(&status); err != nil {
		t.Fatalf("query persisted order: %v", err)
	}
	if status != string(domain.StatePending) {
		t.Fatalf("persisted status = %s, want %s", status, domain.StatePending)
	}
}

type pgOrderRepo struct{ db *pgxpool.Pool }

func (r *pgOrderRepo) Save(ctx context.Context, order domain.FoodOrder) error {
	_, err := r.db.Exec(ctx, `insert into orders (order_id, customer_id, merchant_id, status, created_at)
values ($1,$2,$3,$4,$5)
on conflict (order_id) do update set
customer_id=excluded.customer_id,
merchant_id=excluded.merchant_id,
status=excluded.status`,
		order.ID, order.UserID, order.MerchantID, string(order.Status), order.CreatedAt)
	return err
}

func (r *pgOrderRepo) GetByID(ctx context.Context, orderID string) (domain.FoodOrder, error) {
	var order domain.FoodOrder
	var status string
	err := r.db.QueryRow(ctx,
		"select order_id, customer_id, merchant_id, status, created_at from orders where order_id=$1", orderID).
		Scan(&order.ID, &order.UserID, &order.MerchantID, &status, &order.CreatedAt)
	if err != nil {
		return domain.FoodOrder{}, err
	}
	order.Status = domain.OrderState(status)
	return order, nil
}

type noopOrderLocation struct{ distance float64 }

func (n noopOrderLocation) GetDistanceKM(_ context.Context, _, _ string) (float64, error) {
	return n.distance, nil
}

type noopOrderAuth struct{}

func (noopOrderAuth) SendOTP(_ context.Context, _ string) error          { return nil }
func (noopOrderAuth) VerifyOTP(_ context.Context, _, _ string) (bool, error) { return true, nil }

type noopOrderPublisher struct{}

func (noopOrderPublisher) Publish(_ context.Context, _ domain.Event) error { return nil }

type fixedOrderIDs struct {
	ids []string
	idx int
}

func (g *fixedOrderIDs) NewID() string {
	id := g.ids[g.idx]
	g.idx++
	return id
}

func initOrderTables(t *testing.T, db *pgxpool.Pool) {
	t.Helper()
	ctx := context.Background()
	_, err := db.Exec(ctx, `create table if not exists orders (
		order_id text primary key,
		customer_id text not null,
		merchant_id text not null,
		shipment_ref text not null default '',
		status text not null,
		created_at timestamptz not null
	)`)
	if err != nil {
		t.Fatalf("initOrderTables: %v", err)
	}
}

func openOrderPG(t *testing.T, dsn string) *pgxpool.Pool {
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
