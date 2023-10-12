package pwaciii

import "fmt"

func (d *Device) OneEntrance() error {

	cmd := LiberaUnaEntrada
	if _, err := sendCommand(d.portserial, cmd.WaitResponse(), cmd.WaitAck(), byte(cmd), nil); err != nil {
		return err
	}
	return nil

}

type EstadoPaso string

const (
	EsperandoPasodeUsuario EstadoPaso = "W"
	PasodeUsuarioDesistido EstadoPaso = "T"
	OcurrioUnPaso          EstadoPaso = "Y"
)

func (d *Device) SolicitaEstadoEntrada() (EstadoPaso, error) {

	cmd := SolicitaInformacionUsuarioPasandoEntrada
	resp, err := sendCommand(d.portserial, cmd.WaitResponse(), cmd.WaitAck(), byte(cmd), nil)
	if err != nil {
		return "", err
	}
	if len(resp) < 2 {
		return "", fmt.Errorf("cmd (%02X) != %02X", byte(cmd), resp[0])
	}

	result := EstadoPaso(resp[1])

	switch EstadoPaso(resp[1]) {
	case EsperandoPasodeUsuario:
	case PasodeUsuarioDesistido:
	case OcurrioUnPaso:
	default:
		return "", fmt.Errorf("unkown response %q", resp)
	}

	return result, nil
}

func (d *Device) SolicitaEstadoSalida() (EstadoPaso, error) {

	cmd := SolicitaInformacionUsuarioPasandoSalida
	resp, err := sendCommand(d.portserial, cmd.WaitResponse(), cmd.WaitAck(), byte(cmd), nil)
	if err != nil {
		return "", err
	}
	if len(resp) < 2 {
		return "", fmt.Errorf("cmd (%02X) != %02X", byte(cmd), resp[0])
	}

	result := EstadoPaso(resp[1])

	switch EstadoPaso(resp[1]) {
	case EsperandoPasodeUsuario:
	case PasodeUsuarioDesistido:
	case OcurrioUnPaso:
	default:
		return "", fmt.Errorf("unkown response %q", resp)
	}

	return result, nil
}

func (d *Device) Solicitacontadores() (EstadoPaso, error) {

	cmd := SolicitaInformacionUsuarioPasandoSalida
	resp, err := sendCommand(d.portserial, cmd.WaitResponse(), cmd.WaitAck(), byte(cmd), nil)
	if err != nil {
		return "", err
	}
	if len(resp) < 2 {
		return "", fmt.Errorf("cmd (%02X) != %02X", byte(cmd), resp[0])
	}

	result := EstadoPaso(resp[1])

	switch EstadoPaso(resp[1]) {
	case EsperandoPasodeUsuario:
	case PasodeUsuarioDesistido:
	case OcurrioUnPaso:
	default:
		return "", fmt.Errorf("unkown response %q", resp)
	}

	return result, nil
}
