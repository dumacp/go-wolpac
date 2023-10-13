package pwaciii

import (
	"encoding/binary"
	"fmt"
)

func (d *Device) OneEntrance() error {

	cmd := LiberaUnaEntrada
	if _, err := sendCommand(d.portserial, cmd.WaitResponse(), cmd.WaitAck(), byte(cmd), []byte("permitir")); err != nil {
		return err
	}
	return nil

}

type EstadoPaso string

const (
	EsperandoPasodeUsuario EstadoPaso = "W"
	PasodeUsuarioDesistido EstadoPaso = "T"
	OcurrioUnPasoDeEntrada EstadoPaso = "Y"
	OcurrioUnPasoDeSalida  EstadoPaso = "X"
)

var estadoNames = map[EstadoPaso]string{
	EsperandoPasodeUsuario: "EsperandoPasodeUsuario",
	PasodeUsuarioDesistido: "PasodeUsuarioDesistido",
	OcurrioUnPasoDeEntrada: "OcurrioUnPasoDeEntrada",
	OcurrioUnPasoDeSalida:  "OcurrioUnPasoDeSalida",
}

func (estado EstadoPaso) String() string {
	if name, ok := estadoNames[estado]; ok {
		return name
	}
	return ""
}

func (d *Device) SolicitaEstadoEntrada() (EstadoPaso, error) {

	cmd := SolicitaInformacionUsuarioPasandoEntrada
	resp, err := sendCommand(d.portserial, cmd.WaitResponse(), cmd.WaitAck(), byte(cmd), nil)
	if err != nil {
		return "", err
	}
	if len(resp) < 2 {
		return "", fmt.Errorf("unknown response [%X]", resp)
	} else if resp[0] != byte(cmd) {
		return "", fmt.Errorf("cmd (%02X) != %02X, unknown response [%X]", byte(cmd), resp[0], resp)
	}

	result := EstadoPaso(resp[1])

	switch EstadoPaso(resp[1]) {
	case EsperandoPasodeUsuario:
	case PasodeUsuarioDesistido:
	case OcurrioUnPasoDeEntrada:
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
		return "", fmt.Errorf("unknown response [%X]", resp)
	} else if resp[0] != byte(cmd) {
		return "", fmt.Errorf("cmd (%02X) != %02X, unknown response [%X]", byte(cmd), resp[0], resp)
	}

	result := EstadoPaso(resp[1])

	switch EstadoPaso(resp[1]) {
	case EsperandoPasodeUsuario:
	case PasodeUsuarioDesistido:
	case OcurrioUnPasoDeSalida:
	default:
		return "", fmt.Errorf("unkown response %q", resp)
	}

	return result, nil
}

func (d *Device) SolicitaContadores() (int64, int64, error) {

	cmd := ColetaContadoresInternosPWAC3
	resp, err := sendCommand(d.portserial, cmd.WaitResponse(), cmd.WaitAck(), byte(cmd), nil)
	if err != nil {
		return 0, 0, err
	}
	if len(resp) < 9 {
		return 0, 0, fmt.Errorf("unknown response [%X]", resp)
	} else if resp[0] != byte(cmd) {
		return 0, 0, fmt.Errorf("cmd (%02X) != %02X, unknown response [%X]", byte(cmd), resp[0], resp)
	}

	entradas := binary.LittleEndian.Uint32(resp[1:5])
	salidas := binary.LittleEndian.Uint32(resp[5:9])

	return int64(entradas), int64(salidas), nil
}
