package memory

import (
	"reflect"
	"testing"

	"github.com/p0pr0ck5/volchestrator/server/client"
	"github.com/p0pr0ck5/volchestrator/server/lease"
	leaserequest "github.com/p0pr0ck5/volchestrator/server/lease_request"
	"github.com/p0pr0ck5/volchestrator/server/model"
	"github.com/p0pr0ck5/volchestrator/server/volume"
)

func TestMemory_Create(t *testing.T) {
	type fields struct {
		ClientMap       ClientMap
		VolumeMap       VolumeMap
		LeaseRequestMap LeaseRequestMap
		LeaseMap        LeaseMap
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
		{
			"create new lease request",
			fields{
				LeaseRequestMap: map[string]*leaserequest.LeaseRequest{
					"foo": {
						ID: "foo",
					},
				},
			},
			args{
				entity: &leaserequest.LeaseRequest{
					ID: "bar",
				},
			},
			false,
		},
		{
			"create existing lease request",
			fields{
				LeaseRequestMap: map[string]*leaserequest.LeaseRequest{
					"foo": {
						ID: "foo",
					},
				},
			},
			args{
				entity: &leaserequest.LeaseRequest{
					ID: "foo",
				},
			},
			true,
		},
		{
			"create new lease",
			fields{
				LeaseMap: map[string]*lease.Lease{
					"foo": {
						ID: "foo",
					},
				},
			},
			args{
				entity: &lease.Lease{
					ID: "bar",
				},
			},
			false,
		},
		{
			"create existing lease",
			fields{
				LeaseMap: map[string]*lease.Lease{
					"foo": {
						ID: "foo",
					},
				},
			},
			args{
				entity: &lease.Lease{
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
				LeaseRequestMap: tt.fields.LeaseRequestMap,
				LeaseMap:        tt.fields.LeaseMap,
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
		LeaseRequestMap LeaseRequestMap
		LeaseMap        LeaseMap
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
		{
			"read lease request",
			fields{
				LeaseRequestMap: map[string]*leaserequest.LeaseRequest{
					"foo": {
						ID: "foo",
					},
				},
			},
			args{
				entity: &leaserequest.LeaseRequest{
					ID: "foo",
				},
			},
			&leaserequest.LeaseRequest{
				ID: "foo",
			},
			false,
		},
		{
			"read non-existing lease request",
			fields{
				LeaseRequestMap: map[string]*leaserequest.LeaseRequest{
					"foo": {
						ID: "foo",
					},
				},
			},
			args{
				entity: &leaserequest.LeaseRequest{
					ID: "bar",
				},
			},
			nil,
			true,
		},
		{
			"read existing lease",
			fields{
				LeaseMap: map[string]*lease.Lease{
					"foo": {
						ID: "foo",
					},
				},
			},
			args{
				entity: &lease.Lease{
					ID: "foo",
				},
			},
			&lease.Lease{
				ID: "foo",
			},
			false,
		},
		{
			"read non-existing lease",
			fields{
				LeaseMap: map[string]*lease.Lease{
					"foo": {
						ID: "foo",
					},
				},
			},
			args{
				entity: &lease.Lease{
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
				LeaseRequestMap: tt.fields.LeaseRequestMap,
				LeaseMap:        tt.fields.LeaseMap,
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
		LeaseRequestMap LeaseRequestMap
		LeaseMap        LeaseMap
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
		{
			"update existing lease request",
			fields{
				LeaseRequestMap: map[string]*leaserequest.LeaseRequest{
					"foo": {
						ID: "foo",
					},
				},
			},
			args{
				entity: &leaserequest.LeaseRequest{
					ID: "foo",
				},
			},
			false,
		},
		{
			"update non-existing lease request",
			fields{
				LeaseRequestMap: map[string]*leaserequest.LeaseRequest{
					"foo": {
						ID: "foo",
					},
				},
			},
			args{
				entity: &leaserequest.LeaseRequest{
					ID: "bar",
				},
			},
			true,
		},
		{
			"update existing lease",
			fields{
				LeaseMap: map[string]*lease.Lease{
					"foo": {
						ID: "foo",
					},
				},
			},
			args{
				entity: &lease.Lease{
					ID: "foo",
				},
			},
			false,
		},
		{
			"update non-existing lease",
			fields{
				LeaseMap: map[string]*lease.Lease{
					"foo": {
						ID: "foo",
					},
				},
			},
			args{
				entity: &lease.Lease{
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
				LeaseRequestMap: tt.fields.LeaseRequestMap,
				LeaseMap:        tt.fields.LeaseMap,
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
		LeaseRequestMap LeaseRequestMap
		LeaseMap        LeaseMap
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
		{
			"delete existing lease request",
			fields{
				LeaseRequestMap: map[string]*leaserequest.LeaseRequest{
					"foo": {
						ID: "foo",
					},
				},
			},
			args{
				entity: &leaserequest.LeaseRequest{
					ID: "foo",
				},
			},
			false,
		},
		{
			"delete non-existing lease request",
			fields{
				LeaseRequestMap: map[string]*leaserequest.LeaseRequest{
					"foo": {
						ID: "foo",
					},
				},
			},
			args{
				entity: &leaserequest.LeaseRequest{
					ID: "bar",
				},
			},
			true,
		},
		{
			"delete existing lease",
			fields{
				LeaseMap: map[string]*lease.Lease{
					"foo": {
						ID: "foo",
					},
				},
			},
			args{
				entity: &lease.Lease{
					ID: "foo",
				},
			},
			false,
		},
		{
			"delete non-existing lease",
			fields{
				LeaseMap: map[string]*lease.Lease{
					"foo": {
						ID: "foo",
					},
				},
			},
			args{
				entity: &lease.Lease{
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
				LeaseRequestMap: tt.fields.LeaseRequestMap,
				LeaseMap:        tt.fields.LeaseMap,
				notificationMap: tt.fields.notificationMap,
			}
			if err := m.Delete(tt.args.entity); (err != nil) != tt.wantErr {
				t.Errorf("Memory.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMemory_List(t *testing.T) {
	type fields struct {
		ClientMap       ClientMap
		VolumeMap       VolumeMap
		LeaseRequestMap LeaseRequestMap
		LeaseMap        LeaseMap
		notificationMap map[string]*ChQueue
	}
	type args struct {
		entityType string
		entities   *[]model.Base
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *[]model.Base
		wantErr bool
	}{
		{
			"one client",
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
				entityType: "Client",
				entities:   &[]model.Base{},
			},
			&[]model.Base{
				&client.Client{
					ID: "foo",
				},
			},
			false,
		},
		{
			"two clients",
			fields{
				ClientMap: map[string]*client.Client{
					"foo": {
						ID: "foo",
					},
					"bar": {
						ID: "bar",
					},
				},
				notificationMap: map[string]*ChQueue{
					"foo": MustNewChQueue(),
				},
			},
			args{
				entityType: "Client",
				entities:   &[]model.Base{},
			},
			&[]model.Base{
				&client.Client{
					ID: "bar",
				},
				&client.Client{
					ID: "foo",
				},
			},
			false,
		},
		{
			"zero clients",
			fields{
				ClientMap:       map[string]*client.Client{},
				notificationMap: map[string]*ChQueue{},
			},
			args{
				entityType: "Client",
				entities:   &[]model.Base{},
			},
			&[]model.Base{},
			false,
		},
		{
			"one volume",
			fields{
				VolumeMap: map[string]*volume.Volume{
					"foo": {
						ID: "foo",
					},
				},
			},
			args{
				entityType: "Volume",
				entities:   &[]model.Base{},
			},
			&[]model.Base{
				&volume.Volume{
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
			args{
				entityType: "Volume",
				entities:   &[]model.Base{},
			},
			&[]model.Base{
				&volume.Volume{
					ID: "bar",
				},
				&volume.Volume{
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
			args{
				entityType: "Volume",
				entities:   &[]model.Base{},
			},
			&[]model.Base{},
			false,
		},
		{
			"one lease request",
			fields{
				LeaseRequestMap: map[string]*leaserequest.LeaseRequest{
					"foo": {
						ID: "foo",
					},
				},
			},
			args{
				entityType: "LeaseRequest",
				entities:   &[]model.Base{},
			},
			&[]model.Base{
				&leaserequest.LeaseRequest{
					ID: "foo",
				},
			},
			false,
		},
		{
			"two lease requests",
			fields{
				LeaseRequestMap: map[string]*leaserequest.LeaseRequest{
					"foo": {
						ID: "foo",
					},
					"bar": {
						ID: "bar",
					},
				},
			},
			args{
				entityType: "LeaseRequest",
				entities:   &[]model.Base{},
			},
			&[]model.Base{
				&leaserequest.LeaseRequest{
					ID: "bar",
				},
				&leaserequest.LeaseRequest{
					ID: "foo",
				},
			},
			false,
		},
		{
			"zero lease requests",
			fields{
				LeaseRequestMap: map[string]*leaserequest.LeaseRequest{},
			},
			args{
				entityType: "LeaseRequest",
				entities:   &[]model.Base{},
			},
			&[]model.Base{},
			false,
		},
		{
			"one lease",
			fields{
				LeaseMap: map[string]*lease.Lease{
					"foo": {
						ID: "foo",
					},
				},
			},
			args{
				entityType: "Lease",
				entities:   &[]model.Base{},
			},
			&[]model.Base{
				&lease.Lease{
					ID: "foo",
				},
			},
			false,
		},
		{
			"two leases",
			fields{
				LeaseMap: map[string]*lease.Lease{
					"foo": {
						ID: "foo",
					},
					"bar": {
						ID: "bar",
					},
				},
			},
			args{
				entityType: "Lease",
				entities:   &[]model.Base{},
			},
			&[]model.Base{
				&lease.Lease{
					ID: "bar",
				},
				&lease.Lease{
					ID: "foo",
				},
			},
			false,
		},
		{
			"zero lease s",
			fields{
				LeaseMap: map[string]*lease.Lease{},
			},
			args{
				entityType: "Lease",
				entities:   &[]model.Base{},
			},
			&[]model.Base{},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Memory{
				ClientMap:       tt.fields.ClientMap,
				VolumeMap:       tt.fields.VolumeMap,
				LeaseRequestMap: tt.fields.LeaseRequestMap,
				LeaseMap:        tt.fields.LeaseMap,
				notificationMap: tt.fields.notificationMap,
			}
			if err := m.List(tt.args.entityType, tt.args.entities); (err != nil) != tt.wantErr {
				t.Errorf("Memory.List() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(tt.args.entities, tt.want) {
				t.Errorf("Memory.List() = %v, want %v", tt.args.entities, tt.want)

			}
		})
	}
}

func TestMemory_getMap(t *testing.T) {
	type fields struct {
		ClientMap       ClientMap
		VolumeMap       VolumeMap
		LeaseRequestMap LeaseRequestMap
		LeaseMap        LeaseMap
	}
	type args struct {
		entityType string
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantPanic bool
	}{
		{
			"client",
			fields{
				ClientMap: make(ClientMap),
			},
			args{
				"Client",
			},
			false,
		},
		{
			"volume",
			fields{
				VolumeMap: make(VolumeMap),
			},
			args{
				"Volume",
			},
			false,
		},
		{
			"lease request",
			fields{
				LeaseRequestMap: make(LeaseRequestMap),
			},
			args{
				"LeaseRequest",
			},
			false,
		},
		{
			"lease",
			fields{
				LeaseMap: make(LeaseMap),
			},
			args{
				"Lease",
			},
			false,
		},
		{
			"invalid",
			fields{
				ClientMap: make(ClientMap),
			},
			args{
				"nope",
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Memory{
				ClientMap:       tt.fields.ClientMap,
				VolumeMap:       tt.fields.VolumeMap,
				LeaseRequestMap: tt.fields.LeaseRequestMap,
				LeaseMap:        tt.fields.LeaseMap,
			}
			defer func() {
				if r := recover(); (r != nil) != tt.wantPanic {
					t.Error("Memory.getMap() did not panic as expected")
				}
			}()

			m.getMap(tt.args.entityType)
		})
	}
}
