package unit_test

import (
	"context"
	"reflect"

	"github.com/felo/felo-backend/services/auth-service/internal/domain"
	"go.uber.org/mock/gomock"
)

type MockOTPStore struct{ ctrl *gomock.Controller; recorder *MockOTPStoreRecorder }
type MockOTPStoreRecorder struct{ mock *MockOTPStore }
func NewMockOTPStore(ctrl *gomock.Controller) *MockOTPStore { m := &MockOTPStore{ctrl: ctrl}; m.recorder = &MockOTPStoreRecorder{m}; return m }
func (m *MockOTPStore) EXPECT() *MockOTPStoreRecorder { return m.recorder }
func (m *MockOTPStore) Save(ctx context.Context, otp domain.OTP) error { ret := m.ctrl.Call(m, "Save", ctx, otp); ret0, _ := ret[0].(error); return ret0 }
func (mr *MockOTPStoreRecorder) Save(ctx, otp any) *gomock.Call { return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockOTPStore)(nil).Save), ctx, otp) }
func (m *MockOTPStore) GetByPhone(ctx context.Context, phone string) (domain.OTP, error) { ret := m.ctrl.Call(m, "GetByPhone", ctx, phone); ret0, _ := ret[0].(domain.OTP); ret1, _ := ret[1].(error); return ret0, ret1 }
func (mr *MockOTPStoreRecorder) GetByPhone(ctx, phone any) *gomock.Call { return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByPhone", reflect.TypeOf((*MockOTPStore)(nil).GetByPhone), ctx, phone) }
func (m *MockOTPStore) Delete(ctx context.Context, phone string) error { ret := m.ctrl.Call(m, "Delete", ctx, phone); ret0, _ := ret[0].(error); return ret0 }
func (mr *MockOTPStoreRecorder) Delete(ctx, phone any) *gomock.Call { return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockOTPStore)(nil).Delete), ctx, phone) }

type MockSessionStore struct{ ctrl *gomock.Controller; recorder *MockSessionStoreRecorder }
type MockSessionStoreRecorder struct{ mock *MockSessionStore }
func NewMockSessionStore(ctrl *gomock.Controller) *MockSessionStore { m := &MockSessionStore{ctrl: ctrl}; m.recorder = &MockSessionStoreRecorder{m}; return m }
func (m *MockSessionStore) EXPECT() *MockSessionStoreRecorder { return m.recorder }
func (m *MockSessionStore) Save(ctx context.Context, session domain.AuthSession) error { ret := m.ctrl.Call(m, "Save", ctx, session); ret0, _ := ret[0].(error); return ret0 }
func (mr *MockSessionStoreRecorder) Save(ctx, session any) *gomock.Call { return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockSessionStore)(nil).Save), ctx, session) }

type MockTokenGenerator struct{ ctrl *gomock.Controller; recorder *MockTokenGeneratorRecorder }
type MockTokenGeneratorRecorder struct{ mock *MockTokenGenerator }
func NewMockTokenGenerator(ctrl *gomock.Controller) *MockTokenGenerator { m := &MockTokenGenerator{ctrl: ctrl}; m.recorder = &MockTokenGeneratorRecorder{m}; return m }
func (m *MockTokenGenerator) EXPECT() *MockTokenGeneratorRecorder { return m.recorder }
func (m *MockTokenGenerator) GenerateToken(userID string, role domain.UserRole) (string, error) { ret := m.ctrl.Call(m, "GenerateToken", userID, role); ret0, _ := ret[0].(string); ret1, _ := ret[1].(error); return ret0, ret1 }
func (mr *MockTokenGeneratorRecorder) GenerateToken(userID, role any) *gomock.Call { return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateToken", reflect.TypeOf((*MockTokenGenerator)(nil).GenerateToken), userID, role) }
