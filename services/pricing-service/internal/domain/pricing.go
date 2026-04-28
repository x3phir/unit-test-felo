package domain

import "time"

// PricingInput holds all factors required for fare calculation.
// PRD 4.5.1: distance, duration, demand level, supply level.
type PricingInput struct {
	DistanceKM   float64
	DurationMins float64
	DemandLevel  int
	SupplyLevel  int
}

// SurgeConfig holds operator-configurable surge parameters.
// PRD 4.5.2: multiplier and threshold values are configurable by ops.
type SurgeConfig struct {
	DemandSupplyThreshold float64
	MaxMultiplier         float64
	BaseFarePerKM         int64
	BaseFarePerMinute     int64
	MinimumFare           int64
}

// FareEstimate is the result of a fare calculation.
type FareEstimate struct {
	TripID          string
	BaseFare        int64
	SurgeMultiplier float64
	FinalFare       int64
	Currency        string
	CalculatedAt    time.Time
}

// FareAuditEntry records the full input set for a fare calculation.
// PRD 4.5.4: each calculation must be logged with its full input set
// for audit and dispute resolution.
type FareAuditEntry struct {
	TripID          string
	Input           PricingInput
	BaseFare        int64
	SurgeMultiplier float64
	FinalFare       int64
	Currency        string
	CalculatedAt    time.Time
}
