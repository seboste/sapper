package core

import (
	"fmt"

	"github.com/seboste/sapper/ports"
)

type Core struct {
}

func (c Core) New() {
	fmt.Println("new")
}

func (c Core) Add() {
	fmt.Println("add")
}

func (c Core) Update() {
	fmt.Println("update")
}

var _ ports.Api = Core{}
