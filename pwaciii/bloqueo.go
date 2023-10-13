package pwaciii

import "fmt"

func (d *Device) SolicitaEstadoBloqueo() (EstadoBloqueo, error) {

	cmd := SolicitaEstadoBloqueo
	resp, err := sendCommand(d.portserial, cmd.WaitResponse(), cmd.WaitAck(), byte(cmd), nil)
	if err != nil {
		return -1, err
	}
	if len(resp) < 2 {
		return -1, fmt.Errorf("cmd (%02X) != %02X", byte(cmd), resp[0])
	}

	result := EstadoBloqueo(resp[1])

	switch EstadoBloqueo(resp[1]) {
	case EntradaLibreSalidaLibre:
	case EntradaLibreSalidaControlada:
	case EntradaLibreSalidaBloqueada:
	case EntradacontroladaSalidaLibre:
	case EntradacontroladaSalidaControlada:
	case EntradacontroladaSalidaBloqueada:
	case EntradaBloqeadadaSalidaLibre:
	case EntradaBloqeadadaSalidaControlada:
	case EntradaBloqeadadaSalidaBloqueada:
	default:
		return -1, fmt.Errorf("unkown response %d (%X)", result, result)
	}

	return result, nil
}

func (d *Device) AlteraEstadoBloqueo(estado EstadoBloqueo) error {

	cmd := AlteraEstadoBloqueo
	if _, err := sendCommand(d.portserial, cmd.WaitResponse(), cmd.WaitAck(), byte(cmd), []byte{byte(estado)}); err != nil {
		return err
	}
	return nil
}
