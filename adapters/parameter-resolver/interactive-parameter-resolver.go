package parameterResolver

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/seboste/sapper/ports"
)

type InteractiveParameterResolver struct {
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

var _ ports.ParameterResolver = InteractiveParameterResolver{}

func (ipr InteractiveParameterResolver) Resolve(name string, defaultValue string) string {
	return resolve(os.Stdin, os.Stdout, name, defaultValue)
}

var _ ports.ParameterResolver = InteractiveParameterResolver{}
