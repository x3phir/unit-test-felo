package domain

import "time"

// Merchant mewakili entitas restoran/toko
type Merchant struct {
	MerchantID  string
	OwnerUserID string
	Name        string
	Status      string
	UpdatedAt   time.Time
}

// Menu mewakili item makanan/minuman yang terikat pada Merchant
type Menu struct {
	ID          string
	MerchantID  string
	Name        string
	Description string
	Price       float64
	IsAvailable bool
}