package server

// NotificationType defines the type of notification sent to the client
// during streaming
type NotificationType int

const (
	// UnknownType is a base value
	UnknownType NotificationType = iota
)

// Notification is a message to be passed to the client
type Notification struct {
	ID      string
	Type    NotificationType
	Message string
}
