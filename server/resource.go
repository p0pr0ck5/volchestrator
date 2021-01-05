package server

import "github.com/p0pr0ck5/volchestrator/lease"

// ResourceManager is responsible for managing the underlying resource represented by a
// Volume, with a given client
type ResourceManager interface {
	Associate(*lease.Lease) error
	Disassociate(*lease.Lease) error
}
