package unit_test

import (
	"context"
	"reflect"

	"github.com/felo/felo-backend/services/send-order-service/internal/domain"
	"go.uber.org/mock/gomock"
)

type MockSendOrderRepository struct{ ctrl *gomock.Controller; recorder *MockSendOrderRepositoryRecorder }
type MockSendOrderRepositoryRecorder struct{ mock *MockSendOrderRepository }
func NewMockSendOrderRepository(ctrl *gomock.Controller) *MockSendOrderRepository { m := &MockSendOrderRepository{ctrl: ctrl}; m.recorder = &MockSendOrderRepositoryRecorder{m}; return m }
func (m *MockSendOrderRepository) EXPECT() *MockSendOrderRepositoryRecorder { return m.recorder }
func (m *MockSendOrderRepository) Save(ctx context.Context, order domain.SendOrder) error { ret := m.ctrl.Call(m, "Save", ctx, order); ret0, _ := ret[0].(error); return ret0 }
func (mr *MockSendOrderRepositoryRecorder) Save(ctx, order any) *gomock.Call { return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockSendOrderRepository)(nil).Save), ctx, order) }

type MockPricingClient struct{ ctrl *gomock.Controller; recorder *MockPricingClientRecorder }
type MockPricingClientRecorder struct{ mock *MockPricingClient }
func NewMockPricingClient(ctrl *gomock.Controller) *MockPricingClient { m := &MockPricingClient{ctrl: ctrl}; m.recorder = &MockPricingClientRecorder{m}; return m }
func (m *MockPricingClient) EXPECT() *MockPricingClientRecorder { return m.recorder }
func (m *MockPricingClient) CalculateShippingFee(ctx context.Context, pkg domain.PackageDetails, origin string, destination string) (int64, error) { ret := m.ctrl.Call(m, "CalculateShippingFee", ctx, pkg, origin, destination); ret0, _ := ret[0].(int64); ret1, _ := ret[1].(error); return ret0, ret1 }
func (mr *MockPricingClientRecorder) CalculateShippingFee(ctx, pkg, origin, destination any) *gomock.Call { return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CalculateShippingFee", reflect.TypeOf((*MockPricingClient)(nil).CalculateShippingFee), ctx, pkg, origin, destination) }

type MockInvoiceClient struct{ ctrl *gomock.Controller; recorder *MockInvoiceClientRecorder }
type MockInvoiceClientRecorder struct{ mock *MockInvoiceClient }
func NewMockSendOrderInvoiceClient(ctrl *gomock.Controller) *MockInvoiceClient { m := &MockInvoiceClient{ctrl: ctrl}; m.recorder = &MockInvoiceClientRecorder{m}; return m }
func (m *MockInvoiceClient) EXPECT() *MockInvoiceClientRecorder { return m.recorder }
func (m *MockInvoiceClient) CreateInvoice(ctx context.Context, orderID string, payerID string, payerType domain.PayerType, amount int64) error { ret := m.ctrl.Call(m, "CreateInvoice", ctx, orderID, payerID, payerType, amount); ret0, _ := ret[0].(error); return ret0 }
func (mr *MockInvoiceClientRecorder) CreateInvoice(ctx, orderID, payerID, payerType, amount any) *gomock.Call { return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateInvoice", reflect.TypeOf((*MockInvoiceClient)(nil).CreateInvoice), ctx, orderID, payerID, payerType, amount) }

type MockSendOrderEventPublisher struct{ ctrl *gomock.Controller; recorder *MockSendOrderEventPublisherRecorder }
type MockSendOrderEventPublisherRecorder struct{ mock *MockSendOrderEventPublisher }
func NewMockSendOrderEventPublisher(ctrl *gomock.Controller) *MockSendOrderEventPublisher { m := &MockSendOrderEventPublisher{ctrl: ctrl}; m.recorder = &MockSendOrderEventPublisherRecorder{m}; return m }
func (m *MockSendOrderEventPublisher) EXPECT() *MockSendOrderEventPublisherRecorder { return m.recorder }
func (m *MockSendOrderEventPublisher) Publish(ctx context.Context, event domain.Event) error { ret := m.ctrl.Call(m, "Publish", ctx, event); ret0, _ := ret[0].(error); return ret0 }
func (mr *MockSendOrderEventPublisherRecorder) Publish(ctx, event any) *gomock.Call { return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Publish", reflect.TypeOf((*MockSendOrderEventPublisher)(nil).Publish), ctx, event) }

type MockSendOrderIDGenerator struct{ ctrl *gomock.Controller; recorder *MockSendOrderIDGeneratorRecorder }
type MockSendOrderIDGeneratorRecorder struct{ mock *MockSendOrderIDGenerator }
func NewMockSendOrderIDGenerator(ctrl *gomock.Controller) *MockSendOrderIDGenerator { m := &MockSendOrderIDGenerator{ctrl: ctrl}; m.recorder = &MockSendOrderIDGeneratorRecorder{m}; return m }
func (m *MockSendOrderIDGenerator) EXPECT() *MockSendOrderIDGeneratorRecorder { return m.recorder }
func (m *MockSendOrderIDGenerator) NewID() string { ret := m.ctrl.Call(m, "NewID"); ret0, _ := ret[0].(string); return ret0 }
func (mr *MockSendOrderIDGeneratorRecorder) NewID() *gomock.Call { return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewID", reflect.TypeOf((*MockSendOrderIDGenerator)(nil).NewID)) }
