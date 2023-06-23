package ports

import (
	"io"
)

type BrickUpgrader interface {
	UpgradeInDB(brickId string, db BrickDB) error
}

type BrickApi interface {
	Add(servicePath string, brickId string, parameterResolver ParameterResolver) error
	Upgrade(brickId string) error
	Describe(brickId string, writer io.Writer) error
	List() []Brick
	Search(term string) []Brick
}
