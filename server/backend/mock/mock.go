package mock

import (
	"errors"
	"reflect"
	"time"

	"github.com/p0pr0ck5/volchestrator/server/client"
	"github.com/p0pr0ck5/volchestrator/server/lease"
	leaserequest "github.com/p0pr0ck5/volchestrator/server/lease_request"
	"github.com/p0pr0ck5/volchestrator/server/model"
	"github.com/p0pr0ck5/volchestrator/server/notification"
	"github.com/p0pr0ck5/volchestrator/server/volume"
)

type BackendOpt func(*MockBackend) error

func WithMocks(mocks map[string]model.Base) BackendOpt {
	return func(m *MockBackend) error {
		m.mocks = mocks

		return nil
	}
}

var mocks = map[string]model.Base{
	"Client": &client.Client{
		ID:         "foo",
		Token:      "mock",
		Registered: nowIsh(),
		LastSeen:   nowIsh(),

		Model: model.Model{
			CreatedAt: nowIsh(),
		},
	},
	"Volume": &volume.Volume{
		ID:     "foo",
		Region: "us-west-2",
		Tag:    "foo",
		Status: volume.Unavailable,

		Model: model.Model{
			CreatedAt: nowIsh(),
		},
	},
	"LeaseRequest": &leaserequest.LeaseRequest{
		ID:       "foo",
		ClientID: "foo",
		Region:   "us-west-2",
		Tag:      "foo",
		Status:   leaserequest.Pending,

		Model: model.Model{
			CreatedAt: nowIsh(),
		},
	},
	"Lease": &lease.Lease{
		ID:       "foo",
		ClientID: "foo",
		VolumeID: "foo",
		Status:   lease.Active,

		Model: model.Model{
			CreatedAt: nowIsh(),
		},
	},
}

func nowIsh() time.Time {
	t := time.Now()
	return t.Round(time.Hour)
}

func NowIsh() time.Time {
	return nowIsh()
}

type MockBackend struct {
	mocks map[string]model.Base
}

func NewMockBackend(opts ...BackendOpt) *MockBackend {
	m := &MockBackend{
		mocks: mocks,
	}

	for _, opt := range opts {
		opt(m)
	}

	return m
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
	id := entity.Identifier()

	if id == "bad" {
		return errors.New("error")
	}

	if _, ok := m.mocks[entityType]; !ok {
		return errors.New("unsupported")
	}

	return nil
}

func (m *MockBackend) Create(entity model.Base) error {
	return m.crud("Create", entity)
}

func (m *MockBackend) Read(entity model.Base) (model.Base, error) {
	entityType := reflect.ValueOf(entity).Elem().Type().Name()
	id := entity.Identifier()

	if id == "" {
		return nil, errors.New("missing")
	}

	if id == "bad" {
		return nil, errors.New("error")
	}

	e := m.mocks[entityType]
	if e == nil {
		return nil, errors.New("unsupported")
	}
	reflect.ValueOf(e).Elem().FieldByName("ID").Set(reflect.ValueOf(id))

	e.Init(model.WithSM(m.BuildSMMap()[entityType]))

	return e, nil
}

func (m *MockBackend) Update(entity model.Base) error {
	return m.crud("Update", entity)
}

func (m *MockBackend) Delete(entity model.Base) error {
	return m.crud("Delete", entity)
}

func (m *MockBackend) List(entityType string, entities *[]model.Base) error {
	e, ok := m.mocks[entityType]
	if !ok {
		return errors.New("unsupported")
	}
	e = e.Clone()
	e.Init(model.WithSM(m.BuildSMMap()[entityType]))

	*entities = append(*entities, e)

	return nil
}

func (m *MockBackend) Find(entityType, fieldName, id string) []model.Base {
	e, ok := m.mocks[entityType]
	if !ok {
		return []model.Base{}
	}
	f := reflect.ValueOf(e).Elem().FieldByName(fieldName).Interface().(string)
	if f == id {
		e = e.Clone()
		e.Init(model.WithSM(m.BuildSMMap()[entityType]))
		return []model.Base{e}
	}
	return nil
}
