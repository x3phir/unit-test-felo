package unit_test

import (
	"context"
	"reflect"
	"time"

	"github.com/felo/felo-backend/services/pricing-service/internal/domain"
	"go.uber.org/mock/gomock"
)

type MockSurgeConfigReader struct{ ctrl *gomock.Controller; recorder *MockSurgeConfigReaderRecorder }
type MockSurgeConfigReaderRecorder struct{ mock *MockSurgeConfigReader }
func NewMockSurgeConfigReader(ctrl *gomock.Controller) *MockSurgeConfigReader { m := &MockSurgeConfigReader{ctrl: ctrl}; m.recorder = &MockSurgeConfigReaderRecorder{m}; return m }
func (m *MockSurgeConfigReader) EXPECT() *MockSurgeConfigReaderRecorder { return m.recorder }
func (m *MockSurgeConfigReader) GetSurgeConfig(ctx context.Context) (domain.SurgeConfig, error) { ret := m.ctrl.Call(m, "GetSurgeConfig", ctx); ret0, _ := ret[0].(domain.SurgeConfig); ret1, _ := ret[1].(error); return ret0, ret1 }
func (mr *MockSurgeConfigReaderRecorder) GetSurgeConfig(ctx any) *gomock.Call { return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSurgeConfig", reflect.TypeOf((*MockSurgeConfigReader)(nil).GetSurgeConfig), ctx) }

type MockFareAuditLog struct{ ctrl *gomock.Controller; recorder *MockFareAuditLogRecorder }
type MockFareAuditLogRecorder struct{ mock *MockFareAuditLog }
func NewMockFareAuditLog(ctrl *gomock.Controller) *MockFareAuditLog { m := &MockFareAuditLog{ctrl: ctrl}; m.recorder = &MockFareAuditLogRecorder{m}; return m }
func (m *MockFareAuditLog) EXPECT() *MockFareAuditLogRecorder { return m.recorder }
func (m *MockFareAuditLog) Save(ctx context.Context, entry domain.FareAuditEntry) error { ret := m.ctrl.Call(m, "Save", ctx, entry); ret0, _ := ret[0].(error); return ret0 }
func (mr *MockFareAuditLogRecorder) Save(ctx, entry any) *gomock.Call { return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockFareAuditLog)(nil).Save), ctx, entry) }
func (m *MockFareAuditLog) GetByTripID(ctx context.Context, tripID string) (domain.FareAuditEntry, bool, error) { ret := m.ctrl.Call(m, "GetByTripID", ctx, tripID); ret0, _ := ret[0].(domain.FareAuditEntry); ret1, _ := ret[1].(bool); ret2, _ := ret[2].(error); return ret0, ret1, ret2 }
func (mr *MockFareAuditLogRecorder) GetByTripID(ctx, tripID any) *gomock.Call { return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByTripID", reflect.TypeOf((*MockFareAuditLog)(nil).GetByTripID), ctx, tripID) }

type MockPricingClock struct{ ctrl *gomock.Controller; recorder *MockPricingClockRecorder }
type MockPricingClockRecorder struct{ mock *MockPricingClock }
func NewMockPricingClock(ctrl *gomock.Controller) *MockPricingClock { m := &MockPricingClock{ctrl: ctrl}; m.recorder = &MockPricingClockRecorder{m}; return m }
func (m *MockPricingClock) EXPECT() *MockPricingClockRecorder { return m.recorder }
func (m *MockPricingClock) Now() time.Time { ret := m.ctrl.Call(m, "Now"); ret0, _ := ret[0].(time.Time); return ret0 }
func (mr *MockPricingClockRecorder) Now() *gomock.Call { return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Now", reflect.TypeOf((*MockPricingClock)(nil).Now)) }
