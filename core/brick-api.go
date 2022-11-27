package core

import (
	"fmt"
	"strings"

	"github.com/seboste/sapper/ports"
)

type BrickApi struct {
	Db                 ports.BrickDB
	ServicePersistence ports.ServicePersistence
}

func removeBricks(bricks []ports.Brick, brickIdsToRemove []ports.BrickDependency) []ports.Brick {
	brickIdsToRemoveMap := make(map[string]bool)
	for _, b := range brickIdsToRemove {
		brickIdsToRemoveMap[b.Id] = true
	}

	filteredBricks := []ports.Brick{}
	for _, b := range bricks {
		if brickIdsToRemoveMap[b.GetId()] == false {
			filteredBricks = append(filteredBricks, b)
		}
	}
	return filteredBricks
}

func (b BrickApi) Add(servicePath string, brickId string, parameterResolver ports.ParameterResolver) error {

	bricks, err := GetBricksRecursive(brickId, b.Db)
	if err != nil {
		return err
	}

	service, err := b.ServicePersistence.Load(servicePath)
	if err != nil {
		return err
	}

	bricks = removeBricks(bricks, service.BrickIds) //remove all bricks that are already there

	if len(bricks) == 0 {
		return fmt.Errorf("brick %s has already been added.", brickId)
	}

	parameters, err := ResolveParameterSlice(bricks, parameterResolver)
	if err != nil {
		return err
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
