package unit_test

import (
	"context"
	"reflect"

	"github.com/felo/felo-backend/services/matching-service/internal/domain"
	"go.uber.org/mock/gomock"
)

type MockAvailabilityReader struct{ ctrl *gomock.Controller; recorder *MockAvailabilityReaderRecorder }
type MockAvailabilityReaderRecorder struct{ mock *MockAvailabilityReader }
func NewMockAvailabilityReader(ctrl *gomock.Controller) *MockAvailabilityReader { m := &MockAvailabilityReader{ctrl: ctrl}; m.recorder = &MockAvailabilityReaderRecorder{m}; return m }
func (m *MockAvailabilityReader) EXPECT() *MockAvailabilityReaderRecorder { return m.recorder }
func (m *MockAvailabilityReader) FindAvailableDrivers(ctx context.Context, pickup domain.Coordinate, radius float64) ([]domain.DriverCandidate, error) {
	ret := m.ctrl.Call(m, "FindAvailableDrivers", ctx, pickup, radius); ret0, _ := ret[0].([]domain.DriverCandidate); ret1, _ := ret[1].(error); return ret0, ret1
}
func (mr *MockAvailabilityReaderRecorder) FindAvailableDrivers(ctx, pickup, radius any) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAvailableDrivers", reflect.TypeOf((*MockAvailabilityReader)(nil).FindAvailableDrivers), ctx, pickup, radius)
}

type MockAssignmentRepository struct{ ctrl *gomock.Controller; recorder *MockAssignmentRepositoryRecorder }
type MockAssignmentRepositoryRecorder struct{ mock *MockAssignmentRepository }
func NewMockAssignmentRepository(ctrl *gomock.Controller) *MockAssignmentRepository { m := &MockAssignmentRepository{ctrl: ctrl}; m.recorder = &MockAssignmentRepositoryRecorder{m}; return m }
func (m *MockAssignmentRepository) EXPECT() *MockAssignmentRepositoryRecorder { return m.recorder }
func (m *MockAssignmentRepository) Save(ctx context.Context, assignment domain.Assignment) error { ret := m.ctrl.Call(m, "Save", ctx, assignment); ret0, _ := ret[0].(error); return ret0 }
func (mr *MockAssignmentRepositoryRecorder) Save(ctx, assignment any) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockAssignmentRepository)(nil).Save), ctx, assignment)
}

type MockEventPublisher struct{ ctrl *gomock.Controller; recorder *MockEventPublisherRecorder }
type MockEventPublisherRecorder struct{ mock *MockEventPublisher }
func NewMockEventPublisher(ctrl *gomock.Controller) *MockEventPublisher { m := &MockEventPublisher{ctrl: ctrl}; m.recorder = &MockEventPublisherRecorder{m}; return m }
func (m *MockEventPublisher) EXPECT() *MockEventPublisherRecorder { return m.recorder }
func (m *MockEventPublisher) Publish(ctx context.Context, event domain.Event) error { ret := m.ctrl.Call(m, "Publish", ctx, event); ret0, _ := ret[0].(error); return ret0 }
func (mr *MockEventPublisherRecorder) Publish(ctx, event any) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Publish", reflect.TypeOf((*MockEventPublisher)(nil).Publish), ctx, event)
}
