package core

import (
	"fmt"

	"github.com/seboste/sapper/ports"
)

type Remote struct {
}

func (r Remote) Add() {
	fmt.Println("add")
}

func (r Remote) List() {
	fmt.Println("list")
}

var _ ports.RemoteApi = Remote{}
