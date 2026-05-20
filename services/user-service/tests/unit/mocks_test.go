package unit_test

import (
	"context"
	"reflect"

	"github.com/felo/felo-backend/services/user-service/internal/domain"
	"go.uber.org/mock/gomock"
)

type MockUserRepository struct{ ctrl *gomock.Controller; recorder *MockUserRepositoryRecorder }
type MockUserRepositoryRecorder struct{ mock *MockUserRepository }
func NewMockUserRepository(ctrl *gomock.Controller) *MockUserRepository { m := &MockUserRepository{ctrl: ctrl}; m.recorder = &MockUserRepositoryRecorder{m}; return m }
func (m *MockUserRepository) EXPECT() *MockUserRepositoryRecorder { return m.recorder }
func (m *MockUserRepository) Save(ctx context.Context, user domain.UserProfile) error { ret := m.ctrl.Call(m, "Save", ctx, user); ret0, _ := ret[0].(error); return ret0 }
func (mr *MockUserRepositoryRecorder) Save(ctx, user any) *gomock.Call { return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockUserRepository)(nil).Save), ctx, user) }
func (m *MockUserRepository) GetByID(ctx context.Context, userID string) (domain.UserProfile, error) { ret := m.ctrl.Call(m, "GetByID", ctx, userID); ret0, _ := ret[0].(domain.UserProfile); ret1, _ := ret[1].(error); return ret0, ret1 }
func (mr *MockUserRepositoryRecorder) GetByID(ctx, userID any) *gomock.Call { return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockUserRepository)(nil).GetByID), ctx, userID) }
