package domain

import "time"

type InvoiceStatus string

const (
	StatusPending InvoiceStatus = "pending"
	StatusPaid    InvoiceStatus = "paid"
	StatusFailed  InvoiceStatus = "failed"
)

type Invoice struct {
	ID               string
	OrderID          string
	Amount           float64
	FinalPayerID     string        // Menyimpan Payer_Type (Sender/Receiver khusus FELO-Send) [cite: 139]
	Status           InvoiceStatus // pending, paid, failed
	PaymentReference string        // Menyimpan ID referensi dari Payment Gateway
	ReceiptURL       string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}