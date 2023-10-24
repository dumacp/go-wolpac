package catraca

import (
	"context"
	"fmt"
	"log"

	"github.com/looplab/fsm"
)

const (
	sStart        = "Sensores_0000"
	sAllow        = "sAllow"
	sWait         = "sWait"
	sCancelAllow  = "sCancelAllow"
	Sensores_1000 = "Sensores_1000"
	Sensores_0100 = "Sensores_0100"
	Sensores_0010 = "Sensores_0010"
	Sensores_0001 = "Sensores_0001"
	sOuput        = "sOuput"
	sInput        = "sInput"
	sAlarmInput   = "sAlarmInput"
	sAlarmOutput  = "sAlarmOutput"
	sAlarmExit    = "sAlarmExit"
	sInvalidInput = "sInvalidInput"
)

const (
	eS1_0            = "eS1_0"
	eS1_1            = "eS1_1"
	eS2_0            = "eS2_0"
	eS2_1            = "eS2_1"
	eOneEntrance     = "eOneEntrance"
	eTimeoutEntrance = "eTimeoutEntrance"
	eTimeoutAlarm    = "eTimeoutAlarm"
	// eExitAlarm       = "eExitAlarm"
)

func NewFSM(ch chan struct{}) *fsm.FSM {

	callbacksfsm := fsm.Callbacks{
		"before_event": func(contxt context.Context, e *fsm.Event) {
			fmt.Println("event fsm catraca: ", e.Event)
			if e.Err != nil {
				// log.Println(e.Err)
				e.Cancel(e.Err)
			}
		},
		"leave_state": func(contxt context.Context, e *fsm.Event) {
			if e.Err != nil {
				// log.Println(e.Err)
				e.Cancel(e.Err)
			}
		},
		"enter_state": func(contxt context.Context, e *fsm.Event) {
			select {
			case ch <- struct{}{}:
			default:
			}
			log.Printf("FSM catraca, state src: %s, state dst: %s", e.Src, e.Dst)
		},

		// "leave_closed": func(contxt context.Context, e *fsm.Event) {
		// },
		// "before_verify": func(contxt context.Context, e *fsm.Event) {
		// },
		// "enter_closed": func(contxt context.Context, e *fsm.Event) {
		// },
	}

	rfsm := fsm.NewFSM(
		sStart,
		fsm.Events{
			// {Name: eS1_0, Src: []string{sStart}, Dst: sInvalidInput},
			{Name: eOneEntrance, Src: []string{sStart}, Dst: sAllow},
			{Name: eS1_1, Src: []string{sAllow}, Dst: sWait},
			{Name: eS1_0, Src: []string{sWait}, Dst: sInput},
			{Name: eS2_1, Src: []string{sInput}, Dst: Sensores_0010},
			{Name: eS2_0, Src: []string{Sensores_0010}, Dst: sStart},

			{Name: eTimeoutEntrance, Src: []string{sAllow}, Dst: sCancelAllow},
			{Name: eTimeoutAlarm, Src: []string{sWait}, Dst: sAlarmInput},

			{Name: eS2_1, Src: []string{sStart}, Dst: Sensores_0001},
			{Name: eS2_0, Src: []string{Sensores_0001}, Dst: Sensores_0010},
			{Name: eS1_1, Src: []string{Sensores_0010}, Dst: sOuput},
			{Name: eS1_0, Src: []string{sOuput}, Dst: sStart},

			{Name: eS1_0, Src: []string{sAlarmInput}, Dst: sAlarmExit},
			// {Name: eS1_0, Src: []string{sAlarmExit}, Dst: sStart},
			{Name: eS2_0, Src: []string{sAlarmExit}, Dst: sStart},
		},
		callbacksfsm,
	)
	return rfsm
}
