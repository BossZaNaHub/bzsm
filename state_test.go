package bzsm

import (
	"context"
	"log"
	"testing"
)

const (
	Avaliable   StateType = "Available"
	Unavaliable StateType = "Unavailable"

	PeriodOpen  EventType = "open"
	PeriodClose EventType = "close"
)

type OpenMarketAction struct{}
type CloseMarketAction struct{}

func (op *OpenMarketAction) Do(ctx EventCtx) EventType {
	log.Println("Do OpenMarketAction: ", ctx)
	return PeriodOpen
}

func (op *CloseMarketAction) Do(ctx EventCtx) EventType {
	log.Println("Do CloseMarketAction: ", ctx)
	return PeriodClose
}

func TestPeriodState(t *testing.T) {
	setNewState := States{
		Default: State{
			Events: Events{
				PeriodOpen: Unavaliable,
			},
		},
		Avaliable: State{
			Events: Events{
				PeriodOpen: Avaliable,
			},
			Action: &OpenMarketAction{},
		},
		Unavaliable: State{
			Events: Events{
				PeriodClose: Unavaliable,
			},
			Action: &CloseMarketAction{},
		},
	}
	ctx := context.TODO()
	stm := NewStateMachine(setNewState)
	err := stm.DoEvent(PeriodOpen, nil)
	t.Logf("init state: %v", err)
	if err != nil {
		t.Errorf("init state err: %v", err)
	}

	err = stm.DoEvent(PeriodClose, ctx)
	t.Logf("Expected event: %v", err)
	if err != nil {
		t.Errorf("Expected event should be rejected: %v", err)
	}

	t.Log("Doing well")
}
