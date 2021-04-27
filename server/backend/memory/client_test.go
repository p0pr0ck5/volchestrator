package memory

import (
	"reflect"
	"testing"
	"time"

	"github.com/p0pr0ck5/volchestrator/server/client"
)

func TestMemory_ReadClient(t *testing.T) {
	type fields struct {
		clientMap clientMap
	}
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *client.Client
		wantErr bool
	}{
		{
			"read client",
			fields{
				clientMap: map[string]*client.Client{
					"foo": {
						ID: "foo",
					},
				},
			},
			args{
				"foo",
			},
			&client.Client{
				ID: "foo",
			},
			false,
		},
		{
			"read non-existent client",
			fields{
				clientMap: map[string]*client.Client{
					"foo": {
						ID: "foo",
					},
				},
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
				clientMap: tt.fields.clientMap,
			}
			got, err := m.ReadClient(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Memory.ReadClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Memory.ReadClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemory_ListClients(t *testing.T) {
	type fields struct {
		clientMap clientMap
	}
	tests := []struct {
		name    string
		fields  fields
		want    []*client.Client
		wantErr bool
	}{
		{
			"one client",
			fields{
				clientMap: map[string]*client.Client{
					"foo": {
						ID: "foo",
					},
				},
			},
			[]*client.Client{
				{
					ID: "foo",
				},
			},
			false,
		},
		{
			"two clients",
			fields{
				clientMap: map[string]*client.Client{
					"foo": {
						ID: "foo",
					},
					"bar": {
						ID: "bar",
					},
				},
			},
			[]*client.Client{
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
			"zero clients",
			fields{
				clientMap: map[string]*client.Client{},
			},
			[]*client.Client{},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Memory{
				clientMap: tt.fields.clientMap,
			}
			got, err := m.ListClients()
			if (err != nil) != tt.wantErr {
				t.Errorf("Memory.ListClients() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Memory.ListClients() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemory_CreateClient(t *testing.T) {
	type fields struct {
		clientMap clientMap
	}
	type args struct {
		client *client.Client
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
				clientMap: map[string]*client.Client{
					"foo": {
						ID: "foo",
					},
				},
			},
			args{
				client: &client.Client{
					ID: "bar",
				},
			},
			false,
		},
		{
			"create existing client",
			fields{
				clientMap: map[string]*client.Client{
					"foo": {
						ID: "foo",
					},
				},
			},
			args{
				client: &client.Client{
					ID: "foo",
				},
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Memory{
				clientMap: tt.fields.clientMap,
			}
			if err := m.CreateClient(tt.args.client); (err != nil) != tt.wantErr {
				t.Errorf("Memory.CreateClient() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMemory_UpdateClient(t *testing.T) {
	type fields struct {
		clientMap map[string]*client.Client
	}
	type args struct {
		client *client.Client
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
				clientMap: map[string]*client.Client{
					"foo": {
						ID: "foo",
					},
				},
			},
			args{
				client: &client.Client{
					ID: "foo",
				},
			},
			false,
		},
		{
			"update nonexistent client",
			fields{
				clientMap: map[string]*client.Client{
					"foo": {
						ID: "foo",
					},
				},
			},
			args{
				client: &client.Client{
					ID: "bar",
				},
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Memory{
				clientMap: tt.fields.clientMap,
			}
			if err := m.UpdateClient(tt.args.client); (err != nil) != tt.wantErr {
				t.Errorf("Memory.UpdateClient() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMemory_DeleteClient(t *testing.T) {
	type fields struct {
		clientMap clientMap
	}
	type args struct {
		client *client.Client
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
				clientMap: map[string]*client.Client{
					"foo": {
						ID:         "foo",
						Registered: time.Now(),
					},
				},
			},
			args{
				client: &client.Client{
					ID: "foo",
				},
			},
			false,
		},
		{
			"delete nonexistent client",
			fields{
				clientMap: map[string]*client.Client{
					"foo": {
						ID:         "foo",
						Registered: time.Now(),
					},
				},
			},
			args{
				client: &client.Client{
					ID: "bar",
				},
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Memory{
				clientMap: tt.fields.clientMap,
			}
			if err := m.DeleteClient(tt.args.client); (err != nil) != tt.wantErr {
				t.Errorf("Memory.DeleteClient() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
