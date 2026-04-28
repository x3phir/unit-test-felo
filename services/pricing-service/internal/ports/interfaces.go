package ports

import (
	"context"
	"time"

	"github.com/felo/felo-backend/services/pricing-service/internal/domain"
)

// SurgeConfigReader reads operator-configurable surge parameters.
// PRD 4.5.2: surge multiplier and threshold values are configurable by ops.
type SurgeConfigReader interface {
	GetSurgeConfig(ctx context.Context) (domain.SurgeConfig, error)
}

// FareAuditLog stores and retrieves fare calculations for audit.
// PRD 4.5.4: each fare calculation must be logged with its full input set.
type FareAuditLog interface {
	Save(ctx context.Context, entry domain.FareAuditEntry) error
	GetByTripID(ctx context.Context, tripID string) (domain.FareAuditEntry, bool, error)
}

// Clock abstracts time for deterministic testing.
type Clock interface {
	Now() time.Time
}
