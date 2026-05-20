//go:build functional

package functional_test

import (
	"context"
	"os"
	"testing"

	"github.com/felo/felo-backend/services/user-service/internal/domain"
	"github.com/felo/felo-backend/services/user-service/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestUserFunctional_CreateAndUpdateUser(t *testing.T) {
	ctx := context.Background()
	db := openPG(t, getenv("FELO_USER_PG_DSN", "postgres://felo:felo@127.0.0.1:54330/user_db?sslmode=disable"))
	t.Cleanup(func() { db.Close() })

	initDB(t, db)

	repo := &pgUserRepository{db: db}
	svc := service.NewUserService(repo)

	userID := "user-ft-001"
	
	// 1. Create User
	_, err := svc.CreateUser(ctx, userID, "Anas", "+628123")
	if err != nil {
		t.Fatalf("CreateUser() error = %v", err)
	}

	// Verify User in DB
	var dbName string
	if err := db.QueryRow(ctx, "SELECT name FROM users WHERE id=$1", userID).Scan(&dbName); err != nil {
		t.Fatalf("query persisted user: %v", err)
	}
	if dbName != "Anas" {
		t.Fatalf("persisted name = %s, want Anas", dbName)
	}

	// 2. Update Profile
	updatedUser, err := svc.UpdateProfile(ctx, userID, service.UpdateProfileInput{
		Email: "anas@test.com",
	})
	if err != nil {
		t.Fatalf("UpdateProfile() error = %v", err)
	}
	
	if updatedUser.Email != "anas@test.com" {
		t.Fatalf("expected updated email anas@test.com, got %s", updatedUser.Email)
	}

	// 3. Add Saved Address
	addr := domain.SavedAddress{
		Name: "Home",
		Coordinate: domain.Coordinate{Latitude: -6.2, Longitude: 106.8},
	}
	userWithAddr, err := svc.AddSavedAddress(ctx, userID, addr)
	if err != nil {
		t.Fatalf("AddSavedAddress() error = %v", err)
	}

	if len(userWithAddr.SavedAddresses) != 1 {
		t.Fatalf("expected 1 saved address, got %d", len(userWithAddr.SavedAddresses))
	}

	// Verify addresses in DB
	var addrName string
	if err := db.QueryRow(ctx, "SELECT name FROM saved_addresses WHERE user_id=$1 LIMIT 1", userID).Scan(&addrName); err != nil {
		t.Fatalf("query persisted address: %v", err)
	}
	if addrName != "Home" {
		t.Fatalf("expected address name Home, got %s", addrName)
	}
}

func initDB(t *testing.T, db *pgxpool.Pool) {
	t.Helper()
	ctx := context.Background()
	_, err := db.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS users (
			id VARCHAR(50) PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			phone VARCHAR(20) NOT NULL,
			email VARCHAR(100),
			photo_url TEXT,
			locale VARCHAR(10),
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		);
		CREATE TABLE IF NOT EXISTS saved_addresses (
			user_id VARCHAR(50) NOT NULL,
			name VARCHAR(100) NOT NULL,
			latitude FLOAT NOT NULL,
			longitude FLOAT NOT NULL
		);
		TRUNCATE users, saved_addresses;
	`)
	if err != nil {
		t.Fatalf("initDB() error = %v", err)
	}
}

type pgUserRepository struct{ db *pgxpool.Pool }

func (r *pgUserRepository) Save(ctx context.Context, user domain.UserProfile) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `
		INSERT INTO users (id, name, phone, email, photo_url, locale, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			phone = EXCLUDED.phone,
			email = EXCLUDED.email,
			photo_url = EXCLUDED.photo_url,
			locale = EXCLUDED.locale,
			updated_at = EXCLUDED.updated_at
	`, user.ID, user.Name, user.Phone, user.Email, user.PhotoURL, user.Locale, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, "DELETE FROM saved_addresses WHERE user_id=$1", user.ID)
	if err != nil {
		return err
	}

	for _, addr := range user.SavedAddresses {
		_, err = tx.Exec(ctx, `
			INSERT INTO saved_addresses (user_id, name, latitude, longitude)
			VALUES ($1, $2, $3, $4)
		`, user.ID, addr.Name, addr.Coordinate.Latitude, addr.Coordinate.Longitude)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *pgUserRepository) GetByID(ctx context.Context, userID string) (domain.UserProfile, error) {
	var user domain.UserProfile
	err := r.db.QueryRow(ctx, "SELECT id, name, phone, email, photo_url, locale, created_at, updated_at FROM users WHERE id=$1", userID).
		Scan(&user.ID, &user.Name, &user.Phone, &user.Email, &user.PhotoURL, &user.Locale, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return domain.UserProfile{}, err
	}

	rows, err := r.db.Query(ctx, "SELECT name, latitude, longitude FROM saved_addresses WHERE user_id=$1", userID)
	if err != nil {
		return domain.UserProfile{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var addr domain.SavedAddress
		if err := rows.Scan(&addr.Name, &addr.Coordinate.Latitude, &addr.Coordinate.Longitude); err != nil {
			return domain.UserProfile{}, err
		}
		user.SavedAddresses = append(user.SavedAddresses, addr)
	}

	return user, nil
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
