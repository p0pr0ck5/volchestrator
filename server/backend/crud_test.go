package backend

import (
	"reflect"
	"testing"

	"github.com/p0pr0ck5/volchestrator/fsm"
	"github.com/p0pr0ck5/volchestrator/server/backend/mock"
	"github.com/p0pr0ck5/volchestrator/server/client"
	leaserequest "github.com/p0pr0ck5/volchestrator/server/lease_request"
	"github.com/p0pr0ck5/volchestrator/server/model"
	"github.com/p0pr0ck5/volchestrator/server/volume"
)

type mockBase struct{}

func (m *mockBase) Init() {}

func (m *mockBase) Identifier() string {
	return ""
}

func (m *mockBase) Validate() error {
	return nil
}

func (m *mockBase) ValidateTransition(i model.Base) error {
	return nil
}

func (m *mockBase) SetStatus(s string) {}

func (m *mockBase) F() *fsm.FSM {
	return nil
}

func (v *mockBase) Clone() model.Base {
	vv := *v
	return &vv
}

func TestBackend_Create(t *testing.T) {
	type fields struct {
		b backend
	}
	type args struct {
		entity model.Base
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"valid client",
			fields{
				b: mock.NewMockBackend(),
			},
			args{
				entity: &client.Client{
					ID:    "foo",
					Token: "mock",
				},
			},
			false,
		},
		{
			"invalid client - missing id",
			fields{
				b: mock.NewMockBackend(),
			},
			args{
				entity: &client.Client{
					Token: "mock",
				},
			},
			true,
		},
		{
			"invalid client - missing token",
			fields{
				b: mock.NewMockBackend(),
			},
			args{
				entity: &client.Client{
					ID: "foo",
				},
			},
			true,
		},
		{
			"invalid client - bad id",
			fields{
				b: mock.NewMockBackend(),
			},
			args{
				entity: &client.Client{
					ID: "bad",
				},
			},
			true,
		},
		{
			"valid volume",
			fields{
				b: mock.NewMockBackend(),
			},
			args{
				entity: &volume.Volume{
					ID:     "foo",
					Region: "us-west-2",
					Tag:    "bar",
					Status: volume.Available,
				},
			},
			false,
		},
		{
			"invalid volume - missing id",
			fields{
				b: mock.NewMockBackend(),
			},
			args{
				entity: &volume.Volume{
					Region: "us-west-2",
					Tag:    "bar",
				},
			},
			true,
		},
		{
			"invalid volume - missing region",
			fields{
				b: mock.NewMockBackend(),
			},
			args{
				entity: &volume.Volume{
					ID:  "foo",
					Tag: "bar",
				},
			},
			true,
		},
		{
			"invalid volume - missing tag",
			fields{
				b: mock.NewMockBackend(),
			},
			args{
				entity: &volume.Volume{
					ID:     "foo",
					Region: "us-west-2",
				},
			},
			true,
		},
		{
			"invalid volume - bad id",
			fields{
				b: mock.NewMockBackend(),
			},
			args{
				entity: &volume.Volume{
					ID:     "bad",
					Region: "us-west-2",
					Tag:    "bar",
				},
			},
			true,
		},
		{
			"valid lease request",
			fields{
				b: mock.NewMockBackend(),
			},
			args{
				entity: &leaserequest.LeaseRequest{
					ID:       "foo",
					ClientID: "foo",
					Region:   "us-west-2",
					Tag:      "foo",
					Status:   leaserequest.Pending,
				},
			},
			false,
		},
		{
			"invalid lease request - missing id",
			fields{
				b: mock.NewMockBackend(),
			},
			args{
				entity: &leaserequest.LeaseRequest{
					ClientID: "foo",
					Region:   "us-west-2",
					Tag:      "foo",
					Status:   leaserequest.Pending,
				},
			},
			true,
		},
		{
			"invalid lease request - missing client id",
			fields{
				b: mock.NewMockBackend(),
			},
			args{
				entity: &leaserequest.LeaseRequest{
					ID:     "foo",
					Region: "us-west-2",
					Tag:    "foo",
					Status: leaserequest.Pending,
				},
			},
			true,
		},
		{
			"invalid lease request - missing region",
			fields{
				b: mock.NewMockBackend(),
			},
			args{
				entity: &leaserequest.LeaseRequest{
					ID:       "foo",
					ClientID: "foo",
					Tag:      "foo",
					Status:   leaserequest.Pending,
				},
			},
			true,
		},
		{
			"invalid lease request - missing tag",
			fields{
				b: mock.NewMockBackend(),
			},
			args{
				entity: &leaserequest.LeaseRequest{
					ID:       "foo",
					ClientID: "foo",
					Region:   "us-west-2",
					Status:   leaserequest.Pending,
				},
			},
			true,
		},
		{
			"invalid lease request - missing status",
			fields{
				b: mock.NewMockBackend(),
			},
			args{
				entity: &leaserequest.LeaseRequest{
					ID:       "foo",
					ClientID: "foo",
					Region:   "us-west-2",
					Tag:      "foo",
				},
			},
			true,
		},
		{
			"invalid lease request - bad client reference",
			fields{
				b: mock.NewMockBackend(),
			},
			args{
				entity: &leaserequest.LeaseRequest{
					ID:       "foo",
					ClientID: "bad",
					Region:   "us-west-2",
					Tag:      "foo",
					Status:   leaserequest.Pending,
				},
			},
			true,
		},
		{
			"invalid lease request - nonexistent client reference",
			fields{
				b: mock.NewMockBackend(),
			},
			args{
				entity: &leaserequest.LeaseRequest{
					ID:       "foo",
					ClientID: "bar",
					Region:   "us-west-2",
					Tag:      "foo",
					Status:   leaserequest.Pending,
				},
			},
			true,
		},
		{
			"unsupported",
			fields{
				b: mock.NewMockBackend(),
			},
			args{
				entity: &mockBase{},
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Backend{
				b: tt.fields.b,
			}
			if err := b.Create(tt.args.entity); (err != nil) != tt.wantErr {
				t.Errorf("Backend.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBackend_Read(t *testing.T) {
	type fields struct {
		b backend
	}
	type args struct {
		entity model.Base
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    model.Base
		wantErr bool
	}{
		{
			"valid client",
			fields{
				b: mock.NewMockBackend(),
			},
			args{
				&client.Client{
					ID: "foo",
				},
			},
			&client.Client{
				ID:         "foo",
				Token:      "mock",
				Registered: mock.NowIsh(),
				LastSeen:   mock.NowIsh(),

				Model: model.Model{
					CreatedAt: mock.NowIsh(),
				},
			},
			false,
		},
		{
			"invalid client",
			fields{
				b: mock.NewMockBackend(),
			},
			args{
				&client.Client{
					ID: "bad",
				},
			},
			nil,
			true,
		},
		{
			"valid volume",
			fields{
				b: mock.NewMockBackend(),
			},
			args{
				&volume.Volume{
					ID: "foo",
				},
			},
			&volume.Volume{
				ID:     "foo",
				Region: "us-west-2",
				Tag:    "foo",
				Status: volume.Unavailable,

				Model: model.Model{
					CreatedAt: mock.NowIsh(),
				},
			},
			false,
		},
		{
			"invalid volume",
			fields{
				b: mock.NewMockBackend(),
			},
			args{
				&volume.Volume{
					ID: "bad",
				},
			},
			nil,
			true,
		},
		{
			"valid lease request",
			fields{
				b: mock.NewMockBackend(),
			},
			args{
				&leaserequest.LeaseRequest{
					ID: "foo",
				},
			},
			&leaserequest.LeaseRequest{
				ID:       "foo",
				ClientID: "foo",
				Region:   "us-west-2",
				Tag:      "foo",
				Status:   leaserequest.Pending,

				Model: model.Model{
					CreatedAt: mock.NowIsh(),
				},
			},
			false,
		},
		{
			"invalid lease request",
			fields{
				b: mock.NewMockBackend(),
			},
			args{
				&leaserequest.LeaseRequest{
					ID: "bad",
				},
			},
			nil,
			true,
		},
		{
			"unsupported",
			fields{
				b: mock.NewMockBackend(),
			},
			args{
				entity: &mockBase{},
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Backend{
				b: tt.fields.b,
			}
			e := tt.args.entity
			if err := b.Read(e); (err != nil) != tt.wantErr {
				t.Errorf("Backend.Read() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.want != nil {
				// mock
				//
				// e.FSM = nil
				f := reflect.ValueOf(e).Elem().FieldByName("FSM")
				f.Set(reflect.Zero(f.Type()))

				if !reflect.DeepEqual(e, tt.want) {
					t.Errorf("Backend.Read() = %v, want %v", e, tt.want)
				}
			}
		})
	}
}

func TestBackend_Update(t *testing.T) {
	type fields struct {
		b backend
	}
	type args struct {
		entity model.Base
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    model.Base
		wantErr bool
	}{
		{
			"valid client - no-op",
			fields{
				b: mock.NewMockBackend(),
			},
			args{
				entity: &client.Client{
					ID: "foo",
				},
			},
			&client.Client{
				ID:         "foo",
				Token:      "mock",
				Registered: mock.NowIsh(),
				LastSeen:   mock.NowIsh(),
			},
			false,
		},
		{
			"valid client - new token",
			fields{
				b: mock.NewMockBackend(),
			},
			args{
				entity: &client.Client{
					ID:    "foo",
					Token: "newtoken",
				},
			},
			&client.Client{
				ID:         "foo",
				Token:      "newtoken",
				Registered: mock.NowIsh(),
				LastSeen:   mock.NowIsh(),
			},
			false,
		},
		{
			"valid client - Deleting status",
			fields{
				b: mock.NewMockBackend(),
			},
			args{
				entity: &client.Client{
					ID:     "foo",
					Status: client.Deleting,
				},
			},
			&client.Client{
				ID:         "foo",
				Token:      "mock",
				Registered: mock.NowIsh(),
				LastSeen:   mock.NowIsh(),
				Status:     client.Deleting,
			},
			false,
		},
		{
			"invalid client",
			fields{
				b: mock.NewMockBackend(),
			},
			args{
				entity: &client.Client{
					ID: "bad",
				},
			},
			nil,
			true,
		},
		{
			"valid volume update - no-op",
			fields{
				b: mock.NewMockBackend(),
			},
			args{
				entity: &volume.Volume{
					ID: "foo",
				},
			},
			&volume.Volume{
				ID:     "foo",
				Region: "us-west-2",
				Tag:    "foo",
				Status: volume.Unavailable,
			},
			false,
		},
		{
			"valid volume - new region",
			fields{
				b: mock.NewMockBackend(),
			},
			args{
				entity: &volume.Volume{
					ID:     "foo",
					Region: "us-east-1",
				},
			},
			&volume.Volume{
				ID:     "foo",
				Region: "us-east-1",
				Tag:    "foo",
				Status: volume.Unavailable,
			},
			false,
		},
		{
			"valid volume - new tag",
			fields{
				b: mock.NewMockBackend(),
			},
			args{
				entity: &volume.Volume{
					ID:  "foo",
					Tag: "bar",
				},
			},
			&volume.Volume{
				ID:     "foo",
				Region: "us-west-2",
				Tag:    "bar",
				Status: volume.Unavailable,
			},
			false,
		},
		{
			"valid volume - new status",
			fields{
				b: mock.NewMockBackend(),
			},
			args{
				entity: &volume.Volume{
					ID:     "foo",
					Status: volume.Available,
				},
			},
			&volume.Volume{
				ID:     "foo",
				Region: "us-west-2",
				Tag:    "foo",
				Status: volume.Available,
			},
			false,
		},
		{
			"invalid volume - bad volume",
			fields{
				b: mock.NewMockBackend(),
			},
			args{
				entity: &volume.Volume{
					ID: "bad",
				},
			},
			nil,
			true,
		},
		{
			"valid lease request - no-op",
			fields{
				b: mock.NewMockBackend(),
			},
			args{
				entity: &leaserequest.LeaseRequest{
					ID:     "foo",
					Status: leaserequest.Pending,
				},
			},
			&leaserequest.LeaseRequest{
				ID:       "foo",
				ClientID: "foo",
				Region:   "us-west-2",
				Tag:      "foo",
				Status:   leaserequest.Pending,
			},
			false,
		},
		{
			"invalid lease request - change client id",
			fields{
				b: mock.NewMockBackend(),
			},
			args{
				entity: &leaserequest.LeaseRequest{
					ID:       "foo",
					ClientID: "bar",
				},
			},
			nil,
			true,
		},
		{
			"invalid lease request - change region",
			fields{
				b: mock.NewMockBackend(),
			},
			args{
				entity: &leaserequest.LeaseRequest{
					ID:     "foo",
					Region: "eu-west-1",
				},
			},
			nil,
			true,
		},
		{
			"invalid lease request - change tag",
			fields{
				b: mock.NewMockBackend(),
			},
			args{
				entity: &leaserequest.LeaseRequest{
					ID:  "foo",
					Tag: "bar",
				},
			},
			nil,
			true,
		},
		{
			"invalid lease request - invalid status",
			fields{
				b: mock.NewMockBackend(),
			},
			args{
				entity: &leaserequest.LeaseRequest{
					ID:     "foo",
					Status: leaserequest.Fulfilled,
				},
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Backend{
				b: tt.fields.b,
			}
			e := tt.args.entity
			if err := b.Update(e); (err != nil) != tt.wantErr {
				t.Errorf("Backend.Update() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.want != nil {
				// mock
				//
				// e.FSM = nil
				// e.CreatedAt = nil
				// e.UpdatedAt = nil

				createdAt := reflect.ValueOf(e).Elem().FieldByName("CreatedAt").Interface()
				updatedAt := reflect.ValueOf(e).Elem().FieldByName("UpdatedAt").Interface()
				if createdAt == updatedAt {
					t.Errorf("Backend.Update() did not modify UpdatedAt")
				}

				for _, field := range []string{"FSM", "CreatedAt", "UpdatedAt"} {
					f := reflect.ValueOf(e).Elem().FieldByName(field)
					f.Set(reflect.Zero(f.Type()))
				}

				if !reflect.DeepEqual(e, tt.want) {
					t.Errorf("Backend.Update() = %v, want %v", e, tt.want)
				}
			}
		})
	}
}

func TestBackend_Delete(t *testing.T) {
	type fields struct {
		b backend
	}
	type args struct {
		entity model.Base
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"valid client",
			fields{
				b: mock.NewMockBackend(),
			},
			args{
				&client.Client{
					ID: "foo",
				},
			},
			false,
		},
		{
			"invalid client",
			fields{
				b: mock.NewMockBackend(),
			},
			args{
				&client.Client{
					ID: "bad",
				},
			},
			true,
		},
		{
			"valid client with associated lease request",
			fields{
				b: mock.NewMockBackend(mock.WithMocks(map[string]model.Base{
					"Client": &client.Client{
						ID:     "foo",
						Status: client.Alive,
					},
					"LeaseRequest": &leaserequest.LeaseRequest{
						ID:       "foo",
						ClientID: "foo",
						Status:   leaserequest.Pending,
					},
				})),
			},
			args{
				&client.Client{
					ID: "foo",
				},
			},
			false,
		},
		{
			"valid client with error deleting associated lease request",
			fields{
				b: mock.NewMockBackend(mock.WithMocks(map[string]model.Base{
					"Client": &client.Client{
						ID:     "foo",
						Status: client.Alive,
					},
					"LeaseRequest": &leaserequest.LeaseRequest{
						ID:       "bad",
						ClientID: "foo",
						Status:   leaserequest.Pending,
					},
				})),
			},
			args{
				&client.Client{
					ID: "foo",
				},
			},
			true,
		},
		{
			"valid client with bad transition deleting associated lease request",
			fields{
				b: mock.NewMockBackend(mock.WithMocks(map[string]model.Base{
					"Client": &client.Client{
						ID:     "foo",
						Status: client.Alive,
					},
					"LeaseRequest": &leaserequest.LeaseRequest{
						ID:       "foo",
						ClientID: "foo",
						Status:   leaserequest.Fulfilled,
					},
				})),
			},
			args{
				&client.Client{
					ID: "foo",
				},
			},
			true,
		},
		{
			"valid volume",
			fields{
				b: mock.NewMockBackend(),
			},
			args{
				&volume.Volume{
					ID: "foo",
				},
			},
			false,
		},
		{
			"invalid volume",
			fields{
				b: mock.NewMockBackend(),
			},
			args{
				&volume.Volume{
					ID: "bad",
				},
			},
			true,
		},
		{
			"valid lease request",
			fields{
				b: mock.NewMockBackend(),
			},
			args{
				&leaserequest.LeaseRequest{
					ID: "foo",
				},
			},
			false,
		},
		{
			"invalid lease request",
			fields{
				b: mock.NewMockBackend(),
			},
			args{
				&leaserequest.LeaseRequest{
					ID: "bad",
				},
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Backend{
				b: tt.fields.b,
			}
			if err := b.Delete(tt.args.entity); (err != nil) != tt.wantErr {
				t.Errorf("Backend.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
