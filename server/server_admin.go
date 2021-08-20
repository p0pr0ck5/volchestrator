package server

import (
	"context"

	"github.com/pkg/errors"

	"github.com/p0pr0ck5/volchestrator/server/client"
	"github.com/p0pr0ck5/volchestrator/server/model"
	"github.com/p0pr0ck5/volchestrator/server/volume"
	"github.com/p0pr0ck5/volchestrator/svc"
)

func (s *Server) GetClient(ctx context.Context, req *svc.GetClientRequest) (*svc.GetClientResponse, error) {
	if req.ClientId == "" {
		return nil, errors.New("empty client id")
	}

	client := &client.Client{
		ID: req.ClientId,
	}

	if err := s.b.Read(client); err != nil {
		return nil, err
	}

	res := &svc.GetClientResponse{
		Client: toProto(client).(*svc.Client),
	}

	return res, nil
}

func (s *Server) ListClients(ctx context.Context, req *svc.ListClientsRequest) (*svc.ListClientsResponse, error) {
	clients := []model.Base{}
	if err := s.b.List("client", &clients); err != nil {
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

	volume := &volume.Volume{
		ID: req.VolumeId,
	}

	if err := s.b.Read(volume); err != nil {
		return nil, err
	}

	res := &svc.GetVolumeResponse{
		Volume: toProto(volume).(*svc.Volume),
	}

	return res, nil
}

func (s *Server) ListVolumes(ctx context.Context, req *svc.ListVolumesRequest) (*svc.ListVolumesReponse, error) {
	volumes := []model.Base{}
	if err := s.b.List("volume", &volumes); err != nil {
		return nil, err
	}

	res := &svc.ListVolumesReponse{}

	for _, volume := range volumes {
		res.Volumes = append(res.Volumes, toProto(volume).(*svc.Volume))
	}

	return res, nil
}

func (s *Server) AddVolume(ctx context.Context, req *svc.AddVolumeRequest) (*svc.AddVolumeResponse, error) {
	v := toStruct(req.Volume).(*volume.Volume)

	if v.Status != volume.Available && v.Status != volume.Unavailable {
		return nil, errors.New("invalid status")
	}

	if err := s.b.Create(v); err != nil {
		return nil, errors.Wrap(err, "create failed")
	}

	return &svc.AddVolumeResponse{}, nil
}

func (s *Server) UpdateVolume(ctx context.Context, req *svc.UpdateVolumeRequest) (*svc.UpdateVolumeResponse, error) {
	v := toStruct(req.Volume).(*volume.Volume)

	if err := s.b.Update(v); err != nil {
		return nil, errors.Wrap(err, "update failed")
	}

	return &svc.UpdateVolumeResponse{}, nil
}

func (s *Server) DeleteVolume(ctx context.Context, req *svc.DeleteVolumeRequest) (*svc.DeleteVolumeResponse, error) {
	v := toStruct(req.Volume).(*volume.Volume)

	if err := s.b.Delete(v); err != nil {
		return nil, errors.Wrap(err, "delete failed")
	}

	return &svc.DeleteVolumeResponse{}, nil
}
