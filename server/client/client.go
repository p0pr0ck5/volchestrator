package client

import "time"

type Client struct {
	ID         string
	Registered time.Time
	LastSeen   time.Time
}
