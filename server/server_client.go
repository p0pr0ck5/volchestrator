package server

import (
	"context"
	"reflect"
	"time"

	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/p0pr0ck5/volchestrator/server/client"
	leaserequest "github.com/p0pr0ck5/volchestrator/server/lease_request"
	"github.com/p0pr0ck5/volchestrator/svc"
)

func (s *Server) authClient(req interface{}) (*client.Client, error) {
	clientID := reflect.ValueOf(req).Elem().FieldByName("ClientId").Interface().(string)

	if clientID == "" {
		return nil, errors.New("empty client id")
	}

	client := &client.Client{
		ID: clientID,
	}

	if err := s.b.Read(client); err != nil {
		return nil, errors.Wrap(err, "get client")
	}

	token := reflect.ValueOf(req).Elem().FieldByName("Token").Interface().(string)

	if client.Token != token {
		return nil, errors.New("invalid token")
	}

	return client, nil
}

func (s *Server) Register(ctx context.Context, req *svc.RegisterRequest) (*svc.RegisterResponse, error) {
	if req.ClientId == "" {
		return nil, errors.New("empty client id")
	}

	client := &client.Client{
		ID:         req.ClientId,
		Token:      s.tokener.Generate(),
		Registered: time.Now(),
	}

	if err := s.b.Create(client); err != nil {
		return nil, errors.Wrap(err, "create failed")
	}

	res := &svc.RegisterResponse{
		ClientId: client.ID,
		Token:    client.Token,
	}

	return res, nil
}

func (s *Server) Deregister(ctx context.Context, req *svc.DeregisterRequest) (*svc.DeregisterResponse, error) {
	var client *client.Client
	var err error

	if client, err = s.authClient(req); err != nil {
		return nil, err
	}

	if err := s.b.Delete(client); err != nil {
		return nil, errors.Wrap(err, "delete failed")
	}

	return &svc.DeregisterResponse{}, nil
}

func (s *Server) Ping(ctx context.Context, req *svc.PingRequest) (*svc.PingResponse, error) {
	var client *client.Client
	var err error

	if client, err = s.authClient(req); err != nil {
		return nil, err
	}

	if err := s.b.Update(client); err != nil {
		return nil, errors.Wrap(err, "update failed")
	}

	return &svc.PingResponse{}, nil
}

func (s *Server) WatchNotifications(req *svc.WatchNotificationsRequest, stream svc.Volchestrator_WatchNotificationsServer) error {
	var client *client.Client
	var err error

	if client, err = s.authClient(req); err != nil {
		return err
	}

	ch, err := s.b.GetNotifications(client.ID)
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

func (s *Server) RequestLease(ctx context.Context, req *svc.RequestLeaseRequest) (*svc.RequestLeaseResponse, error) {
	var client *client.Client
	var err error

	if client, err = s.authClient(req); err != nil {
		return nil, err
	}

	leaseReq := &leaserequest.LeaseRequest{
		ID:       req.LeaseRequestId,
		ClientID: client.ID,
		Region:   req.Region,
		Tag:      req.Tag,
		Status:   leaserequest.Pending,
	}

	if err := s.b.Create(leaseReq); err != nil {
		return nil, err
	}

	return &svc.RequestLeaseResponse{}, nil
}
