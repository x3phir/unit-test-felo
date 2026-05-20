//go:build functional

package functional_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/felo/felo-backend/services/merchant-service/internal/domain"
	"github.com/felo/felo-backend/services/merchant-service/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestMerchantFunctional_CreateMerchant_PersistsToDatabase(t *testing.T) {
	ctx := context.Background()
	db := openPG(t, getenv("FELO_MERCHANT_PG_DSN", "postgres://felo:felo@127.0.0.1:54327/merchant_db?sslmode=disable"))
	t.Cleanup(func() { db.Close() })

	// Setup tables
	_, err := db.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS merchants (
			merchant_id VARCHAR(255) PRIMARY KEY,
			owner_user_id VARCHAR(255) NOT NULL,
			name VARCHAR(255) NOT NULL,
			status VARCHAR(50) NOT NULL,
			updated_at TIMESTAMPTZ NOT NULL
		);
	`)
	if err != nil {
		t.Fatalf("failed to create tables: %v", err)
	}

	merchantID := "merch-ft-001"
	_, _ = db.Exec(ctx, "delete from merchants where merchant_id=$1", merchantID)

	mRepo := &pgMerchantRepo{db: db}
	mnRepo := &pgMenuRepo{db: db}
	
	now := time.Now().UTC()
	clock := func() time.Time { return now }

	svc := service.NewMerchantService(mRepo, mnRepo, clock)

	merch := &domain.Merchant{
		MerchantID: merchantID,
		OwnerUserID: "user-123",
		Name: "Warung Functional",
		Status: "active",
	}

	err = svc.CreateMerchant(ctx, merch)
	if err != nil {
		t.Fatalf("CreateMerchant() error = %v", err)
	}

	var name string
	var status string
	if err := db.QueryRow(ctx, "select name, status from merchants where merchant_id=$1", merchantID).Scan(&name, &status); err != nil {
		t.Fatalf("query persisted merchant: %v", err)
	}
	
	if name != "Warung Functional" {
		t.Fatalf("persisted name = %s, want %s", name, "Warung Functional")
	}
	if status != "active" {
		t.Fatalf("persisted status = %s, want active", status)
	}
}

type pgMerchantRepo struct{ db *pgxpool.Pool }

func (r *pgMerchantRepo) Create(ctx context.Context, m *domain.Merchant) error {
	_, err := r.db.Exec(ctx, `insert into merchants (merchant_id, owner_user_id, name, status, updated_at)
values ($1,$2,$3,$4,$5)`,
		m.MerchantID, m.OwnerUserID, m.Name, m.Status, m.UpdatedAt)
	return err
}
func (r *pgMerchantRepo) GetByID(ctx context.Context, id string) (*domain.Merchant, error) { return nil, nil }
func (r *pgMerchantRepo) UpdateStatus(ctx context.Context, id string, isClosed bool) error { return nil }

type pgMenuRepo struct{ db *pgxpool.Pool }

func (r *pgMenuRepo) Create(ctx context.Context, m *domain.Menu) error {
	_, err := r.db.Exec(ctx, `insert into menus (id, merchant_id, name, description, price, is_available)
values ($1,$2,$3,$4,$5,$6)`,
		m.ID, m.MerchantID, m.Name, m.Description, m.Price, m.IsAvailable)
	return err
}
func (r *pgMenuRepo) GetByID(ctx context.Context, id string) (*domain.Menu, error) { return nil, nil }
func (r *pgMenuRepo) GetByMerchantID(ctx context.Context, merchantID string) ([]domain.Menu, error) { return nil, nil }
func (r *pgMenuRepo) GetByIDs(ctx context.Context, ids []string) ([]domain.Menu, error) { return nil, nil }
func (r *pgMenuRepo) Update(ctx context.Context, m *domain.Menu) error { return nil }
func (r *pgMenuRepo) UpdateAvailability(ctx context.Context, id string, isAvailable bool) error { return nil }

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
