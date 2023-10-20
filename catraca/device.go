package catraca

import (
	"context"
	"fmt"
	"os/exec"
	"sync"
	"time"

	"github.com/dumacp/go-wolpac/gpiosysfs"
)

type Device struct {
	Opts
	pin1               *gpiosysfs.Pin
	pin2               *gpiosysfs.Pin
	ctx                context.Context
	cancel             func()
	activeAllowchannel chan struct{}
	activeStep         bool
	activeAllow        bool
	mux                sync.Mutex
}

func New(opts ...OptsFunc) Device {
	o := DefaultsOptions()
	for _, fn := range opts {
		fn(&o)
	}
	return Device{
		Opts:               o,
		activeAllowchannel: make(chan struct{}),
	}
}

func (d *Device) Open() error {
	if pin1, err := gpiosysfs.OpenPin(d.InputSysfsT1); err != nil {
		return err
	} else {
		d.pin1 = pin1
	}
	if pin2, err := gpiosysfs.OpenPin(d.InputSysfsT2); err != nil {
		return err
	} else {
		d.pin2 = pin2
	}
	return nil
}

func (d *Device) Events(ctx context.Context, autoCancel bool) (chan Event, error) {
	contxt, cancel := context.WithCancel(ctx)
	d.ctx = contxt
	d.cancel = cancel
	return events(contxt, d, autoCancel)
}

func (d *Device) Close() error {
	if d.cancel != nil {
		d.cancel()
	}
	return nil
}

func (d *Device) OneEntrance() error {
	if d.activeStep || d.activeAllow {
		return fmt.Errorf("turnstile already active")
	}
	select {
	case d.activeAllowchannel <- struct{}{}:
	case <-time.After(1 * time.Second):
		return fmt.Errorf("turnstile active timeout")
	}
	cmdEnable := fmt.Sprintf("echo 1 > /sys/class/leds/%s/brightness", d.SignalLed)
	funcCommand := func() ([]byte, error) {
		if out, err := exec.Command("/bin/sh", "-c", cmdEnable).Output(); err != nil {
			return out, err
		} else if len(out) > 0 {
			return out, nil
		}
		return nil, nil
	}
	if out, err := funcCommand(); err != nil {
		return fmt.Errorf("error comand: %q, err: %s, output: %s", cmdEnable, err, out)

	}
	return nil
}

func (d *Device) CancelEntrance() error {
	cmdDisable := fmt.Sprintf("echo 0 > /sys/class/leds/%s/brightness", d.SignalLed)
	funcCommand := func() ([]byte, error) {
		if out, err := exec.Command("/bin/sh", "-c", cmdDisable).Output(); err != nil {
			return out, err
		} else if len(out) > 0 {
			return out, nil
		}
		return nil, nil
	}
	if out, err := funcCommand(); err != nil {
		return fmt.Errorf("error comand: %q, err: %s, output: %s", cmdDisable, err, out)

	}
	return nil
}
