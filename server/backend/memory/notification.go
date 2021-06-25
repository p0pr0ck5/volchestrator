package memory

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/p0pr0ck5/volchestrator/server/notification"
)

func (m *Memory) WriteNotification(n *notification.Notification) error {
	clientID := n.ClientID

	queue, ok := m.notificationMap[clientID]
	if !ok {
		return errors.New(fmt.Sprintf("no notification queue for client %q", clientID))
	}

	return queue.Write(n)
}

func (m *Memory) GetNotifications(id string) (<-chan *notification.Notification, error) {
	queue, ok := m.notificationMap[id]
	if !ok {
		return nil, errors.New(fmt.Sprintf("no notification queue for client %q", id))
	}

	ch, err := queue.Read()
	if err != nil {
		return nil, err
	}

	return ch, nil
}
