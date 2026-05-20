//go:build functional

package functional_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/felo/felo-backend/services/cart-service/internal/domain"
	"github.com/felo/felo-backend/services/cart-service/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestCartFunctional_AddItem_PersistsCartAndItemToDatabase(t *testing.T) {
	ctx := context.Background()
	db := openCartPG(t, getenv("FELO_CART_PG_DSN", "postgres://felo:felo@127.0.0.1:54334/cart_db?sslmode=disable"))
	t.Cleanup(func() { db.Close() })

	initCartTables(t, db)

	_, _ = db.Exec(ctx, "delete from cart_items where cart_id in (select cart_id from carts where customer_id=$1)", "user-ft-001")
	_, _ = db.Exec(ctx, "delete from carts where customer_id=$1", "user-ft-001")

	svc := service.NewCartService(&pgCartRepo{db: db}, &noopCartMerchant{price: 15000, available: true})

	cart, err := svc.AddItem(ctx, service.AddItemInput{
		UserID:     "user-ft-001",
		MerchantID: "resto-ft-001",
		MenuItemID: "item-ft-001",
		Quantity:   2,
	})
	if err != nil {
		t.Fatalf("AddItem() error = %v", err)
	}

	var status string
	if err := db.QueryRow(ctx, "select status from carts where customer_id=$1", cart.UserID).Scan(&status); err != nil {
		t.Fatalf("query persisted cart: %v", err)
	}
	if status != "active" {
		t.Fatalf("cart status = %s, want active", status)
	}

	var qty int
	if err := db.QueryRow(ctx,
		"select ci.quantity from cart_items ci join carts c on ci.cart_id=c.cart_id where c.customer_id=$1",
		cart.UserID).Scan(&qty); err != nil {
		t.Fatalf("query cart_item: %v", err)
	}
	if qty != 2 {
		t.Fatalf("cart_item quantity = %d, want 2", qty)
	}
}

func TestCartFunctional_ClearCart_RemovesCartAndItems(t *testing.T) {
	ctx := context.Background()
	db := openCartPG(t, getenv("FELO_CART_PG_DSN", "postgres://felo:felo@127.0.0.1:54334/cart_db?sslmode=disable"))
	t.Cleanup(func() { db.Close() })

	initCartTables(t, db)

	cartID := "cart-ft-001"
	_, _ = db.Exec(ctx, "delete from cart_items where cart_id=$1", cartID)
	_, _ = db.Exec(ctx, "delete from carts where cart_id=$1", cartID)
	_, _ = db.Exec(ctx, "insert into carts (cart_id, customer_id, merchant_id, status, updated_at) values ($1,'user-ft-002','resto-ft-001','active',$2)",
		cartID, time.Now().UTC())

	svc := service.NewCartService(&pgCartRepo{db: db}, &noopCartMerchant{})
	if err := svc.ClearCart(ctx, "user-ft-002"); err != nil {
		t.Fatalf("ClearCart() error = %v", err)
	}

	var count int
	if err := db.QueryRow(ctx, "select count(*) from carts where customer_id=$1", "user-ft-002").Scan(&count); err != nil {
		t.Fatalf("query deleted cart: %v", err)
	}
	if count != 0 {
		t.Fatal("cart should be deleted from database")
	}
}

type pgCartRepo struct{ db *pgxpool.Pool }

func (r *pgCartRepo) Save(ctx context.Context, cart domain.Cart) error {
	cartID := "cart-" + cart.UserID
	_, err := r.db.Exec(ctx, `insert into carts (cart_id, customer_id, merchant_id, status, updated_at)
values ($1,$2,$3,'active',$4)
on conflict (cart_id) do update set
merchant_id=excluded.merchant_id,
status=excluded.status,
updated_at=excluded.updated_at`,
		cartID, cart.UserID, cart.MerchantID, time.Now().UTC())
	if err != nil {
		return err
	}

	for _, item := range cart.Items {
		itemID := cartID + "-" + item.MenuItemID
		_, err := r.db.Exec(ctx, `insert into cart_items (item_id, cart_id, product_ref, quantity)
values ($1,$2,$3,$4)
on conflict (item_id) do update set
quantity=excluded.quantity`,
			itemID, cartID, item.MenuItemID, item.Quantity)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *pgCartRepo) GetByUserID(ctx context.Context, userID string) (domain.Cart, bool, error) {
	var cart domain.Cart
	var cartID string
	err := r.db.QueryRow(ctx,
		"select cart_id, customer_id, merchant_id from carts where customer_id=$1", userID).
		Scan(&cartID, &cart.UserID, &cart.MerchantID)
	if err != nil {
		return domain.Cart{}, false, nil
	}

	rows, err := r.db.Query(ctx,
		"select product_ref, quantity from cart_items where cart_id=$1", cartID)
	if err != nil {
		return domain.Cart{}, false, err
	}
	defer rows.Close()

	var totalPrice int64
	for rows.Next() {
		var item domain.CartItem
		if err := rows.Scan(&item.MenuItemID, &item.Quantity); err != nil {
			return domain.Cart{}, false, err
		}
		if item.Quantity > 0 {
			item.Price = 15000
		}
		totalPrice += item.Price * int64(item.Quantity)
		cart.Items = append(cart.Items, item)
	}
	cart.TotalPrice = totalPrice

	return cart, true, nil
}

func (r *pgCartRepo) Delete(ctx context.Context, userID string) error {
	cartID := "cart-" + userID
	_, _ = r.db.Exec(ctx, "delete from cart_items where cart_id=$1", cartID)
	_, err := r.db.Exec(ctx, "delete from carts where cart_id=$1", cartID)
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
		cart_id text primary key,
		customer_id text not null,
		merchant_id text not null,
		status text not null default 'active',
		updated_at timestamptz not null
	)`)
	if err != nil {
		t.Fatalf("initCartTables carts: %v", err)
	}
	_, err = db.Exec(ctx, `create table if not exists cart_items (
		item_id text primary key,
		cart_id text not null,
		product_ref text not null,
		quantity integer not null
	)`)
	if err != nil {
		t.Fatalf("initCartTables cart_items: %v", err)
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
