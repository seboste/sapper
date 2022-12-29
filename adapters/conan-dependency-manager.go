package adapters

import "github.com/seboste/sapper/ports"

type ConanDependencyManager struct {
}

func (cdm ConanDependencyManager) Read(s ports.Service) []ports.PackageDependency {
	return []ports.PackageDependency{}
}

func (cdm ConanDependencyManager) Write(s ports.Service) {
}

var _ ports.DependencyReader = ConanDependencyManager{}
var _ ports.DependencyWriter = ConanDependencyManager{}
