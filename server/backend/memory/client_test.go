package memory

import (
	"reflect"
	"sync"
	"testing"

	"github.com/p0pr0ck5/volchestrator/server/client"
)

func TestMemory_ListClients(t *testing.T) {
	type fields struct {
		ClientMap       ClientMap
		notificationMap map[string]*ChQueue
		dataLocks       map[string]*sync.Mutex
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
				ClientMap: map[string]*client.Client{
					"foo": {
						ID: "foo",
					},
				},
				notificationMap: map[string]*ChQueue{
					"foo": MustNewChQueue(),
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
				ClientMap: map[string]*client.Client{
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
					ID: "bar",
				},
				{
					ID: "foo",
				},
			},
			false,
		},
		{
			"zero clients",
			fields{
				ClientMap: map[string]*client.Client{},
			},
			[]*client.Client{},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Memory{
				ClientMap:       tt.fields.ClientMap,
				notificationMap: tt.fields.notificationMap,
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
