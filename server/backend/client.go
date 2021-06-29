package backend

import (
	"reflect"
	"sort"
	"time"

	"github.com/imdario/mergo"
	"github.com/p0pr0ck5/volchestrator/server/client"
	"github.com/pkg/errors"
)

type timeTransformer struct{}

func (t timeTransformer) Transformer(typ reflect.Type) func(dst, src reflect.Value) error {
	if typ == reflect.TypeOf(time.Time{}) {
		return func(dst, src reflect.Value) error {
			if dst.CanSet() {
				isZero := dst.MethodByName("IsZero")
				result := isZero.Call([]reflect.Value{})
				if result[0].Bool() {
					dst.Set(src)
				}
			}
			return nil
		}
	}
	return nil
}

func (b *Backend) ReadClient(id string) (*client.Client, error) {
	return b.b.ReadClient(id)
}

func (b *Backend) ListClients() ([]*client.Client, error) {
	clients, err := b.b.ListClients()

	sort.Slice(clients, func(i, j int) bool {
		return clients[i].ID < clients[j].ID
	})

	return clients, err
}

func (b *Backend) CreateClient(c *client.Client) error {
	now := time.Now()
	c.CreatedAt = now
	c.UpdatedAt = now

	return b.b.CreateClient(c)
}

func (b *Backend) UpdateClient(c *client.Client) error {
	currentClient, err := b.ReadClient(c.ID)
	if err != nil {
		return errors.Wrap(err, "get current client")
	}

	if err := mergo.Merge(c, currentClient, mergo.WithTransformers(timeTransformer{})); err != nil {
		return errors.Wrap(err, "merge client")
	}

	c.UpdatedAt = time.Now()

	return b.b.UpdateClient(c)
}

func (b *Backend) DeleteClient(c *client.Client) error {
	return b.b.DeleteClient(c)
}
