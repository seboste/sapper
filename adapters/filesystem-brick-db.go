package adapters

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/seboste/sapper/ports"
	"gopkg.in/yaml.v3"
)

func makeBrick(path string) (ports.Brick, error) {
	b := ports.Brick{}

	yamlFile, err := ioutil.ReadFile(filepath.Join(path, "manifest.yaml"))
	if err != nil {
		return ports.Brick(b), err
	}
	err = yaml.Unmarshal(yamlFile, &b)
	if err != nil {
		return b, err
	}
	b.BasePath = filepath.Join(path, "")
	err = filepath.Walk(path,
		func(p string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			_, file := filepath.Split(p)

			if file == "manifest.yaml" {
				return nil //skip manifest.yaml
			}

			if info.IsDir() {
				//skip directories
				return nil
			}

			relPath, _ := filepath.Rel(path, p)
			b.Files = append(b.Files, relPath)
			return nil
		})

	return b, nil
}

type FilesystemBrickDB struct {
	bricks []ports.Brick
}

func (db *FilesystemBrickDB) Init(basePath string) error {
	err := filepath.Walk(basePath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			dir, file := filepath.Split(path)

			if file == "manifest.yaml" {
				brick, err := makeBrick(dir)
				if err != nil {
					return err
				}

				db.bricks = append(db.bricks, brick)
			}

			return nil
		})
	return err
}

func (db *FilesystemBrickDB) Bricks(kind ports.BrickKind) []ports.Brick {
	filteredBricks := []ports.Brick{}
	for _, b := range db.bricks {
		if b.Kind == kind {
			filteredBricks = append(filteredBricks, b)
		}
	}
	return filteredBricks
}

func (db *FilesystemBrickDB) Brick(id string) (ports.Brick, error) {
	for _, b := range db.bricks {
		if b.Id == id {
			return b, nil
		}
	}
	return ports.Brick{}, ports.BrickNotFound
}

var _ ports.BrickDB = &FilesystemBrickDB{}
