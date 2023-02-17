package brickDb

import (
	"reflect"
	"testing"

	"github.com/seboste/sapper/ports"
)

type MockBrickDB struct {
	BricksMap map[ports.BrickKind][]ports.Brick //kind->bricks
}

func (db MockBrickDB) Bricks(k ports.BrickKind) []ports.Brick {
	return db.BricksMap[k]
}

func (db MockBrickDB) Brick(id string) (ports.Brick, error) {
	for _, bricks := range db.BricksMap {
		for _, brick := range bricks {
			if brick.Id == id {
				return brick, nil
			}
		}
	}
	return ports.Brick{}, ports.BrickNotFound
}

func (db MockBrickDB) Update() error {
	return nil
}

func (db MockBrickDB) IsModified() (bool, string) {
	return false, ""
}

func TestAggregateBrickDB_Bricks(t *testing.T) {
	bricksTemplate := []ports.Brick{{Id: "TemplateA"}, {Id: "TemplateB"}}
	bricksExtension := []ports.Brick{{Id: "ExtensionA"}, {Id: "ExtensionB"}}
	db1 := MockBrickDB{BricksMap: map[ports.BrickKind][]ports.Brick{ports.Template: bricksTemplate, ports.Extension: bricksExtension}}
	bricksTemplate2 := []ports.Brick{{Id: "TemplateC", Version: "2.0.0"}, {Id: "TemplateD", Version: "2.0.0"}}
	bricksExtension2 := []ports.Brick{{Id: "ExtensionA", Version: "2.0.0"}, {Id: "ExtensionC", Version: "2.0.0"}} //duplicate
	db2 := MockBrickDB{BricksMap: map[ports.BrickKind][]ports.Brick{ports.Template: bricksTemplate2, ports.Extension: bricksExtension2}}

	combinedTemplateBricks := []ports.Brick{{Id: "TemplateA"}, {Id: "TemplateB"}, {Id: "TemplateC", Version: "2.0.0"}, {Id: "TemplateD", Version: "2.0.0"}}
	combinedExtensionBricks := []ports.Brick{{Id: "ExtensionA"}, {Id: "ExtensionB"}, {Id: "ExtensionC", Version: "2.0.0"}}

	type fields struct {
		dbs []ports.BrickDB
	}
	type args struct {
		k ports.BrickKind
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []ports.Brick
	}{
		{name: "single DB", fields: fields{dbs: []ports.BrickDB{db1}}, args: args{k: ports.Template}, want: bricksTemplate},
		{name: "two DBs, disjoint bricks", fields: fields{dbs: []ports.BrickDB{db1, db2}}, args: args{k: ports.Template}, want: combinedTemplateBricks},
		{name: "two DBs, overlapping bricks (first entry has precedence)", fields: fields{dbs: []ports.BrickDB{db1, db2}}, args: args{k: ports.Extension}, want: combinedExtensionBricks},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			abdb := AggregateBrickDB{
				dbs: tt.fields.dbs,
			}
			if got := abdb.Bricks(tt.args.k); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AggregateBrickDB.Bricks() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAggregateBrickDB_Brick(t *testing.T) {

	bricksTemplate := []ports.Brick{{Id: "TemplateA"}, {Id: "TemplateB"}}
	bricksExtension := []ports.Brick{{Id: "ExtensionA"}, {Id: "ExtensionB"}}
	db1 := MockBrickDB{BricksMap: map[ports.BrickKind][]ports.Brick{ports.Template: bricksTemplate, ports.Extension: bricksExtension}}
	bricksTemplate2 := []ports.Brick{{Id: "TemplateC", Version: "2.0.0"}, {Id: "TemplateD", Version: "2.0.0"}}
	bricksExtension2 := []ports.Brick{{Id: "ExtensionA", Version: "2.0.0"}, {Id: "ExtensionC", Version: "2.0.0"}} //duplicate
	db2 := MockBrickDB{BricksMap: map[ports.BrickKind][]ports.Brick{ports.Template: bricksTemplate2, ports.Extension: bricksExtension2}}

	type fields struct {
		dbs []ports.BrickDB
	}
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    ports.Brick
		wantErr bool
	}{
		{name: "single DB", fields: fields{dbs: []ports.BrickDB{db1}}, args: args{id: "ExtensionA"}, want: bricksExtension[0], wantErr: false},
		{name: "missing brick", fields: fields{dbs: []ports.BrickDB{db1}}, args: args{id: "ExtensionZ"}, want: ports.Brick{}, wantErr: true},
		{name: "two DBs, get from db1", fields: fields{dbs: []ports.BrickDB{db1, db2}}, args: args{id: "TemplateA"}, want: bricksTemplate[0], wantErr: false},
		{name: "two DBs, get from db2", fields: fields{dbs: []ports.BrickDB{db1, db2}}, args: args{id: "ExtensionC"}, want: bricksExtension2[1], wantErr: false},
		{name: "two DBs, prefer db1", fields: fields{dbs: []ports.BrickDB{db1, db2}}, args: args{id: "ExtensionA"}, want: bricksExtension[0], wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			abdb := AggregateBrickDB{
				dbs: tt.fields.dbs,
			}
			got, err := abdb.Brick(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("AggregateBrickDB.Brick() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AggregateBrickDB.Brick() = %v, want %v", got, tt.want)
			}
		})
	}
}
