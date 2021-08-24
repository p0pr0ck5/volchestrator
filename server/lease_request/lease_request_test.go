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
