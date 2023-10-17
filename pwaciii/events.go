package pwaciii

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"time"
)

type EventType byte

type Event struct {
	EventType Command
	Data      []byte
}

const (
	InfoEntrada EventType = 0x17
	InfoSalida  EventType = 0x16
)

func (evt EventType) String() string {
	switch evt {
	case InfoEntrada:
		return "InfoEntrada"
	case InfoSalida:
		return "InfoSalida"
	default:
		return "UnknownEvent"
	}
}

func (d *Device) Events(ctx context.Context) chan Event {

	if ctx == nil {
		ctx = context.TODO()
	}
	d.chCmdResp = make(chan Event)

	ch := make(chan Event, 1)
	go func() {
		defer close(ch)
		defer close(d.chCmdResp)

		contxt, cancel := context.WithCancel(ctx)
		defer cancel()

		const maxErrors int = 3
		countErrors := 0

		for {
			data, err := func() ([]byte, error) {
				d.mux.Lock()
				defer d.mux.Unlock()

				t0 := time.Now()

				r := bufio.NewReader(d.portserial)
				buff := make([]byte, 32)
				n, err := r.Read(buff)
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

				fmt.Printf("data raw: %02X\n", buff[:n])
				if n < 1 {
					return nil, nil
				}

				temp := extractData(buff[:n])
				return temp, nil
			}()
			if err != nil {
				fmt.Println(err)
				return
			}
			if len(data) < 2 {
				// fmt.Println("without data")
				continue
			}
			// fmt.Printf("data: %s\n", data)

			evt := Command(data[0])

			select {
			case <-contxt.Done():
				return
			case d.chCmdResp <- Event{
				EventType: evt,
				Data:      data[:],
			}:
			default:
			}
			// fmt.Println("send data")
			select {
			case <-contxt.Done():
				return
			case ch <- Event{
				EventType: evt,
				Data:      data[1:],
			}:
			case <-time.After(10 * time.Millisecond):
				fmt.Println("timeout write listen events")
			}

		}
	}()

	return ch

}
