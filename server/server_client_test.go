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
			&svc.RegisterResponse{},
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
			s, _ := NewServer(WithMockBackend())
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
				},
				nil,
			},
			false,
		},
		{
			"invalid client",
			args{
				&svc.WatchNotificationsRequest{
					ClientId: "bad",
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
