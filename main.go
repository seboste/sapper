package main

import (
	"fmt"
	"os"

	"github.com/seboste/sapper/adapters"
	"github.com/seboste/sapper/cmd"
	"github.com/seboste/sapper/core"
)

func main() {
	brickDb := &adapters.FilesystemBrickDB{}
	err := brickDb.Init("./remote")
	if err != nil {
		fmt.Println(err)
	}

	dependencyManager := adapters.ConanDependencyManager{}
	servicePersistence := adapters.FileSystemServicePersistence{DependencyReader: dependencyManager}
	ServiceBuilder := adapters.CMakeService{}

	serviceApi := core.ServiceApi{
		Db:                 brickDb,
		ServicePersistence: servicePersistence,
		ServiceBuilder:     ServiceBuilder,
		DependencyInfo:     dependencyManager,
		DependencyWriter:   dependencyManager,
		Stdout:             os.Stdout,
		Stderr:             os.Stderr,
	}

	brickApi := core.BrickApi{Db: brickDb,
		PackageDependencyReader: dependencyManager,
		DependencyInfo:          dependencyManager,
		ServicePersistence:      servicePersistence,
		ServiceApi:              serviceApi,
	}
	cmd.SetApis(brickApi, serviceApi, core.RemoteApi{})
	cmd.Execute()
}
