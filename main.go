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

	fsc, err := configuration.MakeFilesystemConfiguration()
	if err != nil {
		panic(err)
	}

	brickDB, err := brickDb.MakeBrickDB(fsc.Remotes(), fsc.DefaultRemotesDir())
	if err != nil {
		panic(err)
	}

	dependencyManager := dependencyManager.ConanDependencyManager{}
	servicePersistence := service.FileSystemServicePersistence{DependencyReader: dependencyManager}
	ServiceBuilder := service.CMakeService{}

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
