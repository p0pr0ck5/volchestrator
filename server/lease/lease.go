package lease

import (
	"github.com/p0pr0ck5/volchestrator/fsm"
	"github.com/p0pr0ck5/volchestrator/server/model"
)

type Status int

var statusMap = map[string]Status{
	"Active":   Active,
	"Deleting": Deleting,
}

const (
	unused Status = iota
	Active
	Deleting
)

var sm = fsm.TransitionMap{
	Active: []fsm.Transition{
		{
			State: Deleting,
		},
	},
}

type Lease struct {
	model.Model

	ID       string `model:"immutable,required"`
	ClientID string `model:"immutable,required,depends=Client:ID"`
	VolumeID string `model:"immutable,required,depends=Volume:ID"`
	Status   Status `model:"required"`
}

func (l *Lease) Init() {
	l.FSM, _ = fsm.NewFSM(l.Status)
	for k, vv := range sm {
		l.FSM.AddTransitions(k, vv)
	}
}

func (l *Lease) Identifier() string {
	return l.ID
}

func (l *Lease) Validate() error {
	return nil
}

func (l *Lease) ValidateTransition(m model.Base) error {
	return nil
}

func (l *Lease) SetStatus(s string) {
	l.Status = statusMap[s]
}

func (l *Lease) UpdateFSM() error {
	return l.FSM.Transition(l.Status)
}

func (l *Lease) Clone() model.Base {
	ll := *l
	return &ll
}
