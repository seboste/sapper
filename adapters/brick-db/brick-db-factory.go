package brickDb

import (
	"path/filepath"

	"github.com/seboste/sapper/ports"
)

func makeBrickDB(r ports.Remote, remotesDir string) (ports.BrickDB, error) {

	switch r.Kind {
	case ports.GitRemote:
		gbdb, err := MakeGitBrickDB(filepath.Join(remotesDir, r.Name), r.Src)
		return &gbdb, err

	case ports.FilesystemRemote:
		fbdb, err := MakeFilesystemBrickDB(r.Src)
		return &fbdb, err
	}

	return nil, nil
}

func MakeBrickDB(remotes []ports.Remote, remotesDir string) (ports.BrickDB, error) {
	abdb := AggregateBrickDB{}
	for _, r := range remotes {
		db, err := makeBrickDB(r, remotesDir)
		if err != nil {
			return abdb, err
		}
		abdb.dbs = append(abdb.dbs, db)
	}
	return abdb, nil
}
