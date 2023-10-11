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
	DLE = 0x10
	STX = 0x02
	ETX = 0x03
	ACK = 0x06
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

	data1 := bytes.ReplaceAll(data0, []byte{0x10, 0x02}, []byte{})

	data2 := bytes.ReplaceAll(data1, []byte{0x10, 0x03}, []byte{})

	return data2

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

	if waitResponse {
		for v := range ch {
			return v.data, v.err
		}
	}

	return nil, nil
}
