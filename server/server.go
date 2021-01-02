package server

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/thanhpk/randstr"

	"github.com/p0pr0ck5/volchestrator/lease"
	svc "github.com/p0pr0ck5/volchestrator/svc"
)

const heartbeatTTL = 5  // 5 seconds
const tombstoneTTL = 10 // 10 seconds

// Server interacts with clients to manage volume leases
type Server struct {
	svc.UnimplementedVolchestratorServer
	svc.UnimplementedVolchestratorAdminServer

	b Backend

	notifChMap map[string]chan Notification

	notifAckChMap map[string]chan struct{}
}

// NewServer creates a new Server with a given Backend
func NewServer(b Backend) *Server {
	return &Server{
		b:             b,
		notifChMap:    make(map[string]chan Notification),
		notifAckChMap: make(map[string]chan struct{}),
	}
}

// Init starts background routines
func (s *Server) Init() {
	// TODO wire up clean shutdown
	go func() {
		t := time.NewTicker(time.Second * heartbeatTTL)

		for {
			select {
			case <-t.C:
				s.Prune()
			}
		}
	}()
}

func (s *Server) pruneClients() {
	now := time.Now()

	deadClients, err := s.b.Clients(ClientFilterByStatus(DeadClientStatus))
	if err != nil {
		log.Println(err)
	}
	for _, client := range deadClients {
		d := now.Sub(client.LastSeen)
		if d > time.Second*tombstoneTTL {
			log.Printf("Removing %s with diff %v", client.ID, d)
			s.b.RemoveClient(client.ID)
			close(s.notifChMap[client.ID])
			delete(s.notifChMap, client.ID)
		}
	}

	aliveClients, err := s.b.Clients(ClientFilterByStatus(AliveClientStatus))
	if err != nil {
		log.Println(err)
	}
	for _, client := range aliveClients {
		d := now.Sub(client.LastSeen)
		if d > time.Second*heartbeatTTL {
			log.Printf("Marking %s as dead with diff %v", client.ID, d)
			s.b.UpdateClient(client.ID, DeadClientStatus)
		}
	}
}

func (s *Server) pruneLeaseRequests() {
	requests, err := s.b.ListLeaseRequests(lease.LeaseRequestFilterAll)
	if err != nil {
		log.Println(err)
	}

	now := time.Now()

	for _, request := range requests {
		log.Printf("Check %+v\n", request)
		if request.Expires.Before(now) {
			log.Println("Expiring", request.LeaseRequestID)

			s.b.DeleteLeaseRequest(request.LeaseRequestID)

			notifCh, exists := s.notifChMap[request.ClientID]
			if !exists {
				log.Println("No notification channel found for", request.ClientID)
				continue
			}

			s.writeNotification(notifCh, NewNotification(
				LeaseRequestExpiredNotificationType,
				request.LeaseRequestID,
			))
		}
	}
}

// Prune cleans up various resources
func (s *Server) Prune() {
	s.pruneClients()
	s.pruneLeaseRequests()
}

// Register adds a new client
func (s *Server) Register(ctx context.Context, req *svc.RegisterMessage) (*svc.Empty, error) {
	err := s.b.AddClient(req.Id)
	if err != nil {
		return nil, err
	}

	ch := make(chan Notification)
	s.notifChMap[req.Id] = ch

	go func() {
		s.writeNotification(ch, NewNotification(
			UnknownNotificationType,
			fmt.Sprintf("Initial notification for client %q", req.Id),
		))
	}()

	return &svc.Empty{}, nil
}

// Heartbeat handles client HeartbeatMessages
func (s *Server) Heartbeat(ctx context.Context, m *svc.HeartbeatMessage) (*svc.HeartbeatResponse, error) {
	log.Println("Seen", m.Id)

	s.b.UpdateClient(m.Id, AliveClientStatus)

	res := &svc.HeartbeatResponse{
		Id: m.Id,
	}

	return res, nil
}

// WatchNotifications is called for a client to watch notifications
func (s *Server) WatchNotifications(msg *svc.NotificationWatchMessage,
	stream svc.Volchestrator_WatchNotificationsServer) error {

	ch := s.notifChMap[msg.Id]

	if ch == nil {
		return fmt.Errorf("Unknown ID %q", msg.Id)
	}

	for notification := range ch {
		n := &svc.Notification{
			Id:      notification.ID,
			Type:    svc.NotificationType(notification.Type),
			Message: notification.Message,
		}

		if err := stream.Send(n); err != nil {
			log.Println(err)
		}
	}

	log.Printf("Notification channel for %q closed\n", msg.Id)

	return nil
}

// Acknowledge handles an acknowledgement from the client of a Notification
func (s *Server) Acknowledge(ctx context.Context, msg *svc.Acknowledgement) (*svc.Empty, error) {
	ch, exists := s.notifAckChMap[msg.Id]
	if !exists {
		err := fmt.Errorf("Failed to acknowledge %q, channel does not exist", msg.Id)
		return nil, err
	}

	log.Println("Received ack", msg.Id)

	close(ch)
	delete(s.notifAckChMap, msg.Id)

	return &svc.Empty{}, nil
}

// ListClients returns the ClientMap info
func (s *Server) ListClients(ctx context.Context, m *svc.Empty) (*svc.ClientList, error) {
	res := &svc.ClientList{}
	infos := []*svc.ClientInfo{}
	clients, err := s.b.Clients(ClientFilterAll)
	if err != nil {
		return nil, err
	}

	for _, client := range clients {
		f, _ := ptypes.TimestampProto(client.FirstSeen)
		l, _ := ptypes.TimestampProto(client.LastSeen)
		infos = append(infos, &svc.ClientInfo{
			Id:           client.ID,
			ClientStatus: svc.ClientStatus(client.Status),
			FirstSeen:    f,
			LastSeen:     l,
		})
	}

	res.Info = infos
	return res, nil
}

// GetVolume returns a Volume for a given ID, or nil if the
// volume ID is not found in the backend
func (s *Server) GetVolume(ctx context.Context, volumeID *svc.VolumeID) (*svc.Volume, error) {
	volume, err := s.b.GetVolume(volumeID.Id)
	if err != nil {
		return nil, err
	}

	if volume == nil {
		return &svc.Volume{}, nil
	}

	v := &svc.Volume{
		Id:               volume.ID,
		Tags:             volume.Tags,
		AvailabilityZone: volume.AvailabilityZone,
	}

	return v, nil
}

// ListVolumes returns all volumes currently in the backend
func (s *Server) ListVolumes(ctx context.Context, e *svc.Empty) (*svc.VolumeList, error) {
	volumes, err := s.b.ListVolumes(VolumeFilterAll)
	if err != nil {
		return nil, err
	}

	volumeList := &svc.VolumeList{
		Volumes: []*svc.Volume{},
	}

	for _, volume := range volumes {
		volumeList.Volumes = append(volumeList.Volumes, &svc.Volume{
			Id:               volume.ID,
			Tags:             volume.Tags,
			AvailabilityZone: volume.AvailabilityZone,
			Status:           svc.VolumeStatus(volume.Status),
		})
	}

	return volumeList, nil
}

// AddVolume adds a new volume to the backend
func (s *Server) AddVolume(ctx context.Context, volume *svc.Volume) (*svc.Volume, error) {
	v := &Volume{
		ID:               volume.Id,
		Tags:             volume.Tags,
		AvailabilityZone: volume.AvailabilityZone,
		Status:           UnknownVolumeStatus,
	}

	err := s.b.AddVolume(v)
	if err != nil {
		return nil, err
	}

	return volume, nil
}

// UpdateVolume performs an in-place update of an existing volume in the backend
func (s *Server) UpdateVolume(ctx context.Context, volume *svc.Volume) (*svc.Volume, error) {
	v := &Volume{
		ID:               volume.Id,
		Tags:             volume.Tags,
		AvailabilityZone: volume.AvailabilityZone,
	}

	err := s.b.UpdateVolume(v)
	if err != nil {
		return nil, err
	}

	return volume, nil
}

// DeleteVolume deletes a volume from the backend
func (s *Server) DeleteVolume(ctx context.Context, volumeID *svc.VolumeID) (*svc.Empty, error) {
	err := s.b.DeleteVolume(volumeID.Id)

	return &svc.Empty{}, err
}

// SubmitLeaseRequest adds a LeaseRequest to the backend
func (s *Server) SubmitLeaseRequest(ctx context.Context, request *svc.LeaseRequest) (*svc.Empty, error) {
	requestID := randstr.Hex(16)
	err := s.b.AddLeaseRequest(&lease.LeaseRequest{
		LeaseRequestID:         requestID,
		ClientID:               request.ClientId,
		VolumeTag:              request.Tag,
		VolumeAvailabilityZone: request.AvailabilityZone,
		Expires:                time.Now().Add(lease.DefaultLeaseTTL),
	})

	if err != nil {
		return nil, err
	}

	notifCh, exists := s.notifChMap[request.ClientId]
	if !exists {
		log.Printf("No notification channel found for %q\n", request.ClientId)
	}

	s.writeNotification(notifCh, NewNotification(
		LeaseRequestAckNotificationType,
		fmt.Sprintf("Received LeaseRequest submission for %+v, ID: %s", request, requestID),
	))

	return &svc.Empty{}, nil
}

func (s *Server) writeNotification(ch chan Notification, n Notification) {
	s.notifAckChMap[n.ID] = make(chan struct{})
	ch <- n
}

func (s *Server) iterateLeaseRequests() {
	volumes, err := s.b.ListVolumes(VolumeFilterByStatus(AvailableVolumeStatus))
	if err != nil {
		log.Println(err)
	}

	for _, volume := range volumes {
		f := func(v *Volume) []*lease.LeaseRequest {
			return []*lease.LeaseRequest{}
		}
		sort := func(v []*lease.LeaseRequest) []*lease.LeaseRequest {
			return v
		}

		// get all requests relevant to this volume
		// (e.g., search by tag and az)
		requests := f(volume)

		// sort by some sort of (configurable) fn
		// filter out blacklisted clients
		// sieve duplicate clients in the slice
		sortedRequests := sort(requests)

		go s.tryLease(volume, sortedRequests)
	}
}

// given a list of LeaseRequest, try to find a lease
func (s *Server) tryLease(volume *Volume, leases []*lease.LeaseRequest) {
	// set the volume status to pending
	volume.Status = LeasePendingVolumeStatus
	err := s.b.UpdateVolume(volume)
	if err != nil {
		log.Println(err)
	}

	for _, l := range leases {
		// notify the client the lease is available
		notifCh, exists := s.notifChMap[l.ClientID]
		if !exists {
			log.Printf("No notification channel found for %q\n", l.ClientID)
			continue
		}

		n := NewNotification(
			LeaseAvailableNotificationType,
			l.LeaseRequestID,
		)
		s.writeNotification(notifCh, n)

		t := time.After(lease.LeaseAvailableAckTTL)
		ackCh := s.notifAckChMap[n.ID]

		select {
		case <-t:
			log.Println("TIMEOUT")
			continue
		case <-ackCh:
			log.Println("we haz lease")

			s.writeNotification(notifCh, NewNotification(
				LeaseNotificationType,
				"", // format a message with the volume and lease id
			))

			volume.Status = LeasedVolumeStatus
			err = s.b.UpdateVolume(volume)
			if err != nil {
				log.Println(err)
			}

			// TODO setup heartbeat watcher for the lease

			// TODO update the LeaseRequest in the backend
		}
	}
}
