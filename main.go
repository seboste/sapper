package main

import (
	handler "github.com/seboste/sapper/adapters/cli_handler"
	"github.com/seboste/sapper/core"
)

func main() {
	c := core.Core{}
	h := handler.CliHandler{c}
	h.Handle()
}
