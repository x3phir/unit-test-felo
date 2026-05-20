package unit_test

import (
	"context"
	"reflect"

	"github.com/felo/felo-backend/services/merchant-service/internal/domain"
	"go.uber.org/mock/gomock"
)

type MockMerchantRepository struct {
	ctrl     *gomock.Controller
	recorder *MockMerchantRepositoryRecorder
}
type MockMerchantRepositoryRecorder struct{ mock *MockMerchantRepository }

func NewMockMerchantRepository(ctrl *gomock.Controller) *MockMerchantRepository {
	m := &MockMerchantRepository{ctrl: ctrl}
	m.recorder = &MockMerchantRepositoryRecorder{m}
	return m
}
func (m *MockMerchantRepository) EXPECT() *MockMerchantRepositoryRecorder { return m.recorder }
func (m *MockMerchantRepository) GetByID(ctx context.Context, id string) (*domain.Merchant, error) {
	ret := m.ctrl.Call(m, "GetByID", ctx, id)
	ret0, _ := ret[0].(*domain.Merchant)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}
func (mr *MockMerchantRepositoryRecorder) GetByID(ctx, id any) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockMerchantRepository)(nil).GetByID), ctx, id)
}
func (m *MockMerchantRepository) Create(ctx context.Context, merchant *domain.Merchant) error {
	ret := m.ctrl.Call(m, "Create", ctx, merchant)
	ret0, _ := ret[0].(error)
	return ret0
}
func (mr *MockMerchantRepositoryRecorder) Create(ctx, merchant any) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockMerchantRepository)(nil).Create), ctx, merchant)
}
func (m *MockMerchantRepository) UpdateStatus(ctx context.Context, id string, isClosed bool) error {
	ret := m.ctrl.Call(m, "UpdateStatus", ctx, id, isClosed)
	ret0, _ := ret[0].(error)
	return ret0
}
func (mr *MockMerchantRepositoryRecorder) UpdateStatus(ctx, id, isClosed any) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateStatus", reflect.TypeOf((*MockMerchantRepository)(nil).UpdateStatus), ctx, id, isClosed)
}

type MockMenuRepository struct {
	ctrl     *gomock.Controller
	recorder *MockMenuRepositoryRecorder
}
type MockMenuRepositoryRecorder struct{ mock *MockMenuRepository }

func NewMockMenuRepository(ctrl *gomock.Controller) *MockMenuRepository {
	m := &MockMenuRepository{ctrl: ctrl}
	m.recorder = &MockMenuRepositoryRecorder{m}
	return m
}
func (m *MockMenuRepository) EXPECT() *MockMenuRepositoryRecorder { return m.recorder }
func (m *MockMenuRepository) GetByMerchantID(ctx context.Context, merchantID string) ([]domain.Menu, error) {
	ret := m.ctrl.Call(m, "GetByMerchantID", ctx, merchantID)
	ret0, _ := ret[0].([]domain.Menu)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}
func (mr *MockMenuRepositoryRecorder) GetByMerchantID(ctx, merchantID any) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByMerchantID", reflect.TypeOf((*MockMenuRepository)(nil).GetByMerchantID), ctx, merchantID)
}
func (m *MockMenuRepository) GetByID(ctx context.Context, id string) (*domain.Menu, error) {
	ret := m.ctrl.Call(m, "GetByID", ctx, id)
	ret0, _ := ret[0].(*domain.Menu)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}
func (mr *MockMenuRepositoryRecorder) GetByID(ctx, id any) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockMenuRepository)(nil).GetByID), ctx, id)
}
func (m *MockMenuRepository) GetByIDs(ctx context.Context, ids []string) ([]domain.Menu, error) {
	ret := m.ctrl.Call(m, "GetByIDs", ctx, ids)
	ret0, _ := ret[0].([]domain.Menu)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}
func (mr *MockMenuRepositoryRecorder) GetByIDs(ctx, ids any) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByIDs", reflect.TypeOf((*MockMenuRepository)(nil).GetByIDs), ctx, ids)
}
func (m *MockMenuRepository) Create(ctx context.Context, menu *domain.Menu) error {
	ret := m.ctrl.Call(m, "Create", ctx, menu)
	ret0, _ := ret[0].(error)
	return ret0
}
func (mr *MockMenuRepositoryRecorder) Create(ctx, menu any) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockMenuRepository)(nil).Create), ctx, menu)
}
func (m *MockMenuRepository) Update(ctx context.Context, menu *domain.Menu) error {
	ret := m.ctrl.Call(m, "Update", ctx, menu)
	ret0, _ := ret[0].(error)
	return ret0
}
func (mr *MockMenuRepositoryRecorder) Update(ctx, menu any) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockMenuRepository)(nil).Update), ctx, menu)
}
func (m *MockMenuRepository) UpdateAvailability(ctx context.Context, id string, isAvailable bool) error {
	ret := m.ctrl.Call(m, "UpdateAvailability", ctx, id, isAvailable)
	ret0, _ := ret[0].(error)
	return ret0
}
func (mr *MockMenuRepositoryRecorder) UpdateAvailability(ctx, id, isAvailable any) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateAvailability", reflect.TypeOf((*MockMenuRepository)(nil).UpdateAvailability), ctx, id, isAvailable)
}
