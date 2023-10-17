package pwaciii

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"time"
)

const (
	DLE byte = 0x10
	STX byte = 0x02
	ETX byte = 0x03
	ACK byte = 0x06
	NAK byte = 0x15
)

func formatapdu(data []byte) []byte {
	check := byte(0)
	message := []byte{DLE, STX}

	for _, b := range data {
		check ^= b
		// fmt.Printf("check: %X\n", check)
		if b == DLE {
			message = append(message, DLE)
		}
		message = append(message, b)
	}

	if check == DLE {
		message = append(message, DLE)
	}
	message = append(message, check, DLE, ETX)

	return message
}
func extractData(apdu []byte) []byte {

	data0 := bytes.ReplaceAll(apdu, []byte{0x10, 0x10}, []byte{0x10})

	data1 := bytes.Replace(data0, []byte{0x10, 0x02}, []byte{}, 1)

	if len(data1) > 2 {
		if bytes.Equal(data1[len(data1)-2:], []byte{0x10, 0x03}) {
			return data1[:len(data1)-2]
		}
	}
	return data1

}

func sendCommand(d *Device, cmd Command, data []byte) ([]byte, error) {

	port := d.portserial

	// w := bufio.NewWriter(d.Port)
	// r := bufio.NewReader(d.Port)
	// w := d.Port
	// r := d.Port

	fmt.Printf("cmd to send: %02X, [%X]\n", cmd, data)

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
		data []byte
		err  error
	}

	ch := make(chan response)

	go func() {

		defer close(ch)
		// r := bufio.NewReader(d.Port)
		var r *bufio.Reader

		const maxErrors int = 3
		countErrors := 0
		funcResp := func() ([]byte, error) {
			if r == nil {
				r = bufio.NewReader(port)
			}
			// resp, err := r.ReadString('\n')
			t0 := time.Now()
			buff := make([]byte, 32)
			n, err := r.Read(buff)
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
			} else {
				countErrors = 0
			}

			fmt.Printf("data raw: %02X\n", buff[:n])
			if n < 1 {
				return nil, nil
			}

			temp := extractData(buff[:n])
			return temp, nil

		}

		chresponse, err := func() ([]byte, error) {
			if (cmd.WaitResponse() || cmd.WaitAck()) && (d.chCmdResp == nil) {
				if resp, err := funcResp(); err != nil {
					return resp, err
				} else {
					return resp, nil
				}
			}
			select {

			case v, ok := <-d.chCmdResp:
				if !ok && (cmd.WithResponse() || cmd.WithAck()) {
					if resp, err := funcResp(); err != nil {
						return resp, err
					} else {
						return resp, nil
					}
				}
				if len(v.Data) < 2 {
					break
				}
				if v.Data[0] != byte(cmd) {
					return nil, fmt.Errorf("cmd response different to cmd: %q != %q, data: %02X",
						byte(cmd), byte(v.EventType), v.Data)
				}
				return []byte(v.Data[1:]), nil
			case <-time.After(300 * time.Millisecond):
				return nil, fmt.Errorf("timeout read")
			}
			return nil, nil
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

	apdu := make([]byte, 0)
	apdu = append(apdu, byte(cmd))
	apdu = append(apdu, data...)
	if n, err := port.Write(apdu); err != nil {
		return nil, err
	} else if n <= 0 {
		return nil, fmt.Errorf("write 0 bytes to serial port")
	}
	for v := range ch {
		return v.data, v.err
	}
	return nil, fmt.Errorf("unkown error")
	/**/
}
