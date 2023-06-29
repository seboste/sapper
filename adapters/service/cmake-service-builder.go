package service

import (
	"io"
	"os/exec"

	"github.com/seboste/sapper/ports"
)

type CMakeService struct {
}

func (cms CMakeService) Build(s ports.Service, output io.Writer) error {
	cmd := exec.Command("make", "build", "-B")
	cmd.Dir = s.Path
	cmd.Stdout = output
	cmd.Stderr = output
	return cmd.Run()
}

func (cms CMakeService) Test(s ports.Service, output io.Writer) error {
	cmd := exec.Command("make", "test")
	cmd.Dir = s.Path
	cmd.Stdout = output
	cmd.Stderr = output
	return cmd.Run()
}

func (cms CMakeService) Run(s ports.Service, output io.Writer) error {
	cmd := exec.Command("make", "run")
	cmd.Dir = s.Path
	cmd.Stdout = output
	cmd.Stderr = output
	return cmd.Run()
}

func (cms CMakeService) Deploy(s ports.Service, output io.Writer) error {
	cmd := exec.Command("make", "deploy")
	cmd.Dir = s.Path
	cmd.Stdout = output
	cmd.Stderr = output
	return cmd.Run()
}

var _ ports.ServiceBuilder = CMakeService{}
