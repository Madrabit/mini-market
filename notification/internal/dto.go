package internal

import (
	"github.com/google/uuid"
	"time"
)

type Type string
type Status string

const (
	Email  Type   = "email"
	Push   Type   = "push"
	SMS    Type   = "sms"
	Sent   Status = "sent"
	Queued Status = "queued"
	Failed Status = "failed"
)

type NotificationRequest struct {
	UserID    uuid.UUID
	From      string
	To        string
	Type      Type // email, sms, push
	Subject   string
	Text      string
	CreatedAt time.Time
}

type NotificationResponse struct {
	Success   bool
	MessageID string // ID уведомления в вашей системе
	Status    string // "sent", "queued", "failed"
	Timestamp time.Time
}

type UserNotificationsResponse struct {
	UserID        uuid.UUID
	Notifications []NotificationResponse
}
