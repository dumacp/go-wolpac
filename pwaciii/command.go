package pwaciii

import "fmt"

type Command byte

const (
	LiberaUnaEntrada                         Command = 0x03
	LiberaUnaSalida                          Command = 0x25
	InformacionEstadoPasillo                 Command = 0x09
	SolicitaInformacionUsuarioPasandoSalida  Command = 0x27
	SolicitaInformacionUsuarioPasandoEntrada Command = 0x04
	AlteraContadoresInternosPWAC3            Command = 0x10
	AlteraEstadoBloqueo                      Command = 0x0d
	SolicitaEstadoBloqueo                    Command = 0x0e
	InformacionDeSalida                      Command = 0x16
	InformacionDeEntrada                     Command = 0x15
	AvisoMudancaEstadoLocalBloqueo           Command = 0x17
	ComandoLiberarAberturaCofre              Command = 0x30
	ComandoBloquearAberturaCofre             Command = 0x31
	SolicitaEstadoTampaCofre                 Command = 0x2e
	SolicitacaoEstadoChave                   Command = 0x13
	SolicitaEstadoTampas                     Command = 0x07
	ComandoColetarDatosReloj                 Command = 0x2B
	ComandoAjustarRelogio                    Command = 0x2C
	SolicitaAlteracionPictogramaAmarelo2     Command = 0x24
	SolicitaAlteracionPictogramaAmarelo1     Command = 0x23
	SolicitaAlteracionPictogramaBlanco1      Command = 0x21
	SolicitaAlteracionPictogramaBlanco2      Command = 0x22
	ColetaContadoresInternosPWAC3            Command = 0x0f
)

func (c Command) String() string {
	constantNames := map[Command]string{
		LiberaUnaEntrada:                         "LiberaUnaEntrada",
		LiberaUnaSalida:                          "LiberaUnaSalida",
		SolicitaInformacionUsuarioPasandoSalida:  "SolicitaInformacionUsuarioPasandoSalida",
		SolicitaInformacionUsuarioPasandoEntrada: "SolicitaInformacionUsuarioPasandoEntrada",
		AlteraContadoresInternosPWAC3:            "AlteraContadoresInternosPWAC3",
		AlteraEstadoBloqueo:                      "AlteraEstadoBloqueo",
		SolicitaEstadoBloqueo:                    "SolicitaEstadoBloqueo",
		InformacionDeSalida:                      "InformacionDeSalida",
		InformacionDeEntrada:                     "InformacionDeEntrada",
		AvisoMudancaEstadoLocalBloqueo:           "AvisoMudancaEstadoLocalBloqueo",
		ComandoLiberarAberturaCofre:              "ComandoLiberarAberturaCofre",
		ComandoBloquearAberturaCofre:             "ComandoBloquearAberturaCofre",
		SolicitaEstadoTampaCofre:                 "SolicitaEstadoTampaCofre",
		SolicitacaoEstadoChave:                   "SolicitacaoEstadoChave",
		SolicitaEstadoTampas:                     "SolicitaEstadoTampas",
		ComandoColetarDatosReloj:                 "ComandoColetarDadosReloj",
		ComandoAjustarRelogio:                    "ComandoAjustarRelogio",
		SolicitaAlteracionPictogramaAmarelo2:     "SolicitaAlteracionPictogramaAmarelo2",
		SolicitaAlteracionPictogramaAmarelo1:     "SolicitaAlteracionPictogramaAmarelo1",
		SolicitaAlteracionPictogramaBlanco1:      "SolicitaAlteracionPictogramaBlanco1",
		SolicitaAlteracionPictogramaBlanco2:      "SolicitaAlteracionPictogramaBlanco2",
		ColetaContadoresInternosPWAC3:            "ColetaContadoresInternosPWAC3",
	}

	name, found := constantNames[c]
	if !found {
		return fmt.Sprintf("Unknown Command value: 0x%02X", byte(c))
	}

	return name
}

func (c Command) WaitResponse() bool {

	return false
}

func (c Command) WaitAck() bool {

	return false
}

func (c Command) WithAck() bool {
	switch c {
	case LiberaUnaEntrada:
	case LiberaUnaSalida:
	case AlteraEstadoBloqueo:
	case SolicitaAlteracionPictogramaAmarelo1:
	case SolicitaAlteracionPictogramaAmarelo2:
	case SolicitaAlteracionPictogramaBlanco1:
	case SolicitaAlteracionPictogramaBlanco2:
	default:
		return false
	}
	return true
}

func (c Command) WithResponse() bool {
	switch c {
	case SolicitaEstadoBloqueo:
	case SolicitaInformacionUsuarioPasandoEntrada:
	case SolicitaInformacionUsuarioPasandoSalida:
	case ColetaContadoresInternosPWAC3:
	case InformacionEstadoPasillo:
	default:
		return false
	}
	return true
}
