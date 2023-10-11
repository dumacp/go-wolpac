package pwaciii

import (
	"reflect"
	"testing"
)

func TestDefaultsOptions(t *testing.T) {
	tests := []struct {
		name string
		want Opts
	}{
		// TODO: Add test cases.
		{
			name: "test1",
			want: Opts{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DefaultsOptions(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DefaultsOptions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEstadoBloqueo_String(t *testing.T) {
	tests := []struct {
		name string
		s    EstadoBloqueo
		want string
	}{
		// TODO: Add test cases.
		{
			name: "test1",
			s:    EntradaLibreSalidaBloqueada,
			want: "EntradaLibreSalidaBloqueada",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.String(); got != tt.want {
				t.Errorf("EstadoBloqueo.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
