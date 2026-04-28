package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/felo/felo-backend/services/notification-service/internal/domain"
	"github.com/felo/felo-backend/services/notification-service/internal/service"
)

func TestSendNotification_SpecificRouting(t *testing.T) {
	provider := &notificationProviderFake{}
	svc := service.NewNotificationService(provider)

	req := domain.NotificationRequest{ReceiverID: "USR-1", Channel: domain.ChannelWhatsApp}
	status, err := svc.SendNotification(context.Background(), req)

	if err != nil {
		t.Fatalf("Gagal mengirim WhatsApp: %v", err)
	}
	if status.Channel != domain.ChannelWhatsApp {
		t.Errorf("Ekspektasi channel WhatsApp, mendapat %s", status.Channel)
	}
	if provider.waCount != 1 {
		t.Errorf("Ekspektasi WhatsApp dipanggil 1 kali, mendapat %d", provider.waCount)
	}
}

func TestSendWithFallback_SuccessOnFirstTry(t *testing.T) {
	provider := &notificationProviderFake{}
	svc := service.NewNotificationService(provider)

	status, err := svc.SendWithFallback(context.Background(), domain.NotificationRequest{ReceiverID: "USR-1"})

	if err != nil {
		t.Fatalf("Fallback gagal: %v", err)
	}
	if status.Channel != domain.ChannelPush {
		t.Errorf("Ekspektasi berhenti di Push, mendapat %s", status.Channel)
	}
	if provider.pushCount != 1 || provider.waCount != 0 {
		t.Errorf("Sistem tidak seharusnya mencoba WA jika Push berhasil")
	}
}

func TestSendWithFallback_DegradesToWhatsApp(t *testing.T) {
	provider := &notificationProviderFake{failPush: true} // Simulasikan Push gagal
	svc := service.NewNotificationService(provider)

	status, err := svc.SendWithFallback(context.Background(), domain.NotificationRequest{ReceiverID: "USR-2"})

	if err != nil {
		t.Fatalf("Fallback gagal menembus WA: %v", err)
	}
	if status.Channel != domain.ChannelWhatsApp {
		t.Errorf("Ekspektasi berlanjut ke WhatsApp, mendapat %s", status.Channel)
	}
	if provider.pushCount != 1 || provider.waCount != 1 || provider.smsCount != 0 {
		t.Errorf("Urutan pemanggilan provider tidak sesuai")
	}
}

func TestSendWithFallback_FailsAllChannels(t *testing.T) {
	// Simulasikan semua down
	provider := &notificationProviderFake{failPush: true, failWA: true, failSMS: true}
	svc := service.NewNotificationService(provider)

	_, err := svc.SendWithFallback(context.Background(), domain.NotificationRequest{ReceiverID: "USR-3"})

	if !errors.Is(err, service.ErrAllChannelsFailed) {
		t.Errorf("Ekspektasi ErrAllChannelsFailed, mendapat %v", err)
	}
	if provider.pushCount != 1 || provider.waCount != 1 || provider.smsCount != 1 {
		t.Errorf("Sistem tidak mencoba seluruh armada channel")
	}
}

// ==========================================
// FAKE IMPLEMENTATIONS
// ==========================================

type notificationProviderFake struct {
	failPush, failWA, failSMS bool
	pushCount, waCount, smsCount int
}

func (f *notificationProviderFake) SendPush(_ context.Context, _ domain.NotificationRequest) (domain.DeliveryStatus, error) {
	f.pushCount++
	if f.failPush {
		return domain.DeliveryStatus{Status: "Failed", Channel: domain.ChannelPush}, errors.New("fcm timeout")
	}
	return domain.DeliveryStatus{Status: "Delivered", Channel: domain.ChannelPush}, nil
}

func (f *notificationProviderFake) SendWhatsApp(_ context.Context, _ domain.NotificationRequest) (domain.DeliveryStatus, error) {
	f.waCount++
	if f.failWA {
		return domain.DeliveryStatus{Status: "Failed", Channel: domain.ChannelWhatsApp}, errors.New("wa api offline")
	}
	return domain.DeliveryStatus{Status: "Sent", Channel: domain.ChannelWhatsApp}, nil
}

func (f *notificationProviderFake) SendSMS(_ context.Context, _ domain.NotificationRequest) (domain.DeliveryStatus, error) {
	f.smsCount++
	if f.failSMS {
		return domain.DeliveryStatus{Status: "Failed", Channel: domain.ChannelSMS}, errors.New("sms provider error")
	}
	return domain.DeliveryStatus{Status: "Sent", Channel: domain.ChannelSMS}, nil
}