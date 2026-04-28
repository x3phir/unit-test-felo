package service

import (
	"context"
	"errors"
	"time"

	"github.com/felo/felo-backend/services/send-order-service/internal/domain"
	"github.com/felo/felo-backend/services/send-order-service/internal/ports"
)

var (
	ErrInvalidInput = errors.New("invalid input")
)

type SendOrderService struct {
	repo      ports.SendOrderRepository
	pricing   ports.PricingClient
	invoice   ports.InvoiceClient
	publisher ports.EventPublisher
	ids       ports.IDGenerator
	now       func() time.Time
}

func NewSendOrderService(repo ports.SendOrderRepository, pricing ports.PricingClient, invoice ports.InvoiceClient, publisher ports.EventPublisher, ids ports.IDGenerator) *SendOrderService {
	return &SendOrderService{
		repo:      repo,
		pricing:   pricing,
		invoice:   invoice,
		publisher: publisher,
		ids:       ids,
		now:       time.Now,
	}
}

type CreateSendOrderInput struct {
	SenderID       string
	ReceiverPhone  string
	Origin         string
	Destination    string
	PackageDetails domain.PackageDetails
	PayerType      domain.PayerType
}

func (s *SendOrderService) CreateSendOrder(ctx context.Context, input CreateSendOrderInput) (domain.SendOrder, error) {
	if input.SenderID == "" || input.ReceiverPhone == "" || input.PackageDetails.WeightKG <= 0 {
		return domain.SendOrder{}, ErrInvalidInput
	}
	if input.PayerType != domain.PayerSender && input.PayerType != domain.PayerReceiver {
		return domain.SendOrder{}, ErrInvalidInput
	}

	fee, err := s.pricing.CalculateShippingFee(ctx, input.PackageDetails, input.Origin, input.Destination)
	if err != nil {
		return domain.SendOrder{}, err
	}

	now := s.now()
	order := domain.SendOrder{
		ID:             s.ids.NewID(),
		SenderID:       input.SenderID,
		ReceiverPhone:  input.ReceiverPhone,
		PackageDetails: input.PackageDetails,
		PayerType:      input.PayerType,
		ShippingFee:    fee,
		Status:         "created",
		CreatedAt:      now,
	}

	if err := s.repo.Save(ctx, order); err != nil {
		return domain.SendOrder{}, err
	}

	payerID := input.SenderID
	if input.PayerType == domain.PayerReceiver {
		payerID = input.ReceiverPhone
	}

	if err := s.invoice.CreateInvoice(ctx, order.ID, payerID, input.PayerType, fee); err != nil {
		return domain.SendOrder{}, err
	}

	if err := s.publisher.Publish(ctx, domain.Event{Name: "shipment.created.v1", OrderID: order.ID, OccurredAt: now}); err != nil {
		return domain.SendOrder{}, err
	}

	return order, nil
}
