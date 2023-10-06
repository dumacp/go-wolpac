package device

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/tarm/serial"
)

type Device struct {
	Opts
	conf      *serial.Config
	port      *serial.Port
	chCmdResp chan Event
	mux       sync.Mutex
	muxRwrite sync.Mutex
}

func New(port string, opts ...OptsFunc) *Device {

	c := &serial.Config{
		Name:        port,
		Baud:        19200,
		ReadTimeout: 30 * time.Millisecond,
	}

	o := DefaultsOptions()
	for _, fn := range opts {
		fn(&o)
	}
	dev := &Device{
		Opts: o,
		conf: c,
	}

	return dev
}

func (d *Device) Open() error {

	p, err := serial.OpenPort(d.conf)
	if err != nil {
		return err
	}

	d.port = p

	return nil
}

func command(d *Device, cmd CommandType, data string) (string, error) {
	if d.port == nil {
		return "", io.EOF
	}
	d.muxRwrite.Lock()
	defer d.muxRwrite.Unlock()

	w := bufio.NewWriter(d.port)
	s := fmt.Sprintf("$%s%s\n", cmd.Code(), data)
	fmt.Printf("cmd to send: %q\n", s)

	if n, err := w.WriteString(s); err != nil {
		return "", err
	} else if n <= 0 {
		return "", fmt.Errorf("write 0 bytes to serial port")
	}

	if d.chCmdResp == nil {
		r := bufio.NewReader(d.port)
		resp, err := r.ReadString('\n')
		if err != nil {
			return "", err
		}
		if len(resp) < 3 {
			return resp, fmt.Errorf("unkown response: %q", resp)
		}
		if len(resp) < 1 && resp[0] != '!' {
			return resp, fmt.Errorf("unkown response: %q", resp)
		}
		return resp, nil
	}

	select {
	case v, ok := <-d.chCmdResp:
		if !ok {
			r := bufio.NewReader(d.port)
			resp, err := r.ReadString('\n')
			if err != nil {
				return "", err
			}
			if len(resp) < 3 {
				return resp, fmt.Errorf("unkown response: %q", resp)
			}
			if len(resp) < 1 && resp[0] != '!' {
				return resp, fmt.Errorf("unkown response: %q", resp)
			}
			return resp, nil
		}
		if !strings.EqualFold(v.EventType.Code(), cmd.Code()) {
			return data, fmt.Errorf("cmd response different to cmd: %q != %q, data: %q",
				cmd, v.EventType, fmt.Sprintf("%s%s", v.EventType.Code(), v.Data))
		}
		return v.Data, nil
	case <-time.After(600 * time.Millisecond):
	}

	return "", nil
}
