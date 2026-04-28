package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/felo/felo-backend/services/order-service/internal/domain"
	"github.com/felo/felo-backend/services/order-service/internal/service"
)

func TestOrderService_CreateOrder_DistanceGreaterThan1KMCash_ReturnsError(t *testing.T) {
	svc := service.NewOrderService(
		&orderRepoFake{},
		&locationClientFake{distance: 1.5},
		&authClientFake{},
		&eventPublisherFake{},
		sequenceIDs("order-1"),
	)

	_, err := svc.CreateOrder(context.Background(), service.CreateOrderInput{
		UserID:        "user-1",
		MerchantID:    "resto-1",
		PaymentMethod: domain.PayCash,
		TotalAmount:   50000,
	})

	if !errors.Is(err, service.ErrDistanceTooFarForCash) {
		t.Fatalf("CreateOrder() error = %v, want ErrDistanceTooFarForCash", err)
	}
}

func TestOrderService_CreateOrder_DistanceGreaterThan1KMWallet_Success(t *testing.T) {
	publisher := &eventPublisherFake{}
	svc := service.NewOrderService(
		&orderRepoFake{orders: map[string]domain.FoodOrder{}},
		&locationClientFake{distance: 1.5},
		&authClientFake{},
		publisher,
		sequenceIDs("order-1"),
	)

	order, err := svc.CreateOrder(context.Background(), service.CreateOrderInput{
		UserID:        "user-1",
		MerchantID:    "resto-1",
		PaymentMethod: domain.PayWallet,
		TotalAmount:   50000,
	})

	if err != nil {
		t.Fatalf("CreateOrder() error = %v", err)
	}
	if order.Status != domain.StateConfirmed {
		t.Fatalf("Status = %s, want confirmed", order.Status)
	}
	if len(publisher.events) != 1 || publisher.events[0].Name != "order.created.v1" {
		t.Fatalf("publisher.events = %#v, want order.created.v1", publisher.events)
	}
}

func TestOrderService_CreateOrder_DistanceLessThan1KMCash_TriggersOTP(t *testing.T) {
	auth := &authClientFake{}
	svc := service.NewOrderService(
		&orderRepoFake{orders: map[string]domain.FoodOrder{}},
		&locationClientFake{distance: 0.8},
		auth,
		&eventPublisherFake{},
		sequenceIDs("order-1"),
	)

	order, err := svc.CreateOrder(context.Background(), service.CreateOrderInput{
		UserID:        "user-1",
		MerchantID:    "resto-1",
		PaymentMethod: domain.PayCash,
		TotalAmount:   50000,
	})

	if err != nil {
		t.Fatalf("CreateOrder() error = %v", err)
	}
	if order.Status != domain.StatePending {
		t.Fatalf("Status = %s, want pending", order.Status)
	}
	if !order.OTPTriggered {
		t.Fatalf("OTPTriggered = false, want true")
	}
	if auth.sendCalls != 1 {
		t.Fatalf("auth.sendCalls = %d, want 1", auth.sendCalls)
	}
}

func TestOrderService_ConfirmCashOrder_ValidOTP_Success(t *testing.T) {
	repo := &orderRepoFake{orders: map[string]domain.FoodOrder{
		"order-1": {ID: "order-1", UserID: "user-1", Status: domain.StatePending, PaymentMethod: domain.PayCash},
	}}
	publisher := &eventPublisherFake{}
	svc := service.NewOrderService(
		repo,
		&locationClientFake{},
		&authClientFake{valid: true},
		publisher,
		sequenceIDs("order-1"),
	)

	order, err := svc.ConfirmCashOrder(context.Background(), "order-1", "1234")
	if err != nil {
		t.Fatalf("ConfirmCashOrder() error = %v", err)
	}
	if order.Status != domain.StateConfirmed {
		t.Fatalf("Status = %s, want confirmed", order.Status)
	}
	if len(publisher.events) != 1 || publisher.events[0].Name != "order.created.v1" {
		t.Fatalf("publisher.events = %#v, want order.created.v1", publisher.events)
	}
}

func TestOrderService_ConfirmCashOrder_InvalidOTP_ReturnsError(t *testing.T) {
	repo := &orderRepoFake{orders: map[string]domain.FoodOrder{
		"order-1": {ID: "order-1", UserID: "user-1", Status: domain.StatePending, PaymentMethod: domain.PayCash},
	}}
	svc := service.NewOrderService(
		repo,
		&locationClientFake{},
		&authClientFake{valid: false},
		&eventPublisherFake{},
		sequenceIDs("order-1"),
	)

	_, err := svc.ConfirmCashOrder(context.Background(), "order-1", "0000")
	if !errors.Is(err, service.ErrInvalidOTP) {
		t.Fatalf("ConfirmCashOrder() error = %v, want ErrInvalidOTP", err)
	}
}

type orderRepoFake struct {
	orders map[string]domain.FoodOrder
}

func (f *orderRepoFake) Save(_ context.Context, order domain.FoodOrder) error {
	f.orders[order.ID] = order
	return nil
}

func (f *orderRepoFake) GetByID(_ context.Context, orderID string) (domain.FoodOrder, error) {
	order, ok := f.orders[orderID]
	if !ok {
		return domain.FoodOrder{}, errors.New("not found")
	}
	return order, nil
}

type locationClientFake struct {
	distance float64
	err      error
}

func (f *locationClientFake) GetDistanceKM(_ context.Context, _ string, _ string) (float64, error) {
	return f.distance, f.err
}

type authClientFake struct {
	sendCalls   int
	verifyCalls int
	valid       bool
	err         error
}

func (f *authClientFake) SendOTP(_ context.Context, _ string) error {
	f.sendCalls++
	return f.err
}

func (f *authClientFake) VerifyOTP(_ context.Context, _ string, _ string) (bool, error) {
	f.verifyCalls++
	return f.valid, f.err
}

type eventPublisherFake struct {
	events []domain.Event
}

func (f *eventPublisherFake) Publish(_ context.Context, event domain.Event) error {
	f.events = append(f.events, event)
	return nil
}

type sequenceIDGenerator struct {
	values []string
	index  int
}

func sequenceIDs(values ...string) *sequenceIDGenerator {
	return &sequenceIDGenerator{values: values}
}

func (g *sequenceIDGenerator) NewID() string {
	if g.index >= len(g.values) {
		return "generated-id"
	}
	value := g.values[g.index]
	g.index++
	return value
}
