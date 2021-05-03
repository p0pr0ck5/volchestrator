package memory

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/p0pr0ck5/volchestrator/server/notification"
)

func (m *Memory) WriteNotification(n *notification.Notification) error {
	clientID := n.ClientID

	ch, exists := m.notificationChMap[clientID]
	if !exists {
		return errors.New(fmt.Sprintf("no notification channel for client %q", clientID))
	}

	ch <- n

	return nil
}

func (m *Memory) GetNotifications(id string) (<-chan *notification.Notification, error) {
	return m.notificationChMap[id], nil
}
