package server

import (
	"github.com/p0pr0ck5/volchestrator/server/backend/memory"
)

type ServerOpt func(*Server) error

func WithNewMemoryBackend() ServerOpt {
	return func(s *Server) error {
		b, err := memory.NewMemoryBackend()
		if err != nil {
			return err
		}

		s.b = b

		return nil
	}
}

func WithMemoryBackend(b *memory.Memory) ServerOpt {
	return func(s *Server) error {
		s.b = b

		return nil
	}
}
