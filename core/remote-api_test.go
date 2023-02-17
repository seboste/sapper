package core

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/seboste/sapper/ports"
)

func TestAdd(t *testing.T) {

	abc := []ports.Remote{{Name: "a"}, {Name: "b"}, {Name: "c"}}

	type args struct {
		remotes []ports.Remote
		r       ports.Remote
		pos     int
	}
	tests := []struct {
		name string
		args args
		want []ports.Remote
	}{
		{name: "negative position", args: args{remotes: abc, r: ports.Remote{Name: "new"}, pos: -1}, want: []ports.Remote{{Name: "a"}, {Name: "b"}, {Name: "c"}, {Name: "new"}}},
		{name: "first position", args: args{remotes: abc, r: ports.Remote{Name: "new"}, pos: 0}, want: []ports.Remote{{Name: "new"}, {Name: "a"}, {Name: "b"}, {Name: "c"}}},
		{name: "last position", args: args{remotes: abc, r: ports.Remote{Name: "new"}, pos: 3}, want: []ports.Remote{{Name: "a"}, {Name: "b"}, {Name: "c"}, {Name: "new"}}},
		{name: "somewhere in between", args: args{remotes: abc, r: ports.Remote{Name: "new"}, pos: 1}, want: []ports.Remote{{Name: "a"}, {Name: "new"}, {Name: "b"}, {Name: "c"}}},
		{name: "pos greater than size", args: args{remotes: abc, r: ports.Remote{Name: "new"}, pos: 28}, want: []ports.Remote{{Name: "a"}, {Name: "b"}, {Name: "c"}, {Name: "new"}}},
		{name: "empty", args: args{remotes: nil, r: ports.Remote{Name: "new"}, pos: 0}, want: []ports.Remote{{Name: "new"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Add(tt.args.remotes, tt.args.r, tt.args.pos); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Add() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRemoteApi_Add(t *testing.T) {
	type fields struct {
		Configuration  ports.Configuration
		BrickDBFactory ports.BrickDBFactory
	}
	type args struct {
		name     string
		src      string
		position int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := RemoteApi{
				Configuration:  tt.fields.Configuration,
				BrickDBFactory: tt.fields.BrickDBFactory,
			}
			if err := r.Add(tt.args.name, tt.args.src, tt.args.position); (err != nil) != tt.wantErr {
				t.Errorf("RemoteApi.Add() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_inferKind(t *testing.T) {

	dir, _ := os.MkdirTemp("", "testInferKind*")
	defer os.RemoveAll(dir)

	file, _ := os.CreateTemp("", "testInferKind*.txt")
	defer os.Remove(file.Name())
	defer file.Close()

	type args struct {
		src string
	}
	tests := []struct {
		name     string
		args     args
		wantKind ports.RemoteKind
		wantErr  bool
	}{
		{name: "valid Git URL", args: args{src: "https://github.com/seboste/sapper-bricks.git"}, wantKind: ports.GitRemote, wantErr: false},
		{name: "valid dir path", args: args{src: dir}, wantKind: ports.FilesystemRemote, wantErr: false},
		{name: "non existing path", args: args{src: filepath.Join(dir, "non-existent")}, wantKind: -1, wantErr: true},
		{name: "some file", args: args{src: file.Name()}, wantKind: -1, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotKind, err := inferKind(tt.args.src)
			if (err != nil) != tt.wantErr {
				t.Errorf("inferKind() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotKind, tt.wantKind) {
				t.Errorf("inferKind() = %v, want %v", gotKind, tt.wantKind)
			}
		})
	}
}
