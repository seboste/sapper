package brickDb

import (
	"os"
	"os/exec"

	"github.com/seboste/sapper/ports"
)

type GitBrickDB struct {
	FilesystemBrickDB
	Path string
	Url  string
}

func (gbdb GitBrickDB) Clone() error {
	cmd := exec.Command("git", "clone", gbdb.Url, gbdb.Path)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return err
}

func (gbdb GitBrickDB) Update() error {
	cmd := exec.Command("git", "pull")
	cmd.Path = gbdb.Path
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return err
}

func MakeGitBrickDB(path string, url string) (GitBrickDB, error) {
	db := GitBrickDB{Path: path, Url: url}
	if _, err := os.Stat(db.Path); os.IsNotExist(err) {
		err = db.Clone()
		if err != nil {
			return db, err
		}
	}
	var err error
	db.FilesystemBrickDB, err = MakeFilesystemBrickDB(db.Path)
	return db, err
}

var _ ports.BrickDB = &GitBrickDB{}
