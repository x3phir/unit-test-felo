package service

import (
	"context"
	"errors"
	"time"

	"github.com/felo/felo-backend/services/auth-service/internal/domain"
	"github.com/felo/felo-backend/services/auth-service/internal/ports"
)

var (
	ErrInvalidInput = errors.New("invalid input")
	ErrOTPNotFound  = errors.New("OTP not found")
	ErrInvalidOTP   = errors.New("invalid OTP code")
	ErrOTPExpired   = errors.New("OTP is expired")
)

type AuthService struct {
	otpStore ports.OTPStore
	sessions ports.SessionStore
	tokens   ports.TokenGenerator
	now      func() time.Time
	otpTTL   time.Duration
}

// NewAuthService membuat instance baru dari AuthService dengan dependensi yang diberikan.
func NewAuthService(otpStore ports.OTPStore, sessions ports.SessionStore, tokens ports.TokenGenerator) *AuthService {
	return &AuthService{
		otpStore: otpStore,
		sessions: sessions,
		tokens:   tokens,
		now:      time.Now,
		otpTTL:   5 * time.Minute,
	}
}

// RequestOTP membuat dan menyimpan One-Time Password (OTP) baru untuk nomor telepon yang diberikan.
// OTP akan kedaluwarsa setelah batas waktu Time-To-Live (TTL) yang ditentukan.
func (s *AuthService) RequestOTP(ctx context.Context, phone string) (domain.OTP, error) {
	if phone == "" {
		return domain.OTP{}, ErrInvalidInput
	}

	otp := domain.OTP{
		ID:        "otp-123", // Ideally generated ID
		Phone:     phone,
		Code:      "123456", // Ideally random 6-digit code
		ExpiresAt: s.now().Add(s.otpTTL),
	}

	if err := s.otpStore.Save(ctx, otp); err != nil {
		return domain.OTP{}, err
	}

	return otp, nil
}

// VerifyOTP memvalidasi kode OTP yang diberikan dengan OTP yang tersimpan untuk nomor telepon tersebut.
// Mengembalikan ErrOTPExpired jika OTP telah kedaluwarsa, atau ErrInvalidOTP jika kode salah.
func (s *AuthService) VerifyOTP(ctx context.Context, phone, code string) (domain.VerificationStatus, error) {
	if phone == "" || code == "" {
		return "", ErrInvalidInput
	}

	otp, err := s.otpStore.GetByPhone(ctx, phone)
	if err != nil {
		return "", err
	}

	if s.now().After(otp.ExpiresAt) {
		return domain.StatusExpired, ErrOTPExpired
	}

	if otp.Code != code {
		return domain.StatusInvalid, ErrInvalidOTP
	}

	_ = s.otpStore.Delete(ctx, phone)

	return domain.StatusValid, nil
}

// LoginWithOTP memverifikasi OTP dan, jika berhasil, membuat sesi terautentikasi baru
// dengan token JWT yang dibuat untuk pengguna tersebut.
func (s *AuthService) LoginWithOTP(ctx context.Context, userID string, phone, code string, role domain.UserRole) (domain.AuthSession, error) {
	status, err := s.VerifyOTP(ctx, phone, code)
	if err != nil {
		return domain.AuthSession{}, err
	}
	if status != domain.StatusValid {
		return domain.AuthSession{}, ErrInvalidOTP
	}

	token, err := s.tokens.GenerateToken(userID, role)
	if err != nil {
		return domain.AuthSession{}, err
	}

	session := domain.AuthSession{
		ID:        "session-123",
		UserID:    userID,
		Role:      role,
		Token:     token,
		ExpiresAt: s.now().Add(24 * time.Hour),
	}

	if err := s.sessions.Save(ctx, session); err != nil {
		return domain.AuthSession{}, err
	}

	return session, nil
}
