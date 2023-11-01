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
	Error     error
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
			t0 := time.Now()
			data, err := func() ([]byte, error) {
				d.mux.Lock()
				defer d.mux.Unlock()

				s := bufio.NewReader(d.Port)
				b0, err := s.ReadByte()
				if err != nil {
					return nil, fmt.Errorf("error read listen events: %w", err)
				} else {
					countErrors = 0
				}
				if b0 == '@' {
					return []byte{'@'}, nil
				} else if b0 == '%' {
					return []byte{'%'}, nil
				} else if b0 == '\n' {
					return nil, nil
				} else if b0 == '\r' {
					return nil, nil
				} else if b0 != '!' {
					return nil, fmt.Errorf("wrong prefix in response: %q, %w", b0, ErrorRecv)
				}
				datawithdelimiter, err := s.ReadBytes('\n')
				if err != nil {
					return nil, fmt.Errorf("error read listen events: %w", err)
				} else {
					countErrors = 0
				}
				// fmt.Printf("turnstile data raw: %02X%02X\n", b0, datawithdelimiter)
				if len(datawithdelimiter) < 1 {
					return nil, nil
				}
				temp := make([]byte, 0)
				// temp = append(temp, b0)
				temp = append(temp, datawithdelimiter[:len(datawithdelimiter)-1]...)
				return temp, nil
			}()
			if err != nil {

				select {
				case <-contxt.Done():
					return
				case d.chCmdResp <- Event{
					EventType: 0,
					Data:      "",
					Error:     err,
				}:
				default:
				}
				if errors.Is(err, ErrorRecv) {
					fmt.Println(err)
				} else if !errors.Is(err, io.EOF) {
					fmt.Println(err)
					return
				} else if time.Since(t0) < READTIMEOUT/10 {
					fmt.Println(err)
					countErrors++
					if countErrors > maxErrors {
						return
					}
				}
				continue
			}
			if len(data) > 0 {
				// fmt.Printf("turnstile data raw: %q\n", data)
				if data[0] == '@' {
					select {
					case <-contxt.Done():
						return
					case d.chCmdAck <- 1:
					default:
					}
					continue
				} else if data[0] == '%' {
					select {
					case <-contxt.Done():
						return
					case d.chCmdAck <- 0:
					default:
					}
					continue
				}
			} else {
				// fmt.Println("without data")
				continue
			}
			// fmt.Printf("data: %s\n", data)

			if len(data) < 2 {
				select {
				case <-contxt.Done():
					return
				case d.chCmdResp <- Event{
					EventType: 0,
					Data:      "",
					Error:     fmt.Errorf("unkown data: [%X] (%q)", data, data),
				}:
				default:
				}
				continue
			}
			func() {
				var evt EventType
				dataevt := string(data[:])
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
					if len(data[:]) > 2 {
						switch string(data[:2]) {
						case ConfResponse.Code():
							evt = ConfResponse
							select {
							case <-contxt.Done():
								return
							case d.chCmdResp <- Event{
								EventType: evt,
								Data:      dataevt,
							}:
								return
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
								return
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
				case <-time.After(100 * time.Millisecond):
					fmt.Println("timeout write listen events")
				}
			}()

		}
	}()

	return ch

}
