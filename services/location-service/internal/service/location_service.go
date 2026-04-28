package service

import (
	"context"
	"errors"
	"time"

	"github.com/felo/felo-backend/services/location-service/internal/domain"
	"github.com/felo/felo-backend/services/location-service/internal/ports"
)

var (
	ErrInvalidLocationSample = errors.New("invalid location sample")
	ErrLocationNotFound      = errors.New("location not found")
	ErrInvalidTimeWindow     = errors.New("invalid time window")
)

type LocationService struct {
	history ports.HistoryStore
	cache   ports.LatestCache
	router  ports.Router
	now     func() time.Time
}

func NewLocationService(history ports.HistoryStore, cache ports.LatestCache, router ports.Router) *LocationService {
	return &LocationService{
		history: history,
		cache:   cache,
		router:  router,
		now:     time.Now,
	}
}

func (s *LocationService) StoreDriverLocation(ctx context.Context, sample domain.LocationSample) error {
	if sample.DriverID == "" || sample.RecordedAt.IsZero() {
		return ErrInvalidLocationSample
	}
	if err := s.history.Append(ctx, sample); err != nil {
		return err
	}
	if err := s.cache.SetLatest(ctx, sample); err != nil {
		return err
	}
	return nil
}

func (s *LocationService) GetLatestDriverLocation(ctx context.Context, driverID string) (domain.LocationSample, error) {
	sample, found, err := s.cache.GetLatest(ctx, driverID)
	if err != nil {
		return domain.LocationSample{}, err
	}
	if found {
		return sample, nil
	}

	sample, found, err = s.history.LatestByDriver(ctx, driverID)
	if err != nil {
		return domain.LocationSample{}, err
	}
	if !found {
		return domain.LocationSample{}, ErrLocationNotFound
	}
	return sample, nil
}

func (s *LocationService) GetLocationHistory(ctx context.Context, driverID string, from time.Time, to time.Time) ([]domain.LocationSample, error) {
	if driverID == "" || !from.Before(to) {
		return nil, ErrInvalidTimeWindow
	}
	return s.history.ListByDriver(ctx, driverID, from, to)
}

func (s *LocationService) EstimateRoute(ctx context.Context, request domain.RouteRequest) (domain.RouteEstimate, error) {
	estimate, err := s.router.Estimate(ctx, request)
	if err != nil {
		return domain.RouteEstimate{}, err
	}
	if estimate.CalculatedAt.IsZero() {
		estimate.CalculatedAt = s.now()
	}
	return estimate, nil
}
