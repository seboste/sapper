package cmd

import (
	"github.com/spf13/cobra"
)

var serviceCmd = &cobra.Command{
	Use:   "service",
	Short: "Manage C++ microservices",
}

var addServiceCmd = &cobra.Command{
	Use:   "add [template]",
	Short: "Adds a new C++ microservice",
	Run: func(cmd *cobra.Command, args []string) {
		serviceApi.Add()
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
		serviceApi.Build()
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
