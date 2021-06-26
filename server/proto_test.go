package server

import (
	"reflect"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes"

	"github.com/p0pr0ck5/volchestrator/server/client"
	"github.com/p0pr0ck5/volchestrator/server/notification"
	"github.com/p0pr0ck5/volchestrator/server/volume"
	"github.com/p0pr0ck5/volchestrator/svc"
)

func Test_toProto(t *testing.T) {
	mockRegistered := time.Now()
	mockLastSeen := time.Now()

	registeredTS, _ := ptypes.TimestampProto(mockRegistered)
	lastSeenTS, _ := ptypes.TimestampProto(mockLastSeen)

	type args struct {
		from interface{}
	}
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		{
			"client",
			args{
				from: &client.Client{
					ID:         "foo",
					Registered: mockRegistered,
					LastSeen:   mockLastSeen,
				},
			},
			&svc.Client{
				ClientId:   "foo",
				Registered: registeredTS,
				LastSeen:   lastSeenTS,
			},
		},
		{
			"client without timestamps",
			args{
				from: &client.Client{
					ID: "foo",
				},
			},
			&svc.Client{
				ClientId: "foo",
			},
		},
		{
			"volume",
			args{
				from: &volume.Volume{
					ID:     "foo",
					Region: "us-west-2",
					Tag:    "foo",
					Status: volume.Attached,
				},
			},
			&svc.Volume{
				VolumeId: "foo",
				Region:   "us-west-2",
				Tag:      "foo",
				Status:   svc.Volume_Attached,
			},
		},
		{
			"volume without status defined",
			args{
				from: &volume.Volume{
					ID:     "foo",
					Region: "us-west-2",
					Tag:    "foo",
				},
			},
			&svc.Volume{
				VolumeId: "foo",
				Region:   "us-west-2",
				Tag:      "foo",
				Status:   svc.Volume_Available,
			},
		},
		{
			"notification",
			args{
				from: &notification.Notification{
					ClientID:  "foo",
					Message:   "bar",
					MessageID: 1,
				},
			},
			&svc.Notification{
				Message:   "bar",
				MessageId: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := toProto(tt.args.from); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("toProto() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_toStruct(t *testing.T) {

	type args struct {
		from interface{}
	}
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		{
			"volume",
			args{
				from: &svc.Volume{
					VolumeId: "foo",
					Region:   "us-west-2",
					Tag:      "bar",
					Status:   svc.Volume_Attached,
				},
			},
			&volume.Volume{
				ID:     "foo",
				Region: "us-west-2",
				Tag:    "bar",
				Status: volume.Attached,
			},
		},
		{
			"volume without status defined",
			args{
				from: &svc.Volume{
					VolumeId: "foo",
					Region:   "us-west-2",
					Tag:      "bar",
				},
			},
			&volume.Volume{
				ID:     "foo",
				Region: "us-west-2",
				Tag:    "bar",
				Status: volume.Available,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := toStruct(tt.args.from); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("toProto() = %v, want %v", got, tt.want)
			}
		})
	}
}
