package server

import (
	"github.com/thanhpk/randstr"
)

// NotificationType defines the type of notification sent to the client
// during streaming
type NotificationType int

const (
	// UnknownNotificationType is a base value
	UnknownNotificationType NotificationType = iota

	// LeaseRequestAckNotificationType is an acknowledgement of receipt of a LeaseRequest submission
	LeaseRequestAckNotificationType
)

// Notification is a message to be passed to the client
type Notification struct {
	ID      string
	Type    NotificationType
	Message string
}

// NewNotification returns a new Notification with a given type and message
func NewNotification(t NotificationType, message string) Notification {
	return Notification{
		ID:      randstr.Hex(16),
		Type:    t,
		Message: message,
	}
}
