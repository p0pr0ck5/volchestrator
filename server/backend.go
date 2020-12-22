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
	AddClient(string) error
	UpdateClient(string, ClientStatus) error
	RemoveClient(string) error
	Clients(ClientFilterFunc) ([]ClientInfo, error)
}

// LeaseInterface defines functions for managing volume leases
type LeaseInterface interface {
	AddLeaseRequest(*lease.LeaseRequest) error
}

// VolumeInterface defines functions for managing volumes
type VolumeInterface interface {
	GetVolume(string) (*Volume, error)
	ListVolumes() ([]*Volume, error)
	AddVolume(*Volume) error
	UpdateVolume(*Volume) error
	DeleteVolume(string) error
}
