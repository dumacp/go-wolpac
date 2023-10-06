package device

import "testing"

func TestOpts_OptsToString(t *testing.T) {

	tests := []struct {
		name string
		opt  Opts
		want string
	}{
		// TODO: Add test cases.
		{
			name: "test1",
			opt:  DefaultsOptions(),
			want: "03110410100015000000011SS",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := tt.opt
			if got := opt.OptsToString(); got != tt.want {
				t.Errorf("Opts.OptsToString() = %v, want %v", got, tt.want)
			}
		})
	}
}
