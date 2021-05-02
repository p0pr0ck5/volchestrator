package memory

import (
	"fmt"

	"github.com/p0pr0ck5/volchestrator/server/client"
	"github.com/p0pr0ck5/volchestrator/server/volume"
)

type dataMap interface {
	Get(string) (interface{}, bool)
	List() interface{}
	Set(string, interface{})
	Delete(string)
}

type Memory struct {
	clientMap clientMap
	volumeMap volumeMap
}

func NewMemoryBackend() *Memory {
	m := &Memory{
		clientMap: make(map[string]*client.Client),
		volumeMap: make(map[string]*volume.Volume),
	}

	return m
}

func (m *Memory) getMap(entityType string) dataMap {
	var e dataMap

	switch entityType {
	case "client":
		e = m.clientMap
	case "volume":
		e = m.volumeMap
	default:
		panic(fmt.Sprintf("invalid entity type %q", entityType))
	}

	return e
}

func (m *Memory) read(id, entityType string) (interface{}, error) {
	entity, exists := m.getMap(entityType).Get(id)
	if !exists {
		return nil, fmt.Errorf("%s %q not found", entityType, id)
	}

	return entity, nil
}

func (m *Memory) list(entityType string) interface{} {
	return m.getMap(entityType).List()
}

func (m *Memory) cud(op string, entity interface{}) error {
	var e dataMap
	var i, ctx string

	switch t := entity.(type) {
	case *client.Client:
		e = m.clientMap
		i = entity.(*client.Client).ID
		ctx = "client"
	case *volume.Volume:
		e = m.volumeMap
		i = entity.(*volume.Volume).ID
		ctx = "volume"
	default:
		panic(fmt.Sprintf("invalid entity type %q", t))
	}

	_, exists := e.Get(i)

	switch op {
	case "create":
		if exists {
			return fmt.Errorf("%s %q already exists", ctx, i)
		}

		e.Set(i, entity)
	case "update":
		if !exists {
			return fmt.Errorf("%s %q not found", ctx, i)
		}

		e.Set(i, entity)
	case "delete":
		if !exists {
			return fmt.Errorf("%s %q not found", ctx, i)
		}

		e.Delete(i)
	}

	return nil
}
