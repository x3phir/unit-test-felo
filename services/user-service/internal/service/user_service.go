package service

import (
	"context"
	"errors"
	"time"

	"github.com/felo/felo-backend/services/user-service/internal/domain"
	"github.com/felo/felo-backend/services/user-service/internal/ports"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrInvalidInput = errors.New("invalid user input")
)

type UpdateProfileInput struct {
	Name     string
	Email    string
	PhotoURL string
	Locale   string
}

type UserService struct {
	repo ports.UserRepository
	now  func() time.Time
}

// NewUserService membuat instance baru dari UserService dengan repository yang diberikan.
func NewUserService(repo ports.UserRepository) *UserService {
	return &UserService{
		repo: repo,
		now:  time.Now,
	}
}

// CreateUser mendaftarkan profil pengguna baru dengan ID, nama, dan nomor telepon yang diberikan.
// Mengembalikan ErrInvalidInput jika ada field wajib yang kosong.
func (s *UserService) CreateUser(ctx context.Context, id, name, phone string) (domain.UserProfile, error) {
	if id == "" || name == "" || phone == "" {
		return domain.UserProfile{}, ErrInvalidInput
	}

	user := domain.UserProfile{
		ID:        id,
		Name:      name,
		Phone:     phone,
		Locale:    "id_ID", // Default to Indonesian per PRD
		CreatedAt: s.now(),
		UpdatedAt: s.now(),
	}

	if err := s.repo.Save(ctx, user); err != nil {
		return domain.UserProfile{}, err
	}

	return user, nil
}

// UpdateProfile memodifikasi data profil pengguna yang sudah ada berdasarkan input yang diberikan.
// Hanya field yang tidak kosong pada input yang akan diperbarui.
func (s *UserService) UpdateProfile(ctx context.Context, id string, input UpdateProfileInput) (domain.UserProfile, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return domain.UserProfile{}, err
	}

	if input.Name != "" {
		user.Name = input.Name
	}
	if input.Email != "" {
		user.Email = input.Email
	}
	if input.PhotoURL != "" {
		user.PhotoURL = input.PhotoURL
	}
	if input.Locale != "" {
		user.Locale = input.Locale
	}

	user.UpdatedAt = s.now()

	if err := s.repo.Save(ctx, user); err != nil {
		return domain.UserProfile{}, err
	}

	return user, nil
}

// AddSavedAddress menambahkan alamat tersimpan baru ke profil pengguna untuk akses cepat.
// Mengembalikan ErrInvalidInput jika nama alamat atau koordinat tidak ada.
func (s *UserService) AddSavedAddress(ctx context.Context, id string, address domain.SavedAddress) (domain.UserProfile, error) {
	if address.Name == "" || address.Coordinate.Latitude == 0 || address.Coordinate.Longitude == 0 {
		return domain.UserProfile{}, ErrInvalidInput
	}

	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return domain.UserProfile{}, err
	}

	user.SavedAddresses = append(user.SavedAddresses, address)
	user.UpdatedAt = s.now()

	if err := s.repo.Save(ctx, user); err != nil {
		return domain.UserProfile{}, err
	}

	return user, nil
}
