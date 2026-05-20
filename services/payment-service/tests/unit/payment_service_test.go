package unit_test

import (
	"context"
	"testing"

	"github.com/felo/felo-backend/services/payment-service/internal/domain"
	"github.com/felo/felo-backend/services/payment-service/internal/service"
	"go.uber.org/mock/gomock"
)

func TestPaymentService_HandleRideCompleted_CompletesPaymentWithGomock(t *testing.T) {
	ctrl := gomock.NewController(t)
	wallets := NewMockWalletClient(ctrl)
	invoices := NewMockInvoiceClient(ctrl)
	processed := NewMockProcessedEventStore(ctrl)
	publisher := NewMockPaymentEventPublisher(ctrl)
	svc := service.NewPaymentService(wallets, invoices, processed, publisher)

	processed.EXPECT().Get(gomock.Any(), "evt-1").Return(domain.PaymentResult{}, false, nil)
	wallets.EXPECT().DebitCustomer(gomock.Any(), "cust-1", int64(30000), "evt-1").Return(nil)
	invoices.EXPECT().IssueRideInvoice(gomock.Any(), "ride-1", "cust-1", int64(30000), "IDR").Return("inv-1", nil)
	processed.EXPECT().Save(gomock.Any(), gomock.AssignableToTypeOf(domain.PaymentResult{})).Return(nil)
	publisher.EXPECT().Publish(gomock.Any(), gomock.AssignableToTypeOf(domain.Event{})).Return(nil)

	result, err := svc.HandleRideCompleted(context.Background(), domain.RideCompletedEvent{
		EventID: "evt-1", TripID: "ride-1", CustomerID: "cust-1", Amount: 30000, Currency: "IDR",
	})
	if err != nil {
		t.Fatalf("HandleRideCompleted() error = %v", err)
	}
	if result.InvoiceID != "inv-1" {
		t.Fatalf("result.InvoiceID = %s, want inv-1", result.InvoiceID)
	}
}
