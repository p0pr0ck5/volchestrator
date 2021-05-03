package memory

import "github.com/p0pr0ck5/volchestrator/server/notification"

func (m *Memory) GetNotifications(id string) <-chan *notification.Notification {
	return m.notificationChMap[id]
}
