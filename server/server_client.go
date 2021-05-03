package server

import (
	"context"
	"time"

	"github.com/pkg/errors"

	"github.com/p0pr0ck5/volchestrator/server/client"
	"github.com/p0pr0ck5/volchestrator/svc"
)

func (s *Server) Register(ctx context.Context, req *svc.RegisterRequest) (*svc.RegisterResponse, error) {
	if req.ClientId == "" {
		return nil, errors.New("empty client id")
	}

	client := &client.Client{
		ID:         req.ClientId,
		Registered: time.Now(),
	}

	if err := s.b.CreateClient(client); err != nil {
		return nil, errors.Wrap(err, "create failed")
	}

	return &svc.RegisterResponse{}, nil
}

func (s *Server) Ping(ctx context.Context, req *svc.PingRequest) (*svc.PingResponse, error) {
	if req.ClientId == "" {
		return nil, errors.New("empty client id")
	}

	client := &client.Client{
		ID:       req.ClientId,
		LastSeen: time.Now(),
	}

	if err := s.b.UpdateClient(client); err != nil {
		return nil, errors.Wrap(err, "update failed")
	}

	return &svc.PingResponse{}, nil
}

func (s *Server) WatchNotifications(req *svc.WatchNotificationsRequest, stream svc.Volchestrator_WatchNotificationsServer) error {
	if req.ClientId == "" {
		return errors.New("empty client id")
	}

	ch := s.b.GetNotifications(req.ClientId)
	if ch == nil {
		return errors.New("no notifications channel")
	}

	for notif := range ch {
		stream.Send(&svc.WatchNotificationsResponse{
			Notification: toProto(notif).(*svc.Notification),
		})
	}

	return nil
}
