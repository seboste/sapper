package main

import (
	"os"

	"github.com/seboste/sapper/adapters"
	"github.com/seboste/sapper/cmd"
	"github.com/seboste/sapper/core"
)

func main() {
	// brickDb := &adapters.FilesystemBrickDB{}
	// err := brickDb.Init("./remote")
	// if err != nil {
	// 	panic(err)
	// }

	homedir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	fsc, err := adapters.MakeFilesystemConfiguration(homedir)
	if err != nil {
		panic(err)
	}

	brickDB := adapters.MakeAggregateBrickDB(fsc.Remotes())

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
