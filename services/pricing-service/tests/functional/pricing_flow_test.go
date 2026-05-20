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

func TestPricingFunctional_CalculateEstimate_ReturnsCorrectFare(t *testing.T) {
	ctx := context.Background()
	db := openPricingPG(t, getenv("FELO_PRICING_PG_DSN", "postgres://felo:felo@127.0.0.1:54337/pricing_db?sslmode=disable"))
	t.Cleanup(func() { db.Close() })

	audit := &memFareAudit{entries: map[string]domain.FareAuditEntry{}}
	svc := service.NewPricingService(
		&fixedPricingConfig{},
		audit,
		functionalPricingClock{now: time.Date(2026, 5, 1, 10, 0, 0, 0, time.UTC)},
	)

	estimate, err := svc.CalculateEstimate(ctx, "pricing-ft-001", domain.PricingInput{
		DistanceKM:   10,
		DurationMins: 15,
		DemandLevel:  5,
		SupplyLevel:  10,
	})
	if err != nil {
		t.Fatalf("CalculateEstimate() error = %v", err)
	}
	if estimate.FinalFare != 32500 {
		t.Fatalf("estimate.FinalFare = %d, want 32500", estimate.FinalFare)
	}
	if estimate.SurgeMultiplier != 1.0 {
		t.Fatalf("estimate.SurgeMultiplier = %v, want 1.0", estimate.SurgeMultiplier)
	}
}

func TestPricingFunctional_CalculateFinalFare_ReadsExistingAudit(t *testing.T) {
	ctx := context.Background()
	db := openPricingPG(t, getenv("FELO_PRICING_PG_DSN", "postgres://felo:felo@127.0.0.1:54337/pricing_db?sslmode=disable"))
	t.Cleanup(func() { db.Close() })

	audit := &memFareAudit{entries: map[string]domain.FareAuditEntry{}}
	audit.entries["pricing-ft-002"] = domain.FareAuditEntry{
		TripID:          "pricing-ft-002",
		Input:           domain.PricingInput{DistanceKM: 10, DurationMins: 15, DemandLevel: 5, SupplyLevel: 10},
		BaseFare:        32500,
		SurgeMultiplier: 1.0,
		FinalFare:       32500,
		Currency:        "IDR",
		CalculatedAt:    time.Date(2026, 5, 1, 9, 0, 0, 0, time.UTC),
	}

	svc := service.NewPricingService(
		&fixedPricingConfig{},
		audit,
		functionalPricingClock{now: time.Date(2026, 5, 1, 10, 0, 0, 0, time.UTC)},
	)

	estimate, err := svc.CalculateFinalFare(ctx, "pricing-ft-002", domain.PricingInput{
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
}

type memFareAudit struct {
	entries map[string]domain.FareAuditEntry
}

func (m *memFareAudit) Save(_ context.Context, entry domain.FareAuditEntry) error {
	m.entries[entry.TripID] = entry
	return nil
}

func (m *memFareAudit) GetByTripID(_ context.Context, tripID string) (domain.FareAuditEntry, bool, error) {
	entry, ok := m.entries[tripID]
	return entry, ok, nil
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
