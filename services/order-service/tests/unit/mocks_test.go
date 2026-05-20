package unit_test

import (
	"context"
	"reflect"

	"github.com/felo/felo-backend/services/order-service/internal/domain"
	"go.uber.org/mock/gomock"
)

type MockOrderRepository struct{ ctrl *gomock.Controller; recorder *MockOrderRepositoryRecorder }
type MockOrderRepositoryRecorder struct{ mock *MockOrderRepository }
func NewMockOrderRepository(ctrl *gomock.Controller) *MockOrderRepository { m := &MockOrderRepository{ctrl: ctrl}; m.recorder = &MockOrderRepositoryRecorder{m}; return m }
func (m *MockOrderRepository) EXPECT() *MockOrderRepositoryRecorder { return m.recorder }
func (m *MockOrderRepository) Save(ctx context.Context, order domain.FoodOrder) error { ret := m.ctrl.Call(m, "Save", ctx, order); ret0, _ := ret[0].(error); return ret0 }
func (mr *MockOrderRepositoryRecorder) Save(ctx, order any) *gomock.Call { return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockOrderRepository)(nil).Save), ctx, order) }
func (m *MockOrderRepository) GetByID(ctx context.Context, orderID string) (domain.FoodOrder, error) { ret := m.ctrl.Call(m, "GetByID", ctx, orderID); ret0, _ := ret[0].(domain.FoodOrder); ret1, _ := ret[1].(error); return ret0, ret1 }
func (mr *MockOrderRepositoryRecorder) GetByID(ctx, orderID any) *gomock.Call { return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockOrderRepository)(nil).GetByID), ctx, orderID) }

type MockOrderLocationClient struct{ ctrl *gomock.Controller; recorder *MockOrderLocationClientRecorder }
type MockOrderLocationClientRecorder struct{ mock *MockOrderLocationClient }
func NewMockOrderLocationClient(ctrl *gomock.Controller) *MockOrderLocationClient { m := &MockOrderLocationClient{ctrl: ctrl}; m.recorder = &MockOrderLocationClientRecorder{m}; return m }
func (m *MockOrderLocationClient) EXPECT() *MockOrderLocationClientRecorder { return m.recorder }
func (m *MockOrderLocationClient) GetDistanceKM(ctx context.Context, origin string, destination string) (float64, error) { ret := m.ctrl.Call(m, "GetDistanceKM", ctx, origin, destination); ret0, _ := ret[0].(float64); ret1, _ := ret[1].(error); return ret0, ret1 }
func (mr *MockOrderLocationClientRecorder) GetDistanceKM(ctx, origin, destination any) *gomock.Call { return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDistanceKM", reflect.TypeOf((*MockOrderLocationClient)(nil).GetDistanceKM), ctx, origin, destination) }

type MockAuthClient struct{ ctrl *gomock.Controller; recorder *MockAuthClientRecorder }
type MockAuthClientRecorder struct{ mock *MockAuthClient }
func NewMockAuthClient(ctrl *gomock.Controller) *MockAuthClient { m := &MockAuthClient{ctrl: ctrl}; m.recorder = &MockAuthClientRecorder{m}; return m }
func (m *MockAuthClient) EXPECT() *MockAuthClientRecorder { return m.recorder }
func (m *MockAuthClient) SendOTP(ctx context.Context, userID string) error { ret := m.ctrl.Call(m, "SendOTP", ctx, userID); ret0, _ := ret[0].(error); return ret0 }
func (mr *MockAuthClientRecorder) SendOTP(ctx, userID any) *gomock.Call { return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendOTP", reflect.TypeOf((*MockAuthClient)(nil).SendOTP), ctx, userID) }
func (m *MockAuthClient) VerifyOTP(ctx context.Context, userID string, otp string) (bool, error) { ret := m.ctrl.Call(m, "VerifyOTP", ctx, userID, otp); ret0, _ := ret[0].(bool); ret1, _ := ret[1].(error); return ret0, ret1 }
func (mr *MockAuthClientRecorder) VerifyOTP(ctx, userID, otp any) *gomock.Call { return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VerifyOTP", reflect.TypeOf((*MockAuthClient)(nil).VerifyOTP), ctx, userID, otp) }

type MockOrderEventPublisher struct{ ctrl *gomock.Controller; recorder *MockOrderEventPublisherRecorder }
type MockOrderEventPublisherRecorder struct{ mock *MockOrderEventPublisher }
func NewMockOrderEventPublisher(ctrl *gomock.Controller) *MockOrderEventPublisher { m := &MockOrderEventPublisher{ctrl: ctrl}; m.recorder = &MockOrderEventPublisherRecorder{m}; return m }
func (m *MockOrderEventPublisher) EXPECT() *MockOrderEventPublisherRecorder { return m.recorder }
func (m *MockOrderEventPublisher) Publish(ctx context.Context, event domain.Event) error { ret := m.ctrl.Call(m, "Publish", ctx, event); ret0, _ := ret[0].(error); return ret0 }
func (mr *MockOrderEventPublisherRecorder) Publish(ctx, event any) *gomock.Call { return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Publish", reflect.TypeOf((*MockOrderEventPublisher)(nil).Publish), ctx, event) }

type MockOrderIDGenerator struct{ ctrl *gomock.Controller; recorder *MockOrderIDGeneratorRecorder }
type MockOrderIDGeneratorRecorder struct{ mock *MockOrderIDGenerator }
func NewMockOrderIDGenerator(ctrl *gomock.Controller) *MockOrderIDGenerator { m := &MockOrderIDGenerator{ctrl: ctrl}; m.recorder = &MockOrderIDGeneratorRecorder{m}; return m }
func (m *MockOrderIDGenerator) EXPECT() *MockOrderIDGeneratorRecorder { return m.recorder }
func (m *MockOrderIDGenerator) NewID() string { ret := m.ctrl.Call(m, "NewID"); ret0, _ := ret[0].(string); return ret0 }
func (mr *MockOrderIDGeneratorRecorder) NewID() *gomock.Call { return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewID", reflect.TypeOf((*MockOrderIDGenerator)(nil).NewID)) }
