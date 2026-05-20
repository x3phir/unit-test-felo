package unit_test

import (
	"context"
	"reflect"

	"github.com/felo/felo-backend/services/driver-service/internal/domain"
	"go.uber.org/mock/gomock"
)

type MockDriverRepository struct{ ctrl *gomock.Controller; recorder *MockDriverRepositoryRecorder }
type MockDriverRepositoryRecorder struct{ mock *MockDriverRepository }
func NewMockDriverRepository(ctrl *gomock.Controller) *MockDriverRepository { m := &MockDriverRepository{ctrl: ctrl}; m.recorder = &MockDriverRepositoryRecorder{m}; return m }
func (m *MockDriverRepository) EXPECT() *MockDriverRepositoryRecorder { return m.recorder }
func (m *MockDriverRepository) Save(ctx context.Context, driver domain.DriverProfile) error { ret := m.ctrl.Call(m, "Save", ctx, driver); ret0, _ := ret[0].(error); return ret0 }
func (mr *MockDriverRepositoryRecorder) Save(ctx, driver any) *gomock.Call { return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockDriverRepository)(nil).Save), ctx, driver) }
func (m *MockDriverRepository) GetByID(ctx context.Context, driverID string) (domain.DriverProfile, error) { ret := m.ctrl.Call(m, "GetByID", ctx, driverID); ret0, _ := ret[0].(domain.DriverProfile); ret1, _ := ret[1].(error); return ret0, ret1 }
func (mr *MockDriverRepositoryRecorder) GetByID(ctx, driverID any) *gomock.Call { return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockDriverRepository)(nil).GetByID), ctx, driverID) }
