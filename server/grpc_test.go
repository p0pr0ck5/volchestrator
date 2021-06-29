package server

import (
	"context"
	"io"
	"log"
	"net"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/proto"

	"github.com/p0pr0ck5/volchestrator/server/backend"
	"github.com/p0pr0ck5/volchestrator/server/client"
	"github.com/p0pr0ck5/volchestrator/server/notification"
	"github.com/p0pr0ck5/volchestrator/server/volume"
	"github.com/p0pr0ck5/volchestrator/svc"
)

const bufSize = 1024 * 1024

type bufDialFunc func(context.Context, string) (net.Conn, error)

func mockServer() (*Server, bufDialFunc) {
	lis := bufconn.Listen(bufSize)
	s := grpc.NewServer()
	srv, _ := NewServer(WithNewMemoryBackend())
	svc.RegisterVolchestratorServer(s, srv)
	svc.RegisterVolchestratorAdminServer(s, srv)
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()

	return srv, func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}
}

func Test_Register(t *testing.T) {
	type args struct {
		ctx context.Context
		req *svc.RegisterRequest
	}
	tests := []struct {
		name    string
		args    []args
		want    []*svc.RegisterResponse
		wantErr []bool
	}{
		{
			"one valid registration",
			[]args{
				{
					context.Background(),
					&svc.RegisterRequest{
						ClientId: "foo",
					},
				},
			},
			[]*svc.RegisterResponse{
				{},
			},
			[]bool{false},
		},
		{
			"one invalid registration",
			[]args{
				{
					context.Background(),
					&svc.RegisterRequest{},
				},
			},
			[]*svc.RegisterResponse{
				nil,
			},
			[]bool{true},
		},
		{
			"two valid registrations",
			[]args{
				{
					context.Background(),
					&svc.RegisterRequest{
						ClientId: "foo",
					},
				},
				{
					context.Background(),
					&svc.RegisterRequest{
						ClientId: "bar",
					},
				},
			},
			[]*svc.RegisterResponse{
				{},
				{},
			},
			[]bool{false, false},
		},
		{
			"duplicate registration",
			[]args{
				{
					context.Background(),
					&svc.RegisterRequest{
						ClientId: "foo",
					},
				},
				{
					context.Background(),
					&svc.RegisterRequest{
						ClientId: "foo",
					},
				},
			},
			[]*svc.RegisterResponse{
				{},
				nil,
			},
			[]bool{false, true},
		},
		{
			"duplicate registration after valid registration",
			[]args{
				{
					context.Background(),
					&svc.RegisterRequest{
						ClientId: "foo",
					},
				},
				{
					context.Background(),
					&svc.RegisterRequest{
						ClientId: "bar",
					},
				},
				{
					context.Background(),
					&svc.RegisterRequest{
						ClientId: "foo",
					},
				},
			},
			[]*svc.RegisterResponse{
				{},
				{},
				nil,
			},
			[]bool{false, false, true},
		},
		{
			"duplicate registration following valid registration",
			[]args{
				{
					context.Background(),
					&svc.RegisterRequest{
						ClientId: "bar",
					},
				},
				{
					context.Background(),
					&svc.RegisterRequest{
						ClientId: "foo",
					},
				},
				{
					context.Background(),
					&svc.RegisterRequest{
						ClientId: "foo",
					},
				},
			},
			[]*svc.RegisterResponse{
				{},
				{},
				nil,
			},
			[]bool{false, false, true},
		},
		{
			"valid registration following duplicate registration",
			[]args{
				{
					context.Background(),
					&svc.RegisterRequest{
						ClientId: "foo",
					},
				},
				{
					context.Background(),
					&svc.RegisterRequest{
						ClientId: "foo",
					},
				},
				{
					context.Background(),
					&svc.RegisterRequest{
						ClientId: "bar",
					},
				},
			},
			[]*svc.RegisterResponse{
				{},
				nil,
				{},
			},
			[]bool{false, true, false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv, bufDialer := mockServer()
			srv.b = backend.NewMemoryBackend()

			ctx := context.Background()
			conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
			if err != nil {
				t.Fatalf("Failed to dial bufnet: %v", err)
			}
			defer conn.Close()
			client := svc.NewVolchestratorClient(conn)

			for i, req := range tt.args {
				got, err := client.Register(req.ctx, req.req)
				if (err != nil) != tt.wantErr[i] {
					t.Errorf("Server.Register() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !proto.Equal(got, tt.want[i]) {
					t.Errorf("Server.Register() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func Test_Deregister(t *testing.T) {
	mockNow := time.Now()

	mockClients := []*client.Client{
		{
			ID:         "foo",
			Registered: mockNow,
			LastSeen:   mockNow,
		},
		{
			ID:         "bar",
			Registered: mockNow,
			LastSeen:   mockNow,
		},
	}

	type args struct {
		ctx context.Context
		req *svc.DeregisterRequest
	}
	tests := []struct {
		name    string
		args    []args
		want    []*svc.DeregisterResponse
		wantErr []bool
	}{
		{
			"one valid deregistration",
			[]args{
				{
					context.Background(),
					&svc.DeregisterRequest{
						ClientId: "foo",
					},
				},
			},
			[]*svc.DeregisterResponse{
				{},
			},
			[]bool{false},
		},
		{
			"two valid deregistrations",
			[]args{
				{
					context.Background(),
					&svc.DeregisterRequest{
						ClientId: "foo",
					},
				},
				{
					context.Background(),
					&svc.DeregisterRequest{
						ClientId: "bar",
					},
				},
			},
			[]*svc.DeregisterResponse{
				{},
				{},
			},
			[]bool{false, false},
		},
		{
			"one invalid deregistration (no registration)",
			[]args{
				{
					context.Background(),
					&svc.DeregisterRequest{
						ClientId: "baz",
					},
				},
			},
			[]*svc.DeregisterResponse{
				nil,
			},
			[]bool{true},
		},
		{
			"one invalid deregistration (empty request)",
			[]args{
				{
					context.Background(),
					&svc.DeregisterRequest{},
				},
			},
			[]*svc.DeregisterResponse{
				nil,
			},
			[]bool{true},
		},
		{
			"one valid and one invalid registration",
			[]args{
				{
					context.Background(),
					&svc.DeregisterRequest{
						ClientId: "foo",
					},
				},
				{
					context.Background(),
					&svc.DeregisterRequest{
						ClientId: "baz",
					},
				},
			},
			[]*svc.DeregisterResponse{
				{},
				nil,
			},
			[]bool{false, true},
		},
		{
			"one invalid and one valid registration",
			[]args{
				{
					context.Background(),
					&svc.DeregisterRequest{
						ClientId: "baz",
					},
				},
				{
					context.Background(),
					&svc.DeregisterRequest{
						ClientId: "foo",
					},
				},
			},
			[]*svc.DeregisterResponse{
				nil,
				{},
			},
			[]bool{true, false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv, bufDialer := mockServer()
			srv.b = backend.NewMemoryBackend(backend.WithClients(mockClients))

			ctx := context.Background()
			conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
			if err != nil {
				t.Fatalf("Failed to dial bufnet: %v", err)
			}
			defer conn.Close()
			client := svc.NewVolchestratorClient(conn)

			for i, req := range tt.args {
				got, err := client.Deregister(req.ctx, req.req)
				if (err != nil) != tt.wantErr[i] {
					t.Errorf("Server.Deregister() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !proto.Equal(got, tt.want[i]) {
					t.Errorf("Server.Deregister() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func Test_Ping(t *testing.T) {
	mockNow := time.Now()

	mockClients := []*client.Client{
		{
			ID:         "foo",
			Registered: mockNow,
			LastSeen:   mockNow,
		},
		{
			ID:         "bar",
			Registered: mockNow,
			LastSeen:   mockNow,
		},
	}

	type args struct {
		ctx context.Context
		req *svc.PingRequest
	}
	tests := []struct {
		name    string
		args    []args
		want    []*svc.PingResponse
		wantErr []bool
	}{
		{
			"one valid ping",
			[]args{
				{
					context.Background(),
					&svc.PingRequest{
						ClientId: "foo",
					},
				},
			},
			[]*svc.PingResponse{
				{},
			},
			[]bool{false},
		},
		{
			"one invalid ping (no registration)",
			[]args{
				{
					context.Background(),
					&svc.PingRequest{
						ClientId: "baz",
					},
				},
			},
			[]*svc.PingResponse{
				nil,
			},
			[]bool{true},
		},
		{
			"two valid pings (same client)",
			[]args{
				{
					context.Background(),
					&svc.PingRequest{
						ClientId: "foo",
					},
				},
				{
					context.Background(),
					&svc.PingRequest{
						ClientId: "foo",
					},
				},
			},
			[]*svc.PingResponse{
				{},
				{},
			},
			[]bool{false, false},
		},
		{
			"two valid pings (different clients)",
			[]args{
				{
					context.Background(),
					&svc.PingRequest{
						ClientId: "foo",
					},
				},
				{
					context.Background(),
					&svc.PingRequest{
						ClientId: "bar",
					},
				},
			},
			[]*svc.PingResponse{
				{},
				{},
			},
			[]bool{false, false},
		},
		{
			"one valid and one invalid ping",
			[]args{
				{
					context.Background(),
					&svc.PingRequest{
						ClientId: "foo",
					},
				},
				{
					context.Background(),
					&svc.PingRequest{
						ClientId: "baz",
					},
				},
			},
			[]*svc.PingResponse{
				{},
				nil,
			},
			[]bool{false, true},
		},
		{
			"one invalid and one valid ping",
			[]args{
				{
					context.Background(),
					&svc.PingRequest{
						ClientId: "baz",
					},
				},
				{
					context.Background(),
					&svc.PingRequest{
						ClientId: "bar",
					},
				},
			},
			[]*svc.PingResponse{
				nil,
				{},
			},
			[]bool{true, false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv, bufDialer := mockServer()
			srv.b = backend.NewMemoryBackend(backend.WithClients(mockClients))

			ctx := context.Background()
			conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
			if err != nil {
				t.Fatalf("Failed to dial bufnet: %v", err)
			}
			defer conn.Close()
			client := svc.NewVolchestratorClient(conn)

			for i, req := range tt.args {
				got, err := client.Ping(req.ctx, req.req)
				if (err != nil) != tt.wantErr[i] {
					t.Errorf("Server.Ping() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !proto.Equal(got, tt.want[i]) {
					t.Errorf("Server.Ping() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func Test_GetClient(t *testing.T) {
	mockNow := time.Now()
	mockTS, _ := ptypes.TimestampProto(mockNow)

	mockClients := []*client.Client{
		{
			ID:         "foo",
			Registered: mockNow,
			LastSeen:   mockNow,
		},
		{
			ID:         "bar",
			Registered: mockNow,
			LastSeen:   mockNow,
		},
	}

	type args struct {
		ctx context.Context
		req *svc.GetClientRequest
	}
	tests := []struct {
		name    string
		args    []args
		want    []*svc.GetClientResponse
		wantErr []bool
	}{
		{
			"registered client",
			[]args{
				{
					context.Background(),
					&svc.GetClientRequest{
						ClientId: "foo",
					},
				},
			},
			[]*svc.GetClientResponse{
				{
					Client: &svc.Client{
						ClientId:   "foo",
						Registered: mockTS,
						LastSeen:   mockTS,
					},
				},
			},
			[]bool{false},
		},
		{
			"different registered client",
			[]args{
				{
					context.Background(),
					&svc.GetClientRequest{
						ClientId: "bar",
					},
				},
			},
			[]*svc.GetClientResponse{
				{
					Client: &svc.Client{
						ClientId:   "bar",
						Registered: mockTS,
						LastSeen:   mockTS,
					},
				},
			},
			[]bool{false},
		},
		{
			"nonexistent client",
			[]args{
				{
					context.Background(),
					&svc.GetClientRequest{
						ClientId: "baz",
					},
				},
			},
			[]*svc.GetClientResponse{
				nil,
			},
			[]bool{true},
		},
		{
			"empty client id",
			[]args{
				{
					context.Background(),
					&svc.GetClientRequest{},
				},
			},
			[]*svc.GetClientResponse{
				nil,
			},
			[]bool{true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv, bufDialer := mockServer()
			srv.b = backend.NewMemoryBackend(backend.WithClients(mockClients))

			ctx := context.Background()
			conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
			if err != nil {
				t.Fatalf("Failed to dial bufnet: %v", err)
			}
			defer conn.Close()
			client := svc.NewVolchestratorAdminClient(conn)

			for i, req := range tt.args {
				got, err := client.GetClient(req.ctx, req.req)
				if (err != nil) != tt.wantErr[i] {
					t.Errorf("Server.GetClient() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !proto.Equal(got, tt.want[i]) {
					t.Errorf("Server.GetClient() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func Test_ListClients(t *testing.T) {
	mockNow := time.Now()
	mockTS, _ := ptypes.TimestampProto(mockNow)

	mockClients := []*client.Client{
		{
			ID:         "foo",
			Registered: mockNow,
			LastSeen:   mockNow,
		},
		{
			ID:         "bar",
			Registered: mockNow,
			LastSeen:   mockNow,
		},
	}

	type args struct {
		ctx context.Context
		req *svc.ListClientsRequest
	}
	tests := []struct {
		name    string
		args    []args
		want    []*svc.ListClientsResponse
		wantErr []bool
	}{
		{
			"list clients",
			[]args{
				{
					context.Background(),
					&svc.ListClientsRequest{},
				},
			},
			[]*svc.ListClientsResponse{
				{
					Clients: []*svc.Client{
						{
							ClientId:   "bar",
							Registered: mockTS,
							LastSeen:   mockTS,
						},
						{
							ClientId:   "foo",
							Registered: mockTS,
							LastSeen:   mockTS,
						},
					},
				},
			},
			[]bool{false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv, bufDialer := mockServer()
			srv.b = backend.NewMemoryBackend(backend.WithClients(mockClients))

			ctx := context.Background()
			conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
			if err != nil {
				t.Fatalf("Failed to dial bufnet: %v", err)
			}
			defer conn.Close()
			client := svc.NewVolchestratorAdminClient(conn)

			for i, req := range tt.args {
				got, err := client.ListClients(req.ctx, req.req)
				if (err != nil) != tt.wantErr[i] {
					t.Errorf("Server.ListClients() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !proto.Equal(got, tt.want[i]) {
					t.Errorf("Server.ListClients() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func Test_GetVolume(t *testing.T) {
	mockVolumes := []*volume.Volume{
		{
			ID:     "foo",
			Region: "us-west-2",
			Tag:    "bar",
			Status: volume.Available,
		},
		{
			ID:     "bar",
			Region: "us-west-1",
			Tag:    "baz",
			Status: volume.Unavailable,
		},
	}

	type args struct {
		ctx context.Context
		req *svc.GetVolumeRequest
	}
	tests := []struct {
		name    string
		args    []args
		want    []*svc.GetVolumeResponse
		wantErr []bool
	}{
		{
			"valid volume",
			[]args{
				{
					context.Background(),
					&svc.GetVolumeRequest{
						VolumeId: "foo",
					},
				},
			},
			[]*svc.GetVolumeResponse{
				{
					Volume: &svc.Volume{
						VolumeId: "foo",
						Region:   "us-west-2",
						Tag:      "bar",
						Status:   svc.Volume_Status(volume.Available),
					},
				},
			},
			[]bool{false},
		},
		{
			"different valid volume",
			[]args{
				{
					context.Background(),
					&svc.GetVolumeRequest{
						VolumeId: "bar",
					},
				},
			},
			[]*svc.GetVolumeResponse{
				{
					Volume: &svc.Volume{
						VolumeId: "bar",
						Region:   "us-west-1",
						Tag:      "baz",
						Status:   svc.Volume_Status(volume.Unavailable),
					},
				},
			},
			[]bool{false},
		},
		{
			"nonexistent volume",
			[]args{
				{
					context.Background(),
					&svc.GetVolumeRequest{
						VolumeId: "baz",
					},
				},
			},
			[]*svc.GetVolumeResponse{
				nil,
			},
			[]bool{true},
		},
		{
			"empty volume id",
			[]args{
				{
					context.Background(),
					&svc.GetVolumeRequest{},
				},
			},
			[]*svc.GetVolumeResponse{
				nil,
			},
			[]bool{true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv, bufDialer := mockServer()
			srv.b = backend.NewMemoryBackend(backend.WithVolumes(mockVolumes))

			ctx := context.Background()
			conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
			if err != nil {
				t.Fatalf("Failed to dial bufnet: %v", err)
			}
			defer conn.Close()
			client := svc.NewVolchestratorAdminClient(conn)

			for i, req := range tt.args {
				got, err := client.GetVolume(req.ctx, req.req)
				if (err != nil) != tt.wantErr[i] {
					t.Errorf("Server.GetVolume() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !proto.Equal(got, tt.want[i]) {
					t.Errorf("Server.GetVolume() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func Test_ListVolumes(t *testing.T) {
	mockVolumes := []*volume.Volume{
		{
			ID:     "foo",
			Region: "us-west-2",
			Tag:    "bar",
			Status: volume.Available,
		},
		{
			ID:     "bar",
			Region: "us-west-1",
			Tag:    "baz",
			Status: volume.Unavailable,
		},
	}

	type args struct {
		ctx context.Context
		req *svc.ListVolumesRequest
	}
	tests := []struct {
		name    string
		args    []args
		want    []*svc.ListVolumesReponse
		wantErr []bool
	}{
		{
			"list volumes",
			[]args{
				{
					context.Background(),
					&svc.ListVolumesRequest{},
				},
			},
			[]*svc.ListVolumesReponse{
				{
					Volumes: []*svc.Volume{
						{
							VolumeId: "bar",
							Region:   "us-west-1",
							Tag:      "baz",
							Status:   svc.Volume_Status(volume.Unavailable),
						},
						{
							VolumeId: "foo",
							Region:   "us-west-2",
							Tag:      "bar",
							Status:   svc.Volume_Status(volume.Available),
						},
					},
				},
			},
			[]bool{false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv, bufDialer := mockServer()
			srv.b = backend.NewMemoryBackend(backend.WithVolumes(mockVolumes))

			ctx := context.Background()
			conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
			if err != nil {
				t.Fatalf("Failed to dial bufnet: %v", err)
			}
			defer conn.Close()
			client := svc.NewVolchestratorAdminClient(conn)

			for i, req := range tt.args {
				got, err := client.ListVolumes(req.ctx, req.req)
				if (err != nil) != tt.wantErr[i] {
					t.Errorf("Server.ListVolumes() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !proto.Equal(got, tt.want[i]) {
					t.Errorf("Server.ListVolumes() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func Test_AddVolume(t *testing.T) {
	mockVolumes := []*volume.Volume{
		{
			ID:     "foo",
			Region: "us-west-2",
			Tag:    "bar",
			Status: volume.Available,
		},
		{
			ID:     "bar",
			Region: "us-west-1",
			Tag:    "baz",
			Status: volume.Unavailable,
		},
	}

	type args struct {
		ctx context.Context
		req *svc.AddVolumeRequest
	}
	tests := []struct {
		name    string
		args    []args
		want    []*svc.AddVolumeResponse
		wantErr []bool
	}{
		{
			"valid volume",
			[]args{
				{
					context.Background(),
					&svc.AddVolumeRequest{
						Volume: &svc.Volume{
							VolumeId: "baz",
							Region:   "us-west-2",
							Tag:      "foo",
						},
					},
				},
			},
			[]*svc.AddVolumeResponse{
				{},
			},
			[]bool{false},
		},
		{
			"valid volume with unavailable status",
			[]args{
				{
					context.Background(),
					&svc.AddVolumeRequest{
						Volume: &svc.Volume{
							VolumeId: "baz",
							Region:   "us-west-2",
							Tag:      "foo",
							Status:   svc.Volume_Unavailable,
						},
					},
				},
			},
			[]*svc.AddVolumeResponse{
				{},
			},
			[]bool{false},
		},
		{
			"different valid volume",
			[]args{
				{
					context.Background(),
					&svc.AddVolumeRequest{
						Volume: &svc.Volume{
							VolumeId: "bat",
							Region:   "us-west-2",
							Tag:      "foo",
						},
					},
				},
			},
			[]*svc.AddVolumeResponse{
				{},
			},
			[]bool{false},
		},
		{
			"conflicting volume",
			[]args{
				{
					context.Background(),
					&svc.AddVolumeRequest{
						Volume: &svc.Volume{
							VolumeId: "foo",
							Region:   "us-west-2",
							Tag:      "foo",
						},
					},
				},
			},
			[]*svc.AddVolumeResponse{
				nil,
			},
			[]bool{true},
		},
		{
			"empty volume id",
			[]args{
				{
					context.Background(),
					&svc.AddVolumeRequest{
						Volume: &svc.Volume{
							Region: "us-west-2",
							Tag:    "foo",
						},
					},
				},
			},
			[]*svc.AddVolumeResponse{
				nil,
			},
			[]bool{true},
		},
		{
			"empty volume region",
			[]args{
				{
					context.Background(),
					&svc.AddVolumeRequest{
						Volume: &svc.Volume{
							VolumeId: "bat",
							Tag:      "bar",
						},
					},
				},
			},
			[]*svc.AddVolumeResponse{
				nil,
			},
			[]bool{true},
		},
		{
			"empty volume tag",
			[]args{
				{
					context.Background(),
					&svc.AddVolumeRequest{
						Volume: &svc.Volume{
							VolumeId: "bat",
							Region:   "us-west-2",
						},
					},
				},
			},
			[]*svc.AddVolumeResponse{
				nil,
			},
			[]bool{true},
		},
		{
			"invalid volume status",
			[]args{
				{
					context.Background(),
					&svc.AddVolumeRequest{
						Volume: &svc.Volume{
							VolumeId: "bat",
							Region:   "us-west-2",
							Tag:      "baz",
							Status:   svc.Volume_Attached,
						},
					},
				},
			},
			[]*svc.AddVolumeResponse{
				nil,
			},
			[]bool{true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv, bufDialer := mockServer()
			srv.b = backend.NewMemoryBackend(backend.WithVolumes(mockVolumes))

			ctx := context.Background()
			conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
			if err != nil {
				t.Fatalf("Failed to dial bufnet: %v", err)
			}
			defer conn.Close()
			client := svc.NewVolchestratorAdminClient(conn)

			for i, req := range tt.args {
				got, err := client.AddVolume(req.ctx, req.req)
				if (err != nil) != tt.wantErr[i] {
					t.Errorf("Server.AddVolume() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !proto.Equal(got, tt.want[i]) {
					t.Errorf("Server.AddVolume() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func Test_UpdateVolume(t *testing.T) {
	mockVolumes := []*volume.Volume{
		{
			ID:     "foo",
			Region: "us-west-2",
			Tag:    "bar",
			Status: volume.Available,
		},
		{
			ID:     "bar",
			Region: "us-west-1",
			Tag:    "baz",
			Status: volume.Unavailable,
		},
	}

	type args struct {
		ctx context.Context
		req *svc.UpdateVolumeRequest
	}
	tests := []struct {
		name    string
		args    []args
		want    []*svc.UpdateVolumeResponse
		wantErr []bool
	}{
		{
			"valid update",
			[]args{
				{
					context.Background(),
					&svc.UpdateVolumeRequest{
						Volume: &svc.Volume{
							VolumeId: "foo",
							Region:   "us-west-2",
							Tag:      "baz",
							Status:   svc.Volume_Status(volume.Available),
						},
					},
				},
			},
			[]*svc.UpdateVolumeResponse{
				{},
			},
			[]bool{false},
		},
		{
			"another valid update",
			[]args{
				{
					context.Background(),
					&svc.UpdateVolumeRequest{
						Volume: &svc.Volume{
							VolumeId: "bar",
							Region:   "us-west-2",
							Tag:      "foo",
							Status:   svc.Volume_Status(volume.Available),
						},
					},
				},
			},
			[]*svc.UpdateVolumeResponse{
				{},
			},
			[]bool{false},
		},
		{
			"partial update",
			[]args{
				{
					context.Background(),
					&svc.UpdateVolumeRequest{
						Volume: &svc.Volume{
							VolumeId: "foo",
							Tag:      "baz",
						},
					},
				},
			},
			[]*svc.UpdateVolumeResponse{
				{},
			},
			[]bool{false},
		},
		{
			"nonexistent volume",
			[]args{
				{
					context.Background(),
					&svc.UpdateVolumeRequest{
						Volume: &svc.Volume{
							VolumeId: "bat",
							Region:   "us-west-2",
							Tag:      "baz",
							Status:   svc.Volume_Status(volume.Available),
						},
					},
				},
			},
			[]*svc.UpdateVolumeResponse{
				nil,
			},
			[]bool{true},
		},
		{
			"invalid transition",
			[]args{
				{
					context.Background(),
					&svc.UpdateVolumeRequest{
						Volume: &svc.Volume{
							VolumeId: "foo",
							Region:   "us-west-2",
							Tag:      "bar",
							Status:   svc.Volume_Status(volume.Attached),
						},
					},
				},
			},
			[]*svc.UpdateVolumeResponse{
				nil,
			},
			[]bool{true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv, bufDialer := mockServer()
			srv.b = backend.NewMemoryBackend(backend.WithVolumes(mockVolumes))

			ctx := context.Background()
			conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
			if err != nil {
				t.Fatalf("Failed to dial bufnet: %v", err)
			}
			defer conn.Close()
			client := svc.NewVolchestratorAdminClient(conn)

			for i, req := range tt.args {
				got, err := client.UpdateVolume(req.ctx, req.req)
				if (err != nil) != tt.wantErr[i] {
					t.Errorf("Server.UpdateVolume() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !proto.Equal(got, tt.want[i]) {
					t.Errorf("Server.UpdateVolume() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func Test_DeleteVolume(t *testing.T) {
	mockVolumes := []*volume.Volume{
		{
			ID:     "foo",
			Region: "us-west-2",
			Tag:    "bar",
			Status: volume.Unavailable,
		},
		{
			ID:     "bar",
			Region: "us-west-1",
			Tag:    "baz",
			Status: volume.Available,
		},
	}

	type args struct {
		ctx context.Context
		req *svc.DeleteVolumeRequest
	}
	tests := []struct {
		name    string
		args    []args
		want    []*svc.DeleteVolumeResponse
		wantErr []bool
	}{
		{
			"unavailable volume",
			[]args{
				{
					context.Background(),
					&svc.DeleteVolumeRequest{
						Volume: &svc.Volume{
							VolumeId: "foo",
						},
					},
				},
			},
			[]*svc.DeleteVolumeResponse{
				{},
			},
			[]bool{false},
		},
		{
			"available volume",
			[]args{
				{
					context.Background(),
					&svc.DeleteVolumeRequest{
						Volume: &svc.Volume{
							VolumeId: "bar",
						},
					},
				},
			},
			[]*svc.DeleteVolumeResponse{
				nil,
			},
			[]bool{true},
		},
		{
			"nonexistent volume",
			[]args{
				{
					context.Background(),
					&svc.DeleteVolumeRequest{
						Volume: &svc.Volume{
							VolumeId: "baz",
						},
					},
				},
			},
			[]*svc.DeleteVolumeResponse{
				nil,
			},
			[]bool{true},
		},
		{
			"delete same volume twice",
			[]args{
				{
					context.Background(),
					&svc.DeleteVolumeRequest{
						Volume: &svc.Volume{
							VolumeId: "foo",
						},
					},
				},
				{
					context.Background(),
					&svc.DeleteVolumeRequest{
						Volume: &svc.Volume{
							VolumeId: "foo",
						},
					},
				},
			},
			[]*svc.DeleteVolumeResponse{
				{},
				nil,
			},
			[]bool{false, true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv, bufDialer := mockServer()
			srv.b = backend.NewMemoryBackend(backend.WithVolumes(mockVolumes))

			ctx := context.Background()
			conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
			if err != nil {
				t.Fatalf("Failed to dial bufnet: %v", err)
			}
			defer conn.Close()
			client := svc.NewVolchestratorAdminClient(conn)

			for i, req := range tt.args {
				got, err := client.DeleteVolume(req.ctx, req.req)
				if (err != nil) != tt.wantErr[i] {
					t.Errorf("Server.DeleteVolume() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !proto.Equal(got, tt.want[i]) {
					t.Errorf("Server.DeleteVolume() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func Test_WatchNotifications(t *testing.T) {
	mockNow := time.Now()

	mockClients := []*client.Client{
		{
			ID:         "foo",
			Registered: mockNow,
			LastSeen:   mockNow,
		},
		{
			ID:         "bar",
			Registered: mockNow,
			LastSeen:   mockNow,
		},
	}

	type args struct {
		ctx context.Context
		req *svc.WatchNotificationsRequest
	}
	tests := []struct {
		name    string
		args    args
		send    []*notification.Notification
		want    []*svc.WatchNotificationsResponse
		wantErr bool
	}{
		{
			"one message for one client",
			args{
				context.Background(),
				&svc.WatchNotificationsRequest{
					ClientId: "foo",
				},
			},
			[]*notification.Notification{
				{
					ClientID: "foo",
					Message:  "bar",
				},
			},
			[]*svc.WatchNotificationsResponse{
				{
					Notification: &svc.Notification{
						Message:   "bar",
						MessageId: 1,
					},
				},
			},
			false,
		},
		{
			"one message for each client",
			args{
				context.Background(),
				&svc.WatchNotificationsRequest{
					ClientId: "foo",
				},
			},
			[]*notification.Notification{
				{
					ClientID: "foo",
					Message:  "bar",
				},
				{
					ClientID: "bar",
					Message:  "baz",
				},
			},
			[]*svc.WatchNotificationsResponse{
				{
					Notification: &svc.Notification{
						Message:   "bar",
						MessageId: 1,
					},
				},
			},
			false,
		},
		{
			"two messages for one client",
			args{
				context.Background(),
				&svc.WatchNotificationsRequest{
					ClientId: "foo",
				},
			},
			[]*notification.Notification{
				{
					ClientID: "foo",
					Message:  "bar",
				},
				{
					ClientID: "foo",
					Message:  "baz",
				},
			},
			[]*svc.WatchNotificationsResponse{
				{
					Notification: &svc.Notification{
						Message:   "bar",
						MessageId: 1,
					},
				},
				{
					Notification: &svc.Notification{
						Message:   "baz",
						MessageId: 2,
					},
				},
			},
			false,
		},
		{
			"two messages for two clients",
			args{
				context.Background(),
				&svc.WatchNotificationsRequest{
					ClientId: "foo",
				},
			},
			[]*notification.Notification{
				{
					ClientID: "foo",
					Message:  "bar",
				},
				{
					ClientID: "bar",
					Message:  "bar",
				},
				{
					ClientID: "foo",
					Message:  "baz",
				},
				{
					ClientID: "bar",
					Message:  "baz",
				},
			},
			[]*svc.WatchNotificationsResponse{
				{
					Notification: &svc.Notification{
						Message:   "bar",
						MessageId: 1,
					},
				},
				{
					Notification: &svc.Notification{
						Message:   "baz",
						MessageId: 2,
					},
				},
			},
			false,
		},
		{
			"three messages for one client",
			args{
				context.Background(),
				&svc.WatchNotificationsRequest{
					ClientId: "foo",
				},
			},
			[]*notification.Notification{
				{
					ClientID: "foo",
					Message:  "bar",
				},
				{
					ClientID: "bar",
					Message:  "bar",
				},
				{
					ClientID: "foo",
					Message:  "baz",
				},
				{
					ClientID: "foo",
					Message:  "bat",
				},
			},
			[]*svc.WatchNotificationsResponse{
				{
					Notification: &svc.Notification{
						Message:   "bar",
						MessageId: 1,
					},
				},
				{
					Notification: &svc.Notification{
						Message:   "baz",
						MessageId: 2,
					},
				},
				{
					Notification: &svc.Notification{
						Message:   "bat",
						MessageId: 3,
					},
				},
			},
			false,
		},
		{
			"nonexistent client",
			args{
				context.Background(),
				&svc.WatchNotificationsRequest{
					ClientId: "fdsafsd",
				},
			},
			nil,
			nil,
			true,
		},
		{
			"empty request",
			args{
				context.Background(),
				&svc.WatchNotificationsRequest{
					ClientId: "",
				},
			},
			nil,
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv, bufDialer := mockServer()
			srv.b = backend.NewMemoryBackend(backend.WithClients(mockClients))

			for _, n := range tt.send {
				srv.b.WriteNotification(n)
				time.Sleep(time.Millisecond * 50) // lil nap to get notifications in order
			}

			ctx := context.Background()
			conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
			if err != nil {
				t.Fatalf("Failed to dial bufnet: %v", err)
			}
			defer conn.Close()
			client := svc.NewVolchestratorClient(conn)

			stream, err := client.WatchNotifications(tt.args.ctx, tt.args.req)

			notifications := []*svc.WatchNotificationsResponse{}
			var e error

			for {
				got, err := stream.Recv()
				if err != nil {
					if err != io.EOF {
						e = err
					}

					break
				}

				notifications = append(notifications, got)

				if len(notifications) == len(tt.want) {
					break
				}
			}

			for i, want := range tt.want {
				got := notifications[i]
				if !proto.Equal(got, want) {
					t.Errorf("Server.WatchNotifications() = %v, want %v", got, want)
				}
			}

			if (e != nil) != tt.wantErr {
				t.Errorf("Server.WatchNotifications() error = %v, wantErr %v", e, tt.wantErr)
			}
		})
	}
}

func Test_WatchNotifications_Shutdown(t *testing.T) {
	mockNow := time.Now()

	mockClients := []*client.Client{
		{
			ID:         "foo",
			Registered: mockNow,
			LastSeen:   mockNow,
		},
		{
			ID:         "bar",
			Registered: mockNow,
			LastSeen:   mockNow,
		},
	}

	type args struct {
		ctx context.Context
		req *svc.WatchNotificationsRequest
	}
	tests := []struct {
		name    string
		args    args
		send    []*notification.Notification
		wantErr bool
	}{
		{
			"one message pending",
			args{
				context.Background(),
				&svc.WatchNotificationsRequest{
					ClientId: "foo",
				},
			},
			[]*notification.Notification{
				{
					ClientID: "foo",
					Message:  "bar",
				},
			},
			true,
		},
		{
			"two messages pending",
			args{
				context.Background(),
				&svc.WatchNotificationsRequest{
					ClientId: "foo",
				},
			},
			[]*notification.Notification{
				{
					ClientID: "foo",
					Message:  "bar",
				},
				{
					ClientID: "foo",
					Message:  "bar",
				},
			},
			true,
		},
		{
			"zero messages pending",
			args{
				context.Background(),
				&svc.WatchNotificationsRequest{
					ClientId: "foo",
				},
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv, bufDialer := mockServer()
			srv.b = backend.NewMemoryBackend(backend.WithClients(mockClients))

			ctx := context.Background()
			conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
			if err != nil {
				t.Fatalf("Failed to dial bufnet: %v", err)
			}
			defer conn.Close()
			client := svc.NewVolchestratorClient(conn)

			stream, _ := client.WatchNotifications(tt.args.ctx, tt.args.req)

			for _, n := range tt.send {
				srv.b.WriteNotification(n)
			}

			srv.Shutdown()

			msg, err := stream.Recv()
			if (err != nil) != tt.wantErr || err == io.EOF {
				t.Errorf("Server.WatchNotifications() error = %v, wantErr %v", err, tt.wantErr)
				t.Logf("%+v\n", msg)
			}
		})
	}
}

func Test_WatchNotifications_Deregister(t *testing.T) {
	mockNow := time.Now()

	mockClients := []*client.Client{
		{
			ID:         "foo",
			Registered: mockNow,
			LastSeen:   mockNow,
		},
		{
			ID:         "bar",
			Registered: mockNow,
			LastSeen:   mockNow,
		},
	}

	srv, bufDialer := mockServer()
	srv.b = backend.NewMemoryBackend(backend.WithClients(mockClients))

	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()
	client := svc.NewVolchestratorClient(conn)

	stream, _ := client.WatchNotifications(context.Background(), &svc.WatchNotificationsRequest{ClientId: "foo"})

	srv.b.WriteNotification(&notification.Notification{
		ClientID: "foo",
		Message:  "bar",
	})

	got, err := stream.Recv()
	if err != nil {
		t.Errorf("Server.WatchNotifications() error = %v, wantErr %v", err, nil)
	}
	want := &svc.WatchNotificationsResponse{
		Notification: &svc.Notification{
			Message:   "bar",
			MessageId: 1,
		},
	}
	if !proto.Equal(got, want) {
		t.Errorf("Server.WatchNotifications() = %v, want %v", got, want)
	}

	_, err = client.Deregister(context.Background(), &svc.DeregisterRequest{ClientId: "foo"})
	if err != nil {
		t.Errorf("Server.Deregister() error = %v, wantErr %v", err, nil)
	}
	_, err = stream.Recv()
	e, ok := status.FromError(err)
	if !ok {
		t.Errorf("stream.Recv() error = %v", err)
	}
	switch e.Code() {
	case codes.Aborted:
	default:
		t.Errorf("stream.Recv() error = %v", err)
	}
}
