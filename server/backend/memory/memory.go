package memory

import "sync"

type Memory struct {
	ClientMap ClientMap
	VolumeMap VolumeMap

	notificationMap map[string]*ChQueue

	l sync.RWMutex
}

func NewMemoryBackend() *Memory {
	m := &Memory{
		ClientMap:       make(ClientMap),
		VolumeMap:       make(VolumeMap),
		notificationMap: make(map[string]*ChQueue),
	}

	return m
}
