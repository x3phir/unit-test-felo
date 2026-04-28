package ports

import (
	"context"

	"github.com/felo/felo-backend/services/driver-service/internal/domain"
)

type DriverRepository interface {
	Save(ctx context.Context, driver domain.DriverProfile) error
	GetByID(ctx context.Context, driverID string) (domain.DriverProfile, error)
}
