package unit_test

import (
	"context"
	"errors"
	"testing"

	"github.com/felo/felo-backend/services/payment-service/internal/domain"
	"github.com/felo/felo-backend/services/payment-service/internal/service"
)

func TestPaymentService_HandleRideCompleted_DebitsWalletAndIssuesInvoice(t *testing.T) {
	wallets := &walletClientFake{}
	invoices := &invoiceClientFake{invoiceID: "inv-1"}
	processed := &processedStoreFake{results: map[string]domain.PaymentResult{}}
	publisher := &paymentPublisherFake{}
	svc := service.NewPaymentService(wallets, invoices, processed, publisher)

	result, err := svc.HandleRideCompleted(context.Background(), domain.RideCompletedEvent{
		EventID:    "evt-1",
		TripID:     "trip-1",
		CustomerID: "cust-1",
		Amount:     30000,
		Currency:   "IDR",
	})
	if err != nil {
		t.Fatalf("HandleRideCompleted() error = %v", err)
	}

	if result.InvoiceID != "inv-1" {
		t.Fatalf("result.InvoiceID = %s, want inv-1", result.InvoiceID)
	}
	if wallets.calls != 1 {
		t.Fatalf("wallets.calls = %d, want 1", wallets.calls)
	}
	if len(publisher.events) != 1 || publisher.events[0].Name != "payment.completed.v1" {
		t.Fatalf("publisher.events = %#v, want payment.completed.v1", publisher.events)
	}
}

func TestPaymentService_HandleRideCompleted_WalletFailurePublishesFailedEvent(t *testing.T) {
	wallets := &walletClientFake{err: errors.New("insufficient balance")}
	invoices := &invoiceClientFake{invoiceID: "inv-1"}
	publisher := &paymentPublisherFake{}
	svc := service.NewPaymentService(wallets, invoices, &processedStoreFake{results: map[string]domain.PaymentResult{}}, publisher)

	_, err := svc.HandleRideCompleted(context.Background(), domain.RideCompletedEvent{
		EventID:    "evt-1",
		TripID:     "trip-1",
		CustomerID: "cust-1",
		Amount:     30000,
		Currency:   "IDR",
	})
	if err == nil {
		t.Fatal("HandleRideCompleted() error = nil, want error")
	}

	if invoices.calls != 0 {
		t.Fatalf("invoices.calls = %d, want 0", invoices.calls)
	}
	if len(publisher.events) != 1 || publisher.events[0].Name != "payment.failed.v1" {
		t.Fatalf("publisher.events = %#v, want payment.failed.v1", publisher.events)
	}
}

func TestPaymentService_HandleRideCompleted_DuplicateEventReturnsStoredResult(t *testing.T) {
	processed := &processedStoreFake{results: map[string]domain.PaymentResult{
		"evt-1": {EventID: "evt-1", TripID: "trip-1", InvoiceID: "inv-1"},
	}}
	wallets := &walletClientFake{}
	invoices := &invoiceClientFake{invoiceID: "inv-2"}
	svc := service.NewPaymentService(wallets, invoices, processed, &paymentPublisherFake{})

	result, err := svc.HandleRideCompleted(context.Background(), domain.RideCompletedEvent{
		EventID:    "evt-1",
		TripID:     "trip-1",
		CustomerID: "cust-1",
		Amount:     30000,
		Currency:   "IDR",
	})
	if err != nil {
		t.Fatalf("HandleRideCompleted() error = %v", err)
	}

	if result.InvoiceID != "inv-1" {
		t.Fatalf("result.InvoiceID = %s, want inv-1", result.InvoiceID)
	}
	if wallets.calls != 0 {
		t.Fatalf("wallets.calls = %d, want 0", wallets.calls)
	}
}

func TestPaymentService_HandleRideCompleted_InvalidEventReturnsError(t *testing.T) {
	svc := service.NewPaymentService(&walletClientFake{}, &invoiceClientFake{}, &processedStoreFake{results: map[string]domain.PaymentResult{}}, &paymentPublisherFake{})

	_, err := svc.HandleRideCompleted(context.Background(), domain.RideCompletedEvent{})
	if !errors.Is(err, service.ErrInvalidRideCompletedEvent) {
		t.Fatalf("HandleRideCompleted() error = %v, want ErrInvalidRideCompletedEvent", err)
	}
}

type walletClientFake struct {
	calls int
	err   error
}

func (f *walletClientFake) DebitCustomer(_ context.Context, _ string, _ int64, _ string) error {
	f.calls++
	return f.err
}

type invoiceClientFake struct {
	calls     int
	invoiceID string
	err       error
}

func (f *invoiceClientFake) IssueRideInvoice(_ context.Context, _ string, _ string, _ int64, _ string) (string, error) {
	f.calls++
	return f.invoiceID, f.err
}

type processedStoreFake struct {
	results map[string]domain.PaymentResult
}

func (f *processedStoreFake) Get(_ context.Context, eventID string) (domain.PaymentResult, bool, error) {
	result, ok := f.results[eventID]
	return result, ok, nil
}

func (f *processedStoreFake) Save(_ context.Context, result domain.PaymentResult) error {
	f.results[result.EventID] = result
	return nil
}

type paymentPublisherFake struct {
	events []domain.Event
}

func (f *paymentPublisherFake) Publish(_ context.Context, event domain.Event) error {
	f.events = append(f.events, event)
	return nil
}
