package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/felo/felo-backend/services/merchant-service/internal/domain"
	"github.com/felo/felo-backend/services/merchant-service/internal/service"
)

// ==========================================
// TES OPERASI MERCHANT
// ==========================================

func TestIsMerchantOpen(t *testing.T) {
	mRepo := &merchantRepoFake{
		merchants: map[string]*domain.Merchant{
			"MCH-1": {ID: "MCH-1", OpenTime: "08:00", CloseTime: "22:00", IsManuallyClosed: false},
			"MCH-2": {ID: "MCH-2", OpenTime: "08:00", CloseTime: "22:00", IsManuallyClosed: true},
		},
	}
	mnRepo := &menuRepoFake{}
	
	mockTime := func() time.Time {
		t, _ := time.Parse("15:04", "12:00") // Jam 12 siang
		return t
	}
	svc := service.NewMerchantService(mRepo, mnRepo, mockTime)

	// Kasus 1: Buka normal
	isOpen, _ := svc.IsMerchantOpen(context.Background(), "MCH-1")
	if !isOpen {
		t.Errorf("MCH-1 harusnya buka pada jam 12:00")
	}

	// Kasus 2: Tutup darurat
	isOpen2, _ := svc.IsMerchantOpen(context.Background(), "MCH-2")
	if isOpen2 {
		t.Errorf("MCH-2 harusnya tutup karena IsManuallyClosed = true")
	}
}

// ==========================================
// TES OPERASI KOMBINASI
// ==========================================

func TestValidateMenu(t *testing.T) {
	mRepo := &merchantRepoFake{}
	mnRepo := &menuRepoFake{
		menus: map[string]*domain.Menu{
			"MNU-1": {ID: "MNU-1", MerchantID: "MCH-1", IsAvailable: true},
			"MNU-2": {ID: "MNU-2", MerchantID: "MCH-1", IsAvailable: false},
		},
	}
	svc := service.NewMerchantService(mRepo, mnRepo, time.Now)

	// Kasus 1: Valid
	err := svc.ValidateMenu(context.Background(), "MCH-1", "MNU-1")
	if err != nil {
		t.Errorf("Menu 1 harusnya valid, error: %v", err)
	}

	// Kasus 2: Tidak tersedia (Habis)
	err = svc.ValidateMenu(context.Background(), "MCH-1", "MNU-2")
	if !errors.Is(err, service.ErrMenuUnavailable) {
		t.Errorf("Harusnya mendapat ErrMenuUnavailable")
	}

	// Kasus 3: Beda Merchant (Mencegah keranjang belanja campur aduk)
	err = svc.ValidateMenu(context.Background(), "MCH-999", "MNU-1")
	if !errors.Is(err, service.ErrMenuOwnership) {
		t.Errorf("Harusnya mendapat ErrMenuOwnership")
	}
}

func TestGetMenuPrice(t *testing.T) {
	mnRepo := &menuRepoFake{
		menus: map[string]*domain.Menu{
			"MNU-1": {ID: "MNU-1", Price: 25000},
		},
	}
	svc := service.NewMerchantService(&merchantRepoFake{}, mnRepo, time.Now)

	price, err := svc.GetMenuPrice(context.Background(), "MNU-1")
	if err != nil || price != 25000 {
		t.Errorf("Ekspektasi harga 25000, mendapat %v", price)
	}
}

// ==========================================
// FAKE IMPLEMENTATIONS
// ==========================================

type merchantRepoFake struct {
	merchants map[string]*domain.Merchant
}

func (f *merchantRepoFake) GetByID(_ context.Context, id string) (*domain.Merchant, error) {
	if m, ok := f.merchants[id]; ok {
		return m, nil
	}
	return nil, service.ErrMerchantNotFound
}
func (f *merchantRepoFake) Create(_ context.Context, merchant *domain.Merchant) error { return nil }
func (f *merchantRepoFake) UpdateStatus(_ context.Context, id string, isClosed bool) error { return nil }


type menuRepoFake struct {
	menus map[string]*domain.Menu
}

func (f *menuRepoFake) GetByMerchantID(_ context.Context, merchantID string) ([]domain.Menu, error) { return nil, nil }
func (f *menuRepoFake) GetByID(_ context.Context, id string) (*domain.Menu, error) {
	if m, ok := f.menus[id]; ok {
		return m, nil
	}
	return nil, service.ErrMenuNotFound
}
func (f *menuRepoFake) GetByIDs(_ context.Context, ids []string) ([]domain.Menu, error) { return nil, nil }
func (f *menuRepoFake) Create(_ context.Context, menu *domain.Menu) error { return nil }
func (f *menuRepoFake) Update(_ context.Context, menu *domain.Menu) error { return nil }
func (f *menuRepoFake) UpdateAvailability(_ context.Context, id string, isAvailable bool) error { return nil }