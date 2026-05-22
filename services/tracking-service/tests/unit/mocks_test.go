package unit_test

import (
	"context"
	"reflect"
	"time"

	"github.com/felo/felo-backend/services/tracking-service/internal/domain"
	"go.uber.org/mock/gomock"
)

type MockTrackingSessionRepository struct {
	ctrl     *gomock.Controller
	recorder *MockTrackingSessionRepositoryRecorder
}

type MockTrackingSessionRepositoryRecorder struct{ mock *MockTrackingSessionRepository }

func NewMockTrackingSessionRepository(ctrl *gomock.Controller) *MockTrackingSessionRepository {
	mock := &MockTrackingSessionRepository{ctrl: ctrl}
	mock.recorder = &MockTrackingSessionRepositoryRecorder{mock}
	return mock
}

func (m *MockTrackingSessionRepository) EXPECT() *MockTrackingSessionRepositoryRecorder { return m.recorder }
func (m *MockTrackingSessionRepository) Save(ctx context.Context, session domain.TrackingSession) error {
	ret := m.ctrl.Call(m, "Save", ctx, session)
	ret0, _ := ret[0].(error)
	return ret0
}
func (mr *MockTrackingSessionRepositoryRecorder) Save(ctx, session any) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockTrackingSessionRepository)(nil).Save), ctx, session)
}
func (m *MockTrackingSessionRepository) GetByID(ctx context.Context, sessionID string) (domain.TrackingSession, error) {
	ret := m.ctrl.Call(m, "GetByID", ctx, sessionID)
	ret0, _ := ret[0].(domain.TrackingSession)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}
func (mr *MockTrackingSessionRepositoryRecorder) GetByID(ctx, sessionID any) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockTrackingSessionRepository)(nil).GetByID), ctx, sessionID)
}

type MockTrackingRecordRepository struct {
	ctrl     *gomock.Controller
	recorder *MockTrackingRecordRepositoryRecorder
}

type MockTrackingRecordRepositoryRecorder struct{ mock *MockTrackingRecordRepository }

func NewMockTrackingRecordRepository(ctrl *gomock.Controller) *MockTrackingRecordRepository {
	mock := &MockTrackingRecordRepository{ctrl: ctrl}
	mock.recorder = &MockTrackingRecordRepositoryRecorder{mock}
	return mock
}

func (m *MockTrackingRecordRepository) EXPECT() *MockTrackingRecordRepositoryRecorder { return m.recorder }
func (m *MockTrackingRecordRepository) Save(ctx context.Context, record domain.TrackingRecord) error {
	ret := m.ctrl.Call(m, "Save", ctx, record)
	ret0, _ := ret[0].(error)
	return ret0
}
func (mr *MockTrackingRecordRepositoryRecorder) Save(ctx, record any) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockTrackingRecordRepository)(nil).Save), ctx, record)
}
func (m *MockTrackingRecordRepository) ListBySession(ctx context.Context, sessionID string) ([]domain.TrackingRecord, error) {
	ret := m.ctrl.Call(m, "ListBySession", ctx, sessionID)
	ret0, _ := ret[0].([]domain.TrackingRecord)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}
func (mr *MockTrackingRecordRepositoryRecorder) ListBySession(ctx, sessionID any) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListBySession", reflect.TypeOf((*MockTrackingRecordRepository)(nil).ListBySession), ctx, sessionID)
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
