package unit_test

import (
	"context"
	"reflect"
	"time"

	"github.com/felo/felo-backend/services/ride-service/internal/domain"
	"go.uber.org/mock/gomock"
)

type MockTripRepository struct {
	ctrl     *gomock.Controller
	recorder *MockTripRepositoryRecorder
}

type MockTripRepositoryRecorder struct{ mock *MockTripRepository }

func NewMockTripRepository(ctrl *gomock.Controller) *MockTripRepository {
	mock := &MockTripRepository{ctrl: ctrl}
	mock.recorder = &MockTripRepositoryRecorder{mock}
	return mock
}

func (m *MockTripRepository) EXPECT() *MockTripRepositoryRecorder { return m.recorder }
func (m *MockTripRepository) Save(ctx context.Context, trip domain.Trip) error {
	ret := m.ctrl.Call(m, "Save", ctx, trip)
	ret0, _ := ret[0].(error)
	return ret0
}
func (mr *MockTripRepositoryRecorder) Save(ctx, trip any) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockTripRepository)(nil).Save), ctx, trip)
}
func (m *MockTripRepository) GetByID(ctx context.Context, tripID string) (domain.Trip, error) {
	ret := m.ctrl.Call(m, "GetByID", ctx, tripID)
	ret0, _ := ret[0].(domain.Trip)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}
func (mr *MockTripRepositoryRecorder) GetByID(ctx, tripID any) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockTripRepository)(nil).GetByID), ctx, tripID)
}

type MockEventPublisher struct {
	ctrl     *gomock.Controller
	recorder *MockEventPublisherRecorder
}

type MockEventPublisherRecorder struct{ mock *MockEventPublisher }

func NewMockEventPublisher(ctrl *gomock.Controller) *MockEventPublisher {
	mock := &MockEventPublisher{ctrl: ctrl}
	mock.recorder = &MockEventPublisherRecorder{mock}
	return mock
}

func (m *MockEventPublisher) EXPECT() *MockEventPublisherRecorder { return m.recorder }
func (m *MockEventPublisher) Publish(ctx context.Context, event domain.Event) error {
	ret := m.ctrl.Call(m, "Publish", ctx, event)
	ret0, _ := ret[0].(error)
	return ret0
}
func (mr *MockEventPublisherRecorder) Publish(ctx, event any) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Publish", reflect.TypeOf((*MockEventPublisher)(nil).Publish), ctx, event)
}

type fixedClock struct{ now time.Time }
func (c fixedClock) Now() time.Time { return c.now }

type idGen struct{ ids []string; idx int }
func (g *idGen) NewID() string { id := g.ids[g.idx]; g.idx++; return id }
