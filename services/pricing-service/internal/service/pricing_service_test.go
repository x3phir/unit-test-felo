package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/felo/felo-backend/services/pricing-service/internal/domain"
	"github.com/felo/felo-backend/services/pricing-service/internal/service"
)

// --- CalculateEstimate tests ---

func TestPricingService_CalculateEstimate_CalculatesBaseFareCorrectly(t *testing.T) {
	now := time.Date(2026, 4, 28, 9, 0, 0, 0, time.UTC)
	cfg := defaultSurgeConfig()
	audit := &fareAuditLogFake{entries: map[string]domain.FareAuditEntry{}}
	svc := service.NewPricingService(&surgeConfigFake{config: cfg}, audit, fixedClock{now: now})

	estimate, err := svc.CalculateEstimate(context.Background(), "trip-1", domain.PricingInput{
		DistanceKM:   10,
		DurationMins: 15,
		DemandLevel:  5,
		SupplyLevel:  10,
	})
	if err != nil {
		t.Fatalf("CalculateEstimate() error = %v", err)
	}

	// base = 10 * 2500 + 15 * 500 = 32500, demand/supply = 0.5 <= 1.5 so no surge
	wantFare := int64(32500)
	if estimate.FinalFare != wantFare {
		t.Fatalf("estimate.FinalFare = %d, want %d", estimate.FinalFare, wantFare)
	}
	if estimate.SurgeMultiplier != 1.0 {
		t.Fatalf("estimate.SurgeMultiplier = %v, want 1.0", estimate.SurgeMultiplier)
	}
	if estimate.Currency != "IDR" {
		t.Fatalf("estimate.Currency = %s, want IDR", estimate.Currency)
	}
}

func TestPricingService_CalculateEstimate_AppliesSurgeMultiplier(t *testing.T) {
	now := time.Date(2026, 4, 28, 9, 0, 0, 0, time.UTC)
	cfg := defaultSurgeConfig()
	audit := &fareAuditLogFake{entries: map[string]domain.FareAuditEntry{}}
	svc := service.NewPricingService(&surgeConfigFake{config: cfg}, audit, fixedClock{now: now})

	estimate, err := svc.CalculateEstimate(context.Background(), "trip-1", domain.PricingInput{
		DistanceKM:   10,
		DurationMins: 15,
		DemandLevel:  30,
		SupplyLevel:  10,
	})
	if err != nil {
		t.Fatalf("CalculateEstimate() error = %v", err)
	}

	// demand/supply = 3.0, threshold = 1.5, so multiplier = 3.0 / 1.5 = 2.0
	// base = 32500, final = 32500 * 2.0 = 65000
	if estimate.SurgeMultiplier != 2.0 {
		t.Fatalf("estimate.SurgeMultiplier = %v, want 2.0", estimate.SurgeMultiplier)
	}
	if estimate.FinalFare != 65000 {
		t.Fatalf("estimate.FinalFare = %d, want 65000", estimate.FinalFare)
	}
}

func TestPricingService_CalculateEstimate_CapsMultiplierAtMax(t *testing.T) {
	now := time.Date(2026, 4, 28, 9, 0, 0, 0, time.UTC)
	cfg := defaultSurgeConfig()
	audit := &fareAuditLogFake{entries: map[string]domain.FareAuditEntry{}}
	svc := service.NewPricingService(&surgeConfigFake{config: cfg}, audit, fixedClock{now: now})

	estimate, err := svc.CalculateEstimate(context.Background(), "trip-1", domain.PricingInput{
		DistanceKM:   10,
		DurationMins: 15,
		DemandLevel:  100,
		SupplyLevel:  5,
	})
	if err != nil {
		t.Fatalf("CalculateEstimate() error = %v", err)
	}

	// demand/supply = 20, uncapped multiplier = 20 / 1.5 = 13.3..., capped at 3.0
	if estimate.SurgeMultiplier != 3.0 {
		t.Fatalf("estimate.SurgeMultiplier = %v, want 3.0 (max)", estimate.SurgeMultiplier)
	}
}

func TestPricingService_CalculateEstimate_NoSurgeWhenBelowThreshold(t *testing.T) {
	now := time.Date(2026, 4, 28, 9, 0, 0, 0, time.UTC)
	cfg := defaultSurgeConfig()
	audit := &fareAuditLogFake{entries: map[string]domain.FareAuditEntry{}}
	svc := service.NewPricingService(&surgeConfigFake{config: cfg}, audit, fixedClock{now: now})

	estimate, err := svc.CalculateEstimate(context.Background(), "trip-1", domain.PricingInput{
		DistanceKM:   5,
		DurationMins: 10,
		DemandLevel:  10,
		SupplyLevel:  10,
	})
	if err != nil {
		t.Fatalf("CalculateEstimate() error = %v", err)
	}

	// demand/supply = 1.0, threshold = 1.5, no surge
	if estimate.SurgeMultiplier != 1.0 {
		t.Fatalf("estimate.SurgeMultiplier = %v, want 1.0", estimate.SurgeMultiplier)
	}
}

func TestPricingService_CalculateEstimate_AppliesMinimumFare(t *testing.T) {
	now := time.Date(2026, 4, 28, 9, 0, 0, 0, time.UTC)
	cfg := defaultSurgeConfig()
	audit := &fareAuditLogFake{entries: map[string]domain.FareAuditEntry{}}
	svc := service.NewPricingService(&surgeConfigFake{config: cfg}, audit, fixedClock{now: now})

	estimate, err := svc.CalculateEstimate(context.Background(), "trip-1", domain.PricingInput{
		DistanceKM:   0.5,
		DurationMins: 1,
		DemandLevel:  5,
		SupplyLevel:  10,
	})
	if err != nil {
		t.Fatalf("CalculateEstimate() error = %v", err)
	}

	// base = 0.5 * 2500 + 1 * 500 = 1750, but minimum = 7000
	if estimate.BaseFare != 7000 {
		t.Fatalf("estimate.BaseFare = %d, want 7000 (minimum)", estimate.BaseFare)
	}
}

func TestPricingService_CalculateEstimate_ZeroSupply_DefaultsToNoSurge(t *testing.T) {
	now := time.Date(2026, 4, 28, 9, 0, 0, 0, time.UTC)
	cfg := defaultSurgeConfig()
	audit := &fareAuditLogFake{entries: map[string]domain.FareAuditEntry{}}
	svc := service.NewPricingService(&surgeConfigFake{config: cfg}, audit, fixedClock{now: now})

	estimate, err := svc.CalculateEstimate(context.Background(), "trip-1", domain.PricingInput{
		DistanceKM:   5,
		DurationMins: 10,
		DemandLevel:  20,
		SupplyLevel:  0,
	})
	if err != nil {
		t.Fatalf("CalculateEstimate() error = %v", err)
	}

	if estimate.SurgeMultiplier != 1.0 {
		t.Fatalf("estimate.SurgeMultiplier = %v, want 1.0", estimate.SurgeMultiplier)
	}
}

func TestPricingService_CalculateEstimate_InvalidInput_ReturnsError(t *testing.T) {
	cfg := defaultSurgeConfig()
	svc := service.NewPricingService(&surgeConfigFake{config: cfg}, &fareAuditLogFake{entries: map[string]domain.FareAuditEntry{}}, fixedClock{now: time.Now()})

	_, err := svc.CalculateEstimate(context.Background(), "", domain.PricingInput{DistanceKM: 5, DurationMins: 10})
	if !errors.Is(err, service.ErrInvalidPricingInput) {
		t.Fatalf("CalculateEstimate() error = %v, want ErrInvalidPricingInput", err)
	}

	_, err = svc.CalculateEstimate(context.Background(), "trip-1", domain.PricingInput{})
	if !errors.Is(err, service.ErrInvalidPricingInput) {
		t.Fatalf("CalculateEstimate() error = %v, want ErrInvalidPricingInput", err)
	}
}

func TestPricingService_CalculateEstimate_StoresAuditEntry(t *testing.T) {
	now := time.Date(2026, 4, 28, 9, 0, 0, 0, time.UTC)
	cfg := defaultSurgeConfig()
	audit := &fareAuditLogFake{entries: map[string]domain.FareAuditEntry{}}
	svc := service.NewPricingService(&surgeConfigFake{config: cfg}, audit, fixedClock{now: now})

	_, err := svc.CalculateEstimate(context.Background(), "trip-1", domain.PricingInput{
		DistanceKM:   10,
		DurationMins: 15,
		DemandLevel:  5,
		SupplyLevel:  10,
	})
	if err != nil {
		t.Fatalf("CalculateEstimate() error = %v", err)
	}

	entry, ok := audit.entries["trip-1"]
	if !ok {
		t.Fatal("audit entry not stored for trip-1")
	}
	if entry.Input.DistanceKM != 10 {
		t.Fatalf("entry.Input.DistanceKM = %v, want 10", entry.Input.DistanceKM)
	}
	if !entry.CalculatedAt.Equal(now) {
		t.Fatalf("entry.CalculatedAt = %v, want %v", entry.CalculatedAt, now)
	}
}

func TestPricingService_CalculateEstimate_DeterministicOutput(t *testing.T) {
	now := time.Date(2026, 4, 28, 9, 0, 0, 0, time.UTC)
	cfg := defaultSurgeConfig()
	input := domain.PricingInput{
		DistanceKM:   8,
		DurationMins: 12,
		DemandLevel:  20,
		SupplyLevel:  10,
	}

	svc1 := service.NewPricingService(&surgeConfigFake{config: cfg}, &fareAuditLogFake{entries: map[string]domain.FareAuditEntry{}}, fixedClock{now: now})
	svc2 := service.NewPricingService(&surgeConfigFake{config: cfg}, &fareAuditLogFake{entries: map[string]domain.FareAuditEntry{}}, fixedClock{now: now})

	est1, err := svc1.CalculateEstimate(context.Background(), "trip-1", input)
	if err != nil {
		t.Fatalf("CalculateEstimate() run 1 error = %v", err)
	}
	est2, err := svc2.CalculateEstimate(context.Background(), "trip-1", input)
	if err != nil {
		t.Fatalf("CalculateEstimate() run 2 error = %v", err)
	}

	if est1.FinalFare != est2.FinalFare {
		t.Fatalf("FinalFare not deterministic: %d != %d", est1.FinalFare, est2.FinalFare)
	}
	if est1.SurgeMultiplier != est2.SurgeMultiplier {
		t.Fatalf("SurgeMultiplier not deterministic: %v != %v", est1.SurgeMultiplier, est2.SurgeMultiplier)
	}
}

func TestPricingService_CalculateEstimate_ConfigFailure_ReturnsError(t *testing.T) {
	configErr := errors.New("config unavailable")
	svc := service.NewPricingService(&surgeConfigFake{err: configErr}, &fareAuditLogFake{entries: map[string]domain.FareAuditEntry{}}, fixedClock{now: time.Now()})

	_, err := svc.CalculateEstimate(context.Background(), "trip-1", domain.PricingInput{
		DistanceKM:   5,
		DurationMins: 10,
		DemandLevel:  5,
		SupplyLevel:  10,
	})
	if err == nil {
		t.Fatal("CalculateEstimate() error = nil, want error")
	}
}

// --- CalculateFinalFare tests ---

func TestPricingService_CalculateFinalFare_WithinTolerance_ReturnsRecalculatedFare(t *testing.T) {
	now := time.Date(2026, 4, 28, 10, 0, 0, 0, time.UTC)
	cfg := defaultSurgeConfig()
	// original estimate: base = 10 * 2500 + 15 * 500 = 32500, multiplier 1.0, final = 32500
	audit := &fareAuditLogFake{entries: map[string]domain.FareAuditEntry{
		"trip-1": {
			TripID:          "trip-1",
			Input:           domain.PricingInput{DistanceKM: 10, DurationMins: 15, DemandLevel: 5, SupplyLevel: 10},
			BaseFare:        32500,
			SurgeMultiplier: 1.0,
			FinalFare:       32500,
			Currency:        "IDR",
			CalculatedAt:    now.Add(-30 * time.Minute),
		},
	}}
	svc := service.NewPricingService(&surgeConfigFake{config: cfg}, audit, fixedClock{now: now})

	// actual ride: slightly longer distance but within 10%
	estimate, err := svc.CalculateFinalFare(context.Background(), "trip-1", domain.PricingInput{
		DistanceKM:   10.5,
		DurationMins: 16,
		DemandLevel:  5,
		SupplyLevel:  10,
	})
	if err != nil {
		t.Fatalf("CalculateFinalFare() error = %v", err)
	}

	// new base = 10.5 * 2500 + 16 * 500 = 26250 + 8000 = 34250
	// multiplier from original = 1.0, final = 34250
	// diff = |34250 - 32500| = 1750, tolerance = 32500 * 0.1 = 3250 → OK
	if estimate.FinalFare != 34250 {
		t.Fatalf("estimate.FinalFare = %d, want 34250", estimate.FinalFare)
	}
}

func TestPricingService_CalculateFinalFare_ExceedsTolerance_ReturnsError(t *testing.T) {
	now := time.Date(2026, 4, 28, 10, 0, 0, 0, time.UTC)
	cfg := defaultSurgeConfig()
	audit := &fareAuditLogFake{entries: map[string]domain.FareAuditEntry{
		"trip-1": {
			TripID:          "trip-1",
			BaseFare:        32500,
			SurgeMultiplier: 1.0,
			FinalFare:       32500,
			Currency:        "IDR",
			CalculatedAt:    now.Add(-30 * time.Minute),
		},
	}}
	svc := service.NewPricingService(&surgeConfigFake{config: cfg}, audit, fixedClock{now: now})

	// actual ride: much longer, exceeds 10% tolerance
	_, err := svc.CalculateFinalFare(context.Background(), "trip-1", domain.PricingInput{
		DistanceKM:   20,
		DurationMins: 30,
		DemandLevel:  5,
		SupplyLevel:  10,
	})
	if !errors.Is(err, service.ErrFareExceedsTolerance) {
		t.Fatalf("CalculateFinalFare() error = %v, want ErrFareExceedsTolerance", err)
	}
}

func TestPricingService_CalculateFinalFare_EstimateNotFound_ReturnsError(t *testing.T) {
	now := time.Date(2026, 4, 28, 10, 0, 0, 0, time.UTC)
	cfg := defaultSurgeConfig()
	audit := &fareAuditLogFake{entries: map[string]domain.FareAuditEntry{}}
	svc := service.NewPricingService(&surgeConfigFake{config: cfg}, audit, fixedClock{now: now})

	_, err := svc.CalculateFinalFare(context.Background(), "trip-999", domain.PricingInput{
		DistanceKM:   5,
		DurationMins: 10,
		DemandLevel:  5,
		SupplyLevel:  10,
	})
	if !errors.Is(err, service.ErrEstimateNotFound) {
		t.Fatalf("CalculateFinalFare() error = %v, want ErrEstimateNotFound", err)
	}
}

func TestPricingService_CalculateFinalFare_InvalidInput_ReturnsError(t *testing.T) {
	svc := service.NewPricingService(&surgeConfigFake{config: defaultSurgeConfig()}, &fareAuditLogFake{entries: map[string]domain.FareAuditEntry{}}, fixedClock{now: time.Now()})

	_, err := svc.CalculateFinalFare(context.Background(), "", domain.PricingInput{DistanceKM: 5, DurationMins: 10})
	if !errors.Is(err, service.ErrInvalidPricingInput) {
		t.Fatalf("CalculateFinalFare() error = %v, want ErrInvalidPricingInput", err)
	}
}

func TestPricingService_CalculateFinalFare_UsesSurgeMultiplierFromOriginalEstimate(t *testing.T) {
	now := time.Date(2026, 4, 28, 10, 0, 0, 0, time.UTC)
	cfg := defaultSurgeConfig()
	// original estimate had 2x surge
	audit := &fareAuditLogFake{entries: map[string]domain.FareAuditEntry{
		"trip-1": {
			TripID:          "trip-1",
			BaseFare:        32500,
			SurgeMultiplier: 2.0,
			FinalFare:       65000,
			Currency:        "IDR",
			CalculatedAt:    now.Add(-30 * time.Minute),
		},
	}}
	svc := service.NewPricingService(&surgeConfigFake{config: cfg}, audit, fixedClock{now: now})

	// actual ride: same distance/duration, current demand is low, but original surge must apply
	estimate, err := svc.CalculateFinalFare(context.Background(), "trip-1", domain.PricingInput{
		DistanceKM:   10,
		DurationMins: 15,
		DemandLevel:  5,
		SupplyLevel:  10,
	})
	if err != nil {
		t.Fatalf("CalculateFinalFare() error = %v", err)
	}

	// base = 32500, multiplier from original = 2.0, final = 65000
	if estimate.SurgeMultiplier != 2.0 {
		t.Fatalf("estimate.SurgeMultiplier = %v, want 2.0 (from original)", estimate.SurgeMultiplier)
	}
	if estimate.FinalFare != 65000 {
		t.Fatalf("estimate.FinalFare = %d, want 65000", estimate.FinalFare)
	}
}

// --- Fakes ---

func defaultSurgeConfig() domain.SurgeConfig {
	return domain.SurgeConfig{
		DemandSupplyThreshold: 1.5,
		MaxMultiplier:         3.0,
		BaseFarePerKM:         2500,
		BaseFarePerMinute:     500,
		MinimumFare:           7000,
	}
}

type surgeConfigFake struct {
	config domain.SurgeConfig
	err    error
}

func (f *surgeConfigFake) GetSurgeConfig(_ context.Context) (domain.SurgeConfig, error) {
	return f.config, f.err
}

type fareAuditLogFake struct {
	entries map[string]domain.FareAuditEntry
}

func (f *fareAuditLogFake) Save(_ context.Context, entry domain.FareAuditEntry) error {
	f.entries[entry.TripID] = entry
	return nil
}

func (f *fareAuditLogFake) GetByTripID(_ context.Context, tripID string) (domain.FareAuditEntry, bool, error) {
	entry, ok := f.entries[tripID]
	return entry, ok, nil
}

type fixedClock struct {
	now time.Time
}

func (c fixedClock) Now() time.Time {
	return c.now
}
