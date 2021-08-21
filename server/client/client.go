package client

import (
	"time"

	"github.com/p0pr0ck5/volchestrator/fsm"
	"github.com/p0pr0ck5/volchestrator/server/model"
)

type Status int

type ClientError struct {
	e string
}

func newClientError(err string) ClientError {
	return ClientError{
		e: err,
	}
}

func (e ClientError) Error() string {
	return e.e
}

var statusMap = map[string]Status{
	"Alive":    Alive,
	"Deleting": Deleting,
}

const (
	Alive Status = iota
	Deleting
)

var sm = fsm.TransitionMap{
	Alive: []fsm.Transition{
		{
			State: Deleting,
		},
	},
}

type Client struct {
	model.Model

	ID         string
	Token      string `proto:"ignore"`
	Status     Status
	Registered time.Time
	LastSeen   time.Time
}

func (c *Client) Init() {
	c.FSM, _ = fsm.NewFSM(c.Status)
	for k, v := range sm {
		c.FSM.AddTransitions(k, v)
	}
}

func (c *Client) Validate() error {
	if c.ID == "" {
		return newClientError("missing id")
	}

	if c.Token == "" {
		return newClientError("missing token")
	}

	return nil
}

func (c *Client) ValidateTransition(m model.Base) error {
	newClient := m.(*Client)

	if c.ID != newClient.ID {
		return newClientError("cannot change id")
	}

	if !c.FSM.Can(newClient.Status) {
		return newClientError("invalid status transition")
	}

	return nil
}

func (c *Client) StatusVal() int {
	return int(c.Status)
}

func (c *Client) SetStatus(s string) {
	c.Status = statusMap[s]
}

func (c *Client) F() *fsm.FSM {
	return c.FSM
}
