package memory

import (
	"github.com/p0pr0ck5/volchestrator/server/volume"
)

type volumeMap map[string]*volume.Volume

func (v volumeMap) Get(id string) (interface{}, bool) {
	entity, exists := v[id]
	return entity, exists
}

func (v volumeMap) List() interface{} {
	list := []*volume.Volume{}

	for _, volume := range v {
		list = append(list, volume)
	}

	return list
}

func (v volumeMap) Set(id string, entity interface{}) {
	e := entity.(*volume.Volume)
	v[id] = e
}

func (v volumeMap) Delete(id string) {
	delete(v, id)
}

func (m *Memory) ReadVolume(id string) (*volume.Volume, error) {
	c, err := m.read(id, "volume")
	if err != nil {
		return nil, err
	}

	return c.(*volume.Volume), nil
}

func (m *Memory) ListVolumes() ([]*volume.Volume, error) {
	return m.list("volume").([]*volume.Volume), nil
}

func (m *Memory) CreateVolume(volume *volume.Volume) error {
	return m.cud("create", volume)
}

func (m *Memory) UpdateVolume(volume *volume.Volume) error {
	return m.cud("update", volume)
}

func (m *Memory) DeleteVolume(volume *volume.Volume) error {
	return m.cud("delete", volume)
}
