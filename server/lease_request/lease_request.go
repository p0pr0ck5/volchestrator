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

	ID       string
	ClientID string
	Region   string
	Tag      string
	Status   Status
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
	if l.ID == "" {
		return newLeaseRequestError("missing id")
	}

	if l.ClientID == "" {
		return newLeaseRequestError("missing client id")
	}

	if l.Region == "" {
		return newLeaseRequestError("missing region")
	}

	if l.Tag == "" {
		return newLeaseRequestError("missing tag")
	}

	if l.Status == unused {
		return newLeaseRequestError("invalid status")
	}

	return nil
}

func (l *LeaseRequest) ValidateTransition(m model.Base) error {
	newLeaseRequest := m.(*LeaseRequest)

	if l.ID != newLeaseRequest.ID {
		return newLeaseRequestError("cannot change id")
	}

	if l.ClientID != newLeaseRequest.ClientID {
		return newLeaseRequestError("cannot change client id")
	}

	if l.Region != newLeaseRequest.Region {
		return newLeaseRequestError("cannot change region")
	}

	if l.Tag != newLeaseRequest.Tag {
		return newLeaseRequestError("cannot change tag")
	}

	if !l.FSM.Can(newLeaseRequest.Status) {
		return newLeaseRequestError("invalid status transition")
	}

	return nil
}

func (l *LeaseRequest) SetStatus(s string) {
	l.Status = statusMap[s]
}

func (l *LeaseRequest) F() *fsm.FSM {
	return l.FSM
}
