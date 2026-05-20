package domain

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