package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/felo/felo-backend/services/auth-service/internal/domain"
	"github.com/felo/felo-backend/services/auth-service/internal/service"
)

func TestAuthService_RequestOTP_Success(t *testing.T) {
	otpStore := &otpStoreFake{otps: map[string]domain.OTP{}}
	svc := service.NewAuthService(otpStore, &sessionStoreFake{}, &tokenGenFake{})

	otp, err := svc.RequestOTP(context.Background(), "+628123456789")
	if err != nil {
		t.Fatalf("RequestOTP() error = %v", err)
	}

	if otp.Phone != "+628123456789" || otp.Code == "" {
		t.Fatalf("unexpected OTP data: %+v", otp)
	}
}

func TestAuthService_RequestOTP_InvalidInput(t *testing.T) {
	svc := service.NewAuthService(&otpStoreFake{}, &sessionStoreFake{}, &tokenGenFake{})

	_, err := svc.RequestOTP(context.Background(), "")
	if !errors.Is(err, service.ErrInvalidInput) {
		t.Fatalf("RequestOTP() error = %v, want ErrInvalidInput", err)
	}
}

func TestAuthService_VerifyOTP_Success(t *testing.T) {
	now := time.Now()
	otpStore := &otpStoreFake{otps: map[string]domain.OTP{
		"+628123456789": {Phone: "+628123456789", Code: "123456", ExpiresAt: now.Add(5 * time.Minute)},
	}}
	svc := service.NewAuthService(otpStore, &sessionStoreFake{}, &tokenGenFake{})

	status, err := svc.VerifyOTP(context.Background(), "+628123456789", "123456")
	if err != nil {
		t.Fatalf("VerifyOTP() error = %v", err)
	}

	if status != domain.StatusValid {
		t.Fatalf("status = %v, want %v", status, domain.StatusValid)
	}
}

func TestAuthService_VerifyOTP_Expired(t *testing.T) {
	now := time.Now()
	otpStore := &otpStoreFake{otps: map[string]domain.OTP{
		"+628123456789": {Phone: "+628123456789", Code: "123456", ExpiresAt: now.Add(-5 * time.Minute)},
	}}
	svc := service.NewAuthService(otpStore, &sessionStoreFake{}, &tokenGenFake{})

	_, err := svc.VerifyOTP(context.Background(), "+628123456789", "123456")
	if !errors.Is(err, service.ErrOTPExpired) {
		t.Fatalf("VerifyOTP() error = %v, want ErrOTPExpired", err)
	}
}

func TestAuthService_VerifyOTP_InvalidCode(t *testing.T) {
	now := time.Now()
	otpStore := &otpStoreFake{otps: map[string]domain.OTP{
		"+628123456789": {Phone: "+628123456789", Code: "123456", ExpiresAt: now.Add(5 * time.Minute)},
	}}
	svc := service.NewAuthService(otpStore, &sessionStoreFake{}, &tokenGenFake{})

	_, err := svc.VerifyOTP(context.Background(), "+628123456789", "wrong")
	if !errors.Is(err, service.ErrInvalidOTP) {
		t.Fatalf("VerifyOTP() error = %v, want ErrInvalidOTP", err)
	}
}

func TestAuthService_LoginWithOTP_Success(t *testing.T) {
	now := time.Now()
	otpStore := &otpStoreFake{otps: map[string]domain.OTP{
		"+628123456789": {Phone: "+628123456789", Code: "123456", ExpiresAt: now.Add(5 * time.Minute)},
	}}
	sessionStore := &sessionStoreFake{sessions: map[string]domain.AuthSession{}}
	tokenGen := &tokenGenFake{}

	svc := service.NewAuthService(otpStore, sessionStore, tokenGen)

	session, err := svc.LoginWithOTP(context.Background(), "user-1", "+628123456789", "123456", domain.RolePassenger)
	if err != nil {
		t.Fatalf("LoginWithOTP() error = %v", err)
	}

	if session.UserID != "user-1" || session.Token != "token-user-1" {
		t.Fatalf("unexpected session: %+v", session)
	}
}

type otpStoreFake struct {
	otps map[string]domain.OTP
}

func (s *otpStoreFake) Save(_ context.Context, otp domain.OTP) error {
	s.otps[otp.Phone] = otp
	return nil
}

func (s *otpStoreFake) GetByPhone(_ context.Context, phone string) (domain.OTP, error) {
	otp, ok := s.otps[phone]
	if !ok {
		return domain.OTP{}, service.ErrOTPNotFound
	}
	return otp, nil
}

func (s *otpStoreFake) Delete(_ context.Context, phone string) error {
	delete(s.otps, phone)
	return nil
}

type sessionStoreFake struct {
	sessions map[string]domain.AuthSession
}

func (s *sessionStoreFake) Save(_ context.Context, session domain.AuthSession) error {
	s.sessions[session.ID] = session
	return nil
}

type tokenGenFake struct{}

func (g *tokenGenFake) GenerateToken(userID string, _ domain.UserRole) (string, error) {
	return "token-" + userID, nil
}
