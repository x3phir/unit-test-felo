package service

import (
	"context"
	"errors"

	"github.com/felo/felo-backend/services/notification-service/internal/domain"
	"github.com/felo/felo-backend/services/notification-service/internal/ports"
)

var (
	ErrUnsupportedChannel = errors.New("saluran notifikasi tidak didukung")
	ErrAllChannelsFailed  = errors.New("semua saluran fallback gagal")
)

type NotificationService struct {
	provider ports.NotificationProvider
}

func NewNotificationService(provider ports.NotificationProvider) *NotificationService {
	return &NotificationService{
		provider: provider,
	}
}

// SendPush mengirim notifikasi secara spesifik via Push (FCM/APNs)
func (s *NotificationService) SendPush(ctx context.Context, req domain.NotificationRequest) (domain.DeliveryStatus, error) {
	req.Channel = domain.ChannelPush
	return s.provider.SendPush(ctx, req)
}

// SendWhatsApp mengirim notifikasi secara spesifik via WhatsApp Business API
func (s *NotificationService) SendWhatsApp(ctx context.Context, req domain.NotificationRequest) (domain.DeliveryStatus, error) {
	req.Channel = domain.ChannelWhatsApp
	return s.provider.SendWhatsApp(ctx, req)
}

// SendSMS mengirim notifikasi secara spesifik via SMS Gateway
func (s *NotificationService) SendSMS(ctx context.Context, req domain.NotificationRequest) (domain.DeliveryStatus, error) {
	req.Channel = domain.ChannelSMS
	return s.provider.SendSMS(ctx, req)
}

// SendNotification merutekan permintaan berdasarkan properti Channel
func (s *NotificationService) SendNotification(ctx context.Context, req domain.NotificationRequest) (domain.DeliveryStatus, error) {
	switch req.Channel {
	case domain.ChannelPush:
		return s.SendPush(ctx, req)
	case domain.ChannelWhatsApp:
		return s.SendWhatsApp(ctx, req)
	case domain.ChannelSMS:
		return s.SendSMS(ctx, req)
	default:
		return domain.DeliveryStatus{}, ErrUnsupportedChannel
	}
}

// SendWithFallback mengeksekusi strategi SLA: Push -> WA -> SMS
func (s *NotificationService) SendWithFallback(ctx context.Context, req domain.NotificationRequest) (domain.DeliveryStatus, error) {
	// Percobaan 1: Push Notification
	status, err := s.SendPush(ctx, req)
	if err == nil {
		return status, nil
	}

	// Percobaan 2: Fallback ke WhatsApp
	statusWA, errWA := s.SendWhatsApp(ctx, req)
	if errWA == nil {
		return statusWA, nil
	}

	// Percobaan 3: Fallback terakhir ke SMS
	statusSMS, errSMS := s.SendSMS(ctx, req)
	if errSMS == nil {
		return statusSMS, nil
	}

	return domain.DeliveryStatus{Status: "Failed"}, ErrAllChannelsFailed
}