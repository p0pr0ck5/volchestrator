package server

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes"
	timestamppb "github.com/golang/protobuf/ptypes/timestamp"

	"github.com/p0pr0ck5/volchestrator/server/volume"
	"github.com/p0pr0ck5/volchestrator/svc"
)

func timeToProto(t time.Time) *timestamppb.Timestamp {
	ts, _ := ptypes.TimestampProto(t)
	return ts
}

func TestServer_GetClient(t *testing.T) {
	type args struct {
		ctx context.Context
		req *svc.GetClientRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *svc.GetClientResponse
		wantErr bool
	}{
		{
			"valid client",
			args{
				context.Background(),
				&svc.GetClientRequest{
					ClientId: "foo",
				},
			},
			&svc.GetClientResponse{
				Client: &svc.Client{
					ClientId:   "foo",
					Registered: timeToProto(nowIsh()),
					LastSeen:   timeToProto(nowIsh()),
				},
			},
			false,
		},
		{
			"invalid client",
			args{
				context.Background(),
				&svc.GetClientRequest{
					ClientId: "bad",
				},
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, _ := NewServer(WithMockBackend())
			got, err := s.GetClient(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Server.GetClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Server.GetClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServer_ListClients(t *testing.T) {
	type args struct {
		ctx context.Context
		req *svc.ListClientsRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *svc.ListClientsResponse
		wantErr bool
	}{
		{
			"list clients",
			args{
				context.Background(),
				&svc.ListClientsRequest{},
			},
			&svc.ListClientsResponse{
				Clients: []*svc.Client{
					{
						ClientId:   "foo",
						Registered: timeToProto(nowIsh()),
						LastSeen:   timeToProto(nowIsh()),
					},
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, _ := NewServer(WithMockBackend())
			got, err := s.ListClients(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Server.ListClients() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Server.ListClients() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServer_GetVolume(t *testing.T) {
	type args struct {
		ctx context.Context
		req *svc.GetVolumeRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *svc.GetVolumeResponse
		wantErr bool
	}{
		{
			"valid volume",
			args{
				context.Background(),
				&svc.GetVolumeRequest{
					VolumeId: "foo",
				},
			},
			&svc.GetVolumeResponse{
				Volume: &svc.Volume{
					VolumeId: "foo",
					Region:   "us-west-2",
					Tag:      "foo",
					Status:   svc.Volume_Status(volume.Unavailable),
				},
			},
			false,
		},
		{
			"invalid volume",
			args{
				context.Background(),
				&svc.GetVolumeRequest{
					VolumeId: "bad",
				},
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, _ := NewServer(WithMockBackend())
			got, err := s.GetVolume(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Server.GetVolume() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Server.GetVolume() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServer_AddVolume(t *testing.T) {
	type args struct {
		ctx context.Context
		req *svc.AddVolumeRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *svc.AddVolumeResponse
		wantErr bool
	}{
		{
			"valid volume",
			args{
				context.Background(),
				&svc.AddVolumeRequest{
					Volume: &svc.Volume{
						VolumeId: "foo",
						Region:   "us-west-2",
						Tag:      "bar",
					},
				},
			},
			&svc.AddVolumeResponse{},
			false,
		},
		{
			"valid volume with unavailable status",
			args{
				context.Background(),
				&svc.AddVolumeRequest{
					Volume: &svc.Volume{
						VolumeId: "foo",
						Region:   "us-west-2",
						Tag:      "bar",
						Status:   svc.Volume_Status(volume.Unavailable),
					},
				},
			},
			&svc.AddVolumeResponse{},
			false,
		},
		{
			"invalid volume - bad id in request",
			args{
				context.Background(),
				&svc.AddVolumeRequest{
					Volume: &svc.Volume{
						VolumeId: "bad",
						Region:   "us-west-2",
						Tag:      "bar",
					},
				},
			},
			nil,
			true,
		},
		{
			"invalid volume - missing id in request",
			args{
				context.Background(),
				&svc.AddVolumeRequest{
					Volume: &svc.Volume{
						Region: "us-west-2",
						Tag:    "bar",
					},
				},
			},
			nil,
			true,
		},
		{
			"invalid volume - missing region in request",
			args{
				context.Background(),
				&svc.AddVolumeRequest{
					Volume: &svc.Volume{
						VolumeId: "foo",
						Tag:      "bar",
					},
				},
			},
			nil,
			true,
		},
		{
			"invalid volume - missing tag in request",
			args{
				context.Background(),
				&svc.AddVolumeRequest{
					Volume: &svc.Volume{
						VolumeId: "foo",
						Region:   "us-west-2",
					},
				},
			},
			nil,
			true,
		},
		{
			"invalid volume - invalid status",
			args{
				context.Background(),
				&svc.AddVolumeRequest{
					Volume: &svc.Volume{
						VolumeId: "foo",
						Region:   "us-west-2",
						Status:   svc.Volume_Status(volume.Attaching),
					},
				},
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, _ := NewServer(WithMockBackend())
			got, err := s.AddVolume(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Server.AddVolume() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Server.AddVolume() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServer_UpdateVolume(t *testing.T) {
	type args struct {
		ctx context.Context
		req *svc.UpdateVolumeRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *svc.UpdateVolumeResponse
		wantErr bool
	}{
		{
			"valid volume",
			args{
				context.Background(),
				&svc.UpdateVolumeRequest{
					Volume: &svc.Volume{
						VolumeId: "foo",
						Region:   "us-west-2",
						Tag:      "bar",
					},
				},
			},
			&svc.UpdateVolumeResponse{},
			false,
		},
		{
			"valid volume with unavailable status",
			args{
				context.Background(),
				&svc.UpdateVolumeRequest{
					Volume: &svc.Volume{
						VolumeId: "foo",
						Region:   "us-west-2",
						Tag:      "bar",
						Status:   svc.Volume_Status(volume.Unavailable),
					},
				},
			},
			&svc.UpdateVolumeResponse{},
			false,
		},
		{
			"invalid volume - bad id in request",
			args{
				context.Background(),
				&svc.UpdateVolumeRequest{
					Volume: &svc.Volume{
						VolumeId: "bad",
						Region:   "us-west-2",
						Tag:      "bar",
					},
				},
			},
			nil,
			true,
		},
		{
			"invalid volume - missing id in request",
			args{
				context.Background(),
				&svc.UpdateVolumeRequest{
					Volume: &svc.Volume{
						Region: "us-west-2",
						Tag:    "bar",
					},
				},
			},
			nil,
			true,
		},
		{
			"invalid volume - missing region in request",
			args{
				context.Background(),
				&svc.UpdateVolumeRequest{
					Volume: &svc.Volume{
						VolumeId: "foo",
						Tag:      "bar",
					},
				},
			},
			nil,
			true,
		},
		{
			"invalid volume - missing tag in request",
			args{
				context.Background(),
				&svc.UpdateVolumeRequest{
					Volume: &svc.Volume{
						VolumeId: "foo",
						Region:   "us-west-2",
					},
				},
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, _ := NewServer(WithMockBackend())
			got, err := s.UpdateVolume(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Server.UpdateVolume() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Server.UpdateVolume() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServer_DeleteVolume(t *testing.T) {
	type args struct {
		ctx context.Context
		req *svc.DeleteVolumeRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *svc.DeleteVolumeResponse
		wantErr bool
	}{
		{
			"valid volume",
			args{
				context.Background(),
				&svc.DeleteVolumeRequest{
					Volume: &svc.Volume{
						VolumeId: "foo",
					},
				},
			},
			&svc.DeleteVolumeResponse{},
			false,
		},
		{
			"valid volume",
			args{
				context.Background(),
				&svc.DeleteVolumeRequest{
					Volume: &svc.Volume{
						VolumeId: "bad",
					},
				},
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, _ := NewServer(WithMockBackend())
			got, err := s.DeleteVolume(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Server.DeleteVolume() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Server.DeleteVolume() = %v, want %v", got, tt.want)
			}
		})
	}
}
