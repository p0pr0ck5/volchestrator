package memory

import (
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/p0pr0ck5/volchestrator/server/client"
	"github.com/p0pr0ck5/volchestrator/server/notification"
)

func TestMemory_ReadClient(t *testing.T) {
	type fields struct {
		clientMap         clientMap
		notificationChMap map[string]chan *notification.Notification
		dataLocks         map[string]*sync.Mutex
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
				notificationChMap: map[string]chan *notification.Notification{
					"foo": make(chan *notification.Notification),
				},
				dataLocks: make(map[string]*sync.Mutex),
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
				notificationChMap: map[string]chan *notification.Notification{
					"foo": make(chan *notification.Notification),
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
				clientMap:         tt.fields.clientMap,
				notificationChMap: tt.fields.notificationChMap,
				dataLocks:         tt.fields.dataLocks,
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
		clientMap         clientMap
		notificationChMap map[string]chan *notification.Notification
		dataLocks         map[string]*sync.Mutex
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
				notificationChMap: map[string]chan *notification.Notification{
					"foo": make(chan *notification.Notification),
				},
				dataLocks: make(map[string]*sync.Mutex),
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
				clientMap:         tt.fields.clientMap,
				notificationChMap: tt.fields.notificationChMap,
				dataLocks:         tt.fields.dataLocks,
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
		clientMap         clientMap
		notificationChMap map[string]chan *notification.Notification
		dataLocks         map[string]*sync.Mutex
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
				notificationChMap: map[string]chan *notification.Notification{
					"foo": make(chan *notification.Notification),
				},
				dataLocks: make(map[string]*sync.Mutex),
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
				notificationChMap: map[string]chan *notification.Notification{
					"foo": make(chan *notification.Notification),
				},
				dataLocks: make(map[string]*sync.Mutex),
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
				clientMap:         tt.fields.clientMap,
				notificationChMap: tt.fields.notificationChMap,
				dataLocks:         tt.fields.dataLocks,
			}
			if err := m.CreateClient(tt.args.client); (err != nil) != tt.wantErr {
				t.Errorf("Memory.CreateClient() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMemory_UpdateClient(t *testing.T) {
	type fields struct {
		clientMap         map[string]*client.Client
		notificationChMap map[string]chan *notification.Notification
		dataLocks         map[string]*sync.Mutex
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
				notificationChMap: map[string]chan *notification.Notification{
					"foo": make(chan *notification.Notification),
				},
				dataLocks: make(map[string]*sync.Mutex),
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
				notificationChMap: map[string]chan *notification.Notification{
					"foo": make(chan *notification.Notification),
				},
				dataLocks: make(map[string]*sync.Mutex),
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
				clientMap:         tt.fields.clientMap,
				notificationChMap: tt.fields.notificationChMap,
				dataLocks:         tt.fields.dataLocks,
			}
			if err := m.UpdateClient(tt.args.client); (err != nil) != tt.wantErr {
				t.Errorf("Memory.UpdateClient() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMemory_DeleteClient(t *testing.T) {
	type fields struct {
		clientMap         clientMap
		notificationChMap map[string]chan *notification.Notification
		dataLocks         map[string]*sync.Mutex
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
				notificationChMap: map[string]chan *notification.Notification{
					"foo": make(chan *notification.Notification),
				},
				dataLocks: make(map[string]*sync.Mutex),
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
				clientMap:         tt.fields.clientMap,
				notificationChMap: tt.fields.notificationChMap,
				dataLocks:         tt.fields.dataLocks,
			}
			if err := m.DeleteClient(tt.args.client); (err != nil) != tt.wantErr {
				t.Errorf("Memory.DeleteClient() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
