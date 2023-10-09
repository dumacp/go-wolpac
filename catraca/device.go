package catraca

type Device struct{}

func New(signal, sensor int) Device {
	return Device{}
}

func (d *Device) Open() error {

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
