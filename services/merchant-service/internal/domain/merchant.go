package domain

import "time"

// Merchant mewakili entitas restoran/toko
type Merchant struct {
	ID               string
	Name             string
	Latitude         float64 // Esensial untuk kalkulasi jarak > 1 KM di Order Service [cite: 53-54, 262-263]
	Longitude        float64
	OpenTime         string  // Format "HH:MM"
	CloseTime        string  // Format "HH:MM"
	IsManuallyClosed bool    // Toggle Tutup Darurat
	CreatedAt        time.Time
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