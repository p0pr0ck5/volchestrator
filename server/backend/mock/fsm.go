package mock

import (
	"github.com/p0pr0ck5/volchestrator/fsm"
	"github.com/p0pr0ck5/volchestrator/server/client"
	"github.com/p0pr0ck5/volchestrator/server/lease"
	lease_request "github.com/p0pr0ck5/volchestrator/server/lease_request"
	"github.com/p0pr0ck5/volchestrator/server/volume"
)

var SMMap = map[string]fsm.TransitionMap{
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
