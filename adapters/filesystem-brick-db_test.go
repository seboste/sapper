package adapters

import (
	"log"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/seboste/sapper/ports"
)

func TestParseString(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name  string
		args  args
		want  BrickKind
		want1 bool
	}{
		{name: "template", args: args{str: "template"}, want: BrickKind(ports.Template), want1: true},
		{name: "extension", args: args{str: "extension"}, want: BrickKind(ports.Extension), want1: true},
		{name: "template CaMeLcAsE", args: args{str: "TeMplAtE"}, want: BrickKind(ports.Template), want1: true},
		{name: "unknown", args: args{str: "bla"}, want: BrickKind(0), want1: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := ParseString(tt.args.str)
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
		{name: "template", bk: BrickKind(ports.Template), want: "template"},
		{name: "extension", bk: BrickKind(ports.Extension), want: "extension"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.bk.String(); got != tt.want {
				t.Errorf("BrickKind.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_makeFilesystemBrick(t *testing.T) {
	tests := []struct {
		name    string
		yaml    string
		files   []string
		want    filesystemBrick
		wantErr bool
	}{
		{name: "basic brick",
			yaml: `id : test
kind: extension
description : My test brick
version : 1.0.0
parameters : 
 - name : param1
   default : default
 - name : param2
dependencies :
 - dep1
 - dep2`,
			files: []string{"a", "b/c"},
			want: filesystemBrick{
				Id:           "test",
				Description:  "My test brick",
				Version:      "1.0.0",
				Kind:         BrickKind(ports.Extension),
				Parameters:   []ports.BrickParameters{{Name: "param1", Default: "default"}, {Name: "param2", Default: ""}},
				Dependencies: []string{"dep1", "dep2"},
				Files:        []string{"a", "b/c"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			//1. prepare test
			dir, _ := os.MkdirTemp("", "example")
			defer os.RemoveAll(dir) // clean up
			os.WriteFile(filepath.Join(dir, "manifest.yaml"), []byte(tt.yaml), 0666)

			for _, file := range tt.files {
				abspath := filepath.Join(dir, file)
				filedir, _ := filepath.Split(abspath)

				if err := os.MkdirAll(filedir, 0777); err != nil {
					log.Fatalln(err)
				}
				if err := os.WriteFile(abspath, []byte{}, 0666); err != nil {
					log.Fatalln(err)
				}
			}

			//2. execute test
			got, err := makeFilesystemBrick(dir)
			if (err != nil) != tt.wantErr {
				t.Errorf("makeFilesystemBrick() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("makeFilesystemBrick() = %v, want %v", got, tt.want)
			}
		})
	}
}
