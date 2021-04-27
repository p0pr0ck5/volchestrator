package backend

import (
	"github.com/p0pr0ck5/volchestrator/server/client"
	"github.com/p0pr0ck5/volchestrator/server/volume"
)

type Backend interface {
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
