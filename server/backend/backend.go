package backend

import (
	"github.com/p0pr0ck5/volchestrator/server/backend/memory"
	"github.com/p0pr0ck5/volchestrator/server/backend/mock"
	"github.com/p0pr0ck5/volchestrator/server/client"
	"github.com/p0pr0ck5/volchestrator/server/notification"
	"github.com/p0pr0ck5/volchestrator/server/volume"
)

type backend interface {
	ReadClient(id string) (*client.Client, error)
	ListClients() ([]*client.Client, error)
	CreateClient(*client.Client) error
	UpdateClient(*client.Client) error
	DeleteClient(*client.Client) error

	ReadVolume(id string) (*volume.Volume, error)
	ListVolumes() ([]*volume.Volume, error)
	CreateVolume(*volume.Volume) error
	UpdateVolume(*volume.Volume) error
	DeleteVolume(*volume.Volume) error

	WriteNotification(*notification.Notification) error
	GetNotifications(string) (<-chan *notification.Notification, error)
}

type Backend struct {
	b backend
}

func NewBackend(opts ...BackendOpt) *Backend {
	b := &Backend{}

	for _, opt := range opts {
		opt(b)
	}

	return b
}

func NewMemoryBackend(opts ...BackendOpt) *Backend {
	m := memory.NewMemoryBackend()

	b := &Backend{
		b: m,
	}

	for _, opt := range opts {
		opt(b)
	}

	return b
}

func NewMockBackend(opts ...BackendOpt) *Backend {
	b := &Backend{
		b: mock.NewMockBackend(),
	}

	for _, opt := range opts {
		opt(b)
	}

	return b
}
