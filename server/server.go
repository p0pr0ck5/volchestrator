package server

import (
	"time"

	multierror "github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"

	"github.com/p0pr0ck5/volchestrator/server/backend"
	"github.com/p0pr0ck5/volchestrator/server/config"
	"github.com/p0pr0ck5/volchestrator/svc"
)

type Server struct {
	svc.UnimplementedVolchestratorServer
	svc.UnimplementedVolchestratorAdminServer

	b *backend.Backend

	config *config.Config

	t Tokener

	shutdownCh chan struct{}
}

func NewServer(opts ...ServerOpt) (*Server, error) {
	s := &Server{
		config:     config.DefaultConfig(),
		t:          RandTokener{},
		shutdownCh: make(chan struct{}),
	}

	for _, opt := range opts {
		if err := opt(s); err != nil {
			return nil, errors.Wrap(err, "opt error")
		}
	}

	return s, nil
}

func (s *Server) Shutdown() {
	close(s.shutdownCh)
}

func (s *Server) PruneClients() error {
	var errs *multierror.Error

	clients, err := s.b.ListClients()
	if err != nil {
		errs = multierror.Append(err)
		return errs
	}

	for _, client := range clients {
		if time.Now().Sub(client.LastSeen) > time.Second*time.Duration(s.config.ClientTTL) {
			err := s.b.DeleteClient(client)
			if err != nil {
				errs = multierror.Append(err)
			}
		}
	}

	return errs.ErrorOrNil()
}
