package adapters

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/seboste/sapper/ports"
)

type InteractiveParameterResolver struct {
}

func resolve(rd io.Reader, name string) string {
	reader := bufio.NewReader(rd)
	fmt.Printf("Enter value for parameter %s: ", name)
	value, _ := reader.ReadString('\n')
	return value
}

func (ipr InteractiveParameterResolver) Resolve(name string) string {
	return resolve(os.Stdin, name)
}

var _ ports.ParameterResolver = InteractiveParameterResolver{}
