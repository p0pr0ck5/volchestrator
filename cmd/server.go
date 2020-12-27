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
	"log"
	"net"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"

	"github.com/p0pr0ck5/volchestrator/server"
	"github.com/p0pr0ck5/volchestrator/server/backend/memory"
	svc "github.com/p0pr0ck5/volchestrator/svc"
)

var address string

func run(cmd *cobra.Command, args []string) {
	b := memory.New()
	s := server.NewServer(b)
	s.Init()

	listen, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Println("Starting gRPC server at", address)

	grpcServer := grpc.NewServer([]grpc.ServerOption{}...)
	svc.RegisterVolchestratorServer(grpcServer, s)
	svc.RegisterVolchestratorAdminServer(grpcServer, s)
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

	serverCmd.Flags().StringVarP(&address, "address", "a", "127.0.0.1:50051", "Listen address for the volchestrator server")
}
