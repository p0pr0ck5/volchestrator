package memory

import (
	"reflect"
	"sync"
	"testing"

	"github.com/p0pr0ck5/volchestrator/server/volume"
)

func TestMemory_ReadVolume(t *testing.T) {
	type fields struct {
		volumeMap volumeMap
		dataLocks map[string]*sync.Mutex
	}
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *volume.Volume
		wantErr bool
	}{
		{
			"read volume",
			fields{
				volumeMap: map[string]*volume.Volume{
					"foo": {
						ID: "foo",
					},
				},
				dataLocks: make(map[string]*sync.Mutex),
			},
			args{
				"foo",
			},
			&volume.Volume{
				ID: "foo",
			},
			false,
		},
		{
			"read non-existent volume",
			fields{
				volumeMap: map[string]*volume.Volume{
					"foo": {
						ID: "foo",
					},
				},
				dataLocks: make(map[string]*sync.Mutex),
			},
			args{
				"bar",
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Memory{
				volumeMap: tt.fields.volumeMap,
				dataLocks: tt.fields.dataLocks,
			}
			got, err := m.ReadVolume(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Memory.ReadVolume() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Memory.ReadVolume() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemory_ListVolumes(t *testing.T) {
	type fields struct {
		volumeMap volumeMap
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
				volumeMap: map[string]*volume.Volume{
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
				volumeMap: map[string]*volume.Volume{
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
					ID: "foo",
				},
				{
					ID: "bar",
				},
			},
			false,
		},
		{
			"zero volumes",
			fields{
				volumeMap: map[string]*volume.Volume{},
			},
			[]*volume.Volume{},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Memory{
				volumeMap: tt.fields.volumeMap,
				dataLocks: tt.fields.dataLocks,
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

func TestMemory_CreateVolume(t *testing.T) {
	type fields struct {
		volumeMap volumeMap
		dataLocks map[string]*sync.Mutex
	}
	type args struct {
		volume *volume.Volume
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"create new volume",
			fields{
				volumeMap: map[string]*volume.Volume{
					"foo": {
						ID: "foo",
					},
				},
				dataLocks: make(map[string]*sync.Mutex),
			},
			args{
				volume: &volume.Volume{
					ID: "bar",
				},
			},
			false,
		},
		{
			"create existing volume",
			fields{
				volumeMap: map[string]*volume.Volume{
					"foo": {
						ID: "foo",
					},
				},
				dataLocks: make(map[string]*sync.Mutex),
			},
			args{
				volume: &volume.Volume{
					ID: "foo",
				},
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Memory{
				volumeMap: tt.fields.volumeMap,
				dataLocks: tt.fields.dataLocks,
			}
			if err := m.CreateVolume(tt.args.volume); (err != nil) != tt.wantErr {
				t.Errorf("Memory.CreateVolume() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMemory_UpdateVolume(t *testing.T) {
	type fields struct {
		volumeMap map[string]*volume.Volume
		dataLocks map[string]*sync.Mutex
	}
	type args struct {
		volume *volume.Volume
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"update existing volume",
			fields{
				volumeMap: map[string]*volume.Volume{
					"foo": {
						ID: "foo",
					},
				},
				dataLocks: make(map[string]*sync.Mutex),
			},
			args{
				volume: &volume.Volume{
					ID: "foo",
				},
			},
			false,
		},
		{
			"update nonexistent volume",
			fields{
				volumeMap: map[string]*volume.Volume{
					"foo": {
						ID: "foo",
					},
				},
				dataLocks: make(map[string]*sync.Mutex),
			},
			args{
				volume: &volume.Volume{
					ID: "bar",
				},
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Memory{
				volumeMap: tt.fields.volumeMap,
				dataLocks: tt.fields.dataLocks,
			}
			if err := m.UpdateVolume(tt.args.volume); (err != nil) != tt.wantErr {
				t.Errorf("Memory.UpdateVolume() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMemory_DeleteVolume(t *testing.T) {
	type fields struct {
		volumeMap volumeMap
		dataLocks map[string]*sync.Mutex
	}
	type args struct {
		volume *volume.Volume
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"delete existing volume",
			fields{
				volumeMap: map[string]*volume.Volume{
					"foo": {
						ID: "foo",
					},
				},
				dataLocks: make(map[string]*sync.Mutex),
			},
			args{
				volume: &volume.Volume{
					ID: "foo",
				},
			},
			false,
		},
		{
			"delete nonexistent volume",
			fields{
				volumeMap: map[string]*volume.Volume{
					"foo": {
						ID: "foo",
					},
				},
				dataLocks: make(map[string]*sync.Mutex),
			},
			args{
				volume: &volume.Volume{
					ID: "bar",
				},
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Memory{
				volumeMap: tt.fields.volumeMap,
				dataLocks: tt.fields.dataLocks,
			}
			if err := m.DeleteVolume(tt.args.volume); (err != nil) != tt.wantErr {
				t.Errorf("Memory.DeleteVolume() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
