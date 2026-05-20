//go:build functional

package functional_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/felo/felo-backend/services/shipment-service/internal/domain"
	"github.com/felo/felo-backend/services/shipment-service/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestShipmentFunctional_PickupShipment_PersistsDriverAndStatus(t *testing.T) {
	ctx := context.Background()
	db := openShipmentPG(t, getenv("FELO_SHIPMENT_PG_DSN", "postgres://felo:felo@127.0.0.1:54336/shipment_db?sslmode=disable"))
	t.Cleanup(func() { db.Close() })

	initShipmentTables(t, db)

	shipmentID := "ship-ft-001"
	_, _ = db.Exec(ctx, "delete from shipments where shipment_id=$1", shipmentID)
	_, _ = db.Exec(ctx, `insert into shipments (shipment_id, status, eta_minutes, updated_at)
		values ($1,'packed',0,$2)`, shipmentID, time.Now().UTC())

	svc := service.NewShipmentService(&pgShipmentRepo{db: db}, &noopShipmentPublisher{})

	shipment, err := svc.PickupShipment(ctx, shipmentID, "driver-ft-001")
	if err != nil {
		t.Fatalf("PickupShipment() error = %v", err)
	}

	var status string
	if err := db.QueryRow(ctx, "select status from shipments where shipment_id=$1", shipment.ID).Scan(&status); err != nil {
		t.Fatalf("query persisted shipment: %v", err)
	}
	if status != string(domain.StatusPickedUp) {
		t.Fatalf("persisted status = %s, want %s", status, domain.StatusPickedUp)
	}
}

func TestShipmentFunctional_DeliverShipment_UpdatesStatusToDelivered(t *testing.T) {
	ctx := context.Background()
	db := openShipmentPG(t, getenv("FELO_SHIPMENT_PG_DSN", "postgres://felo:felo@127.0.0.1:54336/shipment_db?sslmode=disable"))
	t.Cleanup(func() { db.Close() })

	initShipmentTables(t, db)

	shipmentID := "ship-ft-002"
	_, _ = db.Exec(ctx, "delete from shipments where shipment_id=$1", shipmentID)
	_, _ = db.Exec(ctx, `insert into shipments (shipment_id, driver_id, status, eta_minutes, updated_at)
		values ($1,'driver-ft-001','packed',15,$2)`, shipmentID, time.Now().UTC())

	svc := service.NewShipmentService(&pgShipmentRepo{db: db}, &noopShipmentPublisher{})

	_, err := svc.DeliverShipment(ctx, shipmentID, domain.ProofOfDelivery{
		PhotoURL: "http://example.com/photo.jpg",
	})
	if err != nil {
		t.Fatalf("DeliverShipment() error = %v", err)
	}

	var status string
	if err := db.QueryRow(ctx, "select status from shipments where shipment_id=$1", shipmentID).Scan(&status); err != nil {
		t.Fatalf("query persisted shipment: %v", err)
	}
	if status != string(domain.StatusDelivered) {
		t.Fatalf("persisted status = %s, want %s", status, domain.StatusDelivered)
	}
}

type pgShipmentRepo struct{ db *pgxpool.Pool }

func (r *pgShipmentRepo) Save(ctx context.Context, shipment domain.Shipment) error {
	_, err := r.db.Exec(ctx, `insert into shipments (shipment_id, order_ref, driver_id, status, eta_minutes, updated_at)
values ($1,$2,$3,$4,$5,$6)
on conflict (shipment_id) do update set
order_ref=excluded.order_ref,
driver_id=excluded.driver_id,
status=excluded.status,
eta_minutes=excluded.eta_minutes,
updated_at=excluded.updated_at`,
		shipment.ID, shipment.SendOrderID, shipment.DriverID, string(shipment.Status),
		0, shipment.UpdatedAt)
	return err
}

func (r *pgShipmentRepo) GetByID(ctx context.Context, shipmentID string) (domain.Shipment, error) {
	var shipment domain.Shipment
	var status string
	err := r.db.QueryRow(ctx,
		`select shipment_id, order_ref, driver_id, status, eta_minutes, updated_at
from shipments where shipment_id=$1`, shipmentID).
		Scan(&shipment.ID, &shipment.SendOrderID, &shipment.DriverID, &status,
			new(int), &shipment.UpdatedAt)
	if err != nil {
		return domain.Shipment{}, err
	}
	shipment.Status = domain.ShipmentStatus(status)
	return shipment, nil
}

type noopShipmentPublisher struct{}

func (noopShipmentPublisher) Publish(_ context.Context, _ domain.Event) error { return nil }

func initShipmentTables(t *testing.T, db *pgxpool.Pool) {
	t.Helper()
	ctx := context.Background()
	_, err := db.Exec(ctx, `create table if not exists shipments (
		shipment_id text primary key,
		order_ref text not null default '',
		driver_id text not null default '',
		status text not null,
		eta_minutes integer not null default 0,
		updated_at timestamptz not null
	)`)
	if err != nil {
		t.Fatalf("initShipmentTables: %v", err)
	}
}

func openShipmentPG(t *testing.T, dsn string) *pgxpool.Pool {
	t.Helper()
	db, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		t.Fatalf("pgxpool.New() error = %v", err)
	}
	return db
}

func getenv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
