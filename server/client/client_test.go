package client

import (
	"testing"
	"time"

	"github.com/p0pr0ck5/volchestrator/server/model"
)

func TestClient_Validate(t *testing.T) {
	type fields struct {
		Model      model.Model
		ID         string
		Token      string
		Status     Status
		Registered time.Time
		LastSeen   time.Time
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			"valid client",
			fields{
				ID:    "foo",
				Token: "mock",
			},
			false,
		},
		{
			"invalid client - missing id",
			fields{
				Token: "mock",
			},
			true,
		},
		{
			"invalid client - missing token",
			fields{
				ID: "foo",
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				Model:      tt.fields.Model,
				ID:         tt.fields.ID,
				Token:      tt.fields.Token,
				Status:     tt.fields.Status,
				Registered: tt.fields.Registered,
				LastSeen:   tt.fields.LastSeen,
			}
			if err := c.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Client.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_ValidateTransition(t *testing.T) {
	type fields struct {
		Model      model.Model
		ID         string
		Token      string
		Status     Status
		Registered time.Time
		LastSeen   time.Time
	}
	type args struct {
		newClient *Client
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"valid transition",
			fields{
				ID:    "foo",
				Token: "mock",
			},
			args{
				newClient: &Client{
					ID:     "foo",
					Status: Deleting,
				},
			},
			false,
		},
		{
			"valid transition - no change",
			fields{
				ID:    "foo",
				Token: "mock",
			},
			args{
				newClient: &Client{
					ID:     "foo",
					Status: Alive,
				},
			},
			false,
		},
		{
			"invalid transition - id change",
			fields{
				ID:    "foo",
				Token: "mock",
			},
			args{
				newClient: &Client{
					ID: "bar",
				},
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				Model:      tt.fields.Model,
				ID:         tt.fields.ID,
				Token:      tt.fields.Token,
				Status:     tt.fields.Status,
				Registered: tt.fields.Registered,
				LastSeen:   tt.fields.LastSeen,
			}
			c.Init()
			if err := c.ValidateTransition(tt.args.newClient); (err != nil) != tt.wantErr {
				t.Errorf("Client.ValidateTransition() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
