package server

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/p0pr0ck5/volchestrator/server/client"
	"github.com/p0pr0ck5/volchestrator/svc"
)

func (s *Server) Register(ctx context.Context, req *svc.RegisterRequest) (*svc.RegisterResponse, error) {
	if req.ClientId == "" {
		return nil, errors.New("empty client id")
	}

	client := &client.Client{
		ID:         req.ClientId,
		Token:      s.t.Generate(),
		Registered: time.Now(),
	}

	if err := s.b.CreateClient(client); err != nil {
		return nil, errors.Wrap(err, "create failed")
	}

	res := &svc.RegisterResponse{
		ClientId: client.ID,
		Token:    client.Token,
	}

	return res, nil
}

func (s *Server) Deregister(ctx context.Context, req *svc.DeregisterRequest) (*svc.DeregisterResponse, error) {
	if req.ClientId == "" {
		return nil, errors.New("empty client id")
	}

	client, err := s.b.ReadClient(req.ClientId)
	if err != nil {
		return nil, errors.Wrap(err, "get client")
	}

	if client.Token != req.Token {
		return nil, errors.New("invalid token")
	}

	if err := s.b.DeleteClient(client); err != nil {
		return nil, errors.Wrap(err, "delete failed")
	}

	return &svc.DeregisterResponse{}, nil
}

func (s *Server) Ping(ctx context.Context, req *svc.PingRequest) (*svc.PingResponse, error) {
	if req.ClientId == "" {
		return nil, errors.New("empty client id")
	}

	client, err := s.b.ReadClient(req.ClientId)
	if err != nil {
		return nil, errors.Wrap(err, "get client")
	}

	if client.Token != req.Token {
		return nil, errors.New("invalid token")
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

	client, err := s.b.ReadClient(req.ClientId)
	if err != nil {
		return errors.Wrap(err, "get client")
	}

	if client.Token != req.Token {
		return errors.New("invalid token")
	}

	ch, err := s.b.GetNotifications(req.ClientId)
	if err != nil {
		return errors.Wrap(err, "get notifications")
	}

	if ch == nil {
		return errors.New("no notifications channel")
	}

	for {
		select {
		case <-s.shutdownCh:
			return status.Error(codes.Unavailable, "shutting down")
		case notif := <-ch:
			select {
			case <-s.shutdownCh:
				return status.Error(codes.Unavailable, "shutting down")
			default:
			}
			if notif == nil {
				return status.Error(codes.Aborted, "channel closed")
			}

			stream.Send(&svc.WatchNotificationsResponse{
				Notification: toProto(notif).(*svc.Notification),
			})
		}
	}
}
