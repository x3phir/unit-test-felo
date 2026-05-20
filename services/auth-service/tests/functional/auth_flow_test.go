//go:build functional

package functional_test

import (
	"context"
	"os"
	"testing"

	"github.com/felo/felo-backend/services/auth-service/internal/domain"
	"github.com/felo/felo-backend/services/auth-service/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestAuthFunctional_RequestAndLoginOTP(t *testing.T) {
	ctx := context.Background()
	db := openPG(t, getenv("FELO_AUTH_PG_DSN", "postgres://felo:felo@127.0.0.1:54329/auth_db?sslmode=disable"))
	t.Cleanup(func() { db.Close() })

	initDB(t, db)

	otpStore := &pgOTPStore{db: db}
	sessionStore := &pgSessionStore{db: db}
	tokenGen := &mockTokenGenerator{}
	svc := service.NewAuthService(otpStore, sessionStore, tokenGen)

	phone := "+628123456789"
	userID := "user-001"
	
	// 1. Request OTP
	otp, err := svc.RequestOTP(ctx, phone)
	if err != nil {
		t.Fatalf("RequestOTP() error = %v", err)
	}
	
	// Verify OTP is in DB
	var dbCode string
	if err := db.QueryRow(ctx, "SELECT code FROM otps WHERE phone=$1", phone).Scan(&dbCode); err != nil {
		t.Fatalf("query persisted otp: %v", err)
	}
	if dbCode != otp.Code {
		t.Fatalf("persisted code = %s, want %s", dbCode, otp.Code)
	}

	// 2. Login with OTP
	session, err := svc.LoginWithOTP(ctx, userID, phone, otp.Code, domain.RolePassenger)
	if err != nil {
		t.Fatalf("LoginWithOTP() error = %v", err)
	}

	// Verify Session is in DB
	var dbUserID string
	if err := db.QueryRow(ctx, "SELECT user_id FROM sessions WHERE id=$1", session.ID).Scan(&dbUserID); err != nil {
		t.Fatalf("query persisted session: %v", err)
	}
	if dbUserID != userID {
		t.Fatalf("persisted session user_id = %s, want %s", dbUserID, userID)
	}
}

func initDB(t *testing.T, db *pgxpool.Pool) {
	t.Helper()
	ctx := context.Background()
	_, err := db.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS otps (
			id VARCHAR(50) PRIMARY KEY,
			phone VARCHAR(20) UNIQUE NOT NULL,
			code VARCHAR(10) NOT NULL,
			expires_at TIMESTAMP NOT NULL
		);
		CREATE TABLE IF NOT EXISTS sessions (
			id VARCHAR(50) PRIMARY KEY,
			user_id VARCHAR(50) NOT NULL,
			role VARCHAR(20) NOT NULL,
			token TEXT NOT NULL,
			expires_at TIMESTAMP NOT NULL
		);
		TRUNCATE otps, sessions;
	`)
	if err != nil {
		t.Fatalf("initDB() error = %v", err)
	}
}

type pgOTPStore struct{ db *pgxpool.Pool }

func (r *pgOTPStore) Save(ctx context.Context, otp domain.OTP) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO otps (id, phone, code, expires_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (phone) DO UPDATE SET
			code = EXCLUDED.code,
			expires_at = EXCLUDED.expires_at,
			id = EXCLUDED.id
	`, otp.ID, otp.Phone, otp.Code, otp.ExpiresAt)
	return err
}

func (r *pgOTPStore) GetByPhone(ctx context.Context, phone string) (domain.OTP, error) {
	var otp domain.OTP
	err := r.db.QueryRow(ctx, "SELECT id, phone, code, expires_at FROM otps WHERE phone=$1", phone).
		Scan(&otp.ID, &otp.Phone, &otp.Code, &otp.ExpiresAt)
	return otp, err
}

func (r *pgOTPStore) Delete(ctx context.Context, phone string) error {
	_, err := r.db.Exec(ctx, "DELETE FROM otps WHERE phone=$1", phone)
	return err
}

type pgSessionStore struct{ db *pgxpool.Pool }

func (r *pgSessionStore) Save(ctx context.Context, session domain.AuthSession) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO sessions (id, user_id, role, token, expires_at)
		VALUES ($1, $2, $3, $4, $5)
	`, session.ID, session.UserID, string(session.Role), session.Token, session.ExpiresAt)
	return err
}

type mockTokenGenerator struct{}

func (mockTokenGenerator) GenerateToken(userID string, role domain.UserRole) (string, error) {
	return "mock-token-for-" + userID, nil
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
