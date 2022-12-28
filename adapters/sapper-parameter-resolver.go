package adapters

import (
	"github.com/seboste/sapper/ports"
	"github.com/spf13/pflag"
)

type SapperParameterResolver struct {
	cpr CompoundParameterResolver
}

func RegisterSapperParameterResolver(flags *pflag.FlagSet) {
	flags.StringArrayP("parameter", "p", []string{}, "Sets parameters of the service (Example: '-p PARAM_NAME=value').")
}

func MakeSapperParameterResolver(flags *pflag.FlagSet, name string) (SapperParameterResolver, error) {

	resolver := []ports.ParameterResolver{}

	//1. read from command line parameters
	parameter, _ := flags.GetStringArray("parameter")
	clipr, err := MakeCommandLineInterfaceParameterResolver(parameter)
	if err != nil {
		return SapperParameterResolver{}, err
	}
	resolver = append(resolver, clipr)

	//2. set name as fixed parameter
	if name != "" {
		resolver = append(resolver, MapBasedParameterResolver{parameters: map[string]string{"NAME": name}})
	}

	//3. ask user for parameters if other methods failed
	resolver = append(resolver, InteractiveParameterResolver{})

	return SapperParameterResolver{cpr: CompoundParameterResolver{resolver: resolver}}, nil

}

func (r SapperParameterResolver) Resolve(key string, defaultValue string) string {
	return r.cpr.Resolve(key, defaultValue) //delegate
}

var _ ports.ParameterResolver = SapperParameterResolver{}
