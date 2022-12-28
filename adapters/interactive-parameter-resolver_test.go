package adapters

import (
	"bytes"
	"io"
	"testing"
)

func Test_resolve(t *testing.T) {
	type args struct {
		rd           io.Reader
		name         string
		defaultValue string
	}
	tests := []struct {
		name       string
		args       args
		wantResult string
		wantOutput string
	}{
		{name: "some parameter", args: args{rd: bytes.NewBufferString("value\n"), name: "param_1"}, wantResult: "value", wantOutput: "Enter value for parameter param_1: "},
		{name: "some parameter with default", args: args{rd: bytes.NewBufferString("value\n"), name: "param_1", defaultValue: "defaultValue"}, wantResult: "value", wantOutput: "Enter value for parameter param_1 or press enter for default defaultValue: "},
		{name: "pressing enter with default", args: args{rd: bytes.NewBufferString("\n"), name: "param_1", defaultValue: "defaultValue"}, wantResult: "defaultValue", wantOutput: "Enter value for parameter param_1 or press enter for default defaultValue: "},
		{name: "pressing enter without default", args: args{rd: bytes.NewBufferString("\n\nvalue\n"), name: "param_1"}, wantResult: "value", wantOutput: "Enter value for parameter param_1: Enter value for parameter param_1: Enter value for parameter param_1: "},
	}
	for _, tt := range tests {
		var b bytes.Buffer
		t.Run(tt.name, func(t *testing.T) {
			if got := resolve(tt.args.rd, &b, tt.args.name, tt.args.defaultValue); got != tt.wantResult {
				t.Errorf("resolve() = %v, want %v", got, tt.wantResult)
			}
			if b.String() != tt.wantOutput {
				t.Errorf("resolve() output = %s, want %s", b.String(), tt.wantOutput)
			}
		})
	}
}
