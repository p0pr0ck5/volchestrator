package wrapper

import (
	"log"
	"net"
	"os"

	"github.com/p0pr0ck5/volchestrator/config"
	"github.com/p0pr0ck5/volchestrator/server"
	"github.com/p0pr0ck5/volchestrator/server/backend/memory"
	"github.com/p0pr0ck5/volchestrator/server/resource/timednop"
	svc "github.com/p0pr0ck5/volchestrator/svc"
	"google.golang.org/grpc"
)

// Wrapper handles the gRPC server elements and coordinating process signals with a Server
type Wrapper struct {
	Config config.ServerConfig

	Server *server.Server

	g *grpc.Server

	log *log.Logger
}

// NewWrapper creates a Wrapper based on a given config
func NewWrapper(c config.ServerConfig) (*Wrapper, error) {
	// TODO by config
	b := memory.New()
	r := timednop.New()
	s := server.NewServer(b, r)
	s.Init()

	w := &Wrapper{
		Config: c,
		Server: s,
		log:    log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile),
	}

	return w, nil
}

// Start sets up the listening routines and exits immediately
func (w *Wrapper) Start() error {
	address := w.Config.Listen.Address
	listen, err := net.Listen("tcp", address)
	if err != nil {
		w.log.Fatalf("failed to listen: %v", err)
	}

	w.log.Println("Starting gRPC server at", address)

	grpcServer := grpc.NewServer([]grpc.ServerOption{}...)
	svc.RegisterVolchestratorServer(grpcServer, w.Server)
	svc.RegisterVolchestratorAdminServer(grpcServer, w.Server)
	w.g = grpcServer
	go grpcServer.Serve(listen) // TODO cleanup

	return nil
}

// Stop gracefully stops the underlying gRPC server
func (w *Wrapper) Stop() error {
	w.log.Println("Stopping server")

	w.g.GracefulStop()

	return nil
}
