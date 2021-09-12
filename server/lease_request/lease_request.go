package leaserequest

import (
	"github.com/p0pr0ck5/volchestrator/fsm"
	"github.com/p0pr0ck5/volchestrator/server/model"
)

type Status int

var statusMap = map[string]Status{
	"Pending":    Pending,
	"Fulfilling": Fulfilling,
	"Fulfilled":  Fulfilled,
	"Deleting":   Deleting,
}

const (
	unused Status = iota
	Pending
	Fulfilling
	Fulfilled
	Deleting
)

type LeaseRequest struct {
	model.Model

	ID       string `model:"immutable,required"`
	ClientID string `model:"immutable,required,depends=Client:ID"`
	Region   string `model:"immutable,required"`
	Tag      string `model:"immutable,required"`
	Status   Status `model:"required"`
}

func (l *LeaseRequest) Init(opts ...model.ModelOpt) {
	l.FSM, _ = fsm.NewFSM(l.Status)

	for _, opt := range opts {
		opt(&l.Model)
	}
}

func (l *LeaseRequest) Identifier() string {
	return l.ID
}

func (l *LeaseRequest) Validate() error {
	return nil
}

func (l *LeaseRequest) ValidateTransition(m model.Base) error {
	return nil
}

func (l *LeaseRequest) SetStatus(s string) {
	l.Status = statusMap[s]
}

func (l *LeaseRequest) UpdateFSM() error {
	return l.FSM.Transition(l.Status, l)
}

func (l *LeaseRequest) Clone() model.Base {
	ll := *l
	return &ll
}
