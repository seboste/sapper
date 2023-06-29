package main

import (
	"os"

	brickDb "github.com/seboste/sapper/adapters/brick-db"
	configuration "github.com/seboste/sapper/adapters/configuration"
	dependencyManager "github.com/seboste/sapper/adapters/dependency-manager"
	"github.com/seboste/sapper/adapters/service"
	"github.com/seboste/sapper/cmd"
	"github.com/seboste/sapper/core"
)

func main() {

	config, err := configuration.MakeFilesystemConfiguration()
	if err != nil {
		panic(err)
	}

	brickDbFactory := brickDb.Factory{}

	dependencyManager := dependencyManager.ConanDependencyManager{}
	servicePersistence := service.FileSystemServicePersistence{DependencyReader: dependencyManager}
	ServiceBuilder := service.MakeService{}

	serviceApi := core.ServiceApi{
		Configuration:      &config,
		BrickDBFactory:     brickDbFactory,
		ServicePersistence: servicePersistence,
		ServiceBuilder:     ServiceBuilder,
		DependencyInfo:     dependencyManager,
		DependencyWriter:   dependencyManager,
		Stdout:             os.Stdout,
		Stderr:             os.Stderr,
	}

	brickApi := core.BrickApi{
		Configuration:           &config,
		BrickDBFactory:          brickDbFactory,
		PackageDependencyReader: dependencyManager,
		PackageDependencyWriter: dependencyManager,
		DependencyInfo:          dependencyManager,
		ServicePersistence:      servicePersistence,
		ServiceApi:              serviceApi,
	}

	remoteApi := core.RemoteApi{
		Configuration:  &config,
		BrickDBFactory: brickDbFactory,
		BrickUpgrader:  brickApi,
	}

	cmd.SetApis(brickApi, serviceApi, remoteApi)
	cmd.SetVersion("0.2.0")
	cmd.Execute()
}
