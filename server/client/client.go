package client

import (
	"time"

	"github.com/p0pr0ck5/volchestrator/server/model"
)

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

type Client struct {
	model.Model

	ID         string
	Token      string `proto:"ignore"`
	Registered time.Time
	LastSeen   time.Time
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
