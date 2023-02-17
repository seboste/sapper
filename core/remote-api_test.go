package core

import (
	"io/ioutil"
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

func Test_inferKind(t *testing.T) {

	dir, _ := ioutil.TempDir("", "testInferKind*")
	defer os.RemoveAll(dir)

	file, _ := ioutil.TempFile("", "testInferKind*.txt")
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

func Test_findRemote(t *testing.T) {
	type args struct {
		remotes []ports.Remote
		name    string
	}
	tests := []struct {
		name       string
		args       args
		wantIndex  int
		wantRemote ports.Remote
		wantOk     bool
	}{
		{name: "present", args: args{remotes: []ports.Remote{{Name: "a"}, {Name: "b"}, {Name: "c"}}, name: "b"}, wantIndex: 1, wantRemote: ports.Remote{Name: "b"}, wantOk: true},
		{name: "absent", args: args{remotes: []ports.Remote{{Name: "a"}, {Name: "b"}, {Name: "c"}}, name: "d"}, wantIndex: -1, wantRemote: ports.Remote{}, wantOk: false},
		{name: "empty", args: args{remotes: []ports.Remote{}, name: "a"}, wantIndex: -1, wantRemote: ports.Remote{}, wantOk: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotIndex, gotRemote, gotOk := findRemote(tt.args.remotes, tt.args.name)
			if gotIndex != tt.wantIndex {
				t.Errorf("findRemote() gotIndex = %v, want %v", gotIndex, tt.wantIndex)
			}
			if !reflect.DeepEqual(gotRemote, tt.wantRemote) {
				t.Errorf("findRemote() gotRemote = %v, want %v", gotRemote, tt.wantRemote)
			}
			if gotOk != tt.wantOk {
				t.Errorf("findRemote() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

type MockConfiguration struct {
	saveCalled bool
	remotes    []ports.Remote
}

func (c *MockConfiguration) Save() error {
	c.saveCalled = true
	return nil
}

func (c MockConfiguration) DefaultRemotesDir() string {
	return os.TempDir()
}

func (c MockConfiguration) Remotes() []ports.Remote {
	return c.remotes
}
func (c *MockConfiguration) UpdateRemotes(remotes []ports.Remote) {
	c.remotes = remotes
}

var _ ports.Configuration = (*MockConfiguration)(nil)

type MockBrickDBFactory struct {
	brickDB TestBrickDB
}

func (mbdbf MockBrickDBFactory) MakeBrickDB(r ports.Remote, remotesDir string) (ports.BrickDB, error) {
	return &mbdbf.brickDB, nil
}
func (mbdbf MockBrickDBFactory) MakeAggregatedBrickDB(r []ports.Remote, remotesDir string) (ports.BrickDB, error) {
	return &mbdbf.brickDB, nil
}

var _ ports.BrickDBFactory = MockBrickDBFactory{}

func TestRemoteApi_Add(t *testing.T) {

	type fields struct {
		InitialRemotes []ports.Remote
		BrickDBFactory ports.BrickDBFactory
	}

	flds := fields{InitialRemotes: []ports.Remote{{Name: "a"}}, BrickDBFactory: MockBrickDBFactory{}}

	type args struct {
		name     string
		src      string
		position int
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		wantErr        bool
		wantSaveCalled bool
		wantRemotes    []ports.Remote
	}{
		{name: "add after", fields: flds, args: args{name: "new", src: "someurl.git", position: -1}, wantErr: false, wantSaveCalled: true, wantRemotes: []ports.Remote{{Name: "a"}, {Name: "new", Src: "someurl.git", Kind: ports.GitRemote}}},
		{name: "add before", fields: flds, args: args{name: "new", src: "someurl.git", position: 0}, wantErr: false, wantSaveCalled: true, wantRemotes: []ports.Remote{{Name: "new", Src: "someurl.git", Kind: ports.GitRemote}, {Name: "a"}}},
		{name: "add invalid", fields: flds, args: args{name: "new", src: "someurl.invalid", position: 0}, wantErr: true, wantSaveCalled: false, wantRemotes: []ports.Remote{{Name: "a"}}},
		{name: "add existing", fields: flds, args: args{name: "a", src: "someurl.git", position: 0}, wantErr: true, wantSaveCalled: false, wantRemotes: []ports.Remote{{Name: "a"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := MockConfiguration{remotes: tt.fields.InitialRemotes}
			r := RemoteApi{
				Configuration:  &mc,
				BrickDBFactory: tt.fields.BrickDBFactory,
			}
			if err := r.Add(tt.args.name, tt.args.src, tt.args.position); (err != nil) != tt.wantErr {
				t.Errorf("RemoteApi.Add() error = %v, wantErr %v", err, tt.wantErr)
			}

			gotRemotes := r.Configuration.Remotes()
			if !reflect.DeepEqual(gotRemotes, tt.wantRemotes) {
				t.Errorf("RemoteApi.Add() remotes = %v, wantRemotes %v", gotRemotes, tt.wantRemotes)
			}

			gotSaveCalled := mc.saveCalled
			if gotSaveCalled != tt.wantSaveCalled {
				t.Errorf("RemoteApi.Add() saveCalled = %v, wantSaveCalled %v", gotSaveCalled, tt.wantSaveCalled)
			}
		})
	}
}

func TestRemoteApi_Remove(t *testing.T) {
	type fields struct {
		InitialRemotes []ports.Remote
		BrickDBFactory ports.BrickDBFactory
	}

	flds := fields{InitialRemotes: []ports.Remote{{Name: "a"}, {Name: "b"}}, BrickDBFactory: MockBrickDBFactory{}}

	type args struct {
		name string
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		wantErr        bool
		wantSaveCalled bool
		wantRemotes    []ports.Remote
	}{
		{name: "present", fields: flds, args: args{name: "b"}, wantErr: false, wantSaveCalled: true, wantRemotes: []ports.Remote{{Name: "a"}}},
		{name: "absent", fields: flds, args: args{name: "missing"}, wantErr: true, wantSaveCalled: false, wantRemotes: []ports.Remote{{Name: "a"}, {Name: "b"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := MockConfiguration{remotes: tt.fields.InitialRemotes}
			r := RemoteApi{
				Configuration:  &mc,
				BrickDBFactory: tt.fields.BrickDBFactory,
			}
			if err := r.Remove(tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("RemoteApi.Remove() error = %v, wantErr %v", err, tt.wantErr)
			}

			gotRemotes := r.Configuration.Remotes()
			if !reflect.DeepEqual(gotRemotes, tt.wantRemotes) {
				t.Errorf("RemoteApi.Remove() remotes = %v, wantRemotes %v", gotRemotes, tt.wantRemotes)
			}

			gotSaveCalled := mc.saveCalled
			if gotSaveCalled != tt.wantSaveCalled {
				t.Errorf("RemoteApi.Remove() saveCalled = %v, wantSaveCalled %v", gotSaveCalled, tt.wantSaveCalled)
			}
		})
	}
}

func TestRemoteApi_Update(t *testing.T) {
	type fields struct {
		InitialRemotes []ports.Remote
	}

	type args struct {
		name string
	}
	tests := []struct {
		name             string
		fields           fields
		args             args
		wantErr          bool
		wantUpdateCalled bool
	}{
		{name: "present", fields: fields{InitialRemotes: []ports.Remote{{Name: "a"}, {Name: "b"}}}, args: args{name: "a"}, wantErr: false, wantUpdateCalled: true},
		{name: "absent", fields: fields{InitialRemotes: []ports.Remote{{Name: "a"}, {Name: "b"}}}, args: args{name: "invalid"}, wantErr: true, wantUpdateCalled: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := MockConfiguration{remotes: tt.fields.InitialRemotes}
			mbdf := MockBrickDBFactory{}
			gotUpdateCalled := false
			mbdf.brickDB.updateCalled = &gotUpdateCalled
			r := RemoteApi{
				Configuration:  &mc,
				BrickDBFactory: &mbdf,
			}

			if err := r.Update(tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("RemoteApi.Update() error = %v, wantErr %v", err, tt.wantErr)
			}

			if gotUpdateCalled != tt.wantUpdateCalled {
				t.Errorf("RemoteApi.Update() updateCalled = %v, wantUpdateCalled %v", gotUpdateCalled, tt.wantUpdateCalled)
			}
		})
	}
}
