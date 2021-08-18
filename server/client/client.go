package client

import (
	"time"

	"github.com/p0pr0ck5/volchestrator/fsm"
	"github.com/p0pr0ck5/volchestrator/server/model"
)

type Status int

func (s Status) Value() string {
	return ""
}

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

func (c *Client) ValidateTransition(newClient *Client) error {
	if c.ID != newClient.ID {
		return newClientError("cannot change id")
	}

	if c.Status != newClient.Status && !c.FSM.Can(newClient.Status) {
		return newClientError("invalid status transition")
	}

	return nil
}
