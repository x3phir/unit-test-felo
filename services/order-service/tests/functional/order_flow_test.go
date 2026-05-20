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
	db := openOrderPG(t, getenv("FELO_ORDER_PG_DSN", "postgres://felo:felo@127.0.0.1:54326/order_db?sslmode=disable"))
	t.Cleanup(func() { db.Close() })

	initOrderTables(t, db)

	orderID := "order-ft-001"
	_, _ = db.Exec(ctx, "delete from food_orders where id=$1", orderID)

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
	var amount int64
	if err := db.QueryRow(ctx, "select status, total_amount from food_orders where id=$1", order.ID).Scan(&status, &amount); err != nil {
		t.Fatalf("query persisted order: %v", err)
	}
	if status != string(domain.StateConfirmed) {
		t.Fatalf("persisted status = %s, want %s", status, domain.StateConfirmed)
	}
	if amount != 50000 {
		t.Fatalf("persisted total_amount = %d, want 50000", amount)
	}
}

func TestOrderFunctional_CreateOrderCash_PersistsPendingOrder(t *testing.T) {
	ctx := context.Background()
	db := openOrderPG(t, getenv("FELO_ORDER_PG_DSN", "postgres://felo:felo@127.0.0.1:54326/order_db?sslmode=disable"))
	t.Cleanup(func() { db.Close() })

	initOrderTables(t, db)

	orderID := "order-ft-002"
	_, _ = db.Exec(ctx, "delete from food_orders where id=$1", orderID)

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
	var otpTriggered bool
	if err := db.QueryRow(ctx, "select status, otp_triggered from food_orders where id=$1", order.ID).Scan(&status, &otpTriggered); err != nil {
		t.Fatalf("query persisted order: %v", err)
	}
	if status != string(domain.StatePending) {
		t.Fatalf("persisted status = %s, want %s", status, domain.StatePending)
	}
	if !otpTriggered {
		t.Fatal("persisted otp_triggered = false, want true")
	}
}

type pgOrderRepo struct{ db *pgxpool.Pool }

func (r *pgOrderRepo) Save(ctx context.Context, order domain.FoodOrder) error {
	_, err := r.db.Exec(ctx, `insert into food_orders (id, user_id, merchant_id, items, payment_method, status, distance_km, total_amount, otp_triggered, created_at, updated_at)
values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
on conflict (id) do update set
user_id=excluded.user_id,
merchant_id=excluded.merchant_id,
items=excluded.items,
payment_method=excluded.payment_method,
status=excluded.status,
distance_km=excluded.distance_km,
total_amount=excluded.total_amount,
otp_triggered=excluded.otp_triggered,
updated_at=excluded.updated_at`,
		order.ID, order.UserID, order.MerchantID, order.Items, string(order.PaymentMethod),
		string(order.Status), order.DistanceKM, order.TotalAmount, order.OTPTriggered,
		order.CreatedAt, order.UpdatedAt)
	return err
}

func (r *pgOrderRepo) GetByID(ctx context.Context, orderID string) (domain.FoodOrder, error) {
	var order domain.FoodOrder
	var paymentMethod, status string
	err := r.db.QueryRow(ctx,
		`select id, user_id, merchant_id, items, payment_method, status, distance_km, total_amount, otp_triggered, created_at, updated_at
from food_orders where id=$1`, orderID).
		Scan(&order.ID, &order.UserID, &order.MerchantID, &order.Items,
			&paymentMethod, &status, &order.DistanceKM, &order.TotalAmount,
			&order.OTPTriggered, &order.CreatedAt, &order.UpdatedAt)
	if err != nil {
		return domain.FoodOrder{}, err
	}
	order.PaymentMethod = domain.PaymentMethod(paymentMethod)
	order.Status = domain.OrderState(status)
	return order, nil
}

type noopOrderLocation struct{ distance float64 }

func (n noopOrderLocation) GetDistanceKM(_ context.Context, _, _ string) (float64, error) {
	return n.distance, nil
}

type noopOrderAuth struct{}

func (noopOrderAuth) SendOTP(_ context.Context, _ string) error { return nil }
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
	_, err := db.Exec(ctx, `create table if not exists food_orders (
		id text primary key,
		user_id text not null,
		merchant_id text not null,
		items text[] not null default '{}',
		payment_method text not null,
		status text not null,
		distance_km double precision not null default 0,
		total_amount bigint not null,
		otp_triggered boolean not null default false,
		created_at timestamptz not null,
		updated_at timestamptz not null
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
