package pwacii

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/tarm/serial"
)

const (
	READTIMEOUT = 1200 * time.Millisecond
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

	for _, fn := range opts {
		fn(&o)
	}

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

		const maxErrors int = 3
		countErrors := 0
		funcResp := func(awaitResponse bool) ([]byte, error) {
			d.mux.Lock()
			defer d.mux.Unlock()
			if r == nil {
				r = bufio.NewReader(d.Port)
			}
			// resp, err := r.ReadString('\n')

			var resp []byte
			for range make([]int, 3) {
				t0 := time.Now()
				var err error
				resp, err = r.ReadBytes('\n')
				// fmt.Printf("%q, %s\n", resp, err)
				if err != nil {
					if !errors.Is(err, io.EOF) {
						return nil, fmt.Errorf("error read listen events: %s", err)
					} else if time.Since(t0) < READTIMEOUT/10 {
						countErrors++
						if countErrors > maxErrors {
							return nil, fmt.Errorf("%d errors io.EOF read listen events", countErrors)
						}
						return nil, fmt.Errorf("readTimeout (%s) nil response", err)
					}
					if len(resp) <= 0 {
						continue
					}
					break
				}
				countErrors = 0
				if len(resp) <= 1 {
					continue
				}
				break
			}

			if len(resp) <= 0 {
				return resp, fmt.Errorf("nil response: %q", resp)
			}
			if resp[0] == '@' && !awaitResponse {
				return nil, nil
			} else if resp[0] == '%' {
				return nil, fmt.Errorf("no ack response")
			} else if resp[0] == '!' && len(resp) > 1 {
				return resp[1:], nil
			}
			return nil, fmt.Errorf("unkown response: %q", resp)
		}

		chresponse, err := func() ([]byte, error) {
			if d.chCmdAck == nil || d.chCmdResp == nil {
				if resp, err := funcResp(cmd.WithResponse()); err != nil {
					return resp, err
				} else {
					return resp, nil
				}
			}
			select {
			case v, ok := <-d.chCmdAck:
				if !ok {
					if resp, err := funcResp(cmd.WithResponse()); err != nil {
						return resp, err
					} else {
						return resp, nil
					}
				} else if v != 1 {
					return nil, fmt.Errorf("error %q response", "%")
				}
				if !cmd.WithResponse() {
					return nil, nil
				}
			case v, ok := <-d.chCmdResp:
				if !ok {
					if resp, err := funcResp(cmd.WithResponse()); err != nil {
						return resp, err
					} else {
						return resp, nil
					}
				}
				if v.Error != nil {
					return nil, v.Error
				}
				if !strings.EqualFold(v.EventType.Code(), cmd.Code()) {
					return nil, fmt.Errorf("cmd response different to cmd: %q != %q, data: %q",
						cmd, v.EventType, fmt.Sprintf("%s%s", v.EventType.Code(), v.Data))
				}
				return []byte(v.Data), nil
			case <-time.After(1200 * time.Millisecond):
				return nil, fmt.Errorf("timeout read")
			}
			return nil, nil
		}()
		datatosend := ""
		if len(chresponse) > 0 {
			datatosend = string(chresponse)
		}
		select {
		case ch <- response{
			data: datatosend,
			err:  err,
		}:
		case <-time.After(1200 * time.Millisecond):
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
