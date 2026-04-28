package ports

import (
	"context"

	"github.com/felo/felo-backend/services/user-service/internal/domain"
)

type UserRepository interface {
	Save(ctx context.Context, user domain.UserProfile) error
	GetByID(ctx context.Context, userID string) (domain.UserProfile, error)
}
