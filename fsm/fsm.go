package fsm

import (
	"fmt"
	"sync"
)

type State interface {
	Value() string
}

type EventArg interface{}

type Event struct {
	f *FSM

	PreviousState State
	NextState     State

	Args []EventArg
}

type Callback func(e Event) error

type Transition struct {
	State    State
	Callback Callback
}

type TransitionMap map[State][]Transition

type FSM struct {
	CurrentState State
	stateMutex   sync.Mutex

	transitionMap      TransitionMap
	transitionMapMutex sync.Mutex

	transitionHistory []Event
}

func NewFSM(defaultState State) (*FSM, error) {
	f := &FSM{
		CurrentState:  defaultState,
		transitionMap: make(TransitionMap),
	}

	return f, nil
}

func (f *FSM) AddTransitions(state State, transitions []Transition) error {
	f.transitionMapMutex.Lock()
	defer f.transitionMapMutex.Unlock()

	f.transitionMap[state] = append(f.transitionMap[state], transitions...)

	return nil
}

func (f *FSM) Can(newState State) bool {
	f.stateMutex.Lock()
	defer f.stateMutex.Unlock()

	if f.CurrentState == newState {
		return true
	}

	for _, candidate := range f.transitionMap[f.CurrentState] {
		if candidate.State == newState {
			return true
		}
	}

	return false
}

func (f *FSM) Transition(newState State, args ...EventArg) error {
	if !f.Can(newState) {
		return fmt.Errorf("invalid state transition %q", newState)
	}

	f.stateMutex.Lock()
	defer f.stateMutex.Unlock()

	var transition Transition
	for _, candidate := range f.transitionMap[f.CurrentState] {
		if candidate.State == newState {
			transition = candidate
		}
	}

	e := Event{
		f:             f,
		PreviousState: f.CurrentState,
		NextState:     transition.State,
		Args:          args,
	}

	f.CurrentState = newState
	f.transitionHistory = append(f.transitionHistory, e)

	if cb := transition.Callback; cb != nil {
		if err := cb(e); err != nil {
			return err
		}
	}

	return nil
}

func (f *FSM) Value() string {
	return f.CurrentState.Value()
}
