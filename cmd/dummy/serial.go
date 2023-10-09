package main

import (
	"bytes"
	"io"
	"time"
)

type serial struct {
	buff bytes.Buffer
}

func NewSerial(data []byte) io.ReadWriteCloser {
	buff := bytes.NewBuffer(data)
	return &serial{buff: *buff}
}

func (s *serial) Read(p []byte) (n int, err error) {

	data, err := s.buff.ReadBytes('\n')
	if err != nil {
		return 0, err
	}
	time.Sleep(30 * time.Millisecond)
	if len(data) > len(p) {
		p = append(p, make([]byte, len(data)-len(p))...)
	}
	// fmt.Printf("data lowlevel: %q\n", data)
	copy(p, data)
	return len(data), nil
}

func (s *serial) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func (s *serial) Close() error {
	return nil
}
