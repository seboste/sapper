package parameterResolver

import "github.com/seboste/sapper/ports"

type MapBasedParameterResolver struct {
	parameters map[string]string
}

func (r MapBasedParameterResolver) Resolve(key string, defaultValue string) string {
	return r.parameters[key]
}

var _ ports.ParameterResolver = MapBasedParameterResolver{}
