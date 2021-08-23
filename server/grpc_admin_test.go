package server

import (
	"context"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"

	"github.com/p0pr0ck5/volchestrator/server/backend"
	"github.com/p0pr0ck5/volchestrator/server/client"
	"github.com/p0pr0ck5/volchestrator/server/volume"
	"github.com/p0pr0ck5/volchestrator/svc"
)

func Test_GetClient(t *testing.T) {
	mockNow := time.Now()
	mockTS, _ := ptypes.TimestampProto(mockNow)

	mockClients := []*client.Client{
		{
			ID:         "foo",
			Token:      "mock",
			Registered: mockNow,
			LastSeen:   mockNow,
		},
		{
			ID:         "bar",
			Token:      "mock",
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
			Token:      "mock",
			Registered: mockNow,
			LastSeen:   mockNow,
		},
		{
			ID:         "bar",
			Token:      "mock",
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
							Status:   svc.Volume_Available,
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
							Status:   svc.Volume_Available,
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

func Test_Volume_Lifecycle(t *testing.T) {
	srv, bufDialer := mockServer()
	srv.b = backend.NewMemoryBackend()

	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()
	client := svc.NewVolchestratorAdminClient(conn)

	// create a volume
	_, err = client.AddVolume(context.Background(), &svc.AddVolumeRequest{
		Volume: &svc.Volume{
			VolumeId: "foo",
			Region:   "us-west-2",
			Tag:      "foo",
			Status:   svc.Volume_Available,
		},
	})
	if err != nil {
		t.Errorf("client.AddVolume() unexpected error = %v", err)
	}

	// update its statuses in a 'regular' lifecycle.
	// each one of these should be a valid state tranistion
	statuses := []svc.Volume_Status{
		svc.Volume_Unavailable,
		svc.Volume_Available,
		svc.Volume_Attaching,
		svc.Volume_Detaching,
		svc.Volume_Available,
		svc.Volume_Attaching,
		svc.Volume_Attached,
		svc.Volume_Detaching,
	}

	for _, status := range statuses {
		_, err = client.UpdateVolume(context.Background(), &svc.UpdateVolumeRequest{
			Volume: &svc.Volume{
				VolumeId: "foo",
				Status:   status,
			},
		})
		if err != nil {
			t.Errorf("client.UpdateVolume() unexpected error = %v", err)
		}
	}

	// attempt to delete while still detaching
	res, err := client.GetVolume(context.Background(), &svc.GetVolumeRequest{
		VolumeId: "foo",
	})
	if err != nil {
		t.Errorf("client.Get() unexpected error = %v", err)

	}
	if res.Volume.Status != svc.Volume_Detaching {
		t.Errorf("client.Get() unexpected status = %v, want %v", err, svc.Volume_Detaching)
	}

	_, err = client.DeleteVolume(context.Background(), &svc.DeleteVolumeRequest{
		Volume: &svc.Volume{
			VolumeId: "foo",
		},
	})
	if err == nil {
		t.Errorf("client.Delete() expected error, got = %v", err)
	}

	// set as state from which we can delete
	_, err = client.UpdateVolume(context.Background(), &svc.UpdateVolumeRequest{
		Volume: &svc.Volume{
			VolumeId: "foo",
			Status:   svc.Volume_Unavailable,
		},
	})
	if err != nil {
		t.Errorf("client.UpdateVolume() unexpected error = %v", err)
	}

	// delete
	_, err = client.DeleteVolume(context.Background(), &svc.DeleteVolumeRequest{
		Volume: &svc.Volume{
			VolumeId: "foo",
		},
	})
	if err != nil {
		t.Errorf("client.Delete() unexpected error = %v", err)
	}

	res, err = client.GetVolume(context.Background(), &svc.GetVolumeRequest{
		VolumeId: "foo",
	})
	if err == nil {
		t.Errorf("client.Get() expected error, got = %v", err)
	}
	if res != nil {
		t.Errorf("client.Get() unexpected res = %v", res)
	}

	volumes, err := client.ListVolumes(context.Background(), &svc.ListVolumesRequest{})
	if err != nil {
		t.Errorf("client.List() unexpected error, got = %v", err)
	}
	if len(volumes.Volumes) != 0 {
		t.Errorf("client.List() unexpected Volumes = %v", volumes)
	}
}
