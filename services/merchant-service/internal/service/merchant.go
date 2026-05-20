package service

import (
	"context"
	"errors"
	"time"

	"github.com/felo/felo-backend/services/merchant-service/internal/domain"
	"github.com/felo/felo-backend/services/merchant-service/internal/ports"
)

var (
	ErrMerchantNotFound = errors.New("merchant tidak ditemukan")
	ErrMenuNotFound     = errors.New("menu tidak ditemukan")
	ErrMenuUnavailable  = errors.New("menu sedang tidak tersedia")
	ErrMenuOwnership    = errors.New("menu tidak berasal dari merchant yang sesuai")
)

type MerchantService struct {
	merchantRepo ports.MerchantRepository
	menuRepo     ports.MenuRepository
	now          func() time.Time
}

func NewMerchantService(mRepo ports.MerchantRepository, mnRepo ports.MenuRepository, now func() time.Time) *MerchantService {
	return &MerchantService{
		merchantRepo: mRepo,
		menuRepo:     mnRepo,
		now:          now,
	}
}

// ==========================================
// OPERASI MERCHANT
// ==========================================

func (s *MerchantService) GetMerchant(ctx context.Context, id string) (*domain.Merchant, error) {
	return s.merchantRepo.GetByID(ctx, id)
}

func (s *MerchantService) CreateMerchant(ctx context.Context, merchant *domain.Merchant) error {
	merchant.UpdatedAt = s.now()
	return s.merchantRepo.Create(ctx, merchant)
}

func (s *MerchantService) SetMerchantStatus(ctx context.Context, id string, isClosed bool) error {
	return s.merchantRepo.UpdateStatus(ctx, id, isClosed)
}

func (s *MerchantService) IsMerchantOpen(ctx context.Context, id string) (bool, error) {
	merchant, err := s.merchantRepo.GetByID(ctx, id)
	if err != nil {
		return false, err
	}

	return merchant.Status == "active", nil
}

// ==========================================
// OPERASI MENU
// ==========================================

func (s *MerchantService) GetMenus(ctx context.Context, merchantID string) ([]domain.Menu, error) {
	return s.menuRepo.GetByMerchantID(ctx, merchantID)
}

func (s *MerchantService) GetMenu(ctx context.Context, menuID string) (*domain.Menu, error) {
	return s.menuRepo.GetByID(ctx, menuID)
}

func (s *MerchantService) CreateMenu(ctx context.Context, menu *domain.Menu) error {
	return s.menuRepo.Create(ctx, menu)
}

func (s *MerchantService) UpdateMenu(ctx context.Context, menu *domain.Menu) error {
	return s.menuRepo.Update(ctx, menu)
}

func (s *MerchantService) SetMenuAvailability(ctx context.Context, menuID string, isAvailable bool) error {
	return s.menuRepo.UpdateAvailability(ctx, menuID, isAvailable)
}

// ==========================================
// OPERASI KOMBINASI
// ==========================================

func (s *MerchantService) ValidateMenu(ctx context.Context, merchantID string, menuID string) error {
	menu, err := s.menuRepo.GetByID(ctx, menuID)
	if err != nil {
		return err
	}
	
	if menu.MerchantID != merchantID {
		return ErrMenuOwnership
	}
	if !menu.IsAvailable {
		return ErrMenuUnavailable
	}
	
	return nil
}

func (s *MerchantService) GetMenuPrice(ctx context.Context, menuID string) (float64, error) {
	menu, err := s.menuRepo.GetByID(ctx, menuID)
	if err != nil {
		return 0, err
	}
	return menu.Price, nil
}

func (s *MerchantService) GetMenusByIDs(ctx context.Context, menuIDs []string) ([]domain.Menu, error) {
	if len(menuIDs) == 0 {
		return []domain.Menu{}, nil
	}
	return s.menuRepo.GetByIDs(ctx, menuIDs)
}