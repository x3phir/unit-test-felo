package service

import (
	"context"
	"errors"
	"time"

	"github.com/felo/felo-backend/services/driver-service/internal/domain"
	"github.com/felo/felo-backend/services/driver-service/internal/ports"
)

var (
	ErrDriverNotFound = errors.New("driver not found")
	ErrInvalidInput   = errors.New("invalid driver input")
	ErrKYCNotApproved = errors.New("driver KYC is not approved")
)

type DriverService struct {
	repo ports.DriverRepository
	now  func() time.Time
}

// NewDriverService membuat instance baru dari DriverService dengan repository yang diberikan.
func NewDriverService(repo ports.DriverRepository) *DriverService {
	return &DriverService{
		repo: repo,
		now:  time.Now,
	}
}

// RegisterDriver mendaftarkan profil driver baru dengan detail dan info kendaraan yang diberikan.
// Status KYC awal driver akan diatur menjadi Pending dan status operasional menjadi Offline.
func (s *DriverService) RegisterDriver(ctx context.Context, id, name, phone string, vehicle domain.VehicleInfo) (domain.DriverProfile, error) {
	if id == "" || name == "" || phone == "" || vehicle.LicensePlate == "" || vehicle.Type == "" {
		return domain.DriverProfile{}, ErrInvalidInput
	}

	driver := domain.DriverProfile{
		ID:                id,
		Name:              name,
		Phone:             phone,
		Vehicle:           vehicle,
		KYCStatus:         domain.KYCPending,
		OperationalStatus: domain.StatusOffline,
		Rating:            0.0,
		CreatedAt:         s.now(),
		UpdatedAt:         s.now(),
	}

	if err := s.repo.Save(ctx, driver); err != nil {
		return domain.DriverProfile{}, err
	}

	return driver, nil
}

// ApproveKYC menyetujui dokumen KYC (Know Your Customer) milik driver,
// memungkinkan mereka untuk mengubah status operasionalnya menjadi Online.
func (s *DriverService) ApproveKYC(ctx context.Context, id string) (domain.DriverProfile, error) {
	driver, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return domain.DriverProfile{}, err
	}

	driver.KYCStatus = domain.KYCApproved
	driver.UpdatedAt = s.now()

	if err := s.repo.Save(ctx, driver); err != nil {
		return domain.DriverProfile{}, err
	}

	return driver, nil
}

// SetOperationalStatus memperbarui status ketersediaan driver saat ini (contoh: Online, Offline, Busy).
// Mengembalikan ErrKYCNotApproved jika driver mencoba untuk Online tanpa KYC yang telah disetujui.
func (s *DriverService) SetOperationalStatus(ctx context.Context, id string, status domain.OperationalStatus) (domain.DriverProfile, error) {
	driver, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return domain.DriverProfile{}, err
	}

	if status == domain.StatusOnline && driver.KYCStatus != domain.KYCApproved {
		return domain.DriverProfile{}, ErrKYCNotApproved
	}

	driver.OperationalStatus = status
	driver.UpdatedAt = s.now()

	if err := s.repo.Save(ctx, driver); err != nil {
		return domain.DriverProfile{}, err
	}

	return driver, nil
}
