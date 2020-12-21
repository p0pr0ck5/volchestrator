package lease

import "time"

const defaultLeaseTTL = 60 // 60 seconds

// LeaseRequest represets a client's desire to lease a given volume
type LeaseRequest struct {
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
