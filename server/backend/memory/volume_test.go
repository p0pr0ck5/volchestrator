package memory

import (
	"reflect"
	"sync"
	"testing"

	"github.com/p0pr0ck5/volchestrator/server/volume"
)

func TestMemory_ListVolumes(t *testing.T) {
	type fields struct {
		VolumeMap VolumeMap
		dataLocks map[string]*sync.Mutex
	}
	tests := []struct {
		name    string
		fields  fields
		want    []*volume.Volume
		wantErr bool
	}{
		{
			"one volume",
			fields{
				VolumeMap: map[string]*volume.Volume{
					"foo": {
						ID: "foo",
					},
				},
				dataLocks: make(map[string]*sync.Mutex),
			},
			[]*volume.Volume{
				{
					ID: "foo",
				},
			},
			false,
		},
		{
			"two volumes",
			fields{
				VolumeMap: map[string]*volume.Volume{
					"foo": {
						ID: "foo",
					},
					"bar": {
						ID: "bar",
					},
				},
			},
			[]*volume.Volume{
				{
					ID: "bar",
				},
				{
					ID: "foo",
				},
			},
			false,
		},
		{
			"zero volumes",
			fields{
				VolumeMap: map[string]*volume.Volume{},
			},
			[]*volume.Volume{},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Memory{
				VolumeMap: tt.fields.VolumeMap,
			}
			got, err := m.ListVolumes()
			if (err != nil) != tt.wantErr {
				t.Errorf("Memory.ListVolumes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Memory.ListVolumes() = %v, want %v", got, tt.want)
			}
		})
	}
}
