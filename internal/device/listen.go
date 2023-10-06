package device

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"time"
)

type EventType int

type Event struct {
	EventType EventType
	Data      string
}

const (
	Input EventType = iota
	Output
	HalfTurnStart
	HalfTurnEnd
	SensorAlarmActived
	SensorAlarmDisabled
	InputCancel
	ConfResponse
	StatusReponse
)

func (evt EventType) String() string {
	switch evt {
	case Input:
		return "Input"
	case Output:
		return "Output"
	case HalfTurnStart:
		return "HalfTurnStart"
	case HalfTurnEnd:
		return "HalfTurnEnd"
	case SensorAlarmActived:
		return "SensorAlarmActived"
	case SensorAlarmDisabled:
		return "SensorAlarmDisabled"
	case InputCancel:
		return "InputCancel"
	case ConfResponse:
		return "ConfResponse"
	case StatusReponse:
		return "StatusReponse"
	default:
		return "UnknownEvent"
	}
}

func (evt EventType) Code() string {
	switch evt {
	case Input:
		return "I1"
	case Output:
		return "I2"
	case HalfTurnStart:
		return "MG1"
	case HalfTurnEnd:
		return "MG0"
	case SensorAlarmActived:
		return "VI1"
	case SensorAlarmDisabled:
		return "VI0"
	case InputCancel:
		return "TMP"
	case ConfResponse:
		return "CG"
	case StatusReponse:
		return "SS"
	default:
		return "UnknownEvent"
	}
}

func (d *Device) Events(ctx context.Context) chan Event {

	if ctx == nil {
		ctx = context.TODO()
	}
	d.chCmdResp = make(chan Event)

	contxt, cancel := context.WithCancel(ctx)
	defer cancel()

	ch := make(chan Event, 1)
	go func() {
		defer close(ch)
		defer close(d.chCmdResp)

		for {
			d.mux.Lock()
			defer d.mux.Unlock()

			s := bufio.NewScanner(d.port)
			for s.Scan() {
				data := s.Bytes()
				if len(data) < 3 || data[0] != '!' {
					fmt.Printf("unkown data: [%X] (%q)\n", data[0], data[0])
					continue
				}
				var evt EventType
				var dataevt string
				switch string(data[1:]) {
				case Input.Code():
					evt = Input
				case Output.Code():
					evt = Output
				case HalfTurnStart.Code():
					evt = HalfTurnStart
				case HalfTurnEnd.Code():
					evt = HalfTurnEnd
				case SensorAlarmActived.Code():
					evt = SensorAlarmActived
				case SensorAlarmDisabled.Code():
					evt = SensorAlarmDisabled
				case InputCancel.Code():
					evt = InputCancel
				default:
					if len(data[1:]) > 2 {
						switch string(data[1:3]) {
						case ConfResponse.Code():
							evt = ConfResponse
							dataevt = string(data[3:])
							select {
							case d.chCmdResp <- Event{
								EventType: evt,
								Data:      dataevt,
							}:
							default:
							}
						case StatusReponse.Code():
							evt = StatusReponse
							dataevt = string(data[3:])
							select {
							case d.chCmdResp <- Event{
								EventType: evt,
								Data:      dataevt,
							}:
							default:
							}
						default:
						}
					}
				}
				select {
				case <-contxt.Done():
					return
				case ch <- Event{
					EventType: evt,
					Data:      dataevt,
				}:
				case <-time.After(10 * time.Millisecond):
					fmt.Println("timeout write listen events")
				}
			}
			err := s.Err()
			fmt.Printf("error scan events: %s", err)
			if s.Err() == io.EOF || s.Err() == nil {
				fmt.Println("io.EOF write listen events")
				return
			}

		}
	}()

	return ch

}
