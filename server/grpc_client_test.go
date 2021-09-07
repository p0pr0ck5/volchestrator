package server

import (
	"context"
	"io"
	"testing"
	"time"

	"github.com/p0pr0ck5/volchestrator/server/backend"
	"github.com/p0pr0ck5/volchestrator/server/client"
	"github.com/p0pr0ck5/volchestrator/server/notification"
	"github.com/p0pr0ck5/volchestrator/svc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

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
				{
					ClientId: "foo",
					Token:    "mock",
				},
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
				{
					ClientId: "foo",
					Token:    "mock",
				},
				{
					ClientId: "bar",
					Token:    "mock",
				},
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
				{
					ClientId: "foo",
					Token:    "mock",
				},
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
				{
					ClientId: "foo",
					Token:    "mock",
				},
				{
					ClientId: "bar",
					Token:    "mock",
				},
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
				{
					ClientId: "bar",
					Token:    "mock",
				},
				{
					ClientId: "foo",
					Token:    "mock",
				},
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
				{
					ClientId: "foo",
					Token:    "mock",
				},
				nil,
				{
					ClientId: "bar",
					Token:    "mock",
				},
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
						Token:    "mock",
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
						Token:    "mock",
					},
				},
				{
					context.Background(),
					&svc.DeregisterRequest{
						ClientId: "bar",
						Token:    "mock",
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
			"one invalid deregistration (invalid token)",
			[]args{
				{
					context.Background(),
					&svc.DeregisterRequest{
						ClientId: "foo",
						Token:    "nope",
					},
				},
			},
			[]*svc.DeregisterResponse{
				nil,
			},
			[]bool{true},
		},
		{
			"one invalid deregistration (empty token)",
			[]args{
				{
					context.Background(),
					&svc.DeregisterRequest{
						ClientId: "foo",
					},
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
						Token:    "mock",
					},
				},
				{
					context.Background(),
					&svc.DeregisterRequest{
						ClientId: "baz",
						Token:    "mock",
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
						Token:    "mock",
					},
				},
				{
					context.Background(),
					&svc.DeregisterRequest{
						ClientId: "foo",
						Token:    "mock",
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
						Token:    "mock",
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
			"one invalid ping (invalid token)",
			[]args{
				{
					context.Background(),
					&svc.PingRequest{
						ClientId: "foo",
						Token:    "nope",
					},
				},
			},
			[]*svc.PingResponse{
				nil,
			},
			[]bool{true},
		},
		{
			"one invalid ping (no token)",
			[]args{
				{
					context.Background(),
					&svc.PingRequest{
						ClientId: "foo",
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
						Token:    "mock",
					},
				},
				{
					context.Background(),
					&svc.PingRequest{
						ClientId: "foo",
						Token:    "mock",
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
						Token:    "mock",
					},
				},
				{
					context.Background(),
					&svc.PingRequest{
						ClientId: "bar",
						Token:    "mock",
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
						Token:    "mock",
					},
				},
				{
					context.Background(),
					&svc.PingRequest{
						ClientId: "baz",
						Token:    "mock",
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
						Token:    "mock",
					},
				},
				{
					context.Background(),
					&svc.PingRequest{
						ClientId: "bar",
						Token:    "mock",
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

func Test_WatchNotifications(t *testing.T) {
	mockNow := time.Now()

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
					Token:    "mock",
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
					Token:    "mock",
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
					Token:    "mock",
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
					Token:    "mock",
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
					Token:    "mock",
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
		{
			"bad token",
			args{
				context.Background(),
				&svc.WatchNotificationsRequest{
					ClientId: "foo",
					Token:    "nope",
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

			stream, _ := client.WatchNotifications(tt.args.ctx, tt.args.req)

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
					Token:    "mock",
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
					Token:    "mock",
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
					Token:    "mock",
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

	srv, bufDialer := mockServer()
	srv.b = backend.NewMemoryBackend(backend.WithClients(mockClients))

	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()
	client := svc.NewVolchestratorClient(conn)

	watch := &svc.WatchNotificationsRequest{
		ClientId: "foo",
		Token:    "mock",
	}
	stream, _ := client.WatchNotifications(context.Background(), watch)

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

	dereg := &svc.DeregisterRequest{
		ClientId: "foo",
		Token:    "mock",
	}
	_, err = client.Deregister(context.Background(), dereg)
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

func Test_Client_Lifecycle(t *testing.T) {
	srv, bufDialer := mockServer()
	srv.b = backend.NewMemoryBackend()

	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()
	client := svc.NewVolchestratorClient(conn)

	id := "foo"

	// register
	res, err := client.Register(context.Background(), &svc.RegisterRequest{
		ClientId: id,
	})
	if err != nil {
		t.Errorf("client.Register() unexpected error = %v", err)
	}

	token := res.Token

	// ping several times in a row
	for i := 0; i < 3; i++ {
		_, err = client.Ping(context.Background(), &svc.PingRequest{
			ClientId: id,
			Token:    token,
		})
		if err != nil {
			t.Errorf("client.Ping() unexpected error = %v", err)
		}
	}

	// deregister
	_, err = client.Deregister(context.Background(), &svc.DeregisterRequest{
		ClientId: id,
		Token:    token,
	})
	if err != nil {
		t.Errorf("client.Deregister() unexpected error = %v", err)
	}

	// cannot ping again
	for i := 0; i < 3; i++ {
		_, err = client.Ping(context.Background(), &svc.PingRequest{
			ClientId: id,
			Token:    token,
		})
		if err == nil {
			t.Errorf("client.Ping() expected error, got = %v", err)
		}
	}
}

func Test_RequestLease(t *testing.T) {
	mockNow := time.Now()

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
		req *svc.RequestLeaseRequest
	}
	tests := []struct {
		name    string
		args    []args
		want    []*svc.RequestLeaseResponse
		wantErr []bool
	}{
		{
			"one valid lease request",
			[]args{
				{
					context.Background(),
					&svc.RequestLeaseRequest{
						LeaseRequestId: "foo",
						ClientId:       "foo",
						Token:          "mock",
						Region:         "us-west-2",
						Tag:            "foo",
					},
				},
			},
			[]*svc.RequestLeaseResponse{
				{},
			},
			[]bool{false},
		},
		{
			"two valid lease requests",
			[]args{
				{
					context.Background(),
					&svc.RequestLeaseRequest{
						LeaseRequestId: "foo",
						ClientId:       "foo",
						Token:          "mock",
						Region:         "us-west-2",
						Tag:            "foo",
					},
				},
				{
					context.Background(),
					&svc.RequestLeaseRequest{
						LeaseRequestId: "bar",
						ClientId:       "foo",
						Token:          "mock",
						Region:         "us-west-2",
						Tag:            "foo",
					},
				},
			},
			[]*svc.RequestLeaseResponse{
				{},
				{},
			},
			[]bool{false, false},
		},
		{
			"one invalid lease request - bad client id",
			[]args{
				{
					context.Background(),
					&svc.RequestLeaseRequest{
						LeaseRequestId: "dne",
						ClientId:       "bad",
						Token:          "mock",
						Region:         "us-west-2",
						Tag:            "foo",
					},
				},
			},
			[]*svc.RequestLeaseResponse{
				nil,
			},
			[]bool{true},
		},
		{
			"one invalid lease request - bad client token",
			[]args{
				{
					context.Background(),
					&svc.RequestLeaseRequest{
						LeaseRequestId: "foo",
						ClientId:       "foo",
						Token:          "nope",
						Region:         "us-west-2",
						Tag:            "foo",
					},
				},
			},
			[]*svc.RequestLeaseResponse{
				nil,
			},
			[]bool{true},
		},
		{
			"one invalid lease request - missing client token",
			[]args{
				{
					context.Background(),
					&svc.RequestLeaseRequest{
						LeaseRequestId: "foo",
						ClientId:       "foo",
						Region:         "us-west-2",
						Tag:            "foo",
					},
				},
			},
			[]*svc.RequestLeaseResponse{
				nil,
			},
			[]bool{true},
		},
		{
			"one invalid lease request - missing region",
			[]args{
				{
					context.Background(),
					&svc.RequestLeaseRequest{
						LeaseRequestId: "foo",
						ClientId:       "foo",
						Token:          "mock",
						Tag:            "foo",
					},
				},
			},
			[]*svc.RequestLeaseResponse{
				nil,
			},
			[]bool{true},
		},
		{
			"one invalid lease request - missing tag",
			[]args{
				{
					context.Background(),
					&svc.RequestLeaseRequest{
						LeaseRequestId: "foo",
						ClientId:       "foo",
						Token:          "mock",
						Region:         "us-west-2",
					},
				},
			},
			[]*svc.RequestLeaseResponse{
				nil,
			},
			[]bool{true},
		},
		{
			"one valid and one invalid lease request",
			[]args{
				{
					context.Background(),
					&svc.RequestLeaseRequest{
						LeaseRequestId: "foo",
						ClientId:       "foo",
						Token:          "mock",
						Region:         "us-west-2",
						Tag:            "foo",
					},
				},
				{
					context.Background(),
					&svc.RequestLeaseRequest{
						LeaseRequestId: "foo",
						ClientId:       "dne",
						Token:          "mock",
						Region:         "us-west-2",
						Tag:            "foo",
					},
				},
			},
			[]*svc.RequestLeaseResponse{
				{},
				nil,
			},
			[]bool{false, true},
		},
		{
			"duplicate lease requests",
			[]args{
				{
					context.Background(),
					&svc.RequestLeaseRequest{
						LeaseRequestId: "foo",
						ClientId:       "foo",
						Token:          "mock",
						Region:         "us-west-2",
						Tag:            "foo",
					},
				},
				{
					context.Background(),
					&svc.RequestLeaseRequest{
						LeaseRequestId: "foo",
						ClientId:       "foo",
						Token:          "mock",
						Region:         "us-west-2",
						Tag:            "foo",
					},
				},
			},
			[]*svc.RequestLeaseResponse{
				{},
				nil,
			},
			[]bool{false, true},
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
				got, err := client.RequestLease(req.ctx, req.req)
				if (err != nil) != tt.wantErr[i] {
					t.Errorf("Server.RequestLease() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !proto.Equal(got, tt.want[i]) {
					t.Errorf("Server.RequestLease() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
