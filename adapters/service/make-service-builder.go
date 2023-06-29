package service

import (
	"io"
	"os/exec"

	"github.com/seboste/sapper/ports"
)

type MakeService struct {
}

func execMake(s ports.Service, target string, output io.Writer) error {
	cmd := exec.Command("make", target)
	cmd.Dir = s.Path
	cmd.Stdout = output
	cmd.Stderr = output
	return cmd.Run()
}

func (ms MakeService) Build(s ports.Service, output io.Writer) error {
	return execMake(s, "build -B", output)
}
func (ms MakeService) Test(s ports.Service, output io.Writer) error {
	return execMake(s, "test", output)
}
func (ms MakeService) Run(s ports.Service, output io.Writer) error {
	return execMake(s, "run", output)
}
func (ms MakeService) Deploy(s ports.Service, output io.Writer) error {
	return execMake(s, "deploy", output)
}

var _ ports.ServiceBuilder = MakeService{}
