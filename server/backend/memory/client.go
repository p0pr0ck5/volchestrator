package memory

import (
	"sort"

	"github.com/p0pr0ck5/volchestrator/server/client"
)

type clientMap map[string]*client.Client

func (c clientMap) Get(id string) (interface{}, bool) {
	entity, exists := c[id]
	return entity, exists
}

func (c clientMap) List() interface{} {
	list := []*client.Client{}

	for _, client := range c {
		list = append(list, client)
	}

	return list
}

func (c clientMap) Set(id string, entity interface{}) {
	e := entity.(*client.Client)
	c[id] = e
}

func (c clientMap) Delete(id string) {
	delete(c, id)
}

func (m *Memory) ReadClient(id string) (*client.Client, error) {
	c, err := m.read(id, "client")
	if err != nil {
		return nil, err
	}

	return c.(*client.Client), nil
}

func (m *Memory) ListClients() ([]*client.Client, error) {
	clients := m.list("client").([]*client.Client)

	sort.Slice(clients, func(i, j int) bool {
		return clients[i].ID < clients[j].ID
	})

	return clients, nil
}

func (m *Memory) CreateClient(client *client.Client) error {
	if err := m.cud("create", client); err != nil {
		return err
	}

	queue, err := NewChQueue()
	if err != nil {
		return err
	}
	m.notificationMap[client.ID] = queue

	return nil
}

func (m *Memory) UpdateClient(client *client.Client) error {
	return m.cud("update", client)
}

func (m *Memory) DeleteClient(client *client.Client) error {
	if err := m.cud("delete", client); err != nil {
		return err
	}

	queue := m.notificationMap[client.ID]
	if err := queue.Close(); err != nil {
		return err
	}
	delete(m.notificationMap, client.ID)

	return nil
}
