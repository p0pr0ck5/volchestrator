package client

import (
	"context"
	"log"
	"time"

	"github.com/thanhpk/randstr"
	"google.golang.org/grpc"

	"github.com/p0pr0ck5/volchestrator/config"
	svc "github.com/p0pr0ck5/volchestrator/svc"
)

type notificationHandler func(*Client, *svc.Notification) error

var notifHandlers map[svc.NotificationType][]notificationHandler

func registerNotificationHandlers() {
	notifHandlers[svc.NotificationType_NOTIFICATIONLEASEREQUESTEXPIRED] = []notificationHandler{
		func(client *Client, msg *svc.Notification) error {
			log.Println("Handling renewal of", msg.Message, "TODO")
			return nil
		},
	}

	notifHandlers[svc.NotificationType_NOTIFICATIONLEASEAVAILABLE] = []notificationHandler{
		func(client *Client, msg *svc.Notification) error {
			log.Println("I shouldn't need to do anything here right?") // because we ack the notification
			return nil
		},
	}

	notifHandlers[svc.NotificationType_NOTIFICATIONLEASE] = []notificationHandler{
		func(client *Client, msg *svc.Notification) error {
			log.Printf("I haz lease! %+v\n", msg)
			// DO THE THING
			return nil
		},
	}
}

// Client represents a volchestrator client
type Client struct {
	Config config.ClientConfig

	ClientID string

	svcClient svc.VolchestratorClient
	conn      *grpc.ClientConn
}

// NewClient create a new client object based on a given configuration
func NewClient(c config.ClientConfig) (*Client, error) {
	client := &Client{
		Config: c,
	}

	if c.ClientID == "" {
		client.ClientID = randstr.Hex(16)
	}

	return client, nil
}

// Run handles registering the client and setting up notification handlers
func (c *Client) Run() error {
	notifHandlers = make(map[svc.NotificationType][]notificationHandler)
	registerNotificationHandlers()

	conn, err := grpc.Dial(c.Config.ServerAddress, []grpc.DialOption{grpc.WithBlock(), grpc.WithInsecure()}...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}

	c.conn = conn
	c.svcClient = svc.NewVolchestratorClient(conn)

	_, err = c.svcClient.Register(context.Background(), &svc.RegisterMessage{
		Id: c.ClientID,
	})
	if err != nil {
		return err
	}

	go c.WatchNotifications()

	go func() {
		for _, request := range c.Config.LeaseRequests {
			c.svcClient.SubmitLeaseRequest(context.Background(), &svc.LeaseRequest{
				ClientId:         c.ClientID,
				Tag:              request.Tag,
				AvailabilityZone: request.AvailabilityZone,
			})
		}
	}()

	go c.SendHeartbeats()

	return nil
}

// SendHeartbeats sends heartbeats on a regular basis
func (c *Client) SendHeartbeats() {
	t := time.NewTicker(time.Millisecond * 500)

	for {
		select {
		case <-t.C:
			_, err := c.svcClient.Heartbeat(context.Background(), &svc.HeartbeatMessage{Id: c.ClientID})
			if err != nil {
				log.Println(err)
			}
		}
	}
}

// WatchNotifications handles notifications and runs callback functions
func (c *Client) WatchNotifications() {
	log.Println("Watching for notifications for client", c.ClientID)

	stream, err := c.svcClient.WatchNotifications(context.Background(), &svc.NotificationWatchMessage{
		Id: c.ClientID,
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
			err = f(c, msg)
			if err != nil {
				log.Println("Error executing callback:", err)
				continue
			}
		}

		_, err = c.svcClient.Acknowledge(context.Background(), &svc.Acknowledgement{
			Id: msg.Id,
		})
		if err != nil {
			log.Fatalln(err)
		}
		log.Println("Acknowledged", msg.Id)
	}
}

// Stop stops a client
func (c *Client) Stop() error {
	c.svcClient.Deregister(context.Background(), &svc.DeregisterMessage{
		Id: c.ClientID,
	})
	c.conn.Close()
	return nil
}
