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
