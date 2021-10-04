package server

import (
	"reflect"
	"testing"
	"time"

	"github.com/p0pr0ck5/volchestrator/server/backend"
	"github.com/p0pr0ck5/volchestrator/server/backend/mock"
	"github.com/p0pr0ck5/volchestrator/server/client"
	"github.com/p0pr0ck5/volchestrator/server/config"
	leaserequest "github.com/p0pr0ck5/volchestrator/server/lease_request"
	"github.com/p0pr0ck5/volchestrator/server/model"
	"github.com/p0pr0ck5/volchestrator/server/volume"
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
			Token:      "mock",
			Registered: mockNow,
			LastSeen:   mockNow,
		},
		{
			ID:         "bar",
			Token:      "mock",
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

			clients := []model.Base{}
			s.b.List(backend.Client, &clients)

			c := []*client.Client{}
			for _, cc := range clients {
				c = append(c, cc.(*client.Client))
			}
			if !reflect.DeepEqual(c, tt.want) {
				t.Errorf("Server.List() = %v, want %v", c, tt.want)
			}
		})
	}
}

func TestServer_PruneClientsReturn(t *testing.T) {
	ttl := time.Duration(30)
	mockThen := time.Now().Add(time.Second * -ttl)

	tests := []struct {
		name    string
		mocks   map[string]model.Base
		wantErr bool
	}{
		{
			"no error",
			map[string]model.Base{
				"Client": &client.Client{
					ID:         "foo",
					Registered: mockThen,
					LastSeen:   mockThen,
				},
			},
			false,
		},
		{
			"error during delete client",
			map[string]model.Base{
				"Client": &client.Client{
					ID:         "bad",
					Registered: mockThen,
					LastSeen:   mockThen,
				},
			},
			true,
		},
		{
			"error during list client",
			map[string]model.Base{
				"Client": &client.Client{
					Registered: mockThen,
					LastSeen:   mockThen,
				},
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &config.Config{
				ClientTTL: int(ttl * 2 / 3),
			}

			b := mock.NewMockBackend(mock.WithMocks(tt.mocks))

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

func TestServer_FindVolumes(t *testing.T) {
	mockVolumes := []*volume.Volume{
		{
			ID:     "foo",
			Region: "us-west-2",
			Tag:    "bar",
			Status: volume.Available,
		},
		{
			ID:     "bar",
			Region: "us-west-1",
			Tag:    "baz",
			Status: volume.Unavailable,
		},
		{
			ID:     "baz",
			Region: "us-west-2",
			Tag:    "bar",
			Status: volume.Available,
		},
		{
			ID:     "bat",
			Region: "us-west-2",
			Tag:    "foo",
			Status: volume.Available,
		},
	}

	type args struct {
		l *leaserequest.LeaseRequest
	}
	tests := []struct {
		name    string
		args    args
		want    []*volume.Volume
		wantErr bool
	}{
		{
			"one matching volume",
			args{
				&leaserequest.LeaseRequest{
					Region: "us-west-2",
					Tag:    "foo",
				},
			},
			[]*volume.Volume{
				{
					ID:     "bat",
					Region: "us-west-2",
					Tag:    "foo",
					Status: volume.Available,
				},
			},
			false,
		},
		{
			"two matching volumes",
			args{
				&leaserequest.LeaseRequest{
					Region: "us-west-2",
					Tag:    "bar",
				},
			},
			[]*volume.Volume{
				{
					ID:     "baz",
					Region: "us-west-2",
					Tag:    "bar",
					Status: volume.Available,
				},
				{
					ID:     "foo",
					Region: "us-west-2",
					Tag:    "bar",
					Status: volume.Available,
				},
			},
			false,
		},
		{
			"one matching volume (unavailable)",
			args{
				&leaserequest.LeaseRequest{
					Region: "us-west-1",
					Tag:    "baz",
				},
			},
			[]*volume.Volume{},
			false,
		},
		{
			"no matches",
			args{
				&leaserequest.LeaseRequest{
					Region: "us-west-3",
					Tag:    "foo",
				},
			},
			[]*volume.Volume{},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, _ := NewServer(WithBackend(backend.NewMemoryBackend(backend.WithVolumes(mockVolumes))))

			got, err := s.FindVolumes(tt.args.l)

			if (err != nil) != tt.wantErr {
				t.Errorf("Server.FindVolumes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want != nil {
				// mock
				for _, g := range got {
					g.CreatedAt = time.Time{}
					g.UpdatedAt = time.Time{}
					g.FSM = nil
				}
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Server.FindVolumes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServer_FindLeaseRequests(t *testing.T) {
	mockNow := time.Now()

	mockClients := []*client.Client{
		{
			ID:         "foo",
			Token:      "mock",
			Registered: mockNow,
			LastSeen:   mockNow,
		},
		{
			ID:         "bar",
			Token:      "mock",
			Registered: mockNow,
			LastSeen:   mockNow,
		},
		{
			ID:         "baz",
			Token:      "mock",
			Registered: mockNow,
			LastSeen:   mockNow,
		},
		{
			ID:         "bat",
			Token:      "mock",
			Registered: mockNow,
			LastSeen:   mockNow,
		},
	}

	mockLeaseRequests := []*leaserequest.LeaseRequest{
		{
			ID:       "foo",
			ClientID: "foo",
			Region:   "us-west-2",
			Tag:      "foo",
			Status:   leaserequest.Pending,
		},
		{
			ID:       "bar",
			ClientID: "bar",
			Region:   "us-west-2",
			Tag:      "bar",
			Status:   leaserequest.Fulfilling,
		},
		{
			ID:       "baz",
			ClientID: "baz",
			Region:   "us-west-2",
			Tag:      "foo",
			Status:   leaserequest.Pending,
		},
		{
			ID:       "bat",
			ClientID: "bat",
			Region:   "us-west-1",
			Tag:      "bat",
			Status:   leaserequest.Pending,
		},
	}
	type args struct {
		v *volume.Volume
	}
	tests := []struct {
		name    string
		args    args
		want    []*leaserequest.LeaseRequest
		wantErr bool
	}{
		{
			"one matching lease request",
			args{
				v: &volume.Volume{
					ID:     "foo",
					Region: "us-west-1",
					Tag:    "bat",
					Status: volume.Available,
				},
			},
			[]*leaserequest.LeaseRequest{
				{
					ID:       "bat",
					ClientID: "bat",
					Region:   "us-west-1",
					Tag:      "bat",
					Status:   leaserequest.Pending,
				},
			},
			false,
		},
		{
			"two matching lease requests",
			args{
				v: &volume.Volume{
					ID:     "foo",
					Region: "us-west-2",
					Tag:    "foo",
					Status: volume.Available,
				},
			},
			[]*leaserequest.LeaseRequest{
				{
					ID:       "baz",
					ClientID: "baz",
					Region:   "us-west-2",
					Tag:      "foo",
					Status:   leaserequest.Pending,
				},
				{
					ID:       "foo",
					ClientID: "foo",
					Region:   "us-west-2",
					Tag:      "foo",
					Status:   leaserequest.Pending,
				},
			},
			false,
		},
		{
			"one matching lease request (status)",
			args{
				v: &volume.Volume{
					ID:     "foo",
					Region: "us-west-2",
					Tag:    "bar",
					Status: volume.Available,
				},
			},
			[]*leaserequest.LeaseRequest{},
			false,
		},
		{
			"no matching lease requests",
			args{
				v: &volume.Volume{
					ID:     "foo",
					Region: "us-west-2",
					Tag:    "bib",
					Status: volume.Available,
				},
			},
			[]*leaserequest.LeaseRequest{},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := []model.Base{}
			for _, c := range mockClients {
				b = append(b, c)
			}
			for _, e := range mockLeaseRequests {
				b = append(b, e)
			}
			s, _ := NewServer(WithBackend(backend.NewMemoryBackend(backend.WithEntities(b))))

			got, err := s.FindLeaseRequests(tt.args.v)
			if (err != nil) != tt.wantErr {
				t.Errorf("Server.FindLeaseRequests() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want != nil {
				// mock
				for _, g := range got {
					g.CreatedAt = time.Time{}
					g.UpdatedAt = time.Time{}
					g.FSM = nil
				}
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Server.FindLeaseRequests() = %v, want %v", got, tt.want)
			}
		})
	}
}
