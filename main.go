package main

import (
	"fmt"

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

	cmd.SetApis(
		core.BrickApi{Db: brickDb, ServicePersistence: servicePersistence},
		core.ServiceApi{
			Db:                 brickDb,
			ServicePersistence: servicePersistence,
			ServiceBuilder:     ServiceBuilder,
			DependencyInfo:     dependencyManager,
			DependencyWriter:   dependencyManager,
		},
		core.RemoteApi{},
	)
	cmd.Execute()
}
