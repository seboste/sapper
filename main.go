package main

import (
	"os"

	"github.com/seboste/sapper/adapters"
	"github.com/seboste/sapper/cmd"
	"github.com/seboste/sapper/core"
)

func main() {

	fsc, err := adapters.MakeFilesystemConfiguration()
	if err != nil {
		panic(err)
	}

	brickDB, err := adapters.MakeBrickDB(fsc.Remotes())
	if err != nil {
		panic(err)
	}

	dependencyManager := adapters.ConanDependencyManager{}
	servicePersistence := adapters.FileSystemServicePersistence{DependencyReader: dependencyManager}
	ServiceBuilder := adapters.CMakeService{}

	serviceApi := core.ServiceApi{
		Db:                 brickDB,
		ServicePersistence: servicePersistence,
		ServiceBuilder:     ServiceBuilder,
		DependencyInfo:     dependencyManager,
		DependencyWriter:   dependencyManager,
		Stdout:             os.Stdout,
		Stderr:             os.Stderr,
	}

	brickApi := core.BrickApi{Db: brickDB,
		PackageDependencyReader: dependencyManager,
		PackageDependencyWriter: dependencyManager,
		DependencyInfo:          dependencyManager,
		ServicePersistence:      servicePersistence,
		ServiceApi:              serviceApi,
	}
	cmd.SetApis(brickApi, serviceApi, core.RemoteApi{})
	cmd.Execute()
}
