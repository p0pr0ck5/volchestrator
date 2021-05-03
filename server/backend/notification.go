package backend

import "github.com/p0pr0ck5/volchestrator/server/notification"

func (b *Backend) GetNotifications(id string) <-chan *notification.Notification {
	return b.b.GetNotifications(id)
}
