package service

import (
	"context"
	"errors"
	"time"

	"github.com/felo/felo-backend/services/feedback-service/internal/domain"
	"github.com/felo/felo-backend/services/feedback-service/internal/ports"
)

var (
	ErrInvalidFeedback = errors.New("invalid feedback input")
	ErrFeedbackExists  = errors.New("feedback already submitted for this trip")
	ErrFeedbackNotFound = errors.New("feedback not found")
)

type FeedbackService struct {
	repo      ports.FeedbackRepository
	publisher ports.EventPublisher
	now       func() time.Time
}

// NewFeedbackService membuat instance baru dari FeedbackService dengan repository dan event publisher yang diberikan.
func NewFeedbackService(repo ports.FeedbackRepository, publisher ports.EventPublisher) *FeedbackService {
	return &FeedbackService{
		repo:      repo,
		publisher: publisher,
		now:       time.Now,
	}
}

// SubmitFeedback merekam rating dan ulasan pengguna baru untuk perjalanan yang telah selesai.
// Fungsi ini juga mempublikasikan event "feedback.submitted.v1" agar layanan lain dapat memperbarui nilai rata-rata rating target.
func (s *FeedbackService) SubmitFeedback(ctx context.Context, tripID, userID, targetID, comment string, score int) (domain.Feedback, error) {
	if tripID == "" || userID == "" || targetID == "" || score < 1 || score > 5 {
		return domain.Feedback{}, ErrInvalidFeedback
	}

	_, err := s.repo.GetByTripID(ctx, tripID)
	if err == nil {
		return domain.Feedback{}, ErrFeedbackExists
	}

	now := s.now()
	feedback := domain.Feedback{
		ID:        "review-" + tripID,
		TripID:    tripID,
		UserID:    userID,
		TargetID:  targetID,
		Score:     score,
		Comment:   comment,
		CreatedAt: now,
	}

	if err := s.repo.Save(ctx, feedback); err != nil {
		return domain.Feedback{}, err
	}

	event := domain.Event{
		Name:       "feedback.submitted.v1",
		ReviewID:   feedback.ID,
		TargetID:   feedback.TargetID,
		Score:      feedback.Score,
		OccurredAt: now,
	}
	_ = s.publisher.Publish(ctx, event)

	return feedback, nil
}
