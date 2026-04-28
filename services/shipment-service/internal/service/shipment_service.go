package service

import (
	"context"
	"errors"
	"time"

	"github.com/felo/felo-backend/services/shipment-service/internal/domain"
	"github.com/felo/felo-backend/services/shipment-service/internal/ports"
)

var (
	ErrInvalidStateTransition = errors.New("invalid state transition")
	ErrInvalidProof           = errors.New("invalid proof of delivery")
)

type ShipmentService struct {
	repo      ports.ShipmentRepository
	publisher ports.EventPublisher
	now       func() time.Time
}

func NewShipmentService(repo ports.ShipmentRepository, publisher ports.EventPublisher) *ShipmentService {
	return &ShipmentService{
		repo:      repo,
		publisher: publisher,
		now:       time.Now,
	}
}

func (s *ShipmentService) PickupShipment(ctx context.Context, shipmentID string, driverID string) (domain.Shipment, error) {
	shipment, err := s.repo.GetByID(ctx, shipmentID)
	if err != nil {
		return domain.Shipment{}, err
	}

	if shipment.Status != domain.StatusCreated && shipment.Status != "" {
		return domain.Shipment{}, ErrInvalidStateTransition
	}

	now := s.now()
	shipment.DriverID = driverID
	shipment.Status = domain.StatusPickedUp
	shipment.UpdatedAt = now

	if err := s.repo.Save(ctx, shipment); err != nil {
		return domain.Shipment{}, err
	}

	if err := s.publisher.Publish(ctx, domain.Event{Name: "shipment.picked_up.v1", ShipmentID: shipment.ID, OccurredAt: now}); err != nil {
		return domain.Shipment{}, err
	}

	return shipment, nil
}

func (s *ShipmentService) DeliverShipment(ctx context.Context, shipmentID string, proof domain.ProofOfDelivery) (domain.Shipment, error) {
	if proof.PhotoURL == "" && proof.Signature == "" {
		return domain.Shipment{}, ErrInvalidProof
	}

	shipment, err := s.repo.GetByID(ctx, shipmentID)
	if err != nil {
		return domain.Shipment{}, err
	}

	if shipment.Status != domain.StatusPickedUp && shipment.Status != domain.StatusInTransit {
		return domain.Shipment{}, ErrInvalidStateTransition
	}

	now := s.now()
	shipment.Status = domain.StatusDelivered
	shipment.ProofOfDelivery = proof
	shipment.UpdatedAt = now

	if err := s.repo.Save(ctx, shipment); err != nil {
		return domain.Shipment{}, err
	}

	if err := s.publisher.Publish(ctx, domain.Event{Name: "shipment.delivered.v1", ShipmentID: shipment.ID, OccurredAt: now}); err != nil {
		return domain.Shipment{}, err
	}

	return shipment, nil
}
