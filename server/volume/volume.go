package volume

import (
	"github.com/p0pr0ck5/volchestrator/fsm"
	"github.com/p0pr0ck5/volchestrator/server/model"
)

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

var statusMap = map[string]Status{
	"Available":   Available,
	"Unavailable": Unavailable,
	"Attaching":   Attaching,
	"Attached":    Attached,
	"Detaching":   Detaching,
	"Deleting":    Deleting,
}

const (
	unused Status = iota
	Available
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

	ID     string `model:"immutable,required"`
	Region string `model:"required"`
	Tag    string `model:"required"`
	Status Status `model:"required"`
}

func (v *Volume) Init() {
	v.FSM, _ = fsm.NewFSM(v.Status)
	for k, vv := range sm {
		v.FSM.AddTransitions(k, vv)
	}
}

func (v *Volume) Identifier() string {
	return v.ID
}

func (v *Volume) Validate() error {
	return nil
}

func (v *Volume) ValidateTransition(m model.Base) error {
	newVolume := m.(*Volume)

	if v.Status != Available && v.Status != Unavailable {
		if v.Region != newVolume.Region || v.Tag != newVolume.Tag {
			return newVolumeError("cannot change region or tag in current state")
		}
	}

	return nil
}

func (v *Volume) SetStatus(s string) {
	v.Status = statusMap[s]
}

func (v *Volume) F() *fsm.FSM {
	return v.FSM
}
