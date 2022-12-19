package ports

import (
	"reflect"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestParseBrickKind(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name  string
		args  args
		want  BrickKind
		want1 bool
	}{
		{name: "template", args: args{str: "template"}, want: BrickKind(Template), want1: true},
		{name: "extension", args: args{str: "extension"}, want: BrickKind(Extension), want1: true},
		{name: "template CaMeLcAsE", args: args{str: "TeMplAtE"}, want: BrickKind(Template), want1: true},
		{name: "unknown", args: args{str: "bla"}, want: BrickKind(0), want1: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := ParseBrickKind(tt.args.str)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseString() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("ParseString() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestBrickKind_String(t *testing.T) {
	tests := []struct {
		name string
		bk   BrickKind
		want string
	}{
		{name: "template", bk: BrickKind(Template), want: "template"},
		{name: "extension", bk: BrickKind(Extension), want: "extension"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.bk.String(); got != tt.want {
				t.Errorf("BrickKind.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBrickKind_UnmarshalYAML(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name    string
		args    args
		wantBk  BrickKind
		wantErr bool
	}{
		{name: "extension", args: args{value: "extension"}, wantBk: Extension, wantErr: false},
		{name: "template", args: args{value: "template"}, wantBk: Template, wantErr: false},
		{name: "invalid", args: args{value: "something_invalid"}, wantBk: Template, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var bk BrickKind
			if err := yaml.Unmarshal([]byte(tt.args.value), &bk); (err != nil) != tt.wantErr {
				t.Errorf("BrickKind.UnmarshalYAML() error = %v, wantErr %v", err, tt.wantErr)
			}
			if bk != tt.wantBk {
				t.Errorf("BrickKind.UnmarshalYAML() brickKind = %v, wantErr %v", bk, tt.wantBk)
			}
		})
	}
}
