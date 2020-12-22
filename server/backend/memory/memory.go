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
	volumeMap *VolumeMap

	clientMap *ClientMap
}

// NewMemoryBackend creates an initialized empty MemoryBackend
func NewMemoryBackend() *MemoryBackend {
	m := &MemoryBackend{
		volumeMap: NewVolumeMap(),
		clientMap: NewClientMap(),
	}

	return m
}

// VolumeMap holds information about registered volumes
type VolumeMap struct {
	m map[string]*server.Volume
	l sync.Mutex
}

// NewVolumeMap returns an initialized VolumeMap
func NewVolumeMap() *VolumeMap {
	m := &VolumeMap{
		m: make(map[string]*server.Volume),
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

// AddClient adds a Client to the backend if it doesn't already exist
func (m *MemoryBackend) AddClient(id string) error {
	m.clientMap.l.Lock()
	defer m.clientMap.l.Unlock()

	if _, exists := m.clientMap.m[id]; exists {
		return fmt.Errorf("Client %q already exists in memory backend", id)
	}

	m.clientMap.m[id] = server.ClientInfo{
		ID:        id,
		Status:    server.UnknownStatus,
		FirstSeen: time.Now(),
	}

	return nil
}

// UpdateClient updates the client info for a given client
func (m *MemoryBackend) UpdateClient(id string, status server.ClientStatus) error {
	m.clientMap.l.Lock()
	defer m.clientMap.l.Unlock()

	var client server.ClientInfo
	var exists bool
	if client, exists = m.clientMap.m[id]; !exists {
		return fmt.Errorf("Client %q does not exist in memory backend", id)
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
	m.volumeMap.l.Lock()
	defer m.volumeMap.l.Unlock()

	v, ok := m.volumeMap.m[id]
	if !ok {
		log.Printf("No volume %q found in memory map\n", id)
	}

	return v, nil
}

// ListVolumes satisfies server.Backend
func (m *MemoryBackend) ListVolumes() ([]*server.Volume, error) {
	m.volumeMap.l.Lock()
	defer m.volumeMap.l.Unlock()

	volumes := []*server.Volume{}

	for _, volume := range m.volumeMap.m {
		volumes = append(volumes, volume)
	}

	return volumes, nil
}

// AddVolume satisfies server.Backend
func (m *MemoryBackend) AddVolume(volume *server.Volume) error {
	m.volumeMap.l.Lock()
	defer m.volumeMap.l.Unlock()

	if _, exists := m.volumeMap.m[volume.ID]; exists {
		return fmt.Errorf("Volume %q already exists in memory backend", volume.ID)
	}

	m.volumeMap.m[volume.ID] = volume

	return nil
}

// UpdateVolume satisfies server.Backend
func (m *MemoryBackend) UpdateVolume(volume *server.Volume) error {
	m.volumeMap.l.Lock()
	defer m.volumeMap.l.Unlock()

	if _, exists := m.volumeMap.m[volume.ID]; !exists {
		return fmt.Errorf("Volume %q does not exist in memory backend", volume.ID)
	}

	m.volumeMap.m[volume.ID] = volume

	return nil
}

// DeleteVolume satisfies server.Backend
func (m *MemoryBackend) DeleteVolume(id string) error {
	m.volumeMap.l.Lock()
	defer m.volumeMap.l.Unlock()

	if _, exists := m.volumeMap.m[id]; !exists {
		return fmt.Errorf("Volume %q does not exist in memory backend", id)
	}

	delete(m.volumeMap.m, id)

	return nil
}
