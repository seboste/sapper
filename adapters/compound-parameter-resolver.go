package adapters

import (
	"github.com/seboste/sapper/ports"
)

type CompoundParameterResolver struct {
	resolver []ports.ParameterResolver
}

func MakeCompoundParameterResolver(resolver []ports.ParameterResolver) CompoundParameterResolver {
	return CompoundParameterResolver{resolver: resolver}
}

func (cpr CompoundParameterResolver) Resolve(name string, defaultValue string) string {
	for _, pr := range cpr.resolver {
		value := pr.Resolve(name, defaultValue)
		if value != "" {
			return value
		}
	}
	return ""
}
