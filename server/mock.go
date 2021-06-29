package server

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/p0pr0ck5/volchestrator/svc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
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

func nowIsh() time.Time {
	t := time.Now()
	return t.Round(time.Hour)
}
