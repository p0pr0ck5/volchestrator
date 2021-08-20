package server

import (
	"time"

	multierror "github.com/hashicorp/go-multierror"

	"github.com/p0pr0ck5/volchestrator/server/backend"
	"github.com/p0pr0ck5/volchestrator/server/client"
	"github.com/p0pr0ck5/volchestrator/server/config"
	"github.com/p0pr0ck5/volchestrator/server/model"
	"github.com/p0pr0ck5/volchestrator/svc"
)

type Server struct {
	svc.UnimplementedVolchestratorServer
	svc.UnimplementedVolchestratorAdminServer

	b *backend.Backend

	config *config.Config

	tokener Tokener

	shutdownCh chan struct{}
}

func NewServer(opts ...ServerOpt) (*Server, error) {
	s := &Server{
		config:     config.DefaultConfig(),
		tokener:    RandTokener{},
		shutdownCh: make(chan struct{}),
	}

	var errs *multierror.Error

	for _, opt := range opts {
		if err := opt(s); err != nil {
			errs = multierror.Append(err)
		}
	}

	return s, errs.ErrorOrNil()
}

func (s *Server) Shutdown() {
	close(s.shutdownCh)
}

func (s *Server) PruneClients() error {
	var errs *multierror.Error

	clients := []model.Base{}
	if err := s.b.List("client", &clients); err != nil {
		errs = multierror.Append(err)
		return errs
	}

	for _, cc := range clients {
		c := cc.(*client.Client)
		if time.Since(c.LastSeen) > time.Second*time.Duration(s.config.ClientTTL) {
			err := s.b.Delete(cc)
			if err != nil {
				errs = multierror.Append(err)
			}
		}
	}

	return errs.ErrorOrNil()
}
