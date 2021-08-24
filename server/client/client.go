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

	ID         string `model:"immutable,required"`
	Token      string `proto:"ignore" model:"required"`
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

func (c *Client) Identifier() string {
	return c.ID
}

func (c *Client) Validate() error {
	return nil
}

func (c *Client) ValidateTransition(m model.Base) error {
	return nil
}

func (c *Client) SetStatus(s string) {
	c.Status = statusMap[s]
}

func (c *Client) F() *fsm.FSM {
	return c.FSM
}
