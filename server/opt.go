package server

import (
	"github.com/p0pr0ck5/volchestrator/server/backend"
	"github.com/p0pr0ck5/volchestrator/server/backend/memory"
)

type ServerOpt func(*Server) error

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
