package handler

import (
	"github.com/seboste/sapper/ports"
	"github.com/spf13/cobra"
)

type CliHandler struct {
	Api ports.Api
}

func (h CliHandler) Handle() error {
	var rootCmd = &cobra.Command{
		Use:   "sapper [command]",
		Short: "A cli tool for the rapid development of C++ microservices",
	}

	var cmdNew = &cobra.Command{
		Use:   "new [template]",
		Short: "Creates a new C++ microservice",
		Run: func(cmd *cobra.Command, args []string) {
			h.Api.New()
		},
	}

	var cmdAdd = &cobra.Command{
		Use:   "add [template]",
		Short: "Adds an extension to the C++ microservice",
		Run: func(cmd *cobra.Command, args []string) {
			h.Api.Add()
		},
	}

	var cmdUpdate = &cobra.Command{
		Use:   "update",
		Short: "Updates the dependencies of the C++ microservice",
		Run: func(cmd *cobra.Command, args []string) {
			h.Api.Update()
		},
	}

	rootCmd.AddCommand(cmdNew, cmdAdd, cmdUpdate)
	return rootCmd.Execute()
}
