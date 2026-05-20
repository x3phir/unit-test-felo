//go:build functional

package functional_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/felo/felo-backend/services/ride-service/internal/domain"
	"github.com/felo/felo-backend/services/ride-service/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestRideFunctional_RequestRide_PersistsTripToDatabase(t *testing.T) {
	ctx := context.Background()
	db := openPG(t, getenv("FELO_RIDE_PG_DSN", "postgres://felo:felo@127.0.0.1:54321/ride_db?sslmode=disable"))
	t.Cleanup(func() { db.Close() })

	initRideTables(t, db)

	rideID := "ride-ft-001"
	_, _ = db.Exec(ctx, "delete from rides where id=$1", rideID)

	repo := &pgRideRepo{db: db}
	publisher := &noopRidePublisher{}
	svc := service.NewTripService(repo, publisher, functionalClock{now: time.Now().UTC()}, &fixedRideIDs{ids: []string{rideID}})

	trip, err := svc.RequestRide(ctx, service.RequestRideInput{
		CustomerID:   "cust-active-001",
		Pickup:       domain.Coordinate{Latitude: -6.2, Longitude: 106.8},
		Destination:  domain.Coordinate{Latitude: -6.3, Longitude: 106.9},
		FareEstimate: 25000,
	})
	if err != nil {
		t.Fatalf("RequestRide() error = %v", err)
	}

	var state string
	if err := db.QueryRow(ctx, "select state from rides where id=$1", trip.ID).Scan(&state); err != nil {
		t.Fatalf("query persisted ride: %v", err)
	}
	if state != string(domain.StateMatching) {
		t.Fatalf("persisted state = %s, want %s", state, domain.StateMatching)
	}
}

type pgRideRepo struct{ db *pgxpool.Pool }

func (r *pgRideRepo) Save(ctx context.Context, trip domain.Trip) error {
	_, err := r.db.Exec(ctx, `insert into rides (id, customer_id, driver_id, pickup_lat, pickup_lng, dest_lat, dest_lng, fare, state, qr_code, qr_expires_at, qr_locked_driver, created_at, updated_at)
values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)
on conflict (id) do update set
customer_id=excluded.customer_id,
driver_id=excluded.driver_id,
pickup_lat=excluded.pickup_lat,
pickup_lng=excluded.pickup_lng,
dest_lat=excluded.dest_lat,
dest_lng=excluded.dest_lng,
fare=excluded.fare,
state=excluded.state,
qr_code=excluded.qr_code,
qr_expires_at=excluded.qr_expires_at,
qr_locked_driver=excluded.qr_locked_driver,
updated_at=excluded.updated_at`,
		trip.ID, trip.CustomerID, trip.DriverID, trip.Pickup.Latitude, trip.Pickup.Longitude, trip.Destination.Latitude, trip.Destination.Longitude,
		trip.FareEstimate, string(trip.State), trip.QRCode, nullableTime(trip.QRExpiresAt), "", trip.CreatedAt, trip.UpdatedAt)
	return err
}

func (r *pgRideRepo) GetByID(ctx context.Context, tripID string) (domain.Trip, error) {
	var trip domain.Trip
	var state string
	err := r.db.QueryRow(ctx, `select id, customer_id, driver_id, pickup_lat, pickup_lng, dest_lat, dest_lng, fare, state, qr_code, created_at, updated_at from rides where id=$1`, tripID).
		Scan(&trip.ID, &trip.CustomerID, &trip.DriverID, &trip.Pickup.Latitude, &trip.Pickup.Longitude, &trip.Destination.Latitude, &trip.Destination.Longitude, &trip.FareEstimate, &state, &trip.QRCode, &trip.CreatedAt, &trip.UpdatedAt)
	if err != nil {
		return domain.Trip{}, err
	}
	trip.State = domain.TripState(state)
	return trip, nil
}

type noopRidePublisher struct{}
func (noopRidePublisher) Publish(context.Context, domain.Event) error { return nil }

type functionalClock struct{ now time.Time }
func (c functionalClock) Now() time.Time { return c.now }

type fixedRideIDs struct{ ids []string; idx int }
func (g *fixedRideIDs) NewID() string { id := g.ids[g.idx]; g.idx++; return id }

func initRideTables(t *testing.T, db *pgxpool.Pool) {
	t.Helper()
	ctx := context.Background()
	_, err := db.Exec(ctx, `create table if not exists rides (
		id text primary key,
		customer_id text not null,
		driver_id text not null default '',
		pickup_lat double precision not null,
		pickup_lng double precision not null,
		dest_lat double precision not null,
		dest_lng double precision not null,
		fare bigint not null,
		state text not null,
		qr_code text not null default '',
		qr_expires_at timestamptz,
		qr_locked_driver text not null default '',
		created_at timestamptz not null,
		updated_at timestamptz not null
	)`)
	if err != nil {
		t.Fatalf("initRideTables: %v", err)
	}
}

func openPG(t *testing.T, dsn string) *pgxpool.Pool {
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

func nullableTime(ts time.Time) any {
	if ts.IsZero() {
		return nil
	}
	return ts
}
