package adapters

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/seboste/sapper/ports"
)

func Test_makeFilesystemBrick(t *testing.T) {
	tempDir, _ := ioutil.TempDir("", "example")
	defer os.RemoveAll(tempDir) // clean up

	tests := []struct {
		name    string
		yaml    string
		files   []string
		want    ports.Brick
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
			want: ports.Brick{
				Id:           "test",
				Description:  "My test brick",
				Version:      "1.0.0",
				Kind:         ports.Extension,
				Parameters:   []ports.BrickParameters{{Name: "param1", Default: "default"}, {Name: "param2", Default: ""}},
				Dependencies: []string{"dep1", "dep2"},
				BasePath:     filepath.Join(tempDir, "test"),
				Files:        []string{"a", "b/c"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			//1. prepare test
			brickDir := filepath.Join(tempDir, tt.want.Id)
			if err := os.MkdirAll(brickDir, 0777); err != nil {
				log.Fatalln(err)
			}
			ioutil.WriteFile(filepath.Join(brickDir, "manifest.yaml"), []byte(tt.yaml), 0666)

			for _, file := range tt.files {
				abspath := filepath.Join(brickDir, file)
				filedir, _ := filepath.Split(abspath)

				if err := os.MkdirAll(filedir, 0777); err != nil {
					log.Fatalln(err)
				}
				if err := ioutil.WriteFile(abspath, []byte{}, 0666); err != nil {
					log.Fatalln(err)
				}
			}

			//2. execute test
			got, err := makeBrick(brickDir)
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

func TestFilesystemBrickDB_Bricks(t *testing.T) {
	brickTemp := ports.Brick{Id: "test1_templ", Description: "desc", Kind: ports.BrickKind(ports.Template)}
	brickExt := ports.Brick{Id: "test2_ext", Description: "desc", Kind: ports.BrickKind(ports.Extension)}

	type fields struct {
		bricks []ports.Brick
	}
	type args struct {
		kind ports.BrickKind
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []ports.Brick
	}{
		{name: "filter extensions",
			fields: fields{bricks: []ports.Brick{brickTemp, brickExt}},
			args:   args{kind: ports.Extension},
			want:   []ports.Brick{brickExt},
		},
		{name: "filter template",
			fields: fields{bricks: []ports.Brick{brickTemp, brickExt}},
			args:   args{kind: ports.Template},
			want:   []ports.Brick{brickTemp},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &FilesystemBrickDB{
				bricks: tt.fields.bricks,
			}
			if got := db.Bricks(tt.args.kind); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FilesystemBrickDB.Bricks() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilesystemBrickDB_Brick(t *testing.T) {
	someBrick := ports.Brick{Id: "some_id"}
	someOtherBrick := ports.Brick{Id: "some_other_id"}
	type fields struct {
		bricks []ports.Brick
	}
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		want    ports.Brick
	}{
		{name: "brick available",
			fields:  fields{bricks: []ports.Brick{someBrick, someOtherBrick}},
			args:    args{id: "some_id"},
			wantErr: false,
			want:    someBrick,
		},
		{name: "brick not available",
			fields:  fields{bricks: []ports.Brick{someBrick, someOtherBrick}},
			args:    args{id: "some_unknown_id"},
			wantErr: true,
			want:    ports.Brick{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &FilesystemBrickDB{
				bricks: tt.fields.bricks,
			}

			got, err := db.Brick(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("FilesystemBrickDB.Brick() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FilesystemBrickDB.Brick() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilesystemBrickDB_Init(t *testing.T) {

	tempDir, _ := ioutil.TempDir("", "example_db")
	defer os.RemoveAll(tempDir) // clean up

	initiallyAvailableBrick := ports.Brick{Id: "initial", Kind: ports.Extension}

	type fields struct {
		bricks []ports.Brick
	}
	type args struct {
		basePath string
	}
	type input struct {
		name string
		yaml string
	}
	tests := []struct {
		name    string
		input   []input
		fields  fields
		args    args
		wantErr bool
		want    []ports.Brick
	}{
		{
			name:   "init",
			fields: fields{bricks: []ports.Brick{initiallyAvailableBrick}},
			args:   args{basePath: filepath.Join(tempDir, "example_db")},
			input: []input{{name: "brick_1", yaml: `id : brick_1
kind: extension
description : test brick 1
`},
				{name: "brick_2", yaml: `id : brick_2
kind: extension
description : test brick 2
`},
			},
			wantErr: false,
			want: []ports.Brick{
				initiallyAvailableBrick,
				{Id: "brick_1", Description: "test brick 1", Kind: ports.Extension, BasePath: filepath.Join(tempDir, "example_db/brick_1")},
				{Id: "brick_2", Description: "test brick 2", Kind: ports.Extension, BasePath: filepath.Join(tempDir, "example_db/brick_2")},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//1. prepare test
			os.Mkdir(tt.args.basePath, 0777)
			for _, in := range tt.input {
				brickDir := filepath.Join(tt.args.basePath, in.name)
				os.Mkdir(brickDir, 0777)
				ioutil.WriteFile(filepath.Join(brickDir, "manifest.yaml"), []byte(in.yaml), 0666)
			}
			//2. execute test
			db := &FilesystemBrickDB{
				bricks: tt.fields.bricks,
			}
			if err := db.Init(tt.args.basePath); (err != nil) != tt.wantErr {
				t.Errorf("FilesystemBrickDB.Init() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(db.bricks, tt.want) {
				t.Errorf("FilesystemBrickDB.bricks = %v, want %v", db.bricks, tt.want)
			}
		})
	}
}
