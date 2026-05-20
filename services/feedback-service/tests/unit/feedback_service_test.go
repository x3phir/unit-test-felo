package unit_test

import (
	"context"
	"errors"
	"testing"

	"github.com/felo/felo-backend/services/feedback-service/internal/domain"
	"github.com/felo/felo-backend/services/feedback-service/internal/service"
)

func TestFeedbackService_SubmitFeedback_Success(t *testing.T) {
	repo := &feedbackRepoFake{feedbacks: map[string]domain.Feedback{}}
	publisher := &eventPublisherFake{}
	svc := service.NewFeedbackService(repo, publisher)

	feedback, err := svc.SubmitFeedback(context.Background(), "trip-1", "user-1", "driver-1", "Great!", 5)
	if err != nil {
		t.Fatalf("SubmitFeedback() error = %v", err)
	}

	if feedback.TripID != "trip-1" || feedback.Score != 5 {
		t.Fatalf("unexpected feedback: %+v", feedback)
	}

	if len(publisher.events) != 1 || publisher.events[0].Name != "feedback.submitted.v1" {
		t.Fatalf("event not published properly: %+v", publisher.events)
	}
}

func TestFeedbackService_SubmitFeedback_InvalidInput(t *testing.T) {
	svc := service.NewFeedbackService(&feedbackRepoFake{}, &eventPublisherFake{})

	_, err := svc.SubmitFeedback(context.Background(), "", "user-1", "driver-1", "Great!", 5)
	if !errors.Is(err, service.ErrInvalidFeedback) {
		t.Fatalf("SubmitFeedback() error = %v, want ErrInvalidFeedback", err)
	}

	_, err = svc.SubmitFeedback(context.Background(), "trip-1", "user-1", "driver-1", "Great!", 6)
	if !errors.Is(err, service.ErrInvalidFeedback) {
		t.Fatalf("SubmitFeedback() error = %v, want ErrInvalidFeedback for score > 5", err)
	}
}

func TestFeedbackService_SubmitFeedback_AlreadyExists(t *testing.T) {
	repo := &feedbackRepoFake{feedbacks: map[string]domain.Feedback{
		"trip-1": {TripID: "trip-1"},
	}}
	svc := service.NewFeedbackService(repo, &eventPublisherFake{})

	_, err := svc.SubmitFeedback(context.Background(), "trip-1", "user-1", "driver-1", "Great!", 5)
	if !errors.Is(err, service.ErrFeedbackExists) {
		t.Fatalf("SubmitFeedback() error = %v, want ErrFeedbackExists", err)
	}
}

type feedbackRepoFake struct {
	feedbacks map[string]domain.Feedback
}

func (r *feedbackRepoFake) Save(_ context.Context, f domain.Feedback) error {
	r.feedbacks[f.TripID] = f
	return nil
}

func (r *feedbackRepoFake) GetByTripID(_ context.Context, tripID string) (domain.Feedback, error) {
	f, ok := r.feedbacks[tripID]
	if !ok {
		return domain.Feedback{}, service.ErrFeedbackNotFound
	}
	return f, nil
}

type eventPublisherFake struct {
	events []domain.Event
}

func (p *eventPublisherFake) Publish(_ context.Context, event domain.Event) error {
	p.events = append(p.events, event)
	return nil
}
