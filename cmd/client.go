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

type notificationHandler func(svc.VolchestratorClient, *svc.Notification) error

var notifHandlers map[svc.NotificationType][]notificationHandler

func watch(client svc.VolchestratorClient) {
	log.Println("Watching for notifications for client", clientID)

	stream, err := client.WatchNotifications(context.Background(), &svc.NotificationWatchMessage{
		Id: clientID,
	})
	if err != nil {
		log.Fatalln(err)
	}

	for {
		msg, err := stream.Recv()
		if err != nil {
			log.Fatalln(err)
		}

		log.Printf("Received notification: '%+v'\n", msg)

		for _, f := range notifHandlers[msg.Type] {
			err = f(client, msg)
			if err != nil {
				log.Println("Error executing callback:", err)
				continue
			}
		}

		_, err = client.Acknowledge(context.Background(), &svc.Acknowledgement{
			Id: msg.Id,
		})
		if err != nil {
			log.Fatalln(err)
		}
		log.Println("Acknowledged", msg.Id)
	}
}

func registerNotificationHandlers() {
	notifHandlers[svc.NotificationType_NOTIFICATIONLEASEREQUESTEXPIRED] = []notificationHandler{
		func(client svc.VolchestratorClient, msg *svc.Notification) error {
			log.Println("Handling renewal of", msg.Message)
			client.SubmitLeaseRequest(context.Background(), &svc.LeaseRequest{
				ClientId:         clientID,
				Tag:              "foo",
				AvailabilityZone: "us-west-2a",
			})
			return nil
		},
	}

	notifHandlers[svc.NotificationType_NOTIFICATIONLEASEAVAILABLE] = []notificationHandler{
		func(client svc.VolchestratorClient, msg *svc.Notification) error {
			log.Println("I shouldn't need to do anything here right?") // because we ack the notification
			return nil
		},
	}

	notifHandlers[svc.NotificationType_NOTIFICATIONLEASE] = []notificationHandler{
		func(client svc.VolchestratorClient, msg *svc.Notification) error {
			log.Printf("I haz lease! %+v\n", msg)
			// DO THE THING
			return nil
		},
	}
}

func clientRun(cmd *cobra.Command, args []string) {
	notifHandlers = make(map[svc.NotificationType][]notificationHandler)
	registerNotificationHandlers()

	if clientID == "" {
		clientID = randstr.Hex(16)
	}

	conn, err := grpc.Dial(serverAddress, []grpc.DialOption{grpc.WithBlock(), grpc.WithInsecure()}...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()

	client := svc.NewVolchestratorClient(conn)

	_, err = client.Register(context.Background(), &svc.RegisterMessage{
		Id: clientID,
	})
	if err != nil {
		log.Fatalln(err)
	}

	go watch(client)

	go func() {
		client.SubmitLeaseRequest(context.Background(), &svc.LeaseRequest{
			ClientId:         clientID,
			Tag:              "foo",
			AvailabilityZone: "us-west-2a",
		})
	}()

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
