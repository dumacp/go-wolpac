package catraca

import (
	"container/list"
	"context"
	"fmt"
	"log"
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
	Input       EventType = "INPUT"
	Output      EventType = "OUTPUT"
	Sensor1UP   EventType = "Sensor1_UP"
	Sensor2UP   EventType = "Sensor2_UP"
	Sensor1DOWN EventType = "Sensor1_DOWN"
	Sensor2DOWN EventType = "Sensor2_DOWN"
	Alarm       EventType = "ALARM"
	Cancel      EventType = "CancelAllow"
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
	ch2, err := dev.pin2.EpollEvents(ctx,
		dev.InputsSysfsActiveLow,
		gpiosysfs.Edge(gpiosysfs.Both))
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

	chEvt := make(chan Event, 1)

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
				listEvts.Remove(listEvts.Front())
			}
		}
		tickTimeoutEntrance := time.NewTicker(dev.TimeoutEntrance)
		if !controlEntrance {
			tickTimeoutEntrance.Stop()
		} else {
			defer tickTimeoutEntrance.Stop()
		}
		tickTimeoutTurn := time.NewTimer(0)
		tickTimeoutTurn.Stop()
		defer tickTimeoutTurn.Stop()
		for {
			if err := func() error {
				dev.mux.Lock()
				defer dev.mux.Unlock()
				select {
				case <-ctx.Done():
					return fmt.Errorf("cancel events")
				case <-dev.activeAllowchannel:
					dev.activeAllow = true
					dev.activeStep = true
					if tickTimeoutTurn.Stop() {
						select {
						case <-tickTimeoutEntrance.C:
						default:
						}
					}
					tickTimeoutTurn.Reset(dev.TimeoutTurnAlarm)
				case <-tickTimeoutTurn.C:
					if dev.activeAllow || dev.activeStep {
						evt := Event{EventType: Alarm, Time: time.Now()}
						// listEvts.PushBack(evt)
						// // fmt.Printf("el list: %v\n", el)
						// funcListRemove()
						select {
						case chEvt <- evt:
						default:
						}
						if controlEntrance && dev.activeStep {
							dev.activeStep = false
							// fmt.Printf("activeStep: %v, activeAllow: %v\n", dev.activeStep, dev.activeAllow)
							if out, err := funcCommand(); err != nil {
								return fmt.Errorf("error comand: %q, err: %s, output: %s", cmdDisable, err, out)
							}
							el := listEvts.Back()
							// fmt.Printf("prev el list: %v\n", prev)
							var lastEvt *Event
							if el != nil && el.Value != nil {
								if evt, ok := el.Value.(Event); ok {
									lastEvt = &evt
								}
							}
							if lastEvt != nil {
								if dev.InputsSysfsEdge == Falling && lastEvt.EventType == Sensor2UP {
									dev.activeAllow = false
								} else if dev.InputsSysfsEdge == Rising && lastEvt.EventType == Sensor2DOWN {
									dev.activeAllow = false
								}
							} else {
								dev.activeAllow = false
							}

						}
					}
				case <-tickTimeoutEntrance.C:
					if controlEntrance && dev.activeStep {
						dev.activeStep = false
						// fmt.Printf("activeStep: %v, activeAllow: %v\n", dev.activeStep, dev.activeAllow)
						if out, err := funcCommand(); err != nil {
							return fmt.Errorf("error comand: %q, err: %s, output: %s", cmdDisable, err, out)
						}
						el := listEvts.Back()
						// fmt.Printf("prev el list: %v\n", prev)
						var lastEvt *Event
						if el != nil && el.Value != nil {
							if evt, ok := el.Value.(Event); ok {
								lastEvt = &evt
							}
						}
						if lastEvt != nil {
							if dev.InputsSysfsEdge == Falling && lastEvt.EventType == Sensor2UP {
								dev.activeAllow = false
							} else if dev.InputsSysfsEdge == Rising && lastEvt.EventType == Sensor2DOWN {
								dev.activeAllow = false
							}
						} else {
							dev.activeAllow = false
						}

						evt := Event{EventType: Cancel, Time: time.Now()}
						// listEvts.PushBack(evt)
						// // fmt.Printf("el list: %v\n", el)
						// funcListRemove()
						select {
						case chEvt <- evt:
						default:
						}
					}
				case v := <-ch1:
					s1 := func() EventType {
						if v.RisingEdge {
							return Sensor1UP
						}
						return Sensor1DOWN

					}()
					evt := Event{EventType: s1, Time: time.Now()}
					el := listEvts.PushBack(evt)
					// fmt.Printf("el list: %v\n", el)
					funcListRemove()
					select {
					case chEvt <- evt:
					default:
					}
					prev := el.Prev()
					// fmt.Printf("prev el list: %v\n", prev)
					var prevEvt *Event
					if prev != nil && prev.Value != nil {
						if evt, ok := prev.Value.(Event); ok {
							prevEvt = &evt
						}
					}
					if controlEntrance && dev.activeStep {
						dev.activeStep = false
						// fmt.Printf("activeStep: %v, activeAllow: %v\n", dev.activeStep, dev.activeAllow)
						if out, err := funcCommand(); err != nil {
							return fmt.Errorf("error comand: %q, err: %s, output: %s", cmdDisable, err, out)
						}
						// fmt.Printf("activeStep: %v, activeAllow: %v\n", dev.activeStep, dev.activeAllow)
					}
					if prevEvt != nil && (prevEvt.EventType == Sensor2UP || prevEvt.EventType == Sensor2DOWN) {
						evt := Event{EventType: Output, Time: time.Now()}
						// listEvts.PushBack(evt)
						// // fmt.Printf("el list: %v\n", prev)
						// funcListRemove()
						select {
						case chEvt <- evt:
						default:
						}
					}
				case v := <-ch2:
					s2 := func() EventType {
						if v.RisingEdge {
							return Sensor2UP
						}
						return Sensor2DOWN

					}()
					evt := Event{EventType: s2, Time: time.Now()}
					el := listEvts.PushBack(evt)
					// fmt.Printf("el list: %v\n", el)
					funcListRemove()
					select {
					case chEvt <- evt:
					default:
					}
					prev := el.Prev()
					// fmt.Printf("prev el list: %v\n", prev)
					var prevEvt *Event
					if prev != nil && prev.Value != nil {
						if evt, ok := prev.Value.(Event); ok {
							prevEvt = &evt
						}
					}

					exitPulse := false
					if v.RisingEdge {
						if string(dev.InputsSysfsEdge) == string(Falling) {
							exitPulse = true
						}
					} else {
						if dev.InputsSysfsEdge == Rising {
							exitPulse = true
						}
					}
					// fmt.Printf("exit pulse sensor 2? (%v), rising: %v, %s\n", exitPulse, v.RisingEdge, dev.InputsSysfsEdge)
					// fmt.Printf("activeStep: %v, activeAllow: %v\n", dev.activeStep, dev.activeAllow)
					// fmt.Printf("prev event: %v\n", prevEvt)
					if exitPulse && prevEvt != nil && (prevEvt.EventType == Sensor2UP || prevEvt.EventType == Sensor2DOWN) {
						dev.activeAllow = false
						// evt := Event{EventType: Output, Time: time.Now()}
						// fmt.Printf("activeStep: %v, activeAllow: %v\n", dev.activeStep, dev.activeAllow)
						evt := Event{EventType: Input, Time: time.Now()}
						// listEvts.PushBack(evt)
						// // fmt.Printf("el list: %v\n", el)
						// funcListRemove()
						// // println("send input")
						select {
						case chEvt <- evt:
						default:
						}
					}
				}
				return nil
			}(); err != nil {
				log.Println(err)
				return
			}
		}
	}()

	return chEvt, nil
}
