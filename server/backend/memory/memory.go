package memory

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/p0pr0ck5/volchestrator/server"
)

// MemoryBackend implements server.Backend
type MemoryBackend struct {
	volumeMap map[string]*server.Volume

	clientMap *ClientMap
}

// NewMemoryBackend creates an initialized empty MemoryBackend
func NewMemoryBackend() *MemoryBackend {
	m := &MemoryBackend{
		volumeMap: make(map[string]*server.Volume),
		clientMap: NewClientMap(),
	}

	return m
}

// ClientMap maps clients to their info
type ClientMap struct {
	m map[string]server.ClientInfo
	l sync.Mutex
}

// NewClientMap returns an initialized ClientMap
func NewClientMap() *ClientMap {
	m := &ClientMap{
		m: make(map[string]server.ClientInfo),
	}

	return m
}

// UpdateClient updates the client info for a given client
func (m *MemoryBackend) UpdateClient(id string, status server.ClientStatus) error {
	m.clientMap.l.Lock()
	defer m.clientMap.l.Unlock()

	var client server.ClientInfo
	var ok bool
	if client, ok = m.clientMap.m[id]; !ok {
		client = server.ClientInfo{
			ID:        id,
			FirstSeen: time.Now(),
		}
	}

	client.LastSeen = time.Now()
	client.Status = status

	m.clientMap.m[id] = client

	return nil
}

// RemoveClient deletes a given client from the ClientMap
func (m *MemoryBackend) RemoveClient(id string) error {
	m.clientMap.l.Lock()
	defer m.clientMap.l.Unlock()

	delete(m.clientMap.m, id)

	return nil
}

// Clients returns a list of server.ClientInfo
func (m *MemoryBackend) Clients(f server.ClientFilterFunc) ([]server.ClientInfo, error) {
	m.clientMap.l.Lock()
	defer m.clientMap.l.Unlock()

	var c []server.ClientInfo
	for _, ci := range m.clientMap.m {
		if f(ci) {
			c = append(c, ci)
		}
	}

	return c, nil
}

// GetVolume satisfies server.Backend
func (m *MemoryBackend) GetVolume(id string) (*server.Volume, error) {
	v, ok := m.volumeMap[id]
	if !ok {
		log.Printf("No volume %q found in memory map\n", id)
	}

	return v, nil
}

// ListVolumes satisfies server.Backend
func (m *MemoryBackend) ListVolumes() ([]*server.Volume, error) {
	volumes := []*server.Volume{}

	for _, volume := range m.volumeMap {
		volumes = append(volumes, volume)
	}

	return volumes, nil
}

// AddVolume satisfies server.Backend
func (m *MemoryBackend) AddVolume(volume *server.Volume) error {
	if _, exists := m.volumeMap[volume.ID]; exists {
		return fmt.Errorf("Volume %q already exists in memory backend", volume.ID)
	}

	m.volumeMap[volume.ID] = volume

	return nil
}

// UpdateVolume satisfies server.Backend
func (m *MemoryBackend) UpdateVolume(volume *server.Volume) error {
	if _, exists := m.volumeMap[volume.ID]; !exists {
		return fmt.Errorf("Volume %q does not exist in memory backend", volume.ID)
	}

	m.volumeMap[volume.ID] = volume

	return nil
}

// DeleteVolume satisfies server.Backend
func (m *MemoryBackend) DeleteVolume(id string) error {
	if _, exists := m.volumeMap[id]; !exists {
		return fmt.Errorf("Volume %q does not exist in memory backend", id)
	}

	delete(m.volumeMap, id)

	return nil
}
