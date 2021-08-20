package server

import (
	"testing"

	"github.com/p0pr0ck5/volchestrator/server/backend/memory"
	"github.com/p0pr0ck5/volchestrator/server/client"
	"github.com/p0pr0ck5/volchestrator/server/model"
)

func TestWithMemoryBackend(t *testing.T) {
	m := memory.NewMemoryBackend()

	type args struct {
		b *memory.Memory
	}
	tests := []struct {
		name      string
		args      []args
		wantEqual bool
		wantErr   bool
	}{
		{
			"two different memory backends",
			[]args{
				{
					b: memory.NewMemoryBackend(),
				},
				{
					b: memory.NewMemoryBackend(),
				},
			},
			false,
			false,
		},
		{
			"the same memory backend",
			[]args{
				{
					b: m,
				},
				{
					b: m,
				},
			},
			true,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s1, err := NewServer(WithMemoryBackend(tt.args[0].b))

			if (err != nil) != tt.wantErr {
				t.Errorf("NewServer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			s2, err := NewServer(WithMemoryBackend(tt.args[1].b))

			if (err != nil) != tt.wantErr {
				t.Errorf("NewServer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			c := &client.Client{
				ID:    "foo",
				Token: "mock",
			}
			s1.b.Create(c)

			l1 := []model.Base{}
			l2 := []model.Base{}
			s1.b.List("client", &l1)
			s2.b.List("client", &l2)
			if (len(l1) == len(l2)) != tt.wantEqual {
				t.Errorf("s1.b == s2.b = %v, want %v", s1.b == s2.b, tt.wantEqual)
			}
		})
	}
}
