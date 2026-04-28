package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/felo/felo-backend/services/location-service/internal/domain"
	"github.com/felo/felo-backend/services/location-service/internal/service"
)

func TestLocationService_StoreDriverLocation_WritesHistoryAndCache(t *testing.T) {
	history := &historyStoreFake{}
	cache := &latestCacheFake{}
	svc := service.NewLocationService(history, cache, &routerFake{})

	sample := domain.LocationSample{
		DriverID:   "driver-1",
		Position:   domain.Coordinate{Latitude: -6.2, Longitude: 106.8},
		RecordedAt: time.Date(2026, 4, 28, 13, 0, 0, 0, time.UTC),
	}

	if err := svc.StoreDriverLocation(context.Background(), sample); err != nil {
		t.Fatalf("StoreDriverLocation() error = %v", err)
	}
	if len(history.samples) != 1 {
		t.Fatalf("len(history.samples) = %d, want 1", len(history.samples))
	}
	if cache.latest.DriverID != "driver-1" {
		t.Fatalf("cache.latest.DriverID = %s, want driver-1", cache.latest.DriverID)
	}
}

func TestLocationService_StoreDriverLocation_InvalidSampleReturnsError(t *testing.T) {
	svc := service.NewLocationService(&historyStoreFake{}, &latestCacheFake{}, &routerFake{})

	err := svc.StoreDriverLocation(context.Background(), domain.LocationSample{})
	if err == nil {
		t.Fatal("StoreDriverLocation() error = nil, want error")
	}
}

func TestLocationService_GetLatestDriverLocation_ReturnsCachedSample(t *testing.T) {
	cache := &latestCacheFake{
		latest: domain.LocationSample{
			DriverID:   "driver-1",
			Position:   domain.Coordinate{Latitude: -6.2, Longitude: 106.8},
			RecordedAt: time.Date(2026, 4, 28, 13, 0, 0, 0, time.UTC),
		},
	}
	svc := service.NewLocationService(&historyStoreFake{}, cache, &routerFake{})

	got, err := svc.GetLatestDriverLocation(context.Background(), "driver-1")
	if err != nil {
		t.Fatalf("GetLatestDriverLocation() error = %v", err)
	}
	if got.DriverID != "driver-1" {
		t.Fatalf("got.DriverID = %s, want driver-1", got.DriverID)
	}
}

func TestLocationService_GetLatestDriverLocation_FallsBackToHistory(t *testing.T) {
	sample := domain.LocationSample{
		DriverID:   "driver-1",
		Position:   domain.Coordinate{Latitude: -6.2, Longitude: 106.8},
		RecordedAt: time.Date(2026, 4, 28, 13, 0, 0, 0, time.UTC),
	}
	history := &historyStoreFake{latest: sample}
	svc := service.NewLocationService(history, &latestCacheFake{}, &routerFake{})

	got, err := svc.GetLatestDriverLocation(context.Background(), "driver-1")
	if err != nil {
		t.Fatalf("GetLatestDriverLocation() error = %v", err)
	}
	if got.DriverID != "driver-1" {
		t.Fatalf("got.DriverID = %s, want driver-1", got.DriverID)
	}
}

func TestLocationService_GetLatestDriverLocation_NotFoundReturnsError(t *testing.T) {
	svc := service.NewLocationService(&historyStoreFake{}, &latestCacheFake{}, &routerFake{})

	_, err := svc.GetLatestDriverLocation(context.Background(), "driver-1")
	if err == nil {
		t.Fatal("GetLatestDriverLocation() error = nil, want error")
	}
}

func TestLocationService_GetLocationHistory_InvalidTimeWindowReturnsError(t *testing.T) {
	svc := service.NewLocationService(&historyStoreFake{}, &latestCacheFake{}, &routerFake{})

	_, err := svc.GetLocationHistory(context.Background(), "driver-1", time.Now(), time.Now())
	if err == nil {
		t.Fatal("GetLocationHistory() error = nil, want error")
	}
}

func TestLocationService_GetLocationHistory_ValidWindowReturnsSamples(t *testing.T) {
	sample := domain.LocationSample{
		DriverID:   "driver-1",
		Position:   domain.Coordinate{Latitude: -6.2, Longitude: 106.8},
		RecordedAt: time.Date(2026, 4, 28, 13, 0, 0, 0, time.UTC),
	}
	history := &historyStoreFake{samples: []domain.LocationSample{sample}}
	svc := service.NewLocationService(history, &latestCacheFake{}, &routerFake{})

	samples, err := svc.GetLocationHistory(context.Background(), "driver-1", sample.RecordedAt.Add(-time.Minute), sample.RecordedAt.Add(time.Minute))
	if err != nil {
		t.Fatalf("GetLocationHistory() error = %v", err)
	}
	if len(samples) != 1 {
		t.Fatalf("len(samples) = %d, want 1", len(samples))
	}
}

func TestLocationService_EstimateRoute_StampsCalculationTimeWhenProviderOmitsIt(t *testing.T) {
	router := &routerFake{estimate: domain.RouteEstimate{DistanceMeters: 1200, DurationSeconds: 360, Polyline: "abc"}}
	svc := service.NewLocationService(&historyStoreFake{}, &latestCacheFake{}, router)

	estimate, err := svc.EstimateRoute(context.Background(), domain.RouteRequest{
		Origin:      domain.Coordinate{Latitude: -6.2, Longitude: 106.8},
		Destination: domain.Coordinate{Latitude: -6.3, Longitude: 106.9},
	})
	if err != nil {
		t.Fatalf("EstimateRoute() error = %v", err)
	}
	if estimate.CalculatedAt.IsZero() {
		t.Fatal("estimate.CalculatedAt = zero, want stamped value")
	}
}

type historyStoreFake struct {
	samples []domain.LocationSample
	latest  domain.LocationSample
}

func (f *historyStoreFake) Append(_ context.Context, sample domain.LocationSample) error {
	f.samples = append(f.samples, sample)
	f.latest = sample
	return nil
}

func (f *historyStoreFake) LatestByDriver(_ context.Context, _ string) (domain.LocationSample, bool, error) {
	if f.latest.DriverID == "" {
		return domain.LocationSample{}, false, nil
	}
	return f.latest, true, nil
}

func (f *historyStoreFake) ListByDriver(_ context.Context, _ string, _ time.Time, _ time.Time) ([]domain.LocationSample, error) {
	return f.samples, nil
}

type latestCacheFake struct {
	latest domain.LocationSample
}

func (f *latestCacheFake) SetLatest(_ context.Context, sample domain.LocationSample) error {
	f.latest = sample
	return nil
}

func (f *latestCacheFake) GetLatest(_ context.Context, _ string) (domain.LocationSample, bool, error) {
	if f.latest.DriverID == "" {
		return domain.LocationSample{}, false, nil
	}
	return f.latest, true, nil
}

type routerFake struct {
	estimate domain.RouteEstimate
}

func (f *routerFake) Estimate(_ context.Context, _ domain.RouteRequest) (domain.RouteEstimate, error) {
	return f.estimate, nil
}
