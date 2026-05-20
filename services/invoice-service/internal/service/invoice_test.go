package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/felo/felo-backend/services/invoice-service/internal/domain"
	"github.com/felo/felo-backend/services/invoice-service/internal/service"
)

func TestCreateInvoice_Success(t *testing.T) {
	repo := &invoiceRepoFake{invoices: make(map[string]*domain.Invoice)}
	pub := &notificationPublisherFake{}
	svc := service.NewInvoiceService(repo, pub, time.Now)

	inv, err := svc.CreateInvoice(context.Background(), "ORD-123", 50000, "USR-001")

	if err != nil {
		t.Fatalf("Ekspektasi tidak ada error, mendapat: %v", err)
	}
	if inv.Status != domain.StatusIssued {
		t.Errorf("Ekspektasi status pending, mendapat %s", inv.Status)
	}
	if inv.Amount != 50000 {
		t.Errorf("Ekspektasi amount 50000, mendapat %d", inv.Amount)
	}
}

func TestUpdateInvoiceStatus(t *testing.T) {
	repo := &invoiceRepoFake{
		invoices: map[string]*domain.Invoice{
			"INV-1": {InvoiceID: "INV-1", Status: domain.StatusIssued},
		},
	}
	svc := service.NewInvoiceService(repo, &notificationPublisherFake{}, time.Now)

	errStatus := svc.UpdateInvoiceStatus(context.Background(), "INV-1", domain.StatusPaid)

	if errStatus != nil {
		t.Fatalf("Gagal melakukan update")
	}

	inv, _ := svc.GetInvoice(context.Background(), "INV-1")
	if inv.Status != domain.StatusPaid {
		t.Errorf("Status gagal diperbarui menjadi paid")
	}
}

func TestSendInvoiceNotification(t *testing.T) {
	repo := &invoiceRepoFake{
		invoices: map[string]*domain.Invoice{
			"INV-1": {InvoiceID: "INV-1", SubjectRef: "ORD-1"},
		},
	}
	pub := &notificationPublisherFake{}
	svc := service.NewInvoiceService(repo, pub, time.Now)

	err := svc.SendInvoiceNotification(context.Background(), "INV-1")
	
	if err != nil {
		t.Fatalf("Gagal mengirim notifikasi")
	}
	if len(pub.publishedInvoices) != 1 {
		t.Errorf("Ekspektasi 1 event terpublikasi, mendapat %d", len(pub.publishedInvoices))
	}
}

// ==========================================
// FAKE IMPLEMENTATIONS
// ==========================================

type invoiceRepoFake struct {
	invoices map[string]*domain.Invoice
}

func (f *invoiceRepoFake) Create(_ context.Context, inv *domain.Invoice) error {
	f.invoices[inv.InvoiceID] = inv
	return nil
}

func (f *invoiceRepoFake) GetByID(_ context.Context, id string) (*domain.Invoice, error) {
	if inv, ok := f.invoices[id]; ok {
		return inv, nil
	}
	return nil, service.ErrInvoiceNotFound
}

func (f *invoiceRepoFake) GetByOrderID(_ context.Context, orderID string) ([]domain.Invoice, error) {
	var result []domain.Invoice
	for _, inv := range f.invoices {
		if inv.SubjectRef == orderID {
			result = append(result, *inv)
		}
	}
	return result, nil
}

func (f *invoiceRepoFake) UpdateStatus(_ context.Context, id string, status domain.InvoiceStatus) error {
	if inv, ok := f.invoices[id]; ok {
		inv.Status = status
		return nil
	}
	return service.ErrInvoiceNotFound
}



type notificationPublisherFake struct {
	publishedInvoices []*domain.Invoice
}

func (f *notificationPublisherFake) PublishInvoiceNotification(_ context.Context, inv *domain.Invoice) error {
	f.publishedInvoices = append(f.publishedInvoices, inv)
	return nil
}