package unit_test

import (
	"context"
	"reflect"

	"github.com/felo/felo-backend/services/shipment-service/internal/domain"
	"go.uber.org/mock/gomock"
)

type MockShipmentRepository struct{ ctrl *gomock.Controller; recorder *MockShipmentRepositoryRecorder }
type MockShipmentRepositoryRecorder struct{ mock *MockShipmentRepository }
func NewMockShipmentRepository(ctrl *gomock.Controller) *MockShipmentRepository { m := &MockShipmentRepository{ctrl: ctrl}; m.recorder = &MockShipmentRepositoryRecorder{m}; return m }
func (m *MockShipmentRepository) EXPECT() *MockShipmentRepositoryRecorder { return m.recorder }
func (m *MockShipmentRepository) Save(ctx context.Context, shipment domain.Shipment) error { ret := m.ctrl.Call(m, "Save", ctx, shipment); ret0, _ := ret[0].(error); return ret0 }
func (mr *MockShipmentRepositoryRecorder) Save(ctx, shipment any) *gomock.Call { return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockShipmentRepository)(nil).Save), ctx, shipment) }
func (m *MockShipmentRepository) GetByID(ctx context.Context, shipmentID string) (domain.Shipment, error) { ret := m.ctrl.Call(m, "GetByID", ctx, shipmentID); ret0, _ := ret[0].(domain.Shipment); ret1, _ := ret[1].(error); return ret0, ret1 }
func (mr *MockShipmentRepositoryRecorder) GetByID(ctx, shipmentID any) *gomock.Call { return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockShipmentRepository)(nil).GetByID), ctx, shipmentID) }

type MockShipmentEventPublisher struct{ ctrl *gomock.Controller; recorder *MockShipmentEventPublisherRecorder }
type MockShipmentEventPublisherRecorder struct{ mock *MockShipmentEventPublisher }
func NewMockShipmentEventPublisher(ctrl *gomock.Controller) *MockShipmentEventPublisher { m := &MockShipmentEventPublisher{ctrl: ctrl}; m.recorder = &MockShipmentEventPublisherRecorder{m}; return m }
func (m *MockShipmentEventPublisher) EXPECT() *MockShipmentEventPublisherRecorder { return m.recorder }
func (m *MockShipmentEventPublisher) Publish(ctx context.Context, event domain.Event) error { ret := m.ctrl.Call(m, "Publish", ctx, event); ret0, _ := ret[0].(error); return ret0 }
func (mr *MockShipmentEventPublisherRecorder) Publish(ctx, event any) *gomock.Call { return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Publish", reflect.TypeOf((*MockShipmentEventPublisher)(nil).Publish), ctx, event) }
