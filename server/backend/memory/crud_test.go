package memory

import (
	"reflect"
	"testing"

	"github.com/p0pr0ck5/volchestrator/server/client"
	"github.com/p0pr0ck5/volchestrator/server/model"
	"github.com/p0pr0ck5/volchestrator/server/volume"
)

func TestMemory_Create(t *testing.T) {
	type fields struct {
		ClientMap       ClientMap
		VolumeMap       VolumeMap
		notificationMap map[string]*ChQueue
	}
	type args struct {
		entity model.Base
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"create new client",
			fields{
				ClientMap:       ClientMap{},
				VolumeMap:       VolumeMap{},
				notificationMap: make(map[string]*ChQueue),
			},
			args{
				&client.Client{
					ID: "foo",
				},
			},
			false,
		},
		{
			"create existing client",
			fields{
				ClientMap: map[string]*client.Client{
					"foo": {
						ID: "foo",
					},
				},
				notificationMap: map[string]*ChQueue{
					"foo": MustNewChQueue(),
				},
			},
			args{
				entity: &client.Client{
					ID: "foo",
				},
			},
			true,
		},
		{
			"create new volume",
			fields{
				VolumeMap: map[string]*volume.Volume{
					"foo": {
						ID: "foo",
					},
				},
			},
			args{
				entity: &volume.Volume{
					ID: "bar",
				},
			},
			false,
		},
		{
			"create existing volume",
			fields{
				VolumeMap: map[string]*volume.Volume{
					"foo": {
						ID: "foo",
					},
				},
			},
			args{
				entity: &volume.Volume{
					ID: "foo",
				},
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Memory{
				ClientMap:       tt.fields.ClientMap,
				VolumeMap:       tt.fields.VolumeMap,
				notificationMap: tt.fields.notificationMap,
			}
			if err := m.Create(tt.args.entity); (err != nil) != tt.wantErr {
				t.Errorf("Memory.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMemory_Read(t *testing.T) {
	type fields struct {
		ClientMap       ClientMap
		VolumeMap       VolumeMap
		notificationMap map[string]*ChQueue
	}
	type args struct {
		entity model.Base
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    model.Base
		wantErr bool
	}{
		{
			"read client",
			fields{
				ClientMap: map[string]*client.Client{
					"foo": {
						ID: "foo",
					},
				},
				notificationMap: map[string]*ChQueue{
					"foo": MustNewChQueue(),
				},
			},
			args{
				entity: &client.Client{
					ID: "foo",
				},
			},
			&client.Client{
				ID: "foo",
			},
			false,
		},
		{
			"read non-existent client",
			fields{
				ClientMap: map[string]*client.Client{
					"foo": {
						ID: "foo",
					},
				},
				notificationMap: map[string]*ChQueue{
					"foo": MustNewChQueue(),
				},
			},
			args{
				entity: &client.Client{
					ID: "bar",
				},
			},
			nil,
			true,
		},
		{
			"read volume",
			fields{
				VolumeMap: map[string]*volume.Volume{
					"foo": {
						ID: "foo",
					},
				},
			},
			args{
				entity: &volume.Volume{
					ID: "foo",
				},
			},
			&volume.Volume{
				ID: "foo",
			},
			false,
		},
		{
			"read non-existent volume",
			fields{
				VolumeMap: map[string]*volume.Volume{
					"foo": {
						ID: "foo",
					},
				},
			},
			args{
				entity: &volume.Volume{
					ID: "bar",
				},
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Memory{
				ClientMap:       tt.fields.ClientMap,
				VolumeMap:       tt.fields.VolumeMap,
				notificationMap: tt.fields.notificationMap,
			}
			got, err := m.Read(tt.args.entity)
			if (err != nil) != tt.wantErr {
				t.Errorf("Memory.Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Memory.Read() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemory_Update(t *testing.T) {
	type fields struct {
		ClientMap       ClientMap
		VolumeMap       VolumeMap
		notificationMap map[string]*ChQueue
	}
	type args struct {
		entity model.Base
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"update existing client",
			fields{
				ClientMap: map[string]*client.Client{
					"foo": {
						ID: "foo",
					},
				},
				notificationMap: map[string]*ChQueue{
					"foo": MustNewChQueue(),
				},
			},
			args{
				entity: &client.Client{
					ID: "foo",
				},
			},
			false,
		},
		{
			"update nonexistent client",
			fields{
				ClientMap: map[string]*client.Client{
					"foo": {
						ID: "foo",
					},
				},
				notificationMap: map[string]*ChQueue{
					"foo": MustNewChQueue(),
				},
			},
			args{
				entity: &client.Client{
					ID: "bar",
				},
			},
			true,
		},
		{
			"update existing volume",
			fields{
				VolumeMap: map[string]*volume.Volume{
					"foo": {
						ID: "foo",
					},
				},
			},
			args{
				entity: &volume.Volume{
					ID: "foo",
				},
			},
			false,
		},
		{
			"update nonexistent volume",
			fields{
				VolumeMap: map[string]*volume.Volume{
					"foo": {
						ID: "foo",
					},
				},
			},
			args{
				entity: &volume.Volume{
					ID: "bar",
				},
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Memory{
				ClientMap:       tt.fields.ClientMap,
				VolumeMap:       tt.fields.VolumeMap,
				notificationMap: tt.fields.notificationMap,
			}
			if err := m.Update(tt.args.entity); (err != nil) != tt.wantErr {
				t.Errorf("Memory.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMemory_Delete(t *testing.T) {
	type fields struct {
		ClientMap       ClientMap
		VolumeMap       VolumeMap
		notificationMap map[string]*ChQueue
	}
	type args struct {
		entity model.Base
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"delete existing client",
			fields{
				ClientMap: map[string]*client.Client{
					"foo": {
						ID: "foo",
					},
				},
				notificationMap: map[string]*ChQueue{
					"foo": MustNewChQueue(),
				},
			},
			args{
				entity: &client.Client{
					ID: "foo",
				},
			},
			false,
		},
		{
			"delete nonexistent client",
			fields{
				ClientMap: map[string]*client.Client{
					"foo": {
						ID: "foo",
					},
				},
			},
			args{
				entity: &client.Client{
					ID: "bar",
				},
			},
			true,
		},
		{
			"delete existing volume",
			fields{
				VolumeMap: map[string]*volume.Volume{
					"foo": {
						ID: "foo",
					},
				},
			},
			args{
				entity: &volume.Volume{
					ID: "foo",
				},
			},
			false,
		},
		{
			"delete nonexistent volume",
			fields{
				VolumeMap: map[string]*volume.Volume{
					"foo": {
						ID: "foo",
					},
				},
			},
			args{
				entity: &volume.Volume{
					ID: "bar",
				},
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Memory{
				ClientMap:       tt.fields.ClientMap,
				VolumeMap:       tt.fields.VolumeMap,
				notificationMap: tt.fields.notificationMap,
			}
			if err := m.Delete(tt.args.entity); (err != nil) != tt.wantErr {
				t.Errorf("Memory.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
