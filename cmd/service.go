package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/seboste/sapper/adapters"
	"github.com/spf13/cobra"
)

var serviceCmd = &cobra.Command{
	Use:   "service",
	Short: "Manage C++ microservices",
}

var addServiceCmd = &cobra.Command{
	Use:   "add [folder]",
	Short: "Adds a new C++ microservice",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("service folder argument is missing")
			return
		}

		path, name := filepath.Split(args[0])
		template, _ := cmd.Flags().GetString("template")

		r, err := adapters.MakeSapperParameterResolver(cmd.Flags(), name)
		if err != nil {
			fmt.Println(err)
			return
		}

		if err := serviceApi.Add(template, path, r); err != nil {
			fmt.Println(err)
		}
	},
}

var describeServiceCmd = &cobra.Command{
	Use:   "describe [service folder]",
	Short: "Prints information about a service",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("service folder argument is missing")
			return
		}
		description, err := serviceApi.Describe(args[0])

		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(description)
	},
}

var updateServiceCmd = &cobra.Command{
	Use:   "update [service folder]",
	Short: "Updates the dependencies of the service",
	Run: func(cmd *cobra.Command, args []string) {
		serviceApi.Update()
	},
}

var buildServiceCmd = &cobra.Command{
	Use:   "build [service folder]",
	Short: "Builds the service",
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) < 1 {
			fmt.Println("service folder argument is missing")
			return
		}
		serviceApi.Build(args[0])
	},
}

var testServiceCmd = &cobra.Command{
	Use:   "test [template]",
	Short: "Tests the service",
	Run: func(cmd *cobra.Command, args []string) {
		serviceApi.Test()
	},
}

var deployServiceCmd = &cobra.Command{
	Use:   "deploy [template]",
	Short: "Deploy the service",
	Run: func(cmd *cobra.Command, args []string) {
		serviceApi.Deploy()
	},
}

func init() {
	serviceCmd.AddCommand(addServiceCmd)
	serviceCmd.AddCommand(describeServiceCmd)
	serviceCmd.AddCommand(updateServiceCmd)
	serviceCmd.AddCommand(buildServiceCmd)
	serviceCmd.AddCommand(testServiceCmd)
	serviceCmd.AddCommand(deployServiceCmd)

	rootCmd.AddCommand(serviceCmd)

	addServiceCmd.PersistentFlags().StringP("template", "t", "base-hexagonal-skeleton", "The id of a service template.")
	adapters.RegisterSapperParameterResolver(addServiceCmd.PersistentFlags())

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// brickCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// brickCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
