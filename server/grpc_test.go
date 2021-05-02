package server

import (
	"context"
	"log"
	"net"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/proto"

	"github.com/golang/protobuf/ptypes"
	"github.com/p0pr0ck5/volchestrator/server/backend"
	"github.com/p0pr0ck5/volchestrator/server/client"
	"github.com/p0pr0ck5/volchestrator/server/volume"
	"github.com/p0pr0ck5/volchestrator/svc"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

var srv *Server

func init() {
	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	srv, _ = NewServer(WithNewMemoryBackend())
	svc.RegisterVolchestratorServer(s, srv)
	svc.RegisterVolchestratorAdminServer(s, srv)
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
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

func Test_Ping(t *testing.T) {
	preregister := func(id string) {
		srv.Register(context.Background(), &svc.RegisterRequest{ClientId: id})
	}

	type args struct {
		ctx context.Context
		req *svc.PingRequest
	}
	type prefunc struct {
		f   func(string)
		arg string
	}
	tests := []struct {
		name         string
		args         []args
		preFunctions []prefunc
		want         []*svc.PingResponse
		wantErr      []bool
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
			[]prefunc{
				{preregister, "foo"},
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
						ClientId: "foo",
					},
				},
			},
			[]prefunc{},
			[]*svc.PingResponse{
				nil,
			},
			[]bool{true},
		},
		{
			"one invalid ping (different registration)",
			[]args{
				{
					context.Background(),
					&svc.PingRequest{
						ClientId: "foo",
					},
				},
			},
			[]prefunc{
				{preregister, "bar"},
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
			[]prefunc{
				{preregister, "foo"},
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
			[]prefunc{
				{preregister, "foo"},
				{preregister, "bar"},
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
						ClientId: "bar",
					},
				},
			},
			[]prefunc{
				{preregister, "foo"},
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
			[]prefunc{
				{preregister, "bar"},
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
			srv.b = backend.NewMemoryBackend()

			for _, prefunc := range tt.preFunctions {
				prefunc.f(prefunc.arg)
			}

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
