package fsm

import (
	"errors"
	"testing"
)

type mockState string

func (m mockState) Value() string {
	return string(m)
}

const fooState = mockState("foo")
const barState = mockState("bar")
const bazState = mockState("baz")
const batState = mockState("bat")

func TestNewFSM(t *testing.T) {
	type args struct {
		defaultState State
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"new fsm",
			args{
				defaultState: fooState,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewFSM(tt.args.defaultState)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewFSM() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestFSM_Value(t *testing.T) {
	type fields struct {
		CurrentState State
	}
	tests := []struct {
		name   string
		fields fields
		want   State
	}{
		{
			"foo",
			fields{
				CurrentState: fooState,
			},
			fooState,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &FSM{
				CurrentState: tt.fields.CurrentState,
			}
			if got := f.CurrentState; got != tt.want {
				t.Errorf("FSM.Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFSM_Can(t *testing.T) {
	type fields struct {
		CurrentState  State
		transitionMap TransitionMap
	}
	type args struct {
		newState State
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			"can transition - single candidate",
			fields{
				CurrentState: fooState,
				transitionMap: TransitionMap{
					fooState: []Transition{
						{
							State: barState,
						},
					},
				},
			},
			args{
				newState: barState,
			},
			true,
		},
		{
			"can transition - multiple candidates (first)",
			fields{
				CurrentState: fooState,
				transitionMap: TransitionMap{
					fooState: []Transition{
						{
							State: barState,
						},
						{
							State: bazState,
						},
					},
				},
			},
			args{
				newState: barState,
			},
			true,
		},
		{
			"can transition - multiple candidates (non-first)",
			fields{
				CurrentState: fooState,
				transitionMap: TransitionMap{
					fooState: []Transition{
						{
							State: barState,
						},
						{
							State: bazState,
						},
					},
				},
			},
			args{
				newState: bazState,
			},
			true,
		},
		{
			"cannot transition - single candidate",
			fields{
				CurrentState: fooState,
				transitionMap: TransitionMap{
					fooState: []Transition{
						{
							State: barState,
						},
					},
				},
			},
			args{
				newState: bazState,
			},
			false,
		},
		{
			"cannot transition - multiple candidates",
			fields{
				CurrentState: fooState,
				transitionMap: TransitionMap{
					fooState: []Transition{
						{
							State: barState,
						},
						{
							State: fooState,
						},
					},
				},
			},
			args{
				newState: bazState,
			},
			false,
		},
		{
			"cannot transition - invalid state",
			fields{
				CurrentState: barState,
				transitionMap: TransitionMap{
					fooState: []Transition{
						{
							State: barState,
						},
						{
							State: fooState,
						},
					},
				},
			},
			args{
				newState: bazState,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &FSM{
				CurrentState:  tt.fields.CurrentState,
				transitionMap: tt.fields.transitionMap,
			}
			if got := f.Can(tt.args.newState); got != tt.want {
				t.Errorf("FSM.Can() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFSM_Transition(t *testing.T) {
	type fields struct {
		CurrentState  State
		transitionMap TransitionMap
	}
	type args struct {
		newState State
		args     []EventArg
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    State
		wantErr bool
	}{
		{
			"simple transition",
			fields{
				CurrentState: fooState,
				transitionMap: TransitionMap{
					fooState: []Transition{
						{
							State: barState,
						},
					},
				},
			},
			args{
				newState: barState,
			},
			barState,
			false,
		},
		{
			"invalid transition",
			fields{
				CurrentState: fooState,
				transitionMap: TransitionMap{
					fooState: []Transition{
						{
							State: barState,
						},
					},
				},
			},
			args{
				newState: bazState,
			},
			fooState,
			true,
		},
		{
			"transition with callback",
			fields{
				CurrentState: fooState,
				transitionMap: TransitionMap{
					fooState: []Transition{
						{
							State:    barState,
							Callback: func(e Event) error { return nil },
						},
					},
				},
			},
			args{
				newState: barState,
			},
			barState,
			false,
		},
		{
			"transition with erroring callback",
			fields{
				CurrentState: fooState,
				transitionMap: TransitionMap{
					fooState: []Transition{
						{
							State:    barState,
							Callback: func(e Event) error { return errors.New("e") },
						},
					},
				},
			},
			args{
				newState: barState,
			},
			barState,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &FSM{
				CurrentState:  tt.fields.CurrentState,
				transitionMap: tt.fields.transitionMap,
			}
			if err := f.Transition(tt.args.newState, tt.args.args...); (err != nil) != tt.wantErr {
				t.Errorf("FSM.Transition() error = %v, wantErr %v", err, tt.wantErr)
			}
			if current := f.CurrentState; current != tt.want {
				t.Errorf("FSM.CurrentState = %v, want %v", current, tt.want)
			}
		})
	}
}

func TestFSM_AddTransitions(t *testing.T) {
	type fields struct {
		CurrentState  State
		transitionMap TransitionMap
	}
	type args struct {
		state       State
		transitions []Transition
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"add single initial transition",
			fields{
				transitionMap: make(TransitionMap),
			},
			args{
				state: fooState,
				transitions: []Transition{
					{
						State: barState,
					},
				},
			},
			false,
		},
		{
			"add multiple initial transition",
			fields{
				transitionMap: make(TransitionMap),
			},
			args{
				state: fooState,
				transitions: []Transition{
					{
						State: barState,
					},
					{
						State: bazState,
					},
				},
			},
			false,
		},
		{
			"add single transition to existing transitions",
			fields{
				transitionMap: TransitionMap{
					fooState: []Transition{
						{
							State: barState,
						},
					},
				},
			},
			args{
				state: fooState,
				transitions: []Transition{
					{
						State: bazState,
					},
					{
						State: batState,
					},
				},
			},
			false,
		},
		{
			"add multiple transitions to existing transitions",
			fields{
				transitionMap: TransitionMap{
					fooState: []Transition{
						{
							State: barState,
						},
					},
				},
			},
			args{
				state: fooState,
				transitions: []Transition{
					{
						State: bazState,
					},
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &FSM{
				CurrentState:  tt.fields.CurrentState,
				transitionMap: tt.fields.transitionMap,
			}
			if err := f.AddTransitions(tt.args.state, tt.args.transitions); (err != nil) != tt.wantErr {
				t.Errorf("FSM.AddTransitions() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFSM_AddCallback(t *testing.T) {
	type fields struct {
		transitionMap TransitionMap
	}
	type args struct {
		state      State
		transition State
		callback   Callback
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"add callback",
			fields{
				transitionMap: TransitionMap{
					fooState: []Transition{
						{
							State: barState,
						},
					},
				},
			},
			args{
				state:      fooState,
				transition: barState,
				callback:   func(e Event) error { return nil },
			},
			false,
		},
		{
			"invalid - source state does not exist",
			fields{
				transitionMap: TransitionMap{
					fooState: []Transition{
						{
							State: barState,
						},
					},
				},
			},
			args{
				state:      bazState,
				transition: barState,
				callback:   func(e Event) error { return nil },
			},
			true,
		},
		{
			"invalid - transition state does not exist",
			fields{
				transitionMap: TransitionMap{
					fooState: []Transition{
						{
							State: barState,
						},
					},
				},
			},
			args{
				state:      fooState,
				transition: bazState,
				callback:   func(e Event) error { return nil },
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &FSM{
				transitionMap: tt.fields.transitionMap,
			}
			if err := f.AddCallback(tt.args.state, tt.args.transition, tt.args.callback); (err != nil) != tt.wantErr {
				t.Errorf("FSM.AddCallback() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
