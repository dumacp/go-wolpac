package catraca

import (
	"fmt"
	"time"

	"github.com/dumacp/go-wolpac/gpiosysfs"
)

type Opts struct {
	SignalLed            string
	InputsSysfsEdge      Edge
	SignalGpio           int
	InputSysfsT1         int
	InputSysfsT2         int
	TimeoutEntrance      time.Duration
	TimeoutTurnAlarm     time.Duration
	InputsSysfsActiveLow bool
}

type Edge gpiosysfs.Edge

const (
	// None means no edge is selected to generate interrupts.
	None Edge = Edge(gpiosysfs.None)
	// Rising edges is is selected to generate interrupts. Rising: level is getting to high from low.
	Rising Edge = Edge(gpiosysfs.Rising)
	// Falling edges is is selected to generate interrupts. Falling: level is getting to low from hight.
	Falling Edge = Edge(gpiosysfs.Falling)
	// Both rising and falling edges are selected to generate interrupts.
	Both Edge = Edge(gpiosysfs.Both)
)

func DefaultsOptions() Opts {
	return Opts{
		SignalLed:            "output1",
		SignalGpio:           0,
		InputSysfsT1:         85,
		InputSysfsT2:         86,
		TimeoutEntrance:      20 * time.Second,
		InputsSysfsEdge:      Falling,
		TimeoutTurnAlarm:     5 * time.Second,
		InputsSysfsActiveLow: true,
	}
}

func (opts Opts) String() string {
	return fmt.Sprintf("%+v", opts)
}

type OptsFunc func(*Opts)

func WithSignalTypePathLed(name string) OptsFunc {
	return func(opts *Opts) {
		opts.SignalLed = name
	}
}
func WithSignalTypeSysfs(gpio int) OptsFunc {
	return func(opts *Opts) {
		opts.SignalGpio = gpio
	}
}
func WithInputSysfsT1(gpio int) OptsFunc {
	return func(opts *Opts) {
		opts.InputSysfsT1 = gpio
	}
}

func WithInputSysfsT2(gpio int) OptsFunc {
	return func(opts *Opts) {
		opts.InputSysfsT2 = gpio
	}
}

func WithInputsSysfsInActiveLow(activelow bool) OptsFunc {
	return func(opts *Opts) {
		opts.InputsSysfsActiveLow = activelow
	}
}

func WithInputsSysfsEdge(edge Edge) OptsFunc {
	return func(opts *Opts) {
		opts.InputsSysfsEdge = edge
	}
}
func WithTimeoutEntrance(timeout time.Duration) OptsFunc {
	return func(opts *Opts) {
		opts.TimeoutEntrance = timeout
	}
}
