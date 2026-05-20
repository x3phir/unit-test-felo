//go:build functional

package functional_test

import (
	"context"
	"os"
	"testing"

	"github.com/felo/felo-backend/services/driver-service/internal/domain"
	"github.com/felo/felo-backend/services/driver-service/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestDriverFunctional_RegisterAndApproveDriver(t *testing.T) {
	ctx := context.Background()
	db := openPG(t, getenv("FELO_DRIVER_PG_DSN", "postgres://felo:felo@127.0.0.1:54331/driver_db?sslmode=disable"))
	t.Cleanup(func() { db.Close() })

	initDB(t, db)

	repo := &pgDriverRepository{db: db}
	svc := service.NewDriverService(repo)

	driverID := "driver-ft-001"
	
	// 1. Register Driver
	_, err := svc.RegisterDriver(ctx, driverID, "Budi", "+628111", domain.VehicleInfo{
		LicensePlate: "B 1234 ABC",
		Type:         "Motor",
	})
	if err != nil {
		t.Fatalf("RegisterDriver() error = %v", err)
	}

	// Verify Driver in DB
	var dbKyc string
	if err := db.QueryRow(ctx, "SELECT kyc_status FROM drivers WHERE id=$1", driverID).Scan(&dbKyc); err != nil {
		t.Fatalf("query persisted driver: %v", err)
	}
	if dbKyc != string(domain.KYCPending) {
		t.Fatalf("persisted kyc = %s, want %s", dbKyc, domain.KYCPending)
	}

	// Try to Set Online (Should fail because KYC not approved)
	_, err = svc.SetOperationalStatus(ctx, driverID, domain.StatusOnline)
	if err == nil {
		t.Fatalf("expected error setting online before KYC approved")
	}

	// 2. Approve KYC
	_, err = svc.ApproveKYC(ctx, driverID)
	if err != nil {
		t.Fatalf("ApproveKYC() error = %v", err)
	}
	
	// 3. Set Operational Status
	updatedDriver, err := svc.SetOperationalStatus(ctx, driverID, domain.StatusOnline)
	if err != nil {
		t.Fatalf("SetOperationalStatus() error = %v", err)
	}

	if updatedDriver.OperationalStatus != domain.StatusOnline {
		t.Fatalf("expected status %s, got %s", domain.StatusOnline, updatedDriver.OperationalStatus)
	}

	// Verify in DB
	var opStatus string
	if err := db.QueryRow(ctx, "SELECT op_status FROM drivers WHERE id=$1", driverID).Scan(&opStatus); err != nil {
		t.Fatalf("query persisted op_status: %v", err)
	}
	if opStatus != string(domain.StatusOnline) {
		t.Fatalf("expected db op_status %s, got %s", domain.StatusOnline, opStatus)
	}
}

func initDB(t *testing.T, db *pgxpool.Pool) {
	t.Helper()
	ctx := context.Background()
	_, err := db.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS drivers (
			id VARCHAR(50) PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			phone VARCHAR(20) NOT NULL,
			vehicle_license VARCHAR(20),
			vehicle_type VARCHAR(20),
			kyc_status VARCHAR(20),
			op_status VARCHAR(20),
			rating FLOAT,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		);
		TRUNCATE drivers;
	`)
	if err != nil {
		t.Fatalf("initDB() error = %v", err)
	}
}

type pgDriverRepository struct{ db *pgxpool.Pool }

func (r *pgDriverRepository) Save(ctx context.Context, driver domain.DriverProfile) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO drivers (id, name, phone, vehicle_license, vehicle_type, kyc_status, op_status, rating, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			phone = EXCLUDED.phone,
			vehicle_license = EXCLUDED.vehicle_license,
			vehicle_type = EXCLUDED.vehicle_type,
			kyc_status = EXCLUDED.kyc_status,
			op_status = EXCLUDED.op_status,
			rating = EXCLUDED.rating,
			updated_at = EXCLUDED.updated_at
	`, driver.ID, driver.Name, driver.Phone, driver.Vehicle.LicensePlate, driver.Vehicle.Type, string(driver.KYCStatus), string(driver.OperationalStatus), driver.Rating, driver.CreatedAt, driver.UpdatedAt)
	return err
}

func (r *pgDriverRepository) GetByID(ctx context.Context, driverID string) (domain.DriverProfile, error) {
	var driver domain.DriverProfile
	var kyc, op string
	err := r.db.QueryRow(ctx, "SELECT id, name, phone, vehicle_license, vehicle_type, kyc_status, op_status, rating, created_at, updated_at FROM drivers WHERE id=$1", driverID).
		Scan(&driver.ID, &driver.Name, &driver.Phone, &driver.Vehicle.LicensePlate, &driver.Vehicle.Type, &kyc, &op, &driver.Rating, &driver.CreatedAt, &driver.UpdatedAt)
	if err != nil {
		return domain.DriverProfile{}, err
	}
	driver.KYCStatus = domain.KYCStatus(kyc)
	driver.OperationalStatus = domain.OperationalStatus(op)
	return driver, nil
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
