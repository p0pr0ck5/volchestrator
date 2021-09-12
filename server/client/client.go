package client

import (
	"time"

	"github.com/p0pr0ck5/volchestrator/fsm"
	"github.com/p0pr0ck5/volchestrator/server/model"
)

type Status int

var statusMap = map[string]Status{
	"Alive":    Alive,
	"Deleting": Deleting,
}

const (
	Alive Status = iota
	Deleting
)

type Client struct {
	model.Model

	ID         string `model:"immutable,required,reference=LeaseRequest:ClientID"`
	Token      string `proto:"ignore" model:"required"`
	LeaseID    string `model:"reference=Lease:ID"`
	Status     Status
	Registered time.Time
	LastSeen   time.Time
}

func (c *Client) Init(opts ...model.ModelOpt) {
	c.FSM, _ = fsm.NewFSM(c.Status)

	for _, opt := range opts {
		opt(&c.Model)
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

func (c *Client) UpdateFSM() error {
	return c.FSM.Transition(c.Status)
}

func (v *Client) Clone() model.Base {
	vv := *v
	return &vv
}
