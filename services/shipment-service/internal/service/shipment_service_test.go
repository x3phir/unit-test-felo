package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/felo/felo-backend/services/shipment-service/internal/domain"
	"github.com/felo/felo-backend/services/shipment-service/internal/service"
)

func TestShipmentService_PickupShipment_ValidTransition_Success(t *testing.T) {
	repo := &shipmentRepoFake{shipments: map[string]domain.Shipment{
		"ship-1": {ID: "ship-1", Status: domain.StatusCreated},
	}}
	publisher := &eventPublisherFake{}
	svc := service.NewShipmentService(repo, publisher)

	shipment, err := svc.PickupShipment(context.Background(), "ship-1", "driver-1")
	if err != nil {
		t.Fatalf("PickupShipment() error = %v", err)
	}

	if shipment.Status != domain.StatusPickedUp {
		t.Fatalf("Status = %s, want picked_up", shipment.Status)
	}
	if shipment.DriverID != "driver-1" {
		t.Fatalf("DriverID = %s, want driver-1", shipment.DriverID)
	}
	if len(publisher.events) != 1 || publisher.events[0].Name != "shipment.picked_up.v1" {
		t.Fatalf("publisher.events = %#v, want shipment.picked_up.v1", publisher.events)
	}
}

func TestShipmentService_PickupShipment_InvalidTransition_ReturnsError(t *testing.T) {
	repo := &shipmentRepoFake{shipments: map[string]domain.Shipment{
		"ship-1": {ID: "ship-1", Status: domain.StatusDelivered},
	}}
	svc := service.NewShipmentService(repo, &eventPublisherFake{})

	_, err := svc.PickupShipment(context.Background(), "ship-1", "driver-1")
	if !errors.Is(err, service.ErrInvalidStateTransition) {
		t.Fatalf("PickupShipment() error = %v, want ErrInvalidStateTransition", err)
	}
}

func TestShipmentService_DeliverShipment_ValidProof_Success(t *testing.T) {
	repo := &shipmentRepoFake{shipments: map[string]domain.Shipment{
		"ship-1": {ID: "ship-1", Status: domain.StatusPickedUp},
	}}
	publisher := &eventPublisherFake{}
	svc := service.NewShipmentService(repo, publisher)

	shipment, err := svc.DeliverShipment(context.Background(), "ship-1", domain.ProofOfDelivery{PhotoURL: "http://photo"})
	if err != nil {
		t.Fatalf("DeliverShipment() error = %v", err)
	}

	if shipment.Status != domain.StatusDelivered {
		t.Fatalf("Status = %s, want delivered", shipment.Status)
	}
	if len(publisher.events) != 1 || publisher.events[0].Name != "shipment.delivered.v1" {
		t.Fatalf("publisher.events = %#v, want shipment.delivered.v1", publisher.events)
	}
}

func TestShipmentService_DeliverShipment_InvalidProof_ReturnsError(t *testing.T) {
	svc := service.NewShipmentService(&shipmentRepoFake{}, &eventPublisherFake{})
	_, err := svc.DeliverShipment(context.Background(), "ship-1", domain.ProofOfDelivery{})
	if !errors.Is(err, service.ErrInvalidProof) {
		t.Fatalf("DeliverShipment() error = %v, want ErrInvalidProof", err)
	}
}

func TestShipmentService_DeliverShipment_InvalidTransition_ReturnsError(t *testing.T) {
	repo := &shipmentRepoFake{shipments: map[string]domain.Shipment{
		"ship-1": {ID: "ship-1", Status: domain.StatusCreated},
	}}
	svc := service.NewShipmentService(repo, &eventPublisherFake{})

	_, err := svc.DeliverShipment(context.Background(), "ship-1", domain.ProofOfDelivery{PhotoURL: "url"})
	if !errors.Is(err, service.ErrInvalidStateTransition) {
		t.Fatalf("DeliverShipment() error = %v, want ErrInvalidStateTransition", err)
	}
}

type shipmentRepoFake struct {
	shipments map[string]domain.Shipment
}

func (f *shipmentRepoFake) Save(_ context.Context, shipment domain.Shipment) error {
	if f.shipments == nil {
		f.shipments = make(map[string]domain.Shipment)
	}
	f.shipments[shipment.ID] = shipment
	return nil
}

func (f *shipmentRepoFake) GetByID(_ context.Context, shipmentID string) (domain.Shipment, error) {
	shipment, ok := f.shipments[shipmentID]
	if !ok {
		return domain.Shipment{}, errors.New("not found")
	}
	return shipment, nil
}

type eventPublisherFake struct {
	events []domain.Event
}

func (f *eventPublisherFake) Publish(_ context.Context, event domain.Event) error {
	f.events = append(f.events, event)
	return nil
}
