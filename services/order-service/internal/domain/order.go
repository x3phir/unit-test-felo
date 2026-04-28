package domain

import "time"

type OrderState string

const (
	StatePending   OrderState = "pending"
	StateConfirmed OrderState = "confirmed"
	StateCooking   OrderState = "cooking"
	StateDelivering OrderState = "delivering"
	StateCompleted OrderState = "completed"
	StateCancelled OrderState = "cancelled"
)

type PaymentMethod string

const (
	PayCash   PaymentMethod = "cash"
	PayWallet PaymentMethod = "wallet"
)

type FoodOrder struct {
	ID             string
	UserID         string
	MerchantID     string
	Items          []string
	PaymentMethod  PaymentMethod
	Status         OrderState
	DistanceKM     float64
	TotalAmount    int64
	OTPTriggered   bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type Event struct {
	Name       string
	OrderID    string
	OccurredAt time.Time
}
