package core

import (
	"fmt"

	"github.com/seboste/sapper/ports"
)

type RemoteApi struct {
}

func (r RemoteApi) Add() {
	fmt.Println("add")
}

func (r RemoteApi) List() {
	fmt.Println("list")
}

var _ ports.RemoteApi = RemoteApi{}
