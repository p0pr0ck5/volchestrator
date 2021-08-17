package server

import (
	"github.com/p0pr0ck5/volchestrator/server/backend"
	"github.com/p0pr0ck5/volchestrator/server/backend/memory"
	"github.com/p0pr0ck5/volchestrator/server/config"
)

type ServerOpt func(*Server) error

func WithBackend(b *backend.Backend) ServerOpt {
	return func(s *Server) error {
		s.b = b

		return nil
	}
}

func WithNewMemoryBackend() ServerOpt {
	return func(s *Server) error {
		s.b = backend.NewMemoryBackend()

		return nil
	}
}

func WithMemoryBackend(b *memory.Memory) ServerOpt {
	return func(s *Server) error {
		s.b = backend.NewBackend(backend.WithMemoryBackend(b))

		return nil
	}
}

func WithMockBackend() ServerOpt {
	return func(s *Server) error {
		s.b = backend.NewMockBackend()

		return nil
	}
}

func WithConfig(c *config.Config) ServerOpt {
	return func(s *Server) error {
		s.config = c

		return nil
	}
}

func WithTokener(t Tokener) ServerOpt {
	return func(s *Server) error {
		s.tokener = t

		return nil
	}
}
