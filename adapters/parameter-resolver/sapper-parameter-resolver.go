package parameterResolver

import (
	"github.com/seboste/sapper/ports"
	upr "github.com/seboste/sapper/utils/parameter-resolver"
	"github.com/spf13/pflag"
)

type SapperParameterResolver struct {
	cpr upr.CompoundParameterResolver
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
		resolver = append(resolver, upr.MakeMapBasedParameterResolver(map[string]string{"NAME": name}))
	}

	//3. ask user for parameters if other methods failed
	resolver = append(resolver, InteractiveParameterResolver{})

	return SapperParameterResolver{cpr: upr.MakeCompoundParameterResolver(resolver)}, nil

}

func (r SapperParameterResolver) Resolve(key string, defaultValue string) string {
	return r.cpr.Resolve(key, defaultValue) //delegate
}

var _ ports.ParameterResolver = SapperParameterResolver{}
