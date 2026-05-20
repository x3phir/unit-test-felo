package unit_test

import (
	"context"
	"reflect"

	"github.com/felo/felo-backend/services/cart-service/internal/domain"
	"go.uber.org/mock/gomock"
)

type MockCartRepository struct{ ctrl *gomock.Controller; recorder *MockCartRepositoryRecorder }
type MockCartRepositoryRecorder struct{ mock *MockCartRepository }
func NewMockCartRepository(ctrl *gomock.Controller) *MockCartRepository { m := &MockCartRepository{ctrl: ctrl}; m.recorder = &MockCartRepositoryRecorder{m}; return m }
func (m *MockCartRepository) EXPECT() *MockCartRepositoryRecorder { return m.recorder }
func (m *MockCartRepository) Save(ctx context.Context, cart domain.Cart) error { ret := m.ctrl.Call(m, "Save", ctx, cart); ret0, _ := ret[0].(error); return ret0 }
func (mr *MockCartRepositoryRecorder) Save(ctx, cart any) *gomock.Call { return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockCartRepository)(nil).Save), ctx, cart) }
func (m *MockCartRepository) GetByUserID(ctx context.Context, userID string) (domain.Cart, bool, error) { ret := m.ctrl.Call(m, "GetByUserID", ctx, userID); ret0, _ := ret[0].(domain.Cart); ret1, _ := ret[1].(bool); ret2, _ := ret[2].(error); return ret0, ret1, ret2 }
func (mr *MockCartRepositoryRecorder) GetByUserID(ctx, userID any) *gomock.Call { return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByUserID", reflect.TypeOf((*MockCartRepository)(nil).GetByUserID), ctx, userID) }
func (m *MockCartRepository) Delete(ctx context.Context, userID string) error { ret := m.ctrl.Call(m, "Delete", ctx, userID); ret0, _ := ret[0].(error); return ret0 }
func (mr *MockCartRepositoryRecorder) Delete(ctx, userID any) *gomock.Call { return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockCartRepository)(nil).Delete), ctx, userID) }

type MockMerchantClient struct{ ctrl *gomock.Controller; recorder *MockMerchantClientRecorder }
type MockMerchantClientRecorder struct{ mock *MockMerchantClient }
func NewMockMerchantClient(ctrl *gomock.Controller) *MockMerchantClient { m := &MockMerchantClient{ctrl: ctrl}; m.recorder = &MockMerchantClientRecorder{m}; return m }
func (m *MockMerchantClient) EXPECT() *MockMerchantClientRecorder { return m.recorder }
func (m *MockMerchantClient) GetItemPriceAndAvailability(ctx context.Context, merchantID string, menuItemID string) (int64, bool, error) {
	ret := m.ctrl.Call(m, "GetItemPriceAndAvailability", ctx, merchantID, menuItemID); ret0, _ := ret[0].(int64); ret1, _ := ret[1].(bool); ret2, _ := ret[2].(error); return ret0, ret1, ret2
}
func (mr *MockMerchantClientRecorder) GetItemPriceAndAvailability(ctx, merchantID, menuItemID any) *gomock.Call { return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetItemPriceAndAvailability", reflect.TypeOf((*MockMerchantClient)(nil).GetItemPriceAndAvailability), ctx, merchantID, menuItemID) }
