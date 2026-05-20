package unit_test

import (
	"context"
	"reflect"

	"github.com/felo/felo-backend/services/notification-service/internal/domain"
	"go.uber.org/mock/gomock"
)

type MockNotificationProvider struct{ ctrl *gomock.Controller; recorder *MockNotificationProviderRecorder }
type MockNotificationProviderRecorder struct{ mock *MockNotificationProvider }
func NewMockNotificationProvider(ctrl *gomock.Controller) *MockNotificationProvider { m := &MockNotificationProvider{ctrl: ctrl}; m.recorder = &MockNotificationProviderRecorder{m}; return m }
func (m *MockNotificationProvider) EXPECT() *MockNotificationProviderRecorder { return m.recorder }
func (m *MockNotificationProvider) SendPush(ctx context.Context, req domain.NotificationRequest) (domain.DeliveryStatus, error) { ret := m.ctrl.Call(m, "SendPush", ctx, req); ret0, _ := ret[0].(domain.DeliveryStatus); ret1, _ := ret[1].(error); return ret0, ret1 }
func (mr *MockNotificationProviderRecorder) SendPush(ctx, req any) *gomock.Call { return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendPush", reflect.TypeOf((*MockNotificationProvider)(nil).SendPush), ctx, req) }
func (m *MockNotificationProvider) SendWhatsApp(ctx context.Context, req domain.NotificationRequest) (domain.DeliveryStatus, error) { ret := m.ctrl.Call(m, "SendWhatsApp", ctx, req); ret0, _ := ret[0].(domain.DeliveryStatus); ret1, _ := ret[1].(error); return ret0, ret1 }
func (mr *MockNotificationProviderRecorder) SendWhatsApp(ctx, req any) *gomock.Call { return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendWhatsApp", reflect.TypeOf((*MockNotificationProvider)(nil).SendWhatsApp), ctx, req) }
func (m *MockNotificationProvider) SendSMS(ctx context.Context, req domain.NotificationRequest) (domain.DeliveryStatus, error) { ret := m.ctrl.Call(m, "SendSMS", ctx, req); ret0, _ := ret[0].(domain.DeliveryStatus); ret1, _ := ret[1].(error); return ret0, ret1 }
func (mr *MockNotificationProviderRecorder) SendSMS(ctx, req any) *gomock.Call { return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendSMS", reflect.TypeOf((*MockNotificationProvider)(nil).SendSMS), ctx, req) }
