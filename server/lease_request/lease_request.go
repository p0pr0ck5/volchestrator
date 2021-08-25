package leaserequest

import (
	"github.com/p0pr0ck5/volchestrator/fsm"
	"github.com/p0pr0ck5/volchestrator/server/model"
)

type Status int

type LeaseRequestError struct {
	e string
}

func newLeaseRequestError(err string) LeaseRequestError {
	return LeaseRequestError{
		e: err,
	}
}

func (e LeaseRequestError) Error() string {
	return e.e
}

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

var sm = fsm.TransitionMap{
	Pending: []fsm.Transition{
		{
			State: Fulfilling,
		},
		{
			State: Deleting,
		},
	},
	Fulfilling: []fsm.Transition{
		{
			State: Fulfilled,
		},
		{
			State: Pending,
		},
	},
	Fulfilled: []fsm.Transition{
		{
			State: Pending,
		},
	},
}

type LeaseRequest struct {
	model.Model

	ID       string `model:"immutable,required"`
	ClientID string `model:"immutable,required,depends=Client:ID"`
	Region   string `model:"immutable,required"`
	Tag      string `model:"immutable,required"`
	Status   Status `model:"required"`
}

func (l *LeaseRequest) Init() {
	l.FSM, _ = fsm.NewFSM(Pending)
	for k, vv := range sm {
		l.FSM.AddTransitions(k, vv)
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

func (l *LeaseRequest) F() *fsm.FSM {
	return l.FSM
}
