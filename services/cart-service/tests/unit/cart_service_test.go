package unit_test

import (
	"context"
	"errors"
	"testing"

	"github.com/felo/felo-backend/services/cart-service/internal/domain"
	"github.com/felo/felo-backend/services/cart-service/internal/service"
)

func TestCartService_AddItem_CreatesNewCart(t *testing.T) {
	repo := &cartRepoFake{carts: map[string]domain.Cart{}}
	merchant := &merchantClientFake{price: 15000, available: true}
	svc := service.NewCartService(repo, merchant)

	cart, err := svc.AddItem(context.Background(), service.AddItemInput{
		UserID:     "user-1",
		MerchantID: "resto-1",
		MenuItemID: "item-1",
		Quantity:   2,
	})
	if err != nil {
		t.Fatalf("AddItem() error = %v", err)
	}

	if cart.TotalPrice != 30000 {
		t.Fatalf("TotalPrice = %d, want 30000", cart.TotalPrice)
	}
	if len(cart.Items) != 1 {
		t.Fatalf("len(cart.Items) = %d, want 1", len(cart.Items))
	}
}

func TestCartService_AddItem_ItemUnavailableReturnsError(t *testing.T) {
	svc := service.NewCartService(&cartRepoFake{carts: map[string]domain.Cart{}}, &merchantClientFake{available: false})

	_, err := svc.AddItem(context.Background(), service.AddItemInput{
		UserID:     "user-1",
		MerchantID: "resto-1",
		MenuItemID: "item-1",
		Quantity:   1,
	})
	if !errors.Is(err, service.ErrItemUnavailable) {
		t.Fatalf("AddItem() error = %v, want ErrItemUnavailable", err)
	}
}

func TestCartService_AddItem_MultipleMerchantsReturnsError(t *testing.T) {
	repo := &cartRepoFake{carts: map[string]domain.Cart{
		"user-1": {UserID: "user-1", MerchantID: "resto-1", Items: []domain.CartItem{{}}, TotalPrice: 15000},
	}}
	svc := service.NewCartService(repo, &merchantClientFake{price: 20000, available: true})

	_, err := svc.AddItem(context.Background(), service.AddItemInput{
		UserID:     "user-1",
		MerchantID: "resto-2",
		MenuItemID: "item-2",
		Quantity:   1,
	})
	if !errors.Is(err, service.ErrMultipleMerchants) {
		t.Fatalf("AddItem() error = %v, want ErrMultipleMerchants", err)
	}
}

func TestCartService_AddItem_InvalidInputReturnsError(t *testing.T) {
	svc := service.NewCartService(&cartRepoFake{}, &merchantClientFake{})
	_, err := svc.AddItem(context.Background(), service.AddItemInput{})
	if !errors.Is(err, service.ErrInvalidInput) {
		t.Fatalf("AddItem() error = %v, want ErrInvalidInput", err)
	}
}

func TestCartService_GetCartSummary_ReturnsCorrectSummary(t *testing.T) {
	repo := &cartRepoFake{carts: map[string]domain.Cart{
		"user-1": {
			UserID:     "user-1",
			MerchantID: "resto-1",
			Items:      []domain.CartItem{{Price: 10000, Quantity: 1}, {Price: 20000, Quantity: 1}},
			TotalPrice: 30000,
		},
	}}
	svc := service.NewCartService(repo, &merchantClientFake{})

	summary, err := svc.GetCartSummary(context.Background(), "user-1")
	if err != nil {
		t.Fatalf("GetCartSummary() error = %v", err)
	}
	if summary.ItemCount != 2 {
		t.Fatalf("ItemCount = %d, want 2", summary.ItemCount)
	}
	if summary.TotalPrice != 30000 {
		t.Fatalf("TotalPrice = %d, want 30000", summary.TotalPrice)
	}
}

func TestCartService_GetCartSummary_EmptyCart(t *testing.T) {
	svc := service.NewCartService(&cartRepoFake{carts: map[string]domain.Cart{}}, &merchantClientFake{})
	summary, err := svc.GetCartSummary(context.Background(), "user-1")
	if err != nil {
		t.Fatalf("GetCartSummary() error = %v", err)
	}
	if summary.ItemCount != 0 {
		t.Fatalf("ItemCount = %d, want 0", summary.ItemCount)
	}
}

func TestCartService_ClearCart_DeletesFromRepo(t *testing.T) {
	repo := &cartRepoFake{carts: map[string]domain.Cart{"user-1": {}}}
	svc := service.NewCartService(repo, &merchantClientFake{})

	if err := svc.ClearCart(context.Background(), "user-1"); err != nil {
		t.Fatalf("ClearCart() error = %v", err)
	}
	if _, ok := repo.carts["user-1"]; ok {
		t.Fatalf("cart should be deleted")
	}
}

type cartRepoFake struct {
	carts map[string]domain.Cart
}

func (f *cartRepoFake) Save(_ context.Context, cart domain.Cart) error {
	f.carts[cart.UserID] = cart
	return nil
}

func (f *cartRepoFake) GetByUserID(_ context.Context, userID string) (domain.Cart, bool, error) {
	cart, ok := f.carts[userID]
	return cart, ok, nil
}

func (f *cartRepoFake) Delete(_ context.Context, userID string) error {
	delete(f.carts, userID)
	return nil
}

type merchantClientFake struct {
	price     int64
	available bool
	err       error
}

func (f *merchantClientFake) GetItemPriceAndAvailability(_ context.Context, _ string, _ string) (int64, bool, error) {
	return f.price, f.available, f.err
}
