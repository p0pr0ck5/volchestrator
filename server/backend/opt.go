package backend

import (
	"github.com/p0pr0ck5/volchestrator/server/client"
	"github.com/p0pr0ck5/volchestrator/server/volume"
)

type BackendOpt func(Backend) error

func WithClients(clients []*client.Client) BackendOpt {
	return func(b Backend) error {
		for _, c := range clients {
			err := b.CreateClient(c)
			if err != nil {
				return err
			}
		}

		return nil
	}
}

func WithVolumes(clients []*volume.Volume) BackendOpt {
	return func(b Backend) error {
		for _, c := range clients {
			err := b.CreateVolume(c)
			if err != nil {
				return err
			}
		}

		return nil
	}
}
