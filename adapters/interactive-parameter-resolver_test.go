package adapters

import (
	"bytes"
	"io"
	"testing"
)

func Test_resolve(t *testing.T) {
	type args struct {
		rd   io.Reader
		name string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "some parameter", args: args{rd: bytes.NewBufferString("value"), name: "param_1"}, want: "value"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := resolve(tt.args.rd, tt.args.name); got != tt.want {
				t.Errorf("resolve() = %v, want %v", got, tt.want)
			}
		})
	}
}
