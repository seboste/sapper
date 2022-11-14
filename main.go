package main

import (
	"fmt"

	"github.com/seboste/sapper/adapters"
	"github.com/seboste/sapper/cmd"
	"github.com/seboste/sapper/core"
)

func main() {
	brick_db := adapters.FilesystemBrickDB{}
	err := brick_db.Init("./remote")
	if err != nil {
		fmt.Println(err)
	}

	cmd.SetApis(core.Brick{}, core.Service{}, core.Remote{})
	cmd.Execute()
}
