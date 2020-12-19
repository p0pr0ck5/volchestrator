package memory

import (
	"fmt"
	"log"

	"github.com/p0pr0ck5/volchestrator/server"
)

// MemoryBackend implements server.Backend
type MemoryBackend struct {
	volumeMap map[string]*server.Volume
}

// NewMemoryBackend creates an initialized empty MemoryBackend
func NewMemoryBackend() *MemoryBackend {
	m := &MemoryBackend{
		volumeMap: make(map[string]*server.Volume),
	}

	return m
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
func (m *MemoryBackend) DeleteVolume(volume *server.Volume) error {
	if _, exists := m.volumeMap[volume.ID]; !exists {
		return fmt.Errorf("Volume %q does not exist in memory backend", volume.ID)
	}

	delete(m.volumeMap, volume.ID)

	return nil
}
