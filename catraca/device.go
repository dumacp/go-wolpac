package catraca

import (
	"context"

	"github.com/dumacp/go-wolpac/gpiosysfs"
)

type Device struct {
	gpioTurnStart int
	gpioTurnEnd   int
}

func New(signal, sensorHalfTurnstart, sensorHalfTurnEnd int) Device {
	return Device{}
}

func (d *Device) Open() error {
	return nil
}

func (d *Device) OpenWithEvents(ctx context.Context) error {
	gpiosysfs.OpenPinWithEvents(ctx, d.gpioTurnStart)
	return nil
}

func (d *Device) Close() error {
	return nil
}

func (d *Device) ValueEntrance() (int, error) {
	return 0, nil
}

func (d *Device) ValueExit() (int, error) {
	return 0, nil
}
