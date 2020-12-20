package server

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/golang/protobuf/ptypes"

	svc "github.com/p0pr0ck5/volchestrator/svc"
)

const heartbeatTTL = 5  // 5 seconds
const tombstoneTTL = 30 // 30 seconds

// Volume represents a definition of an EBS volume in the data store
type Volume struct {
	ID               string
	Tags             []string
	AvailabilityZone string
}

// Server interacts with clients to manage volume leases
type Server struct {
	svc.UnimplementedVolchestratorServer
	svc.UnimplementedVolchestratorAdminServer

	b Backend

	clientMap *ClientMap
}

// Backend defines functions implemented by the data store
type Backend interface {
	GetVolume(id string) (*Volume, error)
	ListVolumes() ([]*Volume, error)
	AddVolume(*Volume) error
	UpdateVolume(*Volume) error
	DeleteVolume(id string) error
}

// NewServer creates a new Server with a given Backend
func NewServer(b Backend) *Server {
	return &Server{
		b:         b,
		clientMap: NewClientMap(),
	}
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

// ClientFilterFunc is a function to filter a list of clients based on a given condition
type ClientFilterFunc func(ClientInfo) bool

// ClientFilterAll returns all clients
func ClientFilterAll(ci ClientInfo) bool {
	return true
}

// ClientFilterByStatus returns clients that match a given status
func ClientFilterByStatus(status ClientStatus) ClientFilterFunc {
	return func(ci ClientInfo) bool {
		return ci.Status == status
	}
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

// RemoveClient deletes a given client from the ClientMap
func (m *ClientMap) RemoveClient(id string) {
	m.l.Lock()
	defer m.l.Unlock()

	delete(m.m, id)
}

// Clients returns a list of ClientInfo
func (m *ClientMap) Clients(f ClientFilterFunc) []ClientInfo {
	m.l.Lock()
	defer m.l.Unlock()

	var c []ClientInfo
	for _, ci := range m.m {
		if f(ci) {
			c = append(c, ci)
		}
	}

	return c
}

// Prune cleans up the client list
func (s *Server) Prune() {
	now := time.Now()

	deadClients := s.clientMap.Clients(ClientFilterByStatus(Dead))
	for _, client := range deadClients {
		d := now.Sub(client.LastSeen)
		if d > time.Second*tombstoneTTL {
			log.Printf("Removing %s with diff %v", client.ID, d)
			s.clientMap.RemoveClient(client.ID)
		}
	}

	aliveClients := s.clientMap.Clients(ClientFilterByStatus(Alive))
	for _, client := range aliveClients {
		d := now.Sub(client.LastSeen)
		if d > time.Second*heartbeatTTL {
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

// ListClients returns the ClientMap info
func (s *Server) ListClients(ctx context.Context, m *svc.Empty) (*svc.ClientList, error) {
	res := &svc.ClientList{}
	infos := []*svc.ClientInfo{}
	clients := s.clientMap.Clients(ClientFilterAll)

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

// Init starts background routines
func (s *Server) Init() {
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
	volumes, err := s.b.ListVolumes()
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
