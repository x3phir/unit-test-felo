package unit_test

import (
	"context"
	"testing"
	"time"

	"github.com/felo/felo-backend/services/location-service/internal/domain"
	"github.com/felo/felo-backend/services/location-service/internal/service"
	"go.uber.org/mock/gomock"
)

func TestLocationService_StoreDriverLocation_WritesHistoryAndCacheWithGomock(t *testing.T) {
	ctrl := gomock.NewController(t)
	history := NewMockHistoryStore(ctrl)
	cache := NewMockLatestCache(ctrl)
	router := NewMockRouter(ctrl)
	svc := service.NewLocationService(history, cache, router)

	sample := domain.LocationSample{
		DriverID:   "driver-1",
		Position:   domain.Coordinate{Latitude: -6.2, Longitude: 106.8},
		RecordedAt: time.Now(),
	}
	history.EXPECT().Append(gomock.Any(), sample).Return(nil)
	cache.EXPECT().SetLatest(gomock.Any(), sample).Return(nil)

	if err := svc.StoreDriverLocation(context.Background(), sample); err != nil {
		t.Fatalf("StoreDriverLocation() error = %v", err)
	}
}
