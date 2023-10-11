package catraca

import (
	"container/list"
	"context"
	"fmt"
	"os/exec"
	"time"

	"github.com/dumacp/go-wolpac/gpiosysfs"
)

type EventType string

type Event struct {
	EventType EventType
	Time      time.Time
}

const (
	Input   EventType = "INPUT"
	Output  EventType = "OUTPUT"
	Sensor1 EventType = "Sensor1"
	Sensor2 EventType = "Sensor2"
	Alarm   EventType = "ALARM"
	Cancel  EventType = "CancelAllow"
)

func (evt EventType) String() string {
	return string(evt)
}

func events(ctx context.Context, dev *Device, controlEntrance bool) (chan Event, error) {

	ch1, err := dev.pin1.EpollEvents(ctx,
		dev.InputsSysfsActiveLow,
		gpiosysfs.Edge(dev.InputsSysfsEdge))
	if err != nil {
		return nil, err
	}
	ch2, err := dev.pin1.EpollEvents(ctx,
		dev.InputsSysfsActiveLow,
		gpiosysfs.Edge(dev.InputsSysfsEdge))
	if err != nil {
		return nil, err
	}

	// 0 -> 1 -> 2 -> 0 : Input
	// 0 -> 2 -> 1 -> 0 : Output
	// 0 -> 1 -> 1 -> 0 :
	cmdDisable := fmt.Sprintf("echo 0 > /sys/class/leds/%s/brightness", dev.SignalLed)
	funcCommand := func() ([]byte, error) {
		if out, err := exec.Command("/bin/sh", "-c", cmdDisable).Output(); err != nil {
			return out, err
		} else if len(out) > 0 {
			return out, nil
		}
		return nil, nil
	}

	chEvt := make(chan Event)

	go func() {
		defer close(chEvt)
		listEvts := list.New()
		funcListRemove := func() {
			if listEvts.Len() > 10 {
				// TODO: test
				/**
				el := listEvts.Back()
				ss := make([]string, 0)
				for el != nil {
					if evt, ok := el.Value.(Event); ok {
						ss = append(ss, fmt.Sprintf("%q", evt))
					}
					el = el.Next()
				}
				fmt.Printf("queue events: %s\n", strings.Join(ss, ", "))
				/**/
				listEvts.Remove(listEvts.Back())
			}
		}
		tickTimeoutEntrance := time.NewTicker(dev.TimeoutEntrance)
		if !controlEntrance {
			tickTimeoutEntrance.Stop()
		} else {
			defer tickTimeoutEntrance.Stop()
		}
		tickTimeoutTurn := time.NewTicker(dev.TimeoutTurnAlarm)
		defer tickTimeoutTurn.Stop()
		for {
			select {
			case <-ctx.Done():
				fmt.Println("cancel events")
				return
			case <-tickTimeoutTurn.C:
				el := listEvts.Front()
				var evt *Event
				if el != nil && el.Value != nil {
					if evtx, ok := el.Value.(Event); ok {
						evt = &evtx
					}
				}
				if evt != nil && time.Since(evt.Time) > dev.TimeoutTurnAlarm {
					evt := Event{EventType: Alarm, Time: time.Now()}
					listEvts.PushFront(evt)
					funcListRemove()
					select {
					case chEvt <- evt:
					default:
					}
				}
			case <-tickTimeoutEntrance.C:
				if controlEntrance && dev.activeStep {
					if out, err := funcCommand(); err != nil {
						fmt.Printf("error comand: %q, err: %s, output: %s\n", cmdDisable, err, out)
						return
					}
					evt := Event{EventType: Cancel, Time: time.Now()}
					listEvts.PushFront(evt)
					funcListRemove()
					select {
					case chEvt <- evt:
					default:
					}
				}
			case <-ch1:
				evt := Event{EventType: Sensor1, Time: time.Now()}
				el := listEvts.PushFront(evt)
				funcListRemove()
				select {
				case chEvt <- evt:
				default:
				}
				prev := el.Prev()
				var prevEvt *Event
				if prev != nil && prev.Value != nil {
					if evt, ok := prev.Value.(Event); ok {
						prevEvt = &evt
					}
				}
				if controlEntrance && dev.activeStep {
					if out, err := funcCommand(); err != nil {
						fmt.Printf("error comand: %q, err: %s, output: %s\n", cmdDisable, err, out)
						return
					}
				}
				if prevEvt != nil && prevEvt.EventType == Sensor2 && time.Since(prevEvt.Time) < dev.TimeoutTurnAlarm {
					evt := Event{EventType: Output, Time: time.Now()}
					listEvts.PushFront(evt)
					funcListRemove()
					select {
					case chEvt <- evt:
					default:
					}
				}
			case <-ch2:
				evt := Event{EventType: Sensor2, Time: time.Now()}
				el := listEvts.PushFront(evt)
				funcListRemove()
				select {
				case chEvt <- evt:
				default:
				}
				prev := el.Prev()
				var prevEvt *Event
				if prev != nil && prev.Value != nil {
					if evt, ok := prev.Value.(Event); ok {
						prevEvt = &evt
					}
				}
				if prevEvt != nil && prevEvt.EventType == Sensor1 && time.Since(prevEvt.Time) < dev.TimeoutTurnAlarm {
					evt := Event{EventType: Input, Time: time.Now()}
					listEvts.PushFront(evt)
					funcListRemove()
					select {
					case chEvt <- evt:
					default:
					}
				}
			}
		}
	}()

	return chEvt, nil
}
