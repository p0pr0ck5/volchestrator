package server

import (
	"context"

	"github.com/pkg/errors"

	"github.com/p0pr0ck5/volchestrator/svc"
)

func (s *Server) GetClient(ctx context.Context, req *svc.GetClientRequest) (*svc.GetClientResponse, error) {
	if req.ClientId == "" {
		return nil, errors.New("empty client id")
	}

	client, err := s.b.ReadClient(req.ClientId)

	if err != nil {
		return nil, err
	}

	res := &svc.GetClientResponse{
		Client: toProto(client).(*svc.Client),
	}

	return res, nil
}

func (s *Server) ListClients(ctx context.Context, req *svc.ListClientsRequest) (*svc.ListClientsResponse, error) {
	clients, err := s.b.ListClients()
	if err != nil {
		return nil, err
	}

	res := &svc.ListClientsResponse{}

	for _, client := range clients {
		res.Clients = append(res.Clients, toProto(client).(*svc.Client))
	}

	return res, nil
}

func (s *Server) GetVolume(ctx context.Context, req *svc.GetVolumeRequest) (*svc.GetVolumeResponse, error) {
	if req.VolumeId == "" {
		return nil, errors.New("empty client id")
	}

	volume, err := s.b.ReadVolume(req.VolumeId)

	if err != nil {
		return nil, err
	}

	res := &svc.GetVolumeResponse{
		Volume: toProto(volume).(*svc.Volume),
	}

	return res, nil
}
