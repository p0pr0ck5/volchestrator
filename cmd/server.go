package cmd

/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

import (
	"context"
	"log"
	"net"
	"sync"
	"time"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"

	svc "github.com/p0pr0ck5/volchestrator/svc"
)

const heartbeatTTL = 5 // 5 seconds

var address string

// Server implements the Volchestrator service
type Server struct {
	svc.UnimplementedVolchestratorServer

	Address string

	clientMap *ClientMap
}

// ClientStatus describes client's current status
type ClientStatus int

const (
	// Unknown indicates the client status is unknown
	Unknown ClientStatus = iota

	// Alive indicates the client is alive
	Alive

	// Dead indicate the client is dead/unresponsive
	Dead

	// Left indicates the client intentionally left
	Left
)

// ClientInfo details information about a given client
type ClientInfo struct {
	ID        string
	Status    ClientStatus
	FirstSeen time.Time
	LastSeen  time.Time
}

// ClientMap maps clients to their info
type ClientMap struct {
	m map[string]ClientInfo
	l sync.Mutex
}

// NewClientMap returns an initialized ClientMap
func NewClientMap() *ClientMap {
	m := &ClientMap{
		m: make(map[string]ClientInfo),
	}

	return m
}

// UpdateClient updates the client info for a given client
func (m *ClientMap) UpdateClient(id string, status ClientStatus) {
	m.l.Lock()
	defer m.l.Unlock()

	var client ClientInfo
	var ok bool
	if client, ok = m.m[id]; !ok {
		client = ClientInfo{
			ID:        id,
			FirstSeen: time.Now(),
		}
	}

	client.LastSeen = time.Now()
	client.Status = status

	m.m[id] = client
}

// Clients returns a list of ClientInfo
func (m *ClientMap) Clients() []ClientInfo {
	m.l.Lock()
	defer m.l.Unlock()

	var c []ClientInfo
	for _, ci := range m.m {
		c = append(c, ci)
	}

	return c
}

// Prune cleans up the client list
func (s *Server) Prune() {
	clients := s.clientMap.Clients()

	now := time.Now()
	for _, client := range clients {
		d := now.Sub(client.LastSeen)
		if d > time.Second*heartbeatTTL && client.Status == Alive {
			log.Printf("Marking %s as dead with diff %v", client.ID, d)
			s.clientMap.UpdateClient(client.ID, Dead)
		}
	}
}

// Heartbeat handles client HeartbeatMessages
func (s *Server) Heartbeat(ctx context.Context, m *svc.HeartbeatMessage) (*svc.HeartbeatResponse, error) {
	log.Println("Seen", m.Id)

	s.clientMap.UpdateClient(m.Id, Alive)

	res := &svc.HeartbeatResponse{
		Id: m.Id,
	}

	return res, nil
}

func run(cmd *cobra.Command, args []string) {
	s := &Server{
		Address:   address,
		clientMap: NewClientMap(),
	}

	go func() {
		t := time.NewTicker(time.Second * heartbeatTTL)

		for {
			select {
			case <-t.C:
				s.Prune()
			}
		}
	}()

	listen, err := net.Listen("tcp", s.Address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Println("Starting gRPC server at", s.Address)

	grpcServer := grpc.NewServer([]grpc.ServerOption{}...)
	svc.RegisterVolchestratorServer(grpcServer, s)
	grpcServer.Serve(listen)
}

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run the volchestrator server",
	Long:  `TBD`,
	Run:   run,
}

func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.Flags().StringVarP(&address, "address", "", "", "Listen address for the volchestrator server")
}
