package volume

import "errors"

type Status int

const (
	Available Status = iota
	Unavailable
	Attaching
	Attached
	Detaching
)

type Volume struct {
	ID     string
	Region string
	Tag    string
	Status Status
}

func Validate(v *Volume) error {
	if v.ID == "" {
		return errors.New("missing id")
	}

	if v.Region == "" {
		return errors.New("missing region")
	}

	if v.Tag == "" {
		return errors.New("missing tag")
	}

	return nil
}
