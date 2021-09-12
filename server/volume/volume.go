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

type Volume struct {
	model.Model

	ID      string `model:"immutable,required"`
	LeaseID string `model:"reference=Lease:ID"`
	Region  string `model:"required"`
	Tag     string `model:"required"`
	Status  Status `model:"required"`
}

func (v *Volume) Init(opts ...model.ModelOpt) {
	v.FSM, _ = fsm.NewFSM(v.Status)

	for _, opt := range opts {
		opt(&v.Model)
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

func (v *Volume) UpdateFSM() error {
	return v.FSM.Transition(v.Status)
}

func (v *Volume) Clone() model.Base {
	vv := *v
	return &vv
}
