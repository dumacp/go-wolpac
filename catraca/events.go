package catraca

import (
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
	Input        EventType = "INPUT"
	Output       EventType = "OUTPUT"
	Sensor1UP    EventType = "Sensor1_UP"
	Sensor2UP    EventType = "Sensor2_UP"
	Sensor1DOWN  EventType = "Sensor1_DOWN"
	Sensor2DOWN  EventType = "Sensor2_DOWN"
	Alarm        EventType = "ALARM"
	ReleaseAlarm EventType = "RELEASE_ALARM"
	Cancel       EventType = "CancelAllow"
)

func (evt EventType) String() string {
	return string(evt)
}

func events(ctx context.Context, dev *Device, controlEntrance bool) (chan Event, error) {

	ch1, err := dev.pin1.EpollEvents(ctx,
		dev.InputsSysfsActiveLow,
		// gpiosysfs.Edge(dev.InputsSysfsEdge))
		gpiosysfs.Edge(gpiosysfs.Both))
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
		// listEvts := list.New()
		// funcListRemove := func() {
		// 	if listEvts.Len() > 10 {
		// 		// TODO: test
		// 		/**
		// 		el := listEvts.Back()
		// 		ss := make([]string, 0)
		// 		for el != nil {
		// 			if evt, ok := el.Value.(Event); ok {
		// 				ss = append(ss, fmt.Sprintf("%q", evt))
		// 			}
		// 			el = el.Next()
		// 		}
		// 		fmt.Printf("queue events: %s\n", strings.Join(ss, ", "))
		// 		/**/
		// 		listEvts.Remove(listEvts.Front())
		// 	}
		// }
		tickTimeoutAlarm := time.NewTimer(0)
		tickTimeoutAlarm.Stop()
		defer tickTimeoutAlarm.Stop()
		tickTimeoutEntrance := time.NewTimer(0)
		tickTimeoutEntrance.Stop()
		defer tickTimeoutEntrance.Stop()

		chChangeState := make(chan string, 2)
		fm := NewFSM(chChangeState)
		for {

			if err := func() error {
				dev.mux.Lock()
				defer dev.mux.Unlock()

				select {
				case <-ctx.Done():
					return fmt.Errorf("cancel events")
				case v := <-chChangeState:
					fmt.Printf("new fms catraca state: %q\n", fm.Current())
					switch v {
					case sStart:
					case sAllow:
						if tickTimeoutEntrance.Stop() {
							select {
							case <-tickTimeoutEntrance.C:
							default:
							}
						}
						tickTimeoutEntrance.Reset(dev.TimeoutEntrance)

					case sWait:
						// dev.activeAllow = false
						if out, err := funcCommand(); err != nil {
							return fmt.Errorf("error comand: %q, err: %s, output: %s", cmdDisable, err, out)
						}
						if tickTimeoutEntrance.Stop() {
							select {
							case <-tickTimeoutEntrance.C:
							default:
							}
						}
						if tickTimeoutAlarm.Stop() {
							select {
							case <-tickTimeoutAlarm.C:
							default:
							}
						}
						tickTimeoutAlarm.Reset(dev.TimeoutTurnAlarm)

					case sInput:
						dev.activeAllow = false
						evt := Event{EventType: Input, Time: time.Now()}
						select {
						case chEvt <- evt:
						default:
						}
					case sOuput:
						evt := Event{EventType: Output, Time: time.Now()}
						select {
						case chEvt <- evt:
						default:
						}
					case sAlarmExit:
						dev.activeAllow = false
						evt := Event{EventType: ReleaseAlarm, Time: time.Now()}
						select {
						case chEvt <- evt:
						default:
						}
					case sAlarmInput:
						evt := Event{EventType: Alarm, Time: time.Now()}
						select {
						case chEvt <- evt:
						default:
						}
					case sAlarmOutput:
					case sCancelAllow:
						evt := Event{EventType: Cancel, Time: time.Now()}
						select {
						case chEvt <- evt:
						default:
						}
						dev.activeAllow = false
						fm.Event(ctx, eExitCancel)
						if out, err := funcCommand(); err != nil {
							return fmt.Errorf("error comand: %q, err: %s, output: %s", cmdDisable, err, out)
						}
					}
				case <-dev.activeAllowchannel:
					fm.Event(ctx, eOneEntrance)
				case <-tickTimeoutEntrance.C:
					fm.Event(ctx, eTimeoutEntrance)
				case <-tickTimeoutAlarm.C:
					fm.Event(ctx, eTimeoutAlarm)
				case v := <-ch1:
					switch {
					case v.RisingEdge && dev.InputsSysfsActiveLow:
						fm.Event(ctx, eS1_0)
					case !v.RisingEdge && dev.InputsSysfsActiveLow:
						fm.Event(ctx, eS1_1)
					case v.RisingEdge && !dev.InputsSysfsActiveLow:
						fm.Event(ctx, eS1_1)
					case !v.RisingEdge && !dev.InputsSysfsActiveLow:
						fm.Event(ctx, eS1_0)
					}

				case v := <-ch2:

					switch {
					case v.RisingEdge && dev.InputsSysfsActiveLow:
						fm.Event(ctx, eS2_0)
					case !v.RisingEdge && dev.InputsSysfsActiveLow:
						fm.Event(ctx, eS2_1)
					case v.RisingEdge && !dev.InputsSysfsActiveLow:
						fm.Event(ctx, eS2_1)
					case !v.RisingEdge && !dev.InputsSysfsActiveLow:
						fm.Event(ctx, eS2_0)
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

func events_newcatraca(ctx context.Context, dev *Device, controlEntrance bool) (chan Event, error) {

	ch1, err := dev.pin1.EpollEvents(ctx,
		dev.InputsSysfsActiveLow,
		// gpiosysfs.Edge(dev.InputsSysfsEdge))
		gpiosysfs.Edge(gpiosysfs.Both))
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
		tickTimeoutAlarm := time.NewTimer(0)
		tickTimeoutAlarm.Stop()
		defer tickTimeoutAlarm.Stop()
		tickTimeoutEntrance := time.NewTimer(0)
		tickTimeoutEntrance.Stop()
		defer tickTimeoutEntrance.Stop()

		chChangeState := make(chan struct{}, 1)
		fm := NewFSM_v2(chChangeState)
		for {

			if err := func() error {
				dev.mux.Lock()
				defer dev.mux.Unlock()

				select {
				case <-ctx.Done():
					return fmt.Errorf("cancel events")
				case <-chChangeState:
					fmt.Printf("new fms catraca state: %q\n", fm.Current())
					switch fm.Current() {
					case sStart:
					case sAllow:
						if tickTimeoutEntrance.Stop() {
							select {
							case <-tickTimeoutEntrance.C:
							default:
							}
						}
						tickTimeoutEntrance.Reset(dev.TimeoutEntrance)

					case sWait:
						// dev.activeAllow = false
						if out, err := funcCommand(); err != nil {
							return fmt.Errorf("error comand: %q, err: %s, output: %s", cmdDisable, err, out)
						}
						if tickTimeoutEntrance.Stop() {
							select {
							case <-tickTimeoutEntrance.C:
							default:
							}
						}
						if tickTimeoutAlarm.Stop() {
							select {
							case <-tickTimeoutAlarm.C:
							default:
							}
						}
						tickTimeoutAlarm.Reset(dev.TimeoutTurnAlarm)

					case sInput:
						dev.activeAllow = false
						evt := Event{EventType: Input, Time: time.Now()}
						select {
						case chEvt <- evt:
						default:
						}
					case sOutput1, sOutput2, sOuput:
						dev.activeAllow = false
						evt := Event{EventType: Output, Time: time.Now()}
						select {
						case chEvt <- evt:
						default:
						}
					case sAlarmExit:
						dev.activeAllow = false
						evt := Event{EventType: ReleaseAlarm, Time: time.Now()}
						select {
						case chEvt <- evt:
						default:
						}
					case sAlarmInput:
						evt := Event{EventType: Alarm, Time: time.Now()}
						select {
						case chEvt <- evt:
						default:
						}
					case sAlarmOutput:
					case sCancelAllow:
						evt := Event{EventType: Cancel, Time: time.Now()}
						select {
						case chEvt <- evt:
						default:
						}
						dev.activeAllow = false
						fm.Event(ctx, eExitCancel)
						if out, err := funcCommand(); err != nil {
							return fmt.Errorf("error comand: %q, err: %s, output: %s", cmdDisable, err, out)
						}
					}
				case <-dev.activeAllowchannel:
					fm.Event(ctx, eOneEntrance)
				case <-tickTimeoutEntrance.C:
					fm.Event(ctx, eTimeoutEntrance)
				case <-tickTimeoutAlarm.C:
					fm.Event(ctx, eTimeoutAlarm)
				case v := <-ch1:
					switch {
					case v.RisingEdge && dev.InputsSysfsActiveLow:
						fm.Event(ctx, eS1_0)
					case !v.RisingEdge && dev.InputsSysfsActiveLow:
						fm.Event(ctx, eS1_1)
					case v.RisingEdge && !dev.InputsSysfsActiveLow:
						fm.Event(ctx, eS1_1)
					case !v.RisingEdge && !dev.InputsSysfsActiveLow:
						fm.Event(ctx, eS1_0)
					}

				case v := <-ch2:

					switch {
					case v.RisingEdge && dev.InputsSysfsActiveLow:
						fm.Event(ctx, eS2_0)
					case !v.RisingEdge && dev.InputsSysfsActiveLow:
						fm.Event(ctx, eS2_1)
					case v.RisingEdge && !dev.InputsSysfsActiveLow:
						fm.Event(ctx, eS2_1)
					case !v.RisingEdge && !dev.InputsSysfsActiveLow:
						fm.Event(ctx, eS2_0)
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
