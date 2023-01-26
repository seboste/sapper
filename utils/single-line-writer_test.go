package utils

import (
	"bytes"
	"reflect"
	"testing"
)

func TestSingleLineWriter_Write(t *testing.T) {
	type args struct {
		p [][]byte
	}
	tests := []struct {
		name      string
		args      args
		wantBytes []byte
		wantN     int
		wantErr   bool
	}{
		{name: "single write", args: args{[][]byte{[]byte("hello")}}, wantBytes: []byte("hello"), wantN: 5, wantErr: false},
		{name: "write with new line at the end", args: args{[][]byte{[]byte("hello\n")}}, wantBytes: []byte("hello"), wantN: 5, wantErr: false},
		{name: "write with new line at the beginning", args: args{[][]byte{[]byte("\nhello")}}, wantBytes: []byte("hello"), wantN: 5, wantErr: false},
		{name: "empty write", args: args{[][]byte{[]byte("")}}, wantBytes: nil, wantN: 0, wantErr: false},
		{name: "single write with new line", args: args{[][]byte{[]byte("hello\nworld")}}, wantBytes: []byte("world"), wantN: 5, wantErr: false},
		{name: "two writes", args: args{[][]byte{[]byte("hello"), []byte("world")}}, wantBytes: []byte("helloworld"), wantN: 10, wantErr: false},
		{name: "two writes with new line end of first", args: args{[][]byte{[]byte("hello\n"), []byte("world")}}, wantBytes: []byte("hello\b\b\b\b\bworld"), wantN: 15, wantErr: false},
		{name: "two writes with new line beginning of second", args: args{[][]byte{[]byte("hello"), []byte("\nworld")}}, wantBytes: []byte("hello\b\b\b\b\bworld"), wantN: 15, wantErr: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var b bytes.Buffer
			slw := MakeSingleLineWriter(&b)
			gotN := 0
			var err error
			for _, p := range tt.args.p {
				n, e := slw.Write(p)
				if e != nil {
					err = e
				}
				gotN = gotN + n

			}
			if (err != nil) != tt.wantErr {
				t.Errorf("Write() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotN != tt.wantN {
				t.Errorf("Write() = %v, want %v", gotN, tt.wantN)
			}
			gotBytes := b.Bytes()
			if !reflect.DeepEqual(gotBytes, tt.wantBytes) {
				t.Errorf("Write() bytes = %v (%s), want %v (%s)", gotBytes, string(gotBytes), tt.wantBytes, string(tt.wantBytes))
			}
		})
	}
}

func TestSingleLineWriter_Cleanup(t *testing.T) {
	type args struct {
		p []byte
	}
	tests := []struct {
		name      string
		args      args
		wantBytes []byte
	}{
		{name: "single char", args: args{p: []byte("a")}, wantBytes: []byte("a\b")},
		{name: "two chars", args: args{p: []byte("ab")}, wantBytes: []byte("ab\b\b")},
		{name: "non standard char", args: args{p: []byte("ü")}, wantBytes: []byte("ü\b")},
		//{name: "tab", args: args{p: []byte("\t")}, wantBytes: []byte("\t\b")}, not supported
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var b bytes.Buffer
			slw := MakeSingleLineWriter(&b)
			slw.Write(tt.args.p)
			slw.Cleanup()

			//fmt.Print("-->")
			//slwStdOut := MakeSingleLineWriter(os.Stdout)
			//slwStdOut.Write(tt.args.p)
			//slwStdOut.Cleanup()
			//fmt.Print("<--")

			gotBytes := b.Bytes()
			if !reflect.DeepEqual(gotBytes, tt.wantBytes) {
				t.Errorf("Write() bytes = %v (%s), want %v (%s)", gotBytes, string(gotBytes), tt.wantBytes, string(tt.wantBytes))
			}
		})
	}
}
