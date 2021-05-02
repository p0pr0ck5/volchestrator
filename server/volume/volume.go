package volume

type Status int

type VolumeError struct {
	e string
}

func newVolumeError(err string) VolumeError {
	return VolumeError{
		e: err,
	}
}

func (e VolumeError) Error() string {
	return e.e
}

const (
	Available Status = iota
	Unavailable
	Attaching
	Attached
	Detaching
)

var validStatusTransition = map[Status][]Status{
	Available:   {Unavailable, Attaching},
	Unavailable: {Available},
	Attaching:   {Attached, Detaching},
	Attached:    {Detaching},
	Detaching:   {Available, Unavailable},
}

func contains(needle Status, haystack []Status) bool {
	for _, s := range haystack {
		if needle == s {
			return true
		}
	}

	return false
}

type Volume struct {
	ID     string
	Region string
	Tag    string
	Status Status
}

func (v *Volume) Validate() error {
	if v.ID == "" {
		return newVolumeError("missing id")
	}

	if v.Region == "" {
		return newVolumeError("missing region")
	}

	if v.Tag == "" {
		return newVolumeError("missing tag")
	}

	return nil
}

func (v *Volume) ValidateTransition(newVolume *Volume) error {
	if v.ID != newVolume.ID {
		return newVolumeError("cannot change id")
	}

	if v.Status != Available && v.Status != Unavailable {
		if v.Region != newVolume.Region || v.Tag != newVolume.Tag {
			return newVolumeError("cannot change region or tag in current state")
		}
	}

	if v.Status != newVolume.Status && !contains(newVolume.Status, validStatusTransition[v.Status]) {
		return newVolumeError("invalid status transition")
	}

	return nil
}
