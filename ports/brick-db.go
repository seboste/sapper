package ports

import (
	"errors"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

var BrickNotFound = errors.New("brick not found")

type BrickParameters struct {
	Name    string
	Default string
}

type BrickKind int

const (
	Template BrickKind = iota
	Extension
)

var BrickKinds = []BrickKind{Template, Extension}

type Brick struct {
	Id           string
	Description  string
	Version      string
	Kind         BrickKind
	Parameters   []BrickParameters
	Dependencies []string
	BasePath     string
	Files        []string
}

type BrickDBFactory interface {
	MakeBrickDB(r Remote, remotesDir string) (BrickDB, error)
	MakeAggregatedBrickDB(r []Remote, remotesDir string) (BrickDB, error)
}

type BrickDB interface {
	Bricks(kind BrickKind) []Brick
	Brick(id string) (Brick, error)
	Update() error
	IsModified() (bool, string)
}

var (
	brickKindMap = map[string]BrickKind{
		"template":  BrickKind(Template),
		"extension": BrickKind(Extension),
	}
)

func ParseBrickKind(str string) (BrickKind, bool) {
	c, ok := brickKindMap[strings.ToLower(str)]
	return c, ok
}

func (bk BrickKind) String() string {
	switch BrickKind(bk) {
	case Template:
		return "template"
	case Extension:
		return "extension"
	default:
		return fmt.Sprintf("%d", int(bk))
	}
}

func (bk *BrickKind) UnmarshalYAML(value *yaml.Node) error {
	ok := false
	*bk, ok = ParseBrickKind(value.Value)
	if !ok {
		return fmt.Errorf("invalid brick kind %s", value.Value)
	}
	return nil
}
