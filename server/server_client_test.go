package server

import (
	"context"
	"reflect"
	"testing"

	"github.com/p0pr0ck5/volchestrator/svc"
)

func TestServer_Register(t *testing.T) {
	type args struct {
		ctx context.Context
		req *svc.RegisterRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *svc.RegisterResponse
		wantErr bool
	}{
		{
			"valid register",
			args{
				context.Background(),
				&svc.RegisterRequest{
					ClientId: "foo",
				},
			},
			&svc.RegisterResponse{
				ClientId: "foo",
				Token:    "mock",
			},
			false,
		},
		{
			"invalid register - missing client id in request",
			args{
				context.Background(),
				&svc.RegisterRequest{},
			},
			nil,
			true,
		},
		{
			"invalid register - bad client id in request",
			args{
				context.Background(),
				&svc.RegisterRequest{
					ClientId: "bad",
				},
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, _ := NewServer(WithMockBackend(), WithTokener(mockTokener{}))
			got, err := s.Register(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Server.Register() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Server.Register() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServer_Deregister(t *testing.T) {
	type args struct {
		ctx context.Context
		req *svc.DeregisterRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *svc.DeregisterResponse
		wantErr bool
	}{
		{
			"valid deregister",
			args{
				context.Background(),
				&svc.DeregisterRequest{
					ClientId: "foo",
					Token:    "mock",
				},
			},
			&svc.DeregisterResponse{},
			false,
		},
		{
			"invalid deregister - missing client id in request",
			args{
				context.Background(),
				&svc.DeregisterRequest{},
			},
			nil,
			true,
		},
		{
			"invalid deregister - bad client id in request",
			args{
				context.Background(),
				&svc.DeregisterRequest{
					ClientId: "bad",
				},
			},
			nil,
			true,
		},
		{
			"invalid deregister - bad token in request",
			args{
				context.Background(),
				&svc.DeregisterRequest{
					ClientId: "foo",
					Token:    "nope",
				},
			},
			nil,
			true,
		},
		{
			"invalid deregister - missing token in request",
			args{
				context.Background(),
				&svc.DeregisterRequest{
					ClientId: "foo",
				},
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, _ := NewServer(WithMockBackend())
			got, err := s.Deregister(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Server.Deregister() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Server.Deregister() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServer_Ping(t *testing.T) {
	type args struct {
		ctx context.Context
		req *svc.PingRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *svc.PingResponse
		wantErr bool
	}{
		{
			"valid ping",
			args{
				context.Background(),
				&svc.PingRequest{
					ClientId: "foo",
					Token:    "mock",
				},
			},
			&svc.PingResponse{},
			false,
		},
		{
			"invalid ping - missing client id in request",
			args{
				context.Background(),
				&svc.PingRequest{},
			},
			nil,
			true,
		},
		{
			"invalid ping - bad client id in request",
			args{
				context.Background(),
				&svc.PingRequest{
					ClientId: "bad",
				},
			},
			nil,
			true,
		},
		{
			"invalid ping - bad token in request",
			args{
				context.Background(),
				&svc.PingRequest{
					ClientId: "foo",
					Token:    "nope",
				},
			},
			nil,
			true,
		},
		{
			"invalid ping - missing token in request",
			args{
				context.Background(),
				&svc.PingRequest{
					ClientId: "foo",
				},
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, _ := NewServer(WithMockBackend())
			got, err := s.Ping(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Server.Ping() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Server.Ping() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServer_WatchNotifications(t *testing.T) {
	type args struct {
		req    *svc.WatchNotificationsRequest
		stream svc.Volchestrator_WatchNotificationsServer
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"stream response",
			args{
				&svc.WatchNotificationsRequest{
					ClientId: "foo",
					Token:    "mock",
				},
				nil,
			},
			true,
		},
		{
			"invalid client",
			args{
				&svc.WatchNotificationsRequest{
					ClientId: "bad",
					Token:    "mock",
				},
				nil,
			},
			true,
		},
		{
			"invalid token",
			args{
				&svc.WatchNotificationsRequest{
					ClientId: "foo",
					Token:    "nope",
				},
				nil,
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, _ := NewServer(WithMockBackend())
			if err := s.WatchNotifications(tt.args.req, tt.args.stream); (err != nil) != tt.wantErr {
				t.Errorf("Server.WatchNotifications() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServer_RequestLease(t *testing.T) {
	type args struct {
		ctx context.Context
		req *svc.RequestLeaseRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *svc.RequestLeaseResponse
		wantErr bool
	}{
		{
			"valid lease request",
			args{
				ctx: context.Background(),
				req: &svc.RequestLeaseRequest{
					LeaseRequestId: "foo",
					ClientId:       "foo",
					Token:          "mock",
					Region:         "us-west-2",
					Tag:            "baz",
				},
			},
			&svc.RequestLeaseResponse{},
			false,
		},
		{
			"invalid lease request - missing lease request id",
			args{
				ctx: context.Background(),
				req: &svc.RequestLeaseRequest{
					ClientId: "foo",
					Token:    "mock",
					Region:   "us-west-2",
					Tag:      "baz",
				},
			},
			nil,
			true,
		},
		{
			"invalid lease request - bad lease request id",
			args{
				ctx: context.Background(),
				req: &svc.RequestLeaseRequest{
					LeaseRequestId: "bad",
					ClientId:       "foo",
					Token:          "mock",
					Region:         "us-west-2",
					Tag:            "baz",
				},
			},
			nil,
			true,
		},
		{
			"invalid lease request - missing client id",
			args{
				ctx: context.Background(),
				req: &svc.RequestLeaseRequest{
					LeaseRequestId: "foo",
					Token:          "mock",
					Region:         "us-west-2",
					Tag:            "baz",
				},
			},
			nil,
			true,
		},
		{
			"invalid lease request - bad client id",
			args{
				ctx: context.Background(),
				req: &svc.RequestLeaseRequest{
					LeaseRequestId: "foo",
					ClientId:       "bad",
					Token:          "mock",
					Region:         "us-west-2",
					Tag:            "baz",
				},
			},
			nil,
			true,
		},
		{
			"invalid lease request - bad client token",
			args{
				ctx: context.Background(),
				req: &svc.RequestLeaseRequest{
					LeaseRequestId: "foo",
					ClientId:       "foo",
					Token:          "nope",
					Region:         "us-west-2",
					Tag:            "baz",
				},
			},
			nil,
			true,
		},
		{
			"invalid lease request - missing client token",
			args{
				ctx: context.Background(),
				req: &svc.RequestLeaseRequest{
					LeaseRequestId: "foo",
					ClientId:       "foo",
					Region:         "us-west-2",
					Tag:            "baz",
				},
			},
			nil,
			true,
		},
		{
			"invalid lease request - missing region",
			args{
				ctx: context.Background(),
				req: &svc.RequestLeaseRequest{
					LeaseRequestId: "foo",
					ClientId:       "foo",
					Token:          "mock",
					Tag:            "baz",
				},
			},
			nil,
			true,
		},
		{
			"invalid lease request - missing tag",
			args{
				ctx: context.Background(),
				req: &svc.RequestLeaseRequest{
					LeaseRequestId: "foo",
					ClientId:       "foo",
					Token:          "mock",
					Region:         "us-west-2",
				},
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, _ := NewServer(WithMockBackend())
			got, err := s.RequestLease(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Server.RequestLease() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Server.RequestLease() = %v, want %v", got, tt.want)
			}
		})
	}
}
