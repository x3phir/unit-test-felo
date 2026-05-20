//go:build functional

package functional_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/felo/felo-backend/services/location-service/internal/domain"
	"github.com/felo/felo-backend/services/location-service/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestLocationFunctional_StoreAndQueryHistory_UsesDatabase(t *testing.T) {
	ctx := context.Background()
	db := openLocationPG(t, getenv("FELO_LOCATION_PG_DSN", "postgres://felo:felo@127.0.0.1:54325/location_db?sslmode=disable"))
	t.Cleanup(func() { db.Close() })

	recordedAt := time.Now().UTC().Truncate(time.Second)
	_, _ = db.Exec(ctx, "delete from driver_locations where driver_id=$1 and recorded_at=$2", "driver-ft-001", recordedAt)

	svc := service.NewLocationService(&pgHistory{db: db}, noopLatestCache{}, noopRouter{})
	sample := domain.LocationSample{
		DriverID:   "driver-ft-001",
		Position:   domain.Coordinate{Latitude: -6.2, Longitude: 106.8},
		RecordedAt: recordedAt,
	}
	if err := svc.StoreDriverLocation(ctx, sample); err != nil {
		t.Fatalf("StoreDriverLocation() error = %v", err)
	}

	items, err := svc.GetLocationHistory(ctx, "driver-ft-001", recordedAt.Add(-time.Minute), recordedAt.Add(time.Minute))
	if err != nil {
		t.Fatalf("GetLocationHistory() error = %v", err)
	}
	if len(items) == 0 {
		t.Fatal("len(items) = 0, want at least 1")
	}
}

type pgHistory struct{ db *pgxpool.Pool }
func (h *pgHistory) Append(ctx context.Context, sample domain.LocationSample) error {
	_, err := h.db.Exec(ctx, "insert into driver_locations (driver_id, lat, lng, recorded_at) values ($1,$2,$3,$4)", sample.DriverID, sample.Position.Latitude, sample.Position.Longitude, sample.RecordedAt)
	return err
}
func (h *pgHistory) LatestByDriver(ctx context.Context, driverID string) (domain.LocationSample, bool, error) {
	var sample domain.LocationSample
	err := h.db.QueryRow(ctx, "select driver_id, lat, lng, recorded_at from driver_locations where driver_id=$1 order by recorded_at desc limit 1", driverID).
		Scan(&sample.DriverID, &sample.Position.Latitude, &sample.Position.Longitude, &sample.RecordedAt)
	if err != nil {
		return domain.LocationSample{}, false, nil
	}
	return sample, true, nil
}
func (h *pgHistory) ListByDriver(ctx context.Context, driverID string, from, to time.Time) ([]domain.LocationSample, error) {
	rows, err := h.db.Query(ctx, "select driver_id, lat, lng, recorded_at from driver_locations where driver_id=$1 and recorded_at between $2 and $3 order by recorded_at", driverID, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []domain.LocationSample
	for rows.Next() {
		var sample domain.LocationSample
		if err := rows.Scan(&sample.DriverID, &sample.Position.Latitude, &sample.Position.Longitude, &sample.RecordedAt); err != nil {
			return nil, err
		}
		items = append(items, sample)
	}
	return items, nil
}

type noopLatestCache struct{}
func (noopLatestCache) SetLatest(context.Context, domain.LocationSample) error { return nil }
func (noopLatestCache) GetLatest(context.Context, string) (domain.LocationSample, bool, error) { return domain.LocationSample{}, false, nil }

type noopRouter struct{}
func (noopRouter) Estimate(context.Context, domain.RouteRequest) (domain.RouteEstimate, error) { return domain.RouteEstimate{}, nil }

func openLocationPG(t *testing.T, dsn string) *pgxpool.Pool {
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
