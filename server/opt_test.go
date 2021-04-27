package server

import (
	"reflect"
	"testing"

	"github.com/p0pr0ck5/volchestrator/server/backend/memory"
)

func TestWithNewMemoryBackend(t *testing.T) {
	tests := []struct {
		name    string
		want    string
		wantErr bool
	}{
		{
			"new memory backend",
			"*memory.Memory",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := NewServer(WithNewMemoryBackend())

			if (err != nil) != tt.wantErr {
				t.Errorf("NewServer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if reflect.TypeOf(s.b).String() != tt.want {
				t.Errorf("NewServer(WithNewMemoryBackend()).b = %v, want %v", reflect.TypeOf(s.b).String(), tt.want)
			}
		})
	}
}

func TestWithMemoryBackend(t *testing.T) {
	m, _ := memory.NewMemoryBackend()

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
					b: &memory.Memory{},
				},
				{
					b: &memory.Memory{},
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

			if (s1.b == s2.b) != tt.wantEqual {
				t.Errorf("s1.b == s2.b = %v, want %v", s1.b == s2.b, tt.wantEqual)
			}
		})
	}
}
