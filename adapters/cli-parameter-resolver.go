package adapters

import (
	"fmt"
	"regexp"

	"github.com/seboste/sapper/ports"
)

type CommandLineInterfaceParameterResolver struct {
	parameters map[string]string
}

func MakeCommandLineInterfaceParameterResolver(parameters []string) (CommandLineInterfaceParameterResolver, error) {
	parmeterExp := regexp.MustCompile(`(.*)=(.*)`)

	resolver := CommandLineInterfaceParameterResolver{}
	resolver.parameters = map[string]string{}
	for _, p := range parameters {

		matches := parmeterExp.FindStringSubmatch(p)
		if len(matches) != 3 {
			return CommandLineInterfaceParameterResolver{}, fmt.Errorf("parameter %s must be of the form 'PARAMETER_NAME=value'", p)
		}
		resolver.parameters[matches[1]] = matches[2]
	}
	return resolver, nil
}

func (clipr CommandLineInterfaceParameterResolver) Resolve(name string) string {
	return clipr.parameters[name]
}

var _ ports.ParameterResolver = CommandLineInterfaceParameterResolver{}
