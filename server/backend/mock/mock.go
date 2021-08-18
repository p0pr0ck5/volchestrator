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

type MockBackend struct {
	ClientLister func() ([]*client.Client, error)
	VolumeLister func() ([]*volume.Volume, error)
}

func NewMockBackend() *MockBackend {
	return &MockBackend{}
}

func (m *MockBackend) ReadClient(id string) (*client.Client, error) {
	if id == "bad" {
		return nil, errors.New("error")
	}

	c := &client.Client{
		ID:         id,
		Token:      "mock",
		Registered: nowIsh(),
		LastSeen:   nowIsh(),
	}
	c.Init()

	return c, nil
}

func (m *MockBackend) ListClients() ([]*client.Client, error) {
	if m.ClientLister != nil {
		return m.ClientLister()
	}

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

	v := &volume.Volume{
		ID:     id,
		Region: "us-west-2",
		Tag:    "foo",
		Status: volume.Unavailable,
	}
	v.Init()

	return v, nil
}

func (m *MockBackend) ListVolumes() ([]*volume.Volume, error) {
	if m.VolumeLister != nil {
		return m.VolumeLister()
	}

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

func (m *MockBackend) WriteNotification(n *notification.Notification) error {
	return nil
}

func (m *MockBackend) GetNotifications(id string) (<-chan *notification.Notification, error) {
	if id == "bad" {
		return nil, nil
	}

	ch := make(chan *notification.Notification)

	go func() {
		time.Sleep(time.Millisecond * 10)
		close(ch)
	}()

	return ch, nil
}
