package service

import (
	"io"
	"os/exec"

	"github.com/seboste/sapper/ports"
)

type MakeService struct {
}

func execMake(s ports.Service, target []string, output io.Writer) error {
	cmd := exec.Command("make", target[:]...)
	cmd.Dir = s.Path
	cmd.Stdout = output
	cmd.Stderr = output
	return cmd.Run()
}

func (ms MakeService) Build(s ports.Service, output io.Writer) error {
	return execMake(s, []string{"build", "-B"}, output)
}
func (ms MakeService) Test(s ports.Service, output io.Writer) error {
	return execMake(s, []string{"test"}, output)
}
func (ms MakeService) Run(s ports.Service, output io.Writer) error {
	return execMake(s, []string{"run"}, output)
}
func (ms MakeService) Deploy(s ports.Service, output io.Writer) error {
	return execMake(s, []string{"deploy"}, output)
}
func (ms MakeService) Stop(s ports.Service, output io.Writer) error {
	return execMake(s, []string{"stop"}, output)
}

var _ ports.ServiceBuilder = MakeService{}
