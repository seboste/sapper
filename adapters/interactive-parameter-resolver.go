package adapters

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/seboste/sapper/ports"
)

type InteractiveParameterResolver struct {
	DefaultResolver *ports.ParameterResolver
}

func resolve(rd io.Reader, wr io.Writer, name string, defaultValue string) string {
	reader := bufio.NewReader(rd)
	value := ""
	if defaultValue != "" {
		fmt.Fprintf(wr, "Enter value for parameter %s or press enter for default %s: ", name, defaultValue)
		value, _ = reader.ReadString('\n')
		value = strings.TrimRight(value, "\n")
		if value == "" {
			value = defaultValue
		}
	} else {
		for value == "" {
			fmt.Fprintf(wr, "Enter value for parameter %s: ", name)
			value, _ = reader.ReadString('\n')
			value = strings.TrimRight(value, "\n")
		}
	}
	return value
}

func (ipr InteractiveParameterResolver) Resolve(name string) string {
	defaultValue := ""
	if ipr.DefaultResolver != nil {
		defaultValue = (*ipr.DefaultResolver).Resolve(name)
	}
	return resolve(os.Stdin, os.Stdout, name, defaultValue)
}

var _ ports.ParameterResolver = InteractiveParameterResolver{}
