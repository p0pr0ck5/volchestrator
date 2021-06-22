package server

import (
	"reflect"
	"testing"
	"time"

	"github.com/p0pr0ck5/volchestrator/server/backend"
	"github.com/p0pr0ck5/volchestrator/server/backend/mock"
	"github.com/p0pr0ck5/volchestrator/server/client"
	"github.com/p0pr0ck5/volchestrator/server/config"
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

func TestServer_PruneClients(t *testing.T) {
	ttl := time.Duration(30)

	mockNow := time.Now()
	mockThen := time.Now().Add(time.Second * -ttl)

	mockClients := []*client.Client{
		{
			ID:         "foo",
			Registered: mockNow,
			LastSeen:   mockNow,
		},
		{
			ID:         "bar",
			Registered: mockThen,
			LastSeen:   mockThen,
		},
	}

	tests := []struct {
		name    string
		clients []*client.Client
		want    []*client.Client
	}{
		{
			"prune one client",
			mockClients,
			[]*client.Client{
				mockClients[0],
			},
		},
		{
			"prune multiple clients",
			[]*client.Client{
				mockClients[0],
				mockClients[1],
				{
					ID:         "baz",
					Registered: mockThen,
					LastSeen:   mockThen,
				},
			},
			[]*client.Client{
				mockClients[0],
			},
		},
		{
			"prune no clients",
			[]*client.Client{
				mockClients[0],
			},
			[]*client.Client{
				mockClients[0],
			},
		},
		{
			"prune empty client list",
			[]*client.Client{},
			[]*client.Client{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &config.Config{
				ClientTTL: int(ttl * 2 / 3),
			}

			s, _ := NewServer(
				WithConfig(config),
				WithBackend(backend.NewMemoryBackend(backend.WithClients(tt.clients))),
			)

			if err := s.PruneClients(); err != nil {
				t.Errorf("Server.PruneClients() error = %v", err)
			}

			clients, _ := s.b.ListClients()
			if !reflect.DeepEqual(clients, tt.want) {
				t.Errorf("Server.ListClients() = %v, want %v", clients, tt.want)
			}
		})
	}
}

func TestServer_PruneClientsReturn(t *testing.T) {
	ttl := time.Duration(30)

	mockThen := time.Now().Add(time.Second * -ttl)

	lister := func(s string) func() ([]*client.Client, error) {
		return func() ([]*client.Client, error) {
			if s == "" {
				return nil, errors.New("bad")
			}

			return []*client.Client{
				{
					ID:         s,
					Registered: mockThen,
					LastSeen:   mockThen,
				},
			}, nil
		}
	}

	tests := []struct {
		name    string
		lister  func() ([]*client.Client, error)
		wantErr bool
	}{
		{
			"no error",
			lister("foo"),
			false,
		},
		{
			"error during delete client",
			lister("bad"),
			true,
		},
		{
			"error during list client",
			lister(""),
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &config.Config{
				ClientTTL: int(ttl * 2 / 3),
			}

			b := mock.NewMockBackend()
			b.ClientLister = tt.lister

			s, _ := NewServer(
				WithConfig(config),
				WithBackend(backend.NewMockBackend(backend.WithMockBackend(b))),
			)

			if err := s.PruneClients(); (err != nil) != tt.wantErr {
				t.Errorf("s.PruneClients() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
