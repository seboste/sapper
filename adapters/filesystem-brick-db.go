package adapters

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/seboste/sapper/ports"
	"gopkg.in/yaml.v3"
)

type BrickKind ports.BrickKind

var (
	brickKindMap = map[string]BrickKind{
		"template":  BrickKind(ports.Template),
		"extension": BrickKind(ports.Extension),
	}
)

func ParseString(str string) (BrickKind, bool) {
	c, ok := brickKindMap[strings.ToLower(str)]
	return c, ok
}

func (bk BrickKind) String() string {
	switch ports.BrickKind(bk) {
	case ports.Template:
		return "template"
	case ports.Extension:
		return "extension"
	default:
		return fmt.Sprintf("%d", int(bk))
	}
}

func (bk *BrickKind) UnmarshalYAML(value *yaml.Node) error {
	ok := false
	*bk, ok = ParseString(value.Value)
	if !ok {
		return fmt.Errorf("invalid brick kind %s", value.Value)
	}
	return nil
}

type filesystemBrick struct {
	Id           string
	Description  string
	Version      string
	Kind         BrickKind
	Parameters   []ports.BrickParameters
	Dependencies []string
	Files        []string
}

func (b filesystemBrick) GetId() string                          { return b.Id }
func (b filesystemBrick) GetDescription() string                 { return b.Description }
func (b filesystemBrick) GetVersion() string                     { return b.Version }
func (b filesystemBrick) GetKind() ports.BrickKind               { return ports.BrickKind(b.Kind) }
func (b filesystemBrick) GetParameters() []ports.BrickParameters { return b.Parameters }
func (b filesystemBrick) GetDependencies() []string              { return b.Dependencies }
func (b filesystemBrick) GetFiles() []string                     { return b.Files }

var _ ports.Brick = filesystemBrick{}

func makeFilesystemBrick(path string) (filesystemBrick, error) {
	b := filesystemBrick{}

	yamlFile, err := ioutil.ReadFile(filepath.Join(path, "manifest.yaml"))
	if err != nil {
		return b, err
	}
	err = yaml.Unmarshal(yamlFile, &b)
	if err != nil {
		return b, err
	}

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
				brick, err := makeFilesystemBrick(dir)
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
		if b.GetKind() == kind {
			filteredBricks = append(filteredBricks, b)
		}
	}
	return filteredBricks
}

func (db *FilesystemBrickDB) Brick(id string) ports.Brick {
	for _, b := range db.bricks {
		if b.GetId() == id {
			return b
		}
	}
	return nil
}

var _ ports.BrickDB = &FilesystemBrickDB{}
