package pwaciii

import (
	"bufio"
	"bytes"
	"context"
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
		fmt.Printf("check: %X\n", check)
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

func sendCommand(port io.ReadWriteCloser, waitResponse, waitAck bool, cmd byte, data []byte) ([]byte, error) {

	ctx, cancel := context.WithTimeout(context.TODO(), READTIMEOT+10*time.Millisecond)
	defer cancel()

	// type response struct {
	// 	data []byte
	// 	err  error
	// }
	var ch chan struct {
		data []byte
		err  error
	}
	if waitResponse {
		r := bufio.NewReader(port)
		ch = make(chan struct {
			data []byte
			err  error
		})

		go func() {
			defer close(ch)

			buff := make([]byte, 32)
			n, err := r.Read(buff)
			if err != nil {
				select {
				case ch <- struct {
					data []byte
					err  error
				}{
					data: nil,
					err:  err,
				}:
				case <-ctx.Done():
				}
				return
			}

			select {
			case ch <- struct {
				data []byte
				err  error
			}{
				data: buff[:n],
				err:  nil,
			}:
			case <-ctx.Done():
			}
		}()
	}

	datacmd := make([]byte, 0)
	datacmd = append(datacmd, cmd)
	datacmd = append(datacmd, data...)

	apdu := formatapdu(datacmd)
	if _, err := port.Write(apdu); err != nil {
		return nil, err
	}

	if waitResponse || waitAck {
		for v := range ch {
			if v.err != nil {
				return nil, v.err
			}
			if len(v.data) <= 0 {
				return nil, v.err
			}
			data := extractData(v.data)

			if len(data) <= 0 {
				if waitAck {
					return nil, fmt.Errorf("without ACK response")
				}
				if waitResponse {
					return nil, fmt.Errorf("without response")
				}
				return nil, nil
			}

			if waitAck {
				if data[len(data)-1] == ACK {
					return nil, nil
				} else if data[len(data)-1] == NAK {
					return nil, fmt.Errorf("NAK response")
				} else {
					return nil, fmt.Errorf("unkown response: [%X]", data)
				}
			}
			if waitResponse {
				if data[0] != byte(cmd) {
					return nil, fmt.Errorf("unkown response: [%X], cmd (%02x) != (%02X)", data, byte(cmd), data[0])
				}
				return data[1:], nil
			}
		}
	}

	return nil, nil
}
