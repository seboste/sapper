package adapters

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

var _ ports.ServiceBuilder = CMakeService{}
