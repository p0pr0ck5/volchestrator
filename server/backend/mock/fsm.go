package mock

import (
	"github.com/p0pr0ck5/volchestrator/fsm"
	"github.com/p0pr0ck5/volchestrator/server/client"
	"github.com/p0pr0ck5/volchestrator/server/lease"
	lease_request "github.com/p0pr0ck5/volchestrator/server/lease_request"
	"github.com/p0pr0ck5/volchestrator/server/model"
	"github.com/p0pr0ck5/volchestrator/server/volume"
)

type SMMap map[string]fsm.TransitionMap

func (b *MockBackend) BuildSMMap() SMMap {
	return map[string]fsm.TransitionMap{
		"Client": {
			client.Alive: []fsm.Transition{
				{
					State: client.Deleting,
				},
			},
		},
		"Lease": {
			lease.Active: []fsm.Transition{
				{
					State: lease.Deleting,
					Callback: func(e fsm.Event) error {
						// reset the volume associated with the lease
						lease := e.Args[0].(*lease.Lease)

						var volumes []model.Base
						if err := b.List("Volume", &volumes); err != nil {
							return err
						}

						for _, v := range volumes {
							v := v.(*volume.Volume)

							if v.LeaseID != lease.ID {
								continue
							}

							v.Status = volume.Detaching
							if err := b.Update(v); err != nil {
								return err
							}

							v.Status = volume.Available
							if err := b.Update(v); err != nil {
								return err
							}
						}

						return nil
					},
				},
			},
		},
		"LeaseRequest": {
			lease_request.Pending: []fsm.Transition{
				{
					State: lease_request.Fulfilling,
				},
				{
					State: lease_request.Deleting,
				},
			},
			lease_request.Fulfilling: []fsm.Transition{
				{
					State: lease_request.Fulfilled,
				},
				{
					State: lease_request.Pending,
				},
			},
			lease_request.Fulfilled: []fsm.Transition{
				{
					State: lease_request.Pending,
				},
			},
		},
		"Volume": {
			volume.Available: []fsm.Transition{
				{
					State: volume.Unavailable,
				},
				{
					State: volume.Attaching,
				},
			},
			volume.Unavailable: []fsm.Transition{
				{
					State: volume.Available,
				},
				{
					State: volume.Deleting,
				},
			},
			volume.Attaching: []fsm.Transition{
				{
					State: volume.Attached,
				},
				{
					State: volume.Detaching,
				},
			},
			volume.Attached: []fsm.Transition{
				{
					State: volume.Detaching,
				},
			},
			volume.Detaching: []fsm.Transition{
				{
					State: volume.Available,
				},
				{
					State: volume.Unavailable,
				},
			},
		},
	}
}
