package service

import (
	"context"
	"errors"
	"math"

	"github.com/felo/felo-backend/services/pricing-service/internal/domain"
	"github.com/felo/felo-backend/services/pricing-service/internal/ports"
)

var (
	ErrInvalidPricingInput  = errors.New("invalid pricing input")
	ErrEstimateNotFound     = errors.New("estimate not found")
	ErrFareExceedsTolerance = errors.New("final fare exceeds tolerance")
)

const (
	defaultCurrency  = "IDR"
	tolerancePercent = 0.10
)

// PricingService implements all fare calculation logic.
// PRD 4.5 GUARDRAIL: all fare calculations must go through pricing-service
// exclusively. No other service may compute fares independently.
type PricingService struct {
	config ports.SurgeConfigReader
	audit  ports.FareAuditLog
	clock  ports.Clock
}

func NewPricingService(config ports.SurgeConfigReader, audit ports.FareAuditLog, clock ports.Clock) *PricingService {
	return &PricingService{
		config: config,
		audit:  audit,
		clock:  clock,
	}
}

// CalculateEstimate computes the fare shown to the customer before order
// confirmation. PRD 4.5.3 point 1.
func (s *PricingService) CalculateEstimate(ctx context.Context, tripID string, input domain.PricingInput) (domain.FareEstimate, error) {
	if tripID == "" || input.DistanceKM <= 0 || input.DurationMins <= 0 {
		return domain.FareEstimate{}, ErrInvalidPricingInput
	}

	cfg, err := s.config.GetSurgeConfig(ctx)
	if err != nil {
		return domain.FareEstimate{}, err
	}

	baseFare := s.calculateBaseFare(input, cfg)
	multiplier := s.calculateSurgeMultiplier(input, cfg)
	finalFare := int64(math.Round(float64(baseFare) * multiplier))

	now := s.clock.Now()
	estimate := domain.FareEstimate{
		TripID:          tripID,
		BaseFare:        baseFare,
		SurgeMultiplier: multiplier,
		FinalFare:       finalFare,
		Currency:        defaultCurrency,
		CalculatedAt:    now,
	}

	entry := domain.FareAuditEntry{
		TripID:          tripID,
		Input:           input,
		BaseFare:        baseFare,
		SurgeMultiplier: multiplier,
		FinalFare:       finalFare,
		Currency:        defaultCurrency,
		CalculatedAt:    now,
	}
	if err := s.audit.Save(ctx, entry); err != nil {
		return domain.FareEstimate{}, err
	}

	return estimate, nil
}

// CalculateFinalFare recalculates the fare at ride completion using actual
// distance and duration. PRD 4.5.3 point 2: final adjustment is allowed
// within ±10% of the estimated fare.
func (s *PricingService) CalculateFinalFare(ctx context.Context, tripID string, actual domain.PricingInput) (domain.FareEstimate, error) {
	if tripID == "" || actual.DistanceKM <= 0 || actual.DurationMins <= 0 {
		return domain.FareEstimate{}, ErrInvalidPricingInput
	}

	original, found, err := s.audit.GetByTripID(ctx, tripID)
	if err != nil {
		return domain.FareEstimate{}, err
	}
	if !found {
		return domain.FareEstimate{}, ErrEstimateNotFound
	}

	cfg, err := s.config.GetSurgeConfig(ctx)
	if err != nil {
		return domain.FareEstimate{}, err
	}

	baseFare := s.calculateBaseFare(actual, cfg)
	multiplier := original.SurgeMultiplier
	finalFare := int64(math.Round(float64(baseFare) * multiplier))

	tolerance := math.Abs(float64(original.FinalFare) * tolerancePercent)
	diff := math.Abs(float64(finalFare) - float64(original.FinalFare))
	if diff > tolerance {
		return domain.FareEstimate{}, ErrFareExceedsTolerance
	}

	now := s.clock.Now()
	estimate := domain.FareEstimate{
		TripID:          tripID,
		BaseFare:        baseFare,
		SurgeMultiplier: multiplier,
		FinalFare:       finalFare,
		Currency:        defaultCurrency,
		CalculatedAt:    now,
	}

	entry := domain.FareAuditEntry{
		TripID:          tripID,
		Input:           actual,
		BaseFare:        baseFare,
		SurgeMultiplier: multiplier,
		FinalFare:       finalFare,
		Currency:        defaultCurrency,
		CalculatedAt:    now,
	}
	if err := s.audit.Save(ctx, entry); err != nil {
		return domain.FareEstimate{}, err
	}

	return estimate, nil
}

// calculateBaseFare computes the base fare from distance and duration.
// If the result is below the configured minimum, the minimum fare is used.
func (s *PricingService) calculateBaseFare(input domain.PricingInput, cfg domain.SurgeConfig) int64 {
	fare := int64(math.Round(input.DistanceKM*float64(cfg.BaseFarePerKM) + input.DurationMins*float64(cfg.BaseFarePerMinute)))
	if fare < cfg.MinimumFare {
		fare = cfg.MinimumFare
	}
	return fare
}

// calculateSurgeMultiplier derives the surge multiplier from demand/supply.
// PRD 4.5.2: multiplier is applied when demand/supply ratio exceeds threshold.
// When supply is zero, no rides can be matched so multiplier defaults to 1.0.
func (s *PricingService) calculateSurgeMultiplier(input domain.PricingInput, cfg domain.SurgeConfig) float64 {
	if input.SupplyLevel <= 0 {
		return 1.0
	}
	ratio := float64(input.DemandLevel) / float64(input.SupplyLevel)
	if ratio <= cfg.DemandSupplyThreshold {
		return 1.0
	}
	multiplier := ratio / cfg.DemandSupplyThreshold
	if multiplier > cfg.MaxMultiplier {
		multiplier = cfg.MaxMultiplier
	}
	return multiplier
}
