package core

import (
	"strings"

	"github.com/seboste/sapper/ports"
)

type BrickApi struct {
	Db                 ports.BrickDB
	ServicePersistence ports.ServicePersistence
}

func (b BrickApi) Add(servicePath string, brickId string, parameterResolver ports.ParameterResolver) error {

	bricks, err := GetBricksRecursive(brickId, b.Db)
	if err != nil {
		return nil
	}

	parameters, err := ResolveParameterSlice(bricks, parameterResolver)
	if err != nil {
		return err
	}

	service, err := b.ServicePersistence.Load(servicePath)
	if err != nil {
		return nil
	}

	for _, brick := range bricks {
		if err := AddBrick(&service, brick, parameters); err != nil {
			return err
		}
	}

	if err := b.ServicePersistence.Save(service); err != nil {
		return err
	}
	return nil
}

func (b BrickApi) List() []ports.Brick {
	return b.Db.Bricks(ports.Extension)
}

func (b BrickApi) Search(term string) []ports.Brick {
	filteredBricks := []ports.Brick{}
	for _, brick := range b.Db.Bricks(ports.Extension) {
		if strings.Contains(brick.GetId(), term) || strings.Contains(brick.GetDescription(), term) {
			filteredBricks = append(filteredBricks, brick)
		}
	}
	return filteredBricks
}

var _ ports.BrickApi = BrickApi{}
