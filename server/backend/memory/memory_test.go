package memory

import (
	"reflect"
	"testing"

	"github.com/p0pr0ck5/volchestrator/server/client"
	"github.com/p0pr0ck5/volchestrator/server/volume"
)

func TestNewMemoryBackend(t *testing.T) {
	tests := []struct {
		name    string
		want    *Memory
		wantErr bool
	}{
		{
			"new memory backend",
			&Memory{
				clientMap: make(map[string]*client.Client),
				volumeMap: make(map[string]*volume.Volume),
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewMemoryBackend()
			if (err != nil) != tt.wantErr {
				t.Errorf("NewMemoryBackend() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMemoryBackend() = %v, want %v", got, tt.want)
			}
		})
	}
}
