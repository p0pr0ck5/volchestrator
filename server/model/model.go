package model

import (
	"time"

	"github.com/p0pr0ck5/volchestrator/fsm"
)

type ModelOpt func(*Model) error

func WithSM(sm fsm.TransitionMap) ModelOpt {
	return func(m *Model) error {
		for k, v := range sm {
			if err := m.FSM.AddTransitions(k, v); err != nil {
				return err
			}
		}

		return nil
	}
}

type Model struct {
	CreatedAt time.Time
	UpdatedAt time.Time

	FSM *fsm.FSM
}

type Base interface {
	Init(...ModelOpt)
	Identifier() string
	Validate() error
	ValidateTransition(Base) error
	SetStatus(string)
	UpdateFSM() error
	Clone() Base
}
