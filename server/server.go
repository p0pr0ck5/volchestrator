package server

import (
	"github.com/pkg/errors"

	"github.com/p0pr0ck5/volchestrator/server/backend"
	"github.com/p0pr0ck5/volchestrator/svc"
)

type Server struct {
	svc.UnimplementedVolchestratorServer
	svc.UnimplementedVolchestratorAdminServer

	b backend.Backend
}

func NewServer(opts ...ServerOpt) (*Server, error) {
	s := &Server{}

	for _, opt := range opts {
		err := opt(s)
		if err != nil {
			return nil, errors.Wrap(err, "opt error")
		}
	}

	return s, nil
}
