package utils

import "github.com/seboste/sapper/ports"

type DummyParameterResolver struct {
}

func (r DummyParameterResolver) Resolve(key string, defaultValue string) string {
	if defaultValue != "" {
		return defaultValue
	}
	return key
}

var _ ports.ParameterResolver = DummyParameterResolver{}
