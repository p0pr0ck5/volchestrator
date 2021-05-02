package backend

import (
	"github.com/p0pr0ck5/volchestrator/server/backend/memory"
	"github.com/p0pr0ck5/volchestrator/server/backend/mock"
	"github.com/pkg/errors"

	"github.com/p0pr0ck5/volchestrator/server/client"
	"github.com/p0pr0ck5/volchestrator/server/volume"
)

type backend interface {
	ReadClient(id string) (*client.Client, error)
	ListClients() ([]*client.Client, error)
	CreateClient(*client.Client) error
	UpdateClient(*client.Client) error
	DeleteClient(*client.Client) error

	ReadVolume(id string) (*volume.Volume, error)
	ListVolumes() ([]*volume.Volume, error)
	CreateVolume(*volume.Volume) error
	UpdateVolume(*volume.Volume) error
	DeleteVolume(*volume.Volume) error
}

type Backend struct {
	b backend
}

func NewBackend(opts ...BackendOpt) *Backend {
	b := &Backend{}

	for _, opt := range opts {
		opt(b)
	}

	return b
}

func NewMemoryBackend(opts ...BackendOpt) *Backend {
	m := memory.NewMemoryBackend()

	b := &Backend{
		b: m,
	}

	for _, opt := range opts {
		opt(b)
	}

	return b
}

func NewMockBackend() *Backend {
	b := &Backend{
		b: mock.NewMockBackend(),
	}

	return b
}

func (b *Backend) ReadClient(id string) (*client.Client, error) {
	return b.b.ReadClient(id)
}

func (b *Backend) ListClients() ([]*client.Client, error) {
	return b.b.ListClients()
}

func (b *Backend) CreateClient(c *client.Client) error {
	return b.b.CreateClient(c)
}

func (b *Backend) UpdateClient(c *client.Client) error {
	return b.b.UpdateClient(c)
}

func (b *Backend) DeleteClient(c *client.Client) error {
	return b.b.DeleteClient(c)
}

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

	return b.b.UpdateVolume(v)
}

func (b *Backend) DeleteVolume(v *volume.Volume) error {
	return b.b.DeleteVolume(v)
}
