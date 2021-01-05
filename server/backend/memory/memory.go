package memory

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/p0pr0ck5/volchestrator/lease"
	"github.com/p0pr0ck5/volchestrator/server"
)

// Backend implements server.Backend
type Backend struct {
	volumeMap *VolumeMap

	clientMap *ClientMap

	leaseRequestMap *LeaseRequestMap

	leaseMap *LeaseMap

	notifChMap    map[string]chan server.Notification
	notifAckChMap map[string]chan struct{}

	log *log.Logger
}

// New creates an initialized empty Backend
func New() *Backend {
	m := &Backend{
		volumeMap:       NewVolumeMap(),
		clientMap:       NewClientMap(),
		leaseRequestMap: NewLeaseRequestMap(),
		leaseMap:        NewLeaseMap(),
		notifChMap:      make(map[string]chan server.Notification),
		notifAckChMap:   make(map[string]chan struct{}),
		log:             log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile),
	}

	return m
}

/*
 *
 * Client
 *
 */

// ClientMap maps clients to their info
type ClientMap struct {
	m map[string]server.ClientInfo
	l sync.Mutex
}

// NewClientMap returns an initialized ClientMap
func NewClientMap() *ClientMap {
	m := &ClientMap{
		m: make(map[string]server.ClientInfo),
	}

	return m
}

// GetClient returns a ClientInfo for a given client id
func (m *Backend) GetClient(id string) (server.ClientInfo, error) {
	m.clientMap.l.Lock()
	defer m.clientMap.l.Unlock()

	client, ok := m.clientMap.m[id]
	if !ok {
		m.log.Printf("No client %q found in memory map\n", id)
	}

	return client, nil
}

// AddClient adds a Client to the backend if it doesn't already exist
func (m *Backend) AddClient(id string) error {
	m.clientMap.l.Lock()
	defer m.clientMap.l.Unlock()

	if _, exists := m.clientMap.m[id]; exists {
		return fmt.Errorf("Client %q already exists in memory backend", id)
	}

	m.clientMap.m[id] = server.ClientInfo{
		ID:        id,
		Status:    server.UnknownClientStatus,
		FirstSeen: time.Now(),
	}

	m.notifChMap[id] = make(chan server.Notification)

	return nil
}

// UpdateClient updates the client info for a given client
func (m *Backend) UpdateClient(id string, status server.ClientStatus) error {
	m.clientMap.l.Lock()
	defer m.clientMap.l.Unlock()

	var client server.ClientInfo
	var exists bool
	if client, exists = m.clientMap.m[id]; !exists {
		return fmt.Errorf("Client %q does not exist in memory backend", id)
	}

	client.LastSeen = time.Now()
	client.Status = status

	m.clientMap.m[id] = client

	return nil
}

// RemoveClient deletes a given client from the ClientMap
func (m *Backend) RemoveClient(id string) error {
	m.clientMap.l.Lock()
	defer m.clientMap.l.Unlock()

	delete(m.clientMap.m, id)

	close(m.notifChMap[id])
	delete(m.notifChMap, id)

	return nil
}

// Clients returns a list of server.ClientInfo
func (m *Backend) Clients(f server.ClientFilterFunc) ([]server.ClientInfo, error) {
	m.clientMap.l.Lock()
	defer m.clientMap.l.Unlock()

	var c []server.ClientInfo
	for _, ci := range m.clientMap.m {
		if f(ci) {
			c = append(c, ci)
		}
	}

	return c, nil
}

// WriteNotification writes a Notification into a channel in a blocking manner
func (m *Backend) WriteNotification(id string, n server.Notification) error {
	ch, exists := m.notifChMap[id]
	if !exists {
		return fmt.Errorf("no notification channel found for %q", id)
	}

	ch <- n
	m.notifAckChMap[n.ID] = make(chan struct{})

	return nil
}

// WatchNotifications writes notifications to a given channel
func (m *Backend) WatchNotifications(id string, ch chan<- server.Notification) error {
	notifCh, exists := m.notifChMap[id]
	if !exists {
		return fmt.Errorf("unknown id %q", id)
	}

	for n := range notifCh {
		ch <- n
	}

	return nil
}

// AckNotification acks a notification
func (m *Backend) AckNotification(id string) error {
	ch, exists := m.notifAckChMap[id]
	if !exists {
		return fmt.Errorf("failed to acknowledge %q, channel does not exist", id)
	}

	close(ch)
	delete(m.notifAckChMap, id)

	return nil
}

// WatchNotification returns a channel that will be closed when the
// notification is acked
func (m *Backend) WatchNotification(id string) (chan struct{}, error) {
	ch, exists := m.notifAckChMap[id]
	if !exists {
		return nil, fmt.Errorf("channel does not exist")
	}

	return ch, nil
}

/*
 *
 * LeaseRequest
 *
 */

// LeaseRequestMap holds information about LeaseRequests
type LeaseRequestMap struct {
	m map[string]*lease.LeaseRequest
	l sync.Mutex
}

// NewLeaseRequestMap returns an initialized LeaseRequestMap
func NewLeaseRequestMap() *LeaseRequestMap {
	m := &LeaseRequestMap{
		m: make(map[string]*lease.LeaseRequest),
	}

	return m
}

// AddLeaseRequest adds a LeaseRequest to the backend
func (m *Backend) AddLeaseRequest(request *lease.LeaseRequest) error {
	m.leaseRequestMap.l.Lock()
	defer m.leaseRequestMap.l.Unlock()

	if _, exists := m.leaseRequestMap.m[request.LeaseRequestID]; exists {
		return fmt.Errorf("Lease request %q already exists in memory backend", request.LeaseRequestID)
	}

	m.leaseRequestMap.m[request.LeaseRequestID] = request

	return nil
}

// ListLeaseRequests returns a list of lease.LeaseRequest
func (m *Backend) ListLeaseRequests(f lease.LeaseRequestFilterFunc) ([]*lease.LeaseRequest, error) {
	m.leaseRequestMap.l.Lock()
	defer m.leaseRequestMap.l.Unlock()

	var l []*lease.LeaseRequest
	for _, lr := range m.leaseRequestMap.m {
		if f(*lr) {
			l = append(l, lr)
		}
	}

	return l, nil
}

// UpdateLeaseRequest updates a LeaseRequest in the backend
func (m *Backend) UpdateLeaseRequest(request *lease.LeaseRequest) error {
	m.leaseRequestMap.l.Lock()
	defer m.leaseRequestMap.l.Unlock()

	if _, exists := m.leaseRequestMap.m[request.LeaseRequestID]; !exists {
		return fmt.Errorf("Lease request %q does not exist in memory backend", request.LeaseRequestID)
	}

	m.leaseRequestMap.m[request.LeaseRequestID] = request

	return nil
}

// DeleteLeaseRequest removes a LeaseRequest from the backend
func (m *Backend) DeleteLeaseRequest(leaseRequestID string) error {
	m.leaseRequestMap.l.Lock()
	defer m.leaseRequestMap.l.Unlock()

	if _, exists := m.leaseRequestMap.m[leaseRequestID]; !exists {
		return fmt.Errorf("Lease request %q does not exist in memory backend", leaseRequestID)
	}

	delete(m.leaseRequestMap.m, leaseRequestID)

	return nil
}

/*
 *
 * Lease
 *
 */

// LeaseMap holds information about leases
type LeaseMap struct {
	m map[string]*lease.Lease
	l sync.Mutex
}

// NewLeaseMap returns an initialized LeaseMap
func NewLeaseMap() *LeaseMap {
	m := &LeaseMap{
		m: make(map[string]*lease.Lease),
	}

	return m
}

// AddLease adds a Lease to the backend
func (m *Backend) AddLease(lease *lease.Lease) error {
	m.leaseMap.l.Lock()
	defer m.leaseMap.l.Unlock()

	if _, exists := m.leaseMap.m[lease.LeaseID]; exists {
		return fmt.Errorf("Lease %q already exists in memory backend", lease.LeaseID)
	}

	m.leaseMap.m[lease.LeaseID] = lease

	return nil
}

// ListLeases returns a list of lease.Lease
func (m *Backend) ListLeases(f lease.LeaseFilterFunc) ([]*lease.Lease, error) {
	m.leaseMap.l.Lock()
	defer m.leaseMap.l.Unlock()

	var l []*lease.Lease
	for _, lr := range m.leaseMap.m {
		if f(*lr) {
			l = append(l, lr)
		}
	}

	return l, nil
}

// UpdateLease updates a Lease in the backend
func (m *Backend) UpdateLease(lease *lease.Lease) error {
	m.leaseMap.l.Lock()
	defer m.leaseMap.l.Unlock()

	if _, exists := m.leaseMap.m[lease.LeaseID]; !exists {
		return fmt.Errorf("Lease %q does not exist in memory backend", lease.LeaseID)
	}

	m.leaseMap.m[lease.LeaseID] = lease

	return nil
}

// DeleteLease removes a Lease from the backend
func (m *Backend) DeleteLease(leaseID string) error {
	m.leaseMap.l.Lock()
	defer m.leaseMap.l.Unlock()

	if _, exists := m.leaseMap.m[leaseID]; !exists {
		return fmt.Errorf("Lease %q does not exist in memory backend", leaseID)
	}

	delete(m.leaseMap.m, leaseID)

	return nil
}

/*
 *
 * Volume
 *
 */

// VolumeMap holds information about registered volumes
type VolumeMap struct {
	m map[string]*server.Volume
	l sync.Mutex
}

// NewVolumeMap returns an initialized VolumeMap
func NewVolumeMap() *VolumeMap {
	m := &VolumeMap{
		m: make(map[string]*server.Volume),
	}

	return m
}

// GetVolume satisfies server.Backend
func (m *Backend) GetVolume(id string) (*server.Volume, error) {
	m.volumeMap.l.Lock()
	defer m.volumeMap.l.Unlock()

	v, ok := m.volumeMap.m[id]
	if !ok {
		m.log.Printf("No volume %q found in memory map\n", id)
	}

	return v, nil
}

// ListVolumes satisfies server.Backend
func (m *Backend) ListVolumes(f server.VolumeFilterFunc) ([]*server.Volume, error) {
	m.volumeMap.l.Lock()
	defer m.volumeMap.l.Unlock()

	volumes := []*server.Volume{}

	for _, volume := range m.volumeMap.m {
		if f(*volume) {
			volumes = append(volumes, volume)
		}
	}

	return volumes, nil
}

// AddVolume satisfies server.Backend
func (m *Backend) AddVolume(volume *server.Volume) error {
	m.volumeMap.l.Lock()
	defer m.volumeMap.l.Unlock()

	if _, exists := m.volumeMap.m[volume.ID]; exists {
		return fmt.Errorf("Volume %q already exists in memory backend", volume.ID)
	}

	m.volumeMap.m[volume.ID] = volume

	return nil
}

// UpdateVolume satisfies server.Backend
func (m *Backend) UpdateVolume(volume *server.Volume) error {
	m.volumeMap.l.Lock()
	defer m.volumeMap.l.Unlock()

	if _, exists := m.volumeMap.m[volume.ID]; !exists {
		return fmt.Errorf("Volume %q does not exist in memory backend", volume.ID)
	}

	m.volumeMap.m[volume.ID] = volume

	return nil
}

// DeleteVolume satisfies server.Backend
func (m *Backend) DeleteVolume(id string) error {
	m.volumeMap.l.Lock()
	defer m.volumeMap.l.Unlock()

	if _, exists := m.volumeMap.m[id]; !exists {
		return fmt.Errorf("Volume %q does not exist in memory backend", id)
	}

	delete(m.volumeMap.m, id)

	return nil
}
