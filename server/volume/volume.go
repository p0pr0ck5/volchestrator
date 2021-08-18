package volume

import (
	"github.com/p0pr0ck5/volchestrator/fsm"
	"github.com/p0pr0ck5/volchestrator/server/model"
)

type Status int

func (s Status) Value() string {
	return ""
}

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
	Deleting
)

var sm = fsm.TransitionMap{
	Available: []fsm.Transition{
		{
			State: Unavailable,
		},
		{
			State: Attaching,
		},
	},
	Unavailable: []fsm.Transition{
		{
			State: Available,
		},
		{
			State: Deleting,
		},
	},
	Attaching: []fsm.Transition{
		{
			State: Attached,
		},
		{
			State: Detaching,
		},
	},
	Attached: []fsm.Transition{
		{
			State: Detaching,
		},
	},
	Detaching: []fsm.Transition{
		{
			State: Available,
		},
		{
			State: Unavailable,
		},
	},
}

type Volume struct {
	model.Model

	ID     string
	Region string
	Tag    string
	Status Status
}

func (v *Volume) Init() {
	v.FSM, _ = fsm.NewFSM(v.Status)
	for k, vv := range sm {
		v.FSM.AddTransitions(k, vv)
	}
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

	if v.Status != newVolume.Status && !v.FSM.Can(newVolume.Status) {
		return newVolumeError("invalid status transition")
	}

	return nil
}
