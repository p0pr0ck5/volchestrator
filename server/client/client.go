package client

import (
	"time"

	"github.com/p0pr0ck5/volchestrator/server/model"
)

type Client struct {
	model.Model

	ID         string
	Registered time.Time
	LastSeen   time.Time
}
