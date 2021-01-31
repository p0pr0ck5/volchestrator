package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
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

	r ResourceManager

	iterateWatch chan struct{}

	log *log.Logger
}

// NewServer creates a new Server with a given Backend
func NewServer(b Backend, r ResourceManager) *Server {
	return &Server{
		b:            b,
		r:            r,
		iterateWatch: make(chan struct{}),
		log:          log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile),
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

	go s.watchLeaseRequestIterations()
}

func (s *Server) pruneClients() {
	now := time.Now()

	deadClients, err := s.b.Clients(ClientFilterByStatus(DeadClientStatus))
	if err != nil {
		s.log.Println(err)
	}
	for _, client := range deadClients {
		d := now.Sub(client.LastSeen)
		if d > time.Second*tombstoneTTL {
			s.log.Printf("Removing %s with diff %v", client.ID, d)
			s.b.RemoveClient(client.ID)
		}
	}

	aliveClients, err := s.b.Clients(ClientFilterByStatus(AliveClientStatus))
	if err != nil {
		s.log.Println(err)
	}
	for _, client := range aliveClients {
		d := now.Sub(client.LastSeen)
		if d > time.Second*heartbeatTTL {
			s.log.Printf("Marking %s as dead with diff %v", client.ID, d)
			s.b.UpdateClient(client.ID, DeadClientStatus)
		}
	}
}

func (s *Server) pruneLeaseRequests() {
	requests, err := s.b.ListLeaseRequests(lease.LeaseRequestFilterAll)
	if err != nil {
		s.log.Println(err)
		return
	}

	now := time.Now()

	for _, request := range requests {
		if request.Expires.Before(now) {
			s.log.Println("Expiring", request.LeaseRequestID)

			s.b.DeleteLeaseRequest(request.LeaseRequestID)

			s.writeNotification(request.ClientID, NewNotification(
				LeaseRequestExpiredNotificationType,
				request.LeaseRequestID,
			))
		}
	}
}

func (s *Server) pruneLeases() {
	leases, err := s.b.ListLeases(lease.LeaseFilterAll)
	if err != nil {
		s.log.Println(err)
		return
	}

	now := time.Now()

	for _, l := range leases {
		if l.Expires.Before(now) {
			s.log.Println("Lease", l.LeaseID, "expired")

			// hook into release -> delete
			err := s.releaseLease(l)
			if err != nil {
				s.log.Println(err)
				continue
			}
		}
	}
}

func (s *Server) assignLease(l *lease.Lease) error {
	l.Status = lease.LeaseStatusAssigning
	err := s.b.UpdateLease(l)
	if err != nil {
		return err
	}

	volume, err := s.b.GetVolume(l.VolumeID)
	if err != nil {
		return err
	}

	err = s.r.Associate(l)
	if err != nil {
		return err
	}

	volume.Status = LeasedVolumeStatus
	err = s.b.UpdateVolume(volume)
	if err != nil {
		return err
	}

	l.Status = lease.LeaseStatusAssigned
	err = s.b.UpdateLease(l)
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) releaseLease(l *lease.Lease) error {
	l.Status = lease.LeaseStatusReleasing
	err := s.b.UpdateLease(l)
	if err != nil {
		return err
	}

	v, err := s.b.GetVolume(l.VolumeID)
	if err != nil {
		return err
	}

	err = s.r.Disassociate(l)
	if err != nil {
		return err
	}

	v.Status = AvailableVolumeStatus
	err = s.b.UpdateVolume(v)
	if err != nil {
		return err
	}

	err = s.b.DeleteLease(l.LeaseID)
	if err != nil {
		return err
	}

	go s.iterateLeaseRequests()

	return nil
}

// Prune cleans up various resources
func (s *Server) Prune() {
	s.pruneClients()
	s.pruneLeaseRequests()
	s.pruneLeases()
}

// Register adds a new client
func (s *Server) Register(ctx context.Context, req *svc.RegisterMessage) (*svc.Empty, error) {
	err := s.b.AddClient(req.Id)
	if err != nil {
		return nil, err
	}

	go func() {
		s.writeNotification(req.Id, NewNotification(
			UnknownNotificationType,
			fmt.Sprintf("Initial notification for client %q", req.Id),
		))
	}()

	return &svc.Empty{}, nil
}

// Deregister removes a clients and all its elements
func (s *Server) Deregister(ctx context.Context, req *svc.DeregisterMessage) (*svc.Empty, error) {
	clientID := req.Id

	// remove the leaserequests first, then all leases, then the client itself
	requests, err := s.b.ListLeaseRequests(lease.LeaseRequestFilterByClient(clientID))
	if err != nil {
		return nil, err
	}

	for _, request := range requests {
		err := s.b.DeleteLeaseRequest(request.LeaseRequestID)
		if err != nil {
			return nil, err
		}
	}

	leases, err := s.b.ListLeases(lease.LeaseFilterByClient(clientID))
	if err != nil {
		return nil, err
	}

	for _, l := range leases {
		err := s.releaseLease(l)
		if err != nil {
			return nil, err
		}
	}

	err = s.b.RemoveClient(clientID)
	if err != nil {
		return nil, err
	}

	return &svc.Empty{}, nil
}

// Heartbeat handles client HeartbeatMessages
func (s *Server) Heartbeat(ctx context.Context, m *svc.HeartbeatMessage) (*svc.HeartbeatResponse, error) {
	//s.log.Println("Seen", m.Id)

	err := s.b.UpdateClient(m.Id, AliveClientStatus)
	if err != nil {
		s.log.Println(err)
		return nil, err
	}

	requests, err := s.b.ListLeaseRequests(lease.LeaseRequestFilterByClient(m.Id))
	if err != nil {
		return nil, err
	}

	for _, request := range requests {
		request.Expires = time.Now().Add(lease.DefaultLeaseTTL)
		err = s.b.UpdateLeaseRequest(request)
		if err != nil {
			return nil, err
		}
	}

	leases, err := s.b.ListLeases(lease.LeaseFilterByClient(m.Id))
	if err != nil {
		return nil, err
	}

	for _, l := range leases {
		l.Expires = time.Now().Add(lease.DefaultLeaseTTL)
		err = s.b.UpdateLease(l)
		if err != nil {
			return nil, err
		}
	}

	res := &svc.HeartbeatResponse{
		Id: m.Id,
	}

	return res, nil
}

// WatchNotifications is called for a client to watch notifications
func (s *Server) WatchNotifications(msg *svc.NotificationWatchMessage,
	stream svc.Volchestrator_WatchNotificationsServer) error {

	ch := make(chan Notification)
	go s.b.WatchNotifications(msg.Id, ch)

	for {
		var notification Notification
		select {
		case <-stream.Context().Done():
			return nil
		case notification = <-ch:
		}

		if err := stream.Context().Err(); err != nil {
			s.log.Println(err)
			break
		}

		n := &svc.Notification{
			Id:      notification.ID,
			Type:    svc.NotificationType(notification.Type),
			Message: notification.Message,
		}

		if err := stream.Send(n); err != nil {
			s.log.Println(err)
		}
	}

	return nil
}

// Acknowledge handles an acknowledgement from the client of a Notification
func (s *Server) Acknowledge(ctx context.Context, msg *svc.Acknowledgement) (*svc.Empty, error) {
	s.log.Println("Received ack", msg.Id)

	err := s.b.AckNotification(msg.Id)
	if err != nil {
		return nil, err
	}

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
		Status:           svc.VolumeStatus(volume.Status),
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
		Status:           VolumeStatus(volume.Status),
	}

	err := s.b.AddVolume(v)
	if err != nil {
		return nil, err
	}

	go s.iterateLeaseRequests()

	return volume, nil
}

// UpdateVolume performs an in-place update of an existing volume in the backend
func (s *Server) UpdateVolume(ctx context.Context, volume *svc.Volume) (*svc.Volume, error) {
	v := &Volume{
		ID:               volume.Id,
		Tags:             volume.Tags,
		AvailabilityZone: volume.AvailabilityZone,
		Status:           VolumeStatus(volume.Status),
	}

	err := s.b.UpdateVolume(v)
	if err != nil {
		return nil, err
	}

	go s.iterateLeaseRequests()

	return volume, nil
}

// DeleteVolume deletes a volume from the backend
func (s *Server) DeleteVolume(ctx context.Context, volumeID *svc.VolumeID) (*svc.Empty, error) {
	err := s.b.DeleteVolume(volumeID.Id)

	// TODO verify status

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

	s.writeNotification(request.ClientId, NewNotification(
		LeaseRequestAckNotificationType,
		fmt.Sprintf("Received LeaseRequest submission for %+v, ID: %s", request, requestID),
	))

	go s.iterateLeaseRequests()

	return &svc.Empty{}, nil
}

// ListLeases returns all Leases in the backend
func (s *Server) ListLeases(ctx context.Context, e *svc.Empty) (*svc.LeaseList, error) {
	leases, err := s.b.ListLeases(lease.LeaseFilterAll)
	if err != nil {
		return nil, err
	}

	res := &svc.LeaseList{}
	l := []*svc.Lease{}

	for _, lease := range leases {
		e, _ := ptypes.TimestampProto(lease.Expires)
		l = append(l, &svc.Lease{
			LeaseId:  lease.LeaseID,
			ClientId: lease.ClientID,
			VolumeId: lease.VolumeID,
			Expires:  e,
			Status:   svc.LeaseStatus(lease.Status),
		})
	}

	res.Leases = l

	return res, nil
}

func (s *Server) writeNotification(id string, n Notification) Notification {
	err := s.b.WriteNotification(id, n)
	if err != nil {
		s.log.Println(err)
		return n
	}

	return n
}

func contains(needle string, haystack []string) bool {
	for _, h := range haystack {
		if needle == h {
			return true
		}
	}

	return false
}

func filterLeaseRequests(volume *Volume, requests []*lease.LeaseRequest) []*lease.LeaseRequest {
	matchedRequests := []*lease.LeaseRequest{}

	for _, request := range requests {
		match := request.VolumeAvailabilityZone == volume.AvailabilityZone &&
			contains(request.VolumeTag, volume.Tags)

		if match {
			matchedRequests = append(matchedRequests, request)
		}
	}

	return matchedRequests
}

func (s *Server) iterateLeaseRequests() {
	s.iterateWatch <- struct{}{}
}

type lrm struct {
	m map[string]bool
	l sync.Mutex
}

func (l *lrm) seen(id string) bool {
	l.l.Lock()
	defer l.l.Unlock()

	return l.m[id]
}

func (l *lrm) mark(id string) {
	l.l.Lock()
	defer l.l.Unlock()

	l.m[id] = true
}

func (s *Server) watchLeaseRequestIterations() {
	for {
		select {
		case <-s.iterateWatch:
			s.log.Println("do iterateLeaseRequests")

			volumes, err := s.b.ListVolumes(VolumeFilterByStatus(AvailableVolumeStatus))
			if err != nil {
				s.log.Println(err)
				return
			}

			requests, err := s.b.ListLeaseRequests(lease.LeaseRequestFilterAll)
			if err != nil {
				s.log.Println(err)
				return
			}

			reqMap := &lrm{
				m: make(map[string]bool),
			}

			for _, volume := range volumes {
				s.log.Println("Try to lease", volume.ID)

				// get all requests relevant to this volume
				// (e.g., search by tag and az)
				filteredRequests := filterLeaseRequests(volume, requests)

				go s.tryLease(volume, filteredRequests, reqMap)
			}
		}
	}
}

// given a list of LeaseRequest, try to find a lease
func (s *Server) tryLease(volume *Volume, requests []*lease.LeaseRequest, reqMap *lrm) {
	// set the volume status to pending
	volume.Status = LeasePendingVolumeStatus
	err := s.b.UpdateVolume(volume)
	if err != nil {
		s.log.Println(err)
	}

	for _, request := range requests {
		if reqMap.seen(request.LeaseRequestID) {
			s.log.Println("already seen", request.LeaseRequestID)
			continue
		}

		reqMap.mark(request.LeaseRequestID)

		// notify the client the lease is available
		n := s.writeNotification(request.ClientID, NewNotification(
			LeaseAvailableNotificationType,
			request.LeaseRequestID,
		))

		t := time.After(lease.LeaseAvailableAckTTL)
		ackCh, err := s.b.WatchNotification(n.ID)
		if err != nil {
			s.log.Println(err)
			continue
		}

		select {
		case <-t:
			s.log.Println("TIMEOUT")
			continue
		case <-ackCh:
			s.log.Println("we haz lease")

			l := &lease.Lease{
				LeaseID:  randstr.Hex(16),
				ClientID: request.ClientID,
				VolumeID: volume.ID,
				Expires:  time.Now().Add(lease.DefaultLeaseTTL),
				Status:   lease.LeaseStatusAssigning,
			}

			err = s.b.AddLease(l)
			if err != nil {
				s.log.Println(err)
				continue
				// TODO we're in a bad state here
			}

			err = s.b.DeleteLeaseRequest(request.LeaseRequestID)
			if err != nil {
				s.log.Println(err)
				continue
			}

			err = s.assignLease(l)
			if err != nil {
				s.log.Println(err)
				continue
			}

			m, _ := json.Marshal(l)

			s.writeNotification(request.ClientID, NewNotification(
				LeaseNotificationType,
				string(m), // format a message with the volume and lease id
			))

			return // lol
		}
	}

	// set the volume status to available as we never found a lease
	s.log.Println("Did not lease", volume.ID)
	volume.Status = AvailableVolumeStatus
	err = s.b.UpdateVolume(volume)
	if err != nil {
		s.log.Println(err)
	}
}
