//go:build functional

package functional_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/felo/felo-backend/services/pricing-service/internal/domain"
	"github.com/felo/felo-backend/services/pricing-service/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestPricingFunctional_CalculateEstimate_PersistsFareAudit(t *testing.T) {
	ctx := context.Background()
	db := openPricingPG(t, getenv("FELO_PRICING_PG_DSN", "postgres://felo:felo@127.0.0.1:54337/pricing_db?sslmode=disable"))
	t.Cleanup(func() { db.Close() })

	initPricingTables(t, db)

	tripID := "pricing-ft-001"
	_, _ = db.Exec(ctx, "delete from fare_audit where trip_id=$1", tripID)

	svc := service.NewPricingService(
		&fixedPricingConfig{},
		&pgFareAudit{db: db},
		functionalPricingClock{now: time.Date(2026, 5, 1, 10, 0, 0, 0, time.UTC)},
	)

	_, err := svc.CalculateEstimate(ctx, tripID, domain.PricingInput{
		DistanceKM:   10,
		DurationMins: 15,
		DemandLevel:  5,
		SupplyLevel:  10,
	})
	if err != nil {
		t.Fatalf("CalculateEstimate() error = %v", err)
	}

	var baseFare int64
	var surgeMultiplier float64
	var finalFare int64
	err = db.QueryRow(ctx,
		"select base_fare, surge_multiplier, final_fare from fare_audit where trip_id=$1", tripID).
		Scan(&baseFare, &surgeMultiplier, &finalFare)
	if err != nil {
		t.Fatalf("query persisted fare_audit: %v", err)
	}
	if baseFare != 32500 {
		t.Fatalf("persisted base_fare = %d, want 32500", baseFare)
	}
	if surgeMultiplier != 1.0 {
		t.Fatalf("persisted surge_multiplier = %v, want 1.0", surgeMultiplier)
	}
	if finalFare != 32500 {
		t.Fatalf("persisted final_fare = %d, want 32500", finalFare)
	}
}

func TestPricingFunctional_CalculateFinalFare_ReadsExistingAuditAndPersistsUpdated(t *testing.T) {
	ctx := context.Background()
	db := openPricingPG(t, getenv("FELO_PRICING_PG_DSN", "postgres://felo:felo@127.0.0.1:54337/pricing_db?sslmode=disable"))
	t.Cleanup(func() { db.Close() })

	initPricingTables(t, db)

	tripID := "pricing-ft-002"
	_, _ = db.Exec(ctx, "delete from fare_audit where trip_id=$1", tripID)
	_, _ = db.Exec(ctx, `insert into fare_audit (trip_id, distance_km, duration_mins, demand_level, supply_level, base_fare, surge_multiplier, final_fare, currency, calculated_at)
		values ($1,10,15,5,10,32500,1.0,32500,'IDR',$2)`, tripID, time.Date(2026, 5, 1, 9, 0, 0, 0, time.UTC))

	svc := service.NewPricingService(
		&fixedPricingConfig{},
		&pgFareAudit{db: db},
		functionalPricingClock{now: time.Date(2026, 5, 1, 10, 0, 0, 0, time.UTC)},
	)

	estimate, err := svc.CalculateFinalFare(ctx, tripID, domain.PricingInput{
		DistanceKM:   10.5,
		DurationMins: 16,
		DemandLevel:  5,
		SupplyLevel:  10,
	})
	if err != nil {
		t.Fatalf("CalculateFinalFare() error = %v", err)
	}
	if estimate.FinalFare != 34250 {
		t.Fatalf("estimate.FinalFare = %d, want 34250", estimate.FinalFare)
	}

	var count int
	if err := db.QueryRow(ctx, "select count(*) from fare_audit where trip_id=$1", tripID).Scan(&count); err != nil {
		t.Fatalf("query fare_audit count: %v", err)
	}
	if count != 2 {
		t.Fatalf("expected 2 fare_audit entries (original + updated), got %d", count)
	}
}

type pgFareAudit struct{ db *pgxpool.Pool }

func (r *pgFareAudit) Save(ctx context.Context, entry domain.FareAuditEntry) error {
	_, err := r.db.Exec(ctx, `insert into fare_audit (trip_id, distance_km, duration_mins, demand_level, supply_level, base_fare, surge_multiplier, final_fare, currency, calculated_at)
values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`,
		entry.TripID, entry.Input.DistanceKM, entry.Input.DurationMins, entry.Input.DemandLevel, entry.Input.SupplyLevel,
		entry.BaseFare, entry.SurgeMultiplier, entry.FinalFare, entry.Currency, entry.CalculatedAt)
	return err
}

func (r *pgFareAudit) GetByTripID(ctx context.Context, tripID string) (domain.FareAuditEntry, bool, error) {
	var entry domain.FareAuditEntry
	var input domain.PricingInput
	err := r.db.QueryRow(ctx,
		`select trip_id, distance_km, duration_mins, demand_level, supply_level, base_fare, surge_multiplier, final_fare, currency, calculated_at
from fare_audit where trip_id=$1 order by calculated_at desc limit 1`, tripID).
		Scan(&entry.TripID, &input.DistanceKM, &input.DurationMins, &input.DemandLevel, &input.SupplyLevel,
			&entry.BaseFare, &entry.SurgeMultiplier, &entry.FinalFare, &entry.Currency, &entry.CalculatedAt)
	if err != nil {
		return domain.FareAuditEntry{}, false, nil
	}
	entry.Input = input
	return entry, true, nil
}

type fixedPricingConfig struct{}

func (fixedPricingConfig) GetSurgeConfig(_ context.Context) (domain.SurgeConfig, error) {
	return domain.SurgeConfig{
		DemandSupplyThreshold: 1.5,
		MaxMultiplier:         3.0,
		BaseFarePerKM:         2500,
		BaseFarePerMinute:     500,
		MinimumFare:           7000,
	}, nil
}

type functionalPricingClock struct{ now time.Time }

func (c functionalPricingClock) Now() time.Time { return c.now }

func initPricingTables(t *testing.T, db *pgxpool.Pool) {
	t.Helper()
	ctx := context.Background()
	_, err := db.Exec(ctx, `create table if not exists fare_audit (
		trip_id text not null,
		distance_km numeric not null,
		duration_mins numeric not null,
		demand_level int not null,
		supply_level int not null,
		base_fare bigint not null,
		surge_multiplier numeric not null,
		final_fare bigint not null,
		currency text not null,
		calculated_at timestamptz not null
	)`)
	if err != nil {
		t.Fatalf("initPricingTables fare_audit: %v", err)
	}
}

func openPricingPG(t *testing.T, dsn string) *pgxpool.Pool {
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
