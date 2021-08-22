package memory

import (
	"sort"

	"github.com/p0pr0ck5/volchestrator/server/volume"
)

type VolumeMap map[string]*volume.Volume

func (v VolumeMap) Get(id string) (interface{}, bool) {
	entity, exists := v[id]
	return entity, exists
}

func (v VolumeMap) List() interface{} {
	list := []*volume.Volume{}

	for _, volume := range v {
		list = append(list, volume)
	}

	return list
}

func (v VolumeMap) Set(id string, entity interface{}) {
	e := entity.(*volume.Volume)
	v[id] = e
}

func (v VolumeMap) Delete(id string) {
	delete(v, id)
}

func (m *Memory) ListVolumes() ([]*volume.Volume, error) {
	volumes := m.list("volume").([]*volume.Volume)

	sort.Slice(volumes, func(i, j int) bool {
		return volumes[i].ID < volumes[j].ID
	})

	return volumes, nil
}
