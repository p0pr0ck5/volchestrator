package backend

import (
	"time"

	"github.com/p0pr0ck5/volchestrator/server/volume"
	"github.com/pkg/errors"
)

func (b *Backend) ReadVolume(id string) (*volume.Volume, error) {
	return b.b.ReadVolume(id)
}

func (b *Backend) ListVolumes() ([]*volume.Volume, error) {
	return b.b.ListVolumes()
}

func (b *Backend) CreateVolume(v *volume.Volume) error {
	if err := v.Validate(); err != nil {
		var errMsg string

		_, ok := err.(volume.VolumeError)
		if ok {
			errMsg = "volume validation"
		} else {
			errMsg = "add volume"
		}

		return errors.Wrap(err, errMsg)
	}

	now := time.Now()
	v.CreatedAt = now
	v.UpdatedAt = now

	return b.b.CreateVolume(v)
}

func (b *Backend) UpdateVolume(v *volume.Volume) error {
	currentVolume, err := b.ReadVolume(v.ID)
	if err != nil {
		return errors.Wrap(err, "get current volume")
	}

	if err := v.Validate(); err != nil {
		var errMsg string

		_, ok := err.(volume.VolumeError)
		if ok {
			errMsg = "volume validation"
		} else {
			errMsg = "update volume"
		}

		return errors.Wrap(err, errMsg)
	}

	if err := currentVolume.ValidateTransition(v); err != nil {
		var errMsg string

		_, ok := err.(volume.VolumeError)
		if ok {
			errMsg = "volume validation"
		} else {
			errMsg = "update volume"
		}

		return errors.Wrap(err, errMsg)
	}

	v.UpdatedAt = time.Now()

	return b.b.UpdateVolume(v)
}

func (b *Backend) DeleteVolume(v *volume.Volume) error {
	currentVolume, err := b.ReadVolume(v.ID)
	if err != nil {
		return errors.Wrap(err, "get current volume")
	}

	if currentVolume.Status != volume.Unavailable {
		return errors.New("cannot delete volume when it is not unavailable")
	}

	return b.b.DeleteVolume(v)
}