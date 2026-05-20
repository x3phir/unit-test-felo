package domain

import "time"

type NotificationChannel string

const (
	ChannelPush     NotificationChannel = "Push"
	ChannelWhatsApp NotificationChannel = "WhatsApp"
	ChannelSMS      NotificationChannel = "SMS"
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
	Channel        NotificationChannel
	TemplateCode   string
	Status         string
	CreatedAt      time.Time
}