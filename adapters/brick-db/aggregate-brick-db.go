package brickDb

import "github.com/seboste/sapper/ports"

type AggregateBrickDB struct {
	dbs []ports.BrickDB
}

func contains(bricks []ports.Brick, b ports.Brick) bool {
	for _, brick := range bricks {
		if brick.Id == b.Id {
			return true
		}
	}
	return false
}

func (abdb AggregateBrickDB) Bricks(k ports.BrickKind) []ports.Brick {
	bricks := []ports.Brick{}
	for _, db := range abdb.dbs {
		for _, b := range db.Bricks(k) {
			if !contains(bricks, b) {
				bricks = append(bricks, b)
			}
		}
	}
	return bricks
}

func (abdb AggregateBrickDB) Brick(id string) (ports.Brick, error) {
	for _, db := range abdb.dbs {
		brick, err := db.Brick(id)
		if err == nil {
			return brick, nil
		}
	}
	return ports.Brick{}, ports.BrickNotFound
}

func (abdb AggregateBrickDB) Update() error {
	for _, db := range abdb.dbs {
		if err := db.Update(); err != nil {
			return err
		}
	}
	return nil
}

func (abdb AggregateBrickDB) IsModified() (bool, string) {
	modified := false
	details := ""
	for _, db := range abdb.dbs {
		if m, d := db.IsModified(); m {
			modified = true
			details = details + d
		}
	}
	return modified, details
}

var _ ports.BrickDB = AggregateBrickDB{}
