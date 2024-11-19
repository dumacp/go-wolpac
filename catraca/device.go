package catraca

import (
	"context"
	"fmt"
	"os/exec"
	"sync"

	"github.com/dumacp/go-wolpac/gpiosysfs"
)

type Device struct {
	Opts
	pin1               *gpiosysfs.Pin
	pin2               *gpiosysfs.Pin
	pathSignalValue    string
	ctx                context.Context
	cancel             func()
	activeAllowchannel chan struct{}
	// activeStep         bool
	activeAllow bool
	mux         sync.Mutex
}

func New(opts ...OptsFunc) Device {
	o := DefaultsOptions()
	for _, fn := range opts {
		fn(&o)
	}
	return Device{
		Opts:               o,
		activeAllowchannel: make(chan struct{}, 1),
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
	if d.SignalGpio > 0 {
		sign, err := gpiosysfs.OpenPin(d.SignalGpio)
		if err != nil {
			return err
		}
		if err := sign.SetDirection(gpiosysfs.Out); err != nil {
			return err
		}
		// if out, err := exec.Command("/bin/sh", "-c", fmt.Sprintf("echo %d > /sys/class/gpio/export", d.SignalGpio)).CombinedOutput(); err != nil {
		// 	return fmt.Errorf("%s, error: %s", out, err)
		// }
		// if out, err := exec.Command("/bin/sh", "-c", fmt.Sprintf("echo out > /sys/class/gpio/gpio%d/direction", d.SignalGpio)).CombinedOutput(); err != nil {
		// 	return fmt.Errorf("%s, error: %s", out, err)
		// }
		d.pathSignalValue = fmt.Sprintf("/sys/class/gpio/gpio%d/value", d.SignalGpio)
	} else if len(d.SignalLed) > 0 {
		d.pathSignalValue = fmt.Sprintf("/sys/class/leds/%s/brightness", d.SignalLed)
	}
	fmt.Printf("////////////// options turnstile ///////////////////////////: %s\n", d.Opts)

	return nil
}

func (d *Device) Events(ctx context.Context, autoCancel bool) (chan Event, error) {
	contxt, cancel := context.WithCancel(ctx)
	d.ctx = contxt
	d.cancel = cancel
	if d.NewEvents {
		return events_newcatraca(ctx, d, autoCancel)
	}
	return events(contxt, d, autoCancel)
}

func (d *Device) Events_newcatraca(ctx context.Context, autoCancel bool) (chan Event, error) {
	contxt, cancel := context.WithCancel(ctx)
	d.ctx = contxt
	d.cancel = cancel
	return events_newcatraca(contxt, d, autoCancel)
}

func (d *Device) Close() error {
	if d.cancel != nil {
		d.cancel()
	}
	return nil
}

func (d *Device) OneEntrance() error {
	if d.activeAllow {
		return fmt.Errorf("turnstile already active")
	}
	select {
	case d.activeAllowchannel <- struct{}{}:
	default:
	}
	d.activeAllow = true
	cmdEnable := fmt.Sprintf("echo 1 > %s", d.pathSignalValue)
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
	d.activeAllow = false
	cmdDisable := fmt.Sprintf("echo 0 > %s", d.pathSignalValue)
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
