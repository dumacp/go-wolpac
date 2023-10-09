package main

import (
	"fmt"
	"io"

	"github.com/dumacp/go-wolpac/pwacii"
)

type Conf struct {
	pwacii.Opts
	port io.ReadWriteCloser
}

func NewConf(port io.ReadWriteCloser, opts ...pwacii.OptsFunc) Conf {

	o := pwacii.DefaultsOptions()

	conf := Conf{
		Opts: o,
		port: port,
	}

	return conf
}

func (c *Conf) Open(opts ...pwacii.OptsFunc) (*pwacii.Device, error) {

	o := &c.Opts
	for _, fn := range opts {
		fn(o)
	}

	dev := &pwacii.Device{}
	dev.Port = c.port

	resp, err := dev.Command(pwacii.ConfRequest, o.OptsToString())
	if err != nil {
		return nil, err
	}
	fmt.Printf("conf response: %q\n", resp)

	return dev, nil
}
