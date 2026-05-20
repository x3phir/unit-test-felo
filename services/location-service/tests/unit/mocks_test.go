package unit_test

import (
	"context"
	"reflect"
	"time"

	"github.com/felo/felo-backend/services/location-service/internal/domain"
	"go.uber.org/mock/gomock"
)

type MockHistoryStore struct{ ctrl *gomock.Controller; recorder *MockHistoryStoreRecorder }
type MockHistoryStoreRecorder struct{ mock *MockHistoryStore }
func NewMockHistoryStore(ctrl *gomock.Controller) *MockHistoryStore { m := &MockHistoryStore{ctrl: ctrl}; m.recorder = &MockHistoryStoreRecorder{m}; return m }
func (m *MockHistoryStore) EXPECT() *MockHistoryStoreRecorder { return m.recorder }
func (m *MockHistoryStore) Append(ctx context.Context, sample domain.LocationSample) error { ret := m.ctrl.Call(m, "Append", ctx, sample); ret0, _ := ret[0].(error); return ret0 }
func (mr *MockHistoryStoreRecorder) Append(ctx, sample any) *gomock.Call { return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Append", reflect.TypeOf((*MockHistoryStore)(nil).Append), ctx, sample) }
func (m *MockHistoryStore) LatestByDriver(ctx context.Context, driverID string) (domain.LocationSample, bool, error) { ret := m.ctrl.Call(m, "LatestByDriver", ctx, driverID); ret0, _ := ret[0].(domain.LocationSample); ret1, _ := ret[1].(bool); ret2, _ := ret[2].(error); return ret0, ret1, ret2 }
func (mr *MockHistoryStoreRecorder) LatestByDriver(ctx, driverID any) *gomock.Call { return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LatestByDriver", reflect.TypeOf((*MockHistoryStore)(nil).LatestByDriver), ctx, driverID) }
func (m *MockHistoryStore) ListByDriver(ctx context.Context, driverID string, from, to time.Time) ([]domain.LocationSample, error) { ret := m.ctrl.Call(m, "ListByDriver", ctx, driverID, from, to); ret0, _ := ret[0].([]domain.LocationSample); ret1, _ := ret[1].(error); return ret0, ret1 }
func (mr *MockHistoryStoreRecorder) ListByDriver(ctx, driverID, from, to any) *gomock.Call { return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListByDriver", reflect.TypeOf((*MockHistoryStore)(nil).ListByDriver), ctx, driverID, from, to) }

type MockLatestCache struct{ ctrl *gomock.Controller; recorder *MockLatestCacheRecorder }
type MockLatestCacheRecorder struct{ mock *MockLatestCache }
func NewMockLatestCache(ctrl *gomock.Controller) *MockLatestCache { m := &MockLatestCache{ctrl: ctrl}; m.recorder = &MockLatestCacheRecorder{m}; return m }
func (m *MockLatestCache) EXPECT() *MockLatestCacheRecorder { return m.recorder }
func (m *MockLatestCache) SetLatest(ctx context.Context, sample domain.LocationSample) error { ret := m.ctrl.Call(m, "SetLatest", ctx, sample); ret0, _ := ret[0].(error); return ret0 }
func (mr *MockLatestCacheRecorder) SetLatest(ctx, sample any) *gomock.Call { return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetLatest", reflect.TypeOf((*MockLatestCache)(nil).SetLatest), ctx, sample) }
func (m *MockLatestCache) GetLatest(ctx context.Context, driverID string) (domain.LocationSample, bool, error) { ret := m.ctrl.Call(m, "GetLatest", ctx, driverID); ret0, _ := ret[0].(domain.LocationSample); ret1, _ := ret[1].(bool); ret2, _ := ret[2].(error); return ret0, ret1, ret2 }
func (mr *MockLatestCacheRecorder) GetLatest(ctx, driverID any) *gomock.Call { return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLatest", reflect.TypeOf((*MockLatestCache)(nil).GetLatest), ctx, driverID) }

type MockRouter struct{ ctrl *gomock.Controller; recorder *MockRouterRecorder }
type MockRouterRecorder struct{ mock *MockRouter }
func NewMockRouter(ctrl *gomock.Controller) *MockRouter { m := &MockRouter{ctrl: ctrl}; m.recorder = &MockRouterRecorder{m}; return m }
func (m *MockRouter) EXPECT() *MockRouterRecorder { return m.recorder }
func (m *MockRouter) Estimate(ctx context.Context, request domain.RouteRequest) (domain.RouteEstimate, error) { ret := m.ctrl.Call(m, "Estimate", ctx, request); ret0, _ := ret[0].(domain.RouteEstimate); ret1, _ := ret[1].(error); return ret0, ret1 }
func (mr *MockRouterRecorder) Estimate(ctx, request any) *gomock.Call { return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Estimate", reflect.TypeOf((*MockRouter)(nil).Estimate), ctx, request) }
