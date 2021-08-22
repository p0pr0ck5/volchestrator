package memory

import (
	"fmt"
)

type dataMap interface {
	Get(string) (interface{}, bool)
	List() interface{}
	Set(string, interface{})
	Delete(string)
}

type Memory struct {
	ClientMap ClientMap
	VolumeMap VolumeMap

	notificationMap map[string]*ChQueue
}

func NewMemoryBackend() *Memory {
	m := &Memory{
		ClientMap:       make(ClientMap),
		VolumeMap:       make(VolumeMap),
		notificationMap: make(map[string]*ChQueue),
	}

	return m
}

func (m *Memory) getMap(entityType string) dataMap {
	var e dataMap

	switch entityType {
	case "client":
		e = m.ClientMap
	case "volume":
		e = m.VolumeMap
	default:
		panic(fmt.Sprintf("invalid entity type %q", entityType))
	}

	return e
}

func (m *Memory) list(entityType string) interface{} {
	return m.getMap(entityType).List()
}
