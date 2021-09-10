package server

import (
	"reflect"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes"

	"github.com/p0pr0ck5/volchestrator/server/client"
	leaserequest "github.com/p0pr0ck5/volchestrator/server/lease_request"
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
					Token:      "bar",
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
			"client in deleting state",
			args{
				from: &client.Client{
					ID:     "foo",
					Status: client.Deleting,
				},
			},
			&svc.Client{
				ClientId: "foo",
				Status:   svc.Client_Deleting,
			},
		},
		{
			"client with a lease",
			args{
				from: &client.Client{
					ID:         "foo",
					Token:      "bar",
					LeaseID:    "foo",
					Registered: mockRegistered,
					LastSeen:   mockLastSeen,
				},
			},
			&svc.Client{
				ClientId:   "foo",
				LeaseId:    "foo",
				Registered: registeredTS,
				LastSeen:   lastSeenTS,
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
				Status:   svc.Volume_unused,
			},
		},
		{
			"volume in deleting state",
			args{
				from: &volume.Volume{
					ID:     "foo",
					Region: "us-west-2",
					Tag:    "foo",
					Status: volume.Deleting,
				},
			},
			&svc.Volume{
				VolumeId: "foo",
				Region:   "us-west-2",
				Tag:      "foo",
				Status:   svc.Volume_Deleting,
			},
		},
		{
			"volume with a lease",
			args{
				from: &volume.Volume{
					ID:      "foo",
					LeaseID: "foo",
					Region:  "us-west-2",
					Tag:     "foo",
					Status:  volume.Attached,
				},
			},
			&svc.Volume{
				VolumeId: "foo",
				LeaseId:  "foo",
				Region:   "us-west-2",
				Tag:      "foo",
				Status:   svc.Volume_Attached,
			},
		},
		{
			"notification",
			args{
				from: &notification.Notification{
					ClientID: "foo",
					Message:  "bar",
				},
			},
			&svc.Notification{
				Message: "bar",
			},
		},
		{
			"lease request",
			args{
				from: &leaserequest.LeaseRequest{
					ID:       "foo",
					ClientID: "bar",
					Region:   "us-west-2",
					Tag:      "baz",
					Status:   leaserequest.Pending,
				},
			},
			&svc.LeaseRequest{
				LeaseRequestId: "foo",
				ClientId:       "bar",
				Region:         "us-west-2",
				Tag:            "baz",
				Status:         svc.LeaseRequest_Pending,
			},
		},
		{
			"lease request without status",
			args{
				from: &leaserequest.LeaseRequest{
					ID:       "foo",
					ClientID: "bar",
					Region:   "us-west-2",
					Tag:      "baz",
				},
			},
			&svc.LeaseRequest{
				LeaseRequestId: "foo",
				ClientId:       "bar",
				Region:         "us-west-2",
				Tag:            "baz",
				Status:         0,
			},
		},
		{
			"lease request in deleting state",
			args{
				from: &leaserequest.LeaseRequest{
					ID:       "foo",
					ClientID: "bar",
					Region:   "us-west-2",
					Tag:      "baz",
					Status:   leaserequest.Deleting,
				},
			},
			&svc.LeaseRequest{
				LeaseRequestId: "foo",
				ClientId:       "bar",
				Region:         "us-west-2",
				Tag:            "baz",
				Status:         svc.LeaseRequest_Deleting,
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
				Status: 0,
			},
		},
		{
			"lease request",
			args{
				from: &svc.LeaseRequest{
					LeaseRequestId: "foo",
					ClientId:       "bar",
					Region:         "us-west-2",
					Tag:            "baz",
					Status:         svc.LeaseRequest_Pending,
				},
			},
			&leaserequest.LeaseRequest{
				ID:       "foo",
				ClientID: "bar",
				Region:   "us-west-2",
				Tag:      "baz",
				Status:   leaserequest.Pending,
			},
		},
		{
			"lease request without status",
			args{
				from: &svc.LeaseRequest{
					LeaseRequestId: "foo",
					ClientId:       "bar",
					Region:         "us-west-2",
					Tag:            "baz",
				},
			},
			&leaserequest.LeaseRequest{
				ID:       "foo",
				ClientID: "bar",
				Region:   "us-west-2",
				Tag:      "baz",
				Status:   0,
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
