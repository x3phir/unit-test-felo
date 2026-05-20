package unit_test

import (
	"context"
	"reflect"

	"github.com/felo/felo-backend/services/invoice-service/internal/domain"
	"go.uber.org/mock/gomock"
)

type MockInvoiceRepository struct{ ctrl *gomock.Controller; recorder *MockInvoiceRepositoryRecorder }
type MockInvoiceRepositoryRecorder struct{ mock *MockInvoiceRepository }
func NewMockInvoiceRepository(ctrl *gomock.Controller) *MockInvoiceRepository { m := &MockInvoiceRepository{ctrl: ctrl}; m.recorder = &MockInvoiceRepositoryRecorder{m}; return m }
func (m *MockInvoiceRepository) EXPECT() *MockInvoiceRepositoryRecorder { return m.recorder }
func (m *MockInvoiceRepository) Create(ctx context.Context, invoice *domain.Invoice) error { ret := m.ctrl.Call(m, "Create", ctx, invoice); ret0, _ := ret[0].(error); return ret0 }
func (mr *MockInvoiceRepositoryRecorder) Create(ctx, invoice any) *gomock.Call { return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockInvoiceRepository)(nil).Create), ctx, invoice) }
func (m *MockInvoiceRepository) GetByID(ctx context.Context, id string) (*domain.Invoice, error) { ret := m.ctrl.Call(m, "GetByID", ctx, id); ret0, _ := ret[0].(*domain.Invoice); ret1, _ := ret[1].(error); return ret0, ret1 }
func (mr *MockInvoiceRepositoryRecorder) GetByID(ctx, id any) *gomock.Call { return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockInvoiceRepository)(nil).GetByID), ctx, id) }
func (m *MockInvoiceRepository) GetByOrderID(ctx context.Context, orderID string) ([]domain.Invoice, error) { ret := m.ctrl.Call(m, "GetByOrderID", ctx, orderID); ret0, _ := ret[0].([]domain.Invoice); ret1, _ := ret[1].(error); return ret0, ret1 }
func (mr *MockInvoiceRepositoryRecorder) GetByOrderID(ctx, orderID any) *gomock.Call { return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByOrderID", reflect.TypeOf((*MockInvoiceRepository)(nil).GetByOrderID), ctx, orderID) }
func (m *MockInvoiceRepository) UpdateStatus(ctx context.Context, id string, status domain.InvoiceStatus) error { ret := m.ctrl.Call(m, "UpdateStatus", ctx, id, status); ret0, _ := ret[0].(error); return ret0 }
func (mr *MockInvoiceRepositoryRecorder) UpdateStatus(ctx, id, status any) *gomock.Call { return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateStatus", reflect.TypeOf((*MockInvoiceRepository)(nil).UpdateStatus), ctx, id, status) }
func (m *MockInvoiceRepository) UpdatePaymentReference(ctx context.Context, id string, reference string) error { ret := m.ctrl.Call(m, "UpdatePaymentReference", ctx, id, reference); ret0, _ := ret[0].(error); return ret0 }
func (mr *MockInvoiceRepositoryRecorder) UpdatePaymentReference(ctx, id, reference any) *gomock.Call { return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdatePaymentReference", reflect.TypeOf((*MockInvoiceRepository)(nil).UpdatePaymentReference), ctx, id, reference) }

type MockNotificationPublisher struct{ ctrl *gomock.Controller; recorder *MockNotificationPublisherRecorder }
type MockNotificationPublisherRecorder struct{ mock *MockNotificationPublisher }
func NewMockNotificationPublisher(ctrl *gomock.Controller) *MockNotificationPublisher { m := &MockNotificationPublisher{ctrl: ctrl}; m.recorder = &MockNotificationPublisherRecorder{m}; return m }
func (m *MockNotificationPublisher) EXPECT() *MockNotificationPublisherRecorder { return m.recorder }
func (m *MockNotificationPublisher) PublishInvoiceNotification(ctx context.Context, invoice *domain.Invoice) error { ret := m.ctrl.Call(m, "PublishInvoiceNotification", ctx, invoice); ret0, _ := ret[0].(error); return ret0 }
func (mr *MockNotificationPublisherRecorder) PublishInvoiceNotification(ctx, invoice any) *gomock.Call { return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PublishInvoiceNotification", reflect.TypeOf((*MockNotificationPublisher)(nil).PublishInvoiceNotification), ctx, invoice) }
