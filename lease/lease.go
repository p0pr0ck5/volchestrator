package lease

import "time"

const defaultLeaseTTL = 60 // 60 seconds

// LeaseRequest represets a client's desire to lease a given volume
type LeaseRequest struct {
	LeaseRequestID         string
	ClientID               string
	VolumeTag              string
	VolumeAvailabilityZone string
	TTL                    time.Duration
}

// Lease represents a lease of a volume to a client, for a given period of time
type Lease struct {
	LeaseID  string
	ClientID string
	VolumeID string
	Expires  time.Time
}

// LeaseRequestFilterFunc is a function to filter a list of LeaseRequests
// based on a given condition
type LeaseRequestFilterFunc func(LeaseRequest) bool

// LeaseRequestFilterAll returns all LeaseRequests
func LeaseRequestFilterAll(l LeaseRequest) bool {
	return true
}
