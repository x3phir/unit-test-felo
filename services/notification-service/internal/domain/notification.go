package domain

import "time"

type NotificationChannel string

const (
	ChannelPush     NotificationChannel = "Push"
	ChannelWhatsApp NotificationChannel = "WhatsApp"
	ChannelSMS      NotificationChannel = "SMS"
	ChannelEmail    NotificationChannel = "Email"
)

type NotificationRequest struct {
	ReceiverID string
	MsgContent string
	Channel    NotificationChannel
	ActionLink string
}

type DeliveryStatus struct {
	Status  string // Sent/Delivered/Failed
	Channel NotificationChannel
}

type NotificationRecord struct {
	NotificationID string
	RecipientRef   string
	Channel        string
	TemplateCode   string
	Status         string
	CreatedAt      time.Time
}