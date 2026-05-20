//go:build functional

package functional_test

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/felo/felo-backend/services/cart-service/internal/domain"
	"github.com/felo/felo-backend/services/cart-service/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestCartFunctional_AddItem_PersistsCartToDatabase(t *testing.T) {
	ctx := context.Background()
	db := openCartPG(t, getenv("FELO_CART_PG_DSN", "postgres://felo:felo@127.0.0.1:54327/cart_db?sslmode=disable"))
	t.Cleanup(func() { db.Close() })

	initCartTables(t, db)

	userID := "user-ft-001"
	_, _ = db.Exec(ctx, "delete from carts where user_id=$1", userID)

	svc := service.NewCartService(&pgCartRepo{db: db}, &noopCartMerchant{price: 15000, available: true})

	cart, err := svc.AddItem(ctx, service.AddItemInput{
		UserID:     userID,
		MerchantID: "resto-ft-001",
		MenuItemID: "item-ft-001",
		Quantity:   2,
	})
	if err != nil {
		t.Fatalf("AddItem() error = %v", err)
	}

	var totalPrice int64
	if err := db.QueryRow(ctx, "select total_price from carts where user_id=$1", cart.UserID).Scan(&totalPrice); err != nil {
		t.Fatalf("query persisted cart: %v", err)
	}
	if totalPrice != 30000 {
		t.Fatalf("persisted total_price = %d, want 30000", totalPrice)
	}
}

func TestCartFunctional_ClearCart_RemovesFromDatabase(t *testing.T) {
	ctx := context.Background()
	db := openCartPG(t, getenv("FELO_CART_PG_DSN", "postgres://felo:felo@127.0.0.1:54327/cart_db?sslmode=disable"))
	t.Cleanup(func() { db.Close() })

	initCartTables(t, db)

	userID := "user-ft-002"
	_, _ = db.Exec(ctx, `insert into carts (user_id, merchant_id, items, total_price)
	values ($1,'resto-ft-001','[]',25000) on conflict (user_id) do nothing`, userID)

	svc := service.NewCartService(&pgCartRepo{db: db}, &noopCartMerchant{})
	if err := svc.ClearCart(ctx, userID); err != nil {
		t.Fatalf("ClearCart() error = %v", err)
	}

	var count int
	if err := db.QueryRow(ctx, "select count(*) from carts where user_id=$1", userID).Scan(&count); err != nil {
		t.Fatalf("query deleted cart: %v", err)
	}
	if count != 0 {
		t.Fatal("cart should be deleted from database")
	}
}

type pgCartRepo struct{ db *pgxpool.Pool }

func (r *pgCartRepo) Save(ctx context.Context, cart domain.Cart) error {
	items, _ := json.Marshal(cart.Items)
	_, err := r.db.Exec(ctx, `insert into carts (user_id, merchant_id, items, total_price)
values ($1,$2,$3,$4)
on conflict (user_id) do update set
merchant_id=excluded.merchant_id,
items=excluded.items,
total_price=excluded.total_price`,
		cart.UserID, cart.MerchantID, string(items), cart.TotalPrice)
	return err
}

func (r *pgCartRepo) GetByUserID(ctx context.Context, userID string) (domain.Cart, bool, error) {
	var cart domain.Cart
	var itemsJSON string
	err := r.db.QueryRow(ctx, "select user_id, merchant_id, items, total_price from carts where user_id=$1", userID).
		Scan(&cart.UserID, &cart.MerchantID, &itemsJSON, &cart.TotalPrice)
	if err != nil {
		return domain.Cart{}, false, nil
	}
	if err := json.Unmarshal([]byte(itemsJSON), &cart.Items); err != nil {
		return domain.Cart{}, false, err
	}
	return cart, true, nil
}

func (r *pgCartRepo) Delete(ctx context.Context, userID string) error {
	_, err := r.db.Exec(ctx, "delete from carts where user_id=$1", userID)
	return err
}

type noopCartMerchant struct {
	price     int64
	available bool
}

func (n noopCartMerchant) GetItemPriceAndAvailability(_ context.Context, _, _ string) (int64, bool, error) {
	return n.price, n.available, nil
}

func initCartTables(t *testing.T, db *pgxpool.Pool) {
	t.Helper()
	ctx := context.Background()
	_, err := db.Exec(ctx, `create table if not exists carts (
		user_id text primary key,
		merchant_id text not null,
		items jsonb not null default '[]',
		total_price bigint not null default 0
	)`)
	if err != nil {
		t.Fatalf("initCartTables: %v", err)
	}
}

func getenv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func openCartPG(t *testing.T, dsn string) *pgxpool.Pool {
	t.Helper()
	db, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		t.Fatalf("pgxpool.New() error = %v", err)
	}
	return db
}
