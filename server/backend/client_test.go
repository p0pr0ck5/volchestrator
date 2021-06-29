package backend

import (
	"reflect"
	"testing"
	"time"

	"github.com/p0pr0ck5/volchestrator/server/client"
)

func TestBackend_UpdateClient(t *testing.T) {
	mockNow := time.Now()
	mockThen := time.Now().Add(time.Second * -30)

	mockClients := []*client.Client{
		{
			ID:         "foo",
			Registered: mockThen,
			LastSeen:   mockThen,
		},
		{
			ID:         "bar",
			Registered: mockThen,
			LastSeen:   mockThen,
		},
	}

	type fields struct {
		b backend
	}
	type args struct {
		c *client.Client
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *client.Client
		wantErr bool
	}{
		{
			"full update",
			fields{
				NewMemoryBackend(WithClients(mockClients)),
			},
			args{
				c: &client.Client{
					ID:         "foo",
					Registered: mockNow,
					LastSeen:   mockNow,
				},
			},
			&client.Client{
				ID:         "foo",
				Registered: mockNow,
				LastSeen:   mockNow,
			},
			false,
		},
		{
			"partial update",
			fields{
				NewMemoryBackend(WithClients(mockClients)),
			},
			args{
				c: &client.Client{
					ID:       "foo",
					LastSeen: mockNow,
				},
			},
			&client.Client{
				ID:         "foo",
				Registered: mockThen,
				LastSeen:   mockNow,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Backend{
				b: tt.fields.b,
			}
			if err := b.UpdateClient(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("Backend.UpdateClient() error = %v, wantErr %v", err, tt.wantErr)
			}

			c, _ := b.ReadClient(tt.want.ID)
			c.UpdatedAt = time.Time{}
			c.CreatedAt = time.Time{}

			if !reflect.DeepEqual(c, tt.want) {
				t.Errorf("Backend.UpdateClient() = %v, want %v", c, tt.want)
			}
		})
	}
}
