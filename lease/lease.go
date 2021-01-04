package lease

import "time"

// DefaultLeaseTTL is the length LeaseRequests can stay alive before they need to be replaced
const DefaultLeaseTTL = time.Duration(15) * time.Second // 60 seconds

// LeaseAvailableAckTTL defines how long to wait for a client to ack a LeaseAvailable notification
const LeaseAvailableAckTTL = time.Duration(5) * time.Second // 5 seconds

// LeaseRequest represets a client's desire to lease a given volume
type LeaseRequest struct {
	LeaseRequestID         string
	ClientID               string
	VolumeTag              string
	VolumeAvailabilityZone string
	Expires                time.Time
}

// LeaseRequestFilterFunc is a function to filter a list of LeaseRequests
// based on a given condition
type LeaseRequestFilterFunc func(LeaseRequest) bool

// LeaseRequestFilterAll returns all LeaseRequests
func LeaseRequestFilterAll(l LeaseRequest) bool {
	return true
}

// LeaseRequestFilterByClient returns all LeaseRequests for a given client
func LeaseRequestFilterByClient(id string) LeaseRequestFilterFunc {
	return func(l LeaseRequest) bool {
		return l.ClientID == id
	}
}

type LeaseStatus int

const (
	LeaseStatusUnknown LeaseStatus = iota

	LeaseStatusAssigning

	LeaseStatusAssigned

	LeaseStatusReleasing
)

// Lease represents a lease of a volume to a client, for a given period of time
type Lease struct {
	LeaseID  string
	ClientID string
	VolumeID string
	Expires  time.Time
	Status   LeaseStatus
}

// LeaseFilterFunc is a function to filter a list of Leases based on a given condition
type LeaseFilterFunc func(Lease) bool

// LeaseFilterAll returns all Leases
func LeaseFilterAll(l Lease) bool {
	return true
}

//LeaseFilterByClient returns all Leases for a given client
func LeaseFilterByClient(id string) LeaseFilterFunc {
	return func(l Lease) bool {
		return l.ClientID == id
	}
}
