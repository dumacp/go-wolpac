package pwaciii

import (
	"io"
	"sync"
	"time"

	"github.com/tarm/serial"
)

const (
	READTIMEOUT = 100 * time.Millisecond
)

type Device struct {
	Opts
	confserial *serial.Config
	portserial io.ReadWriteCloser
	chCmdResp  chan Event
	mux        sync.Mutex
	muxRwrite  sync.Mutex
}

func New(opts ...OptsFunc) Device {
	o := DefaultsOptions()
	for _, fn := range opts {
		fn(&o)
	}

	c := serial.Config{
		Name:        o.Port,
		Baud:        19200,
		ReadTimeout: READTIMEOUT,
	}

	return Device{
		Opts:       o,
		confserial: &c,
	}
}

func (d *Device) Open(opts ...OptsFunc) error {

	for _, fn := range opts {
		fn(&d.Opts)
	}

	p, err := serial.OpenPort(d.confserial)
	if err != nil {
		return err
	}

	d.portserial = p

	return nil
}

func (d *Device) Command(cmd Command, data []byte) ([]byte, error) {
	resp, err := sendCommand(d, cmd, data)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
