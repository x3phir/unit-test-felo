package domain

import "time"

type PayerType string

const (
	PayerSender   PayerType = "sender"
	PayerReceiver PayerType = "receiver"
)

type PackageDetails struct {
	WeightKG    float64
	Dimensions  string
	Description string
}

type SendOrder struct {
	ID             string
	SenderID       string
	ReceiverPhone  string
	PackageDetails PackageDetails
	PayerType      PayerType
	ShippingFee    int64
	Status         string
	CreatedAt      time.Time
}

type Event struct {
	Name       string
	OrderID    string
	OccurredAt time.Time
}
