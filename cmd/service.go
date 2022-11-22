package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
)

var serviceCmd = &cobra.Command{
	Use:   "service",
	Short: "Manage C++ microservices",
}

type MapBasedParameterResolver struct {
	parameters map[string]string
}

func (r MapBasedParameterResolver) Resolve(key string) string {
	return r.parameters[key]
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

		r := MapBasedParameterResolver{parameters: map[string]string{"NAME": name}}
		if err := serviceApi.Add("base-hexagonal-http", path, r); err != nil {
			fmt.Println(err)
		}
	},
}

var updateServiceCmd = &cobra.Command{
	Use:   "update [template]",
	Short: "Updates the dependencies of the service",
	Run: func(cmd *cobra.Command, args []string) {
		serviceApi.Update()
	},
}

var buildServiceCmd = &cobra.Command{
	Use:   "build [template]",
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
	serviceCmd.AddCommand(updateServiceCmd)
	serviceCmd.AddCommand(buildServiceCmd)
	serviceCmd.AddCommand(testServiceCmd)
	serviceCmd.AddCommand(deployServiceCmd)

	rootCmd.AddCommand(serviceCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// brickCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// brickCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
