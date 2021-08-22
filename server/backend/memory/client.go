package memory

import (
	"sort"

	"github.com/p0pr0ck5/volchestrator/server/client"
)

type ClientMap map[string]*client.Client

func (c ClientMap) Get(id string) (interface{}, bool) {
	entity, exists := c[id]
	return entity, exists
}

func (c ClientMap) List() interface{} {
	list := []*client.Client{}

	for _, client := range c {
		list = append(list, client)
	}

	return list
}

func (c ClientMap) Set(id string, entity interface{}) {
	e := entity.(*client.Client)
	c[id] = e
}

func (c ClientMap) Delete(id string) {
	delete(c, id)
}

func (m *Memory) ListClients() ([]*client.Client, error) {
	clients := m.list("client").([]*client.Client)

	sort.Slice(clients, func(i, j int) bool {
		return clients[i].ID < clients[j].ID
	})

	return clients, nil
}
