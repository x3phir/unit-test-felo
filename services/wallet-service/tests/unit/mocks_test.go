package unit_test

import (
	"context"
	"reflect"

	"github.com/felo/felo-backend/services/wallet-service/internal/domain"
	"go.uber.org/mock/gomock"
)

type MockSettlementStore struct{ ctrl *gomock.Controller; recorder *MockSettlementStoreRecorder }
type MockSettlementStoreRecorder struct{ mock *MockSettlementStore }
func NewMockSettlementStore(ctrl *gomock.Controller) *MockSettlementStore { m := &MockSettlementStore{ctrl: ctrl}; m.recorder = &MockSettlementStoreRecorder{m}; return m }
func (m *MockSettlementStore) EXPECT() *MockSettlementStoreRecorder { return m.recorder }
func (m *MockSettlementStore) GetByKey(ctx context.Context, key string) (domain.SettlementRecord, bool, error) { ret := m.ctrl.Call(m, "GetByKey", ctx, key); ret0, _ := ret[0].(domain.SettlementRecord); ret1, _ := ret[1].(bool); ret2, _ := ret[2].(error); return ret0, ret1, ret2 }
func (mr *MockSettlementStoreRecorder) GetByKey(ctx, key any) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByKey", reflect.TypeOf((*MockSettlementStore)(nil).GetByKey), ctx, key)
}
func (m *MockSettlementStore) Save(ctx context.Context, record domain.SettlementRecord) error { ret := m.ctrl.Call(m, "Save", ctx, record); ret0, _ := ret[0].(error); return ret0 }
func (mr *MockSettlementStoreRecorder) Save(ctx, record any) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockSettlementStore)(nil).Save), ctx, record)
}

type MockWalletStore struct{ ctrl *gomock.Controller; recorder *MockWalletStoreRecorder }
type MockWalletStoreRecorder struct{ mock *MockWalletStore }
func NewMockWalletStore(ctrl *gomock.Controller) *MockWalletStore { m := &MockWalletStore{ctrl: ctrl}; m.recorder = &MockWalletStoreRecorder{m}; return m }
func (m *MockWalletStore) EXPECT() *MockWalletStoreRecorder { return m.recorder }
func (m *MockWalletStore) Credit(ctx context.Context, ownerID string, amount int64) (int64, error) { ret := m.ctrl.Call(m, "Credit", ctx, ownerID, amount); ret0, _ := ret[0].(int64); ret1, _ := ret[1].(error); return ret0, ret1 }
func (mr *MockWalletStoreRecorder) Credit(ctx, ownerID, amount any) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Credit", reflect.TypeOf((*MockWalletStore)(nil).Credit), ctx, ownerID, amount)
}
