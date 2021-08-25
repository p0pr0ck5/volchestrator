package backend

import (
	"github.com/p0pr0ck5/volchestrator/server/backend/memory"
	"github.com/p0pr0ck5/volchestrator/server/backend/mock"
	"github.com/p0pr0ck5/volchestrator/server/model"
	"github.com/p0pr0ck5/volchestrator/server/notification"
)

type backend interface {
	Create(model.Base) error
	Read(model.Base) (model.Base, error)
	Update(model.Base) error
	Delete(model.Base) error
	List(string, *[]model.Base) error
	Find(string, string, string) []model.Base

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
