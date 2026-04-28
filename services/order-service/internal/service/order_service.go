package service

import (
	"context"
	"errors"
	"time"

	"github.com/felo/felo-backend/services/order-service/internal/domain"
	"github.com/felo/felo-backend/services/order-service/internal/ports"
)

var (
	ErrDistanceTooFarForCash = errors.New("distance > 1 KM requires Wallet payment")
	ErrInvalidInput          = errors.New("invalid input")
	ErrInvalidOTP            = errors.New("invalid OTP")
	ErrInvalidStateTransition = errors.New("invalid state transition")
)

type OrderService struct {
	repo      ports.OrderRepository
	location  ports.LocationClient
	auth      ports.AuthClient
	publisher ports.EventPublisher
	ids       ports.IDGenerator
	now       func() time.Time
}

func NewOrderService(repo ports.OrderRepository, location ports.LocationClient, auth ports.AuthClient, publisher ports.EventPublisher, ids ports.IDGenerator) *OrderService {
	return &OrderService{
		repo:      repo,
		location:  location,
		auth:      auth,
		publisher: publisher,
		ids:       ids,
		now:       time.Now,
	}
}

type CreateOrderInput struct {
	UserID        string
	MerchantID    string
	UserLocation  string
	RestoLocation string
	PaymentMethod domain.PaymentMethod
	TotalAmount   int64
}

func (s *OrderService) CreateOrder(ctx context.Context, input CreateOrderInput) (domain.FoodOrder, error) {
	if input.UserID == "" || input.MerchantID == "" || input.TotalAmount <= 0 {
		return domain.FoodOrder{}, ErrInvalidInput
	}

	distance, err := s.location.GetDistanceKM(ctx, input.UserLocation, input.RestoLocation)
	if err != nil {
		return domain.FoodOrder{}, err
	}

	if distance > 1.0 && input.PaymentMethod == domain.PayCash {
		return domain.FoodOrder{}, ErrDistanceTooFarForCash
	}

	now := s.now()
	order := domain.FoodOrder{
		ID:            s.ids.NewID(),
		UserID:        input.UserID,
		MerchantID:    input.MerchantID,
		PaymentMethod: input.PaymentMethod,
		Status:        domain.StatePending,
		DistanceKM:    distance,
		TotalAmount:   input.TotalAmount,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	if input.PaymentMethod == domain.PayCash {
		order.OTPTriggered = true
		if err := s.auth.SendOTP(ctx, input.UserID); err != nil {
			return domain.FoodOrder{}, err
		}
	} else {
		order.Status = domain.StateConfirmed
		if err := s.publisher.Publish(ctx, domain.Event{Name: "order.created.v1", OrderID: order.ID, OccurredAt: now}); err != nil {
			return domain.FoodOrder{}, err
		}
	}

	if err := s.repo.Save(ctx, order); err != nil {
		return domain.FoodOrder{}, err
	}

	return order, nil
}

func (s *OrderService) ConfirmCashOrder(ctx context.Context, orderID string, otp string) (domain.FoodOrder, error) {
	order, err := s.repo.GetByID(ctx, orderID)
	if err != nil {
		return domain.FoodOrder{}, err
	}
	if order.Status != domain.StatePending || order.PaymentMethod != domain.PayCash {
		return domain.FoodOrder{}, ErrInvalidStateTransition
	}

	valid, err := s.auth.VerifyOTP(ctx, order.UserID, otp)
	if err != nil {
		return domain.FoodOrder{}, err
	}
	if !valid {
		return domain.FoodOrder{}, ErrInvalidOTP
	}

	now := s.now()
	order.Status = domain.StateConfirmed
	order.UpdatedAt = now

	if err := s.repo.Save(ctx, order); err != nil {
		return domain.FoodOrder{}, err
	}
	if err := s.publisher.Publish(ctx, domain.Event{Name: "order.created.v1", OrderID: order.ID, OccurredAt: now}); err != nil {
		return domain.FoodOrder{}, err
	}

	return order, nil
}
