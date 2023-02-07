package adapters

import "github.com/seboste/sapper/ports"

type AggregateBrickDB struct {
}

func MakeAggregateBrickDB(remotes []ports.Remote) *AggregateBrickDB {
	return &AggregateBrickDB{}
}

func (abdb AggregateBrickDB) Init(Path string) error {
	return nil
}

func (abdb AggregateBrickDB) Bricks(kind ports.BrickKind) []ports.Brick {
	return []ports.Brick{}
}

func (abdb AggregateBrickDB) Brick(id string) (ports.Brick, error) {
	return ports.Brick{}, nil
}

var _ ports.BrickDB = AggregateBrickDB{}
