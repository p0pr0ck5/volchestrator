package mock

import (
	"errors"
	"time"

	"github.com/p0pr0ck5/volchestrator/server/client"
	"github.com/p0pr0ck5/volchestrator/server/notification"
	"github.com/p0pr0ck5/volchestrator/server/volume"
)

func nowIsh() time.Time {
	t := time.Now()
	return t.Round(time.Hour)
}

type MockBackend struct{}

func NewMockBackend() *MockBackend {
	return &MockBackend{}
}

func (m *MockBackend) ReadClient(id string) (*client.Client, error) {
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

func (m *MockBackend) ListClients() ([]*client.Client, error) {
	return []*client.Client{
		{
			ID:         "foo",
			Registered: nowIsh(),
			LastSeen:   nowIsh(),
		},
	}, nil
}

func (m *MockBackend) CreateClient(c *client.Client) error {
	if c.ID == "bad" {
		return errors.New("error")
	}

	return nil
}

func (m *MockBackend) UpdateClient(c *client.Client) error {
	if c.ID == "bad" {
		return errors.New("error")
	}

	return nil
}

func (m *MockBackend) DeleteClient(c *client.Client) error {
	if c.ID == "bad" {
		return errors.New("error")
	}

	return nil
}

func (m *MockBackend) ReadVolume(id string) (*volume.Volume, error) {
	if id == "bad" {
		return nil, errors.New("error")
	}

	c := &volume.Volume{
		ID:     id,
		Region: "us-west-2",
		Tag:    "foo",
		Status: volume.Unavailable,
	}

	return c, nil
}

func (m *MockBackend) ListVolumes() ([]*volume.Volume, error) {
	return []*volume.Volume{
		{
			ID:     "foo",
			Region: "us-west-2",
			Tag:    "foo",
		},
	}, nil
}

func (m *MockBackend) CreateVolume(v *volume.Volume) error {
	if v.ID == "bad" {
		return errors.New("error")
	}

	return nil
}

func (m *MockBackend) UpdateVolume(v *volume.Volume) error {
	if v.ID == "bad" {
		return errors.New("error")
	}

	return nil
}

func (m *MockBackend) DeleteVolume(v *volume.Volume) error {
	if v.ID == "bad" {
		return errors.New("error")
	}

	return nil
}

func (m *MockBackend) GetNotifications(id string) <-chan *notification.Notification {
	if id == "bad" {
		return nil
	}

	ch := make(chan *notification.Notification)

	go func() {
		time.Sleep(time.Millisecond * 10)
		close(ch)
	}()

	return ch
}
