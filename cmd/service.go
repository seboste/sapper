package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	parameterResolver "github.com/seboste/sapper/adapters/parameter-resolver"
	"github.com/spf13/cobra"
)

var serviceCmd = &cobra.Command{
	Use:   "service",
	Short: "Manage C++ microservices",
}

var keepMajorVersion *bool

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

		r, err := parameterResolver.MakeSapperParameterResolver(cmd.Flags(), name)
		if err != nil {
			fmt.Println(err)
			return
		}

		if _, err := serviceApi.Add(template, path, r); err != nil {
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
		err := serviceApi.Describe(args[0], os.Stdout)

		if err != nil {
			fmt.Println(err)
			return
		}
	},
}

var upgradeServiceCmd = &cobra.Command{
	Use:   "upgrade [service folder]",
	Short: "upgrades the dependencies of the service",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.PersistentFlags()
		err := serviceApi.Upgrade(args[0], *keepMajorVersion)
		if err != nil {
			fmt.Println(err)
			return
		}
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

		fmt.Printf("building service...")
		buildLogFilename, err := serviceApi.Build(args[0])
		if err != nil {
			fmt.Printf("failed (see %s for details)\n", buildLogFilename)
		} else {
			fmt.Println("success")
		}
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
	serviceCmd.AddCommand(upgradeServiceCmd)
	serviceCmd.AddCommand(buildServiceCmd)
	serviceCmd.AddCommand(testServiceCmd)
	serviceCmd.AddCommand(deployServiceCmd)

	rootCmd.AddCommand(serviceCmd)

	addServiceCmd.PersistentFlags().StringP("template", "t", "base-hexagonal-skeleton", "The id of a service template.")
	parameterResolver.RegisterSapperParameterResolver(addServiceCmd.PersistentFlags())

	keepMajorVersion = upgradeServiceCmd.PersistentFlags().Bool("keep-major", false, "Upgrades are only conducted within the same major version of a dependency's semantic version")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// brickCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// brickCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
