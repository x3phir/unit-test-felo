package domain

import "time"

type UserRole string

const (
	RolePassenger UserRole = "passenger"
	RoleDriver    UserRole = "driver"
	RoleMerchant  UserRole = "merchant"
)

type VerificationStatus string

const (
	StatusValid   VerificationStatus = "valid"
	StatusInvalid VerificationStatus = "invalid"
	StatusExpired VerificationStatus = "expired"
)

type AuthSession struct {
	ID        string
	UserID    string
	Role      UserRole
	Token     string
	ExpiresAt time.Time
}

type OTP struct {
	ID        string
	Phone     string
	Code      string
	ExpiresAt time.Time
}
