package model

import (
	"time"

	"github.com/p0pr0ck5/volchestrator/fsm"
)

type Model struct {
	CreatedAt time.Time
	UpdatedAt time.Time

	FSM *fsm.FSM
}

type Base interface {
	Init()
	Identifier() string
	Validate() error
	ValidateTransition(Base) error
	SetStatus(string)
	UpdateFSM() error
	Clone() Base
}
