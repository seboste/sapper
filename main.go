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

	servicePersistence := adapters.FileSystemServicePersistence{}

	cmd.SetApis(core.BrickApi{Db: brickDb}, core.ServiceApi{Db: brickDb, ServicePersistence: servicePersistence}, core.RemoteApi{})
	cmd.Execute()
}
