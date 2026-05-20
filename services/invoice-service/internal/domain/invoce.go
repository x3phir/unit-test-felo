package domain

type InvoiceStatus string

const (
	StatusIssued    InvoiceStatus = "issued"
	StatusPaid      InvoiceStatus = "paid"
	StatusCancelled InvoiceStatus = "cancelled"
)

type Invoice struct {
	InvoiceID  string
	SubjectRef string
	CustomerID string
	Amount     int64
	Currency   string
	Status     InvoiceStatus
}