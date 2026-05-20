package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/felo/felo-backend/services/invoice-service/internal/domain"
	"github.com/felo/felo-backend/services/invoice-service/internal/ports"
)

var (
	ErrInvoiceNotFound = errors.New("invoice tidak ditemukan")
	ErrInvalidAmount   = errors.New("jumlah tagihan tidak valid")
)

type InvoiceService struct {
	repo      ports.InvoiceRepository
	publisher ports.NotificationPublisher
	now       func() time.Time
}

func NewInvoiceService(repo ports.InvoiceRepository, pub ports.NotificationPublisher, now func() time.Time) *InvoiceService {
	return &InvoiceService{
		repo:      repo,
		publisher: pub,
		now:       now,
	}
}

// CreateInvoice membuat invoice baru dari order dengan status awal pending
func (s *InvoiceService) CreateInvoice(ctx context.Context, orderID string, amount float64, payerID string) (*domain.Invoice, error) {
	if amount < 0 {
		return nil, ErrInvalidAmount
	}

	invoice := &domain.Invoice{
		ID:           fmt.Sprintf("INV-%s-%d", orderID, s.now().Unix()),
		OrderID:      orderID,
		Amount:       amount,
		FinalPayerID: payerID,
		Status:       domain.StatusPending,
		ReceiptURL:   fmt.Sprintf("https://invoice.felo.app/%s", orderID),
		CreatedAt:    s.now(),
		UpdatedAt:    s.now(),
	}

	if err := s.repo.Create(ctx, invoice); err != nil {
		return nil, err
	}

	return invoice, nil
}

// GetInvoice mengambil detail invoice berdasarkan ID spesifik
func (s *InvoiceService) GetInvoice(ctx context.Context, id string) (*domain.Invoice, error) {
	return s.repo.GetByID(ctx, id)
}

// GetInvoicesByOrderID mengambil semua tagihan yang berkaitan dengan satu pesanan
func (s *InvoiceService) GetInvoicesByOrderID(ctx context.Context, orderID string) ([]domain.Invoice, error) {
	return s.repo.GetByOrderID(ctx, orderID)
}

// UpdateInvoiceStatus memperbarui status pembayaran tagihan
func (s *InvoiceService) UpdateInvoiceStatus(ctx context.Context, id string, status domain.InvoiceStatus) error {
	// Opsional: Lakukan validasi transisi status di sini (misal: paid tidak bisa jadi pending)
	return s.repo.UpdateStatus(ctx, id, status)
}

// SetPaymentReference menyimpan nomor referensi dari payment gateway/wallet
func (s *InvoiceService) SetPaymentReference(ctx context.Context, id string, reference string) error {
	return s.repo.UpdatePaymentReference(ctx, id, reference)
}

// SendInvoiceNotification memicu pengiriman notifikasi/nota digital ke user
func (s *InvoiceService) SendInvoiceNotification(ctx context.Context, id string) error {
	invoice, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	return s.publisher.PublishInvoiceNotification(ctx, invoice)
}