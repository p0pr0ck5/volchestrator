package server

import (
	"github.com/p0pr0ck5/volchestrator/lease"
)

// Backend defines functions implemented by the data store
type Backend interface {
	ClientInterface
	LeaseInterface
	VolumeInterface
}

// ClientInterface defines functions for managing clients
type ClientInterface interface {
	GetClient(string) (ClientInfo, error)
	AddClient(string) error
	UpdateClient(string, ClientStatus) error
	RemoveClient(string) error
	Clients(ClientFilterFunc) ([]ClientInfo, error)

	WriteNotification(string, Notification) error
	WatchNotifications(string, chan<- Notification) error
	AckNotification(string) error
	WatchNotification(string) (chan struct{}, error)
}

// LeaseInterface defines functions for managing volume leases
type LeaseInterface interface {
	ListLeaseRequests(lease.LeaseRequestFilterFunc) ([]*lease.LeaseRequest, error)
	AddLeaseRequest(*lease.LeaseRequest) error
	UpdateLeaseRequest(*lease.LeaseRequest) error
	DeleteLeaseRequest(string) error

	AddLease(*lease.Lease) error
	ListLeases(lease.LeaseFilterFunc) ([]*lease.Lease, error)
	UpdateLease(*lease.Lease) error
	DeleteLease(string) error
}

// VolumeInterface defines functions for managing volumes
type VolumeInterface interface {
	GetVolume(string) (*Volume, error)
	ListVolumes(VolumeFilterFunc) ([]*Volume, error)
	AddVolume(*Volume) error
	UpdateVolume(*Volume) error
	DeleteVolume(string) error
}
