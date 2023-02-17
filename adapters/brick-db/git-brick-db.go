package brickDb

import (
	"bytes"
	"fmt"
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
	cmd.Dir = gbdb.Path
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return err
}

func (gbdb GitBrickDB) IsModified() (bool, string) {

	cmd := exec.Command("git", "status", "-s", "--porcelain")
	cmd.Dir = gbdb.Path
	var buffer bytes.Buffer
	cmd.Stdout = &buffer
	cmd.Stderr = &buffer
	cmd.Run()
	changes := buffer.String()

	details := fmt.Sprintf("the following changes have been detected:\n%smake sure to commit %s", changes, gbdb.Path)
	return changes != "", details
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
