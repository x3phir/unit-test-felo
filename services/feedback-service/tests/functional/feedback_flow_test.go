//go:build functional

package functional_test

import (
	"context"
	"os"
	"testing"

	"github.com/felo/felo-backend/services/feedback-service/internal/domain"
	"github.com/felo/felo-backend/services/feedback-service/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestFeedbackFunctional_SubmitFeedback(t *testing.T) {
	ctx := context.Background()
	db := openPG(t, getenv("FELO_FEEDBACK_PG_DSN", "postgres://felo:felo@127.0.0.1:54332/feedback_db?sslmode=disable"))
	t.Cleanup(func() { db.Close() })

	initDB(t, db)

	repo := &pgFeedbackRepository{db: db}
	pub := &mockEventPublisher{}
	svc := service.NewFeedbackService(repo, pub)

	tripID := "trip-001"
	userID := "user-001"
	targetID := "driver-001"
	
	// 1. Submit Feedback
	feedback, err := svc.SubmitFeedback(ctx, tripID, userID, targetID, "Great trip!", 5)
	if err != nil {
		t.Fatalf("SubmitFeedback() error = %v", err)
	}

	// Verify Feedback in DB
	var dbScore int
	var dbComment string
	if err := db.QueryRow(ctx, "SELECT score, comment FROM feedbacks WHERE id=$1", feedback.ID).Scan(&dbScore, &dbComment); err != nil {
		t.Fatalf("query persisted feedback: %v", err)
	}
	if dbScore != 5 {
		t.Fatalf("persisted score = %d, want 5", dbScore)
	}
	if dbComment != "Great trip!" {
		t.Fatalf("persisted comment = %s, want Great trip!", dbComment)
	}

	// Try submitting again for same trip (Should fail)
	_, err = svc.SubmitFeedback(ctx, tripID, userID, targetID, "Another comment", 4)
	if err == nil {
		t.Fatalf("expected error when submitting feedback for same trip twice")
	}

	// Verify publisher was called
	if !pub.called {
		t.Fatalf("expected EventPublisher to be called")
	}
}

func initDB(t *testing.T, db *pgxpool.Pool) {
	t.Helper()
	ctx := context.Background()
	_, err := db.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS feedbacks (
			id VARCHAR(50) PRIMARY KEY,
			trip_id VARCHAR(50) UNIQUE NOT NULL,
			user_id VARCHAR(50) NOT NULL,
			target_id VARCHAR(50) NOT NULL,
			score INT NOT NULL,
			comment TEXT,
			created_at TIMESTAMP NOT NULL
		);
		TRUNCATE feedbacks;
	`)
	if err != nil {
		t.Fatalf("initDB() error = %v", err)
	}
}

type pgFeedbackRepository struct{ db *pgxpool.Pool }

func (r *pgFeedbackRepository) Save(ctx context.Context, feedback domain.Feedback) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO feedbacks (id, trip_id, user_id, target_id, score, comment, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, feedback.ID, feedback.TripID, feedback.UserID, feedback.TargetID, feedback.Score, feedback.Comment, feedback.CreatedAt)
	return err
}

func (r *pgFeedbackRepository) GetByTripID(ctx context.Context, tripID string) (domain.Feedback, error) {
	var feedback domain.Feedback
	err := r.db.QueryRow(ctx, "SELECT id, trip_id, user_id, target_id, score, comment, created_at FROM feedbacks WHERE trip_id=$1", tripID).
		Scan(&feedback.ID, &feedback.TripID, &feedback.UserID, &feedback.TargetID, &feedback.Score, &feedback.Comment, &feedback.CreatedAt)
	return feedback, err
}

type mockEventPublisher struct {
	called bool
}

func (m *mockEventPublisher) Publish(ctx context.Context, event domain.Event) error {
	m.called = true
	return nil
}

func openPG(t *testing.T, dsn string) *pgxpool.Pool {
	t.Helper()
	db, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		t.Fatalf("pgxpool.New() error = %v", err)
	}
	return db
}

func getenv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
