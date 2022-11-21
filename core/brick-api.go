package core

import (
	"fmt"
	"strings"

	"github.com/seboste/sapper/ports"
)

type BrickApi struct {
	Db ports.BrickDB
}

func (b BrickApi) Add() {
	fmt.Println("add")
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
