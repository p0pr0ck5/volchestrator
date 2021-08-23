package mock

import (
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/p0pr0ck5/volchestrator/server/client"
	"github.com/p0pr0ck5/volchestrator/server/model"
	"github.com/p0pr0ck5/volchestrator/server/notification"
	"github.com/p0pr0ck5/volchestrator/server/volume"
)

func nowIsh() time.Time {
	t := time.Now()
	return t.Round(time.Hour)
}

func NowIsh() time.Time {
	return nowIsh()
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

		Model: model.Model{
			CreatedAt: nowIsh(),
		},
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

		Model: model.Model{
			CreatedAt: nowIsh(),
		},
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

func (m *MockBackend) crud(op string, entity model.Base) error {
	entityType := reflect.ValueOf(entity).Elem().Type().Name()

	fnName := op + entityType

	ff := reflect.ValueOf(m).MethodByName(fnName)

	if !ff.IsValid() {
		return fmt.Errorf("unsupported type %q", entityType)
	}

	res := ff.Call([]reflect.Value{
		reflect.ValueOf(entity),
	})

	if err, ok := res[0].Interface().(error); ok {
		return err
	} else {
		return nil
	}
}

func (m *MockBackend) Create(entity model.Base) error {
	return m.crud("Create", entity)
}

func (m *MockBackend) Read(entity model.Base) (model.Base, error) {
	entityType := reflect.ValueOf(entity).Elem().Type().Name()

	fnName := "Read" + entityType

	ff := reflect.ValueOf(m).MethodByName(fnName)

	if !ff.IsValid() {
		return nil, fmt.Errorf("unsupported type %q", entityType)
	}

	id := reflect.ValueOf(entity).Elem().FieldByName("ID")

	res := ff.Call([]reflect.Value{
		id,
	})

	if err, ok := res[1].Interface().(error); ok {
		return nil, err
	} else {
		return res[0].Interface().(model.Base), nil
	}
}

func (m *MockBackend) Update(entity model.Base) error {
	return m.crud("Update", entity)
}

func (m *MockBackend) Delete(entity model.Base) error {
	return m.crud("Delete", entity)
}

func (m *MockBackend) List(entityType string, entities *[]model.Base) error {
	switch entityType {
	case "Client":
		clients, err := m.ListClients()
		if err != nil {
			return err
		}

		for _, client := range clients {
			*entities = append(*entities, client)
		}

		return nil
	case "Volume":
		volumes, err := m.ListVolumes()
		if err != nil {
			return err
		}

		for _, volume := range volumes {
			*entities = append(*entities, volume)
		}

		return nil
	default:
		return fmt.Errorf("unsupported type %q", entityType)
	}
}
