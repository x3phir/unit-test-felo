package domain

import "time"

type KYCStatus string

const (
	KYCPending  KYCStatus = "pending"
	KYCApproved KYCStatus = "approved"
	KYCRejected KYCStatus = "rejected"
)

type OperationalStatus string

const (
	StatusOnline  OperationalStatus = "online"
	StatusOffline OperationalStatus = "offline"
	StatusBusy    OperationalStatus = "busy"
)

type VehicleInfo struct {
	LicensePlate string
	Type         string // "Motor" or "Mobil"
}

type DriverProfile struct {
	ID                string
	Name              string
	Phone             string
	Vehicle           VehicleInfo
	KYCStatus         KYCStatus
	OperationalStatus OperationalStatus
	Rating            float64
	CreatedAt         time.Time
	UpdatedAt         time.Time
}
