package backend

import (
	"time"

	"github.com/p0pr0ck5/volchestrator/server/client"
)

func (b *Backend) ReadClient(id string) (*client.Client, error) {
	return b.b.ReadClient(id)
}

func (b *Backend) ListClients() ([]*client.Client, error) {
	return b.b.ListClients()
}

func (b *Backend) CreateClient(c *client.Client) error {
	now := time.Now()
	c.CreatedAt = now
	c.UpdatedAt = now

	return b.b.CreateClient(c)
}

func (b *Backend) UpdateClient(c *client.Client) error {
	c.UpdatedAt = time.Now()
	return b.b.UpdateClient(c)
}

func (b *Backend) DeleteClient(c *client.Client) error {
	return b.b.DeleteClient(c)
}
