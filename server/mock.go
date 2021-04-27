package server

import (
	"time"

	"github.com/pkg/errors"

	"github.com/p0pr0ck5/volchestrator/server/client"
	"github.com/p0pr0ck5/volchestrator/server/volume"
)

func nowIsh() time.Time {
	t := time.Now()
	return t.Round(time.Hour)
}

type mockBackend struct{}

func (m *mockBackend) ReadClient(id string) (*client.Client, error) {
	if id == "bad" {
		return nil, errors.New("error")
	}

	c := &client.Client{
		ID:         id,
		Registered: nowIsh(),
		LastSeen:   nowIsh(),
	}

	return c, nil
}

func (m *mockBackend) ListClients() ([]*client.Client, error) {
	return []*client.Client{
		{
			ID:         "foo",
			Registered: nowIsh(),
			LastSeen:   nowIsh(),
		},
	}, nil
}

func (m *mockBackend) CreateClient(c *client.Client) error {
	if c.ID == "bad" {
		return errors.New("error")
	}

	return nil
}

func (m *mockBackend) UpdateClient(c *client.Client) error {
	if c.ID == "bad" {
		return errors.New("error")
	}

	return nil
}

func (m *mockBackend) DeleteClient(c *client.Client) error {
	if c.ID == "bad" {
		return errors.New("error")
	}

	return nil
}

func (m *mockBackend) ReadVolume(id string) (*volume.Volume, error) {
	if id == "bad" {
		return nil, errors.New("error")
	}

	c := &volume.Volume{
		ID:     id,
		Region: "us-west-2",
		Tag:    "foo",
		Status: volume.Available,
	}

	return c, nil
}

func (m *mockBackend) ListVolumes() ([]*volume.Volume, error) {
	return []*volume.Volume{
		{
			ID:     "foo",
			Region: "us-west-2",
			Tag:    "foo",
		},
	}, nil
}

func (m *mockBackend) CreateVolume(c *volume.Volume) error {
	if c.ID == "bad" {
		return errors.New("error")
	}

	return nil
}

func (m *mockBackend) UpdateVolume(c *volume.Volume) error {
	if c.ID == "bad" {
		return errors.New("error")
	}

	return nil
}

func (m *mockBackend) DeleteVolume(c *volume.Volume) error {
	if c.ID == "bad" {
		return errors.New("error")
	}

	return nil
}

func withMockBackend() ServerOpt {
	return func(s *Server) error {
		s.b = &mockBackend{}
		return nil
	}
}
