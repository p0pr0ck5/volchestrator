package backend

import (
	"github.com/p0pr0ck5/volchestrator/server/backend/memory"
	"github.com/p0pr0ck5/volchestrator/server/backend/mock"
	"github.com/p0pr0ck5/volchestrator/server/client"
	"github.com/p0pr0ck5/volchestrator/server/volume"
)

type BackendOpt func(*Backend) error

func WithMemoryBackend(m *memory.Memory) BackendOpt {
	return func(b *Backend) error {
		b.b = m

		return nil
	}
}

func WithMockBackend(m *mock.MockBackend) BackendOpt {
	return func(b *Backend) error {
		b.b = m

		return nil
	}
}

func WithClients(clients []*client.Client) BackendOpt {
	return func(b *Backend) error {
		for _, c := range clients {
			err := b.Create(c)
			if err != nil {
				return err
			}
		}

		return nil
	}
}

func WithVolumes(volumes []*volume.Volume) BackendOpt {
	return func(b *Backend) error {
		for _, v := range volumes {
			err := b.Create(v)
			if err != nil {
				return err
			}
		}

		return nil
	}
}
