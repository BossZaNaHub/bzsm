package bzsm

import (
	"errors"
	"sync"
)

var ErrEventRejected = errors.New("event was rejected")

const (
	Default StateType = ""

	NoOperation EventType = "-"
)

// StateType set state type in state machine
type StateType string

// EventType set event type in state machine
type EventType string

type EventCtx interface{}

type Action interface {
	Do(evtCtx EventCtx) EventType
}

type Events map[EventType]StateType

type State struct {
	Action Action
	Events Events
}

// States map state
type States map[StateType]State

// BZSM BossZa State Machine prepare machine
type BZStateMachine struct {
	PreviousState StateType
	CurrentState  StateType
	States        States

	// Mutex lock state
	mu sync.Mutex
}

func (t StateType) String() string {
	return string(t)
}

func (t EventType) String() string {
	return string(t)
}

func (sm *BZStateMachine) getNextState(evt EventType) (StateType, error) {
	if state, ok := sm.States[sm.CurrentState]; ok {
		if state.Events != nil {
			if next, ok := state.Events[evt]; ok {
				return next, nil
			}
		}
	}
	return Default, nil
}

func (sm *BZStateMachine) DoEvent(evt EventType, ctx EventCtx) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	for {
		nextState, err := sm.getNextState(evt)
		if err != nil {
			return ErrEventRejected
		}

		state, ok := sm.States[nextState]
		if !ok && state.Action == nil {
			// Do state error
			//return errors.New("No State Found")
		}

		sm.PreviousState = sm.CurrentState
		sm.CurrentState = nextState

		//Execute next state action
		nextEvt := state.Action.Do(ctx)
		if nextEvt == NoOperation {
			return nil
		}
		evt = nextEvt
	}
}

func NewStateMachine(states States) *BZStateMachine {
	return &BZStateMachine{
		States: states,
	}
}
