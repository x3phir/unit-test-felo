package unit_test

import (
	"context"
	"reflect"

	"github.com/felo/felo-backend/services/payment-service/internal/domain"
	"go.uber.org/mock/gomock"
)

type MockWalletClient struct{ ctrl *gomock.Controller; recorder *MockWalletClientRecorder }
type MockWalletClientRecorder struct{ mock *MockWalletClient }
func NewMockWalletClient(ctrl *gomock.Controller) *MockWalletClient { m := &MockWalletClient{ctrl: ctrl}; m.recorder = &MockWalletClientRecorder{m}; return m }
func (m *MockWalletClient) EXPECT() *MockWalletClientRecorder { return m.recorder }
func (m *MockWalletClient) DebitCustomer(ctx context.Context, customerID string, amount int64, key string) error { ret := m.ctrl.Call(m, "DebitCustomer", ctx, customerID, amount, key); ret0, _ := ret[0].(error); return ret0 }
func (mr *MockWalletClientRecorder) DebitCustomer(ctx, customerID, amount, key any) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DebitCustomer", reflect.TypeOf((*MockWalletClient)(nil).DebitCustomer), ctx, customerID, amount, key)
}

type MockInvoiceClient struct{ ctrl *gomock.Controller; recorder *MockInvoiceClientRecorder }
type MockInvoiceClientRecorder struct{ mock *MockInvoiceClient }
func NewMockInvoiceClient(ctrl *gomock.Controller) *MockInvoiceClient { m := &MockInvoiceClient{ctrl: ctrl}; m.recorder = &MockInvoiceClientRecorder{m}; return m }
func (m *MockInvoiceClient) EXPECT() *MockInvoiceClientRecorder { return m.recorder }
func (m *MockInvoiceClient) IssueRideInvoice(ctx context.Context, tripID string, customerID string, amount int64, currency string) (string, error) {
	ret := m.ctrl.Call(m, "IssueRideInvoice", ctx, tripID, customerID, amount, currency); ret0, _ := ret[0].(string); ret1, _ := ret[1].(error); return ret0, ret1
}
func (mr *MockInvoiceClientRecorder) IssueRideInvoice(ctx, tripID, customerID, amount, currency any) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IssueRideInvoice", reflect.TypeOf((*MockInvoiceClient)(nil).IssueRideInvoice), ctx, tripID, customerID, amount, currency)
}

type MockProcessedEventStore struct{ ctrl *gomock.Controller; recorder *MockProcessedEventStoreRecorder }
type MockProcessedEventStoreRecorder struct{ mock *MockProcessedEventStore }
func NewMockProcessedEventStore(ctrl *gomock.Controller) *MockProcessedEventStore { m := &MockProcessedEventStore{ctrl: ctrl}; m.recorder = &MockProcessedEventStoreRecorder{m}; return m }
func (m *MockProcessedEventStore) EXPECT() *MockProcessedEventStoreRecorder { return m.recorder }
func (m *MockProcessedEventStore) Get(ctx context.Context, eventID string) (domain.PaymentResult, bool, error) { ret := m.ctrl.Call(m, "Get", ctx, eventID); ret0, _ := ret[0].(domain.PaymentResult); ret1, _ := ret[1].(bool); ret2, _ := ret[2].(error); return ret0, ret1, ret2 }
func (mr *MockProcessedEventStoreRecorder) Get(ctx, eventID any) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockProcessedEventStore)(nil).Get), ctx, eventID)
}
func (m *MockProcessedEventStore) Save(ctx context.Context, result domain.PaymentResult) error { ret := m.ctrl.Call(m, "Save", ctx, result); ret0, _ := ret[0].(error); return ret0 }
func (mr *MockProcessedEventStoreRecorder) Save(ctx, result any) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockProcessedEventStore)(nil).Save), ctx, result)
}

type MockEventPublisher struct{ ctrl *gomock.Controller; recorder *MockEventPublisherRecorder }
type MockEventPublisherRecorder struct{ mock *MockEventPublisher }
func NewMockPaymentEventPublisher(ctrl *gomock.Controller) *MockEventPublisher { m := &MockEventPublisher{ctrl: ctrl}; m.recorder = &MockEventPublisherRecorder{m}; return m }
func (m *MockEventPublisher) EXPECT() *MockEventPublisherRecorder { return m.recorder }
func (m *MockEventPublisher) Publish(ctx context.Context, event domain.Event) error { ret := m.ctrl.Call(m, "Publish", ctx, event); ret0, _ := ret[0].(error); return ret0 }
func (mr *MockEventPublisherRecorder) Publish(ctx, event any) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Publish", reflect.TypeOf((*MockEventPublisher)(nil).Publish), ctx, event)
}
