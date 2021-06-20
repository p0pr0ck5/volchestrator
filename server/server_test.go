package server

import (
	"testing"

	"github.com/pkg/errors"
)

func mockErrorOpt() ServerOpt {
	return func(s *Server) error {
		return errors.New("error")
	}
}

func TestNewServer(t *testing.T) {
	tests := []struct {
		name    string
		opts    []ServerOpt
		wantErr bool
	}{
		{
			"empty server",
			nil,
			false,
		},
		{
			"with memory backend",
			[]ServerOpt{WithNewMemoryBackend()},
			false,
		},
		{
			"with error opt",
			[]ServerOpt{mockErrorOpt()},
			true,
		},
		{
			"with memory backend and error opt",
			[]ServerOpt{WithNewMemoryBackend(), mockErrorOpt()},
			true,
		},
		{
			"with error opt and memory backend",
			[]ServerOpt{mockErrorOpt(), WithNewMemoryBackend()},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewServer(tt.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewServer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestServer_Shutdown(t *testing.T) {
	tests := []struct {
		name     string
		shutdown bool
		want     bool
	}{
		{
			"shutdownCh closed following shutdown",
			true,
			false,
		},
		{
			"shutdownCh open",
			false,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, _ := NewServer()

			if tt.shutdown {
				s.Shutdown()
			}

			ok := true

			select {
			case _, ok = <-s.shutdownCh:
			default:
			}

			if ok != tt.want {
				t.Errorf("Shutdown() shutdownCh = %v, want %v", ok, tt.want)
			}
		})
	}
}
