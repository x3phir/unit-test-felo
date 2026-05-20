package unit_test

import (
	"context"
	"errors"
	"testing"

	"github.com/felo/felo-backend/services/send-order-service/internal/domain"
	"github.com/felo/felo-backend/services/send-order-service/internal/service"
)

func TestSendOrderService_CreateSendOrder_ValidInput_Success(t *testing.T) {
	invoice := &invoiceClientFake{}
	publisher := &eventPublisherFake{}
	svc := service.NewSendOrderService(
		&sendOrderRepoFake{},
		&pricingClientFake{fee: 15000},
		invoice,
		publisher,
		sequenceIDs("order-1"),
	)

	order, err := svc.CreateSendOrder(context.Background(), service.CreateSendOrderInput{
		SenderID:      "sender-1",
		ReceiverPhone: "081234567890",
		Origin:        "loc-a",
		Destination:   "loc-b",
		PackageDetails: domain.PackageDetails{WeightKG: 2.5},
		PayerType:     domain.PayerSender,
	})

	if err != nil {
		t.Fatalf("CreateSendOrder() error = %v", err)
	}

	if order.ShippingFee != 15000 {
		t.Fatalf("ShippingFee = %d, want 15000", order.ShippingFee)
	}
	if invoice.calls != 1 {
		t.Fatalf("invoice.calls = %d, want 1", invoice.calls)
	}
	if invoice.lastPayerID != "sender-1" {
		t.Fatalf("invoice.lastPayerID = %s, want sender-1", invoice.lastPayerID)
	}
	if len(publisher.events) != 1 || publisher.events[0].Name != "shipment.created.v1" {
		t.Fatalf("publisher.events = %#v, want shipment.created.v1", publisher.events)
	}
}

func TestSendOrderService_CreateSendOrder_ReceiverPayer_UsesReceiverPhoneForInvoice(t *testing.T) {
	invoice := &invoiceClientFake{}
	svc := service.NewSendOrderService(
		&sendOrderRepoFake{},
		&pricingClientFake{fee: 20000},
		invoice,
		&eventPublisherFake{},
		sequenceIDs("order-1"),
	)

	_, err := svc.CreateSendOrder(context.Background(), service.CreateSendOrderInput{
		SenderID:      "sender-1",
		ReceiverPhone: "081234567890",
		Origin:        "loc-a",
		Destination:   "loc-b",
		PackageDetails: domain.PackageDetails{WeightKG: 2.5},
		PayerType:     domain.PayerReceiver,
	})

	if err != nil {
		t.Fatalf("CreateSendOrder() error = %v", err)
	}
	if invoice.lastPayerID != "081234567890" {
		t.Fatalf("invoice.lastPayerID = %s, want 081234567890", invoice.lastPayerID)
	}
}

func TestSendOrderService_CreateSendOrder_InvalidInput_ReturnsError(t *testing.T) {
	svc := service.NewSendOrderService(&sendOrderRepoFake{}, &pricingClientFake{}, &invoiceClientFake{}, &eventPublisherFake{}, sequenceIDs("order-1"))

	_, err := svc.CreateSendOrder(context.Background(), service.CreateSendOrderInput{})
	if !errors.Is(err, service.ErrInvalidInput) {
		t.Fatalf("CreateSendOrder() error = %v, want ErrInvalidInput", err)
	}
}

type sendOrderRepoFake struct {
	orders map[string]domain.SendOrder
}

func (f *sendOrderRepoFake) Save(_ context.Context, order domain.SendOrder) error {
	if f.orders == nil {
		f.orders = make(map[string]domain.SendOrder)
	}
	f.orders[order.ID] = order
	return nil
}

type pricingClientFake struct {
	fee int64
	err error
}

func (f *pricingClientFake) CalculateShippingFee(_ context.Context, _ domain.PackageDetails, _ string, _ string) (int64, error) {
	return f.fee, f.err
}

type invoiceClientFake struct {
	calls       int
	lastPayerID string
	err         error
}

func (f *invoiceClientFake) CreateInvoice(_ context.Context, _ string, payerID string, _ domain.PayerType, _ int64) error {
	f.calls++
	f.lastPayerID = payerID
	return f.err
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
