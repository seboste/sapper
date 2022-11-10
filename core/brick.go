package core

import (
	"fmt"

	"github.com/seboste/sapper/ports"
)

type Brick struct {
}

func (b Brick) Add() {
	fmt.Println("add")
}

func (b Brick) List() {
	fmt.Println("list")
}

func (b Brick) Search() {
	fmt.Println("search")
}

var _ ports.BrickApi = Brick{}
