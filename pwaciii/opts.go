package pwaciii

import (
	"encoding/json"
	"fmt"
)

type Opts struct {
	Port          string
	EstadoBloqueo int
}

type EstadoBloqueo int

const (
	EntradaLibreSalidaLibre EstadoBloqueo = iota
	EntradaLibreSalidaControlada
	EntradaLibreSalidaBloqueada
	EntradacontroladaSalidaLibre
	EntradacontroladaSalidaControlada
	EntradacontroladaSalidaBloqueada
	EntradaBloqeadadaSalidaLibre
	EntradaBloqeadadaSalidaControlada
	EntradaBloqeadadaSalidaBloqueada
)

func (e EstadoBloqueo) String() string {
	constantNames := []string{
		"EntradaLibreSalidaLibre",
		"EntradaLibreSalidaControlada",
		"EntradaLibreSalidaBloqueada",
		"EntradacontroladaSalidaLibre",
		"EntradacontroladaSalidaControlada",
		"EntradacontroladaSalidaBloqueada",
		"EntradaBloqeadadaSalidaLibre",
		"EntradaBloqeadadaSalidaControlada",
		"EntradaBloqeadadaSalidaBloqueada",
	}

	if e < 0 || int(e) >= len(constantNames) {
		return fmt.Sprintf("Unknown EstadoBloqueo value: %d", e)
	}

	return constantNames[e]
}

func DefaultsOptions() Opts {
	return Opts{
		Port:          "/dev/ttyUSB2",
		EstadoBloqueo: int(EntradaLibreSalidaBloqueada),
	}
}

func (opts Opts) String() string {

	data, err := json.Marshal(opts)
	if err != nil {
		return ""
	}

	return string(data)
}

type OptsFunc func(*Opts)

func WithPort(port string) func(*Opts) {

	return func(o *Opts) {
		o.Port = port
	}
}

func WithEstadoBloque(estado EstadoBloqueo) func(*Opts) {

	return func(o *Opts) {
		o.EstadoBloqueo = int(estado)
	}
}
