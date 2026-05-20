package unit_test

import (
	"context"
	"reflect"

	"github.com/felo/felo-backend/services/feedback-service/internal/domain"
	"go.uber.org/mock/gomock"
)

type MockFeedbackRepository struct{ ctrl *gomock.Controller; recorder *MockFeedbackRepositoryRecorder }
type MockFeedbackRepositoryRecorder struct{ mock *MockFeedbackRepository }
func NewMockFeedbackRepository(ctrl *gomock.Controller) *MockFeedbackRepository { m := &MockFeedbackRepository{ctrl: ctrl}; m.recorder = &MockFeedbackRepositoryRecorder{m}; return m }
func (m *MockFeedbackRepository) EXPECT() *MockFeedbackRepositoryRecorder { return m.recorder }
func (m *MockFeedbackRepository) Save(ctx context.Context, feedback domain.Feedback) error { ret := m.ctrl.Call(m, "Save", ctx, feedback); ret0, _ := ret[0].(error); return ret0 }
func (mr *MockFeedbackRepositoryRecorder) Save(ctx, feedback any) *gomock.Call { return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockFeedbackRepository)(nil).Save), ctx, feedback) }
func (m *MockFeedbackRepository) GetByTripID(ctx context.Context, tripID string) (domain.Feedback, error) { ret := m.ctrl.Call(m, "GetByTripID", ctx, tripID); ret0, _ := ret[0].(domain.Feedback); ret1, _ := ret[1].(error); return ret0, ret1 }
func (mr *MockFeedbackRepositoryRecorder) GetByTripID(ctx, tripID any) *gomock.Call { return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByTripID", reflect.TypeOf((*MockFeedbackRepository)(nil).GetByTripID), ctx, tripID) }

type MockFeedbackEventPublisher struct{ ctrl *gomock.Controller; recorder *MockFeedbackEventPublisherRecorder }
type MockFeedbackEventPublisherRecorder struct{ mock *MockFeedbackEventPublisher }
func NewMockFeedbackEventPublisher(ctrl *gomock.Controller) *MockFeedbackEventPublisher { m := &MockFeedbackEventPublisher{ctrl: ctrl}; m.recorder = &MockFeedbackEventPublisherRecorder{m}; return m }
func (m *MockFeedbackEventPublisher) EXPECT() *MockFeedbackEventPublisherRecorder { return m.recorder }
func (m *MockFeedbackEventPublisher) Publish(ctx context.Context, event domain.Event) error { ret := m.ctrl.Call(m, "Publish", ctx, event); ret0, _ := ret[0].(error); return ret0 }
func (mr *MockFeedbackEventPublisherRecorder) Publish(ctx, event any) *gomock.Call { return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Publish", reflect.TypeOf((*MockFeedbackEventPublisher)(nil).Publish), ctx, event) }
