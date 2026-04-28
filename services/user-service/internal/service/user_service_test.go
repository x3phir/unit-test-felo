package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/felo/felo-backend/services/user-service/internal/domain"
	"github.com/felo/felo-backend/services/user-service/internal/service"
)

func TestUserService_CreateUser_Success(t *testing.T) {
	repo := &userRepoFake{users: map[string]domain.UserProfile{}}
	svc := service.NewUserService(repo)

	user, err := svc.CreateUser(context.Background(), "user-1", "Anas", "+628123456789")
	if err != nil {
		t.Fatalf("CreateUser() error = %v", err)
	}

	if user.ID != "user-1" || user.Name != "Anas" || user.Phone != "+628123456789" {
		t.Fatalf("unexpected user data: %+v", user)
	}
}

func TestUserService_CreateUser_InvalidInput(t *testing.T) {
	svc := service.NewUserService(&userRepoFake{})

	_, err := svc.CreateUser(context.Background(), "", "Anas", "")
	if !errors.Is(err, service.ErrInvalidInput) {
		t.Fatalf("CreateUser() error = %v, want ErrInvalidInput", err)
	}
}

func TestUserService_UpdateProfile_Success(t *testing.T) {
	repo := &userRepoFake{users: map[string]domain.UserProfile{
		"user-1": {ID: "user-1", Name: "Anas", Phone: "123"},
	}}
	svc := service.NewUserService(repo)

	user, err := svc.UpdateProfile(context.Background(), "user-1", service.UpdateProfileInput{
		Email: "anas@felo.com",
	})
	if err != nil {
		t.Fatalf("UpdateProfile() error = %v", err)
	}

	if user.Email != "anas@felo.com" || user.Name != "Anas" {
		t.Fatalf("unexpected user data: %+v", user)
	}
}

func TestUserService_UpdateProfile_NotFound(t *testing.T) {
	repo := &userRepoFake{users: map[string]domain.UserProfile{}}
	svc := service.NewUserService(repo)

	_, err := svc.UpdateProfile(context.Background(), "user-unknown", service.UpdateProfileInput{})
	if !errors.Is(err, service.ErrUserNotFound) {
		t.Fatalf("UpdateProfile() error = %v, want ErrUserNotFound", err)
	}
}

func TestUserService_AddSavedAddress_Success(t *testing.T) {
	repo := &userRepoFake{users: map[string]domain.UserProfile{
		"user-1": {ID: "user-1", Name: "Anas"},
	}}
	svc := service.NewUserService(repo)

	user, err := svc.AddSavedAddress(context.Background(), "user-1", domain.SavedAddress{
		Name:       "Home",
		Coordinate: domain.Coordinate{Latitude: -6.2, Longitude: 106.8},
	})
	if err != nil {
		t.Fatalf("AddSavedAddress() error = %v", err)
	}

	if len(user.SavedAddresses) != 1 || user.SavedAddresses[0].Name != "Home" {
		t.Fatalf("unexpected saved addresses: %+v", user.SavedAddresses)
	}
}

func TestUserService_AddSavedAddress_InvalidInput(t *testing.T) {
	svc := service.NewUserService(&userRepoFake{})

	_, err := svc.AddSavedAddress(context.Background(), "user-1", domain.SavedAddress{})
	if !errors.Is(err, service.ErrInvalidInput) {
		t.Fatalf("AddSavedAddress() error = %v, want ErrInvalidInput", err)
	}
}

type userRepoFake struct {
	users map[string]domain.UserProfile
}

func (r *userRepoFake) Save(_ context.Context, user domain.UserProfile) error {
	r.users[user.ID] = user
	return nil
}

func (r *userRepoFake) GetByID(_ context.Context, userID string) (domain.UserProfile, error) {
	user, ok := r.users[userID]
	if !ok {
		return domain.UserProfile{}, service.ErrUserNotFound
	}
	return user, nil
}
