package server

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
