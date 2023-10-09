package pwacii

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/tarm/serial"
)

const (
	READTIMEOUT = 100 * time.Millisecond
)

type Conf struct {
	Opts
	serialconf *serial.Config
}

type Device struct {
	Port      io.ReadWriteCloser
	chCmdAck  chan int
	chCmdResp chan Event
	mux       sync.Mutex
	muxRwrite sync.Mutex
}

func NewConf(port string, opts ...OptsFunc) Conf {

	c := &serial.Config{
		Name:        port,
		Baud:        19200,
		ReadTimeout: READTIMEOUT,
	}

	o := DefaultsOptions()

	conf := Conf{
		Opts:       o,
		serialconf: c,
	}

	return conf
}

func (c *Device) Close() error {

	return c.Port.Close()
}

func (c *Conf) Open(opts ...OptsFunc) (*Device, error) {

	p, err := serial.OpenPort(c.serialconf)
	if err != nil {
		return nil, err
	}

	o := &c.Opts
	for _, fn := range opts {
		fn(o)
	}

	dev := &Device{}
	dev.Port = p

	// resp, err := command(dev, ConfStatus, "")
	// if err != nil {
	// 	return nil, err
	// }
	// fmt.Printf("conf status response: %q\n", resp)

	resp, err := command(dev, ConfRequest, o.OptsToString())
	if err != nil {
		return nil, err
	}
	fmt.Printf("conf response: %q\n", resp)

	return dev, nil
}

func command(d *Device, cmd CommandType, data string) (string, error) {
	if d.Port == nil {
		return "", io.EOF
	}
	d.muxRwrite.Lock()
	defer d.muxRwrite.Unlock()

	// w := bufio.NewWriter(d.Port)
	// r := bufio.NewReader(d.Port)
	// w := d.Port
	// r := d.Port
	s := fmt.Sprintf("$%s%s\r\n", cmd.Code(), data)
	fmt.Printf("cmd to send: %q, %X\n", s, s)

	// if n, err := w.Write([]byte(s)); err != nil {
	// 	return "", err
	// } else {
	// 	fmt.Printf("%d bytes writtern\n", n)
	// }

	// buff := make([]byte, 1024)
	// if n, err := r.Read(buff); err != nil && n == 0 {
	// 	return "", err
	// } else {
	// 	fmt.Printf("%d bytes read: %X\n", n, buff[:n])
	// }

	// return string(buff), nil

	/**/

	type response struct {
		data string
		err  error
	}

	ch := make(chan response)

	go func() {

		defer close(ch)
		// r := bufio.NewReader(d.Port)
		var r *bufio.Reader
		funcAck := func() (string, error) {
			if r == nil {
				r = bufio.NewReader(d.Port)
			}
			resp, err := r.ReadString('\n')
			if err != nil && len(resp) == 0 {
				return "", err
			}
			if len(resp) <= 0 {
				return resp, fmt.Errorf("nil response: %q", resp)
			}
			if resp[0] != '@' {
				return resp, fmt.Errorf("unkown response: %q", resp[0])
			}
			return resp, nil
		}

		funcResp := func() (string, error) {
			if r == nil {
				r = bufio.NewReader(d.Port)
			}
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
			return resp[1:], nil
		}

		chresponse, err := func() (string, error) {
			if d.chCmdAck == nil {
				if resp, err := funcAck(); err != nil {
					return resp, err
				}
				if !cmd.WithResponse() {
					return "", nil
				}
			}
			if d.chCmdResp == nil && cmd.WithResponse() {
				if resp, err := funcResp(); err != nil {
					return resp, err
				} else {
					return resp, nil
				}
			}
			select {
			case v, ok := <-d.chCmdAck:
				if !ok {
					if resp, err := funcAck(); err != nil {
						return resp, err
					}
				} else if v != 1 {
					return "", fmt.Errorf("error %q response", "%")
				}
				if !cmd.WithResponse() {
					return "", nil
				}
			case v, ok := <-d.chCmdResp:
				if !ok {
					if resp, err := funcResp(); err != nil {
						return resp, err
					} else {
						return resp, nil
					}
				}
				if !strings.EqualFold(v.EventType.Code(), cmd.Code()) {
					return data, fmt.Errorf("cmd response different to cmd: %q != %q, data: %q",
						cmd, v.EventType, fmt.Sprintf("%s%s", v.EventType.Code(), v.Data))
				}
				return v.Data, nil
			case <-time.After(300 * time.Millisecond):
				return "", fmt.Errorf("timeout read")
			}
			return "", nil
		}()
		select {
		case ch <- response{
			data: chresponse,
			err:  err,
		}:
		case <-time.After(10 * time.Millisecond):
			fmt.Printf("response command (%q) without receiver", cmd.String())
		}
	}()

	if n, err := d.Port.Write([]byte(s)); err != nil {
		return "", err
	} else if n <= 0 {
		return "", fmt.Errorf("write 0 bytes to serial port")
	}
	for v := range ch {
		return v.data, v.err
	}
	return "", fmt.Errorf("unkown error")
	/**/
}
