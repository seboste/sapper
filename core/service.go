package core

import (
	"fmt"

	"github.com/seboste/sapper/ports"
)

type Service struct {
}

func (s Service) Add() {
	fmt.Println("add")
}

func (s Service) Update() {
	fmt.Println("update")
}

func (s Service) Build() {
	fmt.Println("build")
}

func (s Service) Test() {
	fmt.Println("test")
}

func (s Service) Deploy() {
	fmt.Println("deploy")
}

var _ ports.ServiceApi = Service{}
