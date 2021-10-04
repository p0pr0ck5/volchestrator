package server

import (
	"time"

	multierror "github.com/hashicorp/go-multierror"

	"github.com/p0pr0ck5/volchestrator/server/backend"
	"github.com/p0pr0ck5/volchestrator/server/client"
	"github.com/p0pr0ck5/volchestrator/server/config"
	leaserequest "github.com/p0pr0ck5/volchestrator/server/lease_request"
	"github.com/p0pr0ck5/volchestrator/server/model"
	"github.com/p0pr0ck5/volchestrator/server/volume"
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
	if err := s.b.List(backend.Client, &clients); err != nil {
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

func (s *Server) FindVolumes(l *leaserequest.LeaseRequest) ([]*volume.Volume, error) {
	var volumes []model.Base
	if err := s.b.List(backend.Volume, &volumes); err != nil {
		return nil, err
	}

	res := []*volume.Volume{}
	for _, v := range volumes {
		v := v.(*volume.Volume)

		if v.Status != volume.Available {
			continue
		}

		if v.Tag == l.Tag && v.Region == l.Region {
			res = append(res, v)
		}
	}

	return res, nil
}

func (s *Server) FindLeaseRequests(v *volume.Volume) ([]*leaserequest.LeaseRequest, error) {
	var leaseRequests []model.Base
	if err := s.b.List(backend.LeaseRequest, &leaseRequests); err != nil {
		return nil, err
	}

	res := []*leaserequest.LeaseRequest{}
	for _, l := range leaseRequests {
		l := l.(*leaserequest.LeaseRequest)

		if l.Status != leaserequest.Pending {
			continue
		}

		if l.Tag == v.Tag && l.Region == v.Region {
			res = append(res, l)
		}
	}

	return res, nil
}
