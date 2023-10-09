package pwacii

import (
	"bufio"
	"context"
	"errors"
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
	AccessTimeout
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
	case AccessTimeout:
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
	case AccessTimeout:
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
	d.chCmdAck = make(chan int)

	ch := make(chan Event, 1)
	go func() {
		defer close(ch)
		defer close(d.chCmdResp)
		defer close(d.chCmdAck)

		contxt, cancel := context.WithCancel(ctx)
		defer cancel()

		const maxErrors int = 3
		countErrors := 0

		for {
			data, err := func() ([]byte, error) {
				d.mux.Lock()
				defer d.mux.Unlock()

				t0 := time.Now()

				s := bufio.NewReader(d.Port)
				datawithdelimiter, err := s.ReadBytes('\n')
				if err != nil {
					if !errors.Is(err, io.EOF) {
						return nil, fmt.Errorf("error read listen events: %s", err)

					} else if time.Since(t0) < READTIMEOUT/10 {
						countErrors++
						if countErrors > maxErrors {
							return nil, fmt.Errorf("%d errors io.EOF read listen events", countErrors)

						}
						return nil, nil
					}
				} else {
					countErrors = 0
				}
				return datawithdelimiter[:len(datawithdelimiter)-1], nil
			}()
			if err != nil {
				fmt.Println(err)
				return
			}
			if len(data) > 0 {
				if data[0] == '@' {
					select {
					case <-contxt.Done():
						return
					case d.chCmdAck <- 1:
					default:
					}
					continue
				}
				if data[0] == '%' {
					select {
					case <-contxt.Done():
						return
					case d.chCmdAck <- 0:
					default:
					}
					continue
				}
			} else {
				fmt.Println("without data")
				continue
			}
			// fmt.Printf("data: %s\n", data)

			if len(data) < 3 || data[0] != '!' {
				fmt.Printf("unkown data: [%X] (%q)\n", data[0], data[0])
				continue
			}
			var evt EventType
			dataevt := string(data[1:])
			switch dataevt {
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
			case AccessTimeout.Code():
				evt = AccessTimeout
			default:
				if len(data[1:]) > 2 {
					switch string(data[1:3]) {
					case ConfResponse.Code():
						evt = ConfResponse
						select {
						case <-contxt.Done():
							return
						case d.chCmdResp <- Event{
							EventType: evt,
							Data:      dataevt,
						}:
						default:
						}
					case StatusReponse.Code():
						evt = StatusReponse
						select {
						case <-contxt.Done():
							return
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
			// fmt.Println("send data")
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
	}()

	return ch

}
