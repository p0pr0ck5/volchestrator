package leaserequest

import (
	"testing"

	"github.com/p0pr0ck5/volchestrator/server/model"
)

func TestLeaseRequest_Validate(t *testing.T) {
	type fields struct {
		Model    model.Model
		ID       string
		ClientID string
		Region   string
		Tag      string
		Status   Status
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			"valid lease request",
			fields{
				ID:       "foo",
				ClientID: "foo",
				Region:   "us-west-2",
				Tag:      "foo",
				Status:   Pending,
			},
			false,
		},
		{
			"invalid lease request - missing id",
			fields{
				ClientID: "foo",
				Region:   "us-west-2",
				Tag:      "foo",
				Status:   Pending,
			},
			true,
		},
		{
			"invalid lease request - missing client id",
			fields{
				ID:     "foo",
				Region: "us-west-2",
				Tag:    "foo",
				Status: Pending,
			},
			true,
		},
		{
			"invalid lease request - missing region",
			fields{
				ID:       "foo",
				ClientID: "foo",
				Tag:      "foo",
				Status:   Pending,
			},
			true,
		},
		{
			"invalid lease request - missing tag",
			fields{
				ID:       "foo",
				ClientID: "foo",
				Region:   "us-west-2",
				Status:   Pending,
			},
			true,
		},
		{
			"invalid lease request - invalid status",
			fields{
				ID:       "foo",
				ClientID: "foo",
				Region:   "us-west-2",
				Tag:      "foo",
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &LeaseRequest{
				Model:    tt.fields.Model,
				ID:       tt.fields.ID,
				ClientID: tt.fields.ClientID,
				Region:   tt.fields.Region,
				Tag:      tt.fields.Tag,
				Status:   tt.fields.Status,
			}
			if err := l.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("LeaseRequest.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLeaseRequest_ValidateTransition(t *testing.T) {
	type fields struct {
		Model    model.Model
		ID       string
		ClientID string
		Region   string
		Tag      string
		Status   Status
	}
	type args struct {
		m model.Base
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"valid transition - new status",
			fields{
				ID:       "foo",
				ClientID: "foo",
				Region:   "us-west-2",
				Tag:      "foo",
				Status:   Pending,
			},
			args{
				m: &LeaseRequest{
					ID:       "foo",
					ClientID: "foo",
					Region:   "us-west-2",
					Tag:      "foo",
					Status:   Fulfilling,
				},
			},
			false,
		},
		{
			"invalid transition - change id",
			fields{
				ID:       "foo",
				ClientID: "foo",
				Region:   "us-west-2",
				Tag:      "foo",
				Status:   Pending,
			},
			args{
				m: &LeaseRequest{
					ID:       "bar",
					ClientID: "foo",
					Region:   "us-west-2",
					Tag:      "foo",
				},
			},
			true,
		},
		{
			"invalid transition - change client id",
			fields{
				ID:       "foo",
				ClientID: "foo",
				Region:   "us-west-2",
				Tag:      "foo",
				Status:   Pending,
			},
			args{
				m: &LeaseRequest{
					ID:       "foo",
					ClientID: "bar",
					Region:   "us-west-2",
					Tag:      "foo",
				},
			},
			true,
		},
		{
			"invalid transition - change region",
			fields{
				ID:       "foo",
				ClientID: "foo",
				Region:   "us-west-2",
				Tag:      "foo",
				Status:   Pending,
			},
			args{
				m: &LeaseRequest{
					ID:       "foo",
					ClientID: "foo",
					Region:   "us-east-1",
					Tag:      "foo",
				},
			},
			true,
		},
		{
			"invalid transition - change tag",
			fields{
				ID:       "foo",
				ClientID: "foo",
				Region:   "us-west-2",
				Tag:      "foo",
				Status:   Pending,
			},
			args{
				m: &LeaseRequest{
					ID:       "foo",
					ClientID: "foo",
					Region:   "us-west-2",
					Tag:      "bar",
					Status:   Pending,
				},
			},
			true,
		},
		{
			"invalid transition - invalid status",
			fields{
				ID:       "foo",
				ClientID: "foo",
				Region:   "us-west-2",
				Tag:      "foo",
				Status:   Pending,
			},
			args{
				m: &LeaseRequest{
					ID:       "foo",
					ClientID: "foo",
					Region:   "us-west-2",
					Tag:      "bar",
					Status:   Fulfilled,
				},
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &LeaseRequest{
				Model:    tt.fields.Model,
				ID:       tt.fields.ID,
				ClientID: tt.fields.ClientID,
				Region:   tt.fields.Region,
				Tag:      tt.fields.Tag,
				Status:   tt.fields.Status,
			}
			l.Init()
			if err := l.ValidateTransition(tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("LeaseRequest.ValidateTransition() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
