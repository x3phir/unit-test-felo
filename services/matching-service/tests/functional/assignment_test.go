//go:build functional

package functional_test

import (
	"context"
	"os"
	"testing"

	"github.com/felo/felo-backend/services/matching-service/internal/domain"
	"github.com/felo/felo-backend/services/matching-service/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestMatchingFunctional_AssignDriver_PersistsAssignmentToDatabase(t *testing.T) {
	ctx := context.Background()
	db := openMatchingPG(t, getenv("FELO_MATCHING_PG_DSN", "postgres://felo:felo@127.0.0.1:54322/matching_db?sslmode=disable"))
	t.Cleanup(func() { db.Close() })

	initMatchingTables(t, db)

	rideID := "ride-ft-assign-001"
	_, _ = db.Exec(ctx, "delete from assignments where ride_id=$1", rideID)
	_, _ = db.Exec(ctx, "delete from drivers where id=$1", "driver-ft-assign-001")
	_, _ = db.Exec(ctx, "insert into drivers (id, status, lat, lng) values ('driver-ft-assign-001','available',-6.200,106.816)")

	svc := service.NewMatchingService(&pgAvailability{db: db}, &pgAssignments{db: db}, noopMatchPublisher{})
	assignment, err := svc.AssignDriver(ctx, domain.MatchRequest{
		RideID:          rideID,
		Pickup:          domain.Coordinate{Latitude: -6.200, Longitude: 106.816},
		InitialRadiusKM: 1,
	})
	if err != nil {
		t.Fatalf("AssignDriver() error = %v", err)
	}

	var driverID string
	if err := db.QueryRow(ctx, "select driver_id from assignments where ride_id=$1", rideID).Scan(&driverID); err != nil {
		t.Fatalf("query assignment: %v", err)
	}
	if driverID != assignment.DriverID {
		t.Fatalf("persisted driver_id = %s, want %s", driverID, assignment.DriverID)
	}
}

type pgAvailability struct{ db *pgxpool.Pool }
func (a *pgAvailability) FindAvailableDrivers(ctx context.Context, _ domain.Coordinate, _ float64) ([]domain.DriverCandidate, error) {
	rows, err := a.db.Query(ctx, "select id, lat, lng from drivers where status='available' order by id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var drivers []domain.DriverCandidate
	for rows.Next() {
		var id string
		var lat, lng float64
		if err := rows.Scan(&id, &lat, &lng); err != nil {
			return nil, err
		}
		drivers = append(drivers, domain.DriverCandidate{DriverID: id, DistanceKM: lat + lng})
	}
	return drivers, nil
}

type pgAssignments struct{ db *pgxpool.Pool }
func (a *pgAssignments) Save(ctx context.Context, assignment domain.Assignment) error {
	_, err := a.db.Exec(ctx, `insert into assignments (ride_id, driver_id, assigned_at) values ($1,$2,$3)
on conflict (ride_id) do update set driver_id=excluded.driver_id, assigned_at=excluded.assigned_at`, assignment.RideID, assignment.DriverID, assignment.AssignedAt)
	return err
}

type noopMatchPublisher struct{}
func (noopMatchPublisher) Publish(context.Context, domain.Event) error { return nil }

func initMatchingTables(t *testing.T, db *pgxpool.Pool) {
	t.Helper()
	ctx := context.Background()
	_, err := db.Exec(ctx, `create table if not exists drivers (
		id text primary key,
		status text not null,
		lat double precision not null,
		lng double precision not null
	)`)
	if err != nil {
		t.Fatalf("initMatchingTables drivers: %v", err)
	}
	_, err = db.Exec(ctx, `create table if not exists assignments (
		ride_id text primary key,
		driver_id text not null,
		assigned_at timestamptz not null
	)`)
	if err != nil {
		t.Fatalf("initMatchingTables assignments: %v", err)
	}
}

func openMatchingPG(t *testing.T, dsn string) *pgxpool.Pool {
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
