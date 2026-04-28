package ports

import (
	"context"

	"github.com/felo/felo-backend/services/auth-service/internal/domain"
)

type OTPStore interface {
	Save(ctx context.Context, otp domain.OTP) error
	GetByPhone(ctx context.Context, phone string) (domain.OTP, error)
	Delete(ctx context.Context, phone string) error
}

type SessionStore interface {
	Save(ctx context.Context, session domain.AuthSession) error
}

type TokenGenerator interface {
	GenerateToken(userID string, role domain.UserRole) (string, error)
}
