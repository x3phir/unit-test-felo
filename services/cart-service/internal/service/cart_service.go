package service

import (
	"context"
	"errors"

	"github.com/felo/felo-backend/services/cart-service/internal/domain"
	"github.com/felo/felo-backend/services/cart-service/internal/ports"
)

var (
	ErrItemUnavailable   = errors.New("menu item is unavailable")
	ErrInvalidInput      = errors.New("invalid input")
	ErrMultipleMerchants = errors.New("cannot add items from different merchants to the same cart")
)

type CartService struct {
	repo     ports.CartRepository
	merchant ports.MerchantClient
}

func NewCartService(repo ports.CartRepository, merchant ports.MerchantClient) *CartService {
	return &CartService{
		repo:     repo,
		merchant: merchant,
	}
}

type AddItemInput struct {
	UserID     string
	MerchantID string
	MenuItemID string
	Quantity   int
	Notes      string
}

func (s *CartService) AddItem(ctx context.Context, input AddItemInput) (domain.Cart, error) {
	if input.UserID == "" || input.MerchantID == "" || input.MenuItemID == "" || input.Quantity <= 0 {
		return domain.Cart{}, ErrInvalidInput
	}

	price, available, err := s.merchant.GetItemPriceAndAvailability(ctx, input.MerchantID, input.MenuItemID)
	if err != nil {
		return domain.Cart{}, err
	}
	if !available {
		return domain.Cart{}, ErrItemUnavailable
	}

	cart, found, err := s.repo.GetByUserID(ctx, input.UserID)
	if err != nil {
		return domain.Cart{}, err
	}

	if !found {
		cart = domain.Cart{
			UserID:     input.UserID,
			MerchantID: input.MerchantID,
			Items:      []domain.CartItem{},
		}
	} else if cart.MerchantID != input.MerchantID {
		return domain.Cart{}, ErrMultipleMerchants
	}

	cart.Items = append(cart.Items, domain.CartItem{
		MenuItemID: input.MenuItemID,
		Quantity:   input.Quantity,
		Notes:      input.Notes,
		Price:      price,
	})

	cart.TotalPrice += price * int64(input.Quantity)

	if err := s.repo.Save(ctx, cart); err != nil {
		return domain.Cart{}, err
	}

	return cart, nil
}

func (s *CartService) GetCartSummary(ctx context.Context, userID string) (domain.CartSummary, error) {
	cart, found, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return domain.CartSummary{}, err
	}
	if !found {
		return domain.CartSummary{UserID: userID}, nil
	}

	return domain.CartSummary{
		UserID:     cart.UserID,
		MerchantID: cart.MerchantID,
		ItemCount:  len(cart.Items),
		TotalPrice: cart.TotalPrice,
	}, nil
}

func (s *CartService) ClearCart(ctx context.Context, userID string) error {
	return s.repo.Delete(ctx, userID)
}
