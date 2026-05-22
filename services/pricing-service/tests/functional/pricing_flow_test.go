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
	initPricingTables(t, db)
	seedPricingRule(t, db)

	audit := &memFareAudit{entries: map[string]domain.FareAuditEntry{}}
	svc := service.NewPricingService(
		&pgPricingConfig{db: db},
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
	initPricingTables(t, db)
	seedPricingRule(t, db)

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
		&pgPricingConfig{db: db},
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

type pgPricingConfig struct{ db *pgxpool.Pool }

func (c *pgPricingConfig) GetSurgeConfig(ctx context.Context) (domain.SurgeConfig, error) {
	var baseFarePerKM int64
	var maxMultiplier float64
	err := c.db.QueryRow(ctx, `select base_fare, surge_multiplier
from pricing_rules
where service_type='ride'
order by active_from desc
limit 1`).Scan(&baseFarePerKM, &maxMultiplier)
	if err != nil {
		return domain.SurgeConfig{}, err
	}
	return domain.SurgeConfig{
		DemandSupplyThreshold: 1.5,
		MaxMultiplier:         maxMultiplier,
		BaseFarePerKM:         baseFarePerKM,
		BaseFarePerMinute:     500,
		MinimumFare:           7000,
	}, nil
}

type functionalPricingClock struct{ now time.Time }

func (c functionalPricingClock) Now() time.Time { return c.now }

func initPricingTables(t *testing.T, db *pgxpool.Pool) {
	t.Helper()
	ctx := context.Background()
	_, err := db.Exec(ctx, `create table if not exists pricing_rules (
		rule_id text primary key,
		service_type text not null,
		base_fare bigint not null,
		surge_multiplier numeric not null,
		active_from timestamptz not null,
		active_to timestamptz not null
	)`)
	if err != nil {
		t.Fatalf("initPricingTables: %v", err)
	}
}

func seedPricingRule(t *testing.T, db *pgxpool.Pool) {
	t.Helper()
	ctx := context.Background()
	_, err := db.Exec(ctx, `insert into pricing_rules (rule_id, service_type, base_fare, surge_multiplier, active_from, active_to)
values ('pricing-rule-ft-ride','ride',2500,3.0,$1,$2)
on conflict (rule_id) do update set
service_type=excluded.service_type,
base_fare=excluded.base_fare,
surge_multiplier=excluded.surge_multiplier,
active_from=excluded.active_from,
active_to=excluded.active_to`, time.Now().UTC().Add(-time.Hour), time.Now().UTC().Add(time.Hour))
	if err != nil {
		t.Fatalf("seedPricingRule: %v", err)
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
