package core

import (
	"reflect"
	"testing"

	"github.com/seboste/sapper/ports"
)

func Test_removeBricks(t *testing.T) {
	type args struct {
		bricks           []ports.Brick
		brickIdsToRemove []ports.BrickDependency
	}
	tests := []struct {
		name string
		args args
		want []ports.Brick
	}{
		{name: "remove two",
			args: args{
				bricks:           []ports.Brick{{Id: "1"}, {Id: "2"}, {Id: "3"}},
				brickIdsToRemove: []ports.BrickDependency{{Id: "2"}, {Id: "3"}},
			},
			want: []ports.Brick{{Id: "1"}},
		},
		{name: "remove not available",
			args: args{
				bricks:           []ports.Brick{{Id: "1"}, {Id: "2"}, {Id: "3"}},
				brickIdsToRemove: []ports.BrickDependency{{Id: "4"}},
			},
			want: []ports.Brick{{Id: "1"}, {Id: "2"}, {Id: "3"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := removeBricks(tt.args.bricks, tt.args.brickIdsToRemove); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("removeBricks() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBrickApi_Add(t *testing.T) {
	type fields struct {
		Db                 ports.BrickDB
		ServicePersistence ports.ServicePersistence
	}
	type args struct {
		servicePath       string
		brickId           string
		parameterResolver ports.ParameterResolver
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := BrickApi{
				Db:                 tt.fields.Db,
				ServicePersistence: tt.fields.ServicePersistence,
			}
			if err := b.Add(tt.args.servicePath, tt.args.brickId, tt.args.parameterResolver); (err != nil) != tt.wantErr {
				t.Errorf("BrickApi.Add() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
