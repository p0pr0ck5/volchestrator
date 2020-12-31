package server

// VolumeStatus describes a volume's current status
type VolumeStatus int

// Volume represents a definition of an EBS volume in the data store
type Volume struct {
	ID               string
	Tags             []string
	AvailabilityZone string
	Status           VolumeStatus
}

const (
	// UnknownVolumeStatus indicates the volume status is unknown
	UnknownVolumeStatus VolumeStatus = iota

	// AvailableVolumeStatus indicates the volume is available for leasing
	AvailableVolumeStatus

	// LeasePendingVolumeStatus indicates the volume is pending a lease acquisition from a client
	LeasePendingVolumeStatus

	// LeasedVolumeStatus indicates the volume is currently leased by a client
	LeasedVolumeStatus
)

// VolumeFilterFunc is a function to filter a list of Volumes
// based on a given condition
type VolumeFilterFunc func(Volume) bool

// VolumeFilterAll returns all Volumes
func VolumeFilterAll(v Volume) bool {
	return true
}

// VolumeFilterByStatus returns volumes with a status of AvailableVolumeStatus
func VolumeFilterByStatus(status VolumeStatus) VolumeFilterFunc {
	return func(v Volume) bool {
		return v.Status == status
	}
}
