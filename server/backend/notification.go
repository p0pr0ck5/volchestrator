package backend

import "github.com/p0pr0ck5/volchestrator/server/notification"

func (b *Backend) WriteNotification(n *notification.Notification) error {
	return b.b.WriteNotification(n)
}

func (b *Backend) GetNotifications(id string) (<-chan *notification.Notification, error) {
	return b.b.GetNotifications(id)
}
