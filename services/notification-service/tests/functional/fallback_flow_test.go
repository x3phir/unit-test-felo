//go:build functional

package functional_test

import (
	"context"
	"errors"
	"testing"

	"github.com/felo/felo-backend/services/notification-service/internal/domain"
	"github.com/felo/felo-backend/services/notification-service/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
	"time"
)

func TestNotificationFunctional_FallbackSLA_RoutesCorrectly(t *testing.T) {
	ctx := context.Background()

	db := openPG(t, getenv("FELO_NOTIFICATION_PG_DSN", "postgres://felo:felo@127.0.0.1:54328/notification_db?sslmode=disable"))
	t.Cleanup(func() { db.Close() })

	_, err := db.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS notifications (
			notification_id VARCHAR(255) PRIMARY KEY,
			recipient_ref VARCHAR(255) NOT NULL,
			channel VARCHAR(50) NOT NULL,
			template_code VARCHAR(255) NOT NULL,
			status VARCHAR(50) NOT NULL,
			created_at TIMESTAMPTZ NOT NULL
		);
	`)
	if err != nil {
		t.Fatalf("failed to create tables: %v", err)
	}

	_, _ = db.Exec(ctx, "delete from notifications where recipient_ref=$1", "user-ft-001")

	provider := &dummyNotificationProvider{
		db:       db,
		failPush: true,
		failWA:   false,
		failSMS:  false,
	}

	svc := service.NewNotificationService(provider)

	req := domain.NotificationRequest{
		ReceiverID: "user-ft-001",
		MsgContent: "Pesanan Anda sedang diproses",
		ActionLink: "https://felo.app/orders/ord-ft-001",
	}

	status, err := svc.SendWithFallback(ctx, req)
	if err != nil {
		t.Fatalf("SendWithFallback() unexpected error: %v", err)
	}

	// Should have fallen back to WhatsApp since Push was set to fail
	if status.Channel != domain.ChannelWhatsApp {
		t.Fatalf("Expected fallback to WhatsApp, got %s", status.Channel)
	}
	if status.Status != "Sent" {
		t.Errorf("Expected status Sent, got %v", status.Status)
	}

	// Verify attempt history
	if !provider.pushAttempted {
		t.Errorf("Expected Push to be attempted before fallback")
	}
	if !provider.waAttempted {
		t.Errorf("Expected WhatsApp to be attempted")
	}
	// Verify DB record
	var dbChannel, dbStatus string
	if err := db.QueryRow(ctx, "select channel, status from notifications where recipient_ref=$1 order by created_at desc limit 1", req.ReceiverID).Scan(&dbChannel, &dbStatus); err != nil {
		t.Fatalf("query persisted notification: %v", err)
	}
	
	if dbChannel != "WhatsApp" {
		t.Fatalf("persisted channel = %s, want WhatsApp", dbChannel)
	}
}

type dummyNotificationProvider struct {
	db       *pgxpool.Pool
	failPush bool
	failWA   bool
	failSMS  bool

	pushAttempted bool
	waAttempted   bool
	smsAttempted  bool
}

func (p *dummyNotificationProvider) SendPush(ctx context.Context, req domain.NotificationRequest) (domain.DeliveryStatus, error) {
	p.pushAttempted = true
	status := "Sent"
	var err error
	if p.failPush {
		status = "Failed"
		err = errors.New("push gateway timeout")
	}
	p.saveToDB(ctx, req, "Push", status)
	return domain.DeliveryStatus{Status: status, Channel: domain.ChannelPush}, err
}

func (p *dummyNotificationProvider) SendWhatsApp(ctx context.Context, req domain.NotificationRequest) (domain.DeliveryStatus, error) {
	p.waAttempted = true
	status := "Sent"
	var err error
	if p.failWA {
		status = "Failed"
		err = errors.New("wa provider offline")
	}
	p.saveToDB(ctx, req, "WhatsApp", status)
	return domain.DeliveryStatus{Status: status, Channel: domain.ChannelWhatsApp}, err
}

func (p *dummyNotificationProvider) SendSMS(ctx context.Context, req domain.NotificationRequest) (domain.DeliveryStatus, error) {
	p.smsAttempted = true
	status := "Sent"
	var err error
	if p.failSMS {
		status = "Failed"
		err = errors.New("sms provider offline")
	}
	p.saveToDB(ctx, req, "SMS", status)
	return domain.DeliveryStatus{Status: status, Channel: domain.ChannelSMS}, err
}

func (p *dummyNotificationProvider) saveToDB(ctx context.Context, req domain.NotificationRequest, channel string, status string) {
	_, _ = p.db.Exec(ctx, `insert into notifications (notification_id, recipient_ref, channel, template_code, status, created_at)
values ($1,$2,$3,$4,$5,$6)`,
		"notif-"+channel, req.ReceiverID, channel, "DEFAULT_TPL", status, time.Now().UTC())
}

func openPG(t *testing.T, dsn string) *pgxpool.Pool {
	t.Helper()
	db, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		t.Fatalf("pgxpool.New() error = %v", err)
	}
	return db
}

func getenv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
