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
	"time"

	"github.com/spf13/cobra"
	"github.com/thanhpk/randstr"
	"google.golang.org/grpc"

	svc "github.com/p0pr0ck5/volchestrator/svc"
)

var serverAddress string
var clientID string

func clientRun(cmd *cobra.Command, args []string) {
	if clientID == "" {
		clientID = randstr.Hex(16)
	}

	conn, err := grpc.Dial(serverAddress, []grpc.DialOption{grpc.WithBlock(), grpc.WithInsecure()}...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()

	client := svc.NewVolchestratorClient(conn)

	t := time.NewTicker(time.Millisecond * 500)

	for {
		select {
		case <-t.C:
			_, err := client.Heartbeat(context.Background(), &svc.HeartbeatMessage{Id: clientID})
			if err != nil {
				log.Println(err)
			}
		}
	}
}

// clientCmd represents the client command
var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "Dummy volchestrator client",
	Long:  `TBD`,
	Run:   clientRun,
}

func init() {
	rootCmd.AddCommand(clientCmd)

	clientCmd.Flags().StringVarP(&serverAddress, "address", "a", "127.0.0.1:50051", "Address for the volchestrator server")
	clientCmd.Flags().StringVarP(&clientID, "client-id", "c", "", "A unique client ID")
}
