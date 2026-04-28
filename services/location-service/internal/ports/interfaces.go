package ports

import (
	"context"
	"time"

	"github.com/felo/felo-backend/services/location-service/internal/domain"
)

type HistoryStore interface {
	Append(ctx context.Context, sample domain.LocationSample) error
	LatestByDriver(ctx context.Context, driverID string) (domain.LocationSample, bool, error)
	ListByDriver(ctx context.Context, driverID string, from time.Time, to time.Time) ([]domain.LocationSample, error)
}

type LatestCache interface {
	SetLatest(ctx context.Context, sample domain.LocationSample) error
	GetLatest(ctx context.Context, driverID string) (domain.LocationSample, bool, error)
}

type Router interface {
	Estimate(ctx context.Context, request domain.RouteRequest) (domain.RouteEstimate, error)
}
