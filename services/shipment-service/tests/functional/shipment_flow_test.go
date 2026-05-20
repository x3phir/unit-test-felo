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

func TestShipmentFunctional_PickupShipment_PersistsToDatabase(t *testing.T) {
	ctx := context.Background()
	db := openShipmentPG(t, getenv("FELO_SHIPMENT_PG_DSN", "postgres://felo:felo@127.0.0.1:54329/shipment_db?sslmode=disable"))
	t.Cleanup(func() { db.Close() })

	initShipmentTables(t, db)

	shipmentID := "ship-ft-001"
	_, _ = db.Exec(ctx, "delete from shipments where id=$1", shipmentID)
	_, _ = db.Exec(ctx, `insert into shipments (id, send_order_id, status, updated_at)
		values ($1,'sendorder-ft-001','created',$2)`, shipmentID, time.Now().UTC())

	svc := service.NewShipmentService(&pgShipmentRepo{db: db}, &noopShipmentPublisher{})

	shipment, err := svc.PickupShipment(ctx, shipmentID, "driver-ft-001")
	if err != nil {
		t.Fatalf("PickupShipment() error = %v", err)
	}

	var status string
	var driverID string
	if err := db.QueryRow(ctx, "select status, driver_id from shipments where id=$1", shipment.ID).Scan(&status, &driverID); err != nil {
		t.Fatalf("query persisted shipment: %v", err)
	}
	if status != string(domain.StatusPickedUp) {
		t.Fatalf("persisted status = %s, want %s", status, domain.StatusPickedUp)
	}
	if driverID != "driver-ft-001" {
		t.Fatalf("persisted driver_id = %s, want driver-ft-001", driverID)
	}
}

func TestShipmentFunctional_DeliverShipment_PersistsToDatabase(t *testing.T) {
	ctx := context.Background()
	db := openShipmentPG(t, getenv("FELO_SHIPMENT_PG_DSN", "postgres://felo:felo@127.0.0.1:54329/shipment_db?sslmode=disable"))
	t.Cleanup(func() { db.Close() })

	initShipmentTables(t, db)

	shipmentID := "ship-ft-002"
	_, _ = db.Exec(ctx, "delete from shipments where id=$1", shipmentID)
	_, _ = db.Exec(ctx, `insert into shipments (id, send_order_id, driver_id, status, updated_at)
		values ($1,'sendorder-ft-001','driver-ft-001','picked_up',$2)`, shipmentID, time.Now().UTC())

	svc := service.NewShipmentService(&pgShipmentRepo{db: db}, &noopShipmentPublisher{})

	shipment, err := svc.DeliverShipment(ctx, shipmentID, domain.ProofOfDelivery{
		PhotoURL: "http://example.com/photo.jpg",
	})
	if err != nil {
		t.Fatalf("DeliverShipment() error = %v", err)
	}

	var status string
	if err := db.QueryRow(ctx, "select status from shipments where id=$1", shipment.ID).Scan(&status); err != nil {
		t.Fatalf("query persisted shipment: %v", err)
	}
	if status != string(domain.StatusDelivered) {
		t.Fatalf("persisted status = %s, want %s", status, domain.StatusDelivered)
	}
}

type pgShipmentRepo struct{ db *pgxpool.Pool }

func (r *pgShipmentRepo) Save(ctx context.Context, shipment domain.Shipment) error {
	_, err := r.db.Exec(ctx, `insert into shipments (id, send_order_id, driver_id, status, tracking_number, proof_photo_url, proof_signature, updated_at)
values ($1,$2,$3,$4,$5,$6,$7,$8)
on conflict (id) do update set
send_order_id=excluded.send_order_id,
driver_id=excluded.driver_id,
status=excluded.status,
tracking_number=excluded.tracking_number,
proof_photo_url=excluded.proof_photo_url,
proof_signature=excluded.proof_signature,
updated_at=excluded.updated_at`,
		shipment.ID, shipment.SendOrderID, shipment.DriverID, string(shipment.Status),
		shipment.TrackingNumber, shipment.ProofOfDelivery.PhotoURL,
		shipment.ProofOfDelivery.Signature, shipment.UpdatedAt)
	return err
}

func (r *pgShipmentRepo) GetByID(ctx context.Context, shipmentID string) (domain.Shipment, error) {
	var shipment domain.Shipment
	var status string
	err := r.db.QueryRow(ctx,
		`select id, send_order_id, driver_id, status, tracking_number, proof_photo_url, proof_signature, updated_at
from shipments where id=$1`, shipmentID).
		Scan(&shipment.ID, &shipment.SendOrderID, &shipment.DriverID, &status,
			&shipment.TrackingNumber, &shipment.ProofOfDelivery.PhotoURL,
			&shipment.ProofOfDelivery.Signature, &shipment.UpdatedAt)
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
		id text primary key,
		send_order_id text not null,
		driver_id text not null default '',
		status text not null,
		tracking_number text not null default '',
		proof_photo_url text not null default '',
		proof_signature text not null default '',
		updated_at timestamptz not null
	)`)
	if err != nil {
		t.Fatalf("initShipmentTables: %v", err)
	}
}

func getenv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func openShipmentPG(t *testing.T, dsn string) *pgxpool.Pool {
	t.Helper()
	db, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		t.Fatalf("pgxpool.New() error = %v", err)
	}
	return db
}
