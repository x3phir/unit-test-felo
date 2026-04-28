package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/felo/felo-backend/services/driver-service/internal/domain"
	"github.com/felo/felo-backend/services/driver-service/internal/service"
)

func TestDriverService_RegisterDriver_Success(t *testing.T) {
	repo := &driverRepoFake{drivers: map[string]domain.DriverProfile{}}
	svc := service.NewDriverService(repo)

	driver, err := svc.RegisterDriver(context.Background(), "driver-1", "Budi", "123", domain.VehicleInfo{
		LicensePlate: "B 1234 ABC",
		Type:         "Motor",
	})
	if err != nil {
		t.Fatalf("RegisterDriver() error = %v", err)
	}

	if driver.ID != "driver-1" || driver.KYCStatus != domain.KYCPending {
		t.Fatalf("unexpected driver data: %+v", driver)
	}
}

func TestDriverService_RegisterDriver_InvalidInput(t *testing.T) {
	svc := service.NewDriverService(&driverRepoFake{})

	_, err := svc.RegisterDriver(context.Background(), "", "", "", domain.VehicleInfo{})
	if !errors.Is(err, service.ErrInvalidInput) {
		t.Fatalf("RegisterDriver() error = %v, want ErrInvalidInput", err)
	}
}

func TestDriverService_ApproveKYC_Success(t *testing.T) {
	repo := &driverRepoFake{drivers: map[string]domain.DriverProfile{
		"driver-1": {ID: "driver-1", KYCStatus: domain.KYCPending},
	}}
	svc := service.NewDriverService(repo)

	driver, err := svc.ApproveKYC(context.Background(), "driver-1")
	if err != nil {
		t.Fatalf("ApproveKYC() error = %v", err)
	}

	if driver.KYCStatus != domain.KYCApproved {
		t.Fatalf("driver.KYCStatus = %v, want %v", driver.KYCStatus, domain.KYCApproved)
	}
}

func TestDriverService_SetOperationalStatus_Online_RequiresKYCApproved(t *testing.T) {
	repo := &driverRepoFake{drivers: map[string]domain.DriverProfile{
		"driver-1": {ID: "driver-1", KYCStatus: domain.KYCPending},
	}}
	svc := service.NewDriverService(repo)

	_, err := svc.SetOperationalStatus(context.Background(), "driver-1", domain.StatusOnline)
	if !errors.Is(err, service.ErrKYCNotApproved) {
		t.Fatalf("SetOperationalStatus() error = %v, want ErrKYCNotApproved", err)
	}
}

func TestDriverService_SetOperationalStatus_Success(t *testing.T) {
	repo := &driverRepoFake{drivers: map[string]domain.DriverProfile{
		"driver-1": {ID: "driver-1", KYCStatus: domain.KYCApproved, OperationalStatus: domain.StatusOffline},
	}}
	svc := service.NewDriverService(repo)

	driver, err := svc.SetOperationalStatus(context.Background(), "driver-1", domain.StatusOnline)
	if err != nil {
		t.Fatalf("SetOperationalStatus() error = %v", err)
	}

	if driver.OperationalStatus != domain.StatusOnline {
		t.Fatalf("driver.OperationalStatus = %v, want %v", driver.OperationalStatus, domain.StatusOnline)
	}
}

func TestDriverService_SetOperationalStatus_NotFound(t *testing.T) {
	repo := &driverRepoFake{drivers: map[string]domain.DriverProfile{}}
	svc := service.NewDriverService(repo)

	_, err := svc.SetOperationalStatus(context.Background(), "driver-unknown", domain.StatusOnline)
	if !errors.Is(err, service.ErrDriverNotFound) {
		t.Fatalf("SetOperationalStatus() error = %v, want ErrDriverNotFound", err)
	}
}

type driverRepoFake struct {
	drivers map[string]domain.DriverProfile
}

func (r *driverRepoFake) Save(_ context.Context, driver domain.DriverProfile) error {
	r.drivers[driver.ID] = driver
	return nil
}

func (r *driverRepoFake) GetByID(_ context.Context, driverID string) (domain.DriverProfile, error) {
	driver, ok := r.drivers[driverID]
	if !ok {
		return domain.DriverProfile{}, service.ErrDriverNotFound
	}
	return driver, nil
}
