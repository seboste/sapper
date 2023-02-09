package adapters

import (
	"github.com/seboste/sapper/ports"
)

func makeBrickDB(r ports.Remote) (ports.BrickDB, error) {
	fbdb := FilesystemBrickDB{}
	fbdb.Init(r.Path)
	return &fbdb, nil
}

func MakeBrickDB(remotes []ports.Remote) (ports.BrickDB, error) {
	abdb := AggregateBrickDB{}
	for _, r := range remotes {
		db, err := makeBrickDB(r)
		if err != nil {
			return abdb, err
		}
		abdb.dbs = append(abdb.dbs, db)
	}
	return abdb, nil
}
