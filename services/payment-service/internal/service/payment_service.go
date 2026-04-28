package service

import (
	"context"
	"errors"
	"time"

	"github.com/felo/felo-backend/services/payment-service/internal/domain"
	"github.com/felo/felo-backend/services/payment-service/internal/ports"
)

var ErrInvalidRideCompletedEvent = errors.New("invalid ride completed event")

type PaymentService struct {
	wallets   ports.WalletClient
	invoices  ports.InvoiceClient
	processed ports.ProcessedEventStore
	publisher ports.EventPublisher
	now       func() time.Time
}

func NewPaymentService(wallets ports.WalletClient, invoices ports.InvoiceClient, processed ports.ProcessedEventStore, publisher ports.EventPublisher) *PaymentService {
	return &PaymentService{
		wallets:   wallets,
		invoices:  invoices,
		processed: processed,
		publisher: publisher,
		now:       time.Now,
	}
}

func (s *PaymentService) HandleRideCompleted(ctx context.Context, event domain.RideCompletedEvent) (domain.PaymentResult, error) {
	if event.EventID == "" || event.TripID == "" || event.CustomerID == "" || event.Amount <= 0 {
		return domain.PaymentResult{}, ErrInvalidRideCompletedEvent
	}

	existing, found, err := s.processed.Get(ctx, event.EventID)
	if err != nil {
		return domain.PaymentResult{}, err
	}
	if found {
		return existing, nil
	}

	if err := s.wallets.DebitCustomer(ctx, event.CustomerID, event.Amount, event.EventID); err != nil {
		_ = s.publisher.Publish(ctx, domain.Event{Name: "payment.failed.v1", TripID: event.TripID, OccurredAt: s.now()})
		return domain.PaymentResult{}, err
	}

	invoiceID, err := s.invoices.IssueRideInvoice(ctx, event.TripID, event.CustomerID, event.Amount, event.Currency)
	if err != nil {
		return domain.PaymentResult{}, err
	}

	result := domain.PaymentResult{
		EventID:   event.EventID,
		TripID:    event.TripID,
		InvoiceID: invoiceID,
		PaidAt:    s.now(),
	}
	if err := s.processed.Save(ctx, result); err != nil {
		return domain.PaymentResult{}, err
	}
	if err := s.publisher.Publish(ctx, domain.Event{Name: "payment.completed.v1", TripID: event.TripID, OccurredAt: result.PaidAt}); err != nil {
		return domain.PaymentResult{}, err
	}

	return result, nil
}
